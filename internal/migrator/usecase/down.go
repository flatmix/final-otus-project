package usecase

import (
	"context"
	"errors"

	"github.com/flatmix/final-otus-project/internal/migrator/storage"
)

var (
	ErrNotFountMigrationFile = errors.New("not found migration file")
	ErrNotFountMigration     = errors.New("not found migrations in db")
	ErrNothingForDownMigrate = errors.New("nothing for down migrate")
)

func Down(ctx context.Context, downStruct UCContract, all bool, step int) (*Outs, error) {
	filesMap, err := downStruct.GetAllMigrationFileMap()
	if err != nil {
		return nil, err
	}

	migrations, err := downStruct.GetAllMigrationsOrderByVersionDesc(ctx, step)
	if err != nil {
		return nil, err
	}

	if len(migrations) == 0 {
		return nil, ErrNothingForDownMigrate
	}

	outs := make(Outs, 0)

	if all || step > 0 {
		for _, migration := range migrations {
			out, err := DownMigration(ctx, downStruct, migration, filesMap)
			if err != nil {
				return nil, err
			}
			if out != nil {
				outs = append(outs, out)
			}
		}
	}

	if !all && step == 0 {
		out, err := DownMigration(ctx, downStruct, migrations[0], filesMap)
		if err != nil {
			return nil, err
		}
		if out != nil {
			outs = append(outs, out)
		}
	}

	return &outs, nil
}

func DownMigration(ctx context.Context, ds UCContract,
	migration storage.MigrationDBStruct, filesMap FilesMap,
) (*Out, error) {
	file, ok := filesMap[migration.Name]
	out := Out{
		Name:   migration.Name,
		Status: "down start",
	}
	if !ok {
		out.Status = "not found migration file"
		return nil, ErrNotFountMigrationFile
	}
	migrateSQLString, err := ds.GetDownPart(file)
	if err != nil {
		return nil, err
	}

	err = ds.Migrate(ctx, migrateSQLString)
	if err != nil {
		return nil, err
	}
	err = ds.DeleteMigration(ctx, file)
	if err != nil {
		return nil, err
	}

	out.Status = "down ok"

	return &out, nil
}
