package processors

// sqlmigrate list

import (
	"github.com/swims/nadeta-sql/migrations"
)

func RunList(migrationStore *migrations.MigrationStore) error {
	fileList, err := migrationStore.FileListMigrations()
	if err != nil {
		return err
	}
	for _, file := range fileList {
		println(file)
	}
	return nil
}
