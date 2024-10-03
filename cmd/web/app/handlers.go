package app

import (
	"html/template"
	"net/http"

	"github.com/tiwanakd/mythoughts-go/cmd/web/helpers"
)

func (app *Application) home(w http.ResponseWriter, r *http.Request) {

	thoughts, err := app.Thoughts.ListAll()
	if err != nil {
		helpers.ServerError(w, r, err, app.Logger)
		return
	}
	files := []string{
		"./ui/html/base.html",
		"./ui/html/pages/home.html",
	}
	ts, err := template.ParseFiles(files...)
	if err != nil {
		helpers.ServerError(w, r, err, app.Logger)
		return
	}

	err = ts.ExecuteTemplate(w, "base", thoughts)
	if err != nil {
		helpers.ServerError(w, r, err, app.Logger)
		return
	}
}

func (app *Application) newThought(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display page to create new thought"))
}

func (app *Application) newThoughtPost(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Post new thought"))
}
