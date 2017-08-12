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
	"fmt"

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

// GetRolePermissions - Returns a collection containing all rolePermissions.
// Handler for HTTP Get - "/organizations/{organization}/role-permissions"
func GetRolePermissions(w http.ResponseWriter, r *http.Request) {
	// Get ID
	vars := mux.Vars(r)
	orgid := vars["organization"]
	// Get repo
	rolePermissionRepo, err := repo.MakeRolePermissionRepository()
	if err != nil {
		app.ShowError(w, app.ErrEntityNotFound, err, http.StatusInternalServerError)
		return
	}
	// Select
	rolePermissions, err := rolePermissionRepo.GetAll(orgid)
	if err != nil {
		app.ShowError(w, app.ErrEntitySelect, err, http.StatusInternalServerError)
		return
	}
	// Marshal
	j, err := json.Marshal(RolePermissionsResource{Data: rolePermissions})
	if err != nil {
		app.ShowError(w, app.ErrResponseMarshalling, err, http.StatusInternalServerError)
		return
	}
	// Respond
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

// CreateRolePermission - Creates a new RolePermission.
// Handler for HTTP Post - "/organizations/{organization}/role-permissions/create"
func CreateRolePermission(w http.ResponseWriter, r *http.Request) {
	// Decode
	var res RolePermissionResource
	err := json.NewDecoder(r.Body).Decode(&res)
	if err != nil {
		app.ShowError(w, app.ErrRequestParsing, err, http.StatusInternalServerError)
		return
	}
	rolePermission := &res.Data
	// Get repo
	rolePermissionRepo, err := repo.MakeRolePermissionRepository()
	if err != nil {
		app.ShowError(w, app.ErrEntityCreate, err, http.StatusInternalServerError)
		return
	}
	// Set values
	u, _ := sessionUser(r)
	rolePermission.CreatedBy = u.ID
	genRolePermissionNameAndDescription(rolePermission)
	// Persist
	err = rolePermissionRepo.Create(rolePermission)
	if err != nil {
		app.ShowError(w, app.ErrEntityCreate, err, http.StatusInternalServerError)
		return
	}
	// Marshal
	j, err := json.Marshal(RolePermissionResource{Data: *rolePermission})
	if err != nil {
		app.ShowError(w, app.ErrResponseMarshalling, err, http.StatusNoContent)
		return
	}
	// Respond
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(j)
}

// GetRolePermission - Returns a single RolePermission by its id or rolePermissionname.
// Handler for HTTP Get - "/organizations/{organization}/role-permissions/{role-permissions}"
func GetRolePermission(w http.ResponseWriter, r *http.Request) {
	GetRolePermissionByID(w, r)
}

// GetRolePermissionByID - Returns a single RolePermission by its id.
// Handler for HTTP Get - "/organizations/{organization}/role-permissions/{role-permissions}"
func GetRolePermissionByID(w http.ResponseWriter, r *http.Request) {
	// Get IDs
	vars := mux.Vars(r)
	orgid := vars["organization"]
	id := vars["role-permission"]
	// Get repo
	rolePermissionRepo, err := repo.MakeRolePermissionRepository()
	if err != nil {
		app.ShowError(w, app.ErrEntitySelect, err, http.StatusInternalServerError)
		return
	}
	// Select
	rolePermission, err := rolePermissionRepo.GetFromOrganization(id, orgid)
	if err != nil {
		app.ShowError(w, app.ErrEntitySelect, err, http.StatusInternalServerError)
		return
	}
	// Marshal
	j, err := json.Marshal(rolePermission)
	if err != nil {
		app.ShowError(w, app.ErrResponseMarshalling, err, http.StatusInternalServerError)
		return
	}
	// Repsond
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

// UpdateRolePermission - Update an existing RolePermission.
// Handler for HTTP Put - "/organizations/{organization}/role-permissions/{role-permissions}"
func UpdateRolePermission(w http.ResponseWriter, r *http.Request) {
	// Get IDs
	vars := mux.Vars(r)
	orgid := vars["organization"]
	id := vars["role-permission"]
	// Decode
	var res RolePermissionResource
	err := json.NewDecoder(r.Body).Decode(&res)
	if err != nil {
		app.ShowError(w, app.ErrRequestParsing, err, http.StatusInternalServerError)
		return
	}
	rolePermission := &res.Data
	rolePermission.ID = models.ToNullsString(id)
	// Get repo
	rolePermissionRepo, err := repo.MakeRolePermissionRepository()
	if err != nil {
		app.ShowError(w, app.ErrEntityUpdate, err, http.StatusInternalServerError)
		return
	}
	// Check against current rolePermission
	currentRolePermission, err := rolePermissionRepo.GetFromOrganization(id, orgid)
	if err != nil {
		app.ShowError(w, app.ErrEntityUpdate, err, http.StatusUnauthorized)
		return
	}
	// Avoid ID spoofing
	err = verifyID(rolePermission.IdentifiableModel, currentRolePermission.IdentifiableModel)
	if err != nil {
		app.ShowError(w, app.ErrEntityUpdate, err, http.StatusUnauthorized)
		return
	}
	// Set values
	genRolePermissionName(rolePermission)
	// Update
	err = rolePermissionRepo.Update(rolePermission)
	if err != nil {
		app.ShowError(w, app.ErrEntityUpdate, err, http.StatusInternalServerError)
		return
	}
	// Marshal
	j, err := json.Marshal(RolePermissionResource{Data: *rolePermission})
	if err != nil {
		app.ShowError(w, app.ErrResponseMarshalling, err, http.StatusNoContent)
		return
	}
	// Respond
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
	w.Write(j)
}

// DeleteRolePermission - Deletes an existing RolePermission
// Handler for HTTP Delete - "/organizations/{organization}/role-permissions/{id}"
func DeleteRolePermission(w http.ResponseWriter, r *http.Request) {
	// Get IDs
	vars := mux.Vars(r)
	orgid := vars["organization"]
	id := vars["role-permission"]
	// Get repo
	rolePermissionRepo, err := repo.MakeRolePermissionRepository()
	if err != nil {
		app.ShowError(w, app.ErrEntitySelect, err, http.StatusInternalServerError)
		return
	}
	// Delete
	err = rolePermissionRepo.DeleteFromOrganization(id, orgid)
	if err != nil {
		app.ShowError(w, app.ErrEntityDelete, err, http.StatusInternalServerError)
		return
	}
	// Respond
	w.WriteHeader(http.StatusNoContent)
}

// GenName - Generate a name for the RolePermission.
func genRolePermissionNameAndDescription(rolePermission *models.RolePermission) error {
	err := genRolePermissionName(rolePermission)
	if err != nil {
		return err

	}
	genRolePermissionDescription(rolePermission)
	return nil
}

// genName - Generate a name for the RolePermission.
func genRolePermissionName(rolePermission *models.RolePermission) error {
	org, _ := getOrganization(rolePermission.OrganizationID.String)
	role, _ := getRole(rolePermission.RoleID.String)
	perm, _ := getPermission(rolePermission.PermissionID.String)
	if org.Name.String != "" && role.Name.String != "" && perm.Name.String != "" {
		name := fmt.Sprintf("%s::%s::%s", org.Name.String, role.Name.String, perm.Name.String)
		rolePermission.Name = models.ToNullsString(name)
		return nil
	}
	return app.ErrEntitySetProperty
}

// genDescription - Generate a name for the RolePermission.
func genRolePermissionDescription(rp *models.RolePermission) error {
	if rp.Name.String != "" {
		rp.Description = models.ToNullsString(fmt.Sprintf("[%s description]", rp.Name.String))
		return nil
	}
	return app.ErrEntitySetProperty
}

func rolePermissionIDfromURL(r *http.Request) string {
	u, _ := url.Parse(r.URL.Path)
	dir := path.Dir(u.Path)
	id := path.Base(dir)
	logger.Debugf("RolePermission id in url is %s", id)
	return id
}

func rolePermissionnameFromURL(r *http.Request) string {
	u, _ := url.Parse(r.URL.Path)
	dir := path.Dir(u.Path)
	rolePermissionname := path.Base(dir)
	logger.Debugf("RolePermissionname in url is %s", rolePermissionname)
	return rolePermissionname
}
