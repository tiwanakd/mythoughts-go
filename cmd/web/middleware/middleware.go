package middleware

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/justinas/nosurf"
)

type Middleware struct {
	Logger         *slog.Logger
	SessionManager *scs.SessionManager
	Autheticator
}

// create an inteface that has the Autheticate Methods from the Application Helpers
// Application struct will implement this interface which allow to use this funciton in Middleware
type Autheticator interface {
	IsAuthenticated(r *http.Request) bool
	AuthenticateandAddContextKey(id int, w http.ResponseWriter, r *http.Request) *http.Request
}

func New(logger *slog.Logger, sessionManager *scs.SessionManager, auth Autheticator) *Middleware {
	return &Middleware{
		Logger:         logger,
		SessionManager: sessionManager,
		Autheticator:   auth,
	}
}

// To write the HTTP status code of the Response to the logger
// creating a wrappedwriter that extends the http.ResponseWriter
type wrappedWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *wrappedWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.statusCode = statusCode
}

func (m *Middleware) LogReqest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			ip     = r.RemoteAddr
			proto  = r.Proto
			method = r.Method
			uri    = r.URL.RequestURI()
		)

		wrapped := &wrappedWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(wrapped, r)
		m.Logger.Info("recevied request", "ip", ip, "proto", proto, "method", method, "uri", uri, "response code", wrapped.statusCode)
	})
}

func (m *Middleware) RecoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
				m.Logger.Error("panic: " + fmt.Sprintf("error: %s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (m *Middleware) CommonHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy",
			"default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com",
		)
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")

		w.Header().Set("Server", "Go")

		next.ServeHTTP(w, r)
	})
}

func (m *Middleware) RequireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !m.IsAuthenticated(r) {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}

		w.Header().Add("Cache-Control", "no-store")
		next.ServeHTTP(w, r)
	})
}

func (m *Middleware) NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   true,
	})

	return csrfHandler
}

func (m *Middleware) Autheticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := m.SessionManager.GetInt(r.Context(), "authenticatedUserID")
		if id == 0 {
			next.ServeHTTP(w, r)
			return
		}

		rWithContextKey := m.AuthenticateandAddContextKey(id, w, r)
		next.ServeHTTP(w, rWithContextKey)
	})
}
