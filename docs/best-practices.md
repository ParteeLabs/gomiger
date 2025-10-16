# Best Practices

## 1. Migration Naming

- Use descriptive names: `create_users_table`, `add_email_index`
- Follow consistent naming patterns
- Include the purpose in the name

## 2. Migration Content

- Keep migrations focused on a single change
- Always implement both Up and Down methods
- Test migrations thoroughly before deploying
- Use transactions when supported by your database

## 3. Error Handling

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

## 4. Testing Migrations

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
		Version: migrator.Migration_202410151200_create_users_table_Version(),
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
