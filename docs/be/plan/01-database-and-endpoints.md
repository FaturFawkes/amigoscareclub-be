# Plan: Database Design + Migrations + Backend Endpoints (DDD / Clean Architecture)

## Context

Frontend untuk **40% OF HEART RATE RUN – VOL.2** sudah diimplementasi (`docs/fe/plan/01`, `02`) dan
membutuhkan backend sesuai kontrak OpenAPI di `docs/fe/plan/event-registration.yaml`. Backend Go
(`module myapp`, Go 1.25, Gin) sudah punya kerangka DDD/clean architecture
(`domain/`, `application/`, `infrastructure/`, `delivery/`), tetapi:

- Domain model existing salah konteks (running 5K/10K + kategori) — kontrak sebenarnya adalah
  pendaftaran event kopi: `coffee_choice`, `payment_proof`, status verifikasi, admin auth.
- Tidak ada koneksi DB, driver, migration, config, maupun implementasi repo (semua stub `not implemented`).
- Repo handler campur (`gin` di `main.go`, `http.ServeMux` di `router.go`).

**Tujuan:** desain skema database + migration, lalu implementasi **8 endpoint** swagger penuh dengan
DDD + clean architecture, sehingga FE bisa langsung disambungkan via `NEXT_PUBLIC_API_BASE_URL`.

**Keputusan (dikonfirmasi user):** PostgreSQL · golang-migrate · S3-compatible storage · JWT + blacklist DB.

## Endpoint yang harus dipenuhi (sumber: `event-registration.yaml`)

| Method | Path | Auth |
|--------|------|------|
| POST | `/events/{eventSlug}/registrations` (multipart, upload bukti) | publik |
| GET | `/events/{eventSlug}/registrations/{registrationId}` | publik |
| GET | `/events/{eventSlug}` | publik |
| POST | `/admin/auth/login` | publik |
| GET | `/admin/auth/me` | Bearer |
| POST | `/admin/auth/logout` | Bearer |
| GET | `/admin/events/{eventSlug}/registrations` (filter status, pagination) | Bearer |
| PATCH | `/admin/events/{eventSlug}/registrations/{registrationId}/verify` | Bearer |

Routes di-mount di group **`/v1`** (configurable, sesuai server URL swagger & base URL FE plan 02).

---

## Bagian 1 — Desain Database (PostgreSQL)

### Tabel

**`events`**
- `slug` TEXT PK
- `title` TEXT NOT NULL
- `date` DATE NOT NULL
- `time` TIME NOT NULL
- `timezone` TEXT NOT NULL DEFAULT `'Asia/Jakarta'`
- `location` TEXT NOT NULL
- `distance_km` INT NOT NULL
- `pace` TEXT NOT NULL
- `registration_open` BOOLEAN NOT NULL DEFAULT true
- `coffee_options` JSONB NOT NULL  (array string enum kopi)
- `payment_bank` TEXT NOT NULL, `payment_account_number` TEXT NOT NULL, `payment_account_name` TEXT NOT NULL
- `created_at`, `updated_at` TIMESTAMPTZ NOT NULL DEFAULT now()

**`registrations`**
- `id` TEXT PK  (format `reg_<ULID>`)
- `ticket_number` TEXT UNIQUE NOT NULL  (format `40% OHHR-VOL.2-NNNN`, dari sequence)
- `event_slug` TEXT NOT NULL REFERENCES `events(slug)`
- `name` TEXT NOT NULL (≤100), `email` TEXT NOT NULL, `phone` TEXT NOT NULL
- `age` INT NOT NULL CHECK (age BETWEEN 10 AND 100)
- `coffee_choice` TEXT NOT NULL
- `status` registration_status NOT NULL DEFAULT `'pending_verification'`
- `payment_proof_url` TEXT  (nullable; object key/URL S3)
- `note` TEXT  (catatan admin saat reject)
- `registered_at` TIMESTAMPTZ NOT NULL DEFAULT now()
- `verified_at` TIMESTAMPTZ NULL, `ticket_sent_at` TIMESTAMPTZ NULL
- `created_at`, `updated_at` TIMESTAMPTZ NOT NULL DEFAULT now()
- **UNIQUE (`event_slug`, `email`)** → idempotency 1 email / event (→ 409 DUPLICATE_REGISTRATION)
- ENUM `registration_status`: `pending_verification | verified | rejected | ticket_sent`
- SEQUENCE `ticket_number_seq` untuk generate nomor tiket (prefix configurable)

