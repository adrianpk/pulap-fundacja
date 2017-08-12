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

	"github.com/adrianpk/fundacja/db"
	"github.com/adrianpk/fundacja/models"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Import pq without side effects
)

// ResourceRepository - Resource repository manager.
type ResourceRepository struct {
	DB *sqlx.DB
}

// MakeResourceRepository - ResourceRepository constructor.
func MakeResourceRepository() (ResourceRepository, error) {
	db, err := db.GetDbx()
	if err != nil {
		return ResourceRepository{}, err
	}
	return ResourceRepository{DB: db}, nil
}

// GetAll - GetAll RolePermissions in repo.
func (repo *ResourceRepository) GetAll(orgid string) ([]models.Resource, error) {
	resources := []models.Resource{}
	err := repo.DB.Select(&resources, "SELECT * FROM resources WHERE organization_id = $1 ORDER BY name ASC", orgid)
	return resources, err
}

// Create - Persists a Resource in repo.
func (repo *ResourceRepository) Create(resource *models.Resource) error {
	resource.SetID()
	resource.SetCreationValues()
	tx := repo.DB.MustBegin()
	resourceInsertSQL := "INSERT INTO resources (id, name, description, geolocation, created_by, is_active, is_logical_deleted, created_at, updated_at, organization_id) VALUES (:id, :name, :description, :geolocation, :created_by, :is_active, :is_logical_deleted, :created_at, :updated_at, :organization_id)"
	_, err := tx.NamedExec(resourceInsertSQL, resource)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return err
}

// Get - Retrive a Resource in repo by its ID.
func (repo *ResourceRepository) Get(id string) (models.Resource, error) {
	u := models.Resource{}
	err := repo.DB.Get(&u, "SELECT * FROM resources WHERE id = $1", id)
	if err != nil {
		return u, err
	}
	return u, nil
}

// GetByName - Retrive a Resource in repo by its resourcename.
func (repo *ResourceRepository) GetByName(name string) (models.Resource, error) {
	u := models.Resource{}
	err := repo.DB.Get(&u, "SELECT * FROM resources WHERE name = $1", name)
	if err != nil {
		return u, err
	}
	return u, nil
}

// Update - Update a resource in repo.
func (repo *ResourceRepository) Update(resource *models.Resource) error {
	// Update password and audit values
	resource.SetUpdateValues()
	// Current state
	reference, err := repo.Get(resource.ID.String)
	if err != nil {
		return err
	}
	// Customized query
	changes := ResourceChanges(resource, reference)
	number := len(changes)
	pos := 0
	last := number < 2
	var query bytes.Buffer
	query.WriteString("UPDATE resources SET ")
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
	query.WriteString(fmt.Sprintf("WHERE id = '%s';", resource.ID.String))
	//logger.Debug(query.String())
	tx := repo.DB.MustBegin()
	_, err = tx.NamedExec(query.String(), resource)
	if err != nil {
		return err
	}
	err = tx.Commit()
	return err
}

// Delete - Deletes resource from database.
func (repo *ResourceRepository) Delete(id string) error {
	tx := repo.DB.MustBegin()
	resourceDeleteSQL := fmt.Sprintf("DELETE FROM resources WHERE id = '%s'", id)
	_ = tx.MustExec(resourceDeleteSQL)
	err := tx.Commit()
	if err != nil {
		return err
	}
	return nil
}
