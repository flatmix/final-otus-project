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

const downSQLString = "test"

func TestDown_Ok(t *testing.T) {
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
	}

	migrations := storage.MigrationsDBStruct{
		{
			ID:        1,
			Name:      "test1",
			Hash:      "Hash1",
			Version:   0,
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

	outs, err := usecase.Down(ctx, ucMock, false, 0)
	assert.NoError(t, err)
	assert.Equal(t, &expOuts, outs)
}

func TestDown_AllOk(t *testing.T) {
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
	}

	migrations := storage.MigrationsDBStruct{
		{
			ID:        1,
			Name:      "test1",
			Hash:      "Hash1",
			Version:   0,
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

	ucMock.EXPECT().GetDownPart(expfiles["test2"]).Return(downSQLString, nil)
	ucMock.EXPECT().Migrate(ctx, downSQLString).Return(nil)
	ucMock.EXPECT().DeleteMigration(ctx, expfiles["test2"]).Return(nil)

	outs, err := usecase.Down(ctx, ucMock, true, 0)
	assert.NoError(t, err)
	assert.Equal(t, &expOuts, outs)
}

func TestDown_FileNotFoundError(t *testing.T) {
	ctx := context.Background()
	ucMock := mocks.NewUCContract(t)

	now := time.Now()
	nowDB := now.Format(time.RFC3339Nano)

	expfiles := make(usecase.FilesMap)

	migrations := storage.MigrationsDBStruct{
		{
			ID:        1,
			Name:      "test1",
			Hash:      "Hash1",
			Version:   0,
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

	_, err := usecase.Down(ctx, ucMock, false, 0)
	assert.ErrorIs(t, err, usecase.ErrNotFountMigrationFile)
}

func TestDown_GetAllMigrationFileMapError(t *testing.T) {
	ctx := context.Background()
	ucMock := mocks.NewUCContract(t)

	expError := errors.New("test error")

	ucMock.EXPECT().GetAllMigrationFileMap().Return(nil, expError)

	_, err := usecase.Down(ctx, ucMock, false, 0)
	assert.ErrorIs(t, err, expError)
}

func TestDown_GetAllMigrationsOrderByVersionDescError(t *testing.T) {
	ctx := context.Background()
	ucMock := mocks.NewUCContract(t)

	expfiles := make(usecase.FilesMap)

	expError := errors.New("test error")

	ucMock.EXPECT().GetAllMigrationFileMap().Return(expfiles, nil)
	ucMock.EXPECT().GetAllMigrationsOrderByVersionDesc(ctx, 0).Return(nil, expError)

	_, err := usecase.Down(ctx, ucMock, false, 0)
	assert.ErrorIs(t, err, expError)
}

func TestDown_ErrNothingForDownMigrateError(t *testing.T) {
	ctx := context.Background()
	ucMock := mocks.NewUCContract(t)

	expfiles := make(usecase.FilesMap)

	migrations := storage.MigrationsDBStruct{}

	ucMock.EXPECT().GetAllMigrationFileMap().Return(expfiles, nil)
	ucMock.EXPECT().GetAllMigrationsOrderByVersionDesc(ctx, 0).Return(migrations, nil)

	_, err := usecase.Down(ctx, ucMock, false, 0)
	assert.ErrorIs(t, err, usecase.ErrNothingForDownMigrate)
}

func TestDown_GetDownPartError(t *testing.T) {
	ctx := context.Background()
	ucMock := mocks.NewUCContract(t)

	now := time.Now()
	nowDB := now.Format(time.RFC3339Nano)

	fileInfo := mocks.NewFileInfo(t)
	fileInfo2 := mocks.NewFileInfo(t)

	expfiles := make(usecase.FilesMap)
	expfiles["test1"] = usecase.FileStruct{File: fileInfo, Hash: "Hash1"}
	expfiles["test2"] = usecase.FileStruct{File: fileInfo2, Hash: "Hash2"}

	migrations := storage.MigrationsDBStruct{
		{
			ID:        1,
			Name:      "test1",
			Hash:      "Hash1",
			Version:   0,
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
	ucMock.EXPECT().GetDownPart(expfiles["test1"]).Return("", expError)

	_, err := usecase.Down(ctx, ucMock, true, 0)
	assert.ErrorIs(t, err, expError)
}

func TestDown_MigrateError(t *testing.T) {
	ctx := context.Background()
	ucMock := mocks.NewUCContract(t)

	now := time.Now()
	nowDB := now.Format(time.RFC3339Nano)

	fileInfo := mocks.NewFileInfo(t)
	fileInfo2 := mocks.NewFileInfo(t)

	expfiles := make(usecase.FilesMap)
	expfiles["test1"] = usecase.FileStruct{File: fileInfo, Hash: "Hash1"}
	expfiles["test2"] = usecase.FileStruct{File: fileInfo2, Hash: "Hash2"}

	migrations := storage.MigrationsDBStruct{
		{
			ID:        1,
			Name:      "test1",
			Hash:      "Hash1",
			Version:   0,
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
	ucMock.EXPECT().Migrate(ctx, downSQLString).Return(expError)

	_, err := usecase.Down(ctx, ucMock, false, 0)
	assert.ErrorIs(t, err, expError)
}

func TestDown_DeleteMigrationError(t *testing.T) {
	ctx := context.Background()
	ucMock := mocks.NewUCContract(t)

	now := time.Now()
	nowDB := now.Format(time.RFC3339Nano)

	fileInfo := mocks.NewFileInfo(t)
	fileInfo2 := mocks.NewFileInfo(t)

	expfiles := make(usecase.FilesMap)
	expfiles["test1"] = usecase.FileStruct{File: fileInfo, Hash: "Hash1"}
	expfiles["test2"] = usecase.FileStruct{File: fileInfo2, Hash: "Hash2"}

	migrations := storage.MigrationsDBStruct{
		{
			ID:        1,
			Name:      "test1",
			Hash:      "Hash1",
			Version:   0,
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
	ucMock.EXPECT().DeleteMigration(ctx, expfiles["test1"]).Return(expError)

	_, err := usecase.Down(ctx, ucMock, false, 0)
	assert.ErrorIs(t, err, expError)
}
