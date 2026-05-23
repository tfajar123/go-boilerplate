package database

import (
	"database/sql"
	"go-boilerplate/ent"
	"log"
	"time"

	entsql "entgo.io/ent/dialect/sql"
	_ "github.com/lib/pq"
)

// DB is the underlying sql.DB connection pool, exposed for health checks
var DB *sql.DB

func NewEntClient(databaseURL string) *ent.Client {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		log.Fatalf("failed opening connection to postgres: %v", err)
	}

	// Connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(1 * time.Minute)

	// Expose for health checks
	DB = db

	drv := entsql.OpenDB("postgres", db)
	return ent.NewClient(ent.Driver(drv))
}
