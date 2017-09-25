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
	"fmt"

	"github.com/gorilla/mux"

	"github.com/adrianpk/pulap/app"
	"github.com/adrianpk/pulap/logger"
	"github.com/adrianpk/pulap/models"
	"net/http"
	"net/url"
	"path"

	_ "github.com/lib/pq" // Import pq without side effects

	"github.com/adrianpk/pulap/repo"
)

// GetUserRoles - Returns a collection containing all userRoles.
// Handler for HTTP Get - "/organizations/{organization}/user-roles"
func GetUserRoles(w http.ResponseWriter, r *http.Request) {
	// Get ID
	vars := mux.Vars(r)
	orgid := vars["organization"]
	// Get repo
	userRoleRepo, err := repo.MakeUserRoleRepository()
	if err != nil {
		app.ShowError(w, app.ErrEntityNotFound, err, http.StatusInternalServerError)
		return
	}
	// Select
	userRoles, err := userRoleRepo.GetAll(orgid)
	if err != nil {
		app.ShowError(w, app.ErrEntitySelect, err, http.StatusInternalServerError)
		return
	}
	// Marshal
	j, err := json.Marshal(UserRolesResource{Data: userRoles})
	if err != nil {
		app.ShowError(w, app.ErrResponseMarshalling, err, http.StatusInternalServerError)
		return
	}
	// Respond
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

// CreateUserRole - Creates a new UserRole.
// Handler for HTTP Post - "/organizations/{organization}/user-roles/create"
func CreateUserRole(w http.ResponseWriter, r *http.Request) {
	// Decode
	var res UserRoleResource
	err := json.NewDecoder(r.Body).Decode(&res)
	if err != nil {
		app.ShowError(w, app.ErrRequestParsing, err, http.StatusInternalServerError)
		return
	}
	userRole := &res.Data
	// Get repo
	userRoleRepo, err := repo.MakeUserRoleRepository()
	if err != nil {
		app.ShowError(w, app.ErrEntityCreate, err, http.StatusInternalServerError)
		return
	}
	// Set values
	u, _ := sessionUser(r)
	userRole.CreatedBy = u.ID
	genUserRoleNameAndDescription(userRole)
	// Persist
	err = userRoleRepo.Create(userRole)
	if err != nil {
		app.ShowError(w, app.ErrEntityCreate, err, http.StatusInternalServerError)
		return
	}
	// Marshal
	j, err := json.Marshal(UserRoleResource{Data: *userRole})
	if err != nil {
		app.ShowError(w, app.ErrResponseMarshalling, err, http.StatusNoContent)
		return
	}
	// Respond
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(j)
}

// GetUserRole - Returns a single UserRole by its id or userRoleName.
// Handler for HTTP Get - "/organizations/{organization}/user-roles/{user-role}"
func GetUserRole(w http.ResponseWriter, r *http.Request) {
	GetUserRoleByID(w, r)
}

// GetUserRoleByID - Returns a single UserRole by its id.
// Handler for HTTP Get - "/organizations/{organization}/user-roles/{user-role}"
func GetUserRoleByID(w http.ResponseWriter, r *http.Request) {
	// Get IDs
	vars := mux.Vars(r)
	orgid := vars["organization"]
	id := vars["user-role"]
	// Get repo
	userRoleRepo, err := repo.MakeUserRoleRepository()
	if err != nil {
		app.ShowError(w, app.ErrEntitySelect, err, http.StatusInternalServerError)
		return
	}
	// Select
	userRole, err := userRoleRepo.GetFromOrganization(id, orgid)
	if err != nil {
		app.ShowError(w, app.ErrEntitySelect, err, http.StatusInternalServerError)
		return
	}
	// Marshal
	j, err := json.Marshal(userRole)
	if err != nil {
		app.ShowError(w, app.ErrResponseMarshalling, err, http.StatusInternalServerError)
		return
	}
	// Repsond
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

// UpdateUserRole - Update an existing UserRole.
// Handler for HTTP Put - "/organizations/{organization}/user-roles/{user-role}"
func UpdateUserRole(w http.ResponseWriter, r *http.Request) {
	// Get IDs
	vars := mux.Vars(r)
	orgid := vars["organization"]
	id := vars["user-role"]
	// Decode
	var res UserRoleResource
	err := json.NewDecoder(r.Body).Decode(&res)
	if err != nil {
		app.ShowError(w, app.ErrRequestParsing, err, http.StatusInternalServerError)
		return
	}
	userRole := &res.Data
	userRole.ID = models.ToNullsString(id)
	// Get repo
	userRoleRepo, err := repo.MakeUserRoleRepository()
	if err != nil {
		app.ShowError(w, app.ErrEntityUpdate, err, http.StatusInternalServerError)
		return
	}
	// Check against current userRole
	currentUserRole, err := userRoleRepo.GetFromOrganization(id, orgid)
	if err != nil {
		app.ShowError(w, app.ErrEntityUpdate, err, http.StatusUnauthorized)
		return
	}
	// Avoid ID spoofing
	err = verifyID(userRole.IdentifiableModel, currentUserRole.IdentifiableModel)
	if err != nil {
		app.ShowError(w, app.ErrEntityUpdate, err, http.StatusUnauthorized)
		return
	}
	// Set values
	genUserRoleName(userRole)
	// Update
	err = userRoleRepo.Update(userRole)
	if err != nil {
		app.ShowError(w, app.ErrEntityUpdate, err, http.StatusInternalServerError)
		return
	}
	// Marshal
	j, err := json.Marshal(UserRoleResource{Data: *userRole})
	if err != nil {
		app.ShowError(w, app.ErrResponseMarshalling, err, http.StatusNoContent)
		return
	}
	// Respond
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
	w.Write(j)
}

// DeleteUserRole - Deletes an existing UserRole
// Handler for HTTP Delete - "/organizations/{organization}/user-roles/{id}"
func DeleteUserRole(w http.ResponseWriter, r *http.Request) {
	// Get IDs
	vars := mux.Vars(r)
	orgid := vars["organization"]
	id := vars["user-role"]
	// Get repo
	userRoleRepo, err := repo.MakeUserRoleRepository()
	if err != nil {
		app.ShowError(w, app.ErrEntitySelect, err, http.StatusInternalServerError)
		return
	}
	// Delete
	err = userRoleRepo.DeleteFromOrganization(id, orgid)
	if err != nil {
		app.ShowError(w, app.ErrEntityDelete, err, http.StatusInternalServerError)
		return
	}
	// Respond
	w.WriteHeader(http.StatusNoContent)
}

// GenName - Generate a name for the UserRole.
func genUserRoleNameAndDescription(userRole *models.UserRole) error {
	err := genUserRoleName(userRole)
	if err != nil {
		return err

	}
	genUserRoleDescription(userRole)
	return nil
}

// genName - Generate a name for the UserRole.
func genUserRoleName(userRole *models.UserRole) error {
	org, _ := getOrganization(userRole.OrganizationID.String)
	user, _ := getUser(userRole.UserID.String)
	role, _ := getRole(userRole.RoleID.String)
	if org.Name.String != "" && user.Username.String != "" && role.Name.String != "" {
		name := fmt.Sprintf("%s::%s::%s", org.Name.String, user.Username.String, role.Name.String)
		userRole.Name = models.ToNullsString(name)
		return nil
	}
	logger.Debug("Que tal")
	return app.ErrEntitySetProperty
}

// genDescription - Generate a name for the UserRole.
func genUserRoleDescription(rp *models.UserRole) error {
	if rp.Name.String != "" {
		rp.Description = models.ToNullsString(fmt.Sprintf("[%s description]", rp.Name.String))
		return nil
	}
	return app.ErrEntitySetProperty
}

func userRoleIDfromURL(r *http.Request) string {
	u, _ := url.Parse(r.URL.Path)
	dir := path.Dir(u.Path)
	id := path.Base(dir)
	logger.Debugf("UserRole id in url is %s", id)
	return id
}

func userRoleNameFromURL(r *http.Request) string {
	u, _ := url.Parse(r.URL.Path)
	dir := path.Dir(u.Path)
	userRoleName := path.Base(dir)
	logger.Debugf("UserRoleName in url is %s", userRoleName)
	return userRoleName
}
