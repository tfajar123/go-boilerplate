package database

import (
	"go-boilerplate/models"
	"log"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

// Bootstrap performs startup database tasks:
// - Ensures all required indexes exist
//
// MongoDB auto-creates databases and collections on first write,
// so no equivalent of PostgreSQL's EnsureDatabaseExists is needed.
func Bootstrap(db *mongo.Database) {
	if err := models.EnsureIndexes(db); err != nil {
		log.Fatalf("failed to ensure mongodb indexes: %v", err)
	}
}
