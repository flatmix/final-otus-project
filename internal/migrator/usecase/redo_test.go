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

func TestRedo_Ok(t *testing.T) {
	ctx := context.Background()
	dbMock := mocks.NewDBUsecaseContract(t)

	now := time.Now()
	nowDB := now.Format(time.RFC3339Nano)

	fileInfo := mocks.NewFileInfo(t)
	fileInfo2 := mocks.NewFileInfo(t)

	expfiles := make(usecase.FilesMap)
	expfiles["test1"] = usecase.FileStruct{File: fileInfo, Hash: "Hash1"}
	expfiles["test2"] = usecase.FileStruct{File: fileInfo2, Hash: "Hash2"}

	expOuts := usecase.Outs{
		&usecase.Out{
			Name:   "test1",
			Status: "down ok",
		},
		&usecase.Out{
			Name:   "test1",
			Status: "migrate ok",
		},
	}

	actualVersion := 1

	migrations := storage.MigrationsDBStruct{
		{
			ID:        1,
			Name:      "test1",
			Hash:      "Hash1",
			Version:   actualVersion,
			CreatedAt: nowDB,
			UpdatedAt: nowDB,
		},
		{
			ID:        2,
			Name:      "test2",
			Hash:      "Hash2",
			Version:   0,
			CreatedAt: nowDB,
			UpdatedAt: nowDB,
		},
	}

	dbMock.EXPECT().GetAllMigrationFileMap().Return(expfiles, nil)
	dbMock.EXPECT().GetAllMigrationsOrderByVersionDesc(ctx, 0).Return(migrations, nil)

	dbMock.EXPECT().GetDownPart(expfiles["test1"]).Return(downSQLString, nil)
	dbMock.EXPECT().Migrate(ctx, downSQLString).Return(nil)
	dbMock.EXPECT().DeleteMigration(ctx, expfiles["test1"]).Return(nil)

	dbMock.EXPECT().GetMigrationRow(ctx, expfiles["test1"]).Return(nil, nil)
	fileInfo.EXPECT().Name().Return("test1").Once()
	dbMock.EXPECT().GetUpPart(expfiles["test1"]).Return(migrateSQLString, nil)
	dbMock.EXPECT().Migrate(ctx, migrateSQLString).Return(nil)
	dbMock.EXPECT().CreateMigration(ctx, expfiles["test1"], actualVersion).Return(nil)

	outs, err := usecase.Redo(ctx, dbMock, false, 0)
	assert.NoError(t, err)
	assert.Equal(t, &expOuts, outs)
}

