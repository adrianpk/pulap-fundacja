// Copyright (c) 2017 Kuguar <licenses@kuguar.io> Author: Adrian P.K. <apk@kuguar.io>
//
// MIT License
//
// Permission is hereby granted, free of charge, to any person obtaining
// a copy of this software and associated documentation files (the
// "Software"), to deal in the Software without restriction, including
// without limitation the rights to use, copy, modify, merge, publish,
// distribute, sublicense, and/or sell copies of the Software, and to
// permit persons to whom the Software is furnished to do so, subject to
// the following conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
// LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
// OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
// WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package controllers

import (
	htmlTemplate "html/template"

	"github.com/arschles/go-bindata-html-template"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"

	"net/http"
	"net/url"
	"path"

	"github.com/adrianpk/fundacja/app"
	"github.com/adrianpk/fundacja/bootstrap"
	"github.com/adrianpk/fundacja/logger"
	"github.com/adrianpk/fundacja/models"
	"github.com/adrianpk/fundacja/repo"

	_ "github.com/lib/pq" // Import pq without side effects
)

// get "users", to: "users#index"
// # http GET method, In this instance, gets a list of users
// get "users/{user}", to: "users#show"
// # In this instance, gets a specific user via the provided id.  For example: /users/3882
// post "users", to: "users#create"
// # http POST method, In this instance, used for creating a user.
// put "users/{user}"
// # http PUT method, used for updating a resource, In this instance, updates the user.  Older versions of Rails used this for all updates.
// patch "users/{user}"
// #  http PATCH method.  in this instance, used to partially update the user's information.
// delete "users/{user}"
// # http DELETE method.  In this instance, used to delete a user.

const (
	rememberField = "remember"
)

var (
	userAssetsBase    string //resources/templates/user
	userExtAssetsBase string
	userTemplates     map[string]*template.Template
	userExtTemplates  map[string]*htmlTemplate.Template
	userIndex         = "/users"
	userNew           = "/users/new"
	userEdit          = "/users/edit/%s"
	userShow          = "/users/%s"
	userDelete        = "/users/delete/%s"
)

// InitializeUser - Initialize the controller
func InitializeUser() {
	if userTemplates == nil {
		userTemplates = make(map[string]*template.Template)
	}
	if userExtTemplates == nil {
		userExtTemplates = make(map[string]*htmlTemplate.Template)
	}
	parseUserAssets()
	parseUserExtAssets()
}

// ShowSignUp - Shows SignUp form.
// Handler for HTTP Post - "/users/signup"
func ShowSignUp(w http.ResponseWriter, r *http.Request) {
	pageModel := makePage(nil, nil)
	renderUserTemplate(w, r, signupView, layoutView, pageModel)
}

// SignUp - SignUp a new User.
// Handler for HTTP Post - "/users/signup"
func SignUp(w http.ResponseWriter, r *http.Request) {
	logger.Debug("SignUp...")
	// Parse
	err := r.ParseForm()
	if err != nil {
		showUserError(w, r, signupView, layoutView, nil, app.ErrRegistration, warningAlert, err)
		return
	}
	// Decode
	var user models.User
	err = schema.NewDecoder().Decode(&user, r.Form)
	if err != nil {
		showUserError(w, r, signupView, layoutView, nil, app.ErrRegistration, warningAlert, err)
		return
	}
	// Get Repo
	userRepo, err := repo.MakeUserRepository()
	if err != nil {
		showUserError(w, r, signupView, layoutView, nil, app.ErrRegistration, warningAlert, err)
		return
	}
	// Persist
	user.SetID()
	user.CreatedBy = user.ID
	err = userRepo.Create(&user)
	if err != nil {
		user.ClearPassword()
		showUserError(w, r, signupView, layoutView, user, app.ErrRegistration, warningAlert, err)
		return
	}
	// Clear password field
	user.ClearPassword()
	// Respond
	// redirectTo(w, r, userIndex, nil)
	renderUserTemplate(w, r, editView, layoutView, makePage(user, nil))
}

