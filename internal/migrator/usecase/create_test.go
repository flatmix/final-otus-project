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

	cfg := config.Config{
		FolderName: "test",
	}
	createStruct := NewCreateStruct(stringForTest, &cfg)

	resultString := createStruct.toSnakeCase()

	assert.Equal(t, expectedString, resultString)
}

func TestCreateEmptyFile_ok(t *testing.T) {
	nameForTest := "testMigrationFile"

	cfg := config.Config{
		FolderName: "test",
	}
	createStruct := NewCreateStruct(nameForTest, &cfg)

	file, err := createStruct.createEmptyFile()

	defer os.RemoveAll(cfg.FolderName)
	assert.NoError(t, err)
	resultFile, err := os.Stat(file.Name())
	assert.NoError(t, err)
	fileStat, err := file.Stat()
	assert.NoError(t, err)
	assert.Equal(t, fileStat.Size(), resultFile.Size())
}

func TestCreateEmptyFile_fail(t *testing.T) {
	nameForTest := "testMigrationFile"

	cfg := config.Config{
		FolderName: "test",
	}
	createStruct := NewCreateStruct(nameForTest, &cfg)

	_, err := createStruct.createEmptyFile()

	defer os.RemoveAll(cfg.FolderName)
	assert.NoError(t, err)
	_, err = createStruct.createEmptyFile()
	assert.Error(t, err)
}

func TestContinueIfDuplicates_ok(t *testing.T) {
	nameForTest := "testMigrationFile"

	cfg := config.Config{
		FolderName: "test",
	}
	createStruct := NewCreateStruct(nameForTest, &cfg)

	result := createStruct.continueIfDuplicates()
	assert.True(t, result)
}

func TestCreate_ok(t *testing.T) {
	cfg := config.Config{
		FolderName: "test",
	}

	err := Create("test_migration_file", &cfg)
	defer os.RemoveAll(cfg.FolderName)
	assert.NoError(t, err)
	files, _ := os.ReadDir(cfg.FolderName)
	contentFile, err := os.ReadFile(fmt.Sprintf("%s/%s", cfg.FolderName, files[0].Name()))
	assert.NoError(t, err)
	content := string(contentFile)
	assert.Equal(t, migrationFileTemplate, content)
}

func TestCreate_fail(t *testing.T) {
	cfg := config.Config{
		FolderName: "test",
	}

	err := Create("test_migration_file", &cfg)
	defer os.RemoveAll(cfg.FolderName)
	assert.NoError(t, err)
	err = Create("test_migration_file", &cfg)
	assert.Error(t, err)
}
