# gomiger

A Golang framework for data migration. Focused on programmability and flexibility.

## installation

**Step 1: Install gomiger and deps.**

```bash
go get github.com/ParteeLabs/gomiger/core
go get github.com/urfave/cli/v3 # Optional, if you want to use the CLI
```

**Step 2: Add the `gomiger.rc.yaml` file to the root of your project.**

```yaml
path: './migrations' # Path to the migrations folder
pkg_name: 'mgr' # Package name
uri: '' # Database URI
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

**Step 5: Chose the DB plugin you want to use**
Gomiger supports multiple DBs by implementing the `DbPlugin` interface. You can add your plugin to the `migrator.mg.go` file. Use our plugin in this monorepo or build your own plugin.

### MongoDB Plugin

First, install the DB plugin

```bash
go get github.com/ParteeLabs/gomiger/mongomiger
```

Then add the DB plugin to your generate `migrator.mg.go` file.

```diff
func NewMigrator(config *core.GomigerConfig) core.Gomiger {
	m := &Migrator{
		Config: config,
	}
+	m.DB = mongomiger.NewMongomiger(config)
...
}
```

### Postgres Plugin: Coming soon...

## Usage

### Commands

**Generate a new migration file.**

```bash
go run cli.go new migration_name
```

**Run migrations up.**

```bash
go run cli.go up # To run all migrations
go run cli.go up version # To stop at a specific version
```

**Run migrations down.**

```bash
go run cli.go down version
```
