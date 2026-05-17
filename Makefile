APP_NAME     = amigoscareclub-be
DOCKER_IMAGE ?= $(APP_NAME):latest
DB_URL       ?= $(DATABASE_URL)

.PHONY: build run test migrate-up migrate-down migrate-create seed compose-up docker-build

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/$(APP_NAME) ./cmd/api

run:
	go run ./cmd/api

test:
	go test ./...

migrate-up:
	@export PATH=$$(go env GOPATH)/bin:$$PATH && \
	export DB_URL=$$(grep '^DATABASE_URL=' .env | cut -d '=' -f2-) && \
	migrate -path migrations -database "$$DB_URL" up

migrate-down:
	@export PATH=$$(go env GOPATH)/bin:$$PATH && \
	export DB_URL=$$(grep '^DATABASE_URL=' .env | cut -d '=' -f2-) && \
	migrate -path migrations -database "$$DB_URL" down

migrate-create:
	@read -p "Migration name: " name; \
	migrate create -ext sql -dir migrations -seq $$name

seed:
	go run ./cmd/seedadmin

compose-up:
	docker compose up -d

docker-build:
	docker build -f deploy/Dockerfile -t $(DOCKER_IMAGE) .
