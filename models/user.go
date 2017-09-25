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

package models

import (
	"encoding/json"
	"time"

	"github.com/adrianpk/pulap/logger"
	"github.com/markbates/pop/nulls"
	"golang.org/x/crypto/bcrypt"
)

// UpdatePasswordHash - Creates a password hash from current password.
func (user *User) UpdatePasswordHash() {
	if user.Password == "" {
		logger.Debug("New password not provided.")
		return
	}
	hpass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Fatal(err)
	}
	user.PasswordHash = string(hpass)
}

// ClearPassword - Clear password related fields.
func (user *User) ClearPassword() {
	user.Password = ""
	user.PasswordHash = ""
}

// MarshalJSON - Custom MarshalJSON function.
func (user *User) MarshalJSON() ([]byte, error) {
	type Alias User
	return json.Marshal(&struct {
		*Alias
		StartedAt int64 `json:"startedAt"`
		CreatedAt int64 `json:"createdAt"`
		UpdatedAt int64 `json:"updatedAt"`
	}{
		Alias:     (*Alias)(user),
		StartedAt: user.StartedAt.Time.Unix(),
		CreatedAt: user.CreatedAt.Time.Unix(),
		UpdatedAt: user.UpdatedAt.Time.Unix(),
	})
}

// UnmarshalJSON - Custom UnmarshalJSON function.
func (user *User) UnmarshalJSON(data []byte) error {
	type Alias User
	aux := &struct {
		*Alias
		StartedAt int64 `json:"startedAt"`
		CreatedAt int64 `json:"createdAt"`
		UpdatedAt int64 `json:"updatedAt"`
	}{
		Alias:     (*Alias)(user),
		StartedAt: user.StartedAt.Time.Unix(),
		CreatedAt: user.CreatedAt.Time.Unix(),
		UpdatedAt: user.UpdatedAt.Time.Unix(),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	ts := time.Unix(aux.StartedAt, 0)
	tc := time.Unix(aux.CreatedAt, 0)
	tu := time.Unix(aux.UpdatedAt, 0)
	user.StartedAt = nulls.Time{Time: ts}
	user.CreatedAt = nulls.Time{Time: tc}
	user.UpdatedAt = nulls.Time{Time: tu}
	return nil
}
