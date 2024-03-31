package usecase

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/flatmix/final-otus-project/internal/migrator/config"
	"github.com/flatmix/final-otus-project/internal/migrator/storage"
)

type DBContract interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
}

type FileUsecaseContract interface {
	GetHash(fileName string) (string, error)
	GetAllMigrationFile() ([]FileStruct, error)
	GetAllMigrationFileMap() (FilesMap, error)
	GetUpPart(fileStruct FileStruct) (string, error)
	GetDownPart(fileStruct FileStruct) (string, error)
}

type DBUsecaseContract interface {
	GetMigrationRow(ctx context.Context, file FileStruct) (*storage.MigrationDBStruct, error)
	GetAllMigrationsOrderByVersionDesc(ctx context.Context,
		step int,
	) (storage.MigrationsDBStruct, error)
	GetActualVersion(ctx context.Context) int
	Migrate(ctx context.Context, migrateSQLString string) error
	CreateMigration(ctx context.Context, file FileStruct, version int) error
	DeleteMigration(ctx context.Context, file FileStruct) error
	ExistTable(ctx context.Context, schema string, table string) bool
	CreateMigrationsTable(ctx context.Context) error
}

type UCContract interface {
	FileUsecaseContract
	DBUsecaseContract
}

type Usecase struct {
	db  DBContract
	cfg *config.Config
}

func NewUsecase(db DBContract, cfg *config.Config) UCContract {
	return &Usecase{db: db, cfg: cfg}
}

func (uc *Usecase) GetHash(fileName string) (string, error) {
	hash := sha256.New()
	file, err := os.Open(fmt.Sprintf("%s/%s", uc.cfg.FolderName, fileName))
	if err != nil {
		return "", fmt.Errorf("open file: %w", err)
	}
	if _, err := io.Copy(hash, file); err != nil {
		return "", fmt.Errorf("get hash of file: %w", err)
	}
	stringHash := base64.URLEncoding.EncodeToString(hash.Sum(nil))

	return stringHash, nil
}

func (uc *Usecase) GetAllMigrationFile() ([]FileStruct, error) {
	var migrations []FileStruct //nolint:prealloc

	files, _ := os.ReadDir(uc.cfg.FolderName)
	for _, file := range files {
		info, err := file.Info()
		if err != nil {
			return nil, fmt.Errorf("get file info: %w", err)
		}
		hash, err := uc.GetHash(file.Name())
		if err != nil {
			return nil, fmt.Errorf("get file hash: %w", err)
		}

		migrations = append(migrations, FileStruct{
			File: info,
			Hash: hash,
		})
	}

	return migrations, nil
}

func (uc *Usecase) GetAllMigrationFileMap() (FilesMap, error) {
	files, _ := os.ReadDir(uc.cfg.FolderName)
	migrations := make(FilesMap, len(files))
	for _, file := range files {
		info, err := file.Info()
		if err != nil {
			return nil, fmt.Errorf("get file info: %w", err)
		}
		hash, err := uc.GetHash(file.Name())
		if err != nil {
			return nil, fmt.Errorf("get file hash: %w", err)
		}

		migrations[file.Name()] = FileStruct{
			File: info,
			Hash: hash,
		}
	}

	return migrations, nil
}

func (uc *Usecase) GetMigrationRow(ctx context.Context, file FileStruct) (*storage.MigrationDBStruct, error) {
	selectMigrationTableQuery := fmt.Sprintf(selectMigrationTable, config.MigrationTableName)
	var migrationDB storage.MigrationDBStruct
	rows, err := uc.db.QueryContext(ctx, selectMigrationTableQuery, file.File.Name())
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&migrationDB.ID, &migrationDB.Name,
			&migrationDB.Hash, &migrationDB.Version, &migrationDB.CreatedAt, &migrationDB.UpdatedAt); err != nil {
			return nil, fmt.Errorf("%w", err)
		}
	}

	return &migrationDB, nil
}

