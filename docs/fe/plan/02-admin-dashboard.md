# Plan 02: Hidden Admin Login + Dashboard Event Registration

**Status:** ✅ Implemented

## Context

Halaman registrasi peserta (`/event/40-of-heart-rate-run`) sudah selesai dibuat, dan OpenAPI spec di `docs/swagger/event-registration.yaml` sudah mendefinisikan endpoint admin (`adminListRegistrations`, `adminVerifyRegistration`) dengan skema `BearerAuth` (JWT).

Kebutuhan baru: admin perlu **halaman tersembunyi** (tidak ada link/nav menuju sana) untuk:
1. **Login admin** — autentikasi memanggil backend eksternal sesuai swagger (mendapat JWT).
2. **Dashboard event registration** — menampilkan daftar peserta yang mendaftar.
3. **Konfirmasi manual pembayaran** — admin menandai sebuah registrasi sebagai `verified` (atau `rejected`).

## Decisions

| Keputusan | Pilihan |
|-----------|---------|
| Auth | Panggil backend eksternal (Bearer JWT sesuai swagger) |
| Data source | Backend eksternal — UI dibangun sesuai kontrak API |
| Hidden route | Path rahasia + tidak ada link + guard redirect + `noindex` |
| Token storage | `sessionStorage` (terhapus saat tab ditutup) |
| Route path | `/admin-acr-2026` (login) · `/admin-acr-2026/dashboard` |

**Catatan:** Backend belum tersedia saat implementasi — UI sudah siap, sambungkan dengan set `NEXT_PUBLIC_API_BASE_URL` ke backend nyata.

## Implementation

### 1. Konfigurasi API base URL
- `.env.example` — dokumentasi `NEXT_PUBLIC_API_BASE_URL=https://api.amigoscareclub.id/v1`
- Fallback default ke URL production bila env kosong.

### 2. API client + auth helpers
- **`lib/adminAuth.ts`** — `getToken()`, `setToken()`, `clearToken()`, hook `useAdminGuard()` (redirect ke login jika token tidak ada).
- **`lib/adminApi.ts`** — satu fetch wrapper dengan Bearer token injection + 401 handler (clear token → redirect login):
  - `adminLogin(email, password)` → `POST /admin/auth/login`
  - `getAdminMe()` → `GET /admin/auth/me`
  - `adminLogout()` → `POST /admin/auth/logout`
  - `listRegistrations(slug, params)` → `GET /admin/events/{slug}/registrations`
  - `verifyRegistration(slug, id, status, note?)` → `PATCH /admin/events/{slug}/registrations/{id}/verify`

### 3. Hidden route `app/admin-acr-2026/`
- **`layout.tsx`** — `metadata.robots = { index: false, follow: false }` → seluruh segment `noindex, nofollow`.
- **`page.tsx`** (Login) — form Email + Password; sukses → `setToken` → `router.replace("/admin-acr-2026/dashboard")`; gagal → tampilkan pesan error.
- **`dashboard/page.tsx`** (Dashboard):
  - `useAdminGuard()` — redirect ke login jika tidak ada token.
  - `useEffect` dengan pattern `run()` inline (cancellable) + `fetchData` callback untuk imperative refresh.
  - Tabel peserta: No. Tiket, Nama, Email, No. HP, Usia, Kopi, Bukti (link), Status (badge berwarna), Tanggal daftar, Aksi.
  - Filter status dropdown + pagination (page/per_page).
  - Aksi baris `pending_verification`: "Verifikasi" → `verifyRegistration(..., "verified")`; "Tolak" → modal input catatan → `verifyRegistration(..., "rejected", note)`.
  - Tombol Logout → `adminLogout()` + `clearToken()` + redirect.

### 4. Update OpenAPI spec
Tambahan ke `docs/swagger/event-registration.yaml`:
- Tag **`Admin Auth`**.
- `POST /admin/auth/login` — `LoginRequest { email, password }` → `LoginResponse { data: { token, expires_at, admin } }`.
- `GET /admin/auth/me` — validasi token, return `AdminProfile`.
- `POST /admin/auth/logout` — invalidasi token server-side, `204`.
- Schema baru: `LoginRequest`, `LoginResponse`, `AdminProfile`.
- Field `payment_proof_url` (nullable string) ditambah ke schema `Registration`.

## Files Modified

| File | Aksi |
|------|------|
| `.env.example` | **Baru** |
| `lib/adminAuth.ts` | **Baru** |
| `lib/adminApi.ts` | **Baru** |
| `app/admin-acr-2026/layout.tsx` | **Baru** |
| `app/admin-acr-2026/page.tsx` | **Baru** — halaman login |
| `app/admin-acr-2026/dashboard/page.tsx` | **Baru** — dashboard |
| `docs/swagger/event-registration.yaml` | **Edit** — auth endpoints + schema baru |

## Verification

1. `npm run lint` & `npm run build` → lolos.
2. `npx @apidevtools/swagger-cli validate docs/swagger/event-registration.yaml` → "is valid".
3. `npm run dev`:
   - `/admin-acr-2026` via URL langsung → muncul form login (tidak ada link ke sini dari mana pun).
   - `/admin-acr-2026/dashboard` tanpa login → redirect ke login.
   - `<head>` halaman admin mengandung `noindex, nofollow`.
4. **End-to-end**: belum bisa ditest sampai backend tersedia. Set `NEXT_PUBLIC_API_BASE_URL` ke backend nyata untuk pengujian penuh.
