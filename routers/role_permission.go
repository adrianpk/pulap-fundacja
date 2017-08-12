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

package routers

import (
	"github.com/adrianpk/fundacja/bootstrap"
	"github.com/adrianpk/fundacja/controllers"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
)

// InitAPIRolePermissionRouter - Initialize API router for role-permissions.
func InitAPIRolePermissionRouter() *mux.Router {
	// Paths
	rolePermissionsPath := "/organizations/{organization}/role-permissions"
	// Router
	rolePermissionRouter := apiV1Router.PathPrefix(rolePermissionsPath).Subrouter()
	// Organization resource
	rolePermissionRouter.HandleFunc("", controllers.GetRolePermissions).Methods("GET")
	rolePermissionRouter.HandleFunc("", controllers.CreateRolePermission).Methods("POST")
	rolePermissionRouter.HandleFunc("/{role-permission}", controllers.GetRolePermission).Methods("GET")
	rolePermissionRouter.HandleFunc("/{role-permission}", controllers.UpdateRolePermission).Methods("PUT")
	rolePermissionRouter.HandleFunc("/{role-permission}", controllers.DeleteRolePermission).Methods("DELETE")
	// Middleware
	apiV1Router.Handle(rolePermissionsPath, negroni.New(
		negroni.HandlerFunc(bootstrap.Authorize),
		negroni.Wrap(rolePermissionRouter),
	))

	return rolePermissionRouter
}
