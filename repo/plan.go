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

// PlanRepository - Plan repository manager.
type PlanRepository struct {
	DB *sqlx.DB
}

// MakePlanRepository - PlanRepository constructor.
func MakePlanRepository() (PlanRepository, error) {
	db, err := db.GetDbx()
	if err != nil {
		return PlanRepository{}, err
	}
	return PlanRepository{DB: db}, nil
}

// GetAll - GetAll Plans in repo.
func (repo *PlanRepository) GetAll() ([]models.Plan, error) {
	plans := []models.Plan{}
	err := repo.DB.Select(&plans, "SELECT * FROM plans ORDER BY name ASC")
	return plans, err
}

// Create - Persists a Plan in repo.
func (repo *PlanRepository) Create(plan *models.Plan) error {
	plan.SetID()
	plan.SetCreationValues()
	tx := repo.DB.MustBegin()
	planInsertSQL := "INSERT INTO plans (id, name, description, started_at, ends_at, created_by, is_active, is_logical_deleted, created_at, updated_at) VALUES (:id, :name, :description, :started_at, :ends_at, :created_by, :is_active, :is_logical_deleted, :created_at, :updated_at)"
	_, err := tx.NamedExec(planInsertSQL, plan)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return err
}

// Get - Retrive a Plan in repo by its ID.
func (repo *PlanRepository) Get(id string) (models.Plan, error) {
	u := models.Plan{}
	err := repo.DB.Get(&u, "SELECT * FROM plans WHERE id = $1", id)
	if err != nil {
		return u, err
	}
	return u, nil
}

// GetByName - Retrive a Plan in repo by its planname.
func (repo *PlanRepository) GetByName(name string) (models.Plan, error) {
	u := models.Plan{}
	err := repo.DB.Get(&u, "SELECT * FROM plans WHERE name = $1", name)
	if err != nil {
		return u, err
	}
	return u, nil
}

// Update - Update a plan in repo.
func (repo *PlanRepository) Update(plan *models.Plan) error {
	// Update password and audit values
	plan.SetUpdateValues()
	// Current state
	reference, err := repo.Get(plan.ID.String)
	if err != nil {
		return err
	}
	// Customized query
	changes := PlanChanges(plan, reference)
	number := len(changes)
	pos := 0
	last := number < 2
	var query bytes.Buffer
	query.WriteString("UPDATE plans SET ")
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
	query.WriteString(fmt.Sprintf("WHERE id = '%s';", plan.ID.String))
	//logger.Debug(query.String())
	tx := repo.DB.MustBegin()
	_, err = tx.NamedExec(query.String(), plan)
	if err != nil {
		return err
	}
	err = tx.Commit()
	return err
}

// Delete - Deletes plan from database.
func (repo *PlanRepository) Delete(id string) error {
	tx := repo.DB.MustBegin()
	planDeleteSQL := fmt.Sprintf("DELETE FROM plans WHERE id = '%s'", id)
	_ = tx.MustExec(planDeleteSQL)
	err := tx.Commit()
	if err != nil {
		return err
	}
	return nil
}
