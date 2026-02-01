package core

import (
	"context"
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
	GetSchema(ctx context.Context, version string) (*Schema, error)
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
