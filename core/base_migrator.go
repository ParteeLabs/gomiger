package core

import (
	"context"
	"fmt"
)

// BaseMigratorAbstractMethods defines the methods that must be implemented by a concrete migrator.
type BaseMigratorAbstractMethods interface {
	Connect(ctx context.Context) error
	GetSchema(ctx context.Context, version string) (*Schema, error)
	ApplyMigration(ctx context.Context, mi Migration) error
	RevertMigration(ctx context.Context, mi Migration) error
}

// BaseMigrator is the base migrator for controlling flow.
// It does not connect or execute to any database.
type BaseMigrator struct {
	BaseMigratorAbstractMethods
	Migrations []Migration
}

var _ Gomiger = (*BaseMigrator)(nil)

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
		schema, err := b.GetSchema(ctx, mi.Version)
		if err != nil {
			return fmt.Errorf("failed to get schema: %w", err)
		}
		if schema.Status == Applied || schema.Status == Dirty {
			continue
		}
		if err := b.ApplyMigration(ctx, mi); err != nil {
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
	for i := len(b.Migrations) - 1; i >= 0; i-- {
		mi := b.Migrations[i]
		schema, err := b.GetSchema(ctx, mi.Version)
		if err != nil {
			return fmt.Errorf("failed to get schema: %w", err)
		}
		if schema.Status != Applied && schema.Status != Dirty {
			continue
		}
		if err := b.RevertMigration(ctx, mi); err != nil {
			return fmt.Errorf("failed to revert migration %s: %w", mi.Version, err)
		}
		if mi.Version == atVersion {
			return nil
		}
	}
	return nil
}
