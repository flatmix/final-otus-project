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

const migrateSQLString = "test"

func TestUp_Ok(t *testing.T) {
	ctx := context.Background()
	dbMock := mocks.NewDBUsecaseContract(t)

	fileInfo := mocks.NewFileInfo(t)
	fileInfo2 := mocks.NewFileInfo(t)

	expfiles := []usecase.FileStruct{
		{File: fileInfo, Hash: "Hash1"},
		{File: fileInfo2, Hash: "Hash2"},
	}

	actualVersion := 1

	expOuts := usecase.Outs{
		&usecase.Out{
			Name:   "test3",
			Status: "migrate ok",
		}, &usecase.Out{
			Name:   "test4",
			Status: "migrate ok",
		},
	}

	dbMock.EXPECT().GetAllMigrationFile().Return(expfiles, nil)
	dbMock.EXPECT().CreateMigrationsTable(ctx).Return(nil)
	dbMock.EXPECT().GetActualVersion(ctx).Return(actualVersion)

	dbMock.EXPECT().GetMigrationRow(ctx, expfiles[0]).Return(nil, nil)
	fileInfo.EXPECT().Name().Return("test3").Once()
	dbMock.EXPECT().GetUpPart(expfiles[0]).Return(migrateSQLString, nil)
	dbMock.EXPECT().Migrate(ctx, migrateSQLString).Return(nil)
	dbMock.EXPECT().CreateMigration(ctx, expfiles[0], actualVersion).Return(nil)

	dbMock.EXPECT().GetMigrationRow(ctx, expfiles[1]).Return(nil, nil)
	fileInfo2.EXPECT().Name().Return("test4")
	dbMock.EXPECT().GetUpPart(expfiles[1]).Return(migrateSQLString, nil)
	dbMock.EXPECT().Migrate(ctx, migrateSQLString).Return(nil)
	dbMock.EXPECT().CreateMigration(ctx, expfiles[1], actualVersion).Return(nil)

	outs, err := usecase.Up(ctx, dbMock)
	assert.NoError(t, err)
	assert.Equal(t, &expOuts, outs)
}

func TestUpMigration_OK(t *testing.T) {
	ctx := context.Background()
	dbMock := mocks.NewDBUsecaseContract(t)

	actualVersion := 2

	fileInfo := mocks.NewFileInfo(t)
	fileStruct := usecase.FileStruct{File: fileInfo, Hash: "Hash1"}

	dbMock.EXPECT().GetMigrationRow(ctx, fileStruct).Return(nil, nil)
	fileInfo.EXPECT().Name().Return("test5")
	dbMock.EXPECT().GetUpPart(fileStruct).Return(migrateSQLString, nil)

	dbMock.EXPECT().Migrate(ctx, migrateSQLString).Return(nil)
	dbMock.EXPECT().CreateMigration(ctx, fileStruct, actualVersion).Return(nil)

	expOut := usecase.Out{
		Name:   "test5",
		Status: "migrate ok",
	}

	out, err := usecase.UpMigration(ctx, dbMock, fileStruct, actualVersion)
	assert.NoError(t, err)
	assert.Equal(t, &expOut, out)
}

func TestUp_GetAllMigrationFileError(t *testing.T) {
	ctx := context.Background()
	dbMock := mocks.NewDBUsecaseContract(t)

	fileInfo := mocks.NewFileInfo(t)
	fileInfo2 := mocks.NewFileInfo(t)

	expfiles := []usecase.FileStruct{
		{File: fileInfo, Hash: "Hash1"},
		{File: fileInfo2, Hash: "Hash2"},
	}

	expError := errors.New("test error")

	dbMock.EXPECT().GetAllMigrationFile().Return(expfiles, expError)

	_, err := usecase.Up(ctx, dbMock)
	assert.ErrorIs(t, err, expError)
}

func TestUp_CreateMigrationsTableError(t *testing.T) {
	ctx := context.Background()
	dbMock := mocks.NewDBUsecaseContract(t)

	fileInfo := mocks.NewFileInfo(t)
	fileInfo2 := mocks.NewFileInfo(t)

	expfiles := []usecase.FileStruct{
		{File: fileInfo, Hash: "Hash1"},
		{File: fileInfo2, Hash: "Hash2"},
	}

	expError := errors.New("test error")

	dbMock.EXPECT().GetAllMigrationFile().Return(expfiles, nil)
	dbMock.EXPECT().CreateMigrationsTable(ctx).Return(expError)

	_, err := usecase.Up(ctx, dbMock)
	assert.ErrorIs(t, err, expError)
}

func TestUp_UpMigrationError(t *testing.T) {
	ctx := context.Background()
	dbMock := mocks.NewDBUsecaseContract(t)

	fileInfo := mocks.NewFileInfo(t)
	fileInfo2 := mocks.NewFileInfo(t)

	expfiles := []usecase.FileStruct{
		{File: fileInfo, Hash: "Hash1"},
		{File: fileInfo2, Hash: "Hash2"},
	}

	actualVersion := 1

	expError := errors.New("test error")

	dbMock.EXPECT().GetAllMigrationFile().Return(expfiles, nil)
	dbMock.EXPECT().CreateMigrationsTable(ctx).Return(nil)
	dbMock.EXPECT().GetActualVersion(ctx).Return(actualVersion)

	dbMock.EXPECT().GetMigrationRow(ctx, expfiles[0]).Return(nil, expError)

	_, err := usecase.Up(ctx, dbMock)
	assert.ErrorIs(t, err, expError)
}

