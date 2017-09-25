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
	_ "github.com/lib/pq" // Import pq without side effects
)

// ResourcePermissionRepository - ResourcePermission repository manager.
type ResourcePermissionRepository struct {
	DB *sqlx.DB
}

// MakeResourcePermissionRepository - ResourcePermissionRepository constructor.
func MakeResourcePermissionRepository() (ResourcePermissionRepository, error) {
	db, err := db.GetDbx()
	if err != nil {
		return ResourcePermissionRepository{}, err
	}
	return ResourcePermissionRepository{DB: db}, nil
}

// GetAll - GetAll ResourcePermissions in repo.
func (repo *ResourcePermissionRepository) GetAll(orgid string) ([]models.ResourcePermission, error) {
	resourcePermissions := []models.ResourcePermission{}
	err := repo.DB.Select(&resourcePermissions, "SELECT * FROM resource_permissions WHERE organization_id = $1 ORDER BY name ASC", orgid)
	return resourcePermissions, err
}

// Create - Persists a ResourcePermission in repo.
func (repo *ResourcePermissionRepository) Create(resourcePermission *models.ResourcePermission) error {
	resourcePermission.SetID()
	resourcePermission.SetCreationValues()
	tx := repo.DB.MustBegin()
	resourcePermissionInsertSQL := "INSERT INTO resource_permissions (id, name, description, geolocation, started_at, created_by, is_active, is_logical_deleted, created_at, updated_at, organization_id, resource_id, permission:_id) VALUES (:id, :name, :description, :geolocation, :started_at, :created_by, :is_active, :is_logical_deleted, :created_at, :updated_at, :organization_id, :resource_id, :permission_id)"
	_, err := tx.NamedExec(resourcePermissionInsertSQL, resourcePermission)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return err
}

// Get - Retrive a ResourcePermission in repo by its ID.
func (repo *ResourcePermissionRepository) Get(id string) (models.ResourcePermission, error) {
	u := models.ResourcePermission{}
	err := repo.DB.Get(&u, "SELECT * FROM resource_permissions WHERE id = $1", id)
	if err != nil {
		return u, err
	}
	return u, nil
}

// GetFromOrganization - Retrive a ResourcePermission in repo by its ID and Organization ID.
func (repo *ResourcePermissionRepository) GetFromOrganization(id string, orgid string) (models.ResourcePermission, error) {
	rp := models.ResourcePermission{}
	err := repo.DB.Get(&rp, "SELECT * FROM resource_permissions WHERE id = $1 AND organization_id = $2", id, orgid)
	if err != nil {
		return rp, err
	}
	return rp, nil
}

// GetByName - Retrive a ResourcePermission in repo by its Name.
func (repo *ResourcePermissionRepository) GetByName(name string) (models.ResourcePermission, error) {
	u := models.ResourcePermission{}
	err := repo.DB.Get(&u, "SELECT * FROM resource_permissions WHERE name = $1", name)
	if err != nil {
		return u, err
	}
	return u, nil
}

// Update - Update a resourcePermission in repo.
func (repo *ResourcePermissionRepository) Update(resourcePermission *models.ResourcePermission) error {
	// Update password and audit values
	resourcePermission.SetUpdateValues()
	// Current state
	reference, err := repo.Get(resourcePermission.ID.String)
	if err != nil {
		return err
	}
	// Customized query
	changes := ResourcePermissionChanges(resourcePermission, reference)
	number := len(changes)
	pos := 0
	last := number < 2
	var query bytes.Buffer
	query.WriteString("UPDATE resource_permissions SET ")
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
	query.WriteString(fmt.Sprintf("WHERE id = '%s';", resourcePermission.ID.String))
	//logger.Debug(query.String())
	tx := repo.DB.MustBegin()
	_, err = tx.NamedExec(query.String(), resourcePermission)
	if err != nil {
		return err
	}
	err = tx.Commit()
	return err
}

// Delete - Deletes resourcePermission from database.
func (repo *ResourcePermissionRepository) Delete(id string) error {
	tx := repo.DB.MustBegin()
	resourcePermissionDeleteSQL := fmt.Sprintf("DELETE FROM resource_permissions WHERE id = '%s'", id)
	_ = tx.MustExec(resourcePermissionDeleteSQL)
	err := tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

// DeleteFromOrganization - Deletes rolePermission from database.
func (repo *ResourcePermissionRepository) DeleteFromOrganization(id string, orgid string) error {
	tx := repo.DB.MustBegin()
	resourcePermissionDeleteSQL := fmt.Sprintf("DELETE FROM resource_permissions WHERE id = '%s' AND organization_id = '%s'", id, orgid)
	_ = tx.MustExec(resourcePermissionDeleteSQL)
	err := tx.Commit()
	if err != nil {
		return err
	}
	return nil
}
