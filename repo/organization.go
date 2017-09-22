// Copyright (c) 2017 kuguar
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

// OrganizationRepository - Organization repository manager.
type OrganizationRepository struct {
	DB *sqlx.DB
}

// MakeOrganizationRepository - OrganizationRepository constructor.
func MakeOrganizationRepository() (OrganizationRepository, error) {
	db, err := db.GetDbx()
	if err != nil {
		return OrganizationRepository{}, err
	}
	return OrganizationRepository{DB: db}, nil
}

// GetAll - GetAll Organizations in repo.
func (repo *OrganizationRepository) GetAll() ([]models.Organization, error) {
	organizations := []models.Organization{}
	err := repo.DB.Select(&organizations, "SELECT * FROM organizations ORDER BY name ASC")
	return organizations, err
}

// Create - Persists a Organization in repo.
func (repo *OrganizationRepository) Create(organization *models.Organization) error {
	organization.SetID()
	organization.SetCreationValues()
	tx := repo.DB.MustBegin()
	organizationInsertSQL := "INSERT INTO organizations (id, name, description, user_username, user_id, geolocation, started_at, created_by, is_active, is_logical_deleted, created_at, updated_at) VALUES (:id, :name, :description, :user_username, :user_id, :geolocation, :started_at, :created_by, :is_active, :is_logical_deleted, :created_at, :updated_at)"
	_, err := tx.NamedExec(organizationInsertSQL, organization)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return err
}

// Get - Retrive a Organization in repo by its ID.
func (repo *OrganizationRepository) Get(id string) (models.Organization, error) {
	u := models.Organization{}
	err := repo.DB.Get(&u, "SELECT * FROM organizations WHERE id = $1", id)
	if err != nil {
		return u, err
	}
	return u, nil
}

// GetByName - Retrive a Organization in repo by its name.
func (repo *OrganizationRepository) GetByName(name string) (models.Organization, error) {
	u := models.Organization{}
	err := repo.DB.Get(&u, "SELECT * FROM organizations WHERE name = $1", name)
	if err != nil {
		return u, err
	}
	return u, nil
}

// Update - Update a organization in repo.
func (repo *OrganizationRepository) Update(organization *models.Organization) error {
	// Update password and audit values
	organization.SetUpdateValues()
	// Current state
	reference, err := repo.Get(organization.ID.String)
	if err != nil {
		return err
	}
	// Customized query
	changes := OrganizationChanges(organization, reference)
	number := len(changes)
	pos := 0
	last := number < 2
	var query bytes.Buffer
	query.WriteString("UPDATE organizations SET ")
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
	query.WriteString(fmt.Sprintf("WHERE id = '%s';", organization.ID.String))
	//logger.Debug(query.String())
	tx := repo.DB.MustBegin()
	_, err = tx.NamedExec(query.String(), organization)
	if err != nil {
		return err
	}
	err = tx.Commit()
	return err
}

// DeleteByIDString - Deletes organization from database.
func (repo *OrganizationRepository) Delete(id string) error {
	tx := repo.DB.MustBegin()
	organizationDeleteSQL := fmt.Sprintf("DELETE FROM organizations WHERE id = '%s'", id)
	_ = tx.MustExec(organizationDeleteSQL)
	err := tx.Commit()
	if err != nil {
		return err
	}
	return nil
}
