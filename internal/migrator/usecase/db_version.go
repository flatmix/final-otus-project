package usecase

import (
	"context"
	"fmt"

	"github.com/flatmix/final-otus-project/internal/migrator/config"
	"github.com/flatmix/final-otus-project/internal/migrator/storage"
)

func DBVersion(ctx context.Context, dbConf config.Postgres) (*int, error) {
	db, err := storage.NewDB(dbConf)
	if err != nil {
		return nil, fmt.Errorf("usecase.DBVersion: %w", err)
	}

	statusStruct := NewDBStruct(db)

	migrations, err := statusStruct.getAllMigrationsOrderByVersionDesc(ctx, 0)
	if err != nil {
		return nil, fmt.Errorf("getAllMigrationsOrderByVersionDesc: %w", err)
	}

	return &migrations[0].Version, nil
}
