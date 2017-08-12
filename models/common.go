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
	"encoding/base64"
	"encoding/hex"
	"image"
	"image/png"
	"time"

	"github.com/markbates/pop/nulls"
)

// ValidateDate - Validates a NullTime date.
func ValidateDate(date *nulls.Time) {
	if true {
		date.Valid = true
	}
}

func base64Encode(bytes []byte) string {
	return base64.StdEncoding.EncodeToString(bytes)
}

func base64Decode(base64Data string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(base64Data)
}

func hexEncode(str string) nulls.ByteSlice {
	src := []byte(str)
	dst := make([]byte, hex.EncodedLen(len(src)))
	hex.Encode(dst, src)
	return nulls.NewByteSlice(dst)
}

func hexDecode(src []byte) string {
	dst := make([]byte, hex.DecodedLen(len(src)))
	hex.Decode(dst, src)
	return string(dst)
}

func pngEncode(image image.Image) ([]byte, error) {
	buffer := new(bytes.Buffer)
	err := png.Encode(buffer, image)
	if err != nil {
		return []byte{}, err
	}
	return buffer.Bytes(), nil
}

func ToNullsString(str string) nulls.String {
	return nulls.NewString(str)
}

func ToNullsInt64(i int64) nulls.Int64 {
	return nulls.NewInt64(i)
}

func ToNullsFoat64(f float64) nulls.Float64 {
	return nulls.NewFloat64(f)
}

func ToNullsBool(bln bool) nulls.Bool {
	return nulls.NewBool(bln)
}

func ToNullsTime(t time.Time) nulls.Time {
	return nulls.NewTime(t)
}

func NullsTrueBool() nulls.Bool {
	return ToNullsBool(true)
}

func NullsFalseBool() nulls.Bool {
	return ToNullsBool(false)
}

func NullsNowTime() nulls.Time {
	return ToNullsTime(time.Now())
}
