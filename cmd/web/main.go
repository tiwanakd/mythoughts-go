package main

import (
	"github.com/tiwanakd/mythoughts-go/cmd/web/app"
	"github.com/tiwanakd/mythoughts-go/cmd/web/server"
)

func main() {
	a := app.New()
	srv := server.New(a.Logger, a.Routes())
	srv.Start()
}
