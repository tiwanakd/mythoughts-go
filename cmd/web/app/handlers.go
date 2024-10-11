package app

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/tiwanakd/mythoughts-go/internal/validator"
)

func (app *Application) home(w http.ResponseWriter, r *http.Request) {
	thoughts, err := app.thoughts.ListAll()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := app.newTemplateData(r)
	data.Form = newThoughtForm{}
	data.Thoughts = thoughts
	app.render(w, r, http.StatusOK, "home.html", data)
}

type newThoughtForm struct {
	Content string
	validator.Validator
}

func (app *Application) newThoughtPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	content := r.Form.Get("newThought")

	form := newThoughtForm{
		Content: content,
	}

	form.CheckField(validator.NotBlank(form.Content), "content", "This field Cannot be blank")

	if !form.IsValid() {
		w.WriteHeader(http.StatusUnprocessableEntity)

		data := app.newTemplateData(r)
		data.Form = form

		tmpl := app.TemplateCache["home.html"]
		err = tmpl.ExecuteTemplate(w, "content-error-block", data)
		if err != nil {
			app.serverError(w, r, err)
		}

		return
	}

	thought, err := app.thoughts.Insert(content)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	//get the home page template from cache
	tmpl := app.TemplateCache["home.html"]

	//execute the thoughts-list with new retuened thougt
	err = tmpl.ExecuteTemplate(w, "thoughts-list", thought)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

}

func (app *Application) addLikePost(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	newAgreeCount, err := app.thoughts.AddLike(id)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, newAgreeCount)
}

func (app *Application) addDislikePost(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	newDisagreeCount, err := app.thoughts.AddDislike(id)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, newDisagreeCount)
}
