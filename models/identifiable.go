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
	"fmt"

	"github.com/markbates/pop/nulls"
	uuid "github.com/twinj/uuid"
)

// IdentifiableModel - Common properties for identifiable models.
type IdentifiableModel struct {
	ID          nulls.String `db:"id" json:"id"`
	URLID       string       `json:"urlid, omitempty"`
	Name        nulls.String `db:"name" json:"name"`
	Description nulls.String `db:"description" json:"description"`
}

// GetURLID - Get ID for use in URLS.
func (identifiable *IdentifiableModel) genURLID() {
	defer func() {
		if r := recover(); r != nil {
			identifiable.URLID = ""
		}
	}()
	identifiable.URLID = toURLID(identifiable.ID)
}

// SetID - Sets an ID od type UUID if its empty.
func (identifiable *IdentifiableModel) SetID() {
	identifiable.setUUID()
}

// SetUUID - Sets an UUID if its empty.
func (identifiable *IdentifiableModel) setUUID() {
	if identifiable.ID.String == "" {
		identifiable.ID = nulls.NewString(identifiable.generateUUIDString())
	}
}

func (identifiable *IdentifiableModel) generateUUID() uuid.UUID {
	//uuid.Init()
	return uuid.NewV4()
}

func (identifiable *IdentifiableModel) generateUUIDString() string {
	return fmt.Sprintf("%v", identifiable.generateUUID())
}

func toURLID(id nulls.String) string {
	parsed, _ := uuid.Parse(id.String)
	return uuid.Formatter(parsed, uuid.FormatHex)
}
