package usecase

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/flatmix/final-otus-project/internal/migrator/config"
	"github.com/flatmix/final-otus-project/internal/migrator/storage"
)

type DBUsecaseContract interface {
	getHash(fileName string) (string, error)
	getAllMigrationFile() ([]FileStruct, error)
	getAllMigrationFileMap() (FilesMap, error)
	getMigrationRow(ctx context.Context, file FileStruct) (*storage.MigrationDBStruct, error)
	getAllMigrationsOrderByVersionDesc(ctx context.Context,
		step int,
	) (storage.MigrationsDBStruct, error)
	getActualVersion(ctx context.Context) int
	migrate(ctx context.Context, migrateSQLString string) error
	createMigration(ctx context.Context, file FileStruct, version int) error
	deleteMigration(ctx context.Context, file FileStruct) error
	existTable(ctx context.Context, schema string, table string) bool
	createMigrationsTable(ctx context.Context) error
	redoMigration(ctx context.Context, migrations storage.MigrationsDBStruct, filesMap FilesMap) (*Outs, error)
	downMigration(ctx context.Context, migration storage.MigrationDBStruct, filesMap FilesMap) (*Out, error)
	upMigration(ctx context.Context, file FileStruct, actualVersion int) (*Out, error)
}

type DB struct {
	db *sql.DB
}

func NewDBStruct(db *sql.DB) DBUsecaseContract {
	return &DB{db: db}
}

func (ds *DB) getHash(fileName string) (string, error) {
	hash := sha256.New()
	file, err := os.Open(fmt.Sprintf("%s/%s", config.FolderName, fileName))
	if err != nil {
		return "", fmt.Errorf("open file: %w", err)
	}
	if _, err := io.Copy(hash, file); err != nil {
		return "", fmt.Errorf("get hash of file: %w", err)
	}
	stringHash := base64.URLEncoding.EncodeToString(hash.Sum(nil))

	return stringHash, nil
}

func (ds *DB) getAllMigrationFile() ([]FileStruct, error) {
	var migrations []FileStruct //nolint:prealloc

	files, _ := os.ReadDir(config.FolderName)
	for _, file := range files {
		info, err := file.Info()
		if err != nil {
			return nil, fmt.Errorf("get file info: %w", err)
		}
		hash, err := ds.getHash(file.Name())
		if err != nil {
			return nil, fmt.Errorf("get file hash: %w", err)
		}

		migrations = append(migrations, FileStruct{
			file: info,
			hash: hash,
		})
	}

	return migrations, nil
}

func (ds *DB) getAllMigrationFileMap() (FilesMap, error) {
	files, _ := os.ReadDir(config.FolderName)
	migrations := make(FilesMap, len(files))
	for _, file := range files {
		info, err := file.Info()
		if err != nil {
			return nil, fmt.Errorf("get file info: %w", err)
		}
		hash, err := ds.getHash(file.Name())
		if err != nil {
			return nil, fmt.Errorf("get file hash: %w", err)
		}

		migrations[file.Name()] = FileStruct{
			file: info,
			hash: hash,
		}
	}

	return migrations, nil
}

func (ds *DB) getMigrationRow(ctx context.Context, file FileStruct) (*storage.MigrationDBStruct, error) {
	selectMigrationTableQuery := fmt.Sprintf(selectMigrationTable, config.MigrationTableName)
	var migrationDB storage.MigrationDBStruct
	rows, err := ds.db.QueryContext(ctx, selectMigrationTableQuery, file.file.Name())
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		if err := rows.Scan(&migrationDB.ID, &migrationDB.Name,
			&migrationDB.Hash, &migrationDB.Version, &migrationDB.CreatedAt, &migrationDB.UpdatedAt); err != nil {
			log.Fatal(err)
		}
	}

	return &migrationDB, nil
}

func (ds *DB) getAllMigrationsOrderByVersionDesc(ctx context.Context,
	step int,
) (storage.MigrationsDBStruct, error) {
	selectAllMigrationsTableQuery := fmt.Sprintf(selectAllMigrationsTable,
		config.MigrationTableName, config.MigrationTableName)

	var migrationsDB []storage.MigrationDBStruct

	actual := ds.getActualVersion(ctx)

	if step == 0 {
		actual = -1
	}

	rows, err := ds.db.QueryContext(ctx, selectAllMigrationsTableQuery, actual-step)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		migrationDB := storage.MigrationDBStruct{}
		if err := rows.Scan(&migrationDB.ID, &migrationDB.Name,
			&migrationDB.Hash, &migrationDB.Version,
			&migrationDB.CreatedAt, &migrationDB.UpdatedAt); err != nil {
			log.Fatal(err)
		}
		migrationsDB = append(migrationsDB, migrationDB)
	}

	return migrationsDB, nil
}

func (ds *DB) getActualVersion(ctx context.Context) int {
	getActualVersion := fmt.Sprintf(getMaxVersionFromMigrationsTable, config.MigrationTableName)
	var versionDB int
	rows, err := ds.db.QueryContext(ctx, getActualVersion)
	if err != nil {
		return 0
	}

	for rows.Next() {
		if err := rows.Scan(&versionDB); err != nil {
			return 0
		}
	}
	return versionDB + 1
}

func (ds *DB) migrate(ctx context.Context, migrateSQLString string) error {
	_, err := ds.db.ExecContext(ctx, migrateSQLString)
	if err != nil {
		return err
	}

	return nil
}

func (ds *DB) createMigration(ctx context.Context, file FileStruct, version int) error {
	insertMigrationsTableQuery := fmt.Sprintf(insertMigrationsTable, config.MigrationTableName)
	_, err := ds.db.ExecContext(ctx, insertMigrationsTableQuery, file.file.Name(), file.hash, version, time.Now())
	if err != nil {
		return err
	}

	return nil
}

func (ds *DB) deleteMigration(ctx context.Context, file FileStruct) error {
	deleteRowFromMigrationsTableQuery := fmt.Sprintf(deleteRowFromMigrationsTable, config.MigrationTableName)
	_, err := ds.db.ExecContext(ctx, deleteRowFromMigrationsTableQuery, file.file.Name())
	if err != nil {
		return err
	}

	return nil
}

func (ds *DB) existTable(ctx context.Context, schema string, table string) bool {
	selectExistsTableQuery := fmt.Sprintf(selectExistsTable, schema, table)

	rows, err := ds.db.QueryContext(ctx, selectExistsTableQuery)
	if err != nil {
		return false
	}
	var exist bool

	for rows.Next() {
		if err := rows.Scan(&exist); err != nil {
			log.Fatal(err)
		}
	}

	return exist
}

func (ds *DB) createMigrationsTable(ctx context.Context) error {
	exists := ds.existTable(ctx, "public", config.MigrationTableName)

	if !exists {
		createMigrationsTableQuery := fmt.Sprintf(createMigrationsTable, config.MigrationTableName)

		_, err := ds.db.ExecContext(ctx, createMigrationsTableQuery)
		if err != nil {
			return err
		}
		fmt.Printf("create migration table - '%s' \n", config.MigrationTableName)
	}

	return nil
}
