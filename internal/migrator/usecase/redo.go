package usecase

import (
	"context"
	"database/sql"
	"fmt"
	"sort"

	"github.com/flatmix/final-otus-project/internal/migrator/storage"
	_ "github.com/lib/pq" //nolint:revive,nolintlint
)

func Redo(ctx context.Context, db *sql.DB, all bool, step int) (*Outs, error) {
	redoStruct := NewDBStruct(db)

	filesMap, err := redoStruct.getAllMigrationFileMap()
	if err != nil {
		return nil, fmt.Errorf("getAllMigrationFileMap: %w", err)
	}

	migrations, err := redoStruct.getAllMigrationsOrderByVersionDesc(ctx, step)
	if err != nil {
		return nil, fmt.Errorf("getAllMigrationsOrderByVersionDesc: %w", err)
	}

	if len(migrations) == 0 {
		return nil, ErrNotFountMigrationFile
	}

	outs := &Outs{} //nolint:ineffassign

	if all || step > 0 {
		outs, err = redoStruct.redoMigration(ctx, migrations, filesMap)
		if err != nil {
			return nil, fmt.Errorf("redoMigration: %w", err)
		}
	} else {
		outs, err = redoStruct.redoMigration(ctx, migrations[0:1], filesMap)
		if err != nil {
			return nil, fmt.Errorf("redoMigration: %w", err)
		}
	}

	return outs, nil
}

func (ds *DB) redoMigration(ctx context.Context,
	migrations storage.MigrationsDBStruct, filesMap FilesMap,
) (*Outs, error) {
	outs := make(Outs, 0)

	for _, migration := range migrations {
		downOut, err := ds.downMigration(ctx, migration, filesMap)
		if err != nil {
			return nil, err
		}
		outs = append(outs, downOut)
	}

	sort.Sort(migrations)

	for _, migration := range migrations {
		file, ok := filesMap[migration.Name]
		if !ok {
			return nil, fmt.Errorf("not found migration file - '%s'", migration.Name)
		}
		upOuts, err := ds.upMigration(ctx, file, migration.Version)
		if err != nil {
			return nil, err
		}
		outs = append(outs, upOuts)
	}
	return &outs, nil
}
