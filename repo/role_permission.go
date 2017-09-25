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
	"github.com/adrianpk/pulap/models"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Import pq without side effect
)

// RolePermissionRepository - RolePermission repository manager.
type RolePermissionRepository struct {
	DB *sqlx.DB
}

// MakeRolePermissionRepository - RolePermissionRepository constructor.
func MakeRolePermissionRepository() (RolePermissionRepository, error) {
	db, err := db.GetDbx()
	if err != nil {
		return RolePermissionRepository{}, err
	}
	return RolePermissionRepository{DB: db}, nil
}

// GetAll - GetAll RolePermissions in repo.
func (repo *RolePermissionRepository) GetAll(orgid string) ([]models.RolePermission, error) {
	rolePermissions := []models.RolePermission{}
	err := repo.DB.Select(&rolePermissions, "SELECT * FROM role_permissions WHERE organization_id = $1 ORDER BY name ASC", orgid)
	return rolePermissions, err
}

// Create - Persists a RolePermission in repo.
func (repo *RolePermissionRepository) Create(rolePermission *models.RolePermission) error {
	rolePermission.SetID()
	rolePermission.SetCreationValues()
	tx := repo.DB.MustBegin()
	rolePermissionInsertSQL := "INSERT INTO role_permissions (id, name, description, created_by, is_active, is_logical_deleted, created_at, updated_at, organization_id, role_id, permission_id) VALUES (:id, :name, :description, :created_by, :is_active, :is_logical_deleted, :created_at, :updated_at, :organization_id, :role_id, :permission_id)"
	_, err := tx.NamedExec(rolePermissionInsertSQL, rolePermission)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

// Get - Retrive a RolePermission in repo by its ID.
func (repo *RolePermissionRepository) Get(id string) (models.RolePermission, error) {
	rp := models.RolePermission{}
	err := repo.DB.Get(&rp, "SELECT * FROM role_permissions WHERE id = $1", id)
	if err != nil {
		return rp, err
	}
	return rp, nil
}

// GetFromOrganization - Retrive a RolePermission in repo by its ID and Organization ID.
func (repo *RolePermissionRepository) GetFromOrganization(id string, orgid string) (models.RolePermission, error) {
	rp := models.RolePermission{}
	err := repo.DB.Get(&rp, "SELECT * FROM role_permissions WHERE id = $1 AND organization_id = $2", id, orgid)
	if err != nil {
		return rp, err
	}
	return rp, nil
}

// GetByName - Retrive a RolePermission in repo by its rolePermissionname.
func (repo *RolePermissionRepository) GetByName(name string) (models.RolePermission, error) {
	rp := models.RolePermission{}
	err := repo.DB.Get(&rp, "SELECT * FROM role_permissions WHERE name = $1", name)
	if err != nil {
		return rp, err
	}
	return rp, nil
}

// Update - Update a rolePermission in repo.
func (repo *RolePermissionRepository) Update(rolePermission *models.RolePermission) error {
	// Update password and audit values
	rolePermission.SetUpdateValues()
	// Current state
	reference, err := repo.Get(rolePermission.ID.String)
	if err != nil {
		return err
	}
	// Customized query
	changes := RolePermissionChanges(rolePermission, reference)
	number := len(changes)
	pos := 0
	last := number < 2
	var query bytes.Buffer
	query.WriteString("UPDATE role_permissions SET ")
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
	query.WriteString(fmt.Sprintf("WHERE id = '%s';", rolePermission.ID.String))
	//logger.Debug(query.String())
	tx := repo.DB.MustBegin()
	_, err = tx.NamedExec(query.String(), rolePermission)
	if err != nil {
		return err
	}
	err = tx.Commit()
	return err
}

// Delete - Deletes rolePermission from database.
func (repo *RolePermissionRepository) Delete(id string) error {
	tx := repo.DB.MustBegin()
	rolePermissionDeleteSQL := fmt.Sprintf("DELETE FROM role_permissions WHERE id = '%s'", id)
	_ = tx.MustExec(rolePermissionDeleteSQL)
	err := tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

// DeleteFromOrganization - Deletes rolePermission from database.
func (repo *RolePermissionRepository) DeleteFromOrganization(id string, orgid string) error {
	tx := repo.DB.MustBegin()
	rolePermissionDeleteSQL := fmt.Sprintf("DELETE FROM role_permissions WHERE id = '%s' AND organization_id = '%s'", id, orgid)
	_ = tx.MustExec(rolePermissionDeleteSQL)
	err := tx.Commit()
	if err != nil {
		return err
	}
	return nil
}
