package database

import (
	"go-boilerplate/ent"
	"log"

	_ "github.com/lib/pq"
)

func NewEntClient(databaseURL string) *ent.Client {
	client, err := ent.Open("postgres", databaseURL)
	if err != nil {
		log.Fatalf("failed opening connection to postgres: %v", err)
	}
	return client
}