func TestUpMigration_GetUpPartError(t *testing.T) {
	ctx := context.Background()
	dbMock := mocks.NewDBUsecaseContract(t)

	actualVersion := 2

	fileInfo := mocks.NewFileInfo(t)
	fileStruct := usecase.FileStruct{File: fileInfo, Hash: "Hash1"}

	expError := errors.New("test error")

	dbMock.EXPECT().GetMigrationRow(ctx, fileStruct).Return(nil, nil)
	fileInfo.EXPECT().Name().Return("test5")
	dbMock.EXPECT().GetUpPart(fileStruct).Return("", expError)

	_, err := usecase.UpMigration(ctx, dbMock, fileStruct, actualVersion)
	assert.ErrorIs(t, err, expError)
}

func TestUpMigration_MigrateError(t *testing.T) {
	ctx := context.Background()
	dbMock := mocks.NewDBUsecaseContract(t)

	actualVersion := 2
	migrateSQLString := "test"

	fileInfo := mocks.NewFileInfo(t)
	fileStruct := usecase.FileStruct{File: fileInfo, Hash: "Hash1"}

	expError := errors.New("test error")

	dbMock.EXPECT().GetMigrationRow(ctx, fileStruct).Return(nil, nil)
	fileInfo.EXPECT().Name().Return("test5")
	dbMock.EXPECT().GetUpPart(fileStruct).Return(migrateSQLString, nil)

	dbMock.EXPECT().Migrate(ctx, migrateSQLString).Return(expError)

	_, err := usecase.UpMigration(ctx, dbMock, fileStruct, actualVersion)
	assert.ErrorIs(t, err, expError)
}

func TestUpMigration_CreateMigrationError(t *testing.T) {
	ctx := context.Background()
	dbMock := mocks.NewDBUsecaseContract(t)

	actualVersion := 2
	migrateSQLString := "test"

	fileInfo := mocks.NewFileInfo(t)
	fileStruct := usecase.FileStruct{File: fileInfo, Hash: "Hash1"}

	expError := errors.New("test error")

	dbMock.EXPECT().GetMigrationRow(ctx, fileStruct).Return(nil, nil)
	fileInfo.EXPECT().Name().Return("test5")
	dbMock.EXPECT().GetUpPart(fileStruct).Return(migrateSQLString, nil)

	dbMock.EXPECT().Migrate(ctx, migrateSQLString).Return(nil)
	dbMock.EXPECT().CreateMigration(ctx, fileStruct, actualVersion).Return(expError)

	_, err := usecase.UpMigration(ctx, dbMock, fileStruct, actualVersion)
	assert.ErrorIs(t, err, expError)
}

func TestUpMigration_FileHasBeenChnged(t *testing.T) {
	ctx := context.Background()
	dbMock := mocks.NewDBUsecaseContract(t)

	actualVersion := 2
	now := time.Now()
	nowDB := now.Format(time.RFC3339Nano)

	fileInfo := mocks.NewFileInfo(t)
	fileStruct := usecase.FileStruct{File: fileInfo, Hash: "Hash1"}

	dbMock.EXPECT().GetMigrationRow(ctx, fileStruct).Return(&storage.MigrationDBStruct{
		ID:        1,
		Name:      "test5",
		Hash:      "Hash2",
		Version:   actualVersion,
		CreatedAt: nowDB,
		UpdatedAt: nowDB,
	}, nil)
	fileInfo.EXPECT().Name().Return("test5")

	expOut := usecase.Out{
		Name:   "test5",
		Status: usecase.FileHasBeenChanged,
	}

	out, err := usecase.UpMigration(ctx, dbMock, fileStruct, actualVersion)
	assert.NoError(t, err)
	assert.Equal(t, &expOut, out)
}

func TestUpMigration_NothingMigrate(t *testing.T) {
	ctx := context.Background()
	dbMock := mocks.NewDBUsecaseContract(t)

	actualVersion := 2
	now := time.Now()
	nowDB := now.Format(time.RFC3339Nano)

	fileInfo := mocks.NewFileInfo(t)
	fileStruct := usecase.FileStruct{File: fileInfo, Hash: "Hash1"}

	dbMock.EXPECT().GetMigrationRow(ctx, fileStruct).Return(&storage.MigrationDBStruct{
		ID:        1,
		Name:      "test5",
		Hash:      "Hash1",
		Version:   actualVersion,
		CreatedAt: nowDB,
		UpdatedAt: nowDB,
	}, nil)
	fileInfo.EXPECT().Name().Return("test5")

	out, err := usecase.UpMigration(ctx, dbMock, fileStruct, actualVersion)
	assert.NoError(t, err)
	assert.Nil(t, out)
}
