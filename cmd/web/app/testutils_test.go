package app

import (
	"html"
	"io"
	"log/slog"
	"regexp"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/tiwanakd/mythoughts-go/cmd/web/templates"
	"github.com/tiwanakd/mythoughts-go/internal/mocks"
)

func newTestApplication(t *testing.T) *Application {
	templateCache, err := templates.NewTemplateCache()
	if err != nil {
		t.Fatal(err)
	}

	sessionManager := scs.New()
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true

	return &Application{
		Logger:         slog.New(slog.NewTextHandler(io.Discard, nil)),
		TemplateCache:  templateCache,
		sessionManager: sessionManager,
		thoughts:       &mocks.ThoughtModel{},
		users:          &mocks.UserModel{},
	}
}

var csrfTokenRX = regexp.MustCompile(`<input type="hidden" name="csrf_token" value="(.+)">`)

func extractCSRFToken(t *testing.T, body string) string {
	matches := csrfTokenRX.FindStringSubmatch(body)
	if len(matches) < 2 {
		t.Fatal("no csrf token found in body")
	}

	return html.UnescapeString(matches[1])
}
