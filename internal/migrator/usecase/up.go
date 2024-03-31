package usecase

import (
	"context"

	_ "github.com/lib/pq" //nolint:revive,nolintlint
)

const FileHasBeenChanged = "the migration file has been changed, attention!!!"

func Up(ctx context.Context, db UCContract) (*Outs, error) {
	files, err := db.GetAllMigrationFile()
	if err != nil {
		return nil, err
	}

	createSQLErr := db.CreateMigrationsTable(ctx)
	if createSQLErr != nil {
		return nil, createSQLErr
	}

	actualVersion := db.GetActualVersion(ctx)

	outUps := make(Outs, 0)
	for _, file := range files {
		outUp, err := UpMigration(ctx, db, file, actualVersion)
		if err != nil {
			return nil, err
		}
		if outUp != nil {
			outUps = append(outUps, outUp)
		}
	}

	return &outUps, nil
}

func UpMigration(ctx context.Context, ds UCContract, file FileStruct, actualVersion int) (*Out, error) {
	migrationDB, err := ds.GetMigrationRow(ctx, file)
	if err != nil {
		return nil, err
	}

	outUp := Out{
		Name:   file.File.Name(),
		Status: "",
	}
	if migrationDB == nil || migrationDB.IsZero() {
		migrateSQLString, err := ds.GetUpPart(file)
		if err != nil {
			return nil, err
		}

		err = ds.Migrate(ctx, migrateSQLString)
		if err != nil {
			return nil, err
		}
		err = ds.CreateMigration(ctx, file, actualVersion)
		if err != nil {
			return nil, err
		}
		outUp.Status = "migrate ok"
		return &outUp, nil
	} else if file.Hash != migrationDB.Hash {
		outUp.Status = FileHasBeenChanged
		return &outUp, nil
	}

	return nil, nil
}
