package mgr

import (
	"github.com/ParteeLabs/gomiger/core"
	"github.com/ParteeLabs/gomiger/mongomiger"
)

// Migrator is the main migrator struct.
type Migrator struct {
	core.BaseMigrator
	*mongomiger.Mongomiger
	Config *core.GomigerConfig
}

// NewMigrator creates a new migrator.
func NewMigrator(config *core.GomigerConfig) core.Gomiger {
	mongomiger := mongomiger.NewMongomiger(config)
	m := &Migrator{
		BaseMigrator: core.BaseMigrator{
			DbPlugin: mongomiger,
		},
		Mongomiger: mongomiger,
		Config:     config,
	}

	// ** Add your migrations here **
	m.Migrations = []core.Migration{
		// {Version: "1.0.0", Up: m.MigrationNameUp, Down: m.MigrationNameDown},
	}
	return m
}
