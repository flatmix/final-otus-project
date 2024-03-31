package storage

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/flatmix/final-otus-project/internal/migrator/config"
)

func NewDB(dbConf config.Postgres) (*sql.DB, error) {
	var psqlInfo string
	if dbConf.DSN != "" {
		psqlInfo = dbConf.DSN
	}

	if dbConf.Username != "" && dbConf.Password != "" {
		sslMode := "disable"
		if dbConf.SslMode {
			sslMode = "enable"
		}
		psqlInfo = fmt.Sprintf("postgres://%s:%s@%s:%v/%s?sslmode=%s",
			dbConf.Username, dbConf.Password, dbConf.Host, dbConf.Port, dbConf.Database, sslMode)
	}
	if psqlInfo == "" {
		if os.Getenv("MIGRATOR_DB_DSN") != "" {
			psqlInfo = os.Getenv("MIGRATOR_DB_DSN")
		}
	}
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
