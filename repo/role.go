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

// RoleRepository - Role repository manager.
type RoleRepository struct {
	DB *sqlx.DB
}

// MakeRoleRepository - RoleRepository constructor.
func MakeRoleRepository() (RoleRepository, error) {
	db, err := db.GetDbx()
	if err != nil {
		return RoleRepository{}, err
	}
	return RoleRepository{DB: db}, nil
}

// GetAll - GetAll Roles in repo.
func (repo *RoleRepository) GetAll(orgid string) ([]models.Role, error) {
	roles := []models.Role{}
	err := repo.DB.Select(&roles, "SELECT * FROM roles WHERE organization_id = $1 ORDER BY name ASC", orgid)
	return roles, err
}

// Create - Persists a Role in repo.
func (repo *RoleRepository) Create(role *models.Role) error {
	role.SetID()
	role.SetCreationValues()
	tx := repo.DB.MustBegin()
	roleInsertSQL := "INSERT INTO roles (id, name, description, geolocation, started_at, created_by, is_active, is_logical_deleted, created_at, updated_at, organization_id) VALUES (:id, :name, :description, :geolocation, :started_at, :created_by, :is_active, :is_logical_deleted, :created_at, :updated_at, :organization_id)"
	_, err := tx.NamedExec(roleInsertSQL, role)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return err
}

// Get - Retrive a Role in repo by its ID.
func (repo *RoleRepository) Get(id string) (models.Role, error) {
	u := models.Role{}
	err := repo.DB.Get(&u, "SELECT * FROM roles WHERE id = $1", id)
	if err != nil {
		return u, err
	}
	return u, nil
}

// GetByRoleName - Retrive a Role in repo by its rolename.
func (repo *RoleRepository) GetByName(name string) (models.Role, error) {
	u := models.Role{}
	err := repo.DB.Get(&u, "SELECT * FROM roles WHERE name = $1", name)
	if err != nil {
		return u, err
	}
	return u, nil
}

// Update - Update a role in repo.
func (repo *RoleRepository) Update(role *models.Role) error {
	// Update password and audit values
	role.SetUpdateValues()
	// Current state
	reference, err := repo.Get(role.ID.String)
	if err != nil {
		return err
	}
	// Customized query
	changes := RoleChanges(role, reference)
	number := len(changes)
	pos := 0
	last := number < 2
	var query bytes.Buffer
	query.WriteString("UPDATE roles SET ")
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
	query.WriteString(fmt.Sprintf("WHERE id = '%s';", role.ID.String))
	//logger.Debug(query.String())
	tx := repo.DB.MustBegin()
	_, err = tx.NamedExec(query.String(), role)
	if err != nil {
		return err
	}
	err = tx.Commit()
	return err
}

// Delete - Deletes role from database.
func (repo *RoleRepository) Delete(id string) error {
	tx := repo.DB.MustBegin()
	roleDeleteSQL := fmt.Sprintf("DELETE FROM roles WHERE id = '%s'", id)
	_ = tx.MustExec(roleDeleteSQL)
	err := tx.Commit()
	if err != nil {
		return err
	}
	return nil
}
