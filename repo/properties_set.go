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

	"github.com/adrianpk/pulap/db"
	"github.com/adrianpk/pulap/models"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Import pq without side effects
)

// PropertiesSetRepository - PropertiesSet repository manager.
type PropertiesSetRepository struct {
	DB *sqlx.DB
}

// MakePropertiesSetRepository - PropertiesSetRepository constructor.
func MakePropertiesSetRepository() (PropertiesSetRepository, error) {
	db, err := db.GetDbx()
	if err != nil {
		return PropertiesSetRepository{}, err
	}
	return PropertiesSetRepository{DB: db}, nil
}

// GetAll - GetAll RolePermissions in repo.
func (repo *PropertiesSetRepository) GetAll(holderID string) ([]models.PropertiesSet, error) {
	propertiesSets := []models.PropertiesSet{}
	err := repo.DB.Select(&propertiesSets, "SELECT * FROM properties_sets WHERE holder_id = $1 ORDER BY name ASC", holderID)
	return propertiesSets, err
}

// Create - Persists a PropertiesSet in repo.
func (repo *PropertiesSetRepository) Create(propertiesSet *models.PropertiesSet) error {
	propertiesSet.SetID()
	propertiesSet.SetCreationValues()
	tx := repo.DB.MustBegin()
	propertiesSetInsertSQL := "INSERT INTO properties_sets (id, name, description, geolocation, created_by, is_active, is_logical_deleted, created_at, updated_at, organization_id) VALUES (:id, :name, :description, :geolocation, :created_by, :is_active, :is_logical_deleted, :created_at, :updated_at, :organization_id)"
	_, err := tx.NamedExec(propertiesSetInsertSQL, propertiesSet)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return err
}

// Get - Retrive a PropertiesSet in repo by its ID.
func (repo *PropertiesSetRepository) Get(id string) (models.PropertiesSet, error) {
	u := models.PropertiesSet{}
	err := repo.DB.Get(&u, "SELECT * FROM properties_sets WHERE id = $1", id)
	if err != nil {
		return u, err
	}
	return u, nil
}

// GetByName - Retrive a PropertiesSet in repo by its propertiesSetname.
func (repo *PropertiesSetRepository) GetByName(name string) (models.PropertiesSet, error) {
	u := models.PropertiesSet{}
	err := repo.DB.Get(&u, "SELECT * FROM properties_sets WHERE name = $1", name)
	if err != nil {
		return u, err
	}
	return u, nil
}

// Update - Update a propertiesSet in repo.
func (repo *PropertiesSetRepository) Update(propertiesSet *models.PropertiesSet) error {
	// Update password and audit values
	propertiesSet.SetUpdateValues()
	// Current state
	reference, err := repo.Get(propertiesSet.ID.String)
	if err != nil {
		return err
	}
	// Customized query
	changes := PropertiesSetChanges(propertiesSet, reference)
	number := len(changes)
	pos := 0
	last := number < 2
	var query bytes.Buffer
	query.WriteString("UPDATE properties_sets SET ")
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
	query.WriteString(fmt.Sprintf("WHERE id = '%s';", propertiesSet.ID.String))
	//logger.Debug(query.String())
	tx := repo.DB.MustBegin()
	_, err = tx.NamedExec(query.String(), propertiesSet)
	if err != nil {
		return err
	}
	err = tx.Commit()
	return err
}

// Delete - Deletes propertiesSet from database.
func (repo *PropertiesSetRepository) Delete(id string) error {
	tx := repo.DB.MustBegin()
	propertiesSetDeleteSQL := fmt.Sprintf("DELETE FROM properties_sets WHERE id = '%s'", id)
	_ = tx.MustExec(propertiesSetDeleteSQL)
	err := tx.Commit()
	if err != nil {
		return err
	}
	return nil
}
