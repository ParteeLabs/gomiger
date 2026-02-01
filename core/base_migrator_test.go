package core

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type MockAbstractMethods struct {
	mock.Mock
}

func (m *MockAbstractMethods) Connect(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockAbstractMethods) GetSchema(ctx context.Context, version string) (*Schema, error) {
	args := m.Called(ctx, version)
	schema, _ := args.Get(0).(*Schema)
	return schema, args.Error(1)
}

func (m *MockAbstractMethods) ApplyMigration(ctx context.Context, mi Migration) error {
	args := m.Called(ctx, mi)
	return args.Error(0)
}

func (m *MockAbstractMethods) RevertMigration(ctx context.Context, mi Migration) error {
	args := m.Called(ctx, mi)
	return args.Error(0)
}

type BaseMigratorTestSuite struct {
	suite.Suite
	migrator *BaseMigrator
}

func (s *BaseMigratorTestSuite) SetupTest() {
	s.migrator = &BaseMigrator{
		BaseMigratorAbstractMethods: &MockAbstractMethods{},
		Migrations: []Migration{
			{Version: "20240101_initial"},
			{Version: "20240201_add_users"},
			{Version: "20240301_add_orders"},
		},
	}
}

func (s *BaseMigratorTestSuite) TestIsVersionExists_ExistingVersion() {
	exists := s.migrator.isVersionExists("20240201_add_users")
	s.True(exists, "Expected version '20240201_add_users' to exist")
}

func (s *BaseMigratorTestSuite) TestIsVersionExists_NonexistentVersion() {
	exists := s.migrator.isVersionExists("20240401_nonexistent")
	s.False(exists, "Expected version '20240401_nonexistent' to not exist")
}

func (s *BaseMigratorTestSuite) TestIsVersionExists_EmptyVersion() {
	exists := s.migrator.isVersionExists("")
	s.False(exists, "Expected empty version to not exist")
}

func (s *BaseMigratorTestSuite) TestUp_NonexistentVersion() {
	err := s.migrator.Up(context.Background(), "20240401_nonexistent")
	s.Error(err)
	s.Contains(err.Error(), "version 20240401_nonexistent does not exist")
}

func (s *BaseMigratorTestSuite) TestUp_GetSchemaError() {
	mockMethods := s.migrator.BaseMigratorAbstractMethods.(*MockAbstractMethods)
	errSchemaNotfound := fmt.Errorf("schema not found")
	mockMethods.On("GetSchema", mock.Anything, mock.Anything).Return(nil, errSchemaNotfound).Once()

	err := s.migrator.Up(context.Background(), "20240201_add_users")
	s.Error(err)
	s.ErrorIs(err, errSchemaNotfound)
	mockMethods.AssertExpectations(s.T())
}

func (s *BaseMigratorTestSuite) TestUp_ApplyMigrationError() {
	mockMethods := s.migrator.BaseMigratorAbstractMethods.(*MockAbstractMethods)
	errApplyMigration := fmt.Errorf("apply migration failed")
	mockMethods.On("GetSchema", mock.Anything, mock.Anything).Return(&Schema{}, nil).Once()
	mockMethods.On("ApplyMigration", mock.Anything, mock.Anything).Return(errApplyMigration).Once()

	err := s.migrator.Up(context.Background(), "20240201_add_users")
	s.Error(err)
	s.ErrorIs(err, errApplyMigration)
	mockMethods.AssertExpectations(s.T())
}

func (s *BaseMigratorTestSuite) TestUp_SuccessfulMigration() {
	mockMethods := s.migrator.BaseMigratorAbstractMethods.(*MockAbstractMethods)
	mockMethods.On("GetSchema", mock.Anything, "20240301_add_orders").Return(&Schema{Status: Applied}, nil).Once()
	mockMethods.On("GetSchema", mock.Anything, mock.Anything).Return(&Schema{Status: InProgress}, nil).Times(2)
	mockMethods.On("ApplyMigration", mock.Anything, mock.Anything).Return(nil).Times(2)

	err := s.migrator.Up(context.Background(), "")
	s.NoError(err)
	mockMethods.AssertExpectations(s.T())
}

func (s *BaseMigratorTestSuite) TestUp_SuccessfulMigrationToVersion() {
	mockMethods := s.migrator.BaseMigratorAbstractMethods.(*MockAbstractMethods)
	mockMethods.On("GetSchema", mock.Anything, mock.Anything).Return(&Schema{Status: InProgress}, nil).Times(2)
	mockMethods.On("ApplyMigration", mock.Anything, mock.Anything).Return(nil).Times(2)

	err := s.migrator.Up(context.Background(), "20240201_add_users")
	s.NoError(err)
	mockMethods.AssertExpectations(s.T())
}

func (s *BaseMigratorTestSuite) TestDown_EmptyVersion() {
	err := s.migrator.Down(context.Background(), "")
	s.Error(err)
	s.Contains(err.Error(), "a version is required")
}

func (s *BaseMigratorTestSuite) TestDown_NonexistentVersion() {
	err := s.migrator.Down(context.Background(), "20240401_nonexistent")
	s.Error(err)
	s.Contains(err.Error(), "version 20240401_nonexistent does not exist")
}

func (s *BaseMigratorTestSuite) TestDown_GetSchemaError() {
	mockMethods := s.migrator.BaseMigratorAbstractMethods.(*MockAbstractMethods)
	errSchemaNotfound := fmt.Errorf("schema not found")
	mockMethods.On("GetSchema", mock.Anything, mock.Anything).Return(nil, errSchemaNotfound).Once()

	err := s.migrator.Down(context.Background(), "20240201_add_users")
	s.Error(err)
	s.ErrorIs(err, errSchemaNotfound)
	mockMethods.AssertExpectations(s.T())
}

func (s *BaseMigratorTestSuite) TestDown_RevertMigrationError() {
	mockMethods := s.migrator.BaseMigratorAbstractMethods.(*MockAbstractMethods)
	errRevertMigration := fmt.Errorf("revert migration failed")
	mockMethods.On("GetSchema", mock.Anything, mock.Anything).Return(&Schema{
		Status: Applied,
	}, nil).Once()
	mockMethods.On("RevertMigration", mock.Anything, mock.Anything).Return(errRevertMigration).Once()

	err := s.migrator.Down(context.Background(), "20240201_add_users")
	s.Error(err)
	s.ErrorIs(err, errRevertMigration)
	mockMethods.AssertExpectations(s.T())
}

func (s *BaseMigratorTestSuite) TestDown_SuccessfulReversion() {
	mockMethods := s.migrator.BaseMigratorAbstractMethods.(*MockAbstractMethods)
	mockMethods.On("GetSchema", mock.Anything, "20240301_add_orders").Return(&Schema{
		Status: InProgress,
	}, nil).Once()
	mockMethods.On("GetSchema", mock.Anything, mock.Anything).Return(&Schema{
		Status: Applied,
	}, nil).Times(2)
	mockMethods.On("RevertMigration", mock.Anything, mock.Anything).Return(nil).Times(2)

	err := s.migrator.Down(context.Background(), "20240101_initial")
	s.NoError(err)
	mockMethods.AssertExpectations(s.T())
}

func TestBaseMigratorTestSuite(t *testing.T) {
	suite.Run(t, new(BaseMigratorTestSuite))
}
