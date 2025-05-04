# gomiger

A Golang framework for data migration.

# instruction

## install

```bash
go get github.com/ParteeLabs/gomiger
```

Add the CLI script to Makefile.

```make
gomiger:
  go run github.com/ParteeLabs/gomiger/cmd/cli
```

Add the `gomiger.rc.yaml` file to the root of your project.

```yaml
path: './migrations' # Path to the migrations folder
pkg_name: 'mgr' # Package name
```

## usage

Initialize the source code.

```bash
make gomiger init
```

Add a migration.

```bash
make gomiger new migration_name
```