// ShowLogin - Shows SignUp form.
// Handler for HTTP Post - "/users/signup"
func ShowLogin(w http.ResponseWriter, r *http.Request) {
	logger.Debug("ShowLogin...")
	pageModel := makePage(models.User{}, nil)
	renderUserTemplate(w, r, loginView, layoutView, pageModel)
}

// Login - Authenticates the HTTP request with username and apssword
// Handler for HTTP Post - "/users/login"
func Login(w http.ResponseWriter, r *http.Request) {
	logger.Debug("Login...")
	debugRequest(r)
	err := r.ParseForm()
	if err != nil {
		showUserError(w, r, loginView, layoutView, nil, app.ErrLogin, warningAlert, err)
		return
	}
	// Decode
	var toLogin models.User
	err = getFormDecoder(true).Decode(&toLogin, r.Form)
	if err != nil {
		showUserError(w, r, loginView, layoutView, toLogin, app.ErrRequestParsing, warningAlert, err)
		return
	}
	userRepo, err := repo.MakeUserRepository()
	if err != nil {
		showUserError(w, r, loginView, layoutView, toLogin, app.ErrLogin, warningAlert, err)
		return
	}
	// Authenticate the logged in user
	user, err := userRepo.Login(toLogin)
	if err != nil {
		showUserError(w, r, loginView, layoutView, toLogin, app.ErrLoginDenied, warningAlert, err)
		return
	}
	toLogin.PasswordHash = ""
	// Create session
	err = setSession(w, r, user, inputIsTrue(r, rememberField))
	if err != nil {
		showUserError(w, r, loginView, layoutView, toLogin, app.ErrLoginSessionCreate, warningAlert, err)
		return
	}
	// redirectTo(w, r, userIndex, nil)
	renderUserTemplate(w, r, editView, layoutView, makePage(user, nil))
}

// IndexUsers - Returns a collection containing all users.
// Handler for HTTP Get - "/users"
func IndexUsers(w http.ResponseWriter, r *http.Request) {
	logger.Debug("IndexUsers...")
	// Get repo
	userRepo, err := repo.MakeUserRepository()
	if err != nil {
		showUserError(w, r, indexView, layoutView, nil, app.ErrEntitySelect, warningAlert, err)
		return
	}
	// Select
	users, err := userRepo.GetAll()
	if err != nil {
		showUserError(w, r, indexView, layoutView, nil, app.ErrEntitySelect, warningAlert, err)
		return
	}
	// Respond
	pageModel := makePage(users, nil)
	renderUserTemplate(w, r, indexView, layoutView, pageModel)
}

// NewUser - Presents new user form.
// Handler for HTTP Get - "/users/new"
func NewUser(w http.ResponseWriter, r *http.Request) {
	logger.Debug("NewUser...")
	renderUserTemplate(w, r, newView, layoutView, makePage(nil, nil))
}

// CreateUser - Creates a new User.
// Handler for HTTP Post - "/users/create"
func CreateUser(w http.ResponseWriter, r *http.Request) {
	logger.Debug("CreateUser...")
	// Parse
	err := r.ParseForm()
	if err != nil {
		showUserError(w, r, newView, layoutView, nil, app.ErrEntityCreate, warningAlert, err)
		return
	}
	// Decode
	var user models.User
	err = schema.NewDecoder().Decode(&user, r.Form)
	if err != nil {
		showUserError(w, r, newView, layoutView, user, app.ErrEntityCreate, warningAlert, err)
		return
	}
	// Get repo
	userRepo, err := repo.MakeUserRepository()
	if err != nil {
		showUserError(w, r, newView, layoutView, user, app.ErrEntityCreate, warningAlert, err)
		return
	}
	// Persist
	userRepo.Create(&user)
	user.ClearPassword()
	if err != nil {
		showUserError(w, r, newView, layoutView, user, app.ErrDataAccess, warningAlert, err)
		return
	}
	// Respond
	redirectTo(w, r, userIndex, makePageAlert("User created", infoAlert))
}

// ShowUser - Returns a single User by its id or username.
// Handler for HTTP Get - "/users/{user}"
func ShowUser(w http.ResponseWriter, r *http.Request) {
	logger.Debug("ShowUser...")
	// Get ID
	vars := mux.Vars(r)
	key := vars["user"]
	if isUUID(key) {
		ShowUserByID(w, r)
	} else {
		ShowUserByUsername(w, r)
	}
}

