package app

import (
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/tiwanakd/mythoughts-go/internal/assert"
	"github.com/tiwanakd/mythoughts-go/internal/testutils"
)

func TestPing(t *testing.T) {
	// rr := httptest.NewRecorder()

	// r, err := http.NewRequest(http.MethodGet, "/", nil)
	// if err != nil {
	// 	t.Fatal(err)
	// }

	// ping(rr, r) p

	app := newTestApplication(t)
	ts := testutils.NewTestServer(t, app.Routes())
	defer ts.Close()

	code, _, body := ts.Get(t, "/ping")

	assert.Equal(t, code, http.StatusOK)
	assert.Equal(t, body, "OK")
}

func TestHomeView(t *testing.T) {
	app := newTestApplication(t)

	ts := testutils.NewTestServer(t, app.Routes())
	defer ts.Close()

	code, _, body := ts.Get(t, "/")

	assert.Equal(t, code, http.StatusOK)
	assert.StringContains(t, body, "You will come over trust me... All this will pass.")
	assert.StringContains(t, body, "You need to forgvie yourself dude... There is a lot more to life.")

	t.Run("unautheticated", func(t *testing.T) {
		if strings.Contains(`<button class="new-tg-btn" id="new-btn">New &plus;</button>`, body) {
			t.Error("from contains new Thought butting; user is authenticated.")
		}
	})

	t.Run("authenticated", func(t *testing.T) {
		_, _, body := ts.Get(t, "/user/login")
		csrfToken := extractCSRFToken(t, body)

		form := url.Values{}
		form.Add("csrf_token", csrfToken)
		form.Add("email", "test@test.com")
		form.Add("password", "pa$$word")

		code, _, _ := ts.PostForm(t, "/user/login", form)
		assert.Equal(t, code, http.StatusSeeOther)

		code, _, body = ts.Get(t, "/")
		assert.Equal(t, code, http.StatusOK)
		assert.StringContains(t, body, `<button class="new-tg-btn" id="new-btn">New &plus;</button>`)
	})
}

func TestNewThoughtPost(t *testing.T) {
	app := newTestApplication(t)
	ts := testutils.NewTestServer(t, app.Routes())
	defer ts.Close()

	_, _, body := ts.Get(t, "/user/login")
	csrfToken := extractCSRFToken(t, body)

	form := url.Values{}
	form.Add("csrf_token", csrfToken)
	form.Add("email", "test@test.com")
	form.Add("password", "pa$$word")

	code, _, _ := ts.PostForm(t, "/user/login", form)
	assert.Equal(t, code, http.StatusSeeOther)

	form.Add("newThought", "This is a test post from a test Application and Test Server.")

	code, _, _ = ts.PostForm(t, "/thought/new", form)
	assert.Equal(t, code, http.StatusCreated)
}

