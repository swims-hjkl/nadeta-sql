package processors

// sqlmigrate status
// sqlmigrate status --pending
// sqlmigrate status --applied

import (
	"fmt"
	"slices"

	"github.com/swims/nadeta-sql/migrations"
	"github.com/swims/nadeta-sql/types"
)

func isElementInArray[T comparable](searchElement T, array []T) bool {
	return slices.Contains(array, searchElement)
}

func RunStatus(statusFlagData *types.StatusFlagData, migrationStore *migrations.MigrationStore) error {
	appliedFlag := *statusFlagData.Applied
	pendingFlag := *statusFlagData.Pending

	if appliedFlag == false && pendingFlag == false {
		appliedFlag = true
		pendingFlag = true
	}

	fmt.Printf("Applied Migrations:\n")

	appliedMigrations, err := migrationStore.DBGetAppliedMigrations(0, false)
	if err != nil {
		return err
	}

	var appliedMigrationNames []string

	for _, migration := range *appliedMigrations {
		if appliedFlag {
			fmt.Println(migration.MigrationName)
		}
		if pendingFlag {
			appliedMigrationNames = append(appliedMigrationNames, migration.MigrationName)
		}
	}

	fmt.Printf("\nPending Migrations:\n")

	if pendingFlag {
		migrationNames, err := migrationStore.FileListMigrations()
		if err != nil {
			return err
		}
		for _, migrationName := range migrationNames {
			if !isElementInArray(migrationName, appliedMigrationNames) {
				fmt.Println(migrationName)
			}
		}

	}

	return nil
}