// ShowUserByID - Returns a single User by its id.
// Handler for HTTP Get - "/users/{user}"
func ShowUserByID(w http.ResponseWriter, r *http.Request) {
	logger.Debug("ShowUserByID...")
	// Get ID
	vars := mux.Vars(r)
	id := vars["user"]
	// Get repo
	userRepo, err := repo.MakeUserRepository()
	if err != nil {
		redirectTo(w, r, userIndex, makePageAlert(app.ErrDataStore.Error(), warningAlert))
		return
	}
	// Select
	user, err := userRepo.Get(id)
	if err != nil {
		redirectTo(w, r, userIndex, makePageAlert(app.ErrDataAccess.Error(), warningAlert))
		return
	}
	renderUserTemplate(w, r, showView, layoutView, makePage(user, nil))
}

// ShowUserByUsername - Returns a single User by its username.
// Handler for HTTP Get - "/users/{user}"
func ShowUserByUsername(w http.ResponseWriter, r *http.Request) {
	logger.Debug("ShowUserByUsername...")
	// Get ID
	vars := mux.Vars(r)
	username := vars["user"]
	// Get repo
	userRepo, err := repo.MakeUserRepository()
	if err != nil {
		redirectTo(w, r, userIndex, makePageAlert(app.ErrDataStore.Error(), warningAlert))
		return
	}
	// Select
	user, err := userRepo.GetByUsername(username)
	if err != nil {
		redirectTo(w, r, userIndex, makePageAlert(app.ErrDataAccess.Error(), warningAlert))
		return
	}
	// Repond
	renderUserTemplate(w, r, showView, layoutView, makePage(user, nil))
}

// EditUser - Presents edit user form.
// Handler for HTTP Get - "/users/edit"
func EditUser(w http.ResponseWriter, r *http.Request) {
	logger.Debug("EditUser...")
	// Get ID
	vars := mux.Vars(r)
	id := vars["user"]
	// Get repo
	userRepo, err := repo.MakeUserRepository()
	if err != nil {
		redirectTo(w, r, userIndex, makePageAlert(app.ErrDataStore.Error(), warningAlert))
		return
	}
	// Select
	user, err := userRepo.Get(id)
	if err != nil {
		redirectTo(w, r, userIndex, makePageAlert(app.ErrDataAccess.Error(), warningAlert))
		return
	}
	renderUserTemplate(w, r, editView, layoutView, makePage(user, nil))
	return
}

// UpdateUser - Update an existing User.
// Handler for HTTP Put - "/users/{user}"
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	logger.Debug("UpdateUser...")
	// Get ID
	vars := mux.Vars(r)
	id := vars["user"]
	// Parse
	err := r.ParseForm()
	if err != nil {
		showUserError(w, r, signupView, layoutView, nil, app.ErrRegistration, warningAlert, err)
		return
	}
	// Decode
	var user models.User
	err = schema.NewDecoder().Decode(&user, r.Form)
	if err != nil {
		showUserError(w, r, editView, layoutView, user, app.ErrRegistration, warningAlert, err)
		return
	}
	user.ID = models.ToNullsString(id)
	// Get repo
	userRepo, err := repo.MakeUserRepository()
	if err != nil {
		showUserError(w, r, editView, layoutView, user, app.ErrEntityUpdate, warningAlert, err)
		return
	}
	// Check against current user
	currentUser, err := userRepo.Get(id)
	if err != nil {
		showUserError(w, r, editView, layoutView, user, app.ErrEntityUpdate, warningAlert, err)
		return
	}
	// Avoid ID spoofing
	err = verifyID(user.IdentifiableModel, currentUser.IdentifiableModel)
	if err != nil {
		showUserError(w, r, editView, layoutView, currentUser, app.ErrEntityUpdate, warningAlert, err)
		return
	}
	// Update
	err = userRepo.Update(&user)
	if err != nil {
		showUserError(w, r, editView, layoutView, currentUser, app.ErrEntityUpdate, warningAlert, err)
		return
	}
	// Respond
	redirectTo(w, r, userIndex, makePageAlert("User updated", warningAlert))
}

