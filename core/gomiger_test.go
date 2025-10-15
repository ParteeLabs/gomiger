package core

import (
	"context"
	"errors"
	"testing"
	"time"
)

// MockDbPlugin implements a mock database plugin for testing
type MockDbPlugin struct {
	schemas    map[string]*Schema
	connected  bool
	shouldFail bool
}

func NewMockDbPlugin() *MockDbPlugin {
	return &MockDbPlugin{
		schemas:   make(map[string]*Schema),
		connected: false,
	}
}

func (m *MockDbPlugin) Connect(_ context.Context) error {
	if m.shouldFail {
		return errors.New("mock connection failed")
	}
	m.connected = true
	return nil
}

func (m *MockDbPlugin) GetSchema(_ context.Context, version string) (*Schema, error) {
	if !m.connected {
		return nil, errors.New("not connected")
	}
	schema, exists := m.schemas[version]
	if !exists {
		return nil, nil // Schema not found
	}
	return schema, nil
}

func (m *MockDbPlugin) SaveSchema(_ context.Context, schema *Schema) error {
	if !m.connected {
		return errors.New("not connected")
	}
	m.schemas[schema.Version] = schema
	return nil
}

func (m *MockDbPlugin) DeleteSchema(_ context.Context, version string) error {
	if !m.connected {
		return errors.New("not connected")
	}
	delete(m.schemas, version)
	return nil
}

func (m *MockDbPlugin) SetShouldFail(fail bool) {
	m.shouldFail = fail
}

func TestSchema(t *testing.T) {
	t.Run("Schema creation", func(t *testing.T) {
		schema := &Schema{
			Version:   "20241015_test",
			Timestamp: time.Now(),
			Status:    InProgress,
		}

		if schema.Version != "20241015_test" {
			t.Errorf("Expected version '20241015_test', got '%s'", schema.Version)
		}

		if schema.Status != InProgress {
			t.Errorf("Expected status 'in_progress', got '%s'", schema.Status)
		}
	})

	t.Run("Schema status constants", func(t *testing.T) {
		if InProgress != "in_progress" {
			t.Errorf("Expected InProgress to be 'in_progress', got '%s'", InProgress)
		}

		if Dirty != "dirty" {
			t.Errorf("Expected Dirty to be 'dirty', got '%s'", Dirty)
		}

		if Applied != "applied" {
			t.Errorf("Expected Applied to be 'applied', got '%s'", Applied)
		}
	})
}

func TestMigration(t *testing.T) {
	t.Run("Migration creation", func(t *testing.T) {
		upCalled := false
		downCalled := false

		migration := Migration{
			Version: "20241015_test_migration",
			Up: func(_ context.Context) error {
				upCalled = true
				return nil
			},
			Down: func(_ context.Context) error {
				downCalled = true
				return nil
			},
		}

		// Test Up migration
		err := migration.Up(context.Background())
		if err != nil {
			t.Errorf("Expected no error for Up migration, got %v", err)
		}

		if !upCalled {
			t.Error("Expected Up function to be called")
		}

		// Test Down migration
		err = migration.Down(context.Background())
		if err != nil {
			t.Errorf("Expected no error for Down migration, got %v", err)
		}

		if !downCalled {
			t.Error("Expected Down function to be called")
		}
	})

	t.Run("Migration with error", func(t *testing.T) {
		expectedErr := errors.New("migration failed")
		migration := Migration{
			Version: "20241015_failing_migration",
			Up: func(_ context.Context) error {
				return expectedErr
			},
		}

		err := migration.Up(context.Background())
		if err == nil {
			t.Error("Expected error for failing migration")
		}

		if err.Error() != expectedErr.Error() {
			t.Errorf("Expected error '%s', got '%s'", expectedErr.Error(), err.Error())
		}
	})
}

// TestBaseMigrator tests the core migration logic
func TestBaseMigrator(t *testing.T) {
	t.Run("BaseMigrator initialization", func(t *testing.T) {
		migrator := &BaseMigrator{
			Migrations: []Migration{},
		}

		if migrator == nil {
			t.Error("Expected BaseMigrator to be created")
		}

		if len(migrator.Migrations) != 0 {
			t.Errorf("Expected empty migrations, got %d", len(migrator.Migrations))
		}
	})

	t.Run("isVersionExists method", func(t *testing.T) {
		migrator := &BaseMigrator{
			Migrations: []Migration{
				{Version: "20241015_test1"},
				{Version: "20241015_test2"},
			},
		}

		if !migrator.isVersionExists("20241015_test1") {
			t.Error("Expected version '20241015_test1' to exist")
		}

		if !migrator.isVersionExists("20241015_test2") {
			t.Error("Expected version '20241015_test2' to exist")
		}

		if migrator.isVersionExists("nonexistent") {
			t.Error("Expected version 'nonexistent' to not exist")
		}
	})

	t.Run("Up method with nonexistent version", func(t *testing.T) {
		migrator := &BaseMigrator{
			Migrations: []Migration{
				{Version: "20241015_test1"},
			},
		}

		err := migrator.Up(context.Background(), "nonexistent")
		if err == nil {
			t.Error("Expected error for nonexistent version")
		}

		if err.Error() != "version nonexistent does not exist" {
			t.Errorf("Expected specific error message, got: %s", err.Error())
		}
	})

	t.Run("Down method with empty version", func(t *testing.T) {
		migrator := &BaseMigrator{}

		err := migrator.Down(context.Background(), "")
		if err == nil {
			t.Error("Expected error for empty version")
		}

		if err.Error() != "a version is required" {
			t.Errorf("Expected specific error message, got: %s", err.Error())
		}
	})

	t.Run("Down method with nonexistent version", func(t *testing.T) {
		migrator := &BaseMigrator{
			Migrations: []Migration{
				{Version: "20241015_test1"},
			},
		}

		err := migrator.Down(context.Background(), "nonexistent")
		if err == nil {
			t.Error("Expected error for nonexistent version")
		}

		if err.Error() != "version nonexistent does not exist" {
			t.Errorf("Expected specific error message, got: %s", err.Error())
		}
	})

	t.Run("Unimplemented methods return errors", func(t *testing.T) {
		migrator := &BaseMigrator{}
		ctx := context.Background()

		// Test Connect
		err := migrator.Connect(ctx)
		if err == nil {
			t.Error("Expected Connect to return 'Not implemented' error")
		}

		// Test GetSchema
		_, err = migrator.GetSchema(ctx, "test")
		if err == nil {
			t.Error("Expected GetSchema to return 'Not implemented' error")
		}

		// Test ApplyMigration
		err = migrator.ApplyMigration(ctx, Migration{})
		if err == nil {
			t.Error("Expected ApplyMigration to return 'Not implemented' error")
		}

		// Test RevertMigration
		err = migrator.RevertMigration(ctx, Migration{})
		if err == nil {
			t.Error("Expected RevertMigration to return 'Not implemented' error")
		}
	})
}
