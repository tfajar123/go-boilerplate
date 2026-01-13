package database

import (
	"context"
	"go-boilerplate/ent"
	"log"

	_ "github.com/lib/pq"
)

func NewEntClient(databaseURL string) *ent.Client {
	client, err := ent.Open("postgres", databaseURL)
	if err != nil {
		log.Fatalf("failed opening connection to postgres: %v", err)
	}

	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	return client
}
