package main

import (
	"context"
)

func (m *Migrator) MigrationNameUp(ctx context.Context) error {
	/** Your migration up code here: */
	return nil
}

func (m *Migrator) MigrationNameDown(ctx context.Context) error {
	/** Your migration down code here: */
	return nil
}

// AUTO GENERATED, DO NOT MODIFY!
func (m *Migrator) MigrationNameVersion() string {
	return "__VERSION__"
}
