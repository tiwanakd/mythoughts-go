package app

import (
	"flag"
	"html/template"
	"log/slog"
	"os"

	"github.com/tiwanakd/mythoughts-go/cmd/web/database"
	"github.com/tiwanakd/mythoughts-go/cmd/web/templates"
	"github.com/tiwanakd/mythoughts-go/internal/models"
)

type Application struct {
	Logger        *slog.Logger
	Thoughts      models.ThoughtModel
	TemplateCache map[string]*template.Template
}

func New() (*Application, *database.Database) {

	dsn := flag.String("dsn", os.Getenv("DATA_SOURCE_NAME"), "PostgreSQL Data Source Name")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	db, err := database.OpenDB(*dsn)
	if err != nil {
		logger.Error("database:" + err.Error())
		os.Exit(1)
	}

	templateCache, err := templates.NewTemplateCache()
	if err != nil {
		logger.Error("templateCache:" + err.Error())
	}

	app := &Application{
		Logger:        logger,
		Thoughts:      models.ThoughtModel{DB: db.DB},
		TemplateCache: templateCache,
	}

	return app, &db
}
