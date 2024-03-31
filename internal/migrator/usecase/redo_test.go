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
	ucMock := mocks.NewUCContract(t)

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

	ucMock.EXPECT().GetAllMigrationFileMap().Return(expfiles, nil)
	ucMock.EXPECT().GetAllMigrationsOrderByVersionDesc(ctx, 0).Return(migrations, nil)

	ucMock.EXPECT().GetDownPart(expfiles["test1"]).Return(downSQLString, nil)
	ucMock.EXPECT().Migrate(ctx, downSQLString).Return(nil)
	ucMock.EXPECT().DeleteMigration(ctx, expfiles["test1"]).Return(nil)

	ucMock.EXPECT().GetMigrationRow(ctx, expfiles["test1"]).Return(nil, nil)
	fileInfo.EXPECT().Name().Return("test1").Once()
	ucMock.EXPECT().GetUpPart(expfiles["test1"]).Return(migrateSQLString, nil)
	ucMock.EXPECT().Migrate(ctx, migrateSQLString).Return(nil)
	ucMock.EXPECT().CreateMigration(ctx, expfiles["test1"], actualVersion).Return(nil)

	outs, err := usecase.Redo(ctx, ucMock, false, 0)
	assert.NoError(t, err)
	assert.Equal(t, &expOuts, outs)
}

func TestRedo_AllOk(t *testing.T) {
	ctx := context.Background()
	ucMock := mocks.NewUCContract(t)

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

	ucMock.EXPECT().GetAllMigrationFileMap().Return(expfiles, nil)
	ucMock.EXPECT().GetAllMigrationsOrderByVersionDesc(ctx, 0).Return(migrations, nil)

	ucMock.EXPECT().GetDownPart(expfiles["test1"]).Return(downSQLString, nil)
	ucMock.EXPECT().Migrate(ctx, downSQLString).Return(nil)
	ucMock.EXPECT().DeleteMigration(ctx, expfiles["test1"]).Return(nil)

	ucMock.EXPECT().GetDownPart(expfiles["test2"]).Return(downSQLString, nil)
	ucMock.EXPECT().Migrate(ctx, downSQLString).Return(nil)
	ucMock.EXPECT().DeleteMigration(ctx, expfiles["test2"]).Return(nil)

	ucMock.EXPECT().GetMigrationRow(ctx, expfiles["test2"]).Return(nil, nil)
	fileInfo2.EXPECT().Name().Return("test2")
	ucMock.EXPECT().GetUpPart(expfiles["test2"]).Return(migrateSQLString, nil)
	ucMock.EXPECT().Migrate(ctx, migrateSQLString).Return(nil)
	ucMock.EXPECT().CreateMigration(ctx, expfiles["test2"], actualVersion).Return(nil)

	ucMock.EXPECT().GetMigrationRow(ctx, expfiles["test1"]).Return(nil, nil)
	fileInfo.EXPECT().Name().Return("test1")
	ucMock.EXPECT().GetUpPart(expfiles["test1"]).Return(migrateSQLString, nil)
	ucMock.EXPECT().Migrate(ctx, migrateSQLString).Return(nil)
	ucMock.EXPECT().CreateMigration(ctx, expfiles["test1"], actualVersion).Return(nil)

	outs, err := usecase.Redo(ctx, ucMock, true, 0)
	assert.NoError(t, err)
	assert.Equal(t, &expOuts, outs)
}

func TestRedo_UpMigrationError(t *testing.T) {
	ctx := context.Background()
	ucMock := mocks.NewUCContract(t)

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

	ucMock.EXPECT().GetAllMigrationFileMap().Return(expfiles, nil)
	ucMock.EXPECT().GetAllMigrationsOrderByVersionDesc(ctx, 0).Return(migrations, nil)

	ucMock.EXPECT().GetDownPart(expfiles["test1"]).Return(downSQLString, nil)
	ucMock.EXPECT().Migrate(ctx, downSQLString).Return(nil)
	ucMock.EXPECT().DeleteMigration(ctx, expfiles["test1"]).Return(nil)

	ucMock.EXPECT().GetMigrationRow(ctx, expfiles["test1"]).Return(nil, expError)

	_, err := usecase.Redo(ctx, ucMock, false, 0)
	assert.ErrorIs(t, err, expError)
}

func TestRedo_NotFountMigrationFileError(t *testing.T) {
	ctx := context.Background()
	ucMock := mocks.NewUCContract(t)

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

	ucMock.EXPECT().GetAllMigrationFileMap().Return(expfilesFail, nil)
	ucMock.EXPECT().GetAllMigrationsOrderByVersionDesc(ctx, 0).Return(migrations, nil)

	_, err := usecase.Redo(ctx, ucMock, true, 0)
	assert.ErrorIs(t, err, usecase.ErrNotFountMigrationFile)
}

func TestRedo_NotFountMigrationError(t *testing.T) {
	ctx := context.Background()
	ucMock := mocks.NewUCContract(t)

	fileInfo := mocks.NewFileInfo(t)
	fileInfo2 := mocks.NewFileInfo(t)

	expfiles := make(usecase.FilesMap)
	expfiles["test1"] = usecase.FileStruct{File: fileInfo, Hash: "Hash1"}
	expfiles["test2"] = usecase.FileStruct{File: fileInfo2, Hash: "Hash2"}

	migrations := storage.MigrationsDBStruct{}

	ucMock.EXPECT().GetAllMigrationFileMap().Return(expfiles, nil)
	ucMock.EXPECT().GetAllMigrationsOrderByVersionDesc(ctx, 0).Return(migrations, nil)

	_, err := usecase.Redo(ctx, ucMock, false, 0)
	assert.ErrorIs(t, err, usecase.ErrNotFountMigration)
}

func TestRedo_GetAllMigrationsOrderByVersionDescError(t *testing.T) {
	ctx := context.Background()
	ucMock := mocks.NewUCContract(t)

	fileInfo := mocks.NewFileInfo(t)
	fileInfo2 := mocks.NewFileInfo(t)

	expfiles := make(usecase.FilesMap)
	expfiles["test1"] = usecase.FileStruct{File: fileInfo, Hash: "Hash1"}
	expfiles["test2"] = usecase.FileStruct{File: fileInfo2, Hash: "Hash2"}

	expError := errors.New("test error")

	ucMock.EXPECT().GetAllMigrationFileMap().Return(expfiles, nil)
	ucMock.EXPECT().GetAllMigrationsOrderByVersionDesc(ctx, 0).Return(nil, expError)

	_, err := usecase.Redo(ctx, ucMock, false, 0)
	assert.ErrorIs(t, err, expError)
}

func TestRedo_GetAllMigrationFileMapError(t *testing.T) {
	ctx := context.Background()
	ucMock := mocks.NewUCContract(t)

	fileInfo := mocks.NewFileInfo(t)

	expfiles := make(usecase.FilesMap)
	expfiles["test1"] = usecase.FileStruct{File: fileInfo, Hash: "Hash1"}

	expError := errors.New("test error")

	ucMock.EXPECT().GetAllMigrationFileMap().Return(nil, expError)

	_, err := usecase.Redo(ctx, ucMock, false, 0)
	assert.ErrorIs(t, err, expError)
}
