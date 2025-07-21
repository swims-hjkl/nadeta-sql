package processors

// sqlmigrate delete --name

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/swims/nadeta-sql/helpers"
	"github.com/swims/nadeta-sql/migrations"
	"github.com/swims/nadeta-sql/types"
)

func RunDelete(deleteFlagData *types.DeleteFlagData, migrationStore *migrations.MigrationStore) error {
	inputMigrationName := *deleteFlagData.Name
	if inputMigrationName == "" {
		return helpers.GetMissingArgError("name")
	}
	repeatFlag := true
	reader := bufio.NewReader(os.Stdin)
	for repeatFlag {
		fmt.Printf("Are you sure you want to delete %s? Did you mean to run 'down'?\nY/y - Yes\nN/n - No\n", inputMigrationName)
		YNInput, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		YNInput = strings.TrimSpace(YNInput)
		fmt.Println("You selected", YNInput)
		if YNInput == "Y" || YNInput == "y" {
			repeatFlag = false
			err = migrationStore.FileDeleteMigration(inputMigrationName)
			if err != nil {
				return err
			}
		}
		if YNInput == "N" || YNInput == "n" {
			repeatFlag = false
		}
	}
	return nil
}
