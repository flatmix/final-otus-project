package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/flatmix/final-otus-project/internal/migrator/storage"
	"github.com/flatmix/final-otus-project/internal/migrator/usecase"
	mocks "github.com/flatmix/final-otus-project/mocks/migrator/usecase"
	"github.com/stretchr/testify/assert"
)

func TestStatus_OK(t *testing.T) {
	ctx := context.Background()
	dbStruct := mocks.NewDBUsecaseContract(t)
	now := time.Now()
	nowDB := now.Format(time.RFC3339Nano)
	nowStruct := now.Format(time.DateTime)

	fileInfo := mocks.NewFileInfo(t)
	fileInfo2 := mocks.NewFileInfo(t)

	expfiles := []usecase.FileStruct{
		{File: fileInfo, Hash: "Hash1"},
		{File: fileInfo2, Hash: "Hash2"},
	}

	expOuts := usecase.Outs{
		&usecase.Out{
			Name:        "Test",
			Status:      "Ok",
			Version:     "1",
			TimeMigrate: nowStruct,
		}, &usecase.Out{
			Name:        "Test 2",
			Status:      "Ok",
			Version:     "2",
			TimeMigrate: nowStruct,
		},
	}

	dbStruct.EXPECT().GetAllMigrationFile().Return(expfiles, nil)
	dbStruct.EXPECT().GetMigrationRow(ctx, expfiles[0]).Return(&storage.MigrationDBStruct{
		ID:        1,
		Name:      "Test",
		Hash:      "Hash1",
		Version:   1,
		CreatedAt: nowDB,
		UpdatedAt: nowDB,
	}, nil)
	fileInfo.EXPECT().Name().Return("Test")

	dbStruct.EXPECT().GetMigrationRow(ctx, expfiles[1]).Return(&storage.MigrationDBStruct{
		ID:        2,
		Name:      "Test 2",
		Hash:      "Hash2",
		Version:   2,
		CreatedAt: nowDB,
		UpdatedAt: nowDB,
	}, nil)
	fileInfo2.EXPECT().Name().Return("Test 2")

	outs, err := usecase.Status(ctx, dbStruct)
	assert.NoError(t, err)
	assert.Equal(t, &expOuts, outs)
}

func TestStatus_NoMigrate(t *testing.T) {
	ctx := context.Background()
	dbStruct := mocks.NewDBUsecaseContract(t)

	fileInfo := mocks.NewFileInfo(t)

	expfiles := []usecase.FileStruct{
		{File: fileInfo, Hash: "Hash1"},
	}

	expOuts := usecase.Outs{
		&usecase.Out{
			Name:        "Test",
			Status:      "No migrate",
			Version:     "-",
			TimeMigrate: "-",
		},
	}

	dbStruct.EXPECT().GetAllMigrationFile().Return(expfiles, nil)
	dbStruct.EXPECT().GetMigrationRow(ctx, expfiles[0]).Return(nil, errors.New("test error"))
	fileInfo.EXPECT().Name().Return("Test")

	outs, err := usecase.Status(ctx, dbStruct)

	assert.NoError(t, err)
	assert.Equal(t, &expOuts, outs)
}

func TestStatus_EmptyMigrate(t *testing.T) {
	ctx := context.Background()
	dbStruct := mocks.NewDBUsecaseContract(t)

	fileInfo := mocks.NewFileInfo(t)

	expfiles := []usecase.FileStruct{
		{File: fileInfo, Hash: "Hash1"},
	}

	expOuts := usecase.Outs{
		&usecase.Out{
			Name:        "Test",
			Status:      "No migrate",
			Version:     "-",
			TimeMigrate: "-",
		},
	}

	dbStruct.EXPECT().GetAllMigrationFile().Return(expfiles, nil)
	dbStruct.EXPECT().GetMigrationRow(ctx, expfiles[0]).Return(nil, nil)
	fileInfo.EXPECT().Name().Return("Test")

	outs, err := usecase.Status(ctx, dbStruct)

	assert.NoError(t, err)
	assert.Equal(t, &expOuts, outs)
}

func TestStatus_Fail(t *testing.T) {
	ctx := context.Background()
	dbStruct := mocks.NewDBUsecaseContract(t)

	expError := errors.New("test error")

	dbStruct.EXPECT().GetAllMigrationFile().Return(nil, expError)

	_, err := usecase.Status(ctx, dbStruct)

	assert.ErrorIs(t, err, expError)
}
