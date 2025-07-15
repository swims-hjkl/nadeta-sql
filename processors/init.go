package processors

import (
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/swims/nadeta-sql/constants"
	"github.com/swims/nadeta-sql/dbutil"
	"github.com/swims/nadeta-sql/helpers"
	"github.com/swims/nadeta-sql/types"
)

func RunInit(initFlagData *types.InitFlagData) error {
	connectionString, directoryName, err := getArguments(initFlagData)
	if err != nil {
		return err
	}
	config := types.Config{
		ConnectionString: *connectionString,
		DirectoryName:    *directoryName,
	}
	err = createConfigFile(&config)
	err = createMigrationFolder(config.DirectoryName)
	dbUtil, err := dbutil.NewDBUtil(config.ConnectionString)
	if err != nil {
		return err
	}
	err = dbUtil.CreateDatabase()
	if err != nil {
		return err
	}
	err = dbUtil.RunExec(constants.SQL_CREATE_TABLE_SCHEMA_MIGRATIONS, nil, false)
	if err != nil {
		return err
	}
	return nil
}

func createConfigFile(config *types.Config) error {
	file, err := os.Create(constants.CONFIG_FILE)
	if err != nil {
		return err
	}
	defer file.Close()
	bytes, err := json.Marshal(map[string]string{
		"connectionString": config.ConnectionString,
		"directoryName":    config.DirectoryName,
	})
	_, err = file.Write(bytes)
	if err != nil {
		return err
	}
	return nil
}

func createMigrationFolder(directoryName string) error {
	workingDirectory, err := os.Getwd()
	folderCreationPath := path.Join(workingDirectory, directoryName)
	fmt.Println(folderCreationPath)
	err = helpers.CreateFolder(folderCreationPath)
	if err != nil {
		return err
	}
	return nil
}

func getArguments(initFlagData *types.InitFlagData) (*string, *string, error) {
	connectionString := *initFlagData.ConnectionString
	directoryName := *initFlagData.DirectoryName
	if connectionString == "" {
		err := helpers.GetMissingArgError("connection-string")
		if err != nil {
			return nil, nil, err
		}
	}
	return &connectionString, &directoryName, nil
}
