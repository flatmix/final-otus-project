package usecase

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/flatmix/final-otus-project/internal/migrator/config"
	"github.com/flatmix/final-otus-project/internal/migrator/storage"
)

func Status(ctx context.Context, dbConf config.Postgres) (*Outs, error) {
	db, err := storage.NewDB(dbConf)
	if err != nil {
		return nil, fmt.Errorf("usecase.Status: %w", err)
	}

	dbStruct := NewDBStruct(db)

	files, err := dbStruct.getAllMigrationFile()
	if err != nil {
		return nil, err
	}

	statusStructs := make(Outs, 0)

	for _, file := range files {
		migrationDB, err := dbStruct.getMigrationRow(ctx, file)
		if err != nil {
			statusStructs = append(statusStructs, &Out{
				Name:        file.file.Name(),
				Status:      "No migrate",
				Version:     "-",
				TimeMigrate: "-",
			})
		}
		status := "Ok"
		version := "-"
		timeMigrate := "-"
		if migrationDB.IsZero() {
			status = "No migrate"
		} else {
			version = strconv.Itoa(migrationDB.Version)
			timeMigrateDB, _ := time.Parse(time.RFC3339Nano, migrationDB.CreatedAt)
			timeMigrate = timeMigrateDB.Format(time.DateTime)
		}

		statusStructs = append(statusStructs, &Out{
			Name:        file.file.Name(),
			Status:      status,
			Version:     version,
			TimeMigrate: timeMigrate,
		})
	}

	return &statusStructs, nil
}
