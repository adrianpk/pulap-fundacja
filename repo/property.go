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

// PropertyRepository - Property repository manager.
type PropertyRepository struct {
	DB *sqlx.DB
}

// MakePropertyRepository - PropertyRepository constructor.
func MakePropertyRepository() (PropertyRepository, error) {
	db, err := db.GetDbx()
	if err != nil {
		return PropertyRepository{}, err
	}
	return PropertyRepository{DB: db}, nil
}

// GetAll - GetAll Properties in repo.
func (repo *PropertyRepository) GetAll(propsetID string) ([]models.Property, error) {
	properties := []models.Property{}
	err := repo.DB.Select(&properties, "SELECT * FROM properties WHERE properties_set_id = $1 ORDER BY name ASC", propsetID)
	return properties, err
}

// Create - Persists a Property in repo.
func (repo *PropertyRepository) Create(property *models.Property) error {
	property.SetID()
	property.SetCreationValues()
	tx := repo.DB.MustBegin()
	propertyInsertSQL := "INSERT INTO properties (id, name, description, string_value, int_value, float_value, boolean_value, timestamp_value, geolocation_value, value_type, position, properties_set_id, created_by, is_active, is_logical_deleted, created_at, updated_at) VALUES (:id, :name, :description, :string_value, :int_value, :float_value, :boolean_value, :timestamp_value, :geolocation_value, :value_type, :position, :properties_set_id, :created_by, :is_active, :is_logical_deleted, :created_at, :updated_at)"
	_, err := tx.NamedExec(propertyInsertSQL, property)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return err
}

// Get - Retrive a Property in repo by its ID.
func (repo *PropertyRepository) Get(id string) (models.Property, error) {
	u := models.Property{}
	err := repo.DB.Get(&u, "SELECT * FROM properties WHERE id = $1", id)
	if err != nil {
		return u, err
	}
	return u, nil
}

// GetByNameInPropertiesSet - Retrive a Property in repo by its propertyname.
func (repo *PropertyRepository) GetByNameInPropertiesSet(name, propertiesSetID string) (models.Property, error) {
	u := models.Property{}
	err := repo.DB.Get(&u, "SELECT * FROM properties WHERE name = $1 and properties_set_id = $2", name, propertiesSetID)
	if err != nil {
		return u, err
	}
	return u, nil
}

// Update - Update a property in repo.
func (repo *PropertyRepository) Update(property *models.Property) error {
	// Update password and audit values
	property.SetUpdateValues()
	// Current state
	reference, err := repo.Get(property.ID.String)
	if err != nil {
		return err
	}
	// Customized query
	changes := PropertyChanges(property, reference)
	number := len(changes)
	pos := 0
	last := number < 2
	var query bytes.Buffer
	query.WriteString("UPDATE properties SET ")
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
	query.WriteString(fmt.Sprintf("WHERE id = '%s';", property.ID.String))
	// logger.Debug(query.String())
	tx := repo.DB.MustBegin()
	_, err = tx.NamedExec(query.String(), property)
	if err != nil {
		return err
	}
	err = tx.Commit()
	return err
}

// Delete - Deletes property from database.
func (repo *PropertyRepository) Delete(id string) error {
	tx := repo.DB.MustBegin()
	propertyDeleteSQL := fmt.Sprintf("DELETE FROM properties WHERE id = '%s'", id)
	_ = tx.MustExec(propertyDeleteSQL)
	err := tx.Commit()
	if err != nil {
		return err
	}
	return nil
}
