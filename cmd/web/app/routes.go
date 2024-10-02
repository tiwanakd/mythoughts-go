package app

import (
	"net/http"

	"github.com/justinas/alice"
	"github.com/tiwanakd/mythoughts-go/cmd/web/middleware"
)

func (app *Application) Routes() http.Handler {
	router := http.NewServeMux()

	router.HandleFunc("GET /{$}", app.home)
	router.HandleFunc("GET /thought/new", app.newThought)
	router.HandleFunc("POST /thought/new", app.newThoughtPost)

	middleware := middleware.New(app.Logger)

	standard := alice.New(middleware.RecoverPanic, middleware.LogReqest, middleware.CommonHeaders)
	return standard.Then(router)
}
