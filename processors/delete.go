package processors

// sqlmigrate delete --name

import (
	"github.com/swims/nadeta-sql/helpers"
	"github.com/swims/nadeta-sql/migrations"
	"github.com/swims/nadeta-sql/types"
)

func RunDelete(deleteFlagData *types.DeleteFlagData, migrationStore *migrations.MigrationStore) error {
	inputMigrationName := *deleteFlagData.Name
	if inputMigrationName == "" {
		return helpers.GetMissingArgError("name")
	}
	err := migrationStore.FileDeleteMigration(inputMigrationName)
	if err != nil {
		return err
	}
	return nil
}
