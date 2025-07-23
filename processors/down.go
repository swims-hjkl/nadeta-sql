package processors

// sqlmigrate down --steps 2
// sqlmigrate down --name migrationname
// sqlmigrate down --steps 2 --dryrun
// sqlmigrate down --name migrationname --dryrun

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/swims/nadeta-sql/migrations"
	"github.com/swims/nadeta-sql/types"
)

func RunDown(downFlagData *types.DownFlagData, migrationStore *migrations.MigrationStore) error {
	isDryRun := *downFlagData.Dryrun
	steps := *downFlagData.Steps

	if isDryRun {
		fmt.Println("Initiating with dryrun")
	}

	if steps == 0 && !isDryRun {
		repeatFlag := true
		reader := bufio.NewReader(os.Stdin)
		for repeatFlag {
			fmt.Println("Are you sure you want to run down all the migrations? (Y/y):")
			YNInput, err := reader.ReadString('\n')
			if err != nil {
				return err
			}
			YNInput = strings.TrimSpace(YNInput)
			fmt.Println("You selected", YNInput)
			if YNInput == "Y" || YNInput == "y" {
				fmt.Printf("Running down all migrations...\n\n")
				repeatFlag = false
			} else {
				fmt.Println("Quitting...")
				repeatFlag = false
				return nil
			}
		}

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
