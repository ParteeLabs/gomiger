package main

import (
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
		BaseMigrator: core.BaseMigrator{
			// DbPlugin: mongomiger.NewMongomiger(config), ** Add your plugin here **
		},
		Config: config,
	}

	// ** Add your migrations here **
	m.Migrations = []core.Migration{
		// {Version: "1.0.0", Up: m.MigrationNameUp, Down: m.MigrationNameDown},
	}
	return m
}
