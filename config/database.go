package config

import (
	"os"
	"strings"
)

const dummyReaderDbUrl = "1d90c12e-5e44-4598-86bb-111f08cf42ad"

type database struct {
	// Create a reader-only DB-connection for the times when we use a multi node
	// e.g. Managed PostgreSQL instances like AWS RDS etc.
	Main   databaseConfig
	Reader databaseConfig
}

type databaseConfig struct {
	Url                      string // The full database URL
	MaxOpenConnections       int
	MaxIdleConnections       int
	ConnMaxLifeTimeInSeconds int
	MigrationFullPath        string // This has to be used when the DB migration is to be done
}

func init() {
	// NOTE: Default values
	config.Database = database{
		Main: databaseConfig{
			Url:                      "postgres://vaibhav:vaibhav@127.0.0.1:5432/vaibhav?sslmode=disable",
			MaxOpenConnections:       30,
			MaxIdleConnections:       5,
			ConnMaxLifeTimeInSeconds: 60,
			MigrationFullPath:        "/app/migrations",
		},
		Reader: databaseConfig{
			Url:                      dummyReaderDbUrl, // Random UUID
			MaxOpenConnections:       30,
			MaxIdleConnections:       5,
			ConnMaxLifeTimeInSeconds: 60,
			MigrationFullPath:        "",
		},
	}
}

func initializeDatabaseConfig() {
	if strings.TrimSpace(os.Getenv("DATABASE_URL")) != "" {
		config.Database.Main.Url = strings.TrimSpace(os.Getenv("DATABASE_URL"))
	} else {
		config.Database.Main.Url = envOrViperOrDefaultString("database.main.url", config.Database.Main.Url)
	}

	config.Database.Main.MaxOpenConnections = int(envOrViperOrDefaultInt64("database.main.maxOpenConnections",
		int64(config.Database.Main.MaxOpenConnections)))
	config.Database.Main.MaxIdleConnections = int(envOrViperOrDefaultInt64("database.main.maxIdleConnections",
		int64(config.Database.Main.MaxIdleConnections)))
	config.Database.Main.ConnMaxLifeTimeInSeconds = int(envOrViperOrDefaultInt64(
		"database.main.connMaxLifeTimeInSeconds", int64(config.Database.Main.ConnMaxLifeTimeInSeconds)))
	config.Database.Main.MigrationFullPath = envOrViperOrDefaultString(
		"database.main.migrationFullPath", config.Database.Main.MigrationFullPath)

	// DATABASE_READER_URL
	config.Database.Reader.Url = envOrViperOrDefaultString("database.reader.url", config.Database.Reader.Url)
	config.Database.Reader.MaxOpenConnections = int(envOrViperOrDefaultInt64("database.reader.maxOpenConnections",
		int64(config.Database.Reader.MaxOpenConnections)))
	config.Database.Reader.MaxIdleConnections = int(envOrViperOrDefaultInt64("database.reader.maxIdleConnections",
		int64(config.Database.Reader.MaxIdleConnections)))
	config.Database.Main.ConnMaxLifeTimeInSeconds = int(envOrViperOrDefaultInt64(
		"database.reader.connMaxLifeTimeInSeconds", int64(config.Database.Reader.ConnMaxLifeTimeInSeconds)))
	config.Database.Reader.MigrationFullPath = envOrViperOrDefaultString(
		"database.reader.migrationFullPath", config.Database.Reader.MigrationFullPath)

	// If the reader DB URL is not supplied, then we use the main database only
	// TODO: use these when we implement the
	if config.Database.Reader.Url == dummyReaderDbUrl {
		config.Database.Reader.Url = config.Database.Main.Url
		config.Database.Reader.MaxOpenConnections = config.Database.Main.MaxOpenConnections
		config.Database.Reader.MaxIdleConnections = config.Database.Main.MaxIdleConnections
		config.Database.Reader.ConnMaxLifeTimeInSeconds = config.Database.Main.ConnMaxLifeTimeInSeconds
	}
}
