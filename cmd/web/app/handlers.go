package app

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/tiwanakd/mythoughts-go/internal/models"
	"github.com/tiwanakd/mythoughts-go/internal/validator"
)

func (app *Application) homeView(w http.ResponseWriter, r *http.Request) {
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
	sortByUserThoughts, ok := app.sessionManager.Get(r.Context(), "userThoughtsSort").(bool)

	sortby := r.PathValue("sortby")

	if sortByUserThoughts && ok {
		userID, ok := app.sessionManager.Get(r.Context(), "authenticatedUserID").(int)
		if !ok {
			app.serverError(w, r, fmt.Errorf("authenticatedUserID: type error"))
			return
		}

		thoughts, err := app.thoughts.UserThoughts(userID, sortby)
		if err != nil {
			app.serverError(w, r, err)
			return
		}
		data := app.newTemplateData(r)
		data.Thoughts = thoughts

		app.render(w, r, http.StatusOK, "userthoughts.html", data)
		return
	} else {
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

	//store the authenticated userID in request context using session manager
	app.sessionManager.Put(r.Context(), "authenticatedUserID", id)

	redirectURI := app.sessionManager.PopString(r.Context(), "redirectURI")
	if redirectURI != "" {
		http.Redirect(w, r, redirectURI, http.StatusSeeOther)
		return
	}

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

func (app *Application) userThoughtsView(w http.ResponseWriter, r *http.Request) {
	userID, ok := app.sessionManager.Get(r.Context(), "authenticatedUserID").(int)
	if !ok {
		app.serverError(w, r, fmt.Errorf("authenticatedUserID: type error"))
		return
	}

	thoughts, err := app.thoughts.UserThoughts(userID, "created")
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	user, err := app.users.Get(userID)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := app.newTemplateData(r)
	data.Thoughts = thoughts
	data.User = user

	//use session manager to add a key to request context and set it to true
	//this will be used in sort handler with assistance of the middleware
	//to check from where the /sort/{sortby} uri is invoked and will
	//sort based on home page or My Thougts page for the logged in user
	app.sessionManager.Put(r.Context(), "userThoughtsSort", true)
	app.render(w, r, http.StatusOK, "userthoughts.html", data)
}

func (app *Application) DeleteThoughtPost(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	err = app.thoughts.DeleteThought(id)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

type userEditForm struct {
	Value string
	validator.Validator
}

func (app *Application) userAccountView(w http.ResponseWriter, r *http.Request) {
	id, ok := app.sessionManager.Get(r.Context(), "authenticatedUserID").(int)
	if !ok {
		app.serverError(w, r, fmt.Errorf("authenticatedUserID: type error"))
		return
	}

	user, err := app.users.Get(id)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := app.newTemplateData(r)
	data.User = user
	data.Form = userEditForm{}

	app.render(w, r, http.StatusOK, "useraccount.html", data)

}

func (app *Application) userAccountUpdate(w http.ResponseWriter, r *http.Request) {
	var form userEditForm
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	columnName := r.PathValue("field")
	form.Value = r.Form.Get(columnName)

	form.Validator.CheckField(validator.NotBlank(form.Value), columnName, "This field cannot be blank")

	if columnName == "username" {
		form.Validator.CheckField(validator.MaxChars(form.Value, 15), columnName, "Username cannot be more than 15 characters long")
		form.Validator.CheckField(validator.MinChars(form.Value, 4), columnName, "Username should have atleast 4 characters")
	}

	if columnName == "email" {
		form.Validator.CheckField(validator.Matches(form.Value, validator.EmailRx), columnName, "This field must be valid email address.")
	}

	if !form.IsValid() {
		w.Header().Set("HX-Reswap", "innerHTML")
		w.WriteHeader(http.StatusUnprocessableEntity)

		data := app.newTemplateData(r)
		data.Form = form

		tmpl := app.TemplateCache["useraccount.html"]

		if columnName == "name" {
			err = tmpl.ExecuteTemplate(w, "content-error-block-name", data)
		} else if columnName == "username" {
			err = tmpl.ExecuteTemplate(w, "content-error-block-username", data)
		} else if columnName == "email" {
			err = tmpl.ExecuteTemplate(w, "content-error-block-email", data)
		}

		if err != nil {
			app.serverError(w, r, err)
		}

		return
	}

	id, ok := app.sessionManager.Get(r.Context(), "authenticatedUserID").(int)
	if !ok {
		app.serverError(w, r, fmt.Errorf("authenticatedUserID: type error"))
		return
	}

	err = app.users.Update(id, columnName, form.Value)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	w.Header().Set("HX-Refresh", "true")
	w.WriteHeader(http.StatusCreated)
}

type changePasswordForm struct {
	currentPassword string
	newPassword     string
	confirmPassword string
	validator.Validator
}

func (app *Application) userPasswordChange(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = changePasswordForm{}
	app.render(w, r, http.StatusOK, "userpassword.html", data)
}

func (app *Application) userPasswordChangePost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	var form changePasswordForm
	form.currentPassword = r.Form.Get("currentPassword")
	form.newPassword = r.Form.Get("newPassword")
	form.confirmPassword = r.Form.Get("confirmPassword")

	form.Validator.CheckField(validator.NotBlank(form.currentPassword), "password", "Password fields cannot be blank")
	form.Validator.CheckField(validator.NotBlank(form.newPassword), "password", " Password fields cannot be blank")
	form.Validator.CheckField(validator.NotBlank(form.confirmPassword), "password", " Password fields cannot be blank")
	form.Validator.CheckField(validator.ValidPassword(form.newPassword), "password", "New Password does not meet the length and complexity requirements")

	if form.currentPassword == form.newPassword {
		form.AddFieldError("password", "New Password cannot be same as Current Password")
	}

	if form.newPassword != form.confirmPassword {
		form.AddFieldError("password", "New Passwords not match")
	}

	if !form.IsValid() {
		w.Header().Set("HX-Reswap", "innerHTML")
		w.WriteHeader(http.StatusUnprocessableEntity)

		data := app.newTemplateData(r)
		data.Form = form

		tmpl := app.TemplateCache["userpassword.html"]
		err = tmpl.ExecuteTemplate(w, "password-error-block", data)
		if err != nil {
			app.serverError(w, r, err)
		}

		return
	}

	id, ok := app.sessionManager.Get(r.Context(), "authenticatedUserID").(int)
	if !ok {
		app.serverError(w, r, fmt.Errorf("authenticatedUserID: type error"))
		return
	}

	err = app.users.ChangePassword(id, form.currentPassword, form.newPassword)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentails) {
			form.AddFieldError("password", "Invalid Current Password")

			w.Header().Set("HX-Reswap", "innerHTML")
			w.WriteHeader(http.StatusUnprocessableEntity)

			data := app.newTemplateData(r)
			data.Form = form

			tmpl := app.TemplateCache["userpassword.html"]
			err = tmpl.ExecuteTemplate(w, "password-error-block", data)
			if err != nil {
				app.serverError(w, r, err)
			}
			return
		}

		app.serverError(w, r, err)
		return
	}

	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Remove(r.Context(), "authenticatedUserID")
	app.sessionManager.Put(r.Context(), "flash", "Password Changed! Please login with your New Password.")

	w.Header().Set("HX-Redirect", "/user/login")
	w.WriteHeader(http.StatusOK)
}

func (app *Application) userAccountDelete(w http.ResponseWriter, r *http.Request) {
	id, ok := app.sessionManager.Get(r.Context(), "authenticatedUserID").(int)
	if !ok {
		app.serverError(w, r, fmt.Errorf("authenticatedUserID: type error"))
		return
	}

	err := app.users.Delete(id)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Remove(r.Context(), "authenticatedUserID")
	app.sessionManager.Put(r.Context(), "flash", "Your User Account has been deleted!")

	w.Header().Set("HX-Redirect", "/")
	w.WriteHeader(http.StatusOK)
}

func (app *Application) ping(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("OK"))
}
