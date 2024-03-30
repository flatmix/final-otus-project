package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/flatmix/final-otus-project/internal/migrator/storage"
	"github.com/flatmix/final-otus-project/internal/migrator/usecase"
	mocks "github.com/flatmix/final-otus-project/mocks/migrator/usecase"
	"github.com/stretchr/testify/assert"
)

func TestDBVersion_OK(t *testing.T) {
	ctx := context.Background()
	dbStruct := mocks.NewDBUsecaseContract(t)

	migrationExp := []storage.MigrationDBStruct{
		{
			ID:        0,
			Name:      "Test",
			Hash:      "hash",
			Version:   35,
			CreatedAt: "",
			UpdatedAt: "",
		},
	}

	dbStruct.EXPECT().GetAllMigrationsOrderByVersionDesc(ctx, 0).Return(migrationExp, nil)

	version, err := usecase.DBVersion(ctx, dbStruct)
	assert.NoError(t, err)
	assert.Equal(t, migrationExp[0].Version, *version)
}

func TestDBVersion_Fail(t *testing.T) {
	ctx := context.Background()
	dbStruct := mocks.NewDBUsecaseContract(t)

	expectedError := errors.New("test errors")

	dbStruct.EXPECT().GetAllMigrationsOrderByVersionDesc(ctx, 0).Return(nil, expectedError)

	_, err := usecase.DBVersion(ctx, dbStruct)
	assert.ErrorIs(t, err, expectedError)
}
