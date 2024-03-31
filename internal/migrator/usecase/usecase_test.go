package usecase_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/flatmix/final-otus-project/internal/migrator/config"
	"github.com/flatmix/final-otus-project/internal/migrator/usecase"
	mocks "github.com/flatmix/final-otus-project/mocks/migrator/usecase"
	"github.com/stretchr/testify/assert"
)

const (
	folderName = "test"
	fileName   = "test"
	fileName2  = "test2"
)

func Test_GetHash_OK(t *testing.T) {
	dbMock := mocks.NewDBContract(t)
	cfg := config.Config{
		FolderName: folderName,
	}

	_, err := createEmptyFile(fileName)
	assert.NoError(t, err)
	defer func() {
		err = deleteTestFile()
		assert.NoError(t, err)
	}()

	expHash := "47DEQpj8HBSa-_TImW-5JCeuQeRkm5NMpJWZG3hSuFU="

	uc := usecase.NewUsecase(dbMock, &cfg)
	hash, err := uc.GetHash(fileName)
	assert.NoError(t, err)
	assert.Equal(t, expHash, hash)
}

func Test_GetAllMigrationFile_OK(t *testing.T) {
	dbMock := mocks.NewDBContract(t)
	cfg := config.Config{
		FolderName: folderName,
	}

	file1, err := createEmptyFile(fileName)
	assert.NoError(t, err)
	file2, err := createEmptyFile(fileName2)
	assert.NoError(t, err)

	defer func() {
		err = deleteTestFile()
		assert.NoError(t, err)
	}()

	uc := usecase.NewUsecase(dbMock, &cfg)

	expFiles := []usecase.FileStruct{
		{
			File: file1,
			Hash: "47DEQpj8HBSa-_TImW-5JCeuQeRkm5NMpJWZG3hSuFU=",
		},
		{
			File: file2,
			Hash: "47DEQpj8HBSa-_TImW-5JCeuQeRkm5NMpJWZG3hSuFU=",
		},
	}

	files, err := uc.GetAllMigrationFile()
	assert.NoError(t, err)
	assert.Equal(t, expFiles, files)
}

func Test_GetAllMigrationFileMap_OK(t *testing.T) {
	dbMock := mocks.NewDBContract(t)
	cfg := config.Config{
		FolderName: folderName,
	}

	file1, err := createEmptyFile(fileName)
	assert.NoError(t, err)
	file2, err := createEmptyFile(fileName2)
	assert.NoError(t, err)

	defer func() {
		err = deleteTestFile()
		assert.NoError(t, err)
	}()

	uc := usecase.NewUsecase(dbMock, &cfg)

	expFiles := make(usecase.FilesMap)

	expFiles[fileName] = usecase.FileStruct{
		File: file1,
		Hash: "47DEQpj8HBSa-_TImW-5JCeuQeRkm5NMpJWZG3hSuFU=",
	}
	expFiles[fileName2] = usecase.FileStruct{
		File: file2,
		Hash: "47DEQpj8HBSa-_TImW-5JCeuQeRkm5NMpJWZG3hSuFU=",
	}

	files, err := uc.GetAllMigrationFileMap()
	assert.NoError(t, err)
	assert.Equal(t, expFiles, files)
}

func createMigrationFolder() bool {
	err := os.MkdirAll(folderName, 0o755)
	return err == nil
}

func createEmptyFile(name string) (os.FileInfo, error) {
	createMigrationFolder()

	filename := fmt.Sprintf("%s/%s", folderName, name)
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, err
	}
	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}
	return fileInfo, nil
}

func deleteTestFile() error {
	err := os.RemoveAll(folderName)
	if err != nil {
		return err
	}
	return nil
}
