package migrations

import (
	"github.com/ParteeLabs/gomiger/core"
	"github.com/ParteeLabs/gomiger/mongomiger"
)

// Migrator is the main migrator struct.
type Migrator struct {
	*mongomiger.Mongomiger

	Config *core.GomigerConfig
}

// NewMigrator creates a new migrator.
func NewMigrator(config *core.GomigerConfig) core.Gomiger {
	m := &Migrator{
		Mongomiger: mongomiger.NewMongomiger(config),

		Config: config,
	}

	// ** Add your migrations here **
	m.Migrations = []core.Migration{
		// {Version: MigrationNameVersion(), Up: m.MigrationNameUp, Down: m.MigrationNameDown},
	}
	return m
}
