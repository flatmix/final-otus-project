package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/flatmix/final-otus-project/internal/migrator/app"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	app.Start(context.Background(), logger)
	// err := app.Start(context.Background(), logger)
	// if err != nil {
	//
	//}
}
