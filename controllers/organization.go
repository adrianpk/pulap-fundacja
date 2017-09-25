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
	htmlTemplate "html/template"

	"github.com/arschles/go-bindata-html-template"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"

	"net/http"
	"net/url"
	"path"

	"github.com/adrianpk/pulap/app"
	"github.com/adrianpk/pulap/bootstrap"
	"github.com/adrianpk/pulap/logger"
	"github.com/adrianpk/pulap/models"
	"github.com/adrianpk/pulap/repo"

	_ "github.com/lib/pq" // Import pq without side effects
)

var (
	organizationAssetsBase    string
	organizationExtAssetsBase string
	organizationTemplates     map[string]*template.Template
	organizationExtTemplates  map[string]*htmlTemplate.Template
	organizationIndex         = "/organizations"
	organizationNew           = "/organizations/new"
	organizationEdit          = "/organizations/edit/%s"
	organizationShow          = "/organizations/%s"
	organizationDelete        = "/organizations/delete/%s"
)

// InitializeOrganization - Initialize the controller
func InitializeOrganization() {
	if organizationTemplates == nil {
		organizationTemplates = make(map[string]*template.Template)
	}
	if organizationExtTemplates == nil {
		organizationExtTemplates = make(map[string]*htmlTemplate.Template)
	}
	parseOrganizationAssets()
	parseOrganizationExtAssets()
}

// IndexOrganizations - Returns a collection containing all organizations.
// Handler for HTTP Get - "/organizations"
func IndexOrganizations(w http.ResponseWriter, r *http.Request) {
	logger.Debug("IndexOrganizations...")
	// Get repo
	organizationRepo, err := repo.MakeOrganizationRepository()
	if err != nil {
		showOrganizationError(w, r, indexView, layoutView, nil, app.ErrEntitySelect, warningAlert, err)
		return
	}
	// Select
	organizations, err := organizationRepo.GetAll()
	if err != nil {
		showOrganizationError(w, r, indexView, layoutView, nil, app.ErrEntitySelect, warningAlert, err)
		return
	}
	// Respond
	pageModel := makePage(organizations, nil)
	renderOrganizationTemplate(w, r, indexView, layoutView, pageModel)
}

// NewOrganization - Presents new organization form.
// Handler for HTTP Get - "/organizations/new"
func NewOrganization(w http.ResponseWriter, r *http.Request) {
	logger.Debug("NewOrganization...")
	renderOrganizationTemplate(w, r, newView, layoutView, makePage(nil, nil))
}

// CreateOrganization - Creates a new Organization.
// Handler for HTTP Post - "/organizations/create"
func CreateOrganization(w http.ResponseWriter, r *http.Request) {
	logger.Debug("CreateOrganization...")
	// Parse
	err := r.ParseForm()
	if err != nil {
		showOrganizationError(w, r, newView, layoutView, nil, app.ErrEntityCreate, warningAlert, err)
		return
	}
	// Decode
	var organization models.Organization
	err = schema.NewDecoder().Decode(&organization, r.Form)
	if err != nil {
		showOrganizationError(w, r, newView, layoutView, organization, app.ErrEntityCreate, warningAlert, err)
		return
	}
	// Get User
	user, err := getUserByUsername(organization.UserUsername.String)
	if err != nil {
		showOrganizationError(w, r, newView, layoutView, organization, app.ErrEntityCreate, warningAlert, err)
		return
	}
	organization.UserID = user.ID
	organization.UserUsername = user.Username
	// Get repo
	organizationRepo, err := repo.MakeOrganizationRepository()
	if err != nil {
		showOrganizationError(w, r, newView, layoutView, organization, app.ErrEntityCreate, warningAlert, err)
		return
	}
	// Persist
	organizationRepo.Create(&organization)
	if err != nil {
		showOrganizationError(w, r, newView, layoutView, organization, app.ErrDataAccess, warningAlert, err)
		return
	}
	// Respond
	redirectTo(w, r, organizationIndex, makePageAlert("Organization created", infoAlert))
}

