package config

const (
	FolderName         = "migrations"
	MigrationTableName = "migrations"
)

type Postgres struct {
	Username string
	Password string
	Host     string
	Port     int
	Database string
	SslMode  bool
	DSN      string
}

type Config struct {
	Step       int
	All        bool
	FolderName string
}
