package helpers

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"

	"github.com/swims/nadeta-sql/constants"
	"github.com/swims/nadeta-sql/types"
)

func CreateFolder(folderCreationPath string) (err error) {
	err = os.Mkdir(folderCreationPath, 0755)
	if err != nil {
		if os.IsExist(err) {
			// do nothing
		} else {
			log.Printf("%v", err)
			return err
		}
	}
	return nil
}

func GetConfig() (*types.Config, error) {
	file, err := os.Open(constants.CONFIG_FILE)
	if os.IsNotExist(err) {
		return nil, errors.New("Config not present")
	}
	if err != nil {
		return nil, err
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	var config types.Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	// Unmarshal does not return error if it can't find suitable keys to put the json into, hence check
	if config.ConnectionString == "" || config.DirectoryName == "" {
		return nil, errors.New("Unmarshal didn't work")
	}
	return &config, nil
}
