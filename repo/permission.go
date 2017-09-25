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

package repo

import (
	"bytes"
	"fmt"

	"github.com/adrianpk/pulap/db"
	"github.com/adrianpk/pulap/logger"
	"github.com/adrianpk/pulap/models"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Import pq without side effects
)

// PermissionRepository - Permission repository manager.
type PermissionRepository struct {
	DB *sqlx.DB
}

// MakePermissionRepository - PermissionRepository constructor.
func MakePermissionRepository() (PermissionRepository, error) {
	db, err := db.GetDbx()
	if err != nil {
		return PermissionRepository{}, err
	}
	return PermissionRepository{DB: db}, nil
}

// GetAll - GetAll Permissions in repo.
func (repo *PermissionRepository) GetAll(orgid string) ([]models.Permission, error) {
	permissions := []models.Permission{}
	err := repo.DB.Select(&permissions, "SELECT * FROM permissions WHERE organization_id = $1 ORDER BY name ASC", orgid)
	return permissions, err
}

// Create - Persists a Permission in repo.
func (repo *PermissionRepository) Create(permission *models.Permission) error {
	permission.SetID()
	permission.SetCreationValues()
	tx := repo.DB.MustBegin()
	permissionInsertSQL := "INSERT INTO permissions (id, name, description, geolocation, created_by, is_active, is_logical_deleted, created_at, updated_at, organization_id) VALUES (:id, :name, :description, :geolocation, :created_by, :is_active, :is_logical_deleted, :created_at, :updated_at, :organization_id)"
	_, err := tx.NamedExec(permissionInsertSQL, permission)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return err
}

// Get - Retrive a Permission in repo by its ID.
func (repo *PermissionRepository) Get(id string) (models.Permission, error) {
	u := models.Permission{}
	err := repo.DB.Get(&u, "SELECT * FROM permissions WHERE id = $1", id)
	if err != nil {
		return u, err
	}
	return u, nil
}

// GetByName - Retrive a Permission in repo by its permissionname.
func (repo *PermissionRepository) GetByName(name string) (models.Permission, error) {
	u := models.Permission{}
	err := repo.DB.Get(&u, "SELECT * FROM permissions WHERE name = $1", name)
	if err != nil {
		return u, err
	}
	return u, nil
}

// Update - Update a permission in repo.
func (repo *PermissionRepository) Update(permission *models.Permission) error {
	// Update password and audit values
	permission.SetUpdateValues()
	// Current state
	reference, err := repo.Get(permission.ID.String)
	if err != nil {
		return err
	}
	// Customized query
	changes := PermissionChanges(permission, reference)
	number := len(changes)
	pos := 0
	last := number < 2
	var query bytes.Buffer
	query.WriteString("UPDATE permissions SET ")
	for field, structField := range changes {
		var partial string
		if last {
			partial = fmt.Sprintf("%v = %v ", field, structField)
		} else {
			partial = fmt.Sprintf("%v = %v, ", field, structField)
		}
		query.WriteString(partial)
		pos = pos + 1
		last = pos == number-1
	}
	query.WriteString(fmt.Sprintf("WHERE id = '%s';", permission.ID.String))
	//logger.Debug(query.String())
	tx := repo.DB.MustBegin()
	_, err = tx.NamedExec(query.String(), permission)
	if err != nil {
		return err
	}
	err = tx.Commit()
	return err
}