**`admins`** (login pakai **email + password** sesuai swagger `LoginRequest`)
- `id` TEXT PK (`adm_<ULID>`)
- `name` TEXT NOT NULL  (untuk `AdminProfile.name`)
- `email` TEXT UNIQUE NOT NULL  (identitas login / "username" — field `email` di `LoginRequest`)
- `password_hash` TEXT NOT NULL  (bcrypt; password plain min 8 char sesuai `LoginRequest.password.minLength`)
- `created_at`, `updated_at` TIMESTAMPTZ NOT NULL DEFAULT now()

> Swagger `LoginRequest` = `{ email, password }`, `AdminProfile` = `{ id, name, email }`.
> Tidak ada field `username` terpisah — `email` berperan sebagai username login.

**`revoked_tokens`** (JWT blacklist untuk logout)
- `jti` TEXT PK, `expires_at` TIMESTAMPTZ NOT NULL, `revoked_at` TIMESTAMPTZ NOT NULL DEFAULT now()
- Index `expires_at` (cleanup token kadaluarsa)

### Migration files (golang-migrate, `migrations/`)

```
000001_create_registration_status_enum.up.sql / .down.sql
000002_create_events.up.sql / .down.sql
000003_create_ticket_number_seq.up.sql / .down.sql
000004_create_registrations.up.sql / .down.sql
000005_create_admins.up.sql / .down.sql
000006_create_revoked_tokens.up.sql / .down.sql
000007_seed_event.up.sql / .down.sql        # seed event 40-of-heart-rate-run + coffee_options
```

Seed admin via CLI Go (`cmd/seedadmin`) karena butuh bcrypt hash (bukan SQL polos).

---

## Bagian 2 — Implementasi Endpoint (DDD + Clean Architecture)

Pertahankan layout existing. Refactor `domain` agar selaras kontrak; folder repo
`infrastructure/repository/mysql` → **`infrastructure/repository/postgres`**.

### domain/ (entitas + aturan bisnis, tanpa dependency luar)
- `registration.go` — rework: aggregate `Registration` (Runner VO: name/email/phone/age/coffee),
  `RegistrationStatus` VO + transisi tervalidasi: `Verify()`, `Reject(note)`, `MarkTicketSent()`
  (tolak transisi ilegal → `ErrInvalidStatusTransition`).
- `event.go` — aggregate `Event` (slug, payment info, coffee_options, `registration_open`).
- `admin.go` — entitas `Admin`.
- `registration_repository.go` — perluas: `Save`, `GetByID`, `FindByEventAndEmail`,
  `List(ctx, eventSlug, filter, page, perPage) ([]*Registration, total)`, `Update`.
- `event_repository.go` — `GetBySlug`.
- `admin_repository.go` — `GetByEmail`, `GetByID`.
- `token_repository.go` — `Revoke(jti, exp)`, `IsRevoked(jti)`.
- `errors.go` — sentinel: `ErrEventNotFound`, `ErrRegistrationNotFound`, `ErrDuplicateRegistration`,
  `ErrInvalidStatusTransition`, `ErrInvalidCredentials`, `ErrUnauthorized`.
- `events.go` — pertahankan domain event `TicketRegistered` (dipakai notifier).

### application/
- `serviceInterface/` — pertahankan `storage.go` (FileStorage) & `notifier.go`; tambah
  `token.go` (TokenService: Issue/Parse JWT+jti), `security.go` (PasswordHasher bcrypt),
  `clock.go`, `idgen.go` (ULID + ticket number).
- `usecase/`:
  - `create_registration.go` — validasi domain, cek duplicate, generate id+ticket,
    upload bukti via FileStorage, simpan, kirim konfirmasi (notifier opsional, non-fatal).
  - `get_registration.go` — by event+id.
  - `get_event.go` — by slug.
  - `admin_login.go` — input `email` + `password` (swagger `LoginRequest`), cari admin by email,
    verifikasi bcrypt, issue JWT → `LoginResponse { token, expires_at, admin: AdminProfile }`;
    gagal → `ErrInvalidCredentials` (401 `INVALID_CREDENTIALS`).
  - `admin_me.go` — profil dari token.
  - `admin_logout.go` — revoke jti.
  - `admin_list_registrations.go` — filter status + pagination → data + PaginationMeta.
  - `admin_verify_registration.go` — transisi status (verified/rejected+note), set timestamp.
