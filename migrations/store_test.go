package migrations

import (
	"fmt"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/swims/nadeta-sql/constants"
	"github.com/swims/nadeta-sql/dbutil"
	"github.com/swims/nadeta-sql/types"
)

func TestMigrationStore(t *testing.T) {
	tempDir := t.TempDir()
	databaseFilePath := path.Join(tempDir, "test.db")
	config := types.Config{
		ConnectionString: fmt.Sprintf("file:%s", databaseFilePath),
		DirectoryName:    tempDir,
	}
	db, err := dbutil.NewDBUtil(config.ConnectionString)
	if err != nil {
		t.Errorf("Expected creation of database, not created")
	}
	store := NewMigrationStore(db, &config)
	store.dbUtil.RunExec(constants.SQL_CREATE_TABLE_SCHEMA_MIGRATIONS, nil, false)

	// test GetConfig
	gotConfig := store.GetConfig()
	if gotConfig.ConnectionString != config.ConnectionString {
		t.Errorf("Expected ConnectionString %s, got %s", config.ConnectionString, gotConfig.ConnectionString)
	}
	if gotConfig.DirectoryName != config.DirectoryName {
		t.Errorf("Expected DirectoryName %s, got %s", config.DirectoryName, gotConfig.DirectoryName)
	}

	// test FileCreateMigration
	migrationName := "CreateA"
	migrationSQL := "CREATE TABLE A (ID TEXT)"
	store.FileCreateMigration(migrationName)

	entries, err := os.ReadDir(config.DirectoryName)
	if err != nil {
		t.Error("Expected to read dir, errored out")
	}
	var matchingFileName string
	for _, entry := range entries {
		fmt.Println("entry", entry.Name())
		if !entry.IsDir() {
			if strings.Contains(entry.Name(), fmt.Sprintf("%s.up.sql", migrationName)) {
				matchingFileName = entry.Name()
			}
		}
	}
	if matchingFileName == "" {
		t.Errorf("Expected to find %s matching migration file, not found", migrationName)
	}

	err = os.WriteFile(path.Join(config.DirectoryName, matchingFileName), []byte(migrationSQL), 644)
	if err != nil {
		t.Fatal("Expected write sql to migration file, errored out")
	}

	// test RunMigration
	err = store.RunMigration(strings.TrimSuffix(matchingFileName, ".up.sql"), types.MIGRATION_TYPE_UP, false)
	if err != nil {
		t.Errorf("%v", err)
		t.Fatal("Expected RunMigration to work, errored out")
	}
	rows, err := store.dbUtil.RunQuery("SELECT count(*) FROM sqlite_master WHERE type='table' AND name=?;", []any{"A"}, false)
	if err != nil {
		t.Errorf("%v", err)
		t.Error("Expected successful query, errored out")
	}
	if rows == nil {
		t.Error("Expected rows not to be nil")
	}
	defer rows.Close()
	found := false
	for rows.Next() {
		var tableACount int
		rows.Scan(&tableACount)
		if tableACount == 1 {
			found = true
		}
	}
	if !found {
		t.Error("Expected table 'A' to be found, not found")
	}

	// test DBGetLatestMigration
	schemaRow, err := store.DBGetLatestMigration(false)
	if err != nil {
		t.Error("Expected succssful fetch of migrations")
	}
	if schemaRow.MigrationName != strings.TrimSuffix(matchingFileName, ".up.sql") {
		t.Error("Expected succssful fetch of migrations, not fetched")
	}

	// test FileListMigrations
	listOfMigrations, err := store.FileListMigrations()
	if err != nil {
		t.Error("Expected successful list of migrations, not listed")
	}
	if !(strings.Contains(listOfMigrations[0], migrationName)) {
		t.Errorf("Expected first migration to be %s, found %s", migrationName, listOfMigrations[0])
	}
	if len(listOfMigrations) != 1 {
		t.Errorf("Expected length of listOfMigrations to be 1, found %d", len(listOfMigrations))
	}

	// test DBGetAppliedMigrations
	listOfAppliedMigrations, err := store.DBGetAppliedMigrations(0, false)
	if err != nil {
		t.Error("Expected successful list of migrations, not listed")
	}
	if len(*listOfAppliedMigrations) != 1 {
		t.Errorf("Expected length of listOfAppliedMigrations to be 1, found %d", len(*listOfAppliedMigrations))
	}
	if (*listOfAppliedMigrations)[0].MigrationName != schemaRow.MigrationName {
		t.Errorf("Expected first migration to be %s, found %s", schemaRow.MigrationName, listOfMigrations[0])
	}

	// test FileDeleteMigration
	err = store.FileDeleteMigration(schemaRow.MigrationName)
	if err != nil {
		t.Error("Expected deletion of migration file, errored out")
	}
	_, err = os.Stat(path.Join(config.DirectoryName, fmt.Sprintf("%s%s", schemaRow.MigrationName, ".up.sql")))
	if !os.IsNotExist(err) {
		t.Error("Expected deletion of migration file, not deleted")
	}
	_, err = os.Stat(path.Join(config.DirectoryName, fmt.Sprintf("%s%s", schemaRow.MigrationName, ".down.sql")))
	if !os.IsNotExist(err) {
		t.Error("Expected deletion of migration file, not deleted")
	}
	// test DBDeleteMigrationData
	err = store.DBDeleteMigrationData(schemaRow.MigrationName)
	if err != nil {
		t.Error("Expected successful deletion of migration data, errored out")
	}

	if appliedRows, _ := store.DBGetAppliedMigrations(0, false); len(*appliedRows) != 0 {
		t.Errorf("Expected 0 applied rows found %d", len(*appliedRows))
	}

	// test PurgeData
	err = store.PurgeData()
	if err != nil {
		t.Errorf("Expected purge data to work, errored out %v", err)
	}
	query := "SELECT count(*) FROM sqlite_master WHERE type='table' AND name=?;"
	rows, err = store.dbUtil.RunQuery(query, []any{migrationName}, false)
	if err != nil {
		t.Error("Expected query data to work, errored out")
	}
	if rows == nil {
		t.Error(ErrNoRows.Error())
	}
	if rows.Next() {
		var count int
		rows.Scan(count)
		if count != 0 {
			t.Errorf("Expected count to be 0, found %d", count)
		}
	}
	if _, err = os.Stat(config.DirectoryName); !os.IsNotExist(err) {
		t.Errorf("Expected config folder %s not to exist, exists", config.DirectoryName)
	}
}
