package app

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/flatmix/final-otus-project/internal/migrator/config"
	"github.com/flatmix/final-otus-project/internal/migrator/storage"
	"github.com/flatmix/final-otus-project/internal/migrator/usecase"
)

const FolderName = "migrations"

var ErrorNoCommand = errors.New("set command first arguments: {up, create, down, redo, status}")

//nolint:gocognit //a lot of error checks
func Start(mainCtx context.Context, logger *slog.Logger) {
	ctx, cancel := context.WithCancel(mainCtx)
	defer cancel()

	var name string

	var postgresConfig config.Postgres
	var configuration config.Config

	argsWithProg, err := setConfig(&name, &postgresConfig, &configuration)
	if err != nil {
		logger.Error("Parse command line", "Error", err)
		return
	}

	db, dbErr := storage.NewDB(postgresConfig)
	if dbErr != nil {
		logger.Error("Connect to DB", "Error", fmt.Errorf("NewDB: %w", dbErr))
	}

	switch argsWithProg[1] {
	case "create":
		if name == "" && len(argsWithProg) >= 3 {
			name = argsWithProg[2]
		}
		if name == "" {
			fmt.Println("Set name for create command: {--name=... or second arguments}")
		}
		err = usecase.Create(name)
		if err != nil {
			logger.Error("Create error", "Error", err)
			return
		}
	case "up":
		if dbErr != nil {
			return
		}
		outs, err := usecase.Up(ctx, db)
		if err != nil {
			logger.Error("Up error", "Error", err)
			return
		}
		err = usecase.TerminalUpOut(outs)
		if err != nil {
			logger.Error("TerminalUpOut", "Error", err)
			return
		}

		fmt.Println("up success!")

	case "down":
		if dbErr != nil {
			return
		}
		outs, err := usecase.Down(ctx, db, configuration.All, configuration.Step)
		if err != nil {
			logger.Error("Down error", "Error", err)
			return
		}

		err = usecase.TerminalUpOut(outs)
		if err != nil {
			logger.Error("TerminalUpOut", "Error", err)
			return
		}

		fmt.Println("down success!")

	case "redo":
		if dbErr != nil {
			return
		}
		outs, err := usecase.Redo(ctx, db, configuration.All, configuration.Step)
		if err != nil {
			logger.Error("Redo error", "Error", err)
			return
		}

		err = usecase.TerminalUpOut(outs)
		if err != nil {
			logger.Error("TerminalUpOut", "Error", err)
			return
		}

		fmt.Println("redo success!")

	case "status":
		if dbErr != nil {
			return
		}
		statuses, err := usecase.Status(ctx, db)
		if err != nil {
			logger.Error("Status error", "Error", err)
			return
		}
		err = usecase.TerminalStatusOut(statuses)
		if err != nil {
			logger.Error("TerminalStatusOut", "Error", err)
		}
	case "dbversion":
		if dbErr != nil {
			return
		}
		dbver, err := usecase.DBVersion(ctx, db)
		if err != nil {
			logger.Error("DBVersion error", "Error", err)
			return
		}
		fmt.Printf("DB Version: %d \n", *dbver)
	default:
		fmt.Println("Set command first arguments: {up, create, down, redo, status}")
	}
}

func setConfig(name *string, postgresConfig *config.Postgres, configuration *config.Config) ([]string, error) {
	fs := flag.NewFlagSet("ExampleFunc", flag.ContinueOnError)
	fs.SetOutput(os.Stdout)
	fs.StringVar(&postgresConfig.DSN, "dbDSN", "", "The url for connect to db")
	fs.StringVar(&postgresConfig.Username, "user", "", "User for connect to db")
	fs.StringVar(&postgresConfig.Password, "pass", "", "Password for connect to db")
	fs.StringVar(&postgresConfig.Host, "host", "localhost", "Host for connect to db, default `localhost`")
	fs.StringVar(&postgresConfig.Database, "db", "postgres", "Database for connect to db, default `postgres`")
	fs.IntVar(&postgresConfig.Port, "port", 5432, "Port for connect to db, default `5432`")
	fs.BoolVar(&postgresConfig.SslMode, "sslMode", false, "Enable sslmode for connect to db, default `false`")
	fs.BoolVar(&configuration.All, "all", false, "All migration: {down, redo}, default `false`")
	fs.StringVar(name, "name", "", "The name of migration")
	fs.IntVar(&configuration.Step, "step", 0, "Step down and redo on version, works for: {down, redo}, default `0`")
	argsWithProg := os.Args
	if len(argsWithProg) < 2 {
		return []string{}, ErrorNoCommand
	}
	err := fs.Parse(argsWithProg[2:]) //nolint:errcheck
	if err != nil {
		return []string{}, err
	}
	return argsWithProg, nil
}
