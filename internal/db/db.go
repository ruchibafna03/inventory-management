package db

import (
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// Connect establishes a PostgreSQL connection using DATABASE_URL env var.
func Connect() *sqlx.DB {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://postgres:postgres@localhost:5432/val_inventory?sslmode=disable"
	}

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatalf("db: failed to connect: %v", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	fmt.Println("db: connected to PostgreSQL")
	// Unsafe allows SELECT * queries where the DB has more columns than the struct.
	return db.Unsafe()
}
