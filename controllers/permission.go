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
	"fmt"
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
	permissionAssetsBase    string
	permissionExtAssetsBase string
	permissionTemplates     map[string]*template.Template
	permissionExtTemplates  map[string]*htmlTemplate.Template
	permissionIndex         = "/permissions"
	permissionNew           = "/permissions/new"
	permissionEdit          = "/permissions/edit/%s"
	permissionShow          = "/permissions/%s"
	permissionDelete        = "/permissions/delete/%s"
)

// InitializePermission - Initialize the controller
func InitializePermission() {
	if permissionTemplates == nil {
		permissionTemplates = make(map[string]*template.Template)
	}
	if permissionExtTemplates == nil {
		permissionExtTemplates = make(map[string]*htmlTemplate.Template)
	}
	parsePermissionAssets()
	parsePermissionExtAssets()
}

// IndexPermissions - Returns a collection containing all permissions.
// Handler for HTTP Get - "/permissions"
func IndexPermissions(w http.ResponseWriter, r *http.Request) {
	logger.Debug("IndexPermissions...")
	// Check permissions
	// defer func() {
	// 	recover()
	// 	showPermissionError(w, r, indexView, layoutView, nil, app.ErrEntitySelect, warningAlert, nil)
	// }()
	// services.IsAllowed("f254cfe5", loggedInUserID(r))
	// Get ID
	vars := mux.Vars(r)
	orgid := vars["organization"]
	// Get Organization
	organization, err := getOrganization(orgid)
	if err != nil {
		showPermissionError(w, r, indexView, layoutView, nil, app.ErrEntityNotFound, warningAlert, err)
		return
	}
	permissionRepo, err := repo.MakePermissionRepository()
	if err != nil {
		showPermissionError(w, r, indexView, layoutView, nil, app.ErrEntitySelect, warningAlert, err)
		return
	}
	// Select
	permissions, err := permissionRepo.GetAll(orgid)
	if err != nil {
		showPermissionError(w, r, indexView, layoutView, nil, app.ErrEntitySelect, warningAlert, err)
		return
	}
	pageModel := makePage(makeParentChild(organization, permissions), nil)
	renderPermissionTemplate(w, r, indexView, layoutView, pageModel)
}

// NewPermission - Presents new permission form.
// Handler for HTTP Get - "/permissions/new"
func NewPermission(w http.ResponseWriter, r *http.Request) {
	logger.Debug("NewPermission...")
	vars := mux.Vars(r)
	orgid := vars["organization"]
	// Get Organization
	organization, err := getOrganization(orgid)
	if err != nil {
		showPermissionError(w, r, indexView, layoutView, nil, app.ErrEntityNotFound, warningAlert, err)
		return
	}
	pageModel := makePage(makeParentChild(organization, nil), nil)
	renderPermissionTemplate(w, r, newView, layoutView, pageModel)
}

// CreatePermission - Creates a new Permission.
// Handler for HTTP Post - "/permissions/create"
func CreatePermission(w http.ResponseWriter, r *http.Request) {
	logger.Debug("CreatePermission...")
	vars := mux.Vars(r)
	orgid := vars["organization"]
	// Parse
	err := r.ParseForm()
	if err != nil {
		showPermissionError(w, r, newView, layoutView, nil, app.ErrEntityCreate, warningAlert, err)
		return
	}
	// Decode
	var permission models.Permission
	err = schema.NewDecoder().Decode(&permission, r.Form)
	if err != nil {
		showPermissionError(w, r, newView, layoutView, permission, app.ErrEntityCreate, warningAlert, err)
		return
	}
	// Get Organization
	organization, err := getOrganization(orgid)
	if err != nil {
		showPermissionError(w, r, indexView, layoutView, nil, app.ErrEntityNotFound, warningAlert, err)
		return
	}
	// Get repo
	permissionRepo, err := repo.MakePermissionRepository()
	if err != nil {
		showPermissionError(w, r, newView, layoutView, permission, app.ErrEntityCreate, warningAlert, err)
		return
	}
	// Set values
	permission.OrganizationName = organization.Name
	permission.OrganizationID = organization.ID
	// Persist
	permissionRepo.Create(&permission)
	if err != nil {
		showPermissionError(w, r, newView, layoutView, permission, app.ErrDataAccess, warningAlert, err)
		return
	}
	// Respond
	indexURL := fmt.Sprintf("/organizations/%s/permissions", organization.ID.String)
	redirectTo(w, r, indexURL, makePageAlert("Permission created", infoAlert))
}

