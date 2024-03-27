package usecase

import (
	"context"
	"database/sql"
	"fmt"
)

func DBVersion(ctx context.Context, db *sql.DB) (*int, error) {
	statusStruct := NewDBStruct(db)

	migrations, err := statusStruct.getAllMigrationsOrderByVersionDesc(ctx, 0)
	if err != nil {
		return nil, fmt.Errorf("getAllMigrationsOrderByVersionDesc: %w", err)
	}

	return &migrations[0].Version, nil
}
