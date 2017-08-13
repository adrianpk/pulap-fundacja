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
	"github.com/adrianpk/fundacja/api"

	"github.com/gorilla/mux"
)

// InitOrganizationRouter - Initialize API router for organizations.
func InitOrganizationRouter() *mux.Router {
	// Paths
	organizationPath := "/api/v1/organizations"
	// Router
	organizationRouter := apiV1Router.PathPrefix(organizationPath).Subrouter()
	// Resource
	organizationRouter.HandleFunc("", api.GetOrganizations).Methods("GET")
	organizationRouter.HandleFunc("", api.CreateOrganization).Methods("POST")
	organizationRouter.HandleFunc("/{organization}", api.GetOrganization).Methods("GET")
	organizationRouter.HandleFunc("/{organization}", api.UpdateOrganization).Methods("PUT")
	organizationRouter.HandleFunc("/{organization}", api.DeleteOrganization).Methods("DELETE")
	// Resource
	organizationRouter.HandleFunc("/{organization}/resources", api.GetResources).Methods("GET")
	organizationRouter.HandleFunc("/{organization}/resources", api.CreateResource).Methods("POST")
	organizationRouter.HandleFunc("/{organization}/resources/{resource}", api.GetResource).Methods("GET")
	organizationRouter.HandleFunc("/{organization}/resources/{resource}", api.UpdateResource).Methods("PUT")
	organizationRouter.HandleFunc("/{organization}/resources/{resource}", api.DeleteResource).Methods("DELETE")
	// Resource
	organizationRouter.HandleFunc("/{organization}/permissions", api.GetPermissions).Methods("GET")
	organizationRouter.HandleFunc("/{organization}/permissions", api.CreatePermission).Methods("POST")
	organizationRouter.HandleFunc("/{organization}/permissions/{permission}", api.GetPermission).Methods("GET")
	organizationRouter.HandleFunc("/{organization}/permissions/{permission}", api.UpdatePermission).Methods("PUT")
	organizationRouter.HandleFunc("/{organization}/permissions/{permission}", api.DeletePermission).Methods("DELETE")
	// Resource
	organizationRouter.HandleFunc("/{organization}/resource-permissions", api.GetResourcePermissions).Methods("GET")
	organizationRouter.HandleFunc("/{organization}/resource-permissions", api.CreateResourcePermission).Methods("POST")
	organizationRouter.HandleFunc("/{organization}/resource-permissions/{resource-permission}", api.GetResourcePermission).Methods("GET")
	organizationRouter.HandleFunc("/{organization}/resource-permissions/{resource-permission}", api.UpdateResourcePermission).Methods("PUT")
	organizationRouter.HandleFunc("/{organization}/resource-permissions/{resource-permission}", api.DeleteResourcePermission).Methods("DELETE")
	// Resource
	organizationRouter.HandleFunc("/{organization}/roles", api.GetRoles).Methods("GET")
	organizationRouter.HandleFunc("/{organization}/roles", api.CreateRole).Methods("POST")
	organizationRouter.HandleFunc("/{organization}/roles/{role}", api.GetRole).Methods("GET")
	organizationRouter.HandleFunc("/{organization}/roles/{role}", api.UpdateRole).Methods("PUT")
	organizationRouter.HandleFunc("/{organization}/roles/{role}", api.DeleteRole).Methods("DELETE")
	// Resource
	organizationRouter.HandleFunc("/{organization}/role-permissions", api.GetRolePermissions).Methods("GET")
	organizationRouter.HandleFunc("/{organization}/role-permissions", api.CreateRolePermission).Methods("POST")
	organizationRouter.HandleFunc("/{organization}/role-permissions/{role-permission}", api.GetRolePermission).Methods("GET")
	organizationRouter.HandleFunc("/{organization}/role-permissions/{role-permission}", api.UpdateRolePermission).Methods("PUT")
	organizationRouter.HandleFunc("/{organization}/role-permissions/{role-permission}", api.DeleteRolePermission).Methods("DELETE")
	// Resource
	organizationRouter.HandleFunc("/{organization}/user-roles", api.GetUserRoles).Methods("GET")
	organizationRouter.HandleFunc("/{organization}/user-roles", api.CreateUserRole).Methods("POST")
	organizationRouter.HandleFunc("/{organization}/user-roles/{user-role}", api.GetUserRole).Methods("GET")
	organizationRouter.HandleFunc("/{organization}/user-roles/{user-role}", api.UpdateUserRole).Methods("PUT")
	organizationRouter.HandleFunc("/{organization}/user-roles/{user-role}", api.DeleteUserRole).Methods("DELETE")
	return organizationRouter
}
