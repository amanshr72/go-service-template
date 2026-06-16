package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

func main() {
	if len(os.Args) < 3 {
		log.Fatal("usage: migrate <ddl|dml> <up|down|status>")
	}

	migType := os.Args[1] // ddl | dml
	action := os.Args[2]  // up | down | status
	_ = godotenv.Load()
	env := os.Getenv("APP_ENV")
	dbString := os.Getenv("GOOSE_DBSTRING")

	if dbString == "" {
		log.Fatal("GOOSE_DBSTRING not set")
	}

	if migType != "ddl" && migType != "dml" {
		log.Fatalf("invalid type %q: must be ddl or dml", migType)
	}

	if migType == "ddl" && env == "prod" {
		log.Fatal("DDL blocked in prod. Use elevated manual credentials directly.")
	}

	db, err := sql.Open("postgres", dbString)
	if err != nil {
		log.Fatalf("db open: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("close db: %v", err)
		}
	}()

	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatal(err)
	}

	dir := fmt.Sprintf("migrations/%s", migType)

	switch action {
	case "up":
		err = goose.Up(db, dir)
	case "down":
		err = goose.Down(db, dir)
	case "status":
		err = goose.Status(db, dir)
	default:
		log.Fatalf("invalid action %q: must be up, down, or status", action)
	}
	if err != nil {
		log.Fatalf("migration failed: %v", err)
	}
	log.Printf("%s %s completed successfully", migType, action)
}
