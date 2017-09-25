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

package api

import (
	"encoding/json"

	"github.com/gorilla/mux"

	"net/http"
	"net/url"
	"path"

	"github.com/adrianpk/pulap/app"
	"github.com/adrianpk/pulap/bootstrap"
	"github.com/adrianpk/pulap/logger"
	"github.com/adrianpk/pulap/models"
	"github.com/adrianpk/pulap/repo"

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

// SignUp - SignUp a new User.
// Handler for HTTP Post - "/users/signup"
func SignUp(w http.ResponseWriter, r *http.Request) {
	// Decode
	var res UserResource
	err := json.NewDecoder(r.Body).Decode(&res)
	if err != nil {
		debugRequest(r)
		app.ShowError(w, app.ErrRequestParsing, err, http.StatusInternalServerError)
		return
	}
	user := &res.Data
	// Get Repo
	userRepo, err := repo.MakeUserRepository()
	if err != nil {
		app.ShowError(w, app.ErrRegistration, err, http.StatusInternalServerError)
		return
	}
	// Persist
	user.SetID()
	user.CreatedBy = user.ID
	err = userRepo.Create(user)
	if err != nil {
		app.ShowError(w, app.ErrRegistration, err, http.StatusInternalServerError)
		return
	}
	// Marshal
	user.ClearPassword()
	j, err := json.Marshal(UserResource{Data: *user})
	if err != nil {
		logger.Dump(err)
		app.ShowError(w, app.ErrResponseMarshalling, err, http.StatusInternalServerError)
		return
	}
	// Respond
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(j)
}

// Login - Authenticates the HTTP request with username and apssword
// Handler for HTTP Post - "/users/login"
func Login(w http.ResponseWriter, r *http.Request) {
	var res LoginResource
	var token string
	// Decode
	err := json.NewDecoder(r.Body).Decode(&res)
	if err != nil {
		app.ShowError(w, app.ErrRequestParsing, err, http.StatusInternalServerError)
		return
	}
	loginModel := res.Data
	loginUser := models.User{
		Username: models.ToNullsString(loginModel.Username),
		Email:    models.ToNullsString(loginModel.Email),
		Password: loginModel.Password,
	}
	// Get repo
	userRepo, err := repo.MakeUserRepository()
	if err != nil {
		app.ShowError(w, app.ErrLogin, err, http.StatusInternalServerError)
		return
	}
	// Authenticate the logged in user
	user, err := userRepo.Login(loginUser)
	if err != nil {
		app.ShowError(w, app.ErrLoginDenied, err, http.StatusUnauthorized)
		return
	}
	// Generate JWT token
	token, err = bootstrap.GenerateJWT(user.ID.String, user.Username.String, "member")
	if err != nil {
		app.ShowError(w, app.ErrLoginTokenCreate, err, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	// Clean-up the hashpassword to eliminate it from response JSON
	user.PasswordHash = ""
	authUser := AuthUserModel{
		User:  user,
		Token: token,
	}
	// Marshal
	j, err := json.Marshal(AuthUserResource{Data: authUser})
	if err != nil {
		app.ShowError(w, app.ErrResponseMarshalling, err, http.StatusInternalServerError)
		return
	}
	// Respond
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

// GetUsers - Returns a collection containing all users.
// Handler for HTTP Get - "/users"
func GetUsers(w http.ResponseWriter, r *http.Request) {
	// Get repo
	userRepo, err := repo.MakeUserRepository()
	if err != nil {
		app.ShowError(w, app.ErrEntityNotFound, err, http.StatusInternalServerError)
		return
	}
	// Select
	users, err := userRepo.GetAll()
	if err != nil {
		app.ShowError(w, app.ErrEntitySelect, err, http.StatusInternalServerError)
		return
	}
	// Marshal
	j, err := json.Marshal(UsersResource{Data: users})
	if err != nil {
		app.ShowError(w, app.ErrResponseMarshalling, err, http.StatusInternalServerError)
		return
	}
	// Respond
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

// CreateUser - Creates a new User.
// Handler for HTTP Post - "/users/create"
func CreateUser(w http.ResponseWriter, r *http.Request) {
	// Decode
	var res UserResource
	err := json.NewDecoder(r.Body).Decode(&res)
	if err != nil {
		app.ShowError(w, app.ErrRequestParsing, err, http.StatusInternalServerError)
		return
	}
	user := &res.Data
	// Get repo
	userRepo, err := repo.MakeUserRepository()
	if err != nil {
		app.ShowError(w, app.ErrEntityCreate, err, http.StatusInternalServerError)
		return
	}
	// Persist
	userRepo.Create(user)
	if err != nil {
		app.ShowError(w, app.ErrEntityCreate, err, http.StatusInternalServerError)
		return
	}
	// Marshal
	user.ClearPassword()
	j, err := json.Marshal(UserResource{Data: *user})
	if err != nil {
		app.ShowError(w, app.ErrResponseMarshalling, err, http.StatusNoContent)
		return
	}
	// Respond
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(j)
}

// GetUser - Returns a single User by its id or username.
// Handler for HTTP Get - "/users/{user}"
func GetUser(w http.ResponseWriter, r *http.Request) {
	// Get ID
	vars := mux.Vars(r)
	key := vars["user"]
	if isUUID(key) {
		GetUserByID(w, r)
	} else {
		GetUserByUsername(w, r)
	}
}

// GetUserByID - Returns a single User by its id.
// Handler for HTTP Get - "/users/{user}"
func GetUserByID(w http.ResponseWriter, r *http.Request) {
	// Get ID
	vars := mux.Vars(r)
	id := vars["user"]
	// Get repo
	userRepo, err := repo.MakeUserRepository()
	if err != nil {
		app.ShowError(w, app.ErrEntitySelect, err, http.StatusInternalServerError)
		return
	}
	// Select
	user, err := userRepo.Get(id)
	if err != nil {
		app.ShowError(w, app.ErrEntitySelect, err, http.StatusInternalServerError)
		return
	}
	// Marshal
	j, err := json.Marshal(user)
	if err != nil {
		app.ShowError(w, app.ErrResponseMarshalling, err, http.StatusInternalServerError)
		return
	}
	// Repsond
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

// GetUserByUsername - Returns a single User by its username.
// Handler for HTTP Get - "/users/{user}"
func GetUserByUsername(w http.ResponseWriter, r *http.Request) {
	// Get ID
	vars := mux.Vars(r)
	username := vars["user"]
	// Get repo
	userRepo, err := repo.MakeUserRepository()
	if err != nil {
		app.ShowError(w, app.ErrEntitySelect, err, http.StatusInternalServerError)
		return
	}
	// Select
	user, err := userRepo.GetByUsername(username)
	if err != nil {
		app.ShowError(w, app.ErrEntitySelect, err, http.StatusInternalServerError)
		return
	}
	// Marshal
	j, err := json.Marshal(user)
	if err != nil {
		app.ShowError(w, app.ErrResponseMarshalling, err, http.StatusInternalServerError)
		return
	}
	// Repond
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

// UpdateUser - Update an existing User.
// Handler for HTTP Put - "/users/{user}"
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	// Get ID
	vars := mux.Vars(r)
	id := vars["user"]
	// Decode
	var res UserResource
	err := json.NewDecoder(r.Body).Decode(&res)
	if err != nil {
		app.ShowError(w, app.ErrRequestParsing, err, http.StatusInternalServerError)
		return
	}
	user := &res.Data
	user.ID = models.ToNullsString(id)
	// Get repo
	userRepo, err := repo.MakeUserRepository()
	if err != nil {
		app.ShowError(w, app.ErrEntityUpdate, err, http.StatusInternalServerError)
		return
	}
	// Check against current user
	currentUser, err := userRepo.Get(id)
	if err != nil {
		app.ShowError(w, app.ErrEntityUpdate, err, http.StatusUnauthorized)
		return
	}
	// Avoid ID spoofing
	err = verifyID(user.IdentifiableModel, currentUser.IdentifiableModel)
	if err != nil {
		app.ShowError(w, app.ErrEntityUpdate, err, http.StatusUnauthorized)
		return
	}
	// Update
	err = userRepo.Update(user)
	if err != nil {
		app.ShowError(w, app.ErrEntityUpdate, err, http.StatusInternalServerError)
		return
	}
	// Marshal
	j, err := json.Marshal(UserResource{Data: *user})
	if err != nil {
		app.ShowError(w, app.ErrResponseMarshalling, err, http.StatusNoContent)
		return
	}
	// Respond
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
	w.Write(j)
}

// DeleteUser - Deletes an existing User
// Handler for HTTP Delete - "/users/{id}"
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	// Get ID
	vars := mux.Vars(r)
	id := vars["user"]
	// Get repo
	userRepo, err := repo.MakeUserRepository()
	if err != nil {
		app.ShowError(w, app.ErrEntitySelect, err, http.StatusInternalServerError)
		return
	}
	// Delete
	err = userRepo.Delete(id)
	if err != nil {
		app.ShowError(w, app.ErrEntityDelete, err, http.StatusInternalServerError)
		return
	}
	// Respond
	w.WriteHeader(http.StatusNoContent)
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
