// Package generator provides functionality for generating and managing database migration files.
// It includes tools for initializing source code, creating migration templates,
// and generating timestamped migration files.
//
// The package handles three main template types:
// - Migration script template - For individual migration files
// - Migrator template - For the migration executor
// - CLI template - For command line interface
//
// Templates are stored as base64 encoded strings and decoded at runtime.
// The package provides functions to:
// - Load and parse templates
// - Initialize source code structure
// - Generate new migration files with timestamps
// - Check initialization status
//
// Usage requires a GomigerConfig that specifies:
// - Package name for generated files
// - Target path for migrations
//
// Generated migration files follow the naming convention:
// YYYYMMDDHHMM_name.mg.go
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

// Template represents a parsed Go source file template.
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
			return nil, fmt.Errorf("failed to decode template content: %w", decodeErr)
		}
		fs := token.NewFileSet()
		node, err := parser.ParseFile(fs, "", content, parser.ParseComments)
		if err != nil {
			return nil, fmt.Errorf("failed to parse template content to ast.File node: %w", err)
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
		return fmt.Errorf("cannot load the templates: %w", err)
	}
	migrator := templates[1]
	cli := templates[2]
	helper.UpdatePackageName(migrator.node, rc.PkgName)
	helper.UpdatePackageName(cli.node, rc.PkgName)
	/// init the migration folder
	//nolint:gosec
	if err := os.MkdirAll(rc.Path, os.ModePerm); err != nil {
		return fmt.Errorf("cannot init the migration folder: %w", err)
	}
	/// init the migrator file
	if err := helper.ExportFile(migrator.node, migrator.fs, rc.Path+"/migrator.mg.go"); err != nil {
		return fmt.Errorf("cannot init the migrator file: %w", err)
	}
	/// init the cli file
	if err := helper.ExportFile(cli.node, cli.fs, rc.Path+"/cli.mg.go"); err != nil {
		return fmt.Errorf("cannot init the cli file: %w", err)
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
	timestamp := time.Now().Format("200601021504")
	filePath := filepath.Join(rc.Path, fmt.Sprintf("%s_%s.mg.go", timestamp, name))

	templates, err := LoadTemplates()
	migration := templates[0]
	if err != nil {
		return fmt.Errorf("cannot load the templates: %w", err)
	}

	helper.UpdatePackageName(migration.node, rc.PkgName)
	helper.UpdateFuncName(migration.node, "MigrationNameUp", fmt.Sprintf("Migration_%s_%s_Up", timestamp, name))
	helper.UpdateFuncName(migration.node, "MigrationNameDown", fmt.Sprintf("Migration_%s_%s_Down", timestamp, name))
	helper.UpdateFuncName(migration.node, "MigrationNameVersion", fmt.Sprintf("Migration_%s_%s_Version", timestamp, name))
	helper.UpdateStringValue(migration.node, "__VERSION__", timestamp)

	err = helper.ExportFile(migration.node, migration.fs, filePath)
	if err != nil {
		return fmt.Errorf("cannot generate the migration file: %w", err)
	}
	return nil
}
