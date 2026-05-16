# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/amigoscareclub-be ./cmd

# Runtime stage
FROM alpine:3.20

RUN adduser -D -g '' appuser
WORKDIR /app

COPY --from=builder /out/amigoscareclub-be /app/amigoscareclub-be

USER appuser
EXPOSE 8080

ENTRYPOINT ["/app/amigoscareclub-be"]
