package server

import (
	"crypto/tls"
	"flag"
	"log/slog"
	"net/http"
	"os"
	"time"
)

type Server struct {
	*http.Server
	log  *slog.Logger
	addr string
}

func New(logger *slog.Logger, routes http.Handler) Server {

	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()

	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}
	srv := http.Server{
		Addr:         *addr,
		Handler:      routes,
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	return Server{&srv, logger, *addr}
}

func (s *Server) Start() {

	s.log.Info("starting server", "addr", s.addr)
	err := s.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	s.log.Error(err.Error())
	os.Exit(1)
}
