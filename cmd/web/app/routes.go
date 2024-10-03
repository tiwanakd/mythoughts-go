package app

import (
	"net/http"

	"github.com/justinas/alice"
	"github.com/tiwanakd/mythoughts-go/cmd/web/middleware"
)

func (app *Application) Routes() http.Handler {
	router := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	router.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	router.HandleFunc("GET /{$}", app.home)
	router.HandleFunc("GET /thought/new", app.newThought)
	router.HandleFunc("POST /thought/new", app.newThoughtPost)

	middleware := middleware.New(app.Logger)

	standard := alice.New(middleware.RecoverPanic, middleware.LogReqest, middleware.CommonHeaders)
	return standard.Then(router)
}
