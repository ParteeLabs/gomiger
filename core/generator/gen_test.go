package generator

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ParteeLabs/gomiger/core"
)

func TestLoadTemplates(t *testing.T) {
	t.Run("successfully loads all templates", func(t *testing.T) {
		templates, err := LoadTemplates()
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		// Should return 3 templates: migration, migrator, cli
		if len(templates) != 3 {
			t.Errorf("Expected 3 templates, got: %d", len(templates))
		}

		// Verify each template has valid AST nodes
		for i, tmpl := range templates {
			if tmpl.fs == nil {
				t.Errorf("Template %d: FileSet is nil", i)
			}
			if tmpl.node == nil {
				t.Errorf("Template %d: AST node is nil", i)
			}
			if tmpl.node != nil && tmpl.node.Name == nil {
				t.Errorf("Template %d: Package name is nil", i)
			}
		}
	})
}

func TestInitSrcCode(t *testing.T) {
	t.Run("creates migration folder and files", func(t *testing.T) {
		tmpDir := t.TempDir()
		migrationPath := filepath.Join(tmpDir, "migrations")

		rc := &core.GomigerConfig{
			Path:    migrationPath,
			PkgName: "testmigrations",
		}

		err := InitSrcCode(rc)
		if err != nil {
			t.Fatalf("InitSrcCode failed: %v", err)
		}

		// Verify migration folder was created
		if _, err := os.Stat(migrationPath); os.IsNotExist(err) {
			t.Error("Migration folder was not created")
		}

		// Verify migrator.mg.go was created
		migratorPath := filepath.Join(migrationPath, "migrator.mg.go")
		if _, err := os.Stat(migratorPath); os.IsNotExist(err) {
			t.Error("migrator.mg.go was not created")
		}

		// Verify cli.mg.go was created
		cliPath := filepath.Join(migrationPath, "cli.mg.go")
		if _, err := os.Stat(cliPath); os.IsNotExist(err) {
			t.Error("cli.mg.go was not created")
		}

		// Verify package name in migrator file
		migratorContent, err := os.ReadFile(migratorPath)
		if err != nil {
			t.Fatalf("Failed to read migrator file: %v", err)
		}
		if !strings.Contains(string(migratorContent), "package testmigrations") {
			t.Error("migrator.mg.go does not have correct package name")
		}

		// Verify package name in cli file
		cliContent, err := os.ReadFile(cliPath)
		if err != nil {
			t.Fatalf("Failed to read cli file: %v", err)
		}
		if !strings.Contains(string(cliContent), "package testmigrations") {
			t.Error("cli.mg.go does not have correct package name")
		}
	})

	t.Run("creates nested path structure", func(t *testing.T) {
		tmpDir := t.TempDir()
		migrationPath := filepath.Join(tmpDir, "db", "migrations", "v1")

		rc := &core.GomigerConfig{
			Path:    migrationPath,
			PkgName: "v1",
		}

		err := InitSrcCode(rc)
		if err != nil {
			t.Fatalf("InitSrcCode with nested path failed: %v", err)
		}

		// Verify all nested directories were created
		if _, err := os.Stat(migrationPath); os.IsNotExist(err) {
			t.Error("Nested migration folder was not created")
		}

		// Verify files exist
		if _, err := os.Stat(filepath.Join(migrationPath, "migrator.mg.go")); os.IsNotExist(err) {
			t.Error("migrator.mg.go was not created in nested path")
		}
	})

	t.Run("returns error when path is invalid", func(t *testing.T) {
		rc := &core.GomigerConfig{
			Path:    "/invalid/\x00/path", // null byte makes it invalid
			PkgName: "test",
		}

		err := InitSrcCode(rc)
		if err == nil {
			t.Error("Expected error for invalid path")
		}
	})
}

