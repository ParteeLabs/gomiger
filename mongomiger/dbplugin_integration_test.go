package mongomiger

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/ParteeLabs/gomiger/core"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type MongomigerTestSuite struct {
	suite.Suite
	config     *core.GomigerConfig
	mongomiger *Mongomiger
	ctx        context.Context
	cancel     context.CancelFunc
}

func (s *MongomigerTestSuite) SetupTest() {
	s.config = &core.GomigerConfig{
		URI:         "mongodb://localhost:27017/mongomiger_test",
		SchemaStore: "schema_migrations",
	}
	s.ctx, s.cancel = context.WithTimeout(context.Background(), 5*time.Second)
	s.mongomiger = NewMongomiger(s.config)
	err := s.mongomiger.Connect(s.ctx)
	s.Require().NoError(err)
}

func (s *MongomigerTestSuite) TearDownTest() {
	if s.mongomiger != nil && s.mongomiger.Db != nil {
		err := s.mongomiger.Db.Drop(s.ctx)
		s.Require().NoError(err)
	}
	if s.cancel != nil {
		s.cancel()
	}
}

func (s *MongomigerTestSuite) TestMongomiger_Connect() {
	s.Require().NotNil(s.mongomiger.Client)
	s.Require().NotNil(s.mongomiger.Db)
	s.Require().NotNil(s.mongomiger.schemaCollection)
}

func (s *MongomigerTestSuite) TestMongomiger_Connect_InvalidURI() {
	invalidConfig := &core.GomigerConfig{
		URI:         "invalid://uri",
		SchemaStore: "schema_migrations",
	}
	mongomiger := NewMongomiger(invalidConfig)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err := mongomiger.Connect(ctx)
	s.Require().Error(err)
}

func (s *MongomigerTestSuite) TestMongomiger_Connect_Unavailable() {
	unavailableConfig := &core.GomigerConfig{
		URI:         "mongodb://localhost:27099/unavailable_db", // Wrong port
		SchemaStore: "schema_migrations",
	}
	mongomiger := NewMongomiger(unavailableConfig)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err := mongomiger.Connect(ctx)
	s.Require().Error(err)
}

func (s *MongomigerTestSuite) TestMongomiger_GetSchema_NotFound() {
	// Seed the database.
	_, err := s.mongomiger.schemaCollection.InsertOne(s.ctx, bson.M{"version": "1.0.0", "status": "applied"})
	s.Require().NoError(err)
	// Try to get a non-existent schema.
	_, err = s.mongomiger.GetSchema(s.ctx, "2.0.0")
	s.Require().Error(err)
}

func (s *MongomigerTestSuite) TestMongomiger_GetSchema_Found() {
	// Seed the database.
	expectedSchema := &core.Schema{
		Version: "1.0.0",
		Status:  core.Applied,
	}
	_, err := s.mongomiger.schemaCollection.InsertOne(s.ctx, bson.M{"version": expectedSchema.Version, "status": expectedSchema.Status})
	s.Require().NoError(err)
	// Try to get the existing schema.
	schema, err := s.mongomiger.GetSchema(s.ctx, "1.0.0")
	s.Require().NoError(err)
	s.Require().Equal(expectedSchema.Version, schema.Version)
	s.Require().Equal(expectedSchema.Status, schema.Status)
}

func (s *MongomigerTestSuite) TestMongomiger_ApplyMigration_UpdateStatusError() {
	// Use a separate instance to avoid affecting suite state
	tempMongomiger := NewMongomiger(s.config)
	tempCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := tempMongomiger.Connect(tempCtx)
	s.Require().NoError(err)
	defer tempMongomiger.Db.Drop(tempCtx)

	// Create a migration that will succeed but then fail to update status
	migration := core.Migration{
		Version: "1.0.0",
		Up: func(ctx context.Context) error {
			// Disconnect after migration runs but before status update
			tempMongomiger.Client.Disconnect(tempCtx)
			return nil
		},
	}

	err = tempMongomiger.ApplyMigration(tempCtx, migration)
	s.Require().Error(err)
	s.ErrorContains(err, "failed to update schema status")
}

func (s *MongomigerTestSuite) TestMongomiger_ApplyMigration_DuplicateVersion() {
	schema := &core.Schema{
		Version: "1.0.0",
		Status:  core.Applied,
	}
	_, err := s.mongomiger.schemaCollection.InsertOne(s.ctx, bson.M{"version": schema.Version, "status": schema.Status})
	s.Require().NoError(err)
	// Try to apply a migration with the same version.
	migration := core.Migration{
		Version: "1.0.0",
		Up:      func(ctx context.Context) error { return nil },
	}
	err = s.mongomiger.ApplyMigration(s.ctx, migration)
	s.Require().Error(err)
	s.ErrorContains(err, "failed to apply migration at version")
}

func (s *MongomigerTestSuite) TestMongomiger_ApplyMigration_FailureMarksDirty() {
	migration := core.Migration{
		Version: "1.0.0",
		Up:      func(ctx context.Context) error { return fmt.Errorf("migration failed") },
	}
	err := s.mongomiger.ApplyMigration(s.ctx, migration)
	s.Require().Error(err)
	s.ErrorContains(err, "failed to apply migration")
	// Verify that the schema status is set to Dirty.
	schema, err := s.mongomiger.GetSchema(s.ctx, "1.0.0")
	s.Require().NoError(err)
	s.Require().Equal(core.Dirty, schema.Status)
}

