package processors

import (
	"fmt"

	"github.com/swims/nadeta-sql/helpers"
	"github.com/swims/nadeta-sql/migrations"
	"github.com/swims/nadeta-sql/types"
)

func RunCreate(createFlagData *types.CreateFlagData, migrationStore *migrations.MigrationStore) error {
	migrationName := *createFlagData.Name
	if migrationName == "" {
		return helpers.GetMissingArgError("name")
	}
	err := migrationStore.FileCreateMigration(migrationName)
	if err != nil {
		return err
	}
	fmt.Println("created migration files")
	return nil
}
