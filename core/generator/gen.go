package generator

import (
	"encoding/base64"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"time"

	"github.com/ParteeLabs/gomiger/core"
	"github.com/ParteeLabs/gomiger/core/generator/helper"
)

type Template struct {
	fs   *token.FileSet
	node *ast.File
}

// LoadTemplates load the preset template strings to ast.Node
func LoadTemplates() ([]Template, error) {
	encodedTemplates := []string{
		MigrationScriptTemplateBase64,
		MigratorTemplateBase64,
		CliTemplateBase64,
	}
	templates := make([]Template, 0)

	for _, encoded := range encodedTemplates {
		content, decodeErr := base64.StdEncoding.DecodeString(encoded)
		if decodeErr != nil {
			return nil, decodeErr
		}
		fs := token.NewFileSet()
		node, err := parser.ParseFile(fs, "", content, parser.ParseComments)
		if err != nil {
			return nil, err
		}
		templates = append(templates, Template{
			fs,
			node,
		})
	}
	return templates, nil
}

// InitSrcCode initializes the source code
func InitSrcCode(rc *core.GomigerConfig) error {
	templates, err := LoadTemplates()
	if err != nil {
		return fmt.Errorf("Cannot load the templates: %w", err)
	}
	migrator := templates[1]
	cli := templates[2]
	helper.UpdatePackageName(migrator.node, rc.PkgName)
	helper.UpdatePackageName(cli.node, rc.PkgName)
	/// init the migration folder
	if err := os.MkdirAll(rc.Path, os.ModePerm); err != nil {
		return fmt.Errorf("Cannot init the migration folder: %w", err)
	}
	/// init the migrator file
	if err := helper.ExportFile(migrator.node, migrator.fs, rc.Path+"/migrator.mg.go"); err != nil {
		return fmt.Errorf("Cannot init the migrator file: %w", err)
	}
	/// init the cli file
	if err := helper.ExportFile(cli.node, cli.fs, rc.Path+"/cli.mg.go"); err != nil {
		return fmt.Errorf("Cannot init the cli file: %w", err)
	}
	return nil
}

// IsSrcCodeInitialized checks if the source code is initialized.
func IsSrcCodeInitialized(rc *core.GomigerConfig) bool {
	_, err := os.Stat(rc.Path + "/migrator.mg.go")
	return err == nil
}

// GenMigrationFile generates a migration file
func GenMigrationFile(rc *core.GomigerConfig, name string) error {
	// Create the migration file path
	timestamp := time.Now().Format("20060102150405")
	filePath := filepath.Join(rc.Path, fmt.Sprintf("%s_%s.mg.go", timestamp, name))

	templates, err := LoadTemplates()
	migration := templates[0]
	if err != nil {
		return fmt.Errorf("Cannot load the templates: %w", err)
	}
	helper.UpdatePackageName(migration.node, rc.PkgName)
	helper.UpdateFuncName(migration.node, "MigrationNameUp", fmt.Sprintf("Migration_%s_Up", name))
	helper.UpdateFuncName(migration.node, "MigrationNameDown", fmt.Sprintf("Migration_%s_Down", name))

	return helper.ExportFile(migration.node, migration.fs, filePath)
}
