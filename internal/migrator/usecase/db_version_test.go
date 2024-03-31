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
	ucMock := mocks.NewUCContract(t)

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

	ucMock.EXPECT().GetAllMigrationsOrderByVersionDesc(ctx, 0).Return(migrationExp, nil)

	version, err := usecase.DBVersion(ctx, ucMock)
	assert.NoError(t, err)
	assert.Equal(t, migrationExp[0].Version, *version)
}

func TestDBVersion_Fail(t *testing.T) {
	ctx := context.Background()
	ucMock := mocks.NewUCContract(t)

	expectedError := errors.New("test errors")

	ucMock.EXPECT().GetAllMigrationsOrderByVersionDesc(ctx, 0).Return(nil, expectedError)

	_, err := usecase.DBVersion(ctx, ucMock)
	assert.ErrorIs(t, err, expectedError)
}
