package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// =================================================================
// USER COLLECTION SCHEMA
// =================================================================
// Collection: "users"
//
// This file is the **source of truth** for the users collection.
// Any field changes here should be reflected in a migration file.
// =================================================================

// User represents the users collection document in MongoDB.
// Profile is embedded as a subdocument (1:1 relationship).
type User struct {
	ID            bson.ObjectID `bson:"_id,omitempty" json:"id"`
	Email         string        `bson:"email" json:"email"`
	Password      string        `bson:"password" json:"-"`
	Role          string        `bson:"role" json:"role"` // "admin" | "user"
	EmailVerified bool          `bson:"email_verified" json:"email_verified"`
	Profile       Profile       `bson:"profile" json:"profile"`
	CreatedAt     time.Time     `bson:"created_at" json:"created_at"`
	UpdatedAt     time.Time     `bson:"updated_at" json:"updated_at"`
}

// Profile is an embedded subdocument within User.
type Profile struct {
	Name      string `bson:"name" json:"name"`
	ImageUrl  string `bson:"image_url,omitempty" json:"image_url"`
	BirthDate string `bson:"birth_date,omitempty" json:"birth_date"`
	Address   string `bson:"address,omitempty" json:"address"`
}

// Collection name constants
const (
	CollectionUsers = "users"
)

// Default role values
const (
	RoleAdmin = "admin"
	RoleUser  = "user"
)

// ValidRoles returns the list of valid roles.
func ValidRoles() []string {
	return []string{RoleAdmin, RoleUser}
}