- `dto/` — ganti DTO existing dengan request/response sesuai schema swagger
  (`RegistrationRequest/Response`, `EventResponse`, `LoginRequest/Response`, `AdminProfile`,
  `PaginationMeta`, envelope error).

### infrastructure/
- `config/config.go` — loader env (DB DSN, JWT secret/TTL, S3 endpoint/bucket/keys, port, base path);
  `joho/godotenv` untuk `.env` dev.
- `db/postgres.go` — `pgxpool.Pool` + ping + pool config.
- `repository/postgres/` — `registration_repo.go`, `event_repo.go`, `admin_repo.go`,
  `revoked_token_repo.go` (SQL nyata, map ⇄ domain, terjemahkan unique violation → `ErrDuplicateRegistration`).
- `auth/jwt.go` — TokenService (`golang-jwt/jwt/v5`, HS256, klaim sub+jti+exp).
- `auth/bcrypt.go` — PasswordHasher (`golang.org/x/crypto/bcrypt`).
- `idgen/` — ULID (`oklog/ulid/v2`) + ticket number via sequence.
- `s3/storage.go` — implement FileStorage S3-compatible (`aws-sdk-go-v2` custom endpoint → AWS S3 / MinIO),
  validasi MIME (jpeg/png/webp) & ukuran ≤5MB (→ 413).

### delivery/
- `http/router.go` — **standarisasi ke Gin** (buang `http.ServeMux`); group `/v1`, daftarkan 8 route,
  pasang middleware auth pada grup admin (kecuali login).
- `http/response.go` — helper envelope sukses `{data, meta}` & error `{error:{code,message,details}}`
  + mapping domain error → HTTP status/kode (VALIDATION_ERROR 400, DUPLICATE_REGISTRATION 409,
  EVENT_NOT_FOUND/REGISTRATION 404, 413, INVALID_CREDENTIALS/UNAUTHORIZED 401).
- `http/registration_handler.go` — rework: parse `multipart/form-data`, validasi field
  (`go-playground/validator` sudah ada), panggil usecase.
- `http/event_handler.go`, `http/admin_auth_handler.go`, `http/admin_registration_handler.go` — baru.
- `middleware/auth.go` — implement: ekstrak Bearer, parse JWT, cek `revoked_tokens`,
  inject admin ke context; 401 bila gagal.
- `middleware/logger.go` — pasang di router.

### cmd/
- `api/main.go` — rework wiring: load config → pgxpool → repos → services (jwt/bcrypt/s3/idgen)
  → usecases → handlers → Gin router (`/v1` group + middleware) → graceful shutdown.
- `seedadmin/main.go` — buat admin awal (email+password → bcrypt → insert).

### Dependencies baru (`go.mod`)
`github.com/jackc/pgx/v5` (+ pgxpool) · `github.com/golang-migrate/migrate/v4` ·
`github.com/golang-jwt/jwt/v5` · `github.com/oklog/ulid/v2` ·
`github.com/aws/aws-sdk-go-v2` (config, service/s3) · `github.com/joho/godotenv` ·
`golang.org/x/crypto/bcrypt` (sudah indirect). `mongodb.org/mongo-driver` tidak dipakai.

### Infra pendukung
- `.env.example` — `DATABASE_URL`, `JWT_SECRET`, `JWT_TTL`, `S3_ENDPOINT/BUCKET/KEY/SECRET/REGION`,
  `PORT`, `API_BASE_PATH=/v1`.
- `docker-compose.yml` — service `postgres`, `minio`, `migrate` (golang-migrate), `app`.
- `Makefile` — target: `run`, `build`, `test`, `migrate-up/down/create`, `seed`, `compose-up`,
  `docker-build`.
- `deploy/Dockerfile` — perbaiki `golang:1.22-alpine` → `golang:1.25-alpine` (samakan dgn go.mod).

---

## Task Detail (Breakdown Eksekusi)

Urut sesuai dependensi. Tiap task punya **DoD** (Definition of Done).

### Fase A — Fondasi (config, DB, infra dev)

- **T1. Dependencies** — `go get` pgx/v5, golang-migrate/v4, golang-jwt/v5, oklog/ulid/v2,
  aws-sdk-go-v2 (config + service/s3), joho/godotenv; `go mod tidy`. Hapus mongo-driver tak terpakai.
  *DoD:* `go mod tidy` bersih, `go build ./...` masih jalan.

