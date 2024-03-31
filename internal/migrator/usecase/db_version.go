package usecase

import (
	"context"
	"fmt"
)

func DBVersion(ctx context.Context, statusStruct UCContract) (*int, error) {
	migrations, err := statusStruct.GetAllMigrationsOrderByVersionDesc(ctx, 0)
	if err != nil {
		return nil, fmt.Errorf("getAllMigrationsOrderByVersionDesc: %w", err)
	}

	return &migrations[0].Version, nil
}
