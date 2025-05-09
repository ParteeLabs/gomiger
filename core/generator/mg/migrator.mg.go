package main

import (
	"context"

	"github.com/ParteeLabs/gomiger/core"
)

// Migrator is the main migrator struct.
type Migrator struct {
	core.BaseMigrator
	Config *core.GomigerConfig
}

// NewMigrator creates a new migrator.
func NewMigrator(config *core.GomigerConfig) core.Gomiger {
	m := &Migrator{
		Config: config,
	}
	// m.DB = NewMongomiger() ** Add your plugin here **

	// ** Add your migrations here **
	m.Migrations = []core.Migration{
		// {Version: "1.0.0", Up: m.MigrationNameUp, Down: m.MigrationNameDown},
	}
	return m
}

// Connect connects database.
func (m *Migrator) Connect(ctx context.Context) (err error) {
	return m.DB.Connect(ctx)
}
