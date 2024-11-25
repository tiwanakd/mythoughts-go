package main

import (
	"github.com/tiwanakd/mythoughts-go/cmd/web/app"
	"github.com/tiwanakd/mythoughts-go/cmd/web/server"
)

func main() {
	a, db := app.New()
	defer db.Close()
	srv := server.New(a.Logger, a.Routes())
	srv.Start()

	// pass, err := bcrypt.GenerateFromPassword([]byte("pa$$word"), 12)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println(string(pass))
}
