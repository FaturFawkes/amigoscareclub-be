APP_NAME=amigoscareclub-be
DOCKER_IMAGE?=$(APP_NAME):latest

.PHONY: build docker-build

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/$(APP_NAME) ./cmd

docker-build:
	docker build -f deploy/Dockerfile -t $(DOCKER_IMAGE) .
