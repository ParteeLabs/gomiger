# Getting Started with Gomiger

This guide will walk you through setting up and using Gomiger for database migrations.

## Prerequisites

- Go 1.21 or later
- A supported database (MongoDB, PostgreSQL coming soon)
- Basic understanding of Go programming

## Installation

### 1. Install the CLI Tool

```bash
go install github.com/ParteeLabs/gomiger/core/cmd/gomiger-init@latest
```

### 2. Install Core Library

```bash
go get github.com/ParteeLabs/gomiger/core
```

### 3. Install Database Plugin

For MongoDB:

```bash
go get github.com/ParteeLabs/gomiger/mongomiger
```

## Project Setup

### 1. Initialize Your Go Project

```bash
mkdir my-app && cd my-app
go mod init github.com/my-app
```

### 2. Create Configuration File

Create `gomiger.rc.yaml` in your project root:

```yaml
path: './migrations'
pkg_name: 'migrations'
schema_store: 'schema_migrations'
```

### 3. Initialize Migration Structure

```bash
gomiger-init
```

This creates:

```
migrations/
├── cli.mg.go
└── migrator.mg.go
```

### 4. Setup Database Plugin

Edit `migrations/migrator.mg.go`:

```go
package migrations

import (
	"github.com/ParteeLabs/gomiger/core"
	"github.com/ParteeLabs/gomiger/mongomiger"
)

type Migrator struct {
	*mongomiger.Mongomiger
	Config *core.GomigerConfig
}

func NewMigrator(config *core.GomigerConfig) core.Gomiger {
	return &Migrator{
		Mongomiger: mongomiger.NewMongomiger(config),
		Config:     config,
	}
}
```

### 5. Create CLI Entry Point

Create `main.go`:

```go
package main

import (
	"github.com/my-app/migrations"
)

func main() {
	migrations.Run()
}
```

## Creating Your First Migration

### 1. Generate Migration File

```bash
go run main.go new create_users_table
```

This creates `migrations/TIMESTAMP_create_users_table.mg.go`:

```go
package migrations

import "context"

func (m *Migrator) Migration_TIMESTAMP_create_users_table_Up(ctx context.Context) error {
	// Your migration up code here
	return nil
}

func (m *Migrator) Migration_TIMESTAMP_create_users_table_Down(ctx context.Context) error {
	// Your migration down code here
	return nil
}

// AUTO GENERATED, DO NOT MODIFY!
func (m *Migrator) Migration_TIMESTAMP_create_users_table_Version() string {
	return "TIMESTAMP"
}
```

### 2. Implement Migration Logic

For MongoDB:

```go
func (m *Migrator) Migration_TIMESTAMP_create_users_table_Up(ctx context.Context) error {
	// Create users collection with validation
	validator := bson.M{
		"$jsonSchema": bson.M{
			"bsonType": "object",
			"required": []string{"email", "created_at"},
			"properties": bson.M{
				"email": bson.M{
					"bsonType":    "string",
					"description": "must be a string and is required",
				},
				"name": bson.M{
					"bsonType":    "string",
					"description": "user's full name",
				},
				"created_at": bson.M{
					"bsonType":    "date",
					"description": "must be a date and is required",
				},
			},
		},
	}

	opts := options.CreateCollection().SetValidator(validator)
	return m.Db.CreateCollection(ctx, "users", opts)
}

func (m *Migrator) Migration_TIMESTAMP_create_users_table_Down(ctx context.Context) error {
	return m.Db.Collection("users").Drop(ctx)
}
```

## Running Migrations

### 1. Set Database Connection

```bash
export GOMIGER_URI="mongodb://localhost:27017/myapp"
```

### 2. Run Migrations Up

```bash
# Run all pending migrations
go run main.go up

# Run up to specific version
go run main.go up 202410151200
```

### 3. Run Migrations Down

```bash
# Rollback to specific version
go run main.go down 202410151200
```

### 4. Check Migration Status

```bash
# List migration status (if implemented in your CLI)
go run main.go status
```

## Best Practices

### 1. Migration Naming

- Use descriptive names: `create_users_table`, `add_email_index`
- Follow consistent naming patterns
- Include the purpose in the name

### 2. Migration Content

- Keep migrations focused on a single change
- Always implement both Up and Down methods
- Test migrations thoroughly before deploying
- Use transactions when supported by your database

### 3. Error Handling

```go
func (m *Migrator) Migration_TIMESTAMP_example_Up(ctx context.Context) error {
	collection := m.Db.Collection("users")

	// Use proper error handling
	if err := collection.CreateIndex(ctx, indexModel); err != nil {
		return fmt.Errorf("failed to create index: %w", err)
	}

	return nil
}
```

### 4. Testing Migrations

Create test files for your migrations:

```go
// migrations_test.go
func TestCreateUsersTable(t *testing.T) {
	// Setup test database
	config := &core.GomigerConfig{
		URI:         "mongodb://localhost:27017/test_db",
		SchemaStore: "test_migrations",
	}

	migrator := NewMigrator(config)
	defer migrator.Client.Database("test_db").Drop(context.Background())

	// Test up migration
	err := migrator.Connect(context.Background())
	require.NoError(t, err)

	// Run specific migration
	migration := core.Migration{
		Version: "202410151200",
		Up:      migrator.Migration_202410151200_create_users_table_Up,
		Down:    migrator.Migration_202410151200_create_users_table_Down,
	}

	err = migrator.ApplyMigration(context.Background(), migration)
	assert.NoError(t, err)

	// Verify changes
	collections, err := migrator.Db.ListCollectionNames(context.Background(), bson.D{})
	assert.NoError(t, err)
	assert.Contains(t, collections, "users")
}
```

## Advanced Usage

### Environment-Specific Configurations

```yaml
# gomiger.rc.yaml
path: './migrations'
pkg_name: 'migrations'
schema_store: 'schema_migrations'

# Development settings
development:
  timeout: 10s

# Production settings
production:
  timeout: 60s
  max_retries: 3
```

### Custom Migration Templates

You can customize the generated migration template by modifying the generator in the core package.

## Troubleshooting

### Common Issues

1. **Migration fails to connect**

   - Check your `GOMIGER_URI` environment variable
   - Ensure the database is running and accessible

2. **Migration already applied**

   - Check your schema store collection/table
   - Use `down` command to rollback if needed

3. **Import errors**
   - Run `go mod tidy` to resolve dependencies
   - Check that all required packages are installed

### Getting Help

- Open an issue on [GitHub](https://github.com/ParteeLabs/gomiger/issues)
- Join our [Discussions](https://github.com/ParteeLabs/gomiger/discussions)

## Next Steps

- Learn about [Plugin Development](plugin-development.md)
- Read our [Best Practices Guide](best-practices.md)
- Explore [Advanced Examples](../examples/)
- Join our community on [GitHub Discussions](https://github.com/ParteeLabs/gomiger/discussions)
