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
	"github.com/adrianpk/pulap/api"

	"github.com/gorilla/mux"
)

// InitOrganizationRouter - Initialize API router for organizations.
func InitAPIOrganizationRouter() *mux.Router {
	// Paths
	organizationPath := "/api/v1/organizations"
	// Router
	organizationAPIRouter := apiV1Router.PathPrefix(organizationPath).Subrouter()
	// Resource
	organizationAPIRouter.HandleFunc("", api.GetOrganizations).Methods("GET")
	organizationAPIRouter.HandleFunc("", api.CreateOrganization).Methods("POST")
	organizationAPIRouter.HandleFunc("/{organization}", api.GetOrganization).Methods("GET")
	organizationAPIRouter.HandleFunc("/{organization}", api.UpdateOrganization).Methods("PUT")
	organizationAPIRouter.HandleFunc("/{organization}", api.DeleteOrganization).Methods("DELETE")
	// Resource
	organizationAPIRouter.HandleFunc("/{organization}/resources", api.GetResources).Methods("GET")
	organizationAPIRouter.HandleFunc("/{organization}/resources", api.CreateResource).Methods("POST")
	organizationAPIRouter.HandleFunc("/{organization}/resources/{resource}", api.GetResource).Methods("GET")
	organizationAPIRouter.HandleFunc("/{organization}/resources/{resource}", api.UpdateResource).Methods("PUT")
	organizationAPIRouter.HandleFunc("/{organization}/resources/{resource}", api.DeleteResource).Methods("DELETE")
	// Resource
	organizationAPIRouter.HandleFunc("/{organization}/permissions", api.GetPermissions).Methods("GET")
	organizationAPIRouter.HandleFunc("/{organization}/permissions", api.CreatePermission).Methods("POST")
	organizationAPIRouter.HandleFunc("/{organization}/permissions/{permission}", api.GetPermission).Methods("GET")
	organizationAPIRouter.HandleFunc("/{organization}/permissions/{permission}", api.UpdatePermission).Methods("PUT")
	organizationAPIRouter.HandleFunc("/{organization}/permissions/{permission}", api.DeletePermission).Methods("DELETE")
	// Resource
	organizationAPIRouter.HandleFunc("/{organization}/resource-permissions", api.GetResourcePermissions).Methods("GET")
	organizationAPIRouter.HandleFunc("/{organization}/resource-permissions", api.CreateResourcePermission).Methods("POST")
	organizationAPIRouter.HandleFunc("/{organization}/resource-permissions/{resource-permission}", api.GetResourcePermission).Methods("GET")
	organizationAPIRouter.HandleFunc("/{organization}/resource-permissions/{resource-permission}", api.UpdateResourcePermission).Methods("PUT")
	organizationAPIRouter.HandleFunc("/{organization}/resource-permissions/{resource-permission}", api.DeleteResourcePermission).Methods("DELETE")
	// Resource
	organizationAPIRouter.HandleFunc("/{organization}/roles", api.GetRoles).Methods("GET")
	organizationAPIRouter.HandleFunc("/{organization}/roles", api.CreateRole).Methods("POST")
	organizationAPIRouter.HandleFunc("/{organization}/roles/{role}", api.GetRole).Methods("GET")
	organizationAPIRouter.HandleFunc("/{organization}/roles/{role}", api.UpdateRole).Methods("PUT")
	organizationAPIRouter.HandleFunc("/{organization}/roles/{role}", api.DeleteRole).Methods("DELETE")
	// Resource
	organizationAPIRouter.HandleFunc("/{organization}/role-permissions", api.GetRolePermissions).Methods("GET")
	organizationAPIRouter.HandleFunc("/{organization}/role-permissions", api.CreateRolePermission).Methods("POST")
	organizationAPIRouter.HandleFunc("/{organization}/role-permissions/{role-permission}", api.GetRolePermission).Methods("GET")
	organizationAPIRouter.HandleFunc("/{organization}/role-permissions/{role-permission}", api.UpdateRolePermission).Methods("PUT")
	organizationAPIRouter.HandleFunc("/{organization}/role-permissions/{role-permission}", api.DeleteRolePermission).Methods("DELETE")
	// Resource
	organizationAPIRouter.HandleFunc("/{organization}/user-roles", api.GetUserRoles).Methods("GET")
	organizationAPIRouter.HandleFunc("/{organization}/user-roles", api.CreateUserRole).Methods("POST")
	organizationAPIRouter.HandleFunc("/{organization}/user-roles/{user-role}", api.GetUserRole).Methods("GET")
	organizationAPIRouter.HandleFunc("/{organization}/user-roles/{user-role}", api.UpdateUserRole).Methods("PUT")
	organizationAPIRouter.HandleFunc("/{organization}/user-roles/{user-role}", api.DeleteUserRole).Methods("DELETE")
	return organizationAPIRouter
}
