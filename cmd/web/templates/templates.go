package templates

import (
	"html/template"
	"path/filepath"
	"time"

	"github.com/tiwanakd/mythoughts-go/internal/models"
)

type TemplateData struct {
	CurrentYear     int
	Form            any
	Flash           string
	Thoughts        []models.Thought
	Thought         models.Thought
	User            models.User
	IsAuthenticated bool
	CSRFToken       string
}

func humanDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

// fucntion that caches/parses the html files in to memory when the program starts
func NewTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob("./ui/html/pages/*.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		//get the name of the html template file
		name := filepath.Base(page)

		//parse the base template
		//adding the FuncMap to use the fucntions in templates
		ts, err := template.New(name).Funcs(functions).ParseFiles("./ui/html/base.html")
		if err != nil {
			return nil, err
		}

		//add any partials to the base template set
		ts, err = ts.ParseGlob("./ui/html/partials/*.html")
		if err != nil {
			return nil, err
		}

		//finally add the html pages
		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}
