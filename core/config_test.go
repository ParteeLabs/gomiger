package core

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var validConfigContent = `path: './test-migrations'
pkg_name: 'testmgr'
schema_store: 'test_schemas'`

var invalidConfigContent = `
path: './test-migrations'
pkg_name: 'testmgr'
schema_store: 'test_schemas'
	invalid_indentation
missing_colon_value
`

func createTempConfigFile(t *testing.T, content string) (string, func()) {
	tempFile, err := os.CreateTemp("", "gomiger-test-*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	if _, err := tempFile.WriteString(content); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
		os.Remove(tempFile.Name())
	}
	tempFile.Close()
	return tempFile.Name(), func() {
		os.Remove(tempFile.Name())
	}
}

func TestGomigerConfig(t *testing.T) {
	t.Run("ParseYAML with valid config", func(t *testing.T) {
		tempFile, cleanup := createTempConfigFile(t, validConfigContent)
		defer cleanup()
		config := &GomigerConfig{}
		if err := config.ParseYAML(tempFile); err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if config.Path != "./test-migrations" {
			t.Errorf("Expected path './test-migrations', got: %s", config.Path)
		}

		if config.PkgName != "testmgr" {
			t.Errorf("Expected pkg_name 'testmgr', got: %s", config.PkgName)
		}

		if config.SchemaStore != "test_schemas" {
			t.Errorf("Expected schema_store 'test_schemas', got: %s", config.SchemaStore)
		}
	})

	t.Run("ParseYAML with empty file path", func(t *testing.T) {
		if err := os.WriteFile("gomiger.rc.yaml", []byte(validConfigContent), 0644); err != nil {
			t.Fatalf("Failed to write default config file: %v", err)
		}
		defer os.Remove("gomiger.rc.yaml")

		config := &GomigerConfig{}
		if err := config.ParseYAML(""); err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}
	})

	t.Run("ParseYAML with invalid config", func(t *testing.T) {
		tempFile, cleanup := createTempConfigFile(t, invalidConfigContent)
		defer cleanup()
		config := &GomigerConfig{}
		err := config.ParseYAML(tempFile)
		if err == nil {
			t.Error("Expected error for invalid config")
		}

		if !strings.Contains(err.Error(), "cannot parse the gomiger.rc file") {
			t.Errorf("Expected specific error message, got: %s", err.Error())
		}
	})

	t.Run("ParseYAML with nonexistent file", func(t *testing.T) {
		config := &GomigerConfig{}
		err := config.ParseYAML("nonexistent.yaml")
		if err == nil {
			t.Error("Expected error for nonexistent file")
		}

		if !strings.Contains(err.Error(), "cannot read the gomiger.rc file") {
			t.Errorf("Expected specific error message, got: %s", err.Error())
		}
	})

	t.Run("PopulateAndValidate with empty path", func(t *testing.T) {
		config := &GomigerConfig{}
		err := config.PopulateAndValidate()
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		// Should set default path and convert to absolute
		absDefaultPath, _ := filepath.Abs("./migrations")
		if config.Path != absDefaultPath {
			t.Errorf("Expected default absolute path, got: %s", config.Path)
		}
	})

	t.Run("PopulateAndValidate with custom path", func(t *testing.T) {
		config := &GomigerConfig{
			Path: "./custom-migrations",
		}
		err := config.PopulateAndValidate()
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		// Should convert to absolute path
		absCustomPath, _ := filepath.Abs("./custom-migrations")
		if config.Path != absCustomPath {
			t.Errorf("Expected absolute custom path, got: %s", config.Path)
		}
	})

	t.Run("GetGomigerRC with URI from environment", func(t *testing.T) {
		tempFile, cleanup := createTempConfigFile(t, validConfigContent)
		defer cleanup()
		// Set environment variable
		os.Setenv("GOMIGER_URI", "mongodb://localhost:27017/test")
		defer os.Unsetenv("GOMIGER_URI")

		config, err := GetGomigerRC(tempFile)
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if config.URI != "mongodb://localhost:27017/test" {
			t.Errorf("Expected URI from environment, got: %s", config.URI)
		}
	})

	t.Run("GetGomigerRC with invalid config file", func(t *testing.T) {
		tempFile, cleanup := createTempConfigFile(t, invalidConfigContent)
		defer cleanup()

		_, err := GetGomigerRC(tempFile)
		if err == nil {
			t.Error("Expected error for invalid config file")
		}

		if !strings.Contains(err.Error(), "cannot parse the gomiger.rc file") {
			t.Errorf("Expected specific error message, got: %s", err.Error())
		}
	})
}

// Integration test to ensure config and migrator work together
func TestConfigIntegration(t *testing.T) {
	t.Run("Config with BaseMigrator", func(t *testing.T) {
		config := &GomigerConfig{
			Path:        "./test-migrations",
			PkgName:     "testmgr",
			URI:         "test://localhost",
			SchemaStore: "test_schema",
		}

		err := config.PopulateAndValidate()
		if err != nil {
			t.Errorf("Config validation failed: %v", err)
		}

		// Create a migrator with this config
		migrator := &BaseMigrator{
			Migrations: []Migration{
				{
					Version: "20241015_test",
					Up:      func(ctx context.Context) error { return nil },
					Down:    func(ctx context.Context) error { return nil },
				},
			},
		}

		// Verify migrator can work with the config
		if !migrator.isVersionExists("20241015_test") {
			t.Error("Expected migration to exist")
		}

		// Test that invalid version operations still fail appropriately
		err = migrator.Up(context.Background(), "invalid")
		if err == nil {
			t.Error("Expected error for invalid version")
		}
	})
}
