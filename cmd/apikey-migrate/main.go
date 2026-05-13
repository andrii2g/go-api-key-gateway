package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/andrii2g/go-api-key-gateway/migrations"
	"github.com/andrii2g/go-api-key-gateway/storage/sqlstore"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("usage: apikey-migrate [up|status]")
	}
	dsn := os.Getenv("APIKEY_DB_DSN")
	if dsn == "" {
		log.Fatal("APIKEY_DB_DSN is required")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	db, err := sqlstore.Open(ctx, sqlstore.Options{DSN: dsn})
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	runner := migrations.NewRunner(db, "")
	switch os.Args[1] {
	case "up":
		if err := runner.Up(ctx); err != nil {
			log.Fatal(err)
		}
	case "status":
		rows, err := runner.Status(ctx)
		if err != nil {
			log.Fatal(err)
		}
		for _, row := range rows {
			state := "pending"
			if row.Applied {
				state = "applied"
			}
			fmt.Printf("%s\t%s\n", row.Version, state)
		}
	default:
		log.Fatalf("unknown command %q", os.Args[1])
	}
}
