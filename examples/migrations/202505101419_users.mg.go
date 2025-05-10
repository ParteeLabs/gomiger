package mgr

import (
	"context"
)

func (m *Migrator) Migration_202505101419_users_Up(ctx context.Context) error {
	/** Your migration up code here: */
	return nil
}

func (m *Migrator) Migration_202505101419_users_Down(ctx context.Context) error {
	/** Your migration down code here: */
	return nil
}

// AUTO GENERATED, DO NOT MODIFY!
func (m *Migrator) Migration_202505101419_users_Version() string {
	return "202505101419"
}
