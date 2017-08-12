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
	//".com/jmoiron/sqlx"

	"github.com/adrianpk/fundacja/app"
	"github.com/adrianpk/fundacja/logger"
	"github.com/adrianpk/fundacja/models"
	"net/http"
	"net/url"
	"path"

	_ "github.com/lib/pq" // Import pq without side effects

	"github.com/adrianpk/fundacja/repo"
)

// GetResourcePermissions - Returns a collection containing all resourcePermissions.
// Handler for HTTP Get - "/organizations/{organization}/resource-permissions"
func GetResourcePermissions(w http.ResponseWriter, r *http.Request) {
	// Get ID
	vars := mux.Vars(r)
	orgid := vars["organization"]
	// Get repo
	resourcePermissionRepo, err := repo.MakeResourcePermissionRepository()
	if err != nil {
		app.ShowError(w, app.ErrEntityNotFound, err, http.StatusInternalServerError)
		return
	}
	// Select
	resourcePermissions, err := resourcePermissionRepo.GetAll(orgid)
	if err != nil {
		app.ShowError(w, app.ErrEntitySelect, err, http.StatusInternalServerError)
		return
	}
	// Marshal
	j, err := json.Marshal(ResourcePermissionsResource{Data: resourcePermissions})
	if err != nil {
		app.ShowError(w, app.ErrResponseMarshalling, err, http.StatusInternalServerError)
		return
	}
	// Respond
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

// CreateResourcePermission - Creates a new ResourcePermission.
// Handler for HTTP Post - "/organizations/{organization}/resource-permissions/create"
func CreateResourcePermission(w http.ResponseWriter, r *http.Request) {
	// Decode
	var res ResourcePermissionResource
	err := json.NewDecoder(r.Body).Decode(&res)
	if err != nil {
		app.ShowError(w, app.ErrRequestParsing, err, http.StatusInternalServerError)
		return
	}
	resourcePermission := &res.Data
	// Get repo
	resourcePermissionRepo, err := repo.MakeResourcePermissionRepository()
	if err != nil {
		app.ShowError(w, app.ErrEntityCreate, err, http.StatusInternalServerError)
		return
	}
	// Set values
	u, _ := sessionUser(r)
	resourcePermission.CreatedBy = u.ID
	genResourcePermissionName(resourcePermission)
	// Persist
	resourcePermissionRepo.Create(resourcePermission)
	if err != nil {
		app.ShowError(w, app.ErrEntityCreate, err, http.StatusInternalServerError)
		return
	}
	// Marshal
	j, err := json.Marshal(ResourcePermissionResource{Data: *resourcePermission})
	if err != nil {
		app.ShowError(w, app.ErrResponseMarshalling, err, http.StatusNoContent)
		return
	}
	// Respond
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(j)
}

// GetResourcePermission - Returns a single ResourcePermission by its id or resourcePermissionname.
// Handler for HTTP Get - "/organizations/{organization}/resource-permissions/:key"
func GetResourcePermission(w http.ResponseWriter, r *http.Request) {
	// Get ID
	vars := mux.Vars(r)
	key := vars["resource-permission"]
	if len(key) == 36 {
		GetResourcePermissionByID(w, r)
	} else {
		GetResourcePermissionByName(w, r)
	}
}

// GetResourcePermissionByID - Returns a single ResourcePermission by its id.
// Handler for HTTP Get - "/organizations/{organization}/resource-permissions/:key"
func GetResourcePermissionByID(w http.ResponseWriter, r *http.Request) {
	// Get ID
	vars := mux.Vars(r)
	id := vars["resource-permission"]
	// Get repo
	resourcePermissionRepo, err := repo.MakeResourcePermissionRepository()
	if err != nil {
		app.ShowError(w, app.ErrEntitySelect, err, http.StatusInternalServerError)
		return
	}
	// Select
	resourcePermission, err := resourcePermissionRepo.Get(id)
	if err != nil {
		app.ShowError(w, app.ErrEntitySelect, err, http.StatusInternalServerError)
		return
	}
	// Marshal
	j, err := json.Marshal(resourcePermission)
	if err != nil {
		app.ShowError(w, app.ErrResponseMarshalling, err, http.StatusInternalServerError)
		return
	}
	// Repsond
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

// GetResourcePermissionByName - Returns a single ResourcePermission by its resourcePermissionname.
// Handler for HTTP Get - "/organizations/{organization}/resource-permissions/:key"
func GetResourcePermissionByName(w http.ResponseWriter, r *http.Request) {
	// Get ID
	vars := mux.Vars(r)
	resourcePermissionname := vars["resource-permission"]
	// Get repo
	resourcePermissionRepo, err := repo.MakeResourcePermissionRepository()
	if err != nil {
		app.ShowError(w, app.ErrEntitySelect, err, http.StatusInternalServerError)
		return
	}
	// Select
	resourcePermission, err := resourcePermissionRepo.GetByName(resourcePermissionname)
	if err != nil {
		app.ShowError(w, app.ErrEntitySelect, err, http.StatusInternalServerError)
		return
	}
	// Marshal
	j, err := json.Marshal(resourcePermission)
	if err != nil {
		app.ShowError(w, app.ErrResponseMarshalling, err, http.StatusInternalServerError)
		return
	}
	// Repond
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

// UpdateResourcePermission - Update an existing ResourcePermission.
// Handler for HTTP Put - "/organizations/{organization}/resource-permissions/:id"
func UpdateResourcePermission(w http.ResponseWriter, r *http.Request) {
	// Get ID
	vars := mux.Vars(r)
	id := vars["resource-permission"]
	// Decode
	var res ResourcePermissionResource
	err := json.NewDecoder(r.Body).Decode(&res)
	if err != nil {
		app.ShowError(w, app.ErrRequestParsing, err, http.StatusInternalServerError)
		return
	}
	resourcePermission := &res.Data
	resourcePermission.ID = models.ToNullsString(id)
	// Get repo
	resourcePermissionRepo, err := repo.MakeResourcePermissionRepository()
	if err != nil {
		app.ShowError(w, app.ErrEntityUpdate, err, http.StatusInternalServerError)
		return
	}
	// Check against current resourcePermission
	currentResourcePermission, err := resourcePermissionRepo.Get(id)
	if err != nil {
		app.ShowError(w, app.ErrEntityUpdate, err, http.StatusUnauthorized)
		return
	}
	// Avoid ID spoofing
	err = verifyID(resourcePermission.IdentifiableModel, currentResourcePermission.IdentifiableModel)
	if err != nil {
		app.ShowError(w, app.ErrEntityUpdate, err, http.StatusUnauthorized)
		return
	}
	// Update
	err = resourcePermissionRepo.Update(resourcePermission)
	if err != nil {
		app.ShowError(w, app.ErrEntityUpdate, err, http.StatusInternalServerError)
		return
	}
	// Marshal
	j, err := json.Marshal(ResourcePermissionResource{Data: *resourcePermission})
	if err != nil {
		app.ShowError(w, app.ErrResponseMarshalling, err, http.StatusNoContent)
		return
	}
	// Respond
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
	w.Write(j)
}

// DeleteResourcePermission - Deletes an existing ResourcePermission
// Handler for HTTP Delete - "/organizations/{organization}/resource-permissions/{resource-permission}"
func DeleteResourcePermission(w http.ResponseWriter, r *http.Request) {
	// Get ID
	vars := mux.Vars(r)
	id := vars["resource-permission"]
	// Get repo
	resourcePermissionRepo, err := repo.MakeResourcePermissionRepository()
	if err != nil {
		app.ShowError(w, app.ErrEntitySelect, err, http.StatusInternalServerError)
		return
	}
	// Delete
	err = resourcePermissionRepo.Delete(id)
	if err != nil {
		app.ShowError(w, app.ErrEntityDelete, err, http.StatusInternalServerError)
		return
	}
	// Respond
	w.WriteHeader(http.StatusNoContent)
}

func genResourcePermissionNameAndDescription(resourcePermission *models.ResourcePermission) error {
	err := genResourcePermissionName(resourcePermission)
	if err != nil {
		return err

	}
	genResourcePermissionDescription(resourcePermission)
	return nil
}

func genResourcePermissionName(resourcePermission *models.ResourcePermission) error {
	org, _ := getOrganization(resourcePermission.OrganizationID.String)
	res, _ := getResource(resourcePermission.ResourceID.String)
	perm, _ := getPermission(resourcePermission.PermissionID.String)
	if org.Name.String != "" && res.Name.String != "" && perm.Name.String != "" {
		name := fmt.Sprintf("%s::%s::%s", org.Name.String, res.Name.String, perm.Name.String)
		resourcePermission.Name = models.ToNullsString(name)
		return nil
	}
	return app.ErrEntitySetProperty
}

func genResourcePermissionDescription(resourcePermission *models.ResourcePermission) error {
	if resourcePermission.Name.String != "" {
		resourcePermission.Description = models.ToNullsString(fmt.Sprintf("[%s description]", resourcePermission.Name.String))
		return nil
	}
	return app.ErrEntitySetProperty
}

func resourcePermissionIDfromURL(r *http.Request) string {
	u, _ := url.Parse(r.URL.Path)
	dir := path.Dir(u.Path)
	id := path.Base(dir)
	logger.Debugf("ResourcePermission id in url is %s", id)
	return id
}

func resourcePermissionnameFromURL(r *http.Request) string {
	u, _ := url.Parse(r.URL.Path)
	dir := path.Dir(u.Path)
	resourcePermissionname := path.Base(dir)
	logger.Debugf("ResourcePermissionname in url is %s", resourcePermissionname)
	return resourcePermissionname
}
