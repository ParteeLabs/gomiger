package generator

var MigrationScriptTemplate = `package main

import "context"

func (m *Migrator) MigrationNameUp(ctx context.Context) error {
	/** Your migration up code here: */
	return nil
}

func (m *Migrator) MigrationNameDown(ctx context.Context) error {
	/** Your migration down code here: */
	return nil
}
`
var MigratorTemplate = `package main

import (
	"context"
	"fmt"

	"github.com/ParteeLabs/gomiger/migrator"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MutationFunc is a function that applies a migration.
type MutationFunc = func(context context.Context) error

// Migration contain a version name and a mutation function.
type Migration struct {
	Version  string
	Mutation MutationFunc
}

// Migrator is the main migrator struct.
type Migrator struct {
	Config      migrator.Config
	MongoClient *mongo.Client

	Ups   []Migration
	Downs []Migration
}

// NewMigrator creates a new migrator.
func NewMigrator(config migrator.Config) *Migrator {
	return &Migrator{
		Config: config,
		Ups:    []Migration{},
		Downs:  []Migration{},
	}
}

// MongoConnect connects to the MongoDB database.
func (m *Migrator) MongoConnect(ctx context.Context) (err error) {
	if m.MongoClient, err = mongo.Connect(ctx, options.Client().ApplyURI(m.Config.URI)); err != nil {
		return err
	}
	return nil
}

// Up applies the migrations.
func (m *Migrator) Up(ctx context.Context) error {
	for _, mi := range m.Ups {
		if err := mi.Mutation(ctx); err != nil {
			return fmt.Errorf("failed to apply migration %s: %w", mi.Version, err)
		}
	}
	return nil
}

// Down reverts the migrations.
func (m *Migrator) Down(ctx context.Context) error {
	for _, mi := range m.Downs {
		if err := mi.Mutation(ctx); err != nil {
			return fmt.Errorf("failed to revert migration %s: %w", mi.Version, err)
		}
	}
	return nil
}
`