// ShowOrganization - Returns a single Organization by its id or organizationname.
// Handler for HTTP Get - "/organizations/{organization}"
func ShowOrganization(w http.ResponseWriter, r *http.Request) {
	logger.Debug("ShowOrganization...")
	ShowOrganizationByID(w, r)
}

// ShowOrganizationByID - Returns a single Organization by its id.
// Handler for HTTP Get - "/organizations/{organization}"
func ShowOrganizationByID(w http.ResponseWriter, r *http.Request) {
	logger.Debug("ShowOrganizationByID...")
	// Get ID
	vars := mux.Vars(r)
	id := vars["organization"]
	// Get repo
	organizationRepo, err := repo.MakeOrganizationRepository()
	if err != nil {
		redirectTo(w, r, organizationIndex, makePageAlert(app.ErrDataStore.Error(), warningAlert))
		return
	}
	// Select
	organization, err := organizationRepo.Get(id)
	if err != nil {
		redirectTo(w, r, organizationIndex, makePageAlert(app.ErrDataAccess.Error(), warningAlert))
		return
	}
	renderOrganizationTemplate(w, r, showView, layoutView, makePage(organization, nil))
}

// EditOrganization - Presents edit organization form.
// Handler for HTTP Get - "/organizations/edit"
func EditOrganization(w http.ResponseWriter, r *http.Request) {
	logger.Debug("EditOrganization...")
	// Get ID
	vars := mux.Vars(r)
	id := vars["organization"]
	// Get repo
	organizationRepo, err := repo.MakeOrganizationRepository()
	if err != nil {
		redirectTo(w, r, organizationIndex, makePageAlert(app.ErrDataStore.Error(), warningAlert))
		return
	}
	// Select
	organization, err := organizationRepo.Get(id)
	if err != nil {
		redirectTo(w, r, organizationIndex, makePageAlert(app.ErrDataAccess.Error(), warningAlert))
		return
	}
	renderOrganizationTemplate(w, r, editView, layoutView, makePage(organization, nil))
	return
}

// UpdateOrganization - Update an existing Organization.
// Handler for HTTP Put - "/organizations/{organization}"
func UpdateOrganization(w http.ResponseWriter, r *http.Request) {
	logger.Debug("UpdateOrganization...")
	// Get ID
	vars := mux.Vars(r)
	id := vars["organization"]
	// Parse
	err := r.ParseForm()
	if err != nil {
		showOrganizationError(w, r, signupView, layoutView, nil, app.ErrRegistration, warningAlert, err)
		return
	}
	// Decode
	var organization models.Organization
	err = schema.NewDecoder().Decode(&organization, r.Form)
	if err != nil {
		showOrganizationError(w, r, editView, layoutView, organization, app.ErrRegistration, warningAlert, err)
		return
	}
	organization.ID = models.ToNullsString(id)
	// Get repo
	organizationRepo, err := repo.MakeOrganizationRepository()
	if err != nil {
		showOrganizationError(w, r, editView, layoutView, organization, app.ErrEntityUpdate, warningAlert, err)
		return
	}
	// Check against current organization
	currentOrganization, err := organizationRepo.Get(id)
	if err != nil {
		showOrganizationError(w, r, editView, layoutView, organization, app.ErrEntityUpdate, warningAlert, err)
		return
	}
	// Avoid ID spoofing
	err = verifyID(organization.IdentifiableModel, currentOrganization.IdentifiableModel)
	if err != nil {
		showOrganizationError(w, r, editView, layoutView, currentOrganization, app.ErrEntityUpdate, warningAlert, err)
		return
	}
	// Update
	err = organizationRepo.Update(&organization)
	if err != nil {
		showOrganizationError(w, r, editView, layoutView, currentOrganization, app.ErrEntityUpdate, warningAlert, err)
		return
	}
	// Respond
	redirectTo(w, r, organizationIndex, makePageAlert("Organization updated", warningAlert))
}