func TestRedo_AllOk(t *testing.T) {
	ctx := context.Background()
	dbMock := mocks.NewDBUsecaseContract(t)

	now := time.Now()
	nowDB := now.Format(time.RFC3339Nano)

	fileInfo := mocks.NewFileInfo(t)
	fileInfo2 := mocks.NewFileInfo(t)

	expfiles := make(usecase.FilesMap)
	expfiles["test1"] = usecase.FileStruct{File: fileInfo, Hash: "Hash1"}
	expfiles["test2"] = usecase.FileStruct{File: fileInfo2, Hash: "Hash2"}

	expOuts := usecase.Outs{
		&usecase.Out{
			Name:   "test1",
			Status: "down ok",
		},
		&usecase.Out{
			Name:   "test2",
			Status: "down ok",
		},
		&usecase.Out{
			Name:   "test1",
			Status: "migrate ok",
		},
		&usecase.Out{
			Name:   "test2",
			Status: "migrate ok",
		},
	}

	actualVersion := 1

	migrations := storage.MigrationsDBStruct{
		{
			ID:        1,
			Name:      "test1",
			Hash:      "Hash1",
			Version:   actualVersion,
			CreatedAt: nowDB,
			UpdatedAt: nowDB,
		},
		{
			ID:        2,
			Name:      "test2",
			Hash:      "Hash2",
			Version:   actualVersion,
			CreatedAt: nowDB,
			UpdatedAt: nowDB,
		},
	}

	dbMock.EXPECT().GetAllMigrationFileMap().Return(expfiles, nil)
	dbMock.EXPECT().GetAllMigrationsOrderByVersionDesc(ctx, 0).Return(migrations, nil)

	dbMock.EXPECT().GetDownPart(expfiles["test1"]).Return(downSQLString, nil)
	dbMock.EXPECT().Migrate(ctx, downSQLString).Return(nil)
	dbMock.EXPECT().DeleteMigration(ctx, expfiles["test1"]).Return(nil)

	dbMock.EXPECT().GetDownPart(expfiles["test2"]).Return(downSQLString, nil)
	dbMock.EXPECT().Migrate(ctx, downSQLString).Return(nil)
	dbMock.EXPECT().DeleteMigration(ctx, expfiles["test2"]).Return(nil)

	dbMock.EXPECT().GetMigrationRow(ctx, expfiles["test2"]).Return(nil, nil)
	fileInfo2.EXPECT().Name().Return("test2")
	dbMock.EXPECT().GetUpPart(expfiles["test2"]).Return(migrateSQLString, nil)
	dbMock.EXPECT().Migrate(ctx, migrateSQLString).Return(nil)
	dbMock.EXPECT().CreateMigration(ctx, expfiles["test2"], actualVersion).Return(nil)

	dbMock.EXPECT().GetMigrationRow(ctx, expfiles["test1"]).Return(nil, nil)
	fileInfo.EXPECT().Name().Return("test1")
	dbMock.EXPECT().GetUpPart(expfiles["test1"]).Return(migrateSQLString, nil)
	dbMock.EXPECT().Migrate(ctx, migrateSQLString).Return(nil)
	dbMock.EXPECT().CreateMigration(ctx, expfiles["test1"], actualVersion).Return(nil)

	outs, err := usecase.Redo(ctx, dbMock, true, 0)
	assert.NoError(t, err)
	assert.Equal(t, &expOuts, outs)
}

func TestRedo_UpMigrationError(t *testing.T) {
	ctx := context.Background()
	dbMock := mocks.NewDBUsecaseContract(t)

	now := time.Now()
	nowDB := now.Format(time.RFC3339Nano)

	fileInfo := mocks.NewFileInfo(t)
	fileInfo2 := mocks.NewFileInfo(t)

	expfiles := make(usecase.FilesMap)
	expfiles["test1"] = usecase.FileStruct{File: fileInfo, Hash: "Hash1"}
	expfiles["test2"] = usecase.FileStruct{File: fileInfo2, Hash: "Hash2"}

	actualVersion := 1

	migrations := storage.MigrationsDBStruct{
		{
			ID:        1,
			Name:      "test1",
			Hash:      "Hash1",
			Version:   actualVersion,
			CreatedAt: nowDB,
			UpdatedAt: nowDB,
		},
		{
			ID:        2,
			Name:      "test2",
			Hash:      "Hash2",
			Version:   0,
			CreatedAt: nowDB,
			UpdatedAt: nowDB,
		},
	}

	expError := errors.New("test error")

	dbMock.EXPECT().GetAllMigrationFileMap().Return(expfiles, nil)
	dbMock.EXPECT().GetAllMigrationsOrderByVersionDesc(ctx, 0).Return(migrations, nil)

	dbMock.EXPECT().GetDownPart(expfiles["test1"]).Return(downSQLString, nil)
	dbMock.EXPECT().Migrate(ctx, downSQLString).Return(nil)
	dbMock.EXPECT().DeleteMigration(ctx, expfiles["test1"]).Return(nil)

	dbMock.EXPECT().GetMigrationRow(ctx, expfiles["test1"]).Return(nil, expError)

	_, err := usecase.Redo(ctx, dbMock, false, 0)
	assert.ErrorIs(t, err, expError)
}

