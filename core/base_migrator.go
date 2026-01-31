package core

import (
	"context"
	"errors"
	"fmt"
)

// BaseMigrator is the base migrator for controlling flow.
// It does not connect or execute to any database.
type BaseMigrator struct {
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

// Connect connects to the database.
func (b *BaseMigrator) Connect(_ context.Context) error {
	return errors.New("not implemented")
}

// GetSchema returns the schema at a specific version.
func (b *BaseMigrator) GetSchema(_ context.Context, _ string) (*Schema, error) {
	return nil, errors.New("not implemented")
}

// ApplyMigration applies a migration.
func (b *BaseMigrator) ApplyMigration(_ context.Context, _ Migration) error {
	return errors.New("not implemented")
}

// RevertMigration reverts a migration.
func (b *BaseMigrator) RevertMigration(_ context.Context, _ Migration) error {
	return errors.New("not implemented")
}
