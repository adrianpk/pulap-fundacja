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

	"github.com/adrianpk/pulap/app"
	"github.com/adrianpk/pulap/models"
	"net/http"

	_ "github.com/lib/pq" // Import pq without side effects

	"github.com/adrianpk/pulap/repo"
)

// GetPropertiesSets - Returns a collection containing all properties from some owner.
// Handler for HTTP Get - "/properties-set"
func GetPropertiesSets(w http.ResponseWriter, r *http.Request) {
	// Get ID
	vars := mux.Vars(r)
	holderID := vars["holder"]
	// Get repo
	propSetRepo, err := repo.MakePropertiesSetRepository()
	if err != nil {
		app.ShowError(w, app.ErrEntityNotFound, err, http.StatusInternalServerError)
		return
	}
	// Select
	propSets, err := propSetRepo.GetAll(holderID)
	if err != nil {
		app.ShowError(w, app.ErrEntitySelect, err, http.StatusInternalServerError)
		return
	}
	// Marshal
	j, err := json.Marshal(PropertiesSetsResource{Data: propSets})
	if err != nil {
		app.ShowError(w, app.ErrResponseMarshalling, err, http.StatusInternalServerError)
		return
	}
	// Respond
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

// CreatePropertiesSet - Creates a new Resource.
// Handler for HTTP Post - "/holders/{holder}/properties-sets/create"
func CreatePropertiesSet(w http.ResponseWriter, r *http.Request) {
	// Decode
	var res PropertiesSetResource
	err := json.NewDecoder(r.Body).Decode(&res)
	if err != nil {
		app.ShowError(w, app.ErrRequestParsing, err, http.StatusInternalServerError)
		return
	}
	propertiesSet := &res.Data
	// Get repo
	propertiesSetRepo, err := repo.MakePropertiesSetRepository()
	if err != nil {
		app.ShowError(w, app.ErrEntityCreate, err, http.StatusInternalServerError)
		return
	}
	// Set values
	propertiesSet.SetID()
	// Persist
	propertiesSetRepo.Create(propertiesSet)
	if err != nil {
		app.ShowError(w, app.ErrEntityCreate, err, http.StatusInternalServerError)
		return
	}
	// Marshal
	j, err := json.Marshal(PropertiesSetResource{Data: *propertiesSet})
	if err != nil {
		app.ShowError(w, app.ErrResponseMarshalling, err, http.StatusNoContent)
		return
	}
	// Respond
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(j)
}

// GetPropertiesSet - Returns a single PropertiesSet by its id or propSetName.
// Handler for HTTP Get - "/holders/{holder}/properties-sets/{properties-sets+}"
func GetPropertiesSet(w http.ResponseWriter, r *http.Request) {
	// Get ID
	vars := mux.Vars(r)
	key := vars["holder"]
	if len(key) == 36 {
		GetPropertiesSetByID(w, r)
	} else {
		GetPropertiesSetByName(w, r)
	}
}

// GetPropertiesSetByID - Returns a single Resource by its id.
// Handler for HTTP Get - "/holders/{holder}/properties-sets/{properties-sets+}"
func GetPropertiesSetByID(w http.ResponseWriter, r *http.Request) {
	// Get ID
	vars := mux.Vars(r)
	id := vars["holder"]
	// Get repo
	propertiesSetRepo, err := repo.MakePropertiesSetRepository()
	if err != nil {
		app.ShowError(w, app.ErrEntitySelect, err, http.StatusInternalServerError)
		return
	}
	// Select
	propertiesSet, err := propertiesSetRepo.Get(id)
	if err != nil {
		app.ShowError(w, app.ErrEntitySelect, err, http.StatusInternalServerError)
		return
	}
	// Marshal
	j, err := json.Marshal(propertiesSet)
	if err != nil {
		app.ShowError(w, app.ErrResponseMarshalling, err, http.StatusInternalServerError)
		return
	}
	// Repsond
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

// GetPropertiesSetByName - Returns a single Resource by its propSetName.
// Handler for HTTP Get - "/holders/{holder}/properties-sets/{properties-sets+}"
func GetPropertiesSetByName(w http.ResponseWriter, r *http.Request) {
	// Get ID
	vars := mux.Vars(r)
	propSetName := vars["holder"]
	// Get repo
	propertiesSetRepo, err := repo.MakePropertiesSetRepository()
	if err != nil {
		app.ShowError(w, app.ErrEntitySelect, err, http.StatusInternalServerError)
		return
	}
	// Select
	propertiesSet, err := propertiesSetRepo.GetByName(propSetName)
	if err != nil {
		app.ShowError(w, app.ErrEntitySelect, err, http.StatusInternalServerError)
		return
	}
	// Marshal
	j, err := json.Marshal(propertiesSet)
	if err != nil {
		app.ShowError(w, app.ErrResponseMarshalling, err, http.StatusInternalServerError)
		return
	}
	// Repond
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

// UpdatePropertiesSet - Update an existing Resource.
// Handler for HTTP Put - "/holders/{holder}/properties-sets/{properties-sets+}"
func UpdatePropertiesSet(w http.ResponseWriter, r *http.Request) {
	// Get ID
	vars := mux.Vars(r)
	id := vars["properties-set"]
	// Decode
	var res PropertiesSetResource
	err := json.NewDecoder(r.Body).Decode(&res)
	if err != nil {
		app.ShowError(w, app.ErrRequestParsing, err, http.StatusInternalServerError)
		return
	}
	propertiesSet := &res.Data
	propertiesSet.ID = models.ToNullsString(id)
	// Get repo
	propertiesSetRepo, err := repo.MakePropertiesSetRepository()
	if err != nil {
		app.ShowError(w, app.ErrEntityUpdate, err, http.StatusInternalServerError)
		return
	}
	// Check against current resource
	currentResource, err := propertiesSetRepo.Get(id)
	if err != nil {
		app.ShowError(w, app.ErrEntityUpdate, err, http.StatusUnauthorized)
		return
	}
	// Avoid ID spoofing
	err = verifyID(propertiesSet.IdentifiableModel, currentResource.IdentifiableModel)
	if err != nil {
		app.ShowError(w, app.ErrEntityUpdate, err, http.StatusUnauthorized)
		return
	}
	// Update
	err = propertiesSetRepo.Update(propertiesSet)
	if err != nil {
		app.ShowError(w, app.ErrEntityUpdate, err, http.StatusInternalServerError)
		return
	}
	// Marshal
	j, err := json.Marshal(PropertiesSetResource{Data: *propertiesSet})
	if err != nil {
		app.ShowError(w, app.ErrResponseMarshalling, err, http.StatusNoContent)
		return
	}
	// Respond
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
	w.Write(j)
}

// DeletePropertiesSet - Deletes an existing Resource
// Handler for HTTP Delete - "/holders/{holder}/properties-sets/{properties-set}"
func DeletePropertiesSet(w http.ResponseWriter, r *http.Request) {
	// Get ID
	vars := mux.Vars(r)
	id := vars["properties-set"]
	// Get repo
	propertiesSetRepo, err := repo.MakePropertiesSetRepository()
	if err != nil {
		app.ShowError(w, app.ErrEntitySelect, err, http.StatusInternalServerError)
		return
	}
	// Delete
	err = propertiesSetRepo.Delete(id)
	if err != nil {
		app.ShowError(w, app.ErrEntityDelete, err, http.StatusInternalServerError)
		return
	}
	// Respond
	w.WriteHeader(http.StatusNoContent)
}
