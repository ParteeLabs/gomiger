package gomiger

import (
	"context"
	"fmt"
	"time"
)

// SchemaStatus is the status of the schema
type SchemaStatus string

var (
	// InProgress is for a running migration
	InProgress SchemaStatus = "in_progress"
	// Dirty is for a failed migration
	Dirty SchemaStatus = "dirty"
	// Applied is for a completed migration
	Applied SchemaStatus = "applied"
)

// Schema is a log for migrations
type Schema struct {
	Version   string       `json:"version" bson:"version" validate:"required"`
	Timestamp time.Time    `json:"timestamp" bson:"timestamp" validate:"required"`
	Status    SchemaStatus `json:"status" bson:"status" validate:"required"`
}

// Gomiger is the interface for the migrator
type Gomiger interface {
	Up(ctx context.Context, toVersion string) error
	Down(ctx context.Context, atVersion string) error
	Connect(ctx context.Context) error
}

// DbPlugin is the interface for the plugin
type DbPlugin interface {
	Connect(ctx context.Context) error
	GetSchema(ctx context.Context, version string) (Schema, error)
	ApplyMigration(ctx context.Context, mi Migration) error
	RevertMigration(ctx context.Context, mi Migration) error
}

// MutationFunc is a function that applies a migration.
type MutationFunc = func(context context.Context) error

// Migration contain a version name and a mutation function.
type Migration struct {
	Version string
	Up      MutationFunc
	Down    MutationFunc
}

// BaseMigrator is the base migrator for control flow.
type BaseMigrator struct {
	DB         DbPlugin
	Migrations []Migration
}

func (b *BaseMigrator) isVersionExists(version string) bool {
	for _, mi := range b.Migrations {
		if mi.Version == version {
			return true
		}
	}
	return false
}

// Up updates the database to a specific version.
func (b *BaseMigrator) Up(ctx context.Context, toVersion string) error {
	if toVersion != "" && !b.isVersionExists(toVersion) {
		return fmt.Errorf("version %s does not exist", toVersion)
	}
	for _, mi := range b.Migrations {
		schema, err := b.DB.GetSchema(ctx, mi.Version)
		if err != nil {
			return fmt.Errorf("failed to get schema: %w", err)
		}
		if schema.Status == Applied || schema.Status == Dirty {
			continue
		}
		if err := b.DB.ApplyMigration(ctx, mi); err != nil {
			return fmt.Errorf("failed to apply migration %s: %w", mi.Version, err)
		}
		if mi.Version == toVersion {
			return nil
		}
	}
	return nil
}

// Down reverts the database to a specific version.
func (b *BaseMigrator) Down(ctx context.Context, atVersion string) error {
	if atVersion == "" {
		return fmt.Errorf("a version is required")
	}
	if !b.isVersionExists(atVersion) {
		return fmt.Errorf("version %s does not exist", atVersion)
	}
	for i := len(b.Migrations); i >= 0; i-- {
		mi := b.Migrations[i]
		schema, err := b.DB.GetSchema(ctx, mi.Version)
		if err != nil {
			return fmt.Errorf("failed to get schema: %w", err)
		}
		if schema.Status == Applied || schema.Status == Dirty {
			continue
		}
		if err := b.DB.RevertMigration(ctx, mi); err != nil {
			return fmt.Errorf("failed to revert migration %s: %w", mi.Version, err)
		}
		if mi.Version == atVersion {
			return nil
		}
	}
	return nil
}