func TestIsSrcCodeInitialized(t *testing.T) {
	t.Run("returns false when not initialized", func(t *testing.T) {
		tmpDir := t.TempDir()

		rc := &core.GomigerConfig{
			Path: tmpDir,
		}

		if IsSrcCodeInitialized(rc) {
			t.Error("Expected IsSrcCodeInitialized to return false")
		}
	})

	t.Run("returns true when initialized", func(t *testing.T) {
		tmpDir := t.TempDir()

		rc := &core.GomigerConfig{
			Path:    tmpDir,
			PkgName: "test",
		}

		// Initialize the source code
		if err := InitSrcCode(rc); err != nil {
			t.Fatalf("Failed to initialize: %v", err)
		}

		if !IsSrcCodeInitialized(rc) {
			t.Error("Expected IsSrcCodeInitialized to return true")
		}
	})

	t.Run("returns true when only migrator.mg.go exists", func(t *testing.T) {
		tmpDir := t.TempDir()

		// Create just the migrator file
		migratorPath := filepath.Join(tmpDir, "migrator.mg.go")
		if err := os.WriteFile(migratorPath, []byte("package test"), 0644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		rc := &core.GomigerConfig{
			Path: tmpDir,
		}

		if !IsSrcCodeInitialized(rc) {
			t.Error("Expected IsSrcCodeInitialized to return true when migrator.mg.go exists")
		}
	})
}

func TestGenMigrationFile(t *testing.T) {
	t.Run("creates migration file with correct naming", func(t *testing.T) {
		tmpDir := t.TempDir()

		rc := &core.GomigerConfig{
			Path:    tmpDir,
			PkgName: "migrations",
		}

		migrationName := "create_users_table"
		err := GenMigrationFile(rc, migrationName)
		if err != nil {
			t.Fatalf("GenMigrationFile failed: %v", err)
		}

		// Find the generated file
		entries, err := os.ReadDir(tmpDir)
		if err != nil {
			t.Fatalf("Failed to read directory: %v", err)
		}

		if len(entries) == 0 {
			t.Fatal("No migration file was created")
		}

		// Verify filename format: YYYYMMDDHHMM_name.mg.go
		filename := entries[0].Name()
		if !strings.HasSuffix(filename, "_create_users_table.mg.go") {
			t.Errorf("Generated file has incorrect name: %s", filename)
		}

		// Verify timestamp format (12 digits)
		parts := strings.Split(filename, "_")
		if len(parts) < 2 {
			t.Error("Filename doesn't contain timestamp")
		} else if len(parts[0]) != 12 {
			t.Errorf("Timestamp should be 12 digits, got: %s (len=%d)", parts[0], len(parts[0]))
		}

		// Read and verify file content
		filePath := filepath.Join(tmpDir, filename)
		content, err := os.ReadFile(filePath)
		if err != nil {
			t.Fatalf("Failed to read generated file: %v", err)
		}

		contentStr := string(content)

		// Verify package name
		if !strings.Contains(contentStr, "package migrations") {
			t.Error("Generated file doesn't have correct package name")
		}

		// Verify function names contain timestamp and migration name
		expectedFuncPrefix := "Migration_" + parts[0] + "_create_users_table"
		if !strings.Contains(contentStr, expectedFuncPrefix+"_Up") {
			t.Errorf("Generated file doesn't contain Up function: %s_Up", expectedFuncPrefix)
		}
		if !strings.Contains(contentStr, expectedFuncPrefix+"_Down") {
			t.Errorf("Generated file doesn't contain Down function: %s_Down", expectedFuncPrefix)
		}
		if !strings.Contains(contentStr, expectedFuncPrefix+"_Version") {
			t.Errorf("Generated file doesn't contain Version function: %s_Version", expectedFuncPrefix)
		}

		// Verify version string
		if !strings.Contains(contentStr, parts[0]) {
			t.Errorf("Generated file doesn't contain version timestamp: %s", parts[0])
		}
	})

	t.Run("handles migration name with underscores", func(t *testing.T) {
		tmpDir := t.TempDir()

		rc := &core.GomigerConfig{
			Path:    tmpDir,
			PkgName: "migrations",
		}

		err := GenMigrationFile(rc, "add_user_email_index")
		if err != nil {
			t.Fatalf("GenMigrationFile failed: %v", err)
		}

		entries, err := os.ReadDir(tmpDir)
		if err != nil || len(entries) == 0 {
			t.Fatal("Migration file was not created")
		}

		// Verify the file contains the migration name
		if !strings.Contains(entries[0].Name(), "add_user_email_index") {
			t.Errorf("Filename doesn't contain migration name: %s", entries[0].Name())
		}
	})

	t.Run("handles empty migration name", func(t *testing.T) {
		tmpDir := t.TempDir()

		rc := &core.GomigerConfig{
			Path:    tmpDir,
			PkgName: "migrations",
		}

		err := GenMigrationFile(rc, "")
		if err != nil {
			t.Fatalf("GenMigrationFile failed: %v", err)
		}

		entries, err := os.ReadDir(tmpDir)
		if err != nil || len(entries) == 0 {
			t.Fatal("Migration file was not created")
		}

		// Should still create a file with just timestamp
		filename := entries[0].Name()
		if !strings.HasSuffix(filename, "_.mg.go") {
			t.Errorf("Expected filename to end with _.mg.go, got: %s", filename)
		}
	})

	t.Run("returns error for invalid path", func(t *testing.T) {
		rc := &core.GomigerConfig{
			Path:    "/invalid/\x00/path",
			PkgName: "test",
		}

		err := GenMigrationFile(rc, "test")
		if err == nil {
			t.Error("Expected error for invalid path")
		}
	})
}

func TestIntegration_FullWorkflow(t *testing.T) {
	t.Run("complete initialization and generation workflow", func(t *testing.T) {
		tmpDir := t.TempDir()
		migrationPath := filepath.Join(tmpDir, "migrations")

		rc := &core.GomigerConfig{
			Path:    migrationPath,
			PkgName: "migrations",
		}

		// Step 1: Check not initialized
		if IsSrcCodeInitialized(rc) {
			t.Error("Should not be initialized initially")
		}

		// Step 2: Initialize
		if err := InitSrcCode(rc); err != nil {
			t.Fatalf("Failed to initialize: %v", err)
		}

		// Step 3: Verify initialized
		if !IsSrcCodeInitialized(rc) {
			t.Error("Should be initialized after InitSrcCode")
		}

		// Step 4: Generate first migration
		if err := GenMigrationFile(rc, "initial_schema"); err != nil {
			t.Fatalf("Failed to generate first migration: %v", err)
		}

		// Step 5: Generate second migration
		if err := GenMigrationFile(rc, "add_indexes"); err != nil {
			t.Fatalf("Failed to generate second migration: %v", err)
		}

		// Step 6: Verify all files exist
		entries, err := os.ReadDir(migrationPath)
		if err != nil {
			t.Fatalf("Failed to read directory: %v", err)
		}

		// Should have: migrator.mg.go, cli.mg.go, and 2 generated migrations
		if len(entries) != 4 {
			t.Errorf("Expected 4 files, got: %d", len(entries))
			for _, entry := range entries {
				t.Logf("  - %s", entry.Name())
			}
		}

		// Verify file types
		var hasMigrator, hasCli, migrationCount int
		for _, entry := range entries {
			name := entry.Name()
			if name == "migrator.mg.go" {
				hasMigrator++
			} else if name == "cli.mg.go" {
				hasCli++
			} else if strings.HasSuffix(name, ".mg.go") {
				migrationCount++
			}
		}

		if hasMigrator != 1 {
			t.Errorf("Expected 1 migrator.mg.go, found %d", hasMigrator)
		}
		if hasCli != 1 {
			t.Errorf("Expected 1 cli.mg.go, found %d", hasCli)
		}
		if migrationCount != 2 {
			t.Errorf("Expected 2 migration files, found %d", migrationCount)
		}
	})
}
