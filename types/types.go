package types

import (
	"flag"
	"time"
)

// Enums
type MigrationType int

const (
	MIGRATION_TYPE_UP MigrationType = iota
	MIGRATION_TYPE_DOWN
)

// Structs

type Config struct {
	ConnectionString string `json:"connectionString"`
	DirectoryName    string `json:"directoryName"`
}

// FlagData

type InitFlagData struct {
	FlagSet          flag.FlagSet
	ConnectionString *string
	DirectoryName    *string
}

type UpFlagData struct {
	FlagSet flag.FlagSet
	Steps   *int
	Dryrun  *bool
	Name    *string
}

type DownFlagData struct {
	FlagSet flag.FlagSet
	Steps   *int
	Dryrun  *bool
	Name    *string
}

type CreateFlagData struct {
	FlagSet flag.FlagSet
	Name    *string
}

type DeleteFlagData struct {
	FlagSet flag.FlagSet
	Name    *string
}

type ListFlagData struct {
	FlagSet flag.FlagSet
}

type StatusFlagData struct {
	FlagSet flag.FlagSet
	Pending *bool
	Applied *bool
}

// db related

type SchemaMigrationRow struct {
	ID            int
	MigrationName string
	AppliedAt     time.Time
	UnixOrder     int
}