- **T2. `infrastructure/config/config.go`** — struct `Config` + `Load()`: `DATABASE_URL`,
  `JWT_SECRET`, `JWT_TTL`, `S3_ENDPOINT/BUCKET/ACCESS_KEY/SECRET_KEY/REGION`, `PORT`,
  `API_BASE_PATH` (default `/v1`), `TICKET_PREFIX` (default `40% OHHR-VOL.2`). godotenv load `.env` jika ada.
  *DoD:* unit test `Load()` dgn env minimal lulus; error jelas bila field wajib kosong.

- **T3. `.env.example`** — semua key dari T2 dengan nilai contoh dev (Postgres & MinIO lokal).

- **T4. `infrastructure/db/postgres.go`** — `NewPool(ctx, cfg) (*pgxpool.Pool, error)` + `Ping`.
  *DoD:* konek ke Postgres compose berhasil; gagal koneksi → error eksplisit.

### Fase B — Migration & infra lokal

- **T5. Migration files** `migrations/000001..000007` (up+down) sesuai skema di Bagian 1
  (enum status, events, sequence tiket, registrations + constraints, admins, revoked_tokens, seed event).
  *DoD:* `migrate up` lalu `migrate down` bersih pada DB kosong.

- **T6. `docker-compose.yml`** — service `postgres`, `minio` (+ bucket init), `migrate`, `app`.
  *DoD:* `docker compose up -d postgres minio` sehat; service `migrate` menjalankan migration T5.

- **T7. `Makefile`** — target: `run`, `build`, `test`, `migrate-up`, `migrate-down`, `migrate-create`,
  `seed`, `compose-up`, `docker-build`.
  *DoD:* tiap target jalan tanpa error path yang hilang.

- **T8. `deploy/Dockerfile`** — perbaiki `golang:1.22-alpine` → `golang:1.25-alpine`.
  *DoD:* `make docker-build` sukses.

### Fase C — Domain (tanpa dependency luar)

- **T9. `domain/errors.go`** — semua sentinel error (lihat daftar di Bagian 2).

- **T10. `domain/registration.go`** — rework aggregate `Registration` + VO `Runner`
  (name/email/phone/age/coffee_choice) + `RegistrationStatus` + transisi tervalidasi:
  `Verify()`, `Reject(note string)`, `MarkTicketSent()` (ilegal → `ErrInvalidStatusTransition`).
  *DoD:* unit test transisi legal & ilegal lulus.

- **T11. `domain/event.go`** — aggregate `Event` (slug, title, date, time, timezone, location,
  distance_km, pace, registration_open, coffee_options, payment info).

- **T12. `domain/admin.go`** — entitas `Admin` (id, name, email, passwordHash).

- **T13. Repo interfaces** — update `registration_repository.go` (Save/GetByID/FindByEventAndEmail/List/Update),
  baru `event_repository.go` (GetBySlug), `admin_repository.go` (GetByEmail/GetByID),
  `token_repository.go` (Revoke/IsRevoked). Pertahankan `domain/events.go`.
  *DoD:* `go build ./domain/...` lulus.

### Fase D — Application: service interfaces + usecase

- **T14. `application/serviceInterface/`** — tambah `token.go` (TokenService: Issue/Parse JWT+jti),
  `security.go` (PasswordHasher bcrypt), `clock.go`, `idgen.go` (ULID + ticket number).
  Pertahankan `storage.go` & `notifier.go`.

- **T15. `application/dto/`** — ganti DTO lama: `RegistrationRequest/Response`, `EventResponse`,
  `LoginRequest/Response`, `AdminProfile`, `PaginationMeta`, envelope error — persis schema swagger.

- **T16. Usecase publik** — `create_registration.go` (validasi, cek duplicate, generate id+tiket,
  upload via FileStorage, save, notifier non-fatal), `get_registration.go`, `get_event.go`.
  *DoD:* unit test create sukses & duplicate → `ErrDuplicateRegistration` pakai repo mock.

- **T17. Usecase admin** — `admin_login.go` (email+password → JWT via swagger `LoginRequest`),
  `admin_me.go`, `admin_logout.go` (revoke jti), `admin_list_registrations.go` (filter+pagination),
  `admin_verify_registration.go` (transisi status + set timestamp).
  *DoD:* unit test login sukses/gagal & verify transisi status.

### Fase E — Infrastructure implementasi

- **T18. Rename** `infrastructure/repository/mysql` → `infrastructure/repository/postgres`
  (update package name & import di `cmd/api/main.go`).

- **T19. Repo Postgres** — `registration_repo.go`, `event_repo.go`, `admin_repo.go`,
  `revoked_token_repo.go`: SQL nyata + mapping ⇄ domain; pgconn error 23505 (unique violation) →
  `ErrDuplicateRegistration`; not found → sentinel error.
  *DoD:* integration test via Postgres dari compose.

