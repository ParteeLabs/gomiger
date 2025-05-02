package generator

import (
	"fmt"
	"os"

	"github.com/ParteeLabs/gomiger/config"
	"github.com/ParteeLabs/gomiger/generator/helper"
)

// InitSrcCode initializes the source code
func InitSrcCode(rc *config.GomigerRC) error {
	_, migratorNode, err := LoadTemplates()
	if err != nil {
		return fmt.Errorf("Cannot load the templates: %w", err)
	}
	helper.UpdatePackageName(migratorNode, rc.PkgName)

	/// init the migration folder
	if err := os.MkdirAll(rc.Path, os.ModePerm); err != nil {
		return fmt.Errorf("Cannot init the migration folder: %w", err)
	}

	/// init the migrator file
	return helper.ExportFile(migratorNode, rc.Path+"/migrator.mg.go")
}

// IsSrcCodeInitialized checks if the source code is initialized.
func IsSrcCodeInitialized(rc *config.GomigerRC) bool {
	_, err := os.Stat(rc.Path + "/migrator.mg.go")
	return err == nil
}
