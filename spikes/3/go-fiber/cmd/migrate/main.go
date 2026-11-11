package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://area:area@localhost:5432/area?sslmode=disable"
	}
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		panic(fmt.Errorf("open: %w", err))
	}
	defer db.Close()

	if err := goose.SetDialect("postgres"); err != nil {
		panic(err)
	}
	cmd := "up"
	if len(os.Args) > 1 {
		cmd = os.Args[1]
	}
	ctx := context.Background()
	if err := goose.RunContext(ctx, cmd, db, "migrations"); err != nil {
		panic(fmt.Errorf("goose %s: %w", cmd, err))
	}
}
