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

	"github.com/adrianpk/fundacja/app"
	"github.com/adrianpk/fundacja/logger"
	"github.com/adrianpk/fundacja/models"
	"net/http"
	"net/url"
	"path"

	_ "github.com/lib/pq" // Import pq without side effects

	"github.com/adrianpk/fundacja/repo"
)

// GetPermissions - Returns a collection containing all permissions.
// Handler for HTTP Get - "/organizations/{organization}/permissions"
func GetPermissions(w http.ResponseWriter, r *http.Request) {
	// Check permissions
	// defer func() {
	// 	recover()
	// 	app.ShowError(w, app.ErrUnauthorized, app.ErrUnauthorized, http.StatusUnauthorized)
	// }()
	// services.IsAllowed("f254cfe5", loggedInUserID(r))

	// Get ID
	vars := mux.Vars(r)
	orgid := vars["organization"]
	// Get repo
	permissionRepo, err := repo.MakePermissionRepository()
	if err != nil {
		app.ShowError(w, app.ErrEntityNotFound, err, http.StatusInternalServerError)
		return
	}
	// Select
	permissions, err := permissionRepo.GetAll(orgid)
	if err != nil {
		app.ShowError(w, app.ErrEntitySelect, err, http.StatusInternalServerError)
		return
	}
	// Marshal
	j, err := json.Marshal(PermissionsResource{Data: permissions})
	if err != nil {
		app.ShowError(w, app.ErrResponseMarshalling, err, http.StatusInternalServerError)
		return
	}
	// Respond
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

// CreatePermission - Creates a new Permission.
// Handler for HTTP Post - "/organizations/{organization}/permissions/create"
func CreatePermission(w http.ResponseWriter, r *http.Request) {
	// Decode
	var res PermissionResource
	err := json.NewDecoder(r.Body).Decode(&res)
	if err != nil {
		app.ShowError(w, app.ErrRequestParsing, err, http.StatusInternalServerError)
		return
	}
	permission := &res.Data
	// Get repo
	permissionRepo, err := repo.MakePermissionRepository()
	if err != nil {
		app.ShowError(w, app.ErrEntityCreate, err, http.StatusInternalServerError)
		return
	}
	// Persist
	permissionRepo.Create(permission)
	if err != nil {
		app.ShowError(w, app.ErrEntityCreate, err, http.StatusInternalServerError)
		return
	}
	// Marshal
	j, err := json.Marshal(PermissionResource{Data: *permission})
	if err != nil {
		app.ShowError(w, app.ErrResponseMarshalling, err, http.StatusNoContent)
		return
	}
	// Respond
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(j)
}

// GetPermission - Returns a single Permission by its id or permissionname.
// Handler for HTTP Get - "/organizations/{organization}/permissions/{permission}"
func GetPermission(w http.ResponseWriter, r *http.Request) {
	// Get ID
	vars := mux.Vars(r)
	key := vars["permission"]
	if len(key) == 36 {
		GetPermissionByID(w, r)
	} else {
		GetPermissionByName(w, r)
	}
}

// GetPermissionByID - Returns a single Permission by its id.
// Handler for HTTP Get - "/organizations/{organization}/permissions/{permission}"
func GetPermissionByID(w http.ResponseWriter, r *http.Request) {
	// Get ID
	vars := mux.Vars(r)
	id := vars["permission"]
	// Get repo
	permissionRepo, err := repo.MakePermissionRepository()
	if err != nil {
		app.ShowError(w, app.ErrEntitySelect, err, http.StatusInternalServerError)
		return
	}
	// Select
	permission, err := permissionRepo.Get(id)
	if err != nil {
		app.ShowError(w, app.ErrEntitySelect, err, http.StatusInternalServerError)
		return
	}
	// Marshal
	j, err := json.Marshal(permission)
	if err != nil {
		app.ShowError(w, app.ErrResponseMarshalling, err, http.StatusInternalServerError)
		return
	}
	// Repsond
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

// GetPermissionByName - Returns a single Permission by its permissionname.
// Handler for HTTP Get - "/organizations/{organization}/permissions/{permission}"
func GetPermissionByName(w http.ResponseWriter, r *http.Request) {
	// Get ID
	vars := mux.Vars(r)
	permissionname := vars["permission"]
	// Get repo
	permissionRepo, err := repo.MakePermissionRepository()
	if err != nil {
		app.ShowError(w, app.ErrEntitySelect, err, http.StatusInternalServerError)
		return
	}
	// Select
	permission, err := permissionRepo.GetByName(permissionname)
	if err != nil {
		app.ShowError(w, app.ErrEntitySelect, err, http.StatusInternalServerError)
		return
	}
	// Marshal
	j, err := json.Marshal(permission)
	if err != nil {
		app.ShowError(w, app.ErrResponseMarshalling, err, http.StatusInternalServerError)
		return
	}
	// Repond
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

// UpdatePermission - Update an existing Permission.
// Handler for HTTP Put - "/organizations/{organization}/permissions/{permission}"
func UpdatePermission(w http.ResponseWriter, r *http.Request) {
	// Get ID
	vars := mux.Vars(r)
	id := vars["permission"]
	// Decode
	var res PermissionResource
	err := json.NewDecoder(r.Body).Decode(&res)
	if err != nil {
		app.ShowError(w, app.ErrRequestParsing, err, http.StatusInternalServerError)
		return
	}
	permission := &res.Data
	permission.ID = models.ToNullsString(id)
	// Get repo
	permissionRepo, err := repo.MakePermissionRepository()
	if err != nil {
		app.ShowError(w, app.ErrEntityUpdate, err, http.StatusInternalServerError)
		return
	}
	// Check against current permission
	currentPermission, err := permissionRepo.Get(id)
	if err != nil {
		app.ShowError(w, app.ErrEntityUpdate, err, http.StatusUnauthorized)
		return
	}
	// Avoid ID spoofing
	err = verifyID(permission.IdentifiableModel, currentPermission.IdentifiableModel)
	if err != nil {
		app.ShowError(w, app.ErrEntityUpdate, err, http.StatusUnauthorized)
		return
	}
	// Update
	err = permissionRepo.Update(permission)
	if err != nil {
		app.ShowError(w, app.ErrEntityUpdate, err, http.StatusInternalServerError)
		return
	}
	// Marshal
	j, err := json.Marshal(PermissionResource{Data: *permission})
	if err != nil {
		app.ShowError(w, app.ErrResponseMarshalling, err, http.StatusNoContent)
		return
	}
	// Respond
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
	w.Write(j)
}

// DeletePermission - Deletes an existing Permission
// Handler for HTTP Delete - "/organizations/{organization}/permissions/{id}"
func DeletePermission(w http.ResponseWriter, r *http.Request) {
	// Get ID
	vars := mux.Vars(r)
	id := vars["permission"]
	// Get repo
	permissionRepo, err := repo.MakePermissionRepository()
	if err != nil {
		app.ShowError(w, app.ErrEntitySelect, err, http.StatusInternalServerError)
		return
	}
	// Delete
	err = permissionRepo.Delete(id)
	if err != nil {
		app.ShowError(w, app.ErrEntityDelete, err, http.StatusInternalServerError)
		return
	}
	// Respond
	w.WriteHeader(http.StatusNoContent)
}

func permissionIDfromURL(r *http.Request) string {
	u, _ := url.Parse(r.URL.Path)
	dir := path.Dir(u.Path)
	id := path.Base(dir)
	logger.Debugf("Permission id in url is %s", id)
	return id
}

func permissionnameFromURL(r *http.Request) string {
	u, _ := url.Parse(r.URL.Path)
	dir := path.Dir(u.Path)
	permissionname := path.Base(dir)
	logger.Debugf("Permissionname in url is %s", permissionname)
	return permissionname
}
