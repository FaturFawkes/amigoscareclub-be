package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/oklog/ulid/v2"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	_ = godotenv.Load()

	dbURL := os.Getenv("DATABASE_URL")
	email := os.Getenv("ADMIN_EMAIL")
	password := os.Getenv("ADMIN_PASSWORD")
	name := os.Getenv("ADMIN_NAME")

	if dbURL == "" || email == "" || password == "" || name == "" {
		log.Fatal("DATABASE_URL, ADMIN_EMAIL, ADMIN_PASSWORD, ADMIN_NAME are required")
	}

	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("connect db: %v", err)
	}
	defer pool.Close()

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("bcrypt: %v", err)
	}

	id := "adm_" + ulid.Make().String()
	_, err = pool.Exec(context.Background(), `
		INSERT INTO admins (id, name, email, password_hash, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		ON CONFLICT (email) DO UPDATE SET name=$2, password_hash=$4, updated_at=NOW()`,
		id, name, email, string(hash))
	if err != nil {
		log.Fatalf("insert admin: %v", err)
	}

	fmt.Printf("Admin seeded: %s (%s)\n", name, email)
}
