package app

import (
	"context"
	"flag"
	"fmt"
	"github.com/flatmix/final-otus-project/internal/migrator/usecase"
	"log/slog"
	"os"
)

var url, name string

func Start(mainCtx context.Context, logger *slog.Logger) error {
	ctx, cancel := context.WithCancel(mainCtx)
	defer cancel()
	defer ctx.Done()

	fs := flag.NewFlagSet("ExampleFunc", flag.ContinueOnError)
	fs.SetOutput(os.Stdout)
	fs.StringVar(&url, "url", "", "The url connect to db")
	fs.StringVar(&name, "name", "", "The name of migration")
	argsWithProg := os.Args
	if len(argsWithProg) < 2 {
		fmt.Println("Set command first arguments: {up, create, down, redo, status}")
		return nil
	}
	err := fs.Parse(argsWithProg[2:]) //nolint:errcheck

	if err != nil {
		logger.Error("Parse command line", "Error", err)
		return err
	}

	switch argsWithProg[1] {
	case "create":
		if name == "" && len(argsWithProg) >= 3 {
			name = argsWithProg[2]
		}
		if name == "" {
			fmt.Println("Set name for create command: {--name=... or second arguments}")
			return nil
		}
		err = usecase.Create(name)
		if err != nil {
			logger.Error("Create error", "Error", err)
		}
	case "up":
		fmt.Println("up")

	case "down":
		fmt.Println("down")

	case "redo":
		fmt.Println("redo")

	case "status":
		fmt.Println("status")

	}

	return nil
}
