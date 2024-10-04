package app

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func (app *Application) home(w http.ResponseWriter, r *http.Request) {
	thoughts, err := app.Thoughts.ListAll()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := app.newTemplateData(r)
	data.Thoughts = thoughts
	app.render(w, r, http.StatusOK, "home.html", data)
}

func (app *Application) newThought(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display page to create new thought"))
}

func (app *Application) newThoughtPost(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Post new thought"))
}

func (app *Application) addLikePost(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	newAgreeCount, err := app.Thoughts.AddLike(id)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	//sending a json reponse with the new Agree count
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	encoder := json.NewEncoder(w)
	err = encoder.Encode(map[string]int{
		"newAgreeCount": newAgreeCount,
	})
}

func (app *Application) addDislikePost(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	newDisagreeCount, err := app.Thoughts.AddDislike(id)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	encoder := json.NewEncoder(w)
	encoder.Encode(map[string]int{
		"newDisagreeCount": newDisagreeCount,
	})
}