func (s *MongomigerTestSuite) TestMongomiger_ApplyMigration_UpdateSchemaToDirtyError() {
	// Use a separate instance to avoid affecting suite state
	tempMongomiger := NewMongomiger(s.config)
	tempCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := tempMongomiger.Connect(tempCtx)
	s.Require().NoError(err)
	defer tempMongomiger.Db.Drop(tempCtx)

	migration := core.Migration{
		Version: "1.0.0",
		Up: func(ctx context.Context) error {
			// Disconnect after migration runs but before update
			tempMongomiger.Client.Disconnect(tempCtx)
			return fmt.Errorf("migration failed")
		},
	}

	err = tempMongomiger.ApplyMigration(tempCtx, migration)
	s.Require().Error(err)
	s.ErrorContains(err, "failed to update schema status")
}

func (s *MongomigerTestSuite) TestMongomiger_ApplyMigration_Success() {
	migration := core.Migration{
		Version: "1.0.0",
		Up:      func(ctx context.Context) error { return nil },
	}
	err := s.mongomiger.ApplyMigration(s.ctx, migration)
	s.Require().NoError(err)
	// Verify that the schema status is set to Applied.
	schema, err := s.mongomiger.GetSchema(s.ctx, migration.Version)
	s.Require().NoError(err)
	s.Require().Equal(core.Applied, schema.Status)
}

func (s *MongomigerTestSuite) TestMongomiger_RevertMigration_FailureMarksDirty() {
	schema := &core.Schema{
		Version:   "1.0.0",
		Status:    core.Applied,
		Timestamp: time.Now(),
	}
	_, err := s.mongomiger.schemaCollection.InsertOne(s.ctx, schema)
	s.Require().NoError(err)

	migration := core.Migration{
		Version: "1.0.0",
		Down:    func(ctx context.Context) error { return fmt.Errorf("revert failed") },
	}

	err = s.mongomiger.RevertMigration(s.ctx, migration)
	s.Require().Error(err)
	s.ErrorContains(err, "failed to revert migration")
	// Verify that the schema status is set to Dirty.
	resultSchema, err := s.mongomiger.GetSchema(s.ctx, "1.0.0")
	s.Require().NoError(err)
	s.Require().Equal(core.Dirty, resultSchema.Status)
}

func (s *MongomigerTestSuite) TestMongomiger_RevertMigration_UpdateSchemaToDirtyError() {
	// Use a separate instance to avoid affecting suite state
	tempMongomiger := NewMongomiger(s.config)
	tempCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := tempMongomiger.Connect(tempCtx)
	s.Require().NoError(err)
	defer tempMongomiger.Db.Drop(tempCtx)

	schema := &core.Schema{
		Version:   "1.0.0",
		Status:    core.Applied,
		Timestamp: time.Now(),
	}
	_, err = tempMongomiger.schemaCollection.InsertOne(tempCtx, schema)
	s.Require().NoError(err)

	migration := core.Migration{
		Version: "1.0.0",
		Down: func(ctx context.Context) error {
			// Disconnect after migration runs but before update
			tempMongomiger.Client.Disconnect(tempCtx)
			return fmt.Errorf("migration failed")
		},
	}

	err = tempMongomiger.RevertMigration(tempCtx, migration)
	s.Require().Error(err)
	s.ErrorContains(err, "failed to update schema status")
}

func (s *MongomigerTestSuite) TestMongomiger_RevertMigration_DeleteSchemaError() {
	// Use a separate instance to avoid affecting suite state
	tempMongomiger := NewMongomiger(s.config)
	tempCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := tempMongomiger.Connect(tempCtx)
	s.Require().NoError(err)
	defer tempMongomiger.Db.Drop(tempCtx)

	schema := &core.Schema{
		Version:   "1.0.0",
		Status:    core.Applied,
		Timestamp: time.Now(),
	}
	_, err = tempMongomiger.schemaCollection.InsertOne(tempCtx, schema)
	s.Require().NoError(err)

	migration := core.Migration{
		Version: "1.0.0",
		Down: func(ctx context.Context) error {
			// Disconnect after migration runs but before delete
			tempMongomiger.Client.Disconnect(tempCtx)
			return nil
		},
	}

	err = tempMongomiger.RevertMigration(tempCtx, migration)
	s.Require().Error(err)
	s.ErrorContains(err, "failed to delete schema at version")
}

func (s *MongomigerTestSuite) TestMongomiger_RevertMigration_Success() {
	// First, apply a migration to have something to revert.
	migration := core.Migration{
		Version: "1.0.0",
		Up:      func(ctx context.Context) error { return nil },
	}
	err := s.mongomiger.ApplyMigration(s.ctx, migration)
	s.Require().NoError(err)

	// Now, revert the migration.
	migration.Down = func(ctx context.Context) error { return nil }
	err = s.mongomiger.RevertMigration(s.ctx, migration)
	s.Require().NoError(err)

	// Verify that the schema is deleted.
	_, err = s.mongomiger.GetSchema(s.ctx, migration.Version)
	s.Require().Error(err)
}

func TestMongomigerTestSuite(t *testing.T) {
	suite.Run(t, new(MongomigerTestSuite))
}