// InitDeleteUser - Show user deletion page.
// Handler for HTTP Get - "/users/init-delete/{user}"
func InitDeleteUser(w http.ResponseWriter, r *http.Request) {
	logger.Debug("InitDeleteUser...")
	// Get ID
	vars := mux.Vars(r)
	id := vars["user"]
	// Get repo
	userRepo, err := repo.MakeUserRepository()
	if err != nil {
		redirectTo(w, r, userIndex, makePageAlert(app.ErrDataStore.Error(), warningAlert))
		return
	}
	// Select
	user, err := userRepo.Get(id)
	if err != nil {
		redirectTo(w, r, userIndex, makePageAlert(app.ErrDataAccess.Error(), warningAlert))
		return
	}
	renderUserTemplate(w, r, deleteView, layoutView, makePage(user, nil))
}

// DeleteUser - Deletes an existing User
// Handler for HTTP Delete - "/users/{id}"
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	logger.Debug("DeleteUser...")
	// Get ID
	vars := mux.Vars(r)
	id := vars["user"]
	// Get repo
	userRepo, err := repo.MakeUserRepository()
	if err != nil {
		showUserError(w, r, indexView, layoutView, nil, app.ErrEntitySelect, warningAlert, err)
		return
	}
	// Delete
	err = userRepo.Delete(id)
	if err != nil {
		showUserError(w, r, indexView, layoutView, nil, app.ErrEntityDelete, warningAlert, err)
		return
	}
	// Respond
	redirectTo(w, r, userIndex, makePageAlert("User updated", warningAlert))
}

func userIDfromURL(r *http.Request) string {
	u, _ := url.Parse(r.URL.Path)
	dir := path.Dir(u.Path)
	id := path.Base(dir)
	logger.Debugf("User id in url is %s", id)
	return id
}

func usernameFromURL(r *http.Request) string {
	u, _ := url.Parse(r.URL.Path)
	dir := path.Dir(u.Path)
	username := path.Base(dir)
	logger.Debugf("Username in url is %s", username)
	return username
}

func parseUserAssets() {
	//logger.Debug("Parsing user assets...")
	assetNames := []string{signupView, loginView, indexView, newView, showView, editView, deleteView}
	parseAssets(&userAssetsBase, "layouts", "user", layoutView, assetNames, userTemplates)
}

func parseUserExtAssets() {
	assetNames := []string{signupView, loginView, indexView, newView, showView, editView, deleteView}
	parseExtAssets(&userExtAssetsBase, "layouts", "user", layoutView, assetNames, userExtTemplates)
}

func renderUserTemplate(w http.ResponseWriter, r *http.Request, groupName string, name string, page *Page) {
	if useExtTemplates {
		renderExtUserTemplate(w, r, groupName, name, page)
		return
	}
	renderIntUserTemplate(w, r, groupName, name, page)
}

// Render templates for the given name, template definition and data object
func renderIntUserTemplate(w http.ResponseWriter, r *http.Request, groupName string, name string, page *Page) {
	renderTemplate(w, r, userAssetsBase, groupName, name, userTemplates, page)
}

// Render templates for the given name, template definition and data object
func renderExtUserTemplate(w http.ResponseWriter, r *http.Request, groupName string, name string, page *Page) {
	//logger.Debugf("Autoreload: %t", bootstrap.AppConfig.IsAutoreloadOn())
	if bootstrap.AppConfig.IsAutoreloadOn() {
		logger.Debug("Reloading templates.")
		userExtTemplates = make(map[string]*htmlTemplate.Template)
		parseUserExtAssets()
	}
	renderExtTemplate(w, r, userExtAssetsBase, groupName, name, userExtTemplates, page)
}

func showUserError(w http.ResponseWriter, r *http.Request, page string, layout string, model interface{}, err error, alertKind string, cause error) {
	logger.Dump(cause)
	pageModel := makePage(model, makePageAlert(err.Error(), alertKind))
	renderUserTemplate(w, r, page, layoutView, pageModel)
}
