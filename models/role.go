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

	"github.com/markbates/pop/nulls"
)

// MarshalJSON - Custom MarshalJSON function.
func (role *Role) MarshalJSON() ([]byte, error) {
	type Alias Role
	return json.Marshal(&struct {
		*Alias
		StartedAt int64 `json:"startedAt"`
		CreatedAt int64 `json:"createdAt"`
		UpdatedAt int64 `json:"updatedAt"`
	}{
		Alias:     (*Alias)(role),
		StartedAt: role.StartedAt.Time.Unix(),
		CreatedAt: role.CreatedAt.Time.Unix(),
		UpdatedAt: role.UpdatedAt.Time.Unix(),
	})
}

// UnmarshalJSON - Custom UnmarshalJSON function.
func (role *Role) UnmarshalJSON(data []byte) error {
	type Alias Role
	aux := &struct {
		*Alias
		StartedAt int64 `json:"startedAt"`
		CreatedAt int64 `json:"createdAt"`
		UpdatedAt int64 `json:"updatedAt"`
	}{
		Alias:     (*Alias)(role),
		StartedAt: role.StartedAt.Time.Unix(),
		CreatedAt: role.CreatedAt.Time.Unix(),
		UpdatedAt: role.UpdatedAt.Time.Unix(),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	ts := time.Unix(aux.StartedAt, 0)
	tc := time.Unix(aux.CreatedAt, 0)
	tu := time.Unix(aux.UpdatedAt, 0)
	role.StartedAt = nulls.Time{Time: ts}
	role.CreatedAt = nulls.Time{Time: tc}
	role.UpdatedAt = nulls.Time{Time: tu}
	return nil
}
