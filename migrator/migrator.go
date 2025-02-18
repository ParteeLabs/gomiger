package migrator

import "time"

// DbType is the type of the database
type DbType string

var (
	Mongo DbType = "mongo"
)

// Config is the configuration for the database
type Config struct {
	Type        DbType `json:"type" yaml:"type" validate:"required,oneof=mongo"`
	SchemaStore string `json:"schema_store" yaml:"schema_store" validate:"required"`
	URI         string `json:"uri" yaml:"uri" validate:"required"`
	Path        string `json:"path" yaml:"path" validate:"required"`
}

// Schema is a log for migrations
type Schema struct {
	Version   string    `json:"version" bson:"version" validate:"required"`
	Timestamp time.Time `json:"timestamp" bson:"timestamp" validate:"required"`
}
