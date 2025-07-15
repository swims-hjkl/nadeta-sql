package processors

// sqlmigrate up --steps 2
// sqlmigrate up --name migrationname
// sqlmigrate up --steps 2 --dryrun
// sqlmigrate up --name migrationname --dryrun

import (
	"errors"
	"fmt"

	"github.com/swims/nadeta-sql/migrations"
	"github.com/swims/nadeta-sql/types"
)

func RunUp(upFlagData *types.UpFlagData, migrationStore *migrations.MigrationStore) error {
	isDryRun := *upFlagData.Dryrun
	steps := *upFlagData.Steps
	stepsFlag := false
	fromHereFlag := false

	if isDryRun {
		fmt.Println("Initiating with dryrun")
	}

	lastRow, err := migrationStore.DBGetLatestMigration(isDryRun)
	if err != nil {
		return err
	}
	if lastRow == nil {
		fmt.Println("found no previous up migrations\n")
		fromHereFlag = true
	} else {
		fmt.Println(lastRow.MigrationName)
	}

	allMigrations, err := migrationStore.FileListMigrations()
	if err != nil {
		return err
	}

	if steps < 0 {
		return errors.New("steps cannot be negative, 0 - all, > 0 - specified amount")
	}

	if steps != 0 {
		stepsFlag = true
	}

	counter := 0
	for _, migration := range allMigrations {
		if fromHereFlag {
			fmt.Println("Applying SQL Migration ...", migration)
			err := migrationStore.RunMigration(migration, types.MIGRATION_TYPE_UP, isDryRun)
			if err != nil {
				return err
			}
			fmt.Print("Applied\n\n")
			counter = counter + 1
			if stepsFlag {
				if counter == steps {
					fmt.Printf("Stopping at %s observing steps\n", migration)
					break
				}
			}
		}
		if lastRow != nil && migration == lastRow.MigrationName {
			fromHereFlag = true
			fmt.Println("Running from after:", lastRow.MigrationName)
		}
	}
	return nil
}
