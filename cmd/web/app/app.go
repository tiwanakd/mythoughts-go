package app

import (
	"flag"
	"log/slog"
	"os"

	"github.com/tiwanakd/mythoughts-go/cmd/web/database"
)

type Application struct {
	Logger *slog.Logger
}

func New() *Application {

	dsn := flag.String("dsn", os.Getenv("DATA_SOURCE_NAME"), "PostgreSQL Data Source Name")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	db, err := database.OpenDB(*dsn)
	if err != nil {
		logger.Error("database error:" + err.Error())
		os.Exit(1)
	}
	defer db.Close()

	return &Application{
		Logger: logger,
	}
}
