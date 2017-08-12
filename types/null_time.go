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

package types

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/lib/pq"
)

// NullTime - Nullable pq.Nulltime
type NullTime struct {
	pq.NullTime
}

// Scan implements the Scanner interface.
func (nt *NullTime) Scan(val interface{}) error {
	if val == nil {
		nt.NullTime, nt.Valid = pq.NullTime{}, false
		return nil
	}
	t := &NullTime{}
	err := t.Scan(val)
	if err != nil {
		nt.NullTime, nt.Valid = pq.NullTime{}, false
		return nil
	}
	nt.NullTime = t.NullTime
	nt.Valid = true
	return nil
}

// Value implements the driver Valuer interface.
func (nt NullTime) Value() (driver.Value, error) {
	if !nt.Valid {
		return pq.NullTime{}, nil
	}
	return nt.Time, nil
}

// MarshalJSON implements the driver json interface.
func (nt NullTime) MarshalJSON() ([]byte, error) {
	if nt.Valid {
		return json.Marshal(nt.NullTime)
	}
	return json.Marshal(nil)
}

// UnmarshalJSON implements the driver json interface.
func (nt *NullTime) UnmarshalJSON(data []byte) error {
	// Unmarshalling into a pointer will let us detect null
	var x *pq.NullTime
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}
	if x != nil {
		nt.Valid = true
		nt.NullTime = *x
	} else {
		nt.Valid = false
	}
	return nil
}

// ToTime - Return Value() time.Time value
func (nt *NullTime) ToTime() time.Time {
	value, err := nt.Value()
	if err != nil {
		return time.Time{}
	}
	return value.(time.Time)
}

// ToNullTime - Returns the NullTime corresponding to the value of the argument.
func ToNullTime(nt time.Time) pq.NullTime {
	return pq.NullTime{Time: nt, Valid: true}
}
