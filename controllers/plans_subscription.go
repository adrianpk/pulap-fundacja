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
	"github.com/adrianpk/fundacja/models"
	"net/http"

	_ "github.com/lib/pq" // Import pq without side effects

	"github.com/adrianpk/fundacja/repo"
)

// CreatePlanSubscription - Creates a new PlanSubscription.
// Handler for HTTP Post - "/resources/create"
func CreatePlanSubscription(w http.ResponseWriter, r *http.Request) {
	// Decode
	var res PlanSubscriptionResource
	err := json.NewDecoder(r.Body).Decode(&res)
	if err != nil {
		app.ShowError(w, app.ErrRequestParsing, err, http.StatusInternalServerError)
		return
	}
	planSubscription := &res.Data
	// Get repo
	planSubscriptionRepo, err := repo.MakePlanSubscriptionRepository()
	if err != nil {
		app.ShowError(w, app.ErrEntityCreate, err, http.StatusInternalServerError)
		return
	}
	// Set values
	planSubscription.SetID()
	// Persist
	planSubscriptionRepo.Create(planSubscription)
	if err != nil {
		app.ShowError(w, app.ErrEntityCreate, err, http.StatusInternalServerError)
		return
	}
	// Marshal
	j, err := json.Marshal(PlanSubscriptionResource{Data: *planSubscription})
	if err != nil {
		app.ShowError(w, app.ErrResponseMarshalling, err, http.StatusNoContent)
		return
	}
	// Respond
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(j)
}

// GetPlanSubscription - Returns a single PlanSubscription by its id or propSetName.
// Handler for HTTP Get - "/resources/:key"
func GetPlanSubscription(w http.ResponseWriter, r *http.Request) {
	// Get ID
	vars := mux.Vars(r)
	key := vars["plan-subscriptor"]
	if len(key) == 36 {
		GetPlanSubscriptionByUserID(w, r)
	} else {
		GetPlanSubscriptionByUsername(w, r)
	}
}

