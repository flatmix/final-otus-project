package usecase

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/flatmix/final-otus-project/internal/migrator/config"
	_ "github.com/lib/pq" //nolint:revive,nolintlint
)

func Up(ctx context.Context, db *sql.DB) (*Outs, error) {
	upStruct := NewDBStruct(db)

	files, err := upStruct.getAllMigrationFile()
	if err != nil {
		return nil, err
	}

	createSQLErr := upStruct.createMigrationsTable(ctx)
	if createSQLErr != nil {
		return nil, createSQLErr
	}

	actualVersion := upStruct.getActualVersion(ctx)

	outUps := make(Outs, 0)
	for _, file := range files {
		outUp, err := upStruct.upMigration(ctx, file, actualVersion)
		if err != nil {
			return nil, err
		}
		if outUp != nil {
			outUps = append(outUps, outUp)
		}
	}

	return &outUps, nil
}

func (ds *DB) upMigration(ctx context.Context, file FileStruct, actualVersion int) (*Out, error) {
	migrationDB, err := ds.getMigrationRow(ctx, file)
	if err != nil {
		return nil, err
	}

	outUp := Out{
		Name:   file.file.Name(),
		Status: "",
	}
	if migrationDB.IsZero() {
		migrateSQLString, err := getUpPart(file)
		if err != nil {
			return nil, err
		}

		err = ds.migrate(ctx, migrateSQLString)
		if err != nil {
			return nil, err
		}
		err = ds.createMigration(ctx, file, actualVersion)
		if err != nil {
			return nil, err
		}
		outUp.Status = "migrate ok"
		return &outUp, nil
	} else if file.hash != migrationDB.Hash {
		outUp.Status = "the migration file has been changed, attention!!!"
		return &outUp, nil
	}

	return nil, nil
}

func getUpPart(fileStruct FileStruct) (string, error) {
	contentFile, err := os.ReadFile(fmt.Sprintf("%s/%s", config.FolderName, fileStruct.file.Name()))
	if err != nil {
		return "", err
	}
	content := string(contentFile)

	re := regexp.MustCompile(regexpUpTemplate)

	res := re.FindAllStringSubmatch(content, -1)

	if len(res) == 0 {
		return "", errors.New("not found template string")
	}

	return strings.TrimSpace(res[0][1]), nil
}