// ShowPermission - Returns a single Permission by its id or permissionname.
// Handler for HTTP Get - "/permissions/{permission}"
func ShowPermission(w http.ResponseWriter, r *http.Request) {
	logger.Debug("ShowPermission...")
	ShowPermissionByID(w, r)
}

// ShowPermissionByID - Returns a single Permission by its id.
// Handler for HTTP Get - "/permissions/{permission}"
func ShowPermissionByID(w http.ResponseWriter, r *http.Request) {
	logger.Debug("ShowPermissionByID...")
	// Get ID
	vars := mux.Vars(r)
	id := vars["permission"]
	// Get repo
	permissionRepo, err := repo.MakePermissionRepository()
	if err != nil {
		redirectTo(w, r, permissionIndex, makePageAlert(app.ErrDataStore.Error(), warningAlert))
		return
	}
	// Select
	permission, err := permissionRepo.Get(id)
	if err != nil {
		redirectTo(w, r, permissionIndex, makePageAlert(app.ErrDataAccess.Error(), warningAlert))
		return
	}
	renderPermissionTemplate(w, r, showView, layoutView, makePage(permission, nil))
}

// EditPermission - Presents edit permission form.
// Handler for HTTP Get - "/permissions/edit"
func EditPermission(w http.ResponseWriter, r *http.Request) {
	logger.Debug("EditPermission...")
	// Get ID
	vars := mux.Vars(r)
	id := vars["permission"]
	// Get repo
	permissionRepo, err := repo.MakePermissionRepository()
	if err != nil {
		redirectTo(w, r, permissionIndex, makePageAlert(app.ErrDataStore.Error(), warningAlert))
		return
	}
	// Select
	permission, err := permissionRepo.Get(id)
	if err != nil {
		redirectTo(w, r, permissionIndex, makePageAlert(app.ErrDataAccess.Error(), warningAlert))
		return
	}
	renderPermissionTemplate(w, r, editView, layoutView, makePage(permission, nil))
	return
}

// UpdatePermission - Update an existing Permission.
// Handler for HTTP Put - "/permissions/{permission}"
func UpdatePermission(w http.ResponseWriter, r *http.Request) {
	logger.Debug("UpdatePermission...")
	// Get ID
	vars := mux.Vars(r)
	id := vars["permission"]
	// Parse
	err := r.ParseForm()
	if err != nil {
		showPermissionError(w, r, signupView, layoutView, nil, app.ErrRegistration, warningAlert, err)
		return
	}
	// Decode
	var permission models.Permission
	err = schema.NewDecoder().Decode(&permission, r.Form)
	if err != nil {
		showPermissionError(w, r, editView, layoutView, permission, app.ErrRegistration, warningAlert, err)
		return
	}
	permission.ID = models.ToNullsString(id)
	// Get repo
	permissionRepo, err := repo.MakePermissionRepository()
	if err != nil {
		showPermissionError(w, r, editView, layoutView, permission, app.ErrEntityUpdate, warningAlert, err)
		return
	}
	// Check against current permission
	currentPermission, err := permissionRepo.Get(id)
	if err != nil {
		showPermissionError(w, r, editView, layoutView, permission, app.ErrEntityUpdate, warningAlert, err)
		return
	}
	// Avoid ID spoofing
	err = verifyID(permission.IdentifiableModel, currentPermission.IdentifiableModel)
	if err != nil {
		showPermissionError(w, r, editView, layoutView, currentPermission, app.ErrEntityUpdate, warningAlert, err)
		return
	}
	// Update
	err = permissionRepo.Update(&permission)
	if err != nil {
		showPermissionError(w, r, editView, layoutView, currentPermission, app.ErrEntityUpdate, warningAlert, err)
		return
	}
	// Respond
	redirectTo(w, r, permissionIndex, makePageAlert("Permission updated", warningAlert))
}