// Delete - Deletes permission from database.
func (repo *PermissionRepository) Delete(id string) error {
	tx := repo.DB.MustBegin()
	permissionDeleteSQL := fmt.Sprintf("DELETE FROM permissions WHERE id = '%s'", id)
	_ = tx.MustExec(permissionDeleteSQL)
	err := tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

// GetUserPermissionIDs - Returns an array of Permissions IDs that are assigned to some User
func (repo *PermissionRepository) GetUserPermissionIDs(userID string) ([]string, error) {
	permisisons := []string{}
	var query bytes.Buffer
	query.WriteString("SELECT permissions.id FROM permissions INNER JOIN role_permissions ")
	query.WriteString("ON permissions.id = role_permissions.permission_id ")
	query.WriteString("INNER JOIN user_roles ")
	query.WriteString("ON role_permissions.role_id = user_roles.role_id ")
	query.WriteString("INNER JOIN users ")
	query.WriteString("ON user_roles.user_id = users.id ")
	query.WriteString(fmt.Sprintf("WHERE users.id = '%s';", userID))
	logger.Debugf("Query: %s", query.String())
	tx := repo.DB.MustBegin()
	err := repo.DB.Select(&permisisons, query.String())
	if err != nil {
		logger.Dump(err)
		return permisisons, err
	}
	err = tx.Commit()
	if err != nil {
		return permisisons, err
	}
	return permisisons, nil
}

// GetEnablingPermissionIDs - Returns an array of Permissions IDs that enable the use of some Resource
func (repo *PermissionRepository) GetEnablingPermissionIDs(resourceIDorTag string) ([]string, error) {
	permisisons := []string{}
	var query bytes.Buffer
	query.WriteString("SELECT permissions.id FROM permissions INNER JOIN resource_permissions ")
	query.WriteString("ON permissions.id = resource_permissions.permission_id ")
	query.WriteString("INNER JOIN resources ")
	query.WriteString("ON resources.id = resource_permissions.resource_id ")
	if len(resourceIDorTag) == 36 {
		query.WriteString(fmt.Sprintf("WHERE resources.id = %s ", resourceIDorTag))
	} else {
		query.WriteString(fmt.Sprintf("WHERE resources.tag = '%s';", resourceIDorTag))
	}
	tx := repo.DB.MustBegin()
	err := repo.DB.Select(&permisisons, query.String())
	if err != nil {
		logger.Dump(err)
		return permisisons, err
	}
	err = tx.Commit()
	if err != nil {
		return permisisons, err
	}
	return permisisons, nil
}

// SELECT count(permissions.id) FROM permissions INNER JOIN resource_permissions ON permissions.id = resource_permissions.permission_id INNER JOIN resources ON resources.id = resource_permissions.resource_id WHERE resources.tag = 'f254cfe4' AND permissions.id IN (SELECT permissions.id FROM permissions INNER JOIN role_permissions ON permissions.id = role_permissions.permission_id INNER JOIN user_roles ON role_permissions.role_id = user_roles.role_id INNER JOIN users ON user_roles.user_id = users.id WHERE users.id = '5958b185-8150-4aae-b53f-0c44771ddec5')

// HasPermission - Returns true if user has at least one permission over a resource.
func (repo *PermissionRepository) HasPermission(resourceIDorTag, userID string) (bool, error) {
	hasPermission := false
	var query bytes.Buffer
	query.WriteString("SELECT COUNT(permissions.id) > 0 FROM permissions INNER JOIN resource_permissions ")
	query.WriteString("ON permissions.id = resource_permissions.permission_id ")
	query.WriteString("INNER JOIN resources ")
	query.WriteString("ON resources.id = resource_permissions.resource_id ")
	if len(resourceIDorTag) == 36 {
		query.WriteString(fmt.Sprintf("WHERE resources.id = %s ", resourceIDorTag))
	} else {
		query.WriteString(fmt.Sprintf("WHERE resources.tag = '%s' ", resourceIDorTag))
	}
	query.WriteString("AND permissions.id ")
	query.WriteString("IN (SELECT permissions.id FROM permissions INNER JOIN role_permissions ")
	query.WriteString("ON permissions.id = role_permissions.permission_id ")
	query.WriteString("INNER JOIN user_roles ")
	query.WriteString("ON role_permissions.role_id = user_roles.role_id ")
	query.WriteString("INNER JOIN users ")
	query.WriteString("ON user_roles.user_id = users.id ")
	query.WriteString(fmt.Sprintf("WHERE users.id = '%s')", userID))

	logger.Debugf("Query: %s", query.String())

	tx := repo.DB.MustBegin()
	err := repo.DB.Get(&hasPermission, query.String())
	if err != nil {
		logger.Dump(err)
		return false, err
	}
	err = tx.Commit()
	if err != nil {
		return false, err
	}
	return hasPermission, nil
}
