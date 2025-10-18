# Plugin Development Guide

This guide explains how to create custom database plugins for Gomiger, the Go migration framework.

## Overview

Gomiger uses a plugin-based architecture that allows you to implement database-specific migration logic while leveraging the core framework's flow control and utilities. Each plugin must implement the `core.Gomiger` interface.

## Core Interface

All database plugins must implement the `core.Gomiger` interface:

```go
type Gomiger interface {
    Up(ctx context.Context, toVersion string) error
    Down(ctx context.Context, atVersion string) error
    Connect(ctx context.Context) error
    GetSchema(ctx context.Context, version string) (*Schema, error)
    ApplyMigration(ctx context.Context, mi Migration) error
    RevertMigration(ctx context.Context, mi Migration) error
}
```

## Plugin Structure

### 1. Create the Plugin Struct

Your plugin should embed `*core.BaseMigrator` to inherit base functionality and add database-specific fields:

```go
package yourdbplugin

import (
    "context"
    "github.com/ParteeLabs/gomiger/core"
    // Your database driver imports
)

// YourDbPlugin implements core.Gomiger for YourDatabase.
type YourDbPlugin struct {
    *core.BaseMigrator
    uri              string
    client           *YourDBClient  // Your database client
    database         *YourDatabase  // Your database instance
    schemaStore      string
    schemaCollection *YourCollection // Schema tracking collection/table
}
```

### 2. Constructor Function

Create a constructor that initializes your plugin:

```go
// NewYourDbPlugin creates a new YourDbPlugin.
func NewYourDbPlugin(cfg *core.GomigerConfig) *YourDbPlugin {
    return &YourDbPlugin{
        BaseMigrator: &core.BaseMigrator{
            Migrations: []core.Migration{},
        },
        uri:         cfg.URI,
        schemaStore: cfg.SchemaStore,
    }
}
```

### 3. Implement Required Methods

#### Connect Method

Establish connection to your database:

```go
// Connect implements core.Gomiger.
func (p *YourDbPlugin) Connect(ctx context.Context) error {
    // Parse connection string
    // Create database client
    // Initialize database and schema storage
    // Example:
    client, err := yourdb.Connect(p.uri)
    if err != nil {
        return err
    }
    p.client = client
    p.database = client.Database("your_db_name")
    p.schemaCollection = p.database.Collection(p.schemaStore)
    return nil
}
```

#### GetSchema Method

Retrieve schema information for a specific version:

```go
// GetSchema implements core.Gomiger.
func (p *YourDbPlugin) GetSchema(ctx context.Context, version string) (*core.Schema, error) {
    var schema *core.Schema
    err := p.schemaCollection.FindOne(ctx, yourdb.Filter{"version": version}).Decode(&schema)
    if err != nil {
        return nil, fmt.Errorf("failed to get schema: %w", err)
    }
    return schema, nil
}
```

#### ApplyMigration Method

Execute a migration and track its status:

```go
// ApplyMigration implements core.Gomiger.
func (p *YourDbPlugin) ApplyMigration(ctx context.Context, mi core.Migration) error {
    // Mark migration as in progress
    schema := &core.Schema{
        Version:   mi.Version,
        Timestamp: time.Now(),
        Status:    core.InProgress,
    }

    if _, err := p.schemaCollection.InsertOne(ctx, schema); err != nil {
        return fmt.Errorf("failed to mark migration as in progress: %w", err)
    }

    // Execute the migration
    if err := mi.Up(ctx); err != nil {
        // Mark as dirty on failure
        if updateErr := p.updateSchemaStatus(ctx, mi, core.Dirty); updateErr != nil {
            return updateErr
        }
        return fmt.Errorf("failed to apply migration %s: %w", mi.Version, err)
    }

    // Mark as applied on success
    if err := p.updateSchemaStatus(ctx, mi, core.Applied); err != nil {
        return err
    }

    return nil
}
```

#### RevertMigration Method

Revert a migration and clean up schema tracking:

