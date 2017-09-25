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
	"net/url"
	"path"

	"github.com/adrianpk/pulap/app"
	"github.com/adrianpk/pulap/logger"
	"github.com/adrianpk/pulap/models"
	"github.com/adrianpk/pulap/repo"

	"github.com/gorilla/mux"

	_ "github.com/lib/pq" // Import pq without side effects
)

// GetRoles - Returns a collection containing all roles.
// Handler for HTTP Get - "/organizations/{organization}/roles"
func GetRoles(w http.ResponseWriter, r *http.Request) {
	// Get ID
	vars := mux.Vars(r)
	orgid := vars["organization"]
	// Get repo
	roleRepo, err := repo.MakeRoleRepository()
	if err != nil {
		app.ShowError(w, app.ErrEntityNotFound, err, http.StatusInternalServerError)
		return
	}
	// Select
	roles, err := roleRepo.GetAll(orgid)
	if err != nil {
		app.ShowError(w, app.ErrEntitySelect, err, http.StatusInternalServerError)
		return
	}
	// Marshal
	j, err := json.Marshal(RolesResource{Data: roles})
	if err != nil {
		app.ShowError(w, app.ErrResponseMarshalling, err, http.StatusInternalServerError)
		return
	}
	// Respond
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

// CreateRole - Creates a new Role.
// Handler for HTTP Post - "/organizations/{organization}/roles/create"
func CreateRole(w http.ResponseWriter, r *http.Request) {
	// Decode
	var res RoleResource
	err := json.NewDecoder(r.Body).Decode(&res)
	if err != nil {
		app.ShowError(w, app.ErrRequestParsing, err, http.StatusInternalServerError)
		return
	}
	role := &res.Data
	// Get repo
	roleRepo, err := repo.MakeRoleRepository()
	if err != nil {
		app.ShowError(w, app.ErrEntityCreate, err, http.StatusInternalServerError)
		return
	}
	// Persist
	roleRepo.Create(role)
	if err != nil {
		app.ShowError(w, app.ErrEntityCreate, err, http.StatusInternalServerError)
		return
	}
	// Marshal
	j, err := json.Marshal(RoleResource{Data: *role})
	if err != nil {
		app.ShowError(w, app.ErrResponseMarshalling, err, http.StatusNoContent)
		return
	}
	// Respond
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(j)
}

// GetRole - Returns a single Role by its id or rolename.
// Handler for HTTP Get - "/organizations/{organization}/roles/{role}"
func GetRole(w http.ResponseWriter, r *http.Request) {
	// Get ID
	vars := mux.Vars(r)
	key := vars["role"]
	if len(key) == 36 {
		GetRoleByID(w, r)
	} else {
		GetRoleByName(w, r)
	}
}

// GetRoleByID - Returns a single Role by its id.
// Handler for HTTP Get - "/organizations/{organization}/roles/{role}"
func GetRoleByID(w http.ResponseWriter, r *http.Request) {
	// Get ID
	vars := mux.Vars(r)
	id := vars["role"]
	// Get repo
	roleRepo, err := repo.MakeRoleRepository()
	if err != nil {
		app.ShowError(w, app.ErrEntitySelect, err, http.StatusInternalServerError)
		return
	}
	// Select
	role, err := roleRepo.Get(id)
	if err != nil {
		app.ShowError(w, app.ErrEntitySelect, err, http.StatusInternalServerError)
		return
	}
	// Marshal
	j, err := json.Marshal(role)
	if err != nil {
		app.ShowError(w, app.ErrResponseMarshalling, err, http.StatusInternalServerError)
		return
	}
	// Repsond
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

// GetRoleByName - Returns a single Role by its rolename.
// Handler for HTTP Get - "/organizations/{organization}/roles/{role}"
func GetRoleByName(w http.ResponseWriter, r *http.Request) {
	// Get ID
	vars := mux.Vars(r)
	rolename := vars["role"]
	// Get repo
	roleRepo, err := repo.MakeRoleRepository()
	if err != nil {
		app.ShowError(w, app.ErrEntitySelect, err, http.StatusInternalServerError)
		return
	}
	// Select
	role, err := roleRepo.GetByName(rolename)
	if err != nil {
		app.ShowError(w, app.ErrEntitySelect, err, http.StatusInternalServerError)
		return
	}
	// Marshal
	j, err := json.Marshal(role)
	if err != nil {
		app.ShowError(w, app.ErrResponseMarshalling, err, http.StatusInternalServerError)
		return
	}
	// Repond
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

// UpdateRole - Update an existing Role.
// Handler for HTTP Put - "/organizations/{organization}/roles/{role}"
func UpdateRole(w http.ResponseWriter, r *http.Request) {
	// Get ID
	vars := mux.Vars(r)
	id := vars["role"]
	// Decode
	var res RoleResource
	err := json.NewDecoder(r.Body).Decode(&res)
	if err != nil {
		app.ShowError(w, app.ErrRequestParsing, err, http.StatusInternalServerError)
		return
	}
	role := &res.Data
	role.ID = models.ToNullsString(id)
	// Get repo
	roleRepo, err := repo.MakeRoleRepository()
	if err != nil {
		app.ShowError(w, app.ErrEntityUpdate, err, http.StatusInternalServerError)
		return
	}
	// Check against current role
	currentRole, err := roleRepo.Get(id)
	if err != nil {
		app.ShowError(w, app.ErrEntityUpdate, err, http.StatusUnauthorized)
		return
	}
	// Avoid ID spoofing
	err = verifyID(role.IdentifiableModel, currentRole.IdentifiableModel)
	if err != nil {
		app.ShowError(w, app.ErrEntityUpdate, err, http.StatusUnauthorized)
		return
	}
	// Update
	err = roleRepo.Update(role)
	if err != nil {
		app.ShowError(w, app.ErrEntityUpdate, err, http.StatusInternalServerError)
		return
	}
	// Marshal
	j, err := json.Marshal(RoleResource{Data: *role})
	if err != nil {
		app.ShowError(w, app.ErrResponseMarshalling, err, http.StatusNoContent)
		return
	}
	// Respond
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
	w.Write(j)
}

// DeleteRole - Deletes an existing Role
// Handler for HTTP Delete - "/organizations/{organization}/roles/{id}"
func DeleteRole(w http.ResponseWriter, r *http.Request) {
	// Get ID
	vars := mux.Vars(r)
	id := vars["role"]
	// Get repo
	roleRepo, err := repo.MakeRoleRepository()
	if err != nil {
		app.ShowError(w, app.ErrEntitySelect, err, http.StatusInternalServerError)
		return
	}
	// Delete
	err = roleRepo.Delete(id)
	if err != nil {
		app.ShowError(w, app.ErrEntityDelete, err, http.StatusInternalServerError)
		return
	}
	// Respond
	w.WriteHeader(http.StatusNoContent)
}

func roleIDfromURL(r *http.Request) string {
	u, _ := url.Parse(r.URL.Path)
	dir := path.Dir(u.Path)
	id := path.Base(dir)
	logger.Debugf("Role id in url is %s", id)
	return id
}

func rolenameFromURL(r *http.Request) string {
	u, _ := url.Parse(r.URL.Path)
	dir := path.Dir(u.Path)
	rolename := path.Base(dir)
	logger.Debugf("RoleName in url is %s", rolename)
	return rolename
}
