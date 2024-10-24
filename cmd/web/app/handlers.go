package app

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/tiwanakd/mythoughts-go/internal/models"
	"github.com/tiwanakd/mythoughts-go/internal/validator"
)

func (app *Application) home(w http.ResponseWriter, r *http.Request) {
	thoughts, err := app.thoughts.List("created")
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := app.newTemplateData(r)
	data.Form = newThoughtForm{}
	data.Thoughts = thoughts
	app.render(w, r, http.StatusOK, "home.html", data)
}

func (app *Application) sort(w http.ResponseWriter, r *http.Request) {
	sortby := r.PathValue("sortby")
	thoughts, err := app.thoughts.List(sortby)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := app.newTemplateData(r)
	data.Form = newThoughtForm{}
	data.Thoughts = thoughts
	app.render(w, r, http.StatusOK, "home.html", data)
}

type newThoughtForm struct {
	Content string
	validator.Validator
}

func (app *Application) newThoughtPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	content := r.Form.Get("newThought")

	form := newThoughtForm{
		Content: content,
	}

	form.CheckField(validator.NotBlank(form.Content), "content", "This field Cannot be blank")
	form.CheckField(validator.MinChars(form.Content, 30), "content", "Your thought should have alteast 50 characters")
	form.CheckField(validator.MaxChars(form.Content, 200), "content", "Your thought cannot have more than 200 Characters")

	if !form.IsValid() {
		w.Header().Set("HX-Reswap", "innerHTML")
		w.WriteHeader(http.StatusUnprocessableEntity)

		data := app.newTemplateData(r)
		data.Form = form

		tmpl := app.TemplateCache["home.html"]
		err = tmpl.ExecuteTemplate(w, "content-error-block", data)
		if err != nil {
			app.serverError(w, r, err)
		}

		return
	}

	userID, ok := app.sessionManager.Get(r.Context(), "authenticatedUserID").(int)
	if !ok {
		app.serverError(w, r, err)
		return
	}

	_, err = app.thoughts.Insert(content, userID)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// //get the home page template from cache
	// tmpl := app.TemplateCache["home.html"]

	// buf := new(bytes.Buffer)

	// //execute the thoughts-list with new retuened thougt
	// err = tmpl.ExecuteTemplate(buf, "thoughts-list", thought)
	// if err != nil {
	// 	app.serverError(w, r, err)
	// 	return
	// }

	//Using HX-Refresh header to trigger page refresh on successful POST
	w.Header().Set("HX-Refresh", "true")
	w.WriteHeader(http.StatusCreated)
	// buf.WriteTo(w)
}

func (app *Application) addLikePost(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	newAgreeCount, err := app.thoughts.AddLike(id)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, newAgreeCount)
}

func (app *Application) addDislikePost(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	newDisagreeCount, err := app.thoughts.AddDislike(id)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, newDisagreeCount)
}

type userSignupForm struct {
	Name     string
	Username string
	Email    string
	Password string
	validator.Validator
}

func (app *Application) userSignup(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userSignupForm{}
	app.render(w, r, http.StatusOK, "signup.html", data)
}

func (app *Application) userSignupPost(w http.ResponseWriter, r *http.Request) {
	var form userSignupForm
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.Name = r.Form.Get("name")
	form.Username = r.Form.Get("username")
	form.Email = r.Form.Get("email")
	form.Password = r.Form.Get("password")

	form.Validator.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be empty")
	form.Validator.CheckField(validator.NotBlank(form.Username), "username", "This field cannot be empty")
	form.Validator.CheckField(validator.MaxChars(form.Username, 15), "username", "Username cannot be more than 15 characters long")
	form.Validator.CheckField(validator.MinChars(form.Username, 4), "username", "Username should have atleast 4 characters")
	form.Validator.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be empty")
	form.Validator.CheckField(validator.Matches(form.Email, validator.EmailRx), "email", "This field must be valid email address.")
	form.Validator.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be empty")
	form.Validator.CheckField(validator.ValidPassword(form.Password), "password", "Password does not meet the length and complexity requirements")

	if !form.IsValid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "signup.html", data)
		return
	}

	err = app.users.Insert(form.Username, form.Email, form.Name, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldError("email", "Email already in use")
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, r, http.StatusUnprocessableEntity, "signup.html", data)
			return
		}

		if errors.Is(err, models.ErrDuplicateUsername) {
			form.AddFieldError("username", "Username already taken")
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, r, http.StatusUnprocessableEntity, "signup.html", data)
			return
		}

		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "Your Signup was successful, please login.")
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

type userLoginForm struct {
	Email    string
	Password string
	validator.Validator
}

func (app *Application) userLogin(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userLoginForm{}
	app.render(w, r, http.StatusOK, "login.html", data)
}

func (app *Application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	var form userLoginForm
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.Email = r.Form.Get("email")
	form.Password = r.Form.Get("password")

	form.Validator.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be empty")
	form.Validator.CheckField(validator.Matches(form.Email, validator.EmailRx), "email", "This field must be valid email address.")
	form.Validator.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be empty")

	if !form.IsValid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "login.html", data)
		return
	}

	id, err := app.users.Authenticate(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentails) {
			form.AddNonFieldError("Invalid Username or Password")
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, r, http.StatusUnprocessableEntity, "login.html", data)
			return
		}
	}

	//this will generate a new session id
	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Put(r.Context(), "authenticatedUserID", id)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *Application) userLogout(w http.ResponseWriter, r *http.Request) {
	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Remove(r.Context(), "authenticatedUserID")
	app.sessionManager.Put(r.Context(), "flash", "You have been logged out successfully.")

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *Application) userAccountView(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Account page")
}
