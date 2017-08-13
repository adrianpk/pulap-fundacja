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
	//"database/sql"
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

// GetResources - Returns a collection containing all resources.
// Handler for HTTP Get - "/organizations/{organization}/resources"
func GetResources(w http.ResponseWriter, r *http.Request) {
	// Get ID
	vars := mux.Vars(r)
	orgid := vars["organization"]
	// Get repo
	resourceRepo, err := repo.MakeResourceRepository()
	if err != nil {
		app.ShowError(w, app.ErrEntityNotFound, err, http.StatusInternalServerError)
		return
	}
	// Select
	resources, err := resourceRepo.GetAll(orgid)
	if err != nil {
		app.ShowError(w, app.ErrEntitySelect, err, http.StatusInternalServerError)
		return
	}
	// Marshal
	j, err := json.Marshal(ResourcesResource{Data: resources})
	if err != nil {
		app.ShowError(w, app.ErrResponseMarshalling, err, http.StatusInternalServerError)
		return
	}
	// Respond
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

// CreateResource - Creates a new Resource.
// Handler for HTTP Post - "/organizations/{organization}/resources/create"
func CreateResource(w http.ResponseWriter, r *http.Request) {
	// Decode
	var res ResourceResource
	err := json.NewDecoder(r.Body).Decode(&res)
	if err != nil {
		app.ShowError(w, app.ErrRequestParsing, err, http.StatusInternalServerError)
		return
	}
	resource := &res.Data
	// Get repo
	resourceRepo, err := repo.MakeResourceRepository()
	if err != nil {
		app.ShowError(w, app.ErrEntityCreate, err, http.StatusInternalServerError)
		return
	}
	// Set values
	resource.SetID()
	resource.GenTag()
	// Persist
	resourceRepo.Create(resource)
	if err != nil {
		app.ShowError(w, app.ErrEntityCreate, err, http.StatusInternalServerError)
		return
	}
	// Marshal
	j, err := json.Marshal(ResourceResource{Data: *resource})
	if err != nil {
		app.ShowError(w, app.ErrResponseMarshalling, err, http.StatusNoContent)
		return
	}
	// Respond
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(j)
}

// GetResource - Returns a single Resource by its id or resourcename.
// Handler for HTTP Get - "/organizations/{organization}/resources/{resource}"
func GetResource(w http.ResponseWriter, r *http.Request) {
	// Get ID
	vars := mux.Vars(r)
	key := vars["resource"]
	if len(key) == 36 {
		GetResourceByID(w, r)
	} else {
		GetResourceByName(w, r)
	}
}

// GetResourceByID - Returns a single Resource by its id.
// Handler for HTTP Get - "/organizations/{organization}/resources/{resource}"
func GetResourceByID(w http.ResponseWriter, r *http.Request) {
	// Get ID
	vars := mux.Vars(r)
	id := vars["resource"]
	// Get repo
	resourceRepo, err := repo.MakeResourceRepository()
	if err != nil {
		app.ShowError(w, app.ErrEntitySelect, err, http.StatusInternalServerError)
		return
	}
	// Select
	resource, err := resourceRepo.Get(id)
	if err != nil {
		app.ShowError(w, app.ErrEntitySelect, err, http.StatusInternalServerError)
		return
	}
	// Marshal
	j, err := json.Marshal(resource)
	if err != nil {
		app.ShowError(w, app.ErrResponseMarshalling, err, http.StatusInternalServerError)
		return
	}
	// Repsond
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

// GetResourceByName - Returns a single Resource by its resourcename.
// Handler for HTTP Get - "/organizations/{organization}/resources/{resource}"
func GetResourceByName(w http.ResponseWriter, r *http.Request) {
	// Get ID
	vars := mux.Vars(r)
	resourcename := vars["resource"]
	// Get repo
	resourceRepo, err := repo.MakeResourceRepository()
	if err != nil {
		app.ShowError(w, app.ErrEntitySelect, err, http.StatusInternalServerError)
		return
	}
	// Select
	resource, err := resourceRepo.GetByName(resourcename)
	if err != nil {
		app.ShowError(w, app.ErrEntitySelect, err, http.StatusInternalServerError)
		return
	}
	// Marshal
	j, err := json.Marshal(resource)
	if err != nil {
		app.ShowError(w, app.ErrResponseMarshalling, err, http.StatusInternalServerError)
		return
	}
	// Repond
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

// UpdateResource - Update an existing Resource.
// Handler for HTTP Put - "/organizations/{organization}/resources/{resource}"
func UpdateResource(w http.ResponseWriter, r *http.Request) {
	// Get ID
	vars := mux.Vars(r)
	id := vars["resource"]
	// Decode
	var res ResourceResource
	err := json.NewDecoder(r.Body).Decode(&res)
	if err != nil {
		app.ShowError(w, app.ErrRequestParsing, err, http.StatusInternalServerError)
		return
	}
	resource := &res.Data
	resource.ID = models.ToNullsString(id)
	// Get repo
	resourceRepo, err := repo.MakeResourceRepository()
	if err != nil {
		app.ShowError(w, app.ErrEntityUpdate, err, http.StatusInternalServerError)
		return
	}
	// Check against current resource
	currentResource, err := resourceRepo.Get(id)
	if err != nil {

		app.ShowError(w, app.ErrEntityUpdate, err, http.StatusUnauthorized)
		return
	}
	// Avoid ID spoofing
	err = verifyID(resource.IdentifiableModel, currentResource.IdentifiableModel)
	if err != nil {

		app.ShowError(w, app.ErrEntityUpdate, err, http.StatusUnauthorized)
		return
	}
	// Update
	err = resourceRepo.Update(resource)
	if err != nil {

		app.ShowError(w, app.ErrEntityUpdate, err, http.StatusInternalServerError)
		return
	}
	// Marshal
	j, err := json.Marshal(ResourceResource{Data: *resource})
	if err != nil {

		app.ShowError(w, app.ErrResponseMarshalling, err, http.StatusNoContent)
		return
	}
	// Respond
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
	w.Write(j)
}

// DeleteResource - Deletes an existing Resource
// Handler for HTTP Delete - "/organizations/{organization}/resources/{id}"
func DeleteResource(w http.ResponseWriter, r *http.Request) {
	// Get ID
	vars := mux.Vars(r)
	id := vars["resource"]
	// Get repo
	resourceRepo, err := repo.MakeResourceRepository()
	if err != nil {
		app.ShowError(w, app.ErrEntitySelect, err, http.StatusInternalServerError)
		return
	}
	// Delete
	err = resourceRepo.Delete(id)
	if err != nil {
		app.ShowError(w, app.ErrEntityDelete, err, http.StatusInternalServerError)
		return
	}
	// Respond
	w.WriteHeader(http.StatusNoContent)
}

func resourceIDfromURL(r *http.Request) string {
	u, _ := url.Parse(r.URL.Path)
	dir := path.Dir(u.Path)
	id := path.Base(dir)
	logger.Debugf("Resource id in url is %s", id)
	return id
}

func resourceNameFromURL(r *http.Request) string {
	u, _ := url.Parse(r.URL.Path)
	dir := path.Dir(u.Path)
	resourcename := path.Base(dir)
	logger.Debugf("Resourcename in url is %s", resourcename)
	return resourcename
}