// GetPlanSubscriptionByUserID - Returns a single PlanSubscription associated to Id..
// Handler for HTTP Get - "/resources/:key"
func GetPlanSubscriptionByUserID(w http.ResponseWriter, r *http.Request) {
	// Get ID
	vars := mux.Vars(r)
	id := vars["plan-subscriptor"]
	// Get repo
	planSubscriptionRepo, err := repo.MakePlanSubscriptionRepository()
	if err != nil {
		app.ShowError(w, app.ErrEntitySelect, err, http.StatusInternalServerError)
		return
	}
	// Select
	planSubscription, err := planSubscriptionRepo.GetByUserID(id)
	if err != nil {
		app.ShowError(w, app.ErrEntitySelect, err, http.StatusInternalServerError)
		return
	}
	// Marshal
	j, err := json.Marshal(planSubscription)
	if err != nil {
		app.ShowError(w, app.ErrResponseMarshalling, err, http.StatusInternalServerError)
		return
	}
	// Repsond
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

// GetPlanSubscriptionByUsername - Returns a single PlanSubscription associated to Username.
// Handler for HTTP Get - "/resources/:key"
func GetPlanSubscriptionByUsername(w http.ResponseWriter, r *http.Request) {
	// Get ID
	vars := mux.Vars(r)
	username := vars["plan-subscriptor"]
	// Get User from username
	user, err := getUserByUsername(username)
	if err != nil {
		app.ShowError(w, app.ErrEntitySelect, err, http.StatusInternalServerError)
		return
	}
	// Get repo
	planSubscriptionRepo, err := repo.MakePlanSubscriptionRepository()
	if err != nil {
		app.ShowError(w, app.ErrEntitySelect, err, http.StatusInternalServerError)
		return
	}
	// Select
	planSubscription, err := planSubscriptionRepo.GetByUserID(user.ID.String)
	if err != nil {
		app.ShowError(w, app.ErrEntitySelect, err, http.StatusInternalServerError)
		return
	}
	// Marshal
	j, err := json.Marshal(planSubscription)
	if err != nil {
		app.ShowError(w, app.ErrResponseMarshalling, err, http.StatusInternalServerError)
		return
	}
	// Repond
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

// UpdatePlanSubscription - Update an existing PlanSubscription.
// Handler for HTTP Put - "/resources/:id"
func UpdatePlanSubscription(w http.ResponseWriter, r *http.Request) {
	// Decode
	var res PlanSubscriptionResource
	err := json.NewDecoder(r.Body).Decode(&res)
	if err != nil {
		app.ShowError(w, app.ErrRequestParsing, err, http.StatusInternalServerError)
		return
	}
	planSubscription := &res.Data
	// Verify ownership
	err = verifyOwnership(planSubscription.UserID.String, r)
	if err != nil {
		app.ShowError(w, app.ErrOwnerOnlyCanManage, err, http.StatusUnauthorized)
		return
	}
	// Current User's plan subscription from repo
	userID, _ := sessionUserID(r)
	planSubscriptionRepo, err := repo.MakePlanSubscriptionRepository()
	if err != nil {
		app.ShowError(w, app.ErrEntityCreate, err, http.StatusInternalServerError)
		return
	}
	// Check against current plan subscription
	currentPlanSubscription, err := planSubscriptionRepo.GetByUserID(userID)
	if err != nil {
		app.ShowError(w, app.ErrEntityUpdate, err, http.StatusUnauthorized)
		return
	}
	// Avoid ID spoofing
	err = verifyID(planSubscription.IdentifiableModel, currentPlanSubscription.IdentifiableModel)
	if err != nil {
		app.ShowError(w, app.ErrEntityUpdate, err, http.StatusUnauthorized)
		return
	}
	// Update
	// planSubscription.ValidableDate.Date = &planSubscription.AnniversaryDate
	// planSubscription.ValidableDate.Validate()
	planSubscriptionRepo.Update(planSubscription)
	if err != nil {
		app.ShowError(w, app.ErrEntityCreate, err, http.StatusInternalServerError)
		return
	}
	// Respond
	w.WriteHeader(http.StatusNoContent)
}

// DeletePlanSubscription - Deletes an existing PlanSubscription
// Handler for HTTP Delete - "/plan-subscriptor/{plan-subscriptor}"
func DeletePlanSubscription(w http.ResponseWriter, r *http.Request) {
	// Get ID
	vars := mux.Vars(r)
	id := vars["plan-subscriptor"]
	// Get repo
	planSubscriptionRepo, err := repo.MakePlanSubscriptionRepository()
	if err != nil {
		app.ShowError(w, app.ErrEntitySelect, err, http.StatusInternalServerError)
		return
	}
	// Select
	planSubscription, err := planSubscriptionRepo.GetByUserID(id)
	if err != nil {
		app.ShowError(w, app.ErrEntitySelect, err, http.StatusInternalServerError)
		return
	}
	// Delete
	err = planSubscriptionRepo.Delete(planSubscription.ID.String)
	if err != nil {
		app.ShowError(w, app.ErrEntityDelete, err, http.StatusInternalServerError)
		return
	}
	// Respond
	w.WriteHeader(http.StatusNoContent)
}

func genPlanSubscriptionNameAndDescription(planSubscription *models.PlanSubscription) error {
	err := genPlanSubscriptionName(planSubscription)
	if err != nil {
		return err

	}
	genPlanSubscriptionDescription(planSubscription)
	return nil
}

func genPlanSubscriptionName(planSubscription *models.PlanSubscription) error {
	user, _ := getUser(planSubscription.UserID.String)
	plan, _ := getPlan(planSubscription.PlanID.String)
	if user.Name.String != "" && plan.Name.String != "" {
		name := fmt.Sprintf("%s::%s", user.Name.String, plan.Name.String)
		planSubscription.Name = models.ToNullsString(name)
		return nil
	}
	return app.ErrEntitySetProperty
}

func genPlanSubscriptionDescription(planSubscription *models.PlanSubscription) error {
	if planSubscription.Name.String != "" {
		planSubscription.Description = models.ToNullsString(fmt.Sprintf("[%s description]", planSubscription.Name.String))
		return nil
	}
	return app.ErrEntitySetProperty
}
