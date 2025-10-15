package migrations

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

//nolint:godoclint,revive
func (m *Migrator) Migration_202510152146_create_users_table_Up(ctx context.Context) error {
	// Create users collection with validation
	validator := bson.M{
		"$jsonSchema": bson.M{
			"bsonType": "object",
			"required": []string{"email", "created_at"},
			"properties": bson.M{
				"email": bson.M{
					"bsonType":    "string",
					"description": "must be a string and is required",
				},
				"name": bson.M{
					"bsonType":    "string",
					"description": "user's full name",
				},
				"created_at": bson.M{
					"bsonType":    "date",
					"description": "must be a date and is required",
				},
			},
		},
	}

	opts := options.CreateCollection().SetValidator(validator)
	return m.Db.CreateCollection(ctx, "users", opts)
}

//nolint:godoclint,revive
func (m *Migrator) Migration_202510152146_create_users_table_Down(ctx context.Context) error {
	return m.Db.Collection("users").Drop(ctx)
}

// AUTO GENERATED, DO NOT MODIFY!
//
//nolint:godoclint
func (m *Migrator) Migration_202510152146_create_users_table_Version() string {
	return "202510152146"
}
