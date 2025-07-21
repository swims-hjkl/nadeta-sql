package helpers

import (
	"os"
	"path"
	"testing"
)

func TestGetConfig(t *testing.T) {
	const testConfigFileName = ".testConfigFile.json"
	dir := t.TempDir()
	testConfigFilePath := path.Join(dir, testConfigFileName)
	_, err := getConfigFromPath("randomstring.json")
	if err != ErrConfigNotFound {
		t.Fatal("Expected ErrConfigNotFound")
	}

	// test for happy scenario
	content := `{"connectionString":"testConnectionString","directoryName":"testDirectoryName"}`
	if err := os.WriteFile(testConfigFilePath, []byte(content), 0644); err != nil {
		t.Fatal("couldn't write to temp config file")
	}
	config, err := getConfigFromPath(testConfigFilePath)
	if err != nil {
		t.Fatal("Expected no errors")
	}
	if config.DirectoryName != "testDirectoryName" {
		t.Fatal("Expected config.DirectoryName = testDirectoryName found config.DirectoryName =", config.DirectoryName)
	}
	if config.ConnectionString != "testConnectionString" {
		t.Fatal("Expected config.ConnectionString = testConnectionString found config.testConnectionString=", config.ConnectionString)
	}

	// test for malformed JSON
	content = `{"connectionString":"testConnectionString","directoryName":"testDirectoryName"},`
	if err := os.WriteFile(testConfigFilePath, []byte(content), 0644); err != nil {
		t.Fatal("couldn't write to temp config file")
	}
	_, err = getConfigFromPath(testConfigFilePath)
	if err == nil {
		t.Fatal("Expected invalid character error")
	}

	// test for keys missing
	content = `{"connectionStrin":"testConnectionString","directoryName":"testDirectoryName"}`
	if err := os.WriteFile(testConfigFilePath, []byte(content), 0644); err != nil {
		t.Fatal("couldn't write to temp config file")
	}
	_, err = getConfigFromPath(testConfigFilePath)
	if err != ErrConfigDataMissing {
		t.Fatal("Expected ErrConfigDataMissing, not found")
	}
}
