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
	"encoding/json"
	"net/http"

	"github.com/adrianpk/fundacja/app"
	"github.com/adrianpk/fundacja/logger"
	"github.com/adrianpk/fundacja/models"
	"github.com/adrianpk/fundacja/repo"

	"github.com/gorilla/mux"
)

// GetUserProfile - Returns the profile for a User referenced by its ID or username.
// Handler for HTTP Get - "/users/{user}/}profile"
func GetUserProfile(w http.ResponseWriter, r *http.Request) {
	// Get ID
	vars := mux.Vars(r)
	key := vars["user"]
	if len(key) == 36 {
		GetProfileByUserID(w, r)
	} else {
		GetProfileByUsername(w, r)
	}
}

// GetProfileByUserID - Returns the profile for a User referenced by its ID.
// Handler for HTTP Get - "/users/{user}/}profile"
func GetProfileByUserID(w http.ResponseWriter, r *http.Request) {
	// Get ID
	vars := mux.Vars(r)
	id := vars["user"]
	// Get repo
	profileRepo, err := repo.MakeProfileRepository()
	if err != nil {
		app.ShowError(w, app.ErrEntitySelect, err, http.StatusInternalServerError)
		return
	}
	// Select
	profile, err := profileRepo.GetByUserID(id)
	if err != nil {
		app.ShowError(w, app.ErrEntitySelect, err, http.StatusInternalServerError)
		return
	}
	// Marshal
	j, err := json.Marshal(profile)
	if err != nil {
		app.ShowError(w, app.ErrResponseMarshalling, err, http.StatusInternalServerError)
		return
	}
	// Output
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

// GetProfileByUsername - Returns the profile for a User referenced by its username or username.
// Handler for HTTP Get - "/users/{user}/}profile"
func GetProfileByUsername(w http.ResponseWriter, r *http.Request) {
	// Get ID
	vars := mux.Vars(r)
	username := vars["user"]
	// Get User from username
	user, err := getUserByUsername(username)
	if err != nil {
		app.ShowError(w, app.ErrEntitySelect, err, http.StatusInternalServerError)
		return
	}
	// Get Repo
	profileRepo, err := repo.MakeProfileRepository()
	if err != nil {
		app.ShowError(w, app.ErrEntitySelect, err, http.StatusInternalServerError)
		return
	}
	// Select
	profile, err := profileRepo.GetByUserID(user.ID.String)
	if err != nil {
		app.ShowError(w, app.ErrEntitySelect, err, http.StatusInternalServerError)
		return
	}
	// Marshal
	j, err := json.Marshal(profile)
	if err != nil {
		app.ShowError(w, app.ErrResponseMarshalling, err, http.StatusInternalServerError)
		return
	}
	// Output
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

func getUserByUsername(username string) (models.User, error) {
	userRepo, err := repo.MakeUserRepository()
	if err != nil {
		return models.User{}, err
	}
	user, err := userRepo.GetByUsername(username)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

// CreateUserProfile creates a Profile for a given User
// Handler for HTTP Post - "/users/{user}/profile"
func CreateUserProfile(w http.ResponseWriter, r *http.Request) {
	// Decode
	var res ProfileResource
	err := json.NewDecoder(r.Body).Decode(&res)
	if err != nil {
		app.ShowError(w, app.ErrRequestParsing, err, http.StatusInternalServerError)
		return
	}
	// Set Owner
	profile := &res.Data
	user, err := sessionUser(r)
	if err != nil {
		logger.Debug("00000")
		app.ShowError(w, app.ErrEntityCreate, err, http.StatusInternalServerError)
		return
	}
	profile.SetOwner(user)
	// Set Values
	profile.ValidableDate.Date = &profile.AnniversaryDate
	profile.ValidableDate.Validate()
	// Get repo
	profileRepo, err := repo.MakeProfileRepository()
	if err != nil {
		app.ShowError(w, app.ErrEntityCreate, err, http.StatusInternalServerError)
		return
	}
	// Persist
	user.SetID()
	profileRepo.Create(profile)
	if err != nil {
		app.ShowError(w, app.ErrEntityCreate, err, http.StatusInternalServerError)
		return
	}
	// Marshal
	j, err := json.Marshal(ProfileResource{Data: *profile})
	if err != nil {
		app.ShowError(w, app.ErrResponseMarshalling, err, http.StatusNoContent)
		return
	}

	// Respond
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(j)
}

// UpdateUserProfile updates a Users's Profile
// Handler for HTTP Put - "/users/{user}/profile"
func UpdateUserProfile(w http.ResponseWriter, r *http.Request) {
	// Decode
	var res ProfileResource
	err := json.NewDecoder(r.Body).Decode(&res)
	if err != nil {
		app.ShowError(w, app.ErrRequestParsing, err, http.StatusInternalServerError)
		return
	}
	profile := &res.Data
	// Verify ownership
	err = verifyOwnership(profile.UserID.String, r)
	if err != nil {
		app.ShowError(w, app.ErrOwnerOnlyCanManage, err, http.StatusUnauthorized)
		return
	}
	// Current User's profile from repo
	userID, _ := sessionUserID(r)
	profileRepo, err := repo.MakeProfileRepository()
	if err != nil {
		app.ShowError(w, app.ErrEntityCreate, err, http.StatusInternalServerError)
		return
	}
	// Check against current profile
	//logger.Debugf("UserID: %s", userID)
	currentProfile, err := profileRepo.GetByUserID(userID)
	if err != nil {
		app.ShowError(w, app.ErrEntityUpdate, err, http.StatusUnauthorized)
		return
	}
	// Avoid ID spoofing
	err = verifyID(profile.IdentifiableModel, currentProfile.IdentifiableModel)
	if err != nil {
		app.ShowError(w, app.ErrEntityUpdate, err, http.StatusUnauthorized)
		return
	}
	// Update
	profile.ValidableDate.Date = &profile.AnniversaryDate
	profile.ValidableDate.Validate()
	profileRepo.Update(profile)
	if err != nil {
		logger.Debug("2222")
		app.ShowError(w, app.ErrEntityUpdate, err, http.StatusInternalServerError)
		return
	}
	// Respond
	w.WriteHeader(http.StatusNoContent)
}

// DeleteUserProfile delete a Users's Profile
// Handler for HTTP Delete - "/users/{user}/profile"
func DeleteUserProfile(w http.ResponseWriter, r *http.Request) {
	// Get ID
	vars := mux.Vars(r)
	id := vars["user"]
	// Get repo
	profileRepo, err := repo.MakeProfileRepository()
	if err != nil {
		app.ShowError(w, app.ErrEntitySelect, err, http.StatusInternalServerError)
		return
	}
	// Delete
	err = profileRepo.DeleteByUserID(id)
	if err != nil {
		app.ShowError(w, app.ErrEntitySelect, err, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
