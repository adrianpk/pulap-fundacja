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

// UserRoleRepository - UserRole repository manager.
type UserRoleRepository struct {
	DB *sqlx.DB
}

// MakeUserRoleRepository - UserRoleRepository constructor.
func MakeUserRoleRepository() (UserRoleRepository, error) {
	db, err := db.GetDbx()
	if err != nil {
		return UserRoleRepository{}, err
	}
	return UserRoleRepository{DB: db}, nil
}

// GetAll - GetAll UserRoles in repo.
func (repo *UserRoleRepository) GetAll(orgid string) ([]models.UserRole, error) {
	userRoles := []models.UserRole{}
	err := repo.DB.Select(&userRoles, "SELECT * FROM user_roles WHERE organization_id = $1 ORDER BY name ASC", orgid)
	return userRoles, err
}

// Create - Persists a UserRole in repo.
func (repo *UserRoleRepository) Create(userRole *models.UserRole) error {
	userRole.SetID()
	userRole.SetCreationValues()
	tx := repo.DB.MustBegin()
	userRoleInsertSQL := "INSERT INTO user_roles (id, name, description, created_by, is_active, is_logical_deleted, created_at, updated_at, organization_id, user_id, role_id) VALUES (:id, :name, :description, :created_by, :is_active, :is_logical_deleted, :created_at, :updated_at, :organization_id, :user_id, :role_id)"
	_, err := tx.NamedExec(userRoleInsertSQL, userRole)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return err
}

// Get - Retrive a UserRole in repo by its ID.
func (repo *UserRoleRepository) Get(id string) (models.UserRole, error) {
	u := models.UserRole{}
	err := repo.DB.Get(&u, "SELECT * FROM user_roles WHERE id = $1", id)
	if err != nil {
		return u, err
	}
	return u, nil
}

// GetFromOrganization - Retrive a RolePermission in repo by its ID and Organization ID.
func (repo *UserRoleRepository) GetFromOrganization(id string, orgid string) (models.UserRole, error) {
	rp := models.UserRole{}
	err := repo.DB.Get(&rp, "SELECT * FROM user_roles WHERE id = $1 AND organization_id = $2", id, orgid)
	if err != nil {
		return rp, err
	}
	return rp, nil
}

// GetByName - Retrive a UserRole in repo by its userRoleName.
func (repo *UserRoleRepository) GetByName(name string) (models.UserRole, error) {
	u := models.UserRole{}
	err := repo.DB.Get(&u, "SELECT * FROM user_roles WHERE name = $1", name)
	if err != nil {
		return u, err
	}
	return u, nil
}

// Update - Update a userRole in repo.
func (repo *UserRoleRepository) Update(userRole *models.UserRole) error {
	// Update password and audit values
	userRole.SetUpdateValues()
	// Current state
	reference, err := repo.Get(userRole.ID.String)
	if err != nil {
		return err
	}
	// Customized query
	changes := UserRoleChanges(userRole, reference)
	number := len(changes)
	pos := 0
	last := number < 2
	var query bytes.Buffer
	query.WriteString("UPDATE user_roles SET ")
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
	query.WriteString(fmt.Sprintf("WHERE id = '%s';", userRole.ID.String))
	//logger.Debug(query.String())
	tx := repo.DB.MustBegin()
	_, err = tx.NamedExec(query.String(), userRole)
	if err != nil {
		return err
	}
	err = tx.Commit()
	return err
}

// Delete - Deletes userRole from database.
func (repo *UserRoleRepository) Delete(id string) error {
	tx := repo.DB.MustBegin()
	userRoleDeleteSQL := fmt.Sprintf("DELETE FROM user_roles WHERE id = '%s'", id)
	_ = tx.MustExec(userRoleDeleteSQL)
	err := tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

// DeleteFromOrganization - Deletes rolePermission from database.
func (repo *UserRoleRepository) DeleteFromOrganization(id string, orgid string) error {
	tx := repo.DB.MustBegin()
	userRoleDeleteSQL := fmt.Sprintf("DELETE FROM user_roles WHERE id = '%s' AND organization_id = '%s'", id, orgid)
	_ = tx.MustExec(userRoleDeleteSQL)
	err := tx.Commit()
	if err != nil {
		return err
	}
	return nil
}
