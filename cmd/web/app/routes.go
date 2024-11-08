package app

import (
	"net/http"

	"github.com/justinas/alice"
	"github.com/tiwanakd/mythoughts-go/cmd/web/middleware"
)

func (app *Application) Routes() http.Handler {
	router := http.NewServeMux()

	middleware := middleware.New(app.Logger, app.sessionManager, app)

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	router.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	dynamic := alice.New(app.sessionManager.LoadAndSave, middleware.Autheticate, middleware.NoSurf, middleware.SortUserThoughts)

	router.Handle("GET /{$}", dynamic.ThenFunc(app.home))
	router.Handle("GET /sort/{sortby}", dynamic.ThenFunc(app.sort))
	router.Handle("POST /like/{id}", dynamic.ThenFunc(app.addLikePost))
	router.Handle("POST /dislike/{id}", dynamic.ThenFunc(app.addDislikePost))
	router.Handle("GET /user/signup", dynamic.ThenFunc(app.userSignup))
	router.Handle("POST /user/signup", dynamic.ThenFunc(app.userSignupPost))
	router.Handle("GET /user/login", dynamic.ThenFunc(app.userLogin))
	router.Handle("POST /user/login", dynamic.ThenFunc(app.userLoginPost))

	authenticated := dynamic.Append(middleware.RequireAuthentication)
	router.Handle("POST /thought/new", authenticated.ThenFunc(app.newThoughtPost))
	router.Handle("POST /user/logout", authenticated.ThenFunc(app.userLogout))
	router.Handle("GET /user/account", authenticated.ThenFunc(app.userAccountView))
	router.Handle("PUT /user/account/edit/{field}", authenticated.ThenFunc(app.userAccountUpdate))
	router.Handle("GET /user/account/password/update", authenticated.ThenFunc(app.userPasswordChange))
	router.Handle("PUT /user/account/password/update", authenticated.ThenFunc(app.userPasswordChangePost))
	router.Handle("GET /user/thoughts/view", authenticated.ThenFunc(app.userThoughtsView))
	router.Handle("DELETE /user/thought/delete/{id}", authenticated.ThenFunc(app.DeleteThoughtPost))
	router.Handle("DELETE /user/account/delete", authenticated.ThenFunc(app.userAccountDelete))

	standard := alice.New(middleware.RecoverPanic, middleware.CommonHeaders, middleware.LogReqest)
	return standard.Then(router)
}
