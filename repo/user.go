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

	"golang.org/x/crypto/bcrypt"

	"github.com/adrianpk/fundacja/db"
	"github.com/adrianpk/fundacja/models"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Import pq without side effect
)

// UserRepository - User repository manager.
type UserRepository struct {
	DB *sqlx.DB
}

// MakeUserRepository - UserRepository constructor.
func MakeUserRepository() (UserRepository, error) {
	db, err := db.GetDbx()
	if err != nil {
		return UserRepository{}, err
	}
	return UserRepository{DB: db}, nil
}

// GetAll - GetAll Users in repo.
func (repo *UserRepository) GetAll() ([]models.User, error) {
	users := []models.User{}
	err := repo.DB.Select(&users, "SELECT * FROM users ORDER BY first_name ASC")
	return users, err
}

// Create - Persists a User in repo.
func (repo *UserRepository) Create(user *models.User) error {
	user.SetID()
	user.UpdatePasswordHash()
	user.SetCreationValues()
	tx := repo.DB.MustBegin()
	userInsertSQL := "INSERT INTO users (id, username, password_hash, email, first_name, middle_names, last_name, geolocation, started_at, created_by, is_active, is_logical_deleted, created_at, updated_at) VALUES (:id, :username, :password_hash, :email, :first_name, :middle_names, :last_name, :geolocation, :started_at, :created_by, :is_active, :is_logical_deleted, :created_at, :updated_at)"
	_, err := tx.NamedExec(userInsertSQL, user)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return err
}

// Login - Retrive a User if username/email and provided
func (repo *UserRepository) Login(user models.User) (models.User, error) {
	u := models.User{}
	err := repo.DB.Get(&u, "SELECT * FROM users WHERE username = $1 OR email=$2 LIMIT 1", user.Username, user.Email)
	if err != nil {
		return user, err
	}
	// Validate password
	err = bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(user.Password))
	if err != nil {
		return user, err
	}
	return u, nil
}

// Get - Retrive a Organization in repo by its ID.
func (repo *UserRepository) Get(id string) (models.User, error) {
	u := models.User{}
	err := repo.DB.Get(&u, "SELECT * FROM users WHERE id = $1", id)
	if err != nil {
		return u, err
	}
	return u, nil
}

// GetByUsername - Retrive a Organization in repo by its owner's username.
func (repo *UserRepository) GetByUsername(username string) (models.User, error) {
	u := models.User{}
	err := repo.DB.Get(&u, "SELECT * FROM users WHERE username = $1", username)
	if err != nil {
		return u, err
	}
	return u, nil
}

// Update - Update a user in repo.
func (repo *UserRepository) Update(user *models.User) error {
	// Update password and audit values
	user.SetUpdateValues()
	user.UpdatePasswordHash()
	// Current state
	reference, err := repo.Get(user.ID.String)
	if err != nil {
		return err
	}
	// Customized query
	changes := UserChanges(user, reference)
	number := len(changes)
	pos := 0
	last := number < 2
	var query bytes.Buffer
	query.WriteString("UPDATE users SET ")
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
	query.WriteString(fmt.Sprintf("WHERE id = '%s';", user.ID.String))
	tx := repo.DB.MustBegin()
	_, err = tx.NamedExec(query.String(), &user)
	if err != nil {
		return err
	}
	err = tx.Commit()
	return err
}

func (repo *UserRepository) Delete(id string) error {
	tx := repo.DB.MustBegin()
	userDeleteSQL := fmt.Sprintf("DELETE FROM users WHERE id = '%s'", id)
	_ = tx.MustExec(userDeleteSQL)
	err := tx.Commit()
	if err != nil {
		return err
	}
	return nil
}
