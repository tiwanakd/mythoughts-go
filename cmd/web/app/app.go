package app

import (
	"flag"
	"html/template"
	"log/slog"
	"os"
	"time"

	"github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/v2"
	"github.com/tiwanakd/mythoughts-go/cmd/web/database"
	"github.com/tiwanakd/mythoughts-go/cmd/web/templates"
	"github.com/tiwanakd/mythoughts-go/internal/models"
)

type Application struct {
	Logger         *slog.Logger
	thoughts       models.ThoughtModel
	users          models.UserModel
	TemplateCache  map[string]*template.Template
	sessionManager *scs.SessionManager
}

func New() (*Application, *database.Database) {

	dsn := flag.String("dsn", os.Getenv("DATA_SOURCE_NAME"), "PostgreSQL Data Source Name")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	db, err := database.Open(*dsn)
	if err != nil {
		logger.Error("database:" + err.Error())
		os.Exit(1)
	}

	templateCache, err := templates.NewTemplateCache()
	if err != nil {
		logger.Error("templateCache:" + err.Error())
		os.Exit(1)
	}

	sessionManger := scs.New()
	sessionManger.Store = postgresstore.New(db.DB)
	sessionManger.Lifetime = 12 * time.Hour
	sessionManger.Cookie.Secure = true

	app := &Application{
		Logger:         logger,
		thoughts:       models.ThoughtModel{DB: db.DB},
		users:          models.UserModel{DB: db.DB},
		TemplateCache:  templateCache,
		sessionManager: sessionManger,
	}

	return app, &db
}
