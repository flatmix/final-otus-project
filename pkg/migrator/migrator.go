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
	dbStruct := usecase.NewDBStruct(db)
	outs, err := usecase.Up(ctx, dbStruct)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	return outs, nil
}

func Down(ctx context.Context, db *sql.DB, all bool, step int) (*usecase.Outs, error) {
	dbStruct := usecase.NewDBStruct(db)
	outs, err := usecase.Down(ctx, dbStruct, all, step)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	return outs, nil
}

func Redo(ctx context.Context, db *sql.DB, all bool, step int) (*usecase.Outs, error) {
	dbStruct := usecase.NewDBStruct(db)
	outs, err := usecase.Redo(ctx, dbStruct, all, step)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	return outs, nil
}

func Status(ctx context.Context, db *sql.DB) (*usecase.Outs, error) {
	dbStruct := usecase.NewDBStruct(db)
	outs, err := usecase.Status(ctx, dbStruct)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	return outs, nil
}

func DBVersion(ctx context.Context, db *sql.DB) (*int, error) {
	dbStruct := usecase.NewDBStruct(db)
	version, err := usecase.DBVersion(ctx, dbStruct)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	return version, nil
}
