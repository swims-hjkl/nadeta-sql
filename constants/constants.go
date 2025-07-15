package constants

import (
	"fmt"
)

const SQL_SCHEMA_MIGRATION_TABLE_NAME string = "SchemaMigrations"

var SQL_CREATE_TABLE_SCHEMA_MIGRATIONS = fmt.Sprintf(`
CREATE TABLE IF NOT EXISTS %s (
	ID INTEGER PRIMARY KEY AUTOINCREMENT,
	MigrationName TEXT NOT NULL,
	AppliedAt DATETIME DEFAULT (datetime('now')),
	UnixOrder INTEGER
);`, SQL_SCHEMA_MIGRATION_TABLE_NAME)

const CONFIG_FILE string = ".nadeta-sql-config.json"

var SQL_DROP_TABLE_SCHEMA_MIGRATIONS = fmt.Sprintf("DROP TABLE %s;", SQL_SCHEMA_MIGRATION_TABLE_NAME)
