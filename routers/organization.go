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
	"github.com/adrianpk/fundacja/controllers"

	"github.com/gorilla/mux"
)

// InitOrganizationRouter - Initialize API router for organizations.
func InitOrganizationRouter() *mux.Router {
	// Paths
	organizationPath := "/api/v1/organizations"
	// Router
	organizationRouter := apiV1Router.PathPrefix(organizationPath).Subrouter()
	// Resource
	organizationRouter.HandleFunc("", controllers.GetOrganizations).Methods("GET")
	organizationRouter.HandleFunc("", controllers.CreateOrganization).Methods("POST")
	organizationRouter.HandleFunc("/{organization}", controllers.GetOrganization).Methods("GET")
	organizationRouter.HandleFunc("/{organization}", controllers.UpdateOrganization).Methods("PUT")
	organizationRouter.HandleFunc("/{organization}", controllers.DeleteOrganization).Methods("DELETE")
	// Resource
	organizationRouter.HandleFunc("/{organization}/resources", controllers.GetResources).Methods("GET")
	organizationRouter.HandleFunc("/{organization}/resources", controllers.CreateResource).Methods("POST")
	organizationRouter.HandleFunc("/{organization}/resources/{resource}", controllers.GetResource).Methods("GET")
	organizationRouter.HandleFunc("/{organization}/resources/{resource}", controllers.UpdateResource).Methods("PUT")
	organizationRouter.HandleFunc("/{organization}/resources/{resource}", controllers.DeleteResource).Methods("DELETE")
	// Resource
	organizationRouter.HandleFunc("/{organization}/permissions", controllers.GetPermissions).Methods("GET")
	organizationRouter.HandleFunc("/{organization}/permissions", controllers.CreatePermission).Methods("POST")
	organizationRouter.HandleFunc("/{organization}/permissions/{permission}", controllers.GetPermission).Methods("GET")
	organizationRouter.HandleFunc("/{organization}/permissions/{permission}", controllers.UpdatePermission).Methods("PUT")
	organizationRouter.HandleFunc("/{organization}/permissions/{permission}", controllers.DeletePermission).Methods("DELETE")
	// Resource
	organizationRouter.HandleFunc("/{organization}/resource-permissions", controllers.GetResourcePermissions).Methods("GET")
	organizationRouter.HandleFunc("/{organization}/resource-permissions", controllers.CreateResourcePermission).Methods("POST")
	organizationRouter.HandleFunc("/{organization}/resource-permissions/{resource-permission}", controllers.GetResourcePermission).Methods("GET")
	organizationRouter.HandleFunc("/{organization}/resource-permissions/{resource-permission}", controllers.UpdateResourcePermission).Methods("PUT")
	organizationRouter.HandleFunc("/{organization}/resource-permissions/{resource-permission}", controllers.DeleteResourcePermission).Methods("DELETE")
	// Resource
	organizationRouter.HandleFunc("/{organization}/roles", controllers.GetRoles).Methods("GET")
	organizationRouter.HandleFunc("/{organization}/roles", controllers.CreateRole).Methods("POST")
	organizationRouter.HandleFunc("/{organization}/roles/{role}", controllers.GetRole).Methods("GET")
	organizationRouter.HandleFunc("/{organization}/roles/{role}", controllers.UpdateRole).Methods("PUT")
	organizationRouter.HandleFunc("/{organization}/roles/{role}", controllers.DeleteRole).Methods("DELETE")
	// Resource
	organizationRouter.HandleFunc("/{organization}/role-permissions", controllers.GetRolePermissions).Methods("GET")
	organizationRouter.HandleFunc("/{organization}/role-permissions", controllers.CreateRolePermission).Methods("POST")
	organizationRouter.HandleFunc("/{organization}/role-permissions/{role-permission}", controllers.GetRolePermission).Methods("GET")
	organizationRouter.HandleFunc("/{organization}/role-permissions/{role-permission}", controllers.UpdateRolePermission).Methods("PUT")
	organizationRouter.HandleFunc("/{organization}/role-permissions/{role-permission}", controllers.DeleteRolePermission).Methods("DELETE")
	// Resource
	organizationRouter.HandleFunc("/{organization}/user-roles", controllers.GetUserRoles).Methods("GET")
	organizationRouter.HandleFunc("/{organization}/user-roles", controllers.CreateUserRole).Methods("POST")
	organizationRouter.HandleFunc("/{organization}/user-roles/{user-role}", controllers.GetUserRole).Methods("GET")
	organizationRouter.HandleFunc("/{organization}/user-roles/{user-role}", controllers.UpdateUserRole).Methods("PUT")
	organizationRouter.HandleFunc("/{organization}/user-roles/{user-role}", controllers.DeleteUserRole).Methods("DELETE")
	return organizationRouter
}
