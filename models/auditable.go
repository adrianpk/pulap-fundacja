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
	"github.com/markbates/pop/nulls"
)

// AuditableModel - Common properties for auditable models.
type AuditableModel struct {
	StartedAt        nulls.Time   `db:"started_at" json:"startedAt, omitempty" schema:"-"`
	CreatedBy        nulls.String `db:"created_by" json:"createdBy, omitempty" schema:"-"`
	IsActive         nulls.Bool   `db:"is_active" json:"isActive, omitempty" schema:"is-active"`
	IsLogicalDeleted nulls.Bool   `db:"is_logical_deleted" json:"isLogicalDeleted, omitempty"`
	CreatedAt        nulls.Time   `db:"created_at" json:"createdAt, omitempty" schema:"-"`
	UpdatedAt        nulls.Time   `db:"updated_at" json:"updatedAt, omitempty" schema:"-"`
}

// SetCreationValues - Default values for models after creation.
func (auditable *AuditableModel) SetCreationValues() {
	now := NullsNowTime()
	auditable.StartedAt = now
	auditable.IsActive = NullsTrueBool()
	auditable.IsLogicalDeleted = NullsFalseBool()
	auditable.CreatedAt = now
	auditable.UpdatedAt = now
}

// SetUpdateValues - Default values for models after update.
func (auditable *AuditableModel) SetUpdateValues() {
	now := NullsNowTime()
	auditable.UpdatedAt = now
}

// // MarshalJSON - Custom MarshalJSON function.
// func (auditable *AuditableModel) MarshalJSON() ([]byte, error) {
// 	type Alias AuditableModel
// 	return json.Marshal(&struct {
// 		*Alias
// 		StartedAt int64 `json:"startedAt"`
// 		CreatedAt int64 `json:"createdAt"`
// 		UpdatedAt int64 `json:"updatedAt"`
// 	}{
// 		Alias:     (*Alias)(auditable),
// 		StartedAt: auditable.StartedAt.Time.Unix(),
// 		CreatedAt: auditable.CreatedAt.Time.Unix(),
// 		UpdatedAt: auditable.UpdatedAt.Time.Unix(),
// 	})
// }
//
// // UnmarshalJSON - Custom UnmarshalJSON function.
// func (auditable *AuditableModel) UnmarshalJSON(data []byte) error {
// 	type Alias AuditableModel
// 	aux := &struct {
// 		*Alias
// 		StartedAt int64 `json:"startedAt"`
// 		CreatedAt int64 `json:"createdAt"`
// 		UpdatedAt int64 `json:"updatedAt"`
// 	}{
// 		Alias:     (*Alias)(auditable),
// 		StartedAt: auditable.StartedAt.Time.Unix(),
// 		CreatedAt: auditable.CreatedAt.Time.Unix(),
// 		UpdatedAt: auditable.UpdatedAt.Time.Unix(),
// 	}
// 	if err := json.Unmarshal(data, &aux); err != nil {
// 		return err
// 	}
// 	ts := time.Unix(aux.StartedAt, 0)
// 	tc := time.Unix(aux.CreatedAt, 0)
// 	tu := time.Unix(aux.UpdatedAt, 0)
// 	auditable.StartedAt = nulls.Time{Time: ts}
// 	auditable.CreatedAt = nulls.Time{Time: tc}
// 	auditable.UpdatedAt = nulls.Time{Time: tu}
// 	return nil
// }
