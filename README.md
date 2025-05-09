# gomiger

A Golang framework for data migration.

# instruction

## install

```bash
go get github.com/ParteeLabs/gomiger/core
```

Add the `gomiger.rc.yaml` file to the root of your project.

```yaml
path: './migrations' # Path to the migrations folder
pkg_name: 'mgr' # Package name
uri: '' # Database URI
schema_store: '' # Database schema store
```

## usage

Initialize the source code.

```bash
go run github.com/ParteeLabs/gomiger/core/cmd/gomiger-init
```

Add a migration.

```bash
go run your-migration-root-path add migration_name
```

Run migrations up.

```bash
go run your-migration-root-path up
```

Run migrations down.

```bash
go run your-migration-root-path down version
```
