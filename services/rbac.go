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

package services

import (
	//"database/sql"
	//httpcontext "github.com/gorilla/context"
	//"github.com/jmoiron/sqlx"

	"github.com/adrianpk/pulap/logger"

	_ "github.com/lib/pq" // Import pq without side effects

	"github.com/adrianpk/pulap/repo"
)

// IsAllowed - Returns true if user has at leeast one enabling permission over some resource.
// Panic if not.
func IsAllowed(resourceIDOrTag, userID string) {
	hasPermission := HasPermission(resourceIDOrTag, userID)
	logger.Debugf("User: %s", userID)
	logger.Debugf("Resource: %s", resourceIDOrTag)
	logger.Debugf("Allowed? %t", hasPermission)
	if hasPermission {
		return
	}
	panic("Forbidden")
}

// HasPermission - Returns true if user has at least one permission over a resource.
func HasPermission(resourceIDorTag, userID string) bool {
	// Get repo
	permissionRepo, err := repo.MakePermissionRepository()
	if err != nil {
		return false
	}
	// Select
	hasPermission, _ := permissionRepo.HasPermission(resourceIDorTag, userID)
	return hasPermission
}

// GetUserPermissionsIDs - Returns an array of Permissions IDs that are assigned to some User
func GetUserPermissionsIDs(userID string) []string {
	// Get repo
	permissionRepo, err := repo.MakePermissionRepository()
	if err != nil {
		return []string{}
	}
	// Select
	permissionIDs, _ := permissionRepo.GetUserPermissionIDs(userID)
	return permissionIDs
}

// GetEnablingPermissionIDs - Returns an array of Permissions IDs that enable the use of some Resource
func GetEnablingPermissionIDs(resourceIDorTag string) []string {
	// Get repo
	permissionRepo, err := repo.MakePermissionRepository()
	if err != nil {
		return []string{}
	}
	// Select
	permissionIDs, _ := permissionRepo.GetEnablingPermissionIDs(resourceIDorTag)
	return permissionIDs
}
