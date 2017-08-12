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
	"bytes"
	"encoding/json"
	"image"
	"time"

	"github.com/markbates/pop/nulls"
)

// SetOwner - Set user.
func (profile *Profile) SetOwner(user User) {
	profile.UserID = user.ID
}

// SetAvatarAsBase64 - Update avatar's Base64 and binary representation using a Base64 encoded string.
func (profile *Profile) SetAvatarAsBase64(base64 string) {
	profile.AvatarBase64 = base64
	profile.GenAvatarHex()
}

// SetAvatarAsHex - Update avatar's Base64 and binary representation using a Base64 encoded string.
func (profile *Profile) SetAvatarAsHex(hexa nulls.ByteSlice) {
	profile.Avatar = hexa
	profile.GenAvatarBase64()
}

// GenAvatarHex - Update avatar's Base64 representation corresponding to its binary form.
func (profile *Profile) GenAvatarHex() {
	size := len(profile.AvatarBase64)
	if size > 0 {
		profile.Avatar = hexEncode(profile.AvatarBase64)
	}
}

// GenAvatarBase64 - Update avatar's Base64 representation corresponding to its binary form.
func (profile *Profile) GenAvatarBase64() {
	size := len(profile.Avatar.ByteSlice)
	if size > 0 {
		profile.AvatarBase64 = hexDecode(profile.Avatar.ByteSlice)
	}
}

// AvatarAsPNG - Get a PNG representation of avatar.
func (profile *Profile) AvatarAsPNG() (png []byte, err error) {
	imgBytes, err := base64Decode(profile.AvatarBase64)
	if err != nil {
		return png, err
	}
	img, _, err := image.Decode(bytes.NewReader(imgBytes))
	if err != nil {
		return png, err
	}
	png, err = pngEncode(img)
	if err != nil {
		return png, err
	}
	return png, nil
}

// MarshalJSON - Custom MarshalJSON function.
func (profile *Profile) MarshalJSON() ([]byte, error) {
	type Alias Profile
	return json.Marshal(&struct {
		*Alias
		AnniversaryDate int64 `json:"anniversaryDate"`
		StartedAt       int64 `json:"startedAt"`
		CreatedAt       int64 `json:"createdAt"`
		UpdatedAt       int64 `json:"updatedAt"`
	}{
		Alias:           (*Alias)(profile),
		AnniversaryDate: profile.AnniversaryDate.Time.Unix(),
		StartedAt:       profile.StartedAt.Time.Unix(),
		CreatedAt:       profile.CreatedAt.Time.Unix(),
		UpdatedAt:       profile.UpdatedAt.Time.Unix(),
	})
}

// UnmarshalJSON - Custom UnmarshalJSON function.
func (profile *Profile) UnmarshalJSON(data []byte) error {
	type Alias Profile
	aux := &struct {
		*Alias
		AnniversaryDate int64 `json:"anniversaryDate"`
		StartedAt       int64 `json:"startedAt"`
		CreatedAt       int64 `json:"createdAt"`
		UpdatedAt       int64 `json:"updatedAt"`
	}{
		Alias:           (*Alias)(profile),
		AnniversaryDate: profile.AnniversaryDate.Time.Unix(),
		StartedAt:       profile.StartedAt.Time.Unix(),
		CreatedAt:       profile.CreatedAt.Time.Unix(),
		UpdatedAt:       profile.UpdatedAt.Time.Unix(),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	ta := time.Unix(aux.AnniversaryDate, 0)
	ts := time.Unix(aux.StartedAt, 0)
	tc := time.Unix(aux.CreatedAt, 0)
	tu := time.Unix(aux.UpdatedAt, 0)
	profile.AnniversaryDate = nulls.Time{Time: ta}
	profile.StartedAt = nulls.Time{Time: ts}
	profile.CreatedAt = nulls.Time{Time: tc}
	profile.UpdatedAt = nulls.Time{Time: tu}
	return nil
}