- **T20. `infrastructure/auth/jwt.go` + `bcrypt.go`** — TokenService HS256 (sub+jti+exp), bcrypt hasher.
  *DoD:* unit test issue→parse roundtrip & hash/compare lulus.

- **T21. `infrastructure/idgen/`** — ULID generator dengan prefix (`reg_`/`adm_`) + ticket number
  dari sequence PostgreSQL + `TICKET_PREFIX` dari config.

- **T22. `infrastructure/s3/storage.go`** — implement FileStorage S3-compatible (`aws-sdk-go-v2`,
  endpoint custom untuk MinIO/S3), validasi MIME jpeg/png/webp & size ≤5MB → return URL/key.
  *DoD:* upload ke MinIO compose berhasil; file >5MB/non-image → error.

### Fase F — Delivery (HTTP)

- **T23. `delivery/http/response.go`** — helper envelope sukses `{data, meta}` + error
  `{error:{code,message,details}}` + mapping domain error → HTTP status+kode
  (400 VALIDATION_ERROR, 409 DUPLICATE_REGISTRATION, 404 *_NOT_FOUND, 413, 401 INVALID_CREDENTIALS/UNAUTHORIZED).

- **T24. `delivery/middleware/auth.go`** — ekstrak Bearer, parse JWT, cek `revoked_tokens`,
  inject admin ke Gin context; gagal → 401. Aktifkan `middleware/logger.go` di router.

- **T25. Handlers** — rework `registration_handler.go` (multipart/form-data + validator), buat baru:
  `event_handler.go`, `admin_auth_handler.go`, `admin_registration_handler.go`.

- **T26. `delivery/http/router.go`** — standarisasi **Gin penuh** (buang `http.ServeMux`), group `/v1`,
  daftarkan 8 route, middleware auth pada grup admin (kecuali login).
  *DoD:* `go build ./...` lulus, semua route terdaftar dan terverifikasi.

### Fase G — Wiring & finalisasi

- **T27. `cmd/api/main.go`** — rework wiring lengkap: config → pgxpool → repos → services
  (jwt/bcrypt/s3/idgen) → usecases → handlers → Gin `/v1` + graceful shutdown.
  *DoD:* `make run` start tanpa panic, health check `/v1/events/...` merespons.

- **T28. `cmd/seedadmin/main.go`** — CLI buat admin awal: email + password → bcrypt → INSERT admins,
  dijalankan via `make seed`.

- **T29. Smoke test end-to-end** — jalankan seluruh skenario di bagian Verifikasi.

## Verifikasi (end-to-end)

1. `docker-compose up -d postgres minio` → `make migrate-up` → `make seed` (admin) sukses.
2. `go build ./...` & `go vet ./...` lolos.
3. `make run`, lalu uji tiap endpoint vs swagger (curl/HTTP client):
   - `GET /v1/events/40-of-heart-rate-run` → detail event sesuai seed.
   - `POST /v1/events/40-of-heart-rate-run/registrations` multipart (field + file gambar)
     → 201 `{data:{id,ticket_number,status:pending_verification}, meta.message}`.
   - Submit ulang email sama → 409 `DUPLICATE_REGISTRATION`.
   - Field invalid / file >5MB / non-image → 400 / 413 dengan envelope error.
   - `GET /v1/events/{slug}/registrations/{id}` → detail registrasi.
   - `POST /v1/admin/auth/login` (admin seed) → token; salah → 401 `INVALID_CREDENTIALS`.
   - `GET /v1/admin/auth/me` Bearer → profil; tanpa/expired → 401.
   - `GET /v1/admin/events/{slug}/registrations?status=&page=&per_page=` → data + PaginationMeta.
   - `PATCH .../{id}/verify` `{status:"verified"}` → status berubah, `verified_at` terisi;
     `{status:"rejected", note}` → status rejected + note; transisi ilegal → 400.
   - `POST /v1/admin/auth/logout` → 204; reuse token sama → 401 (blacklist).
4. Validasi spec tetap konsisten: `npx @apidevtools/swagger-cli validate docs/fe/plan/event-registration.yaml`.
5. Sambungkan FE: set `NEXT_PUBLIC_API_BASE_URL` ke `http://localhost:8080/v1`, jalankan alur
   daftar peserta + login/dashboard admin dari FE plan 01 & 02.
