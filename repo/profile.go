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
	"io/ioutil"
	"os"
	"time"

	"github.com/adrianpk/pulap/bootstrap"
	"github.com/adrianpk/pulap/db"
	"github.com/adrianpk/pulap/logger"
	"github.com/adrianpk/pulap/models"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Import pq without side effects
)

// ProfileRepository - Profile repository manager.
type ProfileRepository struct {
	DB *sqlx.DB
}

// MakeProfileRepository - ProfileRepository constructor.
func MakeProfileRepository() (ProfileRepository, error) {
	db, err := db.GetDbx()
	if err != nil {
		return ProfileRepository{}, err
	}
	return ProfileRepository{DB: db}, nil
}

// GetAll - GetAll Profiles in repo.
func (repo *ProfileRepository) GetAll() ([]models.Profile, error) {
	profiles := []models.Profile{}
	err := repo.DB.Select(&profiles, "SELECT * FROM profiles ORDER BY first_name ASC")
	return profiles, err
}

// Create - Persists a Profile in repo.
func (repo *ProfileRepository) Create(profile *models.Profile) error {
	profile.SetID()
	profile.SetCreationValues()
	tx := repo.DB.MustBegin()
	profileInsertSQL := "INSERT INTO profiles (id, name, email, description, bio, moto, website, anniversary_date, avatar, avatar_uri, header_uri, geolocation, started_at, created_by, is_active, is_logical_deleted, created_at, updated_at, user_id) VALUES (:id, :name, :email, :description, :bio, :moto, :website, :anniversary_date, :avatar, :avatar_uri, :header_uri, :geolocation, :started_at, :created_by, :is_active, :is_logical_deleted, :created_at, :updated_at, :user_id)"
	_, err := tx.NamedExec(profileInsertSQL, profile)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return err
}

// GetByUserID - Retrive a Profile in repo by its user id.
func (repo *ProfileRepository) GetByUserID(id string) (models.Profile, error) {
	profile := models.Profile{}
	err := repo.DB.Get(&profile, "SELECT * FROM profiles WHERE user_id = $1", id)
	if err != nil {
		logger.Dump(err)
		return profile, err
	}
	profile.GenAvatarBase64()
	return profile, nil
}

// Get - Retrive a Profile in repo by its ID.
func (repo *ProfileRepository) Get(id string) (models.Profile, error) {
	profile := models.Profile{}
	err := repo.DB.Get(&profile, "SELECT * FROM profiles WHERE id = $1", id)
	if err != nil {
		return profile, err
	}
	profile.GenAvatarBase64()
	return profile, nil
}

// Update - Update a profile in repo.
func (repo *ProfileRepository) Update(profile *models.Profile) error {
	// Current state
	reference, err := repo.Get(profile.ID.String)
	if err != nil {
		return err
	}
	// Customized query
	changes := ProfileChanges(profile, reference)
	number := len(changes)
	pos := 0
	last := number < 2
	var query bytes.Buffer
	query.WriteString("UPDATE profiles SET ")
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
	query.WriteString(fmt.Sprintf("WHERE id = '%s';", profile.ID.String))
	// logger.Debug(query.String())
	tx := repo.DB.MustBegin()
	_, err = tx.NamedExec(query.String(), profile)
	if err != nil {
		return err
	}
	err = tx.Commit()
	return err
}

// Delete - Deletes profile from database.
func (repo *ProfileRepository) Delete(id string) error {
	tx := repo.DB.MustBegin()
	profileDeleteSQL := fmt.Sprintf("DELETE FROM profiles WHERE id = '%s'", id)
	_ = tx.MustExec(profileDeleteSQL)
	err := tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

// DeleteByUserID - Deletes users profile from database.
func (repo *ProfileRepository) DeleteByUserID(userID string) error {
	tx := repo.DB.MustBegin()
	profileDeleteSQL := fmt.Sprintf("DELETE FROM profiles WHERE user_id = '%s'", userID)
	_ = tx.MustExec(profileDeleteSQL)
	err := tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

// SaveProfileAvatar - Store profile's avatar.
func (repo *ProfileRepository) SaveProfileAvatar(profile *models.Profile) error {
	path := fmt.Sprintf("%s/public/users/%s/profile", bootstrap.BaseDir, profile.ID.String)
	assurePath(path)
	uri := fmt.Sprintf("%s/avatar_%d.png", path, timestamp())
	if png, ok := profile.AvatarAsPNG(); ok == nil {
		defaultFilePermission := os.FileMode(int(0777))
		err := ioutil.WriteFile(uri, png, defaultFilePermission)
		if err != nil {
			return err
		}
		profile.AvatarURI = models.ToNullsString(uri)
		repo.Update(profile)
	}
	return nil
}

// DeleteProfileAvatar - Delete profile's avatar.
func (repo *ProfileRepository) DeleteProfileAvatar(profile *models.Profile) error {
	err := os.Remove(profile.AvatarURI.String)
	if err != nil {
		logger.Errorf("Cannot delete %s", profile.AvatarURI.String)
		return err
	}
	profile.AvatarURI = models.ToNullsString("")
	err = repo.Update(profile)
	if err != nil {
		logger.Errorf("Error deleting profile's avatar for user %s", profile.UserID.String)
		return err
	}
	return nil
}

func storeProfileHeader(profile models.Profile) error {
	return nil
}

func assurePath(path string) error {
	defaultFilePermission := os.FileMode(int(0777))
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.MkdirAll(path, defaultFilePermission)
	}
	return nil
}

func timestamp() int64 {
	now := time.Now()
	nanos := now.UnixNano()
	millis := nanos / 1000000
	return millis
}
