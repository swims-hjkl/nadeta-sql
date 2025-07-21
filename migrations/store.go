package migrations

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/swims/nadeta-sql/constants"
	"github.com/swims/nadeta-sql/dbutil"
	"github.com/swims/nadeta-sql/types"
)

var ErrNoRows = errors.New("no rows found")

// should handle filesystem and database
type MigrationStore struct {
	dbUtil *dbutil.DBUtil
	config *types.Config
}

func NewMigrationStore(dbUtil *dbutil.DBUtil, config *types.Config) *MigrationStore {
	migrationStore := &MigrationStore{}
	migrationStore.dbUtil = dbUtil
	migrationStore.config = config
	return migrationStore
}

func (migrationStore *MigrationStore) GetConfig() *types.Config {
	return migrationStore.config
}

func (migrationStore *MigrationStore) GetDBUtil() *dbutil.DBUtil {
	return migrationStore.dbUtil
}

func isSQLPresentInMigration(query *string) bool {
	for line := range strings.SplitSeq(*query, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "--") {
			continue
		}
		return true
	}
	return false
}

func (migrationStore *MigrationStore) RunMigration(migrationName string, migrationType types.MigrationType, isDryRun bool) error {
	query, err := migrationStore.readMigrationFile(migrationName, migrationType)
	if err != nil {
		return err
	}
	if !isSQLPresentInMigration(query) {
		return errors.New("No SQL found in migration")
	}
	migrationStore.dbUtil.RunExec(*query, nil, isDryRun)
	if !isDryRun {
		if migrationType == types.MIGRATION_TYPE_DOWN {
			err := migrationStore.DBDeleteMigrationData(migrationName)
			if err != nil {
				return err
			}
		}
		if migrationType == types.MIGRATION_TYPE_UP {
			err := migrationStore.dbInsertMigrationData(migrationName)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (migrationStore *MigrationStore) dbInsertMigrationData(migrationName string) error {
	query := fmt.Sprintf(
		"INSERT INTO %s (MigrationName, UnixOrder) VALUES (?, ?)",
		constants.SQL_SCHEMA_MIGRATION_TABLE_NAME,
	)
	err := migrationStore.dbUtil.RunExec(query, []any{migrationName, time.Now().UnixMicro()}, false)
	if err != nil {
		return err
	}
	return nil
}

func (migrationStore *MigrationStore) DBGetLatestMigration(isDryRun bool) (*types.SchemaMigrationRow, error) {
	// return nil if no rows are found
	schemaMigrationRow := &types.SchemaMigrationRow{}
	queryString := fmt.Sprintf(
		"SELECT * FROM %s ORDER BY UnixOrder DESC LIMIT 1;",
		constants.SQL_SCHEMA_MIGRATION_TABLE_NAME,
	)
	rows, err := migrationStore.dbUtil.RunQuery(queryString, nil, isDryRun)
	if err != nil {
		return nil, err
	}
	if rows == nil {
		return nil, ErrNoRows
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(
			&schemaMigrationRow.ID,
			&schemaMigrationRow.MigrationName,
			&schemaMigrationRow.AppliedAt,
			&schemaMigrationRow.UnixOrder,
		)
		if err != nil {
			return nil, err
		} else {
			return schemaMigrationRow, nil
		}
	}
	return nil, nil
}

func (migrationStore *MigrationStore) FileListMigrations() ([]string, error) {

	dirEntries, err := os.ReadDir(migrationStore.config.DirectoryName)
	if err != nil {
		return nil, err
	}

	var migrationFiles []string

	for _, dirEntry := range dirEntries {

		if !dirEntry.Type().IsRegular() {
			continue
		}

		//consider only up sql statements
		if strings.Contains(dirEntry.Name(), ".up.sql") {

			fileName := dirEntry.Name()
			unixTime, _, sepFound := strings.Cut(fileName, "-")
			if !sepFound {
				continue
			}

			_, err := strconv.Atoi(unixTime)
			if err != nil {
				continue
			}

			trimmed := strings.TrimSuffix(dirEntry.Name(), ".up.sql")
			migrationFiles = append(migrationFiles, trimmed)
		}
	}

	sort.Slice(migrationFiles, func(i, j int) bool {
		unixI, _, _ := strings.Cut(migrationFiles[i], "-")
		unixJ, _, _ := strings.Cut(migrationFiles[j], "-")

		iVal, _ := strconv.Atoi(unixI)
		jVal, _ := strconv.Atoi(unixJ)
		return iVal < jVal
	})

	return migrationFiles, nil
}

func (migrationStore *MigrationStore) FileCreateMigration(name string) error {
	currentTime := strconv.FormatInt(time.Now().Unix(), 10)
	upFileName := fmt.Sprintf("%s-%s.up.sql", currentTime, name)
	upFilePath := path.Join(migrationStore.config.DirectoryName, upFileName)
	upFile, err := os.Create(upFilePath)
	if err != nil {
		return err
	}
	defer upFile.Close()
	downFileName := fmt.Sprintf("%s-%s.down.sql", currentTime, name)
	downFilePath := path.Join(migrationStore.config.DirectoryName, downFileName)
	downFile, err := os.Create(downFilePath)
	if err != nil {
		return err
	}
	defer downFile.Close()
	return nil
}

func (migrationStore *MigrationStore) deleteFile(migrationFilePath string) error {
	if _, err := os.Stat(migrationFilePath); err != nil {
		return errors.New("Migration not found")
	}
	err := os.Remove(migrationFilePath)
	if err != nil {
		return err
	}
	return nil
}

func (migrationStore *MigrationStore) FileDeleteMigration(migrationName string) error {

	upMigrationNameWithExt := fmt.Sprintf("%s%s", migrationName, ".up.sql")
	upMigrationPath := path.Join(migrationStore.config.DirectoryName, upMigrationNameWithExt)
	err := migrationStore.deleteFile(upMigrationPath)
	if err != nil {
		return err
	}

	downMigrationNameWithExt := fmt.Sprintf("%s%s", migrationName, ".down.sql")
	downMigrationPath := path.Join(migrationStore.config.DirectoryName, downMigrationNameWithExt)
	err = migrationStore.deleteFile(downMigrationPath)
	if err != nil {
		return err
	}

	return nil
}

func (migrationStore *MigrationStore) readMigrationFile(migrationName string,
	migrationType types.MigrationType) (*string, error) {
	migrationFilePath := path.Join(migrationStore.config.DirectoryName, migrationName)
	var file *os.File
	var err error
	if migrationType == types.MIGRATION_TYPE_UP {
		file, err = os.Open(fmt.Sprintf("%s%s", migrationFilePath, ".up.sql"))
	} else {
		file, err = os.Open(fmt.Sprintf("%s%s", migrationFilePath, ".down.sql"))
	}
	if err != nil {
		return nil, err
	}
	data, err := io.ReadAll(file)
	dataStr := string(data)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return &dataStr, nil
}

func (migrationStore *MigrationStore) DBDeleteMigrationData(migrationName string) error {
	queryString := fmt.Sprintf(
		"DELETE FROM %s WHERE MigrationName = ?",
		constants.SQL_SCHEMA_MIGRATION_TABLE_NAME,
	)
	err := migrationStore.dbUtil.RunExec(queryString, []any{migrationName}, false)
	if err != nil {
		return err
	}
	return nil
}

func (migrationStore *MigrationStore) DBGetAppliedMigrations(steps int, reverse bool) (*[]types.SchemaMigrationRow, error) {
	queryString := fmt.Sprintf(
		"SELECT * FROM %s",
		constants.SQL_SCHEMA_MIGRATION_TABLE_NAME,
	)
	if reverse {
		queryString = queryString + " ORDER BY UnixOrder DESC"
	}
	if steps > 0 {
		queryString = queryString + fmt.Sprintf(" LIMIT %s", strconv.Itoa(steps))
	}
	rows, err := migrationStore.dbUtil.RunQuery(queryString, nil, false)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	schemaMigrationRows := []types.SchemaMigrationRow{}
	for rows.Next() {
		schemaMigrationRow := types.SchemaMigrationRow{}
		err := rows.Scan(
			&schemaMigrationRow.ID,
			&schemaMigrationRow.MigrationName,
			&schemaMigrationRow.AppliedAt,
			&schemaMigrationRow.UnixOrder,
		)
		if err != nil {
			return nil, err
		}
		schemaMigrationRows = append(schemaMigrationRows, schemaMigrationRow)
	}
	return &schemaMigrationRows, nil
}

func (migrationStore *MigrationStore) PurgeData() error {
	query := constants.SQL_DROP_TABLE_SCHEMA_MIGRATIONS
	err := migrationStore.dbUtil.RunExec(query, nil, false)
	if err != nil {
		return err
	}
	err = os.RemoveAll(migrationStore.config.DirectoryName)
	if err != nil {
		return err
	}
	return nil
}
