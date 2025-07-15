package processors

// sqlmigrate down --steps 2
// sqlmigrate down --name migrationname
// sqlmigrate down --steps 2 --dryrun
// sqlmigrate down --name migrationname --dryrun

import (
	"errors"
	"fmt"

	"github.com/swims/nadeta-sql/migrations"
	"github.com/swims/nadeta-sql/types"
)

func RunDown(downFlagData *types.DownFlagData, migrationStore *migrations.MigrationStore) error {
	isDryRun := *downFlagData.Dryrun
	steps := *downFlagData.Steps

	if isDryRun {
		fmt.Println("Initiating with dryrun")
	}

	appliedMigrations, err := migrationStore.DBGetAppliedMigrations(steps, true)
	if err != nil {
		return err
	}

	if steps < 0 {
		return errors.New("steps cannot be negative, 0 - all, > 0 - specified amount")
	}

	for _, appliedMigration := range *appliedMigrations {
		fmt.Println("Reverting SQL Migration ...", appliedMigration.MigrationName)
		err := migrationStore.RunMigration(
			appliedMigration.MigrationName,
			types.MIGRATION_TYPE_DOWN,
			isDryRun,
		)
		if err != nil {
			return err
		}
		fmt.Print("Reverted\n\n")
	}
	return nil
}
