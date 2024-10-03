package middleware

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/tiwanakd/mythoughts-go/cmd/web/helpers"
)

type Middleware struct {
	Logger *slog.Logger
}

func New(logger *slog.Logger) *Middleware {
	return &Middleware{
		Logger: logger,
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
				helpers.ServerError(w, r, fmt.Errorf("error: %s", err), m.Logger)
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
