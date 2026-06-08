package models

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// =================================================================
// INDEX DECLARATIONS
// =================================================================
// This file declares all MongoDB indexes for every collection.
// Called once at application startup via EnsureIndexes().
//
// Adding a new index:
//   1. Add the index model to the appropriate collection section
//   2. Run the application — indexes are created idempotently
// =================================================================

// EnsureIndexes creates all required indexes for all collections.
// This function is idempotent and safe to call on every startup.
func EnsureIndexes(db *mongo.Database) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// ── Users Collection ────────────────────────────────────────
	usersIndexes := []mongo.IndexModel{
		{
			// Unique email index — prevents duplicate registrations
			Keys:    bson.D{{Key: "email", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	}

	_, err := db.Collection(CollectionUsers).Indexes().CreateMany(ctx, usersIndexes)
	if err != nil {
		return err
	}

	return nil
}
