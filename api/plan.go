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

// GetPlans - Returns a collection containing all plans.
// Handler for HTTP Get - "/plans"
func GetPlans(w http.ResponseWriter, r *http.Request) {
	// Get repo
	planRepo, err := repo.MakePlanRepository()
	if err != nil {
		app.ShowError(w, app.ErrEntityNotFound, err, http.StatusInternalServerError)
		return
	}
	// Select
	plans, err := planRepo.GetAll()
	if err != nil {
		app.ShowError(w, app.ErrEntitySelect, err, http.StatusInternalServerError)
		return
	}
	// Marshal
	j, err := json.Marshal(PlansResource{Data: plans})
	if err != nil {
		app.ShowError(w, app.ErrResponseMarshalling, err, http.StatusInternalServerError)
		return
	}
	// Respond
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

// CreatePlan - Creates a new Plan.
// Handler for HTTP Post - "/plans/create"
func CreatePlan(w http.ResponseWriter, r *http.Request) {
	// Decode
	var res PlanResource
	err := json.NewDecoder(r.Body).Decode(&res)
	if err != nil {
		app.ShowError(w, app.ErrRequestParsing, err, http.StatusInternalServerError)
		return
	}
	plan := &res.Data
	// Get repo
	planRepo, err := repo.MakePlanRepository()
	if err != nil {
		app.ShowError(w, app.ErrEntityCreate, err, http.StatusInternalServerError)
		return
	}
	// Persist
	planRepo.Create(plan)
	if err != nil {
		app.ShowError(w, app.ErrEntityCreate, err, http.StatusInternalServerError)
		return
	}
	// Marshal
	j, err := json.Marshal(PlanResource{Data: *plan})
	if err != nil {
		app.ShowError(w, app.ErrResponseMarshalling, err, http.StatusNoContent)
		return
	}
	// Respond
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(j)
}

// GetPlan - Returns a single Plan by its id or planName.
// Handler for HTTP Get - "/plans/:key"
func GetPlan(w http.ResponseWriter, r *http.Request) {
	// Get ID
	vars := mux.Vars(r)
	key := vars["plan"]
	if len(key) == 36 {
		GetPlanByID(w, r)
	} else {
		GetPlanByName(w, r)
	}
}

// GetPlanByID - Returns a single Plan by its id.
// Handler for HTTP Get - "/plans/:key"
func GetPlanByID(w http.ResponseWriter, r *http.Request) {
	// Get ID
	vars := mux.Vars(r)
	id := vars["plan"]
	// Get repo
	planRepo, err := repo.MakePlanRepository()
	if err != nil {
		app.ShowError(w, app.ErrEntitySelect, err, http.StatusInternalServerError)
		return
	}
	// Select
	plan, err := planRepo.Get(id)
	if err != nil {
		app.ShowError(w, app.ErrEntitySelect, err, http.StatusInternalServerError)
		return
	}
	// Marshal
	j, err := json.Marshal(plan)
	if err != nil {
		app.ShowError(w, app.ErrResponseMarshalling, err, http.StatusInternalServerError)
		return
	}
	// Repsond
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

// GetPlanByName - Returns a single Plan by its planName.
// Handler for HTTP Get - "/plans/:key"
func GetPlanByName(w http.ResponseWriter, r *http.Request) {
	// Get ID
	vars := mux.Vars(r)
	planName := vars["plan"]
	// Get repo
	planRepo, err := repo.MakePlanRepository()
	if err != nil {
		app.ShowError(w, app.ErrEntitySelect, err, http.StatusInternalServerError)
		return
	}
	// Select
	plan, err := planRepo.GetByName(planName)
	if err != nil {
		app.ShowError(w, app.ErrEntitySelect, err, http.StatusInternalServerError)
		return
	}
	// Marshal
	j, err := json.Marshal(plan)
	if err != nil {
		app.ShowError(w, app.ErrResponseMarshalling, err, http.StatusInternalServerError)
		return
	}
	// Repond
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

// UpdatePlan - Update an existing Plan.
// Handler for HTTP Put - "/plans/{plan}"
func UpdatePlan(w http.ResponseWriter, r *http.Request) {
	// Get ID
	vars := mux.Vars(r)
	id := vars["plan"]
	// Decode
	var res PlanResource
	err := json.NewDecoder(r.Body).Decode(&res)
	if err != nil {
		app.ShowError(w, app.ErrRequestParsing, err, http.StatusInternalServerError)
		return
	}
	plan := &res.Data
	plan.ID = models.ToNullsString(id)
	// Get repo
	planRepo, err := repo.MakePlanRepository()
	if err != nil {
		app.ShowError(w, app.ErrEntityUpdate, err, http.StatusInternalServerError)
		return
	}
	// Check against current plan
	currentPlan, err := planRepo.Get(id)
	if err != nil {
		app.ShowError(w, app.ErrEntityUpdate, err, http.StatusUnauthorized)
		return
	}
	// Avoid ID spoofing
	err = verifyID(plan.IdentifiableModel, currentPlan.IdentifiableModel)
	if err != nil {
		app.ShowError(w, app.ErrEntityUpdate, err, http.StatusUnauthorized)
		return
	}
	// Update
	err = planRepo.Update(plan)
	if err != nil {
		app.ShowError(w, app.ErrEntityUpdate, err, http.StatusInternalServerError)
		return
	}
	// Marshal
	j, err := json.Marshal(PlanResource{Data: *plan})
	if err != nil {
		app.ShowError(w, app.ErrResponseMarshalling, err, http.StatusNoContent)
		return
	}
	// Respond
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
	w.Write(j)
}

// DeletePlan - Deletes an existing Plan
// Handler for HTTP Delete - "/plans/{plan}"
func DeletePlan(w http.ResponseWriter, r *http.Request) {
	// Get ID
	vars := mux.Vars(r)
	id := vars["plan"]
	// Get repo
	planRepo, err := repo.MakePlanRepository()
	if err != nil {
		app.ShowError(w, app.ErrEntitySelect, err, http.StatusInternalServerError)
		return
	}
	// Delete
	err = planRepo.Delete(id)
	if err != nil {
		app.ShowError(w, app.ErrEntityDelete, err, http.StatusInternalServerError)
		return
	}
	// Respond
	w.WriteHeader(http.StatusNoContent)
}

func planIDfromURL(r *http.Request) string {
	u, _ := url.Parse(r.URL.Path)
	dir := path.Dir(u.Path)
	id := path.Base(dir)
	logger.Debugf("Plan id in url is %s", id)
	return id
}

func planNameFromURL(r *http.Request) string {
	u, _ := url.Parse(r.URL.Path)
	dir := path.Dir(u.Path)
	planName := path.Base(dir)
	logger.Debugf("PlanName in url is %s", planName)
	return planName
}
