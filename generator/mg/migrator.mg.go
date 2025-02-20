package main

import (
	"context"

	"github.com/ParteeLabs/gomiger"
)

// Migrator is the main migrator struct.
type Migrator struct {
	gomiger.BaseMigrator
	Config gomiger.Config
}

// NewMigrator creates a new migrator.
func NewMigrator(config gomiger.Config) gomiger.Gomiger {
	m := &Migrator{
		Config: config,
	}
	// m.DB = NewMongomiger() ** Add your plugin here **

	// ** Add your migrations here **
	m.Migrations = []gomiger.Migration{
		// {Version: "1.0.0", Up: m.MigrationNameUp, Down: m.MigrationNameDown},
	}
	return m
}

// Connect connects database.
func (m *Migrator) Connect(ctx context.Context) (err error) {
	// ** Call your plugin Connect func here **
	return nil
}