func TestUserSignup(t *testing.T) {
	app := newTestApplication(t)
	ts := testutils.NewTestServer(t, app.Routes())
	defer ts.Close()

	_, _, body := ts.Get(t, "/user/signup")
	validCSRFToken := extractCSRFToken(t, body)

	const (
		validName     = "Test User"
		validUsername = "test"
		validEmail    = "test@email.com"
		validPassword = "Password123"
		formTag       = `<form action="/user/signup" method="post" novalidate>`
	)

	tests := []struct {
		name         string
		userName     string
		userUname    string
		userEmail    string
		userPassword string
		csrfToken    string
		wantCode     int
		wantFormTag  string
	}{
		{
			name:         "Valid Submisssion",
			userName:     validName,
			userUname:    validUsername,
			userEmail:    validEmail,
			userPassword: validPassword,
			csrfToken:    validCSRFToken,
			wantCode:     http.StatusSeeOther,
		},
		{
			name:         "Invalid CSRF Token",
			userName:     validName,
			userUname:    validUsername,
			userEmail:    validEmail,
			userPassword: validPassword,
			csrfToken:    "invalidToken",
			wantCode:     http.StatusBadRequest,
		},
		{
			name:         "Empty Name",
			userName:     "",
			userUname:    validUsername,
			userEmail:    validEmail,
			userPassword: validPassword,
			csrfToken:    validCSRFToken,
			wantCode:     http.StatusUnprocessableEntity,
			wantFormTag:  formTag,
		},
		{
			name:         "Empty Username",
			userName:     validName,
			userUname:    "",
			userEmail:    validEmail,
			userPassword: validPassword,
			csrfToken:    validCSRFToken,
			wantCode:     http.StatusUnprocessableEntity,
			wantFormTag:  formTag,
		},
		{
			name:         "Empty Email",
			userName:     validName,
			userUname:    validUsername,
			userEmail:    "",
			userPassword: validPassword,
			csrfToken:    validCSRFToken,
			wantCode:     http.StatusUnprocessableEntity,
			wantFormTag:  formTag,
		},
		{
			name:         "Empty Password",
			userName:     validName,
			userUname:    validUsername,
			userEmail:    validEmail,
			userPassword: "",
			csrfToken:    validCSRFToken,
			wantCode:     http.StatusUnprocessableEntity,
			wantFormTag:  formTag,
		},
		{
			name:         "Invalid Email",
			userName:     validName,
			userUname:    validUsername,
			userEmail:    "invalid@email.",
			userPassword: validPassword,
			csrfToken:    validCSRFToken,
			wantCode:     http.StatusUnprocessableEntity,
			wantFormTag:  formTag,
		},
		{
			name:         "Passowrd Complexity not met",
			userName:     validName,
			userUname:    validUsername,
			userEmail:    validEmail,
			userPassword: "password",
			csrfToken:    validCSRFToken,
			wantCode:     http.StatusUnprocessableEntity,
			wantFormTag:  formTag,
		},
		{
			name:         "Duplicate Username",
			userName:     validName,
			userUname:    "dupeUsername",
			userEmail:    validEmail,
			userPassword: validPassword,
			csrfToken:    validCSRFToken,
			wantCode:     http.StatusUnprocessableEntity,
			wantFormTag:  formTag,
		},
		{
			name:         "Duplicate Email",
			userName:     validName,
			userUname:    validUsername,
			userEmail:    "dupe@email.com",
			userPassword: validPassword,
			csrfToken:    validCSRFToken,
			wantCode:     http.StatusUnprocessableEntity,
			wantFormTag:  formTag,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			form := url.Values{}
			form.Add("name", tt.userName)
			form.Add("username", tt.userUname)
			form.Add("email", tt.userEmail)
			form.Add("password", tt.userPassword)
			form.Add("csrf_token", tt.csrfToken)

			code, _, body := ts.PostForm(t, "/user/signup", form)

			assert.Equal(t, code, tt.wantCode)

			if tt.wantFormTag != "" {
				assert.StringContains(t, body, tt.wantFormTag)
			}
		})
	}
}

func TestAddLikeandDislike(t *testing.T) {
	app := newTestApplication(t)
	ts := testutils.NewTestServer(t, app.Routes())
	defer ts.Close()

	_, _, body := ts.Get(t, "/")

	validCSRFToken := extractCSRFToken(t, body)

	tests := []struct {
		name      string
		urlPath   string
		csrfToken string
		wantCode  int
		wantBody  string
	}{
		{
			name:      "Add like to Mock 1",
			urlPath:   "/like/1",
			csrfToken: validCSRFToken,
			wantCode:  http.StatusCreated,
			wantBody:  "4",
		},
		{
			name:      "Add like to Mock 2",
			urlPath:   "/like/2",
			csrfToken: validCSRFToken,
			wantCode:  http.StatusCreated,
			wantBody:  "2",
		},
		{
			name:      "Add dislike to Mock 1",
			urlPath:   "/dislike/1",
			csrfToken: validCSRFToken,
			wantCode:  http.StatusCreated,
			wantBody:  "1",
		},
		{
			name:      "Add dislike to Mock 2",
			urlPath:   "/dislike/2",
			csrfToken: validCSRFToken,
			wantCode:  http.StatusCreated,
			wantBody:  "6",
		},
		{
			name:      "Invalid CSRF Token: Mock 1",
			urlPath:   "/like/1",
			csrfToken: "invalidToken",
			wantCode:  http.StatusBadRequest,
		},
		{
			name:      "Invalid CSRF Token: Mock 2",
			urlPath:   "/like/2",
			csrfToken: "invalidToken",
			wantCode:  http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			form := url.Values{}
			form.Add("csrf_token", tt.csrfToken)

			code, _, body := ts.PostForm(t, tt.urlPath, form)

			assert.Equal(t, code, tt.wantCode)

			if tt.wantBody != "" {
				assert.Equal(t, body, tt.wantBody)
			}
		})
	}
}

