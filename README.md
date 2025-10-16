# 🚀 Gomiger

[![CI/CD Pipeline](https://github.com/ParteeLabs/gomiger/actions/workflows/ci.yml/badge.svg)](https://github.com/ParteeLabs/gomiger/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/ParteeLabs/gomiger/core)](https://goreportcard.com/report/github.com/ParteeLabs/gomiger/core)
[![Coverage Status](https://codecov.io/gh/ParteeLabs/gomiger/branch/main/graph/badge.svg)](https://codecov.io/gh/ParteeLabs/gomiger)
[![Go Reference](https://pkg.go.dev/badge/github.com/ParteeLabs/gomiger.svg)](https://pkg.go.dev/github.com/ParteeLabs/gomiger)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A powerful, type-safe Go migration framework with plugin architecture. Build reliable database migrations with the full power of Go.

## ✨ Features

- **🔌 Plugin Architecture**: Support multiple databases through plugins
- **🛡️ Type Safety**: Leverage Go's type system for reliable migrations
- **🎯 Programmable**: Write migrations in Go with full language features
- **⚡ Performance**: Fast execution with minimal overhead
- **🔧 Flexible**: Customize behavior through configuration
- **📊 Tracking**: Built-in migration state management

## 🎯 Why Gomiger?

Unlike traditional SQL-based migration tools, Gomiger lets you:

- Use Go's powerful standard library and ecosystem
- Write complex data transformations with loops, conditions, and functions
- Leverage existing Go libraries and business logic
- Maintain type safety throughout your migrations
- Test your migrations like regular Go code

## 🚀 Quick Start

### Prerequisites

- Go 1.21 or later
- Database (MongoDB, PostgreSQL, etc.)

### Installation

**Option 1: Go Install (Recommended)**

```bash
go install github.com/ParteeLabs/gomiger/core/cmd/gomiger-init@latest
```

**Option 2: Manual Installation**

```bash
go get github.com/ParteeLabs/gomiger/core
go get github.com/ParteeLabs/gomiger/mongomiger  # For MongoDB
go get github.com/urfave/cli/v3                 # For CLI support
```

**Step 2: Add the `gomiger.rc.yaml` file to the root of your project.**

```yaml
path: './migrations' # Path to the migrations folder
pkg_name: 'mgr' # Package name
schema_store: '' # Database schema store
```

**Step 3: Initialize the source code by running:**

```bash
go run github.com/ParteeLabs/gomiger/core/cmd/gomiger-init
```

The source code will be initialized in the `path` folder.

```plaintext
migrations/
├── cli.mg.go
└── migrator.mg.go
```

**Step 4: Add you CLI entry point.**

You can add any entry point (e.g. `cli.go` or `gomiger.go`) then import & call the `Run` function in `cli.mg.go`

```go
// cli.go
package main

import (
	mgr "test/migrations"
)

func main() {
	mgr.Run()
}
```

### Database Plugins

Gomiger supports multiple databases through plugins:

#### 🍃 MongoDB Plugin

```bash
go get github.com/ParteeLabs/gomiger/mongomiger
```

Add to your `migrator.mg.go`:

```diff
type Migrator struct {
-	core.BaseMigrator
+	*mongomiger.Mongomiger
	Config *core.GomigerConfig
}

func NewMigrator(config *core.GomigerConfig) core.Gomiger {
	return &Migrator{
+		Mongomiger: mongomiger.NewMongomiger(config),
		Config:     config,
	}
}
```

#### 🐘 PostgreSQL Plugin

_Coming Soon_ - We're working on PostgreSQL support!

#### 🔌 Custom Plugin

Want to add support for your database? Check our [Plugin Development Guide](docs/plugin-development.md).

## Usage

### Commands

**Generate a new migration file.**

```bash
go run cli.go new migration_name
```

**Run migrations up.**

```bash
export GOMIGER_URI="mongodb://localhost:27017"
go run cli.go up # To run all migrations
go run cli.go up version # To stop at a specific version
```

**Run migrations down.**

```bash
export GOMIGER_URI="mongodb://localhost:27017"
go run cli.go down version
```

## 📚 Examples

### Simple User Schema Migration

```go
// migrations/202410151200_add_users.mg.go
package mgr

import (
	"context"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func (m *Migrator) Migration_202410151200_add_users_Up(ctx context.Context) error {
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

func (m *Migrator) Migration_202410151200_add_users_Down(ctx context.Context) error {
	return m.Db.Collection("users").Drop(ctx)
}
```

### Complex Data Transformation

```go
// migrations/202410151300_migrate_user_format.mg.go
package mgr

import (
	"context"
	"time"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func (m *Migrator) Migration_202410151300_migrate_user_format_Up(ctx context.Context) error {
	collection := m.Db.Collection("users")

	// Find all users with old format
	cursor, err := collection.Find(ctx, bson.M{"full_name": bson.M{"$exists": true}})
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	// Transform each user
	for cursor.Next(ctx) {
		var user bson.M
		if err := cursor.Decode(&user); err != nil {
			return err
		}

		// Split full_name into first_name and last_name
		if fullName, ok := user["full_name"].(string); ok {
			names := strings.Fields(fullName)
			update := bson.M{
				"$set": bson.M{
					"first_name": names[0],
					"last_name":  "",
					"updated_at": time.Now(),
				},
				"$unset": bson.M{"full_name": ""},
			}

			if len(names) > 1 {
				update["$set"].(bson.M)["last_name"] = strings.Join(names[1:], " ")
			}

			_, err := collection.UpdateOne(ctx, bson.M{"_id": user["_id"]}, update)
			if err != nil {
				return err
			}
		}
	}

	return cursor.Err()
}
```

## 🏗️ Architecture

Gomiger follows a clean plugin architecture:

```
┌─────────────────┐    ┌──────────────────┐
│   Your App      │────│   Gomiger Core   │
└─────────────────┘    └──────────────────┘
                              │
                    ┌─────────┼─────────┐
                    │                   │
            ┌───────▼────┐    ┌─────────▼────┐
            │  MongoDB   │    │ PostgreSQL   │
            │  Plugin    │    │  Plugin      │
            └────────────┘    └──────────────┘
```

- **Core**: Migration engine and interfaces
- **Plugins**: Database-specific implementations
- **Generated Code**: Type-safe migration scaffolding
- **CLI**: Command-line interface for operations

## 🛠️ Advanced Configuration

### Environment Variables

```bash
export GOMIGER_URI="mongodb://localhost:27017/mydb"
```

### Configuration File Options

```yaml
# gomiger.rc.yaml
path: './migrations'
pkg_name: 'mgr'
schema_store: 'schema_migrations'
```

## 🧪 Testing Your Migrations

```go
// migrations_test.go
func TestUserMigration(t *testing.T) {
	config := &core.GomigerConfig{
		URI:         "mongodb://localhost:27017/test_db",
		SchemaStore: "test_schemas",
	}

	migrator := NewMigrator(config)
	ctx := context.Background()

	// Test up migration
	err := migrator.Up(ctx, "202410151200")
	assert.NoError(t, err)

	// Verify collection exists
	collections, err := migrator.Db.ListCollectionNames(ctx, bson.D{})
	assert.NoError(t, err)
	assert.Contains(t, collections, "users")

	// Test down migration
	err = migrator.Down(ctx, "202410151200")
	assert.NoError(t, err)
}
```

## 📖 Documentation

- [Getting Started Guide](docs/getting-started.md)
- [Plugin Development](docs/plugin-development.md)
- [Best Practices](docs/best-practices.md)
- [API Reference](https://pkg.go.dev/github.com/ParteeLabs/gomiger)
- [Examples Repository](examples/)

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

### Development Setup

```bash
git clone https://github.com/ParteeLabs/gomiger.git
cd gomiger
go work use ./core ./mongomiger ./examples
go test ./...
```

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Inspired by Sequelize Umzug
- Built with the amazing Go ecosystem
- Special thanks to all contributors

## 🔗 Links

- [GitHub Repository](https://github.com/ParteeLabs/gomiger)
- [Go Package Documentation](https://pkg.go.dev/github.com/ParteeLabs/gomiger)
- [Issue Tracker](https://github.com/ParteeLabs/gomiger/issues)
- [Discussions](https://github.com/ParteeLabs/gomiger/discussions)

---

<p align="center">
Made with ❤️ by <a href="https://github.com/ParteeLabs">ParteeLabs</a>
</p>
