package app

import (
	"flag"
	"log/slog"
	"os"

	"github.com/tiwanakd/mythoughts-go/cmd/web/database"
	"github.com/tiwanakd/mythoughts-go/internal/models"
)

type Application struct {
	Logger   *slog.Logger
	Thoughts models.ThoughtModel
}

func New() (*Application, *database.Database) {

	dsn := flag.String("dsn", os.Getenv("DATA_SOURCE_NAME"), "PostgreSQL Data Source Name")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	db, err := database.OpenDB(*dsn)
	if err != nil {
		logger.Error("database error:" + err.Error())
		os.Exit(1)
	}

	app := &Application{
		Logger:   logger,
		Thoughts: models.ThoughtModel{DB: db.DB},
	}

	return app, &db
}
