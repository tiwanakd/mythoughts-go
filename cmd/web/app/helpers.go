package app

import (
	"bytes"
	"fmt"
	"net/http"
	"time"

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
		CurrentYear: time.Now().Year(),
	}
}
