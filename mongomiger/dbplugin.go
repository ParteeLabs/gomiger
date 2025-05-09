package mongomiger

import (
	"context"
	"fmt"

	"github.com/ParteeLabs/gomiger/core"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/x/mongo/driver/connstring"
)

// Mongomiger implements core.DbPlugin for MongoDB.
type Mongomiger struct {
	uri              string
	Client           *mongo.Client
	Db               *mongo.Database
	schemaStore      string
	schemaCollection *mongo.Collection
}

// NewMongomiger creates a new Mongomiger plugin.
func NewMongomiger(cfg *core.GomigerConfig) *Mongomiger {
	return &Mongomiger{
		uri:         cfg.URI,
		schemaStore: cfg.SchemaStore,
	}
}

// Connect implements core.DbPlugin.
func (m *Mongomiger) Connect(_ context.Context) (err error) {
	// Parse the connection string to get the database name.
	connStr, err := connstring.Parse(m.uri)
	if err != nil {
		return
	}
	// Connect and get the schema collection.
	if m.Client, err = mongo.Connect(options.Client().ApplyURI(m.uri)); err != nil {
		return
	}
	m.Db = m.Client.Database(connStr.Database)
	m.schemaCollection = m.Db.Collection(m.schemaStore)
	return
}

// GetSchema implements core.DbPlugin.
func (m *Mongomiger) GetSchema(ctx context.Context, version string) (schema *core.Schema, err error) {
	err = m.schemaCollection.FindOne(ctx, bson.M{"version": version}).Decode(&schema)
	if err != nil {
		return nil, fmt.Errorf("failed to get schema: %w", err)
	}
	return
}

func (m *Mongomiger) updateSchemaStatus(ctx context.Context, mi core.Migration, status core.SchemaStatus) error {
	if _, err := m.schemaCollection.UpdateOne(ctx, bson.M{"version": mi.Version}, bson.M{"status": status}); err != nil {
		return fmt.Errorf("failed to update schema status at version: %s, please manually update it with '%s' then try again, Error: %w", mi.Version, status, err)
	}
	return nil
}

// ApplyMigration implements core.DbPlugin.
func (m *Mongomiger) ApplyMigration(ctx context.Context, mi core.Migration) error {
	// Mark the migration as in progress (create a new schema).
	if _, err := m.schemaCollection.InsertOne(ctx, bson.M{"status": core.InProgress}); err != nil {
		return fmt.Errorf("failed to apply migration at version: %s, Error: %w", mi.Version, err)
	}
	// Run the migration.
	if err := mi.Up(ctx); err != nil {
		// Mark the migration as dirty.
		if err := m.updateSchemaStatus(ctx, mi, core.Dirty); err != nil {
			return err
		}
		return fmt.Errorf("failed to apply migration %s: %w", mi.Version, err)
	}
	// Mark the migration as applied.
	if err := m.updateSchemaStatus(ctx, mi, core.Applied); err != nil {
		return err
	}
	return nil
}

// RevertMigration implements core.DbPlugin.
func (m *Mongomiger) RevertMigration(ctx context.Context, mi core.Migration) error {
	if err := mi.Down(ctx); err != nil {
		// Mark the migration as dirty.
		if err := m.updateSchemaStatus(ctx, mi, core.Dirty); err != nil {
			return err
		}
		return fmt.Errorf("failed to apply migration %s: %w", mi.Version, err)
	}
	// Delete the schema.
	if _, err := m.schemaCollection.DeleteOne(ctx, bson.M{"version": mi.Version}); err != nil {
		return fmt.Errorf("failed to delete schema at version: %s, please manually delete it, Error: %w", mi.Version, err)
	}
	return nil

}

// Interface check
var _ core.DbPlugin = (*Mongomiger)(nil)