func (uc *Usecase) GetAllMigrationsOrderByVersionDesc(ctx context.Context,
	step int,
) (storage.MigrationsDBStruct, error) {
	selectAllMigrationsTableQuery := fmt.Sprintf(selectAllMigrationsTable,
		config.MigrationTableName, config.MigrationTableName)

	var migrationsDB []storage.MigrationDBStruct

	actual := uc.GetActualVersion(ctx)

	if step == 0 {
		actual = -1
	}

	rows, err := uc.db.QueryContext(ctx, selectAllMigrationsTableQuery, actual-step)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		migrationDB := storage.MigrationDBStruct{}
		if err := rows.Scan(&migrationDB.ID, &migrationDB.Name,
			&migrationDB.Hash, &migrationDB.Version,
			&migrationDB.CreatedAt, &migrationDB.UpdatedAt); err != nil {
			return nil, fmt.Errorf("%w", err)
		}
		migrationsDB = append(migrationsDB, migrationDB)
	}
	return migrationsDB, nil
}

func (uc *Usecase) GetActualVersion(ctx context.Context) int {
	getActualVersion := fmt.Sprintf(getMaxVersionFromMigrationsTable, config.MigrationTableName)
	var versionDB int
	rows, err := uc.db.QueryContext(ctx, getActualVersion)
	if err != nil {
		return 0
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&versionDB); err != nil {
			return 0
		}
	}
	return versionDB + 1
}

func (uc *Usecase) Migrate(ctx context.Context, migrateSQLString string) error {
	_, err := uc.db.ExecContext(ctx, migrateSQLString)
	if err != nil {
		return err
	}

	return nil
}

func (uc *Usecase) CreateMigration(ctx context.Context, file FileStruct, version int) error {
	insertMigrationsTableQuery := fmt.Sprintf(insertMigrationsTable, config.MigrationTableName)
	_, err := uc.db.ExecContext(ctx, insertMigrationsTableQuery, file.File.Name(), file.Hash, version, time.Now())
	if err != nil {
		return err
	}

	return nil
}

func (uc *Usecase) DeleteMigration(ctx context.Context, file FileStruct) error {
	deleteRowFromMigrationsTableQuery := fmt.Sprintf(deleteRowFromMigrationsTable, config.MigrationTableName)
	_, err := uc.db.ExecContext(ctx, deleteRowFromMigrationsTableQuery, file.File.Name())
	if err != nil {
		return err
	}

	return nil
}

func (uc *Usecase) ExistTable(ctx context.Context, schema string, table string) bool {
	selectExistsTableQuery := fmt.Sprintf(selectExistsTable, schema, table)

	rows, err := uc.db.QueryContext(ctx, selectExistsTableQuery)
	if err != nil {
		return false
	}
	defer rows.Close()
	var exist bool

	for rows.Next() {
		if err := rows.Scan(&exist); err != nil {
			return false
		}
	}

	return exist
}

func (uc *Usecase) CreateMigrationsTable(ctx context.Context) error {
	exists := uc.ExistTable(ctx, "public", config.MigrationTableName)

	if !exists {
		createMigrationsTableQuery := fmt.Sprintf(createMigrationsTable, config.MigrationTableName)

		_, err := uc.db.ExecContext(ctx, createMigrationsTableQuery)
		if err != nil {
			return err
		}
		fmt.Printf("create migration table - '%s' \n", config.MigrationTableName)
	}

	return nil
}

func (uc *Usecase) GetUpPart(fileStruct FileStruct) (string, error) {
	contentFile, err := os.ReadFile(fmt.Sprintf("%s/%s", uc.cfg.FolderName, fileStruct.File.Name()))
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

func (uc *Usecase) GetDownPart(fileStruct FileStruct) (string, error) {
	contentFile, err := os.ReadFile(fmt.Sprintf("%s/%s", uc.cfg.FolderName, fileStruct.File.Name()))
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
