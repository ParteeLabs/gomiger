package migrator

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// IsMigratedMongo returns true if the migration has already been applied.
func IsMigratedMongo(ctx context.Context, db mongo.Database, dbCfg Config, version string) bool {
	err := db.Collection(dbCfg.SchemaStore).FindOne(ctx, bson.M{
		"version": version,
	}).Err()
	return err == nil
}
