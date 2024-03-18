package usecase

import (
	"fmt"
	"os"
	"testing"

	"github.com/flatmix/final-otus-project/internal/migrator/config"
	"github.com/stretchr/testify/assert"
)

func TestToSnakeCase_ok(t *testing.T) {
	stringForTest := "toSnakeCase"
	expectedString := "to_snake_case"

	resultString := toSnakeCase(stringForTest)

	assert.Equal(t, expectedString, resultString)
}

func TestCreateEmptyFile_ok(t *testing.T) {
	file, err := createEmptyFile("testMigrationFile")

	defer os.RemoveAll(config.FolderName)
	assert.NoError(t, err)
	resultFile, err := os.Stat(file.Name())
	assert.NoError(t, err)
	fileStat, err := file.Stat()
	assert.NoError(t, err)
	assert.Equal(t, fileStat.Size(), resultFile.Size())
}

func TestCreateEmptyFile_fail(t *testing.T) {
	_, err := createEmptyFile("testMigrationFile")

	defer os.RemoveAll(config.FolderName)
	assert.NoError(t, err)
	_, err = createEmptyFile("testMigrationFile")
	assert.Error(t, err)
}

func TestContinueIfDuplicates_ok(t *testing.T) {
	result := continueIfDuplicates("test_migration_file")
	assert.True(t, result)
}

func TestCreate_ok(t *testing.T) {
	err := Create("test_migration_file")
	defer os.RemoveAll(config.FolderName)
	assert.NoError(t, err)
	files, _ := os.ReadDir(config.FolderName)
	contentFile, err := os.ReadFile(fmt.Sprintf("%s/%s", config.FolderName, files[0].Name()))
	assert.NoError(t, err)
	content := string(contentFile)
	assert.Equal(t, migrationFileTemplate, content)
}

func TestCreate_fail(t *testing.T) {
	err := Create("test_migration_file")
	defer os.RemoveAll(config.FolderName)
	assert.NoError(t, err)
	err = Create("test_migration_file")
	assert.Error(t, err)
}
