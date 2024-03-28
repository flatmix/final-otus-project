package migrator

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/flatmix/final-otus-project/internal/migrator/usecase"
)

func Create(name string) error {
	err := usecase.Create(name)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	return nil
}

func Up(ctx context.Context, db *sql.DB) (*usecase.Outs, error) {
	outs, err := usecase.Up(ctx, db)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	return outs, nil
}

func Down(ctx context.Context, db *sql.DB, all bool, step int) (*usecase.Outs, error) {
	outs, err := usecase.Down(ctx, db, all, step)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	return outs, nil
}

func Redo(ctx context.Context, db *sql.DB, all bool, step int) (*usecase.Outs, error) {
	outs, err := usecase.Redo(ctx, db, all, step)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	return outs, nil
}

func Status(ctx context.Context, db *sql.DB) (*usecase.Outs, error) {
	outs, err := usecase.Status(ctx, db)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	return outs, nil
}

func DBVersion(ctx context.Context, db *sql.DB) (*int, error) {
	version, err := usecase.DBVersion(ctx, db)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	return version, nil
}
