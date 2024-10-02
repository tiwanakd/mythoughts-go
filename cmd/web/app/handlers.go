package app

import (
	"net/http"
)

func (app *Application) home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Home"))
}

func (app *Application) newThought(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display page to create new thought"))
}

func (app *Application) newThoughtPost(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Post new thought"))
}
