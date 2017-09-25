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

// PlanSubscriptionRepository - PlanSubscription repository manager.
type PlanSubscriptionRepository struct {
	DB *sqlx.DB
}

// MakePlanSubscriptionRepository - PlanSubscriptionRepository constructor.
func MakePlanSubscriptionRepository() (PlanSubscriptionRepository, error) {
	db, err := db.GetDbx()
	if err != nil {
		return PlanSubscriptionRepository{}, err
	}
	return PlanSubscriptionRepository{DB: db}, nil
}

// GetAll - GetAll RolePermissions in repo.
func (repo *PlanSubscriptionRepository) GetAll(subscriptorID string) ([]models.PlanSubscription, error) {
	planSubscriptions := []models.PlanSubscription{}
	err := repo.DB.Select(&planSubscriptions, "SELECT * FROM plan_subscriptions WHERE user_id = $1 ORDER BY name ASC", subscriptorID)
	return planSubscriptions, err
}

func (repo *PlanSubscriptionRepository) Create(planSubscription *models.PlanSubscription) error {
	planSubscription.SetID()
	planSubscription.SetCreationValues()
	tx := repo.DB.MustBegin()
	planSubscriptionInsertSQL := "INSERT INTO plan_subscriptions (id, name, started_at, ends_at, organization_id, user_id, plan_id, created_by, is_active, is_logical_deleted, created_at, updated_at) VALUES (:id, :name, :started_at, :ends_at, :organization_id, :user_id, :plan_id, :created_by, :is_active, :is_logical_deleted, :created_at, :updated_at)"
	_, err := tx.NamedExec(planSubscriptionInsertSQL, planSubscription)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return err
}

// GetByUserID - Retrive a PlanSubscription in repo by its UserID.
func (repo *PlanSubscriptionRepository) GetByUserID(userID string) (models.PlanSubscription, error) {
	u := models.PlanSubscription{}
	err := repo.DB.Get(&u, "SELECT * FROM plan_subscriptions WHERE user_id = $1 LIMIT 1", userID)
	if err != nil {
		return u, err
	}
	return u, nil
}

// Get - Retrive a PlanSubscription in repo by its ID.
func (repo *PlanSubscriptionRepository) Get(id string) (models.PlanSubscription, error) {
	u := models.PlanSubscription{}
	err := repo.DB.Get(&u, "SELECT * FROM plan_subscriptions WHERE id = $1", id)
	if err != nil {
		return u, err
	}
	return u, nil
}

// GetByName - Retrive a PlanSubscription in repo by its planSubscriptionname.
func (repo *PlanSubscriptionRepository) GetByName(name string) (models.PlanSubscription, error) {
	u := models.PlanSubscription{}
	err := repo.DB.Get(&u, "SELECT * FROM plan_subscriptions WHERE name = $1", name)
	if err != nil {
		return u, err
	}
	return u, nil
}

// Update - Update a planSubscription in repo.
func (repo *PlanSubscriptionRepository) Update(planSubscription *models.PlanSubscription) error {
	// Update password and audit values
	planSubscription.SetUpdateValues()
	// Current state
	reference, err := repo.Get(planSubscription.ID.String)
	if err != nil {
		return err
	}
	// Customized query
	changes := PlanSubscriptionChanges(planSubscription, reference)
	number := len(changes)
	pos := 0
	last := number < 2
	var query bytes.Buffer
	query.WriteString("UPDATE plan_subscriptions SET ")
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
	query.WriteString(fmt.Sprintf("WHERE id = '%s';", planSubscription.ID.String))
	//logger.Debug(query.String())
	tx := repo.DB.MustBegin()
	_, err = tx.NamedExec(query.String(), planSubscription)
	if err != nil {
		return err
	}
	err = tx.Commit()
	return err
}

// Delete - Deletes planSubscription from database.
func (repo *PlanSubscriptionRepository) Delete(id string) error {
	tx := repo.DB.MustBegin()
	planSubscriptionDeleteSQL := fmt.Sprintf("DELETE FROM plan_subscriptions WHERE id = '%s'", id)
	_ = tx.MustExec(planSubscriptionDeleteSQL)
	err := tx.Commit()
	if err != nil {
		return err
	}
	return nil
}
