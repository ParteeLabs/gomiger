// Package core provides the main migration functionality for gomiger
package core

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// GomigerConfig is the global migration module configuration.
type GomigerConfig struct {
	// Path to the migration root folder.
	Path string `yaml:"path"`
	// The package name of the migrator & migrations.
	// Default by the folder name of the path.
	PkgName string `yaml:"pkg_name"`
	// Database connection string.
	URI string `yaml:"uri"`
	// The path to the table / collection schema store.
	SchemaStore string `yaml:"schema_store"`
}

var (
	defaultPath   = "./migrations"
	defaultRcPath = "./gomiger.rc.yaml"
)

// ParseYAML parse the RC file in YAML format.
func (rc *GomigerConfig) ParseYAML(path string) error {
	if path == "" {
		path = defaultRcPath
	}
	data, err := os.ReadFile(path) //nolint:gosec // Path is validated by caller
	if err != nil {
		return fmt.Errorf("cannot read the gomiger.rc file: %w", err)
	}
	err = yaml.Unmarshal(data, rc)
	if err != nil {
		return fmt.Errorf("cannot parse the gomiger.rc file: %w", err)
	}
	return nil
}

// PopulateAndValidate populate data and validate it.
func (rc *GomigerConfig) PopulateAndValidate() error {
	/// options validate & populate
	if rc.Path == "" {
		rc.Path = defaultPath
	}
	absPath, err := filepath.Abs(rc.Path)
	if err != nil {
		return fmt.Errorf("migration root path is not valid: %w", err)
	}
	rc.Path = absPath
	if rc.PkgName == "" {
		rc.PkgName = filepath.Dir(rc.Path)
	}
	return nil
}

// GetGomigerRC returns the global migration module configuration
func GetGomigerRC(rcPath string) (*GomigerConfig, error) {
	rc := &GomigerConfig{}
	if err := rc.ParseYAML(rcPath); err != nil {
		return nil, err
	}
	rc.URI = os.Getenv("GOMIGER_URI")
	if err := rc.PopulateAndValidate(); err != nil {
		return nil, err
	}
	return rc, nil
}