func TestUserLogin(t *testing.T) {
	app := newTestApplication(t)
	ts := testutils.NewTestServer(t, app.Routes())
	defer ts.Close()

	code, _, body := ts.Get(t, "/user/login")
	assert.Equal(t, code, http.StatusOK)
	assert.StringContains(t, body, `<form action="/user/login" method="post" novalidate>`)
}

func TestUserLoginPost(t *testing.T) {
	app := newTestApplication(t)
	ts := testutils.NewTestServer(t, app.Routes())
	defer ts.Close()

	_, _, body := ts.Get(t, "/user/login")
	csrfToken := extractCSRFToken(t, body)

	form := url.Values{}
	form.Add("csrf_token", csrfToken)
	form.Add("email", "test@test.com")
	form.Add("password", "pa$$word")

	code, _, _ := ts.PostForm(t, "/user/login", form)
	assert.Equal(t, code, http.StatusSeeOther)
}

func TestUserLogout(t *testing.T) {
	app := newTestApplication(t)
	ts := testutils.NewTestServer(t, app.Routes())
	defer ts.Close()

	_, _, body := ts.Get(t, "/")
	validCSRFToken := extractCSRFToken(t, body)

	form := url.Values{}
	form.Add("csrf_token", validCSRFToken)

	code, _, body := ts.PostForm(t, "/user/logout", form)
	assert.Equal(t, code, http.StatusSeeOther)
}

func TestDeleteThoughtPost(t *testing.T) {
	app := newTestApplication(t)
	ts := testutils.NewTestServer(t, app.Routes())
	defer ts.Close()

	_, _, body := ts.Get(t, "/user/login")
	csrfToken := extractCSRFToken(t, body)

	form := url.Values{}
	form.Add("csrf_token", csrfToken)
	form.Add("email", "test@test.com")
	form.Add("password", "pa$$word")

	ts.PostForm(t, "/user/login", form)

	code, _, _ := ts.Delete(t, "/user/thought/delete/1", csrfToken)
	assert.Equal(t, code, http.StatusOK)
}

func TestUserAccountDelete(t *testing.T) {
	app := newTestApplication(t)
	ts := testutils.NewTestServer(t, app.Routes())
	defer ts.Close()

	_, _, body := ts.Get(t, "/user/login")
	csrfToken := extractCSRFToken(t, body)

	form := url.Values{}
	form.Add("csrf_token", csrfToken)
	form.Add("email", "test@test.com")
	form.Add("password", "pa$$word")

	ts.PostForm(t, "/user/login", form)

	code, _, _ := ts.Delete(t, "/user/account/delete", csrfToken)
	assert.Equal(t, code, http.StatusOK)
}

func TestUserAccountView(t *testing.T) {
	app := newTestApplication(t)
	ts := testutils.NewTestServer(t, app.Routes())
	defer ts.Close()

	_, _, body := ts.Get(t, "/user/login")
	csrfToken := extractCSRFToken(t, body)

	form := url.Values{}
	form.Add("csrf_token", csrfToken)
	form.Add("email", "test@test.com")
	form.Add("password", "pa$$word")

	ts.PostForm(t, "/user/login", form)

	code, _, body := ts.Get(t, "/user/account")
	assert.Equal(t, code, http.StatusOK)
	assert.StringContains(t, body, `<div class="user-info-container">`)
}