// InitDeleteOrganization - Show organization deletion page.
// Handler for HTTP Get - "/organizations/init-delete/{organization}"
func InitDeleteOrganization(w http.ResponseWriter, r *http.Request) {
	logger.Debug("InitDeleteOrganization...")
	// Get ID
	vars := mux.Vars(r)
	id := vars["organization"]
	// Get repo
	organizationRepo, err := repo.MakeOrganizationRepository()
	if err != nil {
		redirectTo(w, r, organizationIndex, makePageAlert(app.ErrDataStore.Error(), warningAlert))
		return
	}
	// Select
	organization, err := organizationRepo.Get(id)
	if err != nil {
		redirectTo(w, r, organizationIndex, makePageAlert(app.ErrDataAccess.Error(), warningAlert))
		return
	}
	renderOrganizationTemplate(w, r, deleteView, layoutView, makePage(organization, nil))
}

// DeleteOrganization - Deletes an existing Organization
// Handler for HTTP Delete - "/organizations/{id}"
func DeleteOrganization(w http.ResponseWriter, r *http.Request) {
	logger.Debug("DeleteOrganization...")
	// Get ID
	vars := mux.Vars(r)
	id := vars["organization"]
	// Get repo
	organizationRepo, err := repo.MakeOrganizationRepository()
	if err != nil {
		showOrganizationError(w, r, indexView, layoutView, nil, app.ErrEntitySelect, warningAlert, err)
		return
	}
	// Delete
	err = organizationRepo.Delete(id)
	if err != nil {
		showOrganizationError(w, r, indexView, layoutView, nil, app.ErrEntityDelete, warningAlert, err)
		return
	}
	// Respond
	redirectTo(w, r, organizationIndex, makePageAlert("Organization updated", warningAlert))
}

func organizationIDfromURL(r *http.Request) string {
	u, _ := url.Parse(r.URL.Path)
	dir := path.Dir(u.Path)
	id := path.Base(dir)
	logger.Debugf("Organization id in url is %s", id)
	return id
}

func organizationnameFromURL(r *http.Request) string {
	u, _ := url.Parse(r.URL.Path)
	dir := path.Dir(u.Path)
	organizationname := path.Base(dir)
	logger.Debugf("OrganizationName in url is %s", organizationname)
	return organizationname
}

func parseOrganizationAssets() {
	//logger.Debug("Parsing organization assets...")
	assetNames := []string{indexView, newView, showView, editView, deleteView}
	parseAssets(&organizationAssetsBase, "layouts", "organization", layoutView, assetNames, organizationTemplates)
}

func parseOrganizationExtAssets() {
	assetNames := []string{indexView, newView, showView, editView, deleteView}
	parseExtAssets(&organizationExtAssetsBase, "layouts", "organization", layoutView, assetNames, organizationExtTemplates)
}

func renderOrganizationTemplate(w http.ResponseWriter, r *http.Request, groupName string, name string, page *Page) {
	if useExtTemplates {
		renderExtOrganizationTemplate(w, r, groupName, name, page)
		return
	}
	renderIntOrganizationTemplate(w, r, groupName, name, page)
}

// Render templates for the given name, template definition and data object
func renderIntOrganizationTemplate(w http.ResponseWriter, r *http.Request, groupName string, name string, page *Page) {
	renderTemplate(w, r, organizationAssetsBase, groupName, name, organizationTemplates, page)
}

// Render templates for the given name, template definition and data object
func renderExtOrganizationTemplate(w http.ResponseWriter, r *http.Request, groupName string, name string, page *Page) {
	//logger.Debugf("Autoreload: %t", bootstrap.AppConfig.IsAutoreloadOn())
	if bootstrap.AppConfig.IsAutoreloadOn() {
		logger.Debug("Reloading templates.")
		organizationExtTemplates = make(map[string]*htmlTemplate.Template)
		parseOrganizationExtAssets()
	}
	renderExtTemplate(w, r, organizationExtAssetsBase, groupName, name, organizationExtTemplates, page)
}

func showOrganizationError(w http.ResponseWriter, r *http.Request, page string, layout string, model interface{}, err error, alertKind string, cause error) {
	logger.Dump(cause)
	pageModel := makePage(model, makePageAlert(err.Error(), alertKind))
	renderOrganizationTemplate(w, r, page, layoutView, pageModel)
}
