//go:build ignore

package main

import (
	"github.com/ParteeLabs/gomiger/core"
)

// Migrator is the main migrator struct.
type Migrator struct {
	// BaseMigrator doses not involve to any database. Use our plugins to connect to your database.
	// Or override Connect, GetSchema, ApplyMigration, RevertMigration methods to implement with your database.
	*core.BaseMigrator

	// *mongomiger.Mongomiger
	Config *core.GomigerConfig
}

// NewMigrator creates a new migrator.
func NewMigrator(config *core.GomigerConfig) core.Gomiger {
	m := &Migrator{
		// Mongomiger: mongomiger.NewMongomiger(config),
		Config: config,
	}

	// ** Add your migrations here **
	m.Migrations = []core.Migration{
		// {Version: MigrationNameVersion(), Up: m.MigrationNameUp, Down: m.MigrationNameDown},
	}
	return m
}
