package config

import (
	"fmt"
	"os"
)

// PostgresConfig contain postgres database configuration
type PostgresConfig struct {
	DbName       string `required:"true" default:"typical-rest-server"`
	User         string `required:"true" default:"root"`
	Password     string `required:"true" default:"root"`
	Host         string `default:"localhost"`
	Port         int    `default:"5432"`
	MigrationSrc string `default:"file://scripts/migration"`
}

// DataSource return connection string
func (c *PostgresConfig) DataSource() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DB_USERNAME"), os.Getenv("DB_PASS"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))
}

// AdminDataSource return connection string for adminitration script
func (c *PostgresConfig) AdminDataSource() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DB_USERNAME"), os.Getenv("DB_PASS"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), "template1")
}

// DatabaseName return the database name
func (c *PostgresConfig) DatabaseName() string {
	return os.Getenv("DB_NAME")
}

// DriverName return the driver name
func (c *PostgresConfig) DriverName() string {
	return "postgres"
}

// MigrationSource return the migration source
func (c *PostgresConfig) MigrationSource() string {
	return os.Getenv("DB_MIGRATION_SRC")
}
