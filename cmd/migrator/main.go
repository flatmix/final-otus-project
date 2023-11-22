package main

import (
	"context"
	"github.com/flatmix/final-otus-project/internal/migrator/app"
	"log/slog"
	"os"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	err := app.Start(context.Background(), logger)
	if err != nil {
		logger.Error("Error start app: ", err)
	}
}
