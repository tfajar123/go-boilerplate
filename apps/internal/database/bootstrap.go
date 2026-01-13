package database

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"

	_ "github.com/lib/pq"
)

func EnsureDatabaseExists(databaseURL string) {
	u, err := url.Parse(databaseURL)
	if err != nil {
		log.Fatalf("invalid database url: %v", err)
	}

	dbName := u.Path[1:] // remove leading "/"
	u.Path = "/postgres" // connect ke default db

	db, err := sql.Open("postgres", u.String())
	if err != nil {
		log.Fatalf("failed connect postgres: %v", err)
	}
	defer db.Close()

	var exists bool
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM pg_database
			WHERE datname = $1
		)
	`

	if err := db.QueryRow(query, dbName).Scan(&exists); err != nil {
		log.Fatal(err)
	}

	if exists {
		return
	}

	_, err = db.Exec(fmt.Sprintf(`CREATE DATABASE "%s"`, dbName))
	if err != nil {
		log.Fatalf("failed create database %s: %v", dbName, err)
	}

	log.Printf("database %s created", dbName)
}
