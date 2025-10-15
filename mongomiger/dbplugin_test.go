package mongomiger

import (
	"testing"

	"github.com/ParteeLabs/gomiger/core"
)

func TestMongomiger(t *testing.T) {
	t.Run("NewMongomiger initialization", func(t *testing.T) {
		config := &core.GomigerConfig{
			URI:         "mongodb://localhost:27017",
			SchemaStore: "test_migrations",
		}

		mongomiger := NewMongomiger(config)

		if mongomiger == nil {
			t.Error("Expected Mongomiger to be created")
		}

		if mongomiger.uri != config.URI {
			t.Errorf("Expected URI '%s', got '%s'", config.URI, mongomiger.uri)
		}

		if mongomiger.schemaStore != config.SchemaStore {
			t.Errorf("Expected schemaStore '%s', got '%s'", config.SchemaStore, mongomiger.schemaStore)
		}
	})

	t.Run("Mongomiger implements core.Gomiger interface", func(t *testing.T) {
		config := &core.GomigerConfig{
			URI:         "mongodb://localhost:27017",
			SchemaStore: "test_migrations",
		}

		mongomiger := NewMongomiger(config)

		// Verify it implements the interface by assigning to interface type
		var _ core.Gomiger = mongomiger

		// Test that basic structure is correct
		if mongomiger.BaseMigrator == nil {
			t.Error("Expected BaseMigrator to be initialized")
		}

		// Test basic configuration
		if mongomiger.uri != config.URI {
			t.Errorf("Expected URI '%s', got '%s'", config.URI, mongomiger.uri)
		}

		if mongomiger.schemaStore != config.SchemaStore {
			t.Errorf("Expected schemaStore '%s', got '%s'", config.SchemaStore, mongomiger.schemaStore)
		}

		// Test that Client and schemaCollection are nil before Connect is called
		if mongomiger.Client != nil {
			t.Error("Expected Client to be nil before Connect is called")
		}

		if mongomiger.schemaCollection != nil {
			t.Error("Expected schemaCollection to be nil before Connect is called")
		}
	})
}

// Note: Integration tests with real MongoDB would require:
// 1. MongoDB test container or test instance
// 2. Connection setup and teardown
// 3. Test database cleanup
//
// For now, we're testing the structure and interface compliance.
// Full integration tests should be added when setting up CI/CD.

func TestMongomigerConfig(t *testing.T) {
	t.Run("Config validation", func(t *testing.T) {
		tests := []struct {
			name        string
			uri         string
			schemaStore string
			expectValid bool
		}{
			{
				name:        "Valid MongoDB URI",
				uri:         "mongodb://localhost:27017",
				schemaStore: "migrations",
				expectValid: true,
			},
			{
				name:        "Valid MongoDB URI with database",
				uri:         "mongodb://localhost:27017/mydb",
				schemaStore: "schema_migrations",
				expectValid: true,
			},
			{
				name:        "Empty schema store",
				uri:         "mongodb://localhost:27017",
				schemaStore: "",
				expectValid: true, // Should handle empty schema store gracefully
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				config := &core.GomigerConfig{
					URI:         tt.uri,
					SchemaStore: tt.schemaStore,
				}

				mongomiger := NewMongomiger(config)

				if mongomiger == nil && tt.expectValid {
					t.Error("Expected valid Mongomiger to be created")
				}

				if mongomiger != nil {
					if mongomiger.uri != tt.uri {
						t.Errorf("Expected URI '%s', got '%s'", tt.uri, mongomiger.uri)
					}

					if mongomiger.schemaStore != tt.schemaStore {
						t.Errorf("Expected schemaStore '%s', got '%s'", tt.schemaStore, mongomiger.schemaStore)
					}
				}
			})
		}
	})
}