// InitDeletePermission - Show permission deletion page.
// Handler for HTTP Get - "/permissions/init-delete/{permission}"
func InitDeletePermission(w http.ResponseWriter, r *http.Request) {
	logger.Debug("InitDeletePermission...")
	// Get ID
	vars := mux.Vars(r)
	id := vars["permission"]
	// Get repo
	permissionRepo, err := repo.MakePermissionRepository()
	if err != nil {
		redirectTo(w, r, permissionIndex, makePageAlert(app.ErrDataStore.Error(), warningAlert))
		return
	}
	// Select
	permission, err := permissionRepo.Get(id)
	if err != nil {
		redirectTo(w, r, permissionIndex, makePageAlert(app.ErrDataAccess.Error(), warningAlert))
		return
	}
	renderPermissionTemplate(w, r, deleteView, layoutView, makePage(permission, nil))
}

// DeletePermission - Deletes an existing Permission
// Handler for HTTP Delete - "/permissions/{id}"
func DeletePermission(w http.ResponseWriter, r *http.Request) {
	logger.Debug("DeletePermission...")
	// Get ID
	vars := mux.Vars(r)
	id := vars["permission"]
	// Get repo
	permissionRepo, err := repo.MakePermissionRepository()
	if err != nil {
		showPermissionError(w, r, indexView, layoutView, nil, app.ErrEntitySelect, warningAlert, err)
		return
	}
	// Delete
	err = permissionRepo.Delete(id)
	if err != nil {
		showPermissionError(w, r, indexView, layoutView, nil, app.ErrEntityDelete, warningAlert, err)
		return
	}
	// Respond
	redirectTo(w, r, permissionIndex, makePageAlert("Permission updated", warningAlert))
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
	logger.Debugf("PermissionName in url is %s", permissionname)
	return permissionname
}

func parsePermissionAssets() {
	//logger.Debug("Parsing permission assets...")
	assetNames := []string{indexView, newView, showView, editView, deleteView}
	parseAssets(&permissionAssetsBase, "layouts", "permission", layoutView, assetNames, permissionTemplates)
}

func parsePermissionExtAssets() {
	assetNames := []string{indexView, newView, showView, editView, deleteView}
	parseExtAssets(&permissionExtAssetsBase, "layouts", "permission", layoutView, assetNames, permissionExtTemplates)
}

func renderPermissionTemplate(w http.ResponseWriter, r *http.Request, groupName string, name string, page *Page) {
	if useExtTemplates {
		renderExtPermissionTemplate(w, r, groupName, name, page)
		return
	}
	renderIntPermissionTemplate(w, r, groupName, name, page)
}

// Render templates for the given name, template definition and data object
func renderIntPermissionTemplate(w http.ResponseWriter, r *http.Request, groupName string, name string, page *Page) {
	renderTemplate(w, r, permissionAssetsBase, groupName, name, permissionTemplates, page)
}

// Render templates for the given name, template definition and data object
func renderExtPermissionTemplate(w http.ResponseWriter, r *http.Request, groupName string, name string, page *Page) {
	//logger.Debugf("Autoreload: %t", bootstrap.AppConfig.IsAutoreloadOn())
	if bootstrap.AppConfig.IsAutoreloadOn() {
		logger.Debug("Reloading templates.")
		permissionExtTemplates = make(map[string]*htmlTemplate.Template)
		parsePermissionExtAssets()
	}
	renderExtTemplate(w, r, permissionExtAssetsBase, groupName, name, permissionExtTemplates, page)
}

func showPermissionError(w http.ResponseWriter, r *http.Request, page string, layout string, model interface{}, err error, alertKind string, cause error) {
	logger.Dump(cause)
	pageModel := makePage(model, makePageAlert(err.Error(), alertKind))
	renderPermissionTemplate(w, r, page, layoutView, pageModel)
}
