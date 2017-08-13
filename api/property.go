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

	"github.com/adrianpk/fundacja/app"
	"github.com/adrianpk/fundacja/logger"
	"github.com/adrianpk/fundacja/models"

	_ "github.com/lib/pq" // Import pq without side effects

	"github.com/adrianpk/fundacja/repo"
)

// GetProperties - Returns a collection containing all properties.
// Handler for HTTP Get - /properties-set/{properties-set}/properties"
func GetProperties(w http.ResponseWriter, r *http.Request) {
	// Get ID
	vars := mux.Vars(r)
	propsetID := vars["properties-set"]
	// Get repo
	propertyRepo, err := repo.MakePropertyRepository()
	if err != nil {
		app.ShowError(w, app.ErrEntityNotFound, err, http.StatusInternalServerError)
		return
	}
	// Select
	properties, err := propertyRepo.GetAll(propsetID)
	if err != nil {
		app.ShowError(w, app.ErrEntitySelect, err, http.StatusInternalServerError)
		return
	}
	// Marshal
	j, err := json.Marshal(PropertiesResource{Data: properties})
	if err != nil {
		app.ShowError(w, app.ErrResponseMarshalling, err, http.StatusInternalServerError)
		return
	}
	// Respond
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

// CreateProperty - Creates a new Property.
// Handler for HTTP Post - /properties-set/{properties-set}/properties/create"
func CreateProperty(w http.ResponseWriter, r *http.Request) {
	// Get PropertiesSet ID
	vars := mux.Vars(r)
	propsetID := vars["properties-set"]
	// Decode
	var res PropertyResource
	err := json.NewDecoder(r.Body).Decode(&res)
	if err != nil {
		app.ShowError(w, app.ErrRequestParsing, err, http.StatusInternalServerError)
		return
	}
	property := &res.Data
	// Set PropertiesSet - Don't trust JSON value
	property.PropertiesSetID = models.ToNullsString(propsetID)
	// Get repo
	propertyRepo, err := repo.MakePropertyRepository()
	if err != nil {
		app.ShowError(w, app.ErrEntityCreate, err, http.StatusInternalServerError)
		return
	}
	// Persist
	propertyRepo.Create(property)
	if err != nil {
		app.ShowError(w, app.ErrEntityCreate, err, http.StatusInternalServerError)
		return
	}
	// Marshal
	j, err := json.Marshal(PropertyResource{Data: *property})
	if err != nil {
		app.ShowError(w, app.ErrResponseMarshalling, err, http.StatusNoContent)
		return
	}
	// Respond
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(j)
}

// GetProperty - Returns a single Property by its id or propertyname.
// Handler for HTTP Get - /properties-set/{properties-set}/properties/:key"
func GetProperty(w http.ResponseWriter, r *http.Request) {
	// Get ID
	vars := mux.Vars(r)
	key := vars["property"]
	if len(key) == 36 {
		GetPropertyByID(w, r)
	} else {
		GetPropertyByName(w, r)
	}
}

// GetPropertyByID - Returns a single Property by its id.
// Handler for HTTP Get - /properties-set/{properties-set}/properties/:key"
func GetPropertyByID(w http.ResponseWriter, r *http.Request) {
	// Get ID
	vars := mux.Vars(r)
	id := vars["property"]
	// Get repo
	propertyRepo, err := repo.MakePropertyRepository()
	if err != nil {
		app.ShowError(w, app.ErrEntitySelect, err, http.StatusInternalServerError)
		return
	}
	// Select
	property, err := propertyRepo.Get(id)
	if err != nil {
		app.ShowError(w, app.ErrEntitySelect, err, http.StatusInternalServerError)
		return
	}
	// Marshal
	j, err := json.Marshal(property)
	if err != nil {
		app.ShowError(w, app.ErrResponseMarshalling, err, http.StatusInternalServerError)
		return
	}
	// Repsond
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

// GetPropertyByName - Returns a single Property by its propertyname.
// Handler for HTTP Get - /properties-set/{properties-set}/properties/:key"
func GetPropertyByName(w http.ResponseWriter, r *http.Request) {
	// Get PropertiesSet ID
	vars := mux.Vars(r)
	propsetID := vars["properties-set"]
	// Get ID
	propertyname := vars["property"]
	// Get repo
	propertyRepo, err := repo.MakePropertyRepository()
	if err != nil {
		app.ShowError(w, app.ErrEntitySelect, err, http.StatusInternalServerError)
		return
	}
	// Select
	property, err := propertyRepo.GetByNameInPropertiesSet(propertyname, propsetID)
	if err != nil {
		app.ShowError(w, app.ErrEntitySelect, err, http.StatusInternalServerError)
		return
	}
	// Marshal
	j, err := json.Marshal(property)
	if err != nil {
		app.ShowError(w, app.ErrResponseMarshalling, err, http.StatusInternalServerError)
		return
	}
	// Repond
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

// UpdateProperty - Update an existing Property.
// Handler for HTTP Put - /properties-set/{properties-set}/properties/:id"
func UpdateProperty(w http.ResponseWriter, r *http.Request) {
	// Get PropertiesSet ID
	vars := mux.Vars(r)
	propsetID := vars["properties-set"]
	// Get ID
	id := vars["property"]
	// Decode
	var res PropertyResource
	err := json.NewDecoder(r.Body).Decode(&res)
	if err != nil {
		app.ShowError(w, app.ErrRequestParsing, err, http.StatusInternalServerError)
		return
	}
	// Set Property ID and PropertiesSetID- Don't trust JSON value
	property := &res.Data
	property.ID = models.ToNullsString(id)
	property.PropertiesSetID = models.ToNullsString(propsetID)
	// Get repo
	propertyRepo, err := repo.MakePropertyRepository()
	if err != nil {
		app.ShowError(w, app.ErrEntityUpdate, err, http.StatusInternalServerError)
		return
	}
	// Update
	err = propertyRepo.Update(property)
	if err != nil {
		app.ShowError(w, app.ErrEntityUpdate, err, http.StatusInternalServerError)
		return
	}
	// Marshal
	j, err := json.Marshal(PropertyResource{Data: *property})
	if err != nil {
		app.ShowError(w, app.ErrResponseMarshalling, err, http.StatusNoContent)
		return
	}
	// Respond
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
	w.Write(j)
}

// DeleteProperty - Deletes an existing Property
// Handler for HTTP Delete - /properties-set/{properties-set}/properties/{property}"
func DeleteProperty(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// Get ID
	id := vars["property"]
	// Get repo
	propertyRepo, err := repo.MakePropertyRepository()
	if err != nil {
		app.ShowError(w, app.ErrEntitySelect, err, http.StatusInternalServerError)
		return
	}
	// Delete
	err = propertyRepo.Delete(id)
	if err != nil {
		app.ShowError(w, app.ErrEntityDelete, err, http.StatusInternalServerError)
		return
	}
	// Respond
	w.WriteHeader(http.StatusNoContent)
}

func propertyIDfromURL(r *http.Request) string {
	u, _ := url.Parse(r.URL.Path)
	dir := path.Dir(u.Path)
	id := path.Base(dir)
	logger.Debugf("Property id in url is %s", id)
	return id
}

func propertynameFromURL(r *http.Request) string {
	u, _ := url.Parse(r.URL.Path)
	dir := path.Dir(u.Path)
	propertyname := path.Base(dir)
	logger.Debugf("Propertyname in url is %s", propertyname)
	return propertyname
}
