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
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("Cannot read the gomiger.rc file: %w", err)
	}
	return yaml.Unmarshal(data, rc)
}

// PopulateAndValidate populate data and validate it.
func (rc *GomigerConfig) PopulateAndValidate() error {
	/// options validate & populate
	if rc.Path == "" {
		rc.Path = defaultPath
	}
	absPath, err := filepath.Abs(rc.Path)
	if err != nil {
		return fmt.Errorf("Migration root path is not valid: %w", err)
	}
	if rootDir, _ := filepath.Abs("./"); rootDir != absPath {
		rc.Path = defaultPath
	}
	if rc.PkgName == "" {
		rc.PkgName = filepath.Dir(rc.Path)
	}
	return nil
}

// GetGomigerRC returns the global migration module configuration
func GetGomigerRC(rcPath string) (*GomigerConfig, error) {
	if rcPath == "" {
		rcPath = defaultRcPath
	}
	rc := &GomigerConfig{}
	if err := rc.ParseYAML(rcPath); err != nil {
		return nil, err
	}
	if err := rc.PopulateAndValidate(); err != nil {
		return nil, err
	}
	return rc, nil
}