func TestRedo_NotFountMigrationFileError(t *testing.T) {
	ctx := context.Background()
	dbMock := mocks.NewDBUsecaseContract(t)

	now := time.Now()
	nowDB := now.Format(time.RFC3339Nano)

	fileInfo := mocks.NewFileInfo(t)
	fileInfo2 := mocks.NewFileInfo(t)

	expfiles := make(usecase.FilesMap)

	expfilesFail := make(usecase.FilesMap)

	expfiles["test1"] = usecase.FileStruct{File: fileInfo, Hash: "Hash1"}
	expfiles["test2"] = usecase.FileStruct{File: fileInfo2, Hash: "Hash2"}

	actualVersion := 1

	migrations := storage.MigrationsDBStruct{
		{
			ID:        1,
			Name:      "test1",
			Hash:      "Hash1",
			Version:   actualVersion,
			CreatedAt: nowDB,
			UpdatedAt: nowDB,
		},
		{
			ID:        2,
			Name:      "test2",
			Hash:      "Hash2",
			Version:   0,
			CreatedAt: nowDB,
			UpdatedAt: nowDB,
		},
	}

	dbMock.EXPECT().GetAllMigrationFileMap().Return(expfilesFail, nil)
	dbMock.EXPECT().GetAllMigrationsOrderByVersionDesc(ctx, 0).Return(migrations, nil)

	_, err := usecase.Redo(ctx, dbMock, true, 0)
	assert.ErrorIs(t, err, usecase.ErrNotFountMigrationFile)
}

func TestRedo_NotFountMigrationError(t *testing.T) {
	ctx := context.Background()
	dbMock := mocks.NewDBUsecaseContract(t)

	fileInfo := mocks.NewFileInfo(t)
	fileInfo2 := mocks.NewFileInfo(t)

	expfiles := make(usecase.FilesMap)
	expfiles["test1"] = usecase.FileStruct{File: fileInfo, Hash: "Hash1"}
	expfiles["test2"] = usecase.FileStruct{File: fileInfo2, Hash: "Hash2"}

	migrations := storage.MigrationsDBStruct{}

	dbMock.EXPECT().GetAllMigrationFileMap().Return(expfiles, nil)
	dbMock.EXPECT().GetAllMigrationsOrderByVersionDesc(ctx, 0).Return(migrations, nil)

	_, err := usecase.Redo(ctx, dbMock, false, 0)
	assert.ErrorIs(t, err, usecase.ErrNotFountMigration)
}

func TestRedo_GetAllMigrationsOrderByVersionDescError(t *testing.T) {
	ctx := context.Background()
	dbMock := mocks.NewDBUsecaseContract(t)

	fileInfo := mocks.NewFileInfo(t)
	fileInfo2 := mocks.NewFileInfo(t)

	expfiles := make(usecase.FilesMap)
	expfiles["test1"] = usecase.FileStruct{File: fileInfo, Hash: "Hash1"}
	expfiles["test2"] = usecase.FileStruct{File: fileInfo2, Hash: "Hash2"}

	expError := errors.New("test error")

	dbMock.EXPECT().GetAllMigrationFileMap().Return(expfiles, nil)
	dbMock.EXPECT().GetAllMigrationsOrderByVersionDesc(ctx, 0).Return(nil, expError)

	_, err := usecase.Redo(ctx, dbMock, false, 0)
	assert.ErrorIs(t, err, expError)
}

func TestRedo_GetAllMigrationFileMapError(t *testing.T) {
	ctx := context.Background()
	dbMock := mocks.NewDBUsecaseContract(t)

	fileInfo := mocks.NewFileInfo(t)

	expfiles := make(usecase.FilesMap)
	expfiles["test1"] = usecase.FileStruct{File: fileInfo, Hash: "Hash1"}

	expError := errors.New("test error")

	dbMock.EXPECT().GetAllMigrationFileMap().Return(nil, expError)

	_, err := usecase.Redo(ctx, dbMock, false, 0)
	assert.ErrorIs(t, err, expError)
}
