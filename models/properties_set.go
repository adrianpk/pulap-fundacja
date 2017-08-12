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
func (propertiesSet *PropertiesSet) MarshalJSON() ([]byte, error) {
	type Alias PropertiesSet
	return json.Marshal(&struct {
		*Alias
		AnniversaryDate int64 `json:"anniversaryDate"`
		StartedAt       int64 `json:"startedAt"`
		CreatedAt       int64 `json:"createdAt"`
		UpdatedAt       int64 `json:"updatedAt"`
	}{
		Alias:     (*Alias)(propertiesSet),
		StartedAt: propertiesSet.StartedAt.Time.Unix(),
		CreatedAt: propertiesSet.CreatedAt.Time.Unix(),
		UpdatedAt: propertiesSet.UpdatedAt.Time.Unix(),
	})
}

// UnmarshalJSON - Custom UnmarshalJSON function.
func (propertiesSet *PropertiesSet) UnmarshalJSON(data []byte) error {
	type Alias PropertiesSet
	aux := &struct {
		*Alias
		AnniversaryDate int64 `json:"anniversaryDate"`
		StartedAt       int64 `json:"startedAt"`
		CreatedAt       int64 `json:"createdAt"`
		UpdatedAt       int64 `json:"updatedAt"`
	}{
		Alias:     (*Alias)(propertiesSet),
		StartedAt: propertiesSet.StartedAt.Time.Unix(),
		CreatedAt: propertiesSet.CreatedAt.Time.Unix(),
		UpdatedAt: propertiesSet.UpdatedAt.Time.Unix(),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	ts := time.Unix(aux.StartedAt, 0)
	tc := time.Unix(aux.CreatedAt, 0)
	tu := time.Unix(aux.UpdatedAt, 0)
	propertiesSet.StartedAt = nulls.Time{Time: ts}
	propertiesSet.CreatedAt = nulls.Time{Time: tc}
	propertiesSet.UpdatedAt = nulls.Time{Time: tu}
	return nil
}
