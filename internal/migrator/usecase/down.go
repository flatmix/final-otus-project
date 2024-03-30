package usecase

import (
	"context"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/flatmix/final-otus-project/internal/migrator/config"
	"github.com/flatmix/final-otus-project/internal/migrator/storage"
)

var (
	ErrNotFountMigrationFile = errors.New("not found migration file")
	ErrNothingForDownMigrate = errors.New("nothing for down migrate")
)

func Down(ctx context.Context, downStruct DBUsecaseContract, all bool, step int) (*Outs, error) {
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

	if all {
		for _, migration := range migrations {
			out, err := downStruct.DownMigration(ctx, migration, filesMap)
			if err != nil {
				return nil, err
			}
			if out != nil {
				outs = append(outs, out)
			}
		}
	}

	if !all {
		out, err := downStruct.DownMigration(ctx, migrations[0], filesMap)
		if err != nil {
			return nil, err
		}
		if out != nil {
			outs = append(outs, out)
		}
	}

	return &outs, nil
}

func (ds *DB) DownMigration(ctx context.Context, migration storage.MigrationDBStruct, filesMap FilesMap) (*Out, error) {
	file, ok := filesMap[migration.Name]
	out := Out{
		Name:   migration.Name,
		Status: "down start",
	}
	if !ok {
		out.Status = "not found migration file"
		return nil, ErrNotFountMigrationFile
	}
	migrateSQLString, err := getDownPart(file)
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

func getDownPart(fileStruct FileStruct) (string, error) {
	contentFile, err := os.ReadFile(fmt.Sprintf("%s/%s", config.FolderName, fileStruct.File.Name()))
	if err != nil {
		return "", err
	}
	content := string(contentFile)

	re := regexp.MustCompile(regexpDownTemplate)

	res := re.FindAllStringSubmatch(content, -1)

	if len(res) == 0 {
		return "", errors.New("not found template string")
	}

	return strings.TrimSpace(res[0][1]), nil
}
