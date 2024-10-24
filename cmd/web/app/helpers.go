package app

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/justinas/nosurf"
	"github.com/tiwanakd/mythoughts-go/cmd/web/templates"
)

func (app *Application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
	)

	app.Logger.Error(err.Error(), "method", method, "uri", uri)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *Application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// The render function will be used to execute the templates
// provided in the template cache
func (app *Application) render(w http.ResponseWriter, r *http.Request, status int, page string, data templates.TemplateData) {
	ts, ok := app.TemplateCache[page]
	if !ok {
		app.serverError(w, r, fmt.Errorf("no template with name %s", page))
		return
	}

	//create a new buffer to execute the emplate onto
	//this ensures that we do end with half rendered HTML
	buff := new(bytes.Buffer)

	err := ts.ExecuteTemplate(buff, "base", data)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	w.WriteHeader(status)
	buff.WriteTo(w)
}

func (app *Application) newTemplateData(r *http.Request) templates.TemplateData {
	return templates.TemplateData{
		CurrentYear:     time.Now().Year(),
		Flash:           app.sessionManager.PopString(r.Context(), "flash"),
		IsAuthenticated: app.IsAuthenticated(r),
		CSRFToken:       nosurf.Token(r),
	}
}

func (app *Application) IsAuthenticated(r *http.Request) bool {
	isAuthenticated, ok := r.Context().Value(isAuthenticatedContextKey).(bool)
	if !ok {
		return false
	}

	return isAuthenticated
}

// This is helper function will be used in middleware to check if the autheticed user still exists DB on each request
// And ensure the user is not deleted from the DB since they last logged in
// If the user still exists in DB, create the copy of the current context by adding the isAuthenticatedContextKey
// Then create a copy of the request by need the new context and return it
// This Method will be added to the Autenitactor Interface in our middleware so it can be used there
func (app *Application) AuthenticateandAddContextKey(id int, w http.ResponseWriter, r *http.Request) *http.Request {
	exists, err := app.users.Exists(id)
	if err != nil {
		app.serverError(w, r, err)
		return nil
	}

	if exists {
		ctx := context.WithValue(r.Context(), isAuthenticatedContextKey, true)
		r = r.WithContext(ctx)
	}

	return r
}
