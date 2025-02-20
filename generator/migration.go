package generator

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/ParteeLabs/gomiger/config"
	"github.com/ParteeLabs/gomiger/generator/helper"
)

func GenMigrationFile(rc *config.GomigerRC, name string) error {
	// Create the migration file path
	timestamp := time.Now().Format("20060102150405")
	filePath := filepath.Join(rc.Path, fmt.Sprintf("%s_%s.mg.go", timestamp, name))

	migrationNode, _, err := LoadTemplates()
	if err != nil {
		return fmt.Errorf("Cannot load the templates: %w", err)
	}
	helper.UpdatePackageName(migrationNode, rc.PkgName)
	helper.UpdateFuncName(migrationNode, "MigrationNameUp", fmt.Sprintf("Migration_%s_Up", name))
	helper.UpdateFuncName(migrationNode, "MigrationNameDown", fmt.Sprintf("Migration_%s_Down", name))

	return helper.ExportFile(migrationNode, filePath)
}
