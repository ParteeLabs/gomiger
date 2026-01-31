package core

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
)

type BaseMigratorTestSuite struct {
	suite.Suite
	migrator *BaseMigrator
}

func (s *BaseMigratorTestSuite) SetupTest() {
	s.migrator = &BaseMigrator{
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

func TestBaseMigratorTestSuite(t *testing.T) {
	suite.Run(t, new(BaseMigratorTestSuite))
}
