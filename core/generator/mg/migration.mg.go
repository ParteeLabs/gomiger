package main

import (
	"context"
)

//nolint:godoclint,revive
func (m *Migrator) MigrationNameUp(ctx context.Context) error {
	/** Your migration up code here: */
	return nil
}

//nolint:godoclint,revive
func (m *Migrator) MigrationNameDown(ctx context.Context) error {
	/** Your migration down code here: */
	return nil
}

// AUTO GENERATED, DO NOT MODIFY!
//
//nolint:godoclint
func (m *Migrator) MigrationNameVersion() string {
	return "__VERSION__"
}
