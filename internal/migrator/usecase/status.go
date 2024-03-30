package usecase

import (
	"context"
	"strconv"
	"time"
)

func Status(ctx context.Context, dbStruct DBUsecaseContract) (*Outs, error) {
	files, err := dbStruct.GetAllMigrationFile()
	if err != nil {
		return nil, err
	}

	statusStructs := make(Outs, 0)

	for _, file := range files {
		migrationDB, err := dbStruct.GetMigrationRow(ctx, file)
		if err != nil {
			statusStructs = append(statusStructs, &Out{
				Name:        file.File.Name(),
				Status:      "No migrate",
				Version:     "-",
				TimeMigrate: "-",
			})
			continue
		}
		status := "Ok"
		version := "-"
		timeMigrate := "-"
		if migrationDB == nil || migrationDB.IsZero() {
			status = "No migrate"
		} else {
			version = strconv.Itoa(migrationDB.Version)
			timeMigrateDB, _ := time.Parse(time.RFC3339Nano, migrationDB.CreatedAt)
			timeMigrate = timeMigrateDB.Format(time.DateTime)
		}

		statusStructs = append(statusStructs, &Out{
			Name:        file.File.Name(),
			Status:      status,
			Version:     version,
			TimeMigrate: timeMigrate,
		})
	}

	return &statusStructs, nil
}