```go
// RevertMigration implements core.Gomiger.
func (p *YourDbPlugin) RevertMigration(ctx context.Context, mi core.Migration) error {
    // Execute the down migration
    if err := mi.Down(ctx); err != nil {
        // Mark as dirty on failure
        if updateErr := p.updateSchemaStatus(ctx, mi, core.Dirty); updateErr != nil {
            return updateErr
        }
        return fmt.Errorf("failed to revert migration %s: %w", mi.Version, err)
    }

    // Remove schema record
    if _, err := p.schemaCollection.DeleteOne(ctx, yourdb.Filter{"version": mi.Version}); err != nil {
        return fmt.Errorf("failed to delete schema record: %w", err)
    }

    return nil
}
```

### 4. Helper Methods

Add helper methods for common operations:

```go
func (p *YourDbPlugin) updateSchemaStatus(ctx context.Context, mi core.Migration, status core.SchemaStatus) error {
    filter := yourdb.Filter{"version": mi.Version}
    update := yourdb.Update{"$set": yourdb.Document{"status": status, "timestamp": time.Now()}}

    if _, err := p.schemaCollection.UpdateOne(ctx, filter, update); err != nil {
        return fmt.Errorf("failed to update schema status for version %s: %w", mi.Version, err)
    }
    return nil
}
```

### 5. Interface Compliance Check

**Important**: Always add an interface compliance check at the end of your file:

```go
// Interface check - ensures YourDbPlugin implements core.Gomiger
var _ core.Gomiger = (*YourDbPlugin)(nil)
```

This compile-time check ensures your plugin correctly implements all required methods of the `core.Gomiger` interface.

## Complete Example

Here's a minimal but complete plugin structure:

```go
package yourdbplugin

import (
    "context"
    "fmt"
    "time"

    "github.com/ParteeLabs/gomiger/core"
    // Your database imports
)

type YourDbPlugin struct {
    *core.BaseMigrator
    uri         string
    client      *YourDBClient
    database    *YourDatabase
    schemaStore string
}

func NewYourDbPlugin(cfg *core.GomigerConfig) *YourDbPlugin {
    return &YourDbPlugin{
        BaseMigrator: &core.BaseMigrator{
            Migrations: []core.Migration{},
        },
        uri:         cfg.URI,
        schemaStore: cfg.SchemaStore,
    }
}

func (p *YourDbPlugin) Connect(ctx context.Context) error {
    // Implementation
    return nil
}

func (p *YourDbPlugin) GetSchema(ctx context.Context, version string) (*core.Schema, error) {
    // Implementation
    return nil, nil
}

func (p *YourDbPlugin) ApplyMigration(ctx context.Context, mi core.Migration) error {
    // Implementation
    return nil
}

func (p *YourDbPlugin) RevertMigration(ctx context.Context, mi core.Migration) error {
    // Implementation
    return nil
}

// Interface compliance check
var _ core.Gomiger = (*YourDbPlugin)(nil)
```

## Best Practices

1. **Error Handling**: Always wrap errors with context using `fmt.Errorf`
2. **Schema Tracking**: Maintain accurate migration status in your schema store
3. **Transactions**: Use database transactions when possible to ensure atomicity
4. **Logging**: Add appropriate logging for debugging and monitoring
5. **Testing**: Write comprehensive tests for your plugin
6. **Documentation**: Document any database-specific configuration requirements

## Schema States

Your plugin must handle these schema states correctly:

- `core.InProgress`: Migration is currently running
- `core.Applied`: Migration completed successfully
- `core.Dirty`: Migration failed and needs manual intervention

## Configuration

Ensure your plugin works with the standard `core.GomigerConfig` struct:

```go
type GomigerConfig struct {
    URI         string // Database connection string
    SchemaStore string // Schema tracking collection/table name
    // Other configuration fields
}
```

## Testing Your Plugin

Create comprehensive tests covering:

- Connection establishment
- Schema retrieval and updates
- Migration application and reversion
- Error scenarios and dirty state handling
- Interface compliance

Remember to always include the interface compliance check to catch implementation issues at compile time!
