package usecase

import (
	"context"
	"fmt"
	"sort"

	"github.com/flatmix/final-otus-project/internal/migrator/storage"
	_ "github.com/lib/pq" //nolint:revive,nolintlint
)

func Redo(ctx context.Context, redoStruct DBUsecaseContract, all bool, step int) (*Outs, error) {
	filesMap, err := redoStruct.GetAllMigrationFileMap()
	if err != nil {
		return nil, fmt.Errorf("getAllMigrationFileMap: %w", err)
	}

	migrations, err := redoStruct.GetAllMigrationsOrderByVersionDesc(ctx, step)
	if err != nil {
		return nil, fmt.Errorf("getAllMigrationsOrderByVersionDesc: %w", err)
	}

	if len(migrations) == 0 {
		return nil, ErrNotFountMigrationFile
	}

	var outs Outs

	if all || step > 0 {
		outsRedo, err := redoStruct.RedoMigration(ctx, migrations, filesMap)
		if err != nil {
			return nil, fmt.Errorf("redoMigration: %w", err)
		}
		outs = *outsRedo
	} else {
		outsRedo, err := redoStruct.RedoMigration(ctx, migrations[0:1], filesMap)
		if err != nil {
			return nil, fmt.Errorf("redoMigration: %w", err)
		}
		outs = *outsRedo
	}

	return &outs, nil
}

func (ds *DB) RedoMigration(ctx context.Context,
	migrations storage.MigrationsDBStruct, filesMap FilesMap,
) (*Outs, error) {
	outs := make(Outs, 0)

	for _, migration := range migrations {
		downOut, err := ds.DownMigration(ctx, migration, filesMap)
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
		upOuts, err := UpMigration(ctx, ds, file, migration.Version)
		if err != nil {
			return nil, err
		}
		outs = append(outs, upOuts)
	}
	return &outs, nil
}
