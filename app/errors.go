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

package app

import (
	"errors"
)

var (
	// ErrDbConnection - Cannot connect to datastore.
	ErrDbConnection = errors.New("Cannot connect to datastore")
	// ErrTokenExpired - Access Token is expired.
	ErrTokenExpired = errors.New("Access Token is expired")
	// ErrTokenParsing - Error while parsing the Access Token
	ErrTokenParsing = errors.New("Error while parsing the Access Token")
	// ErrTokenInvalid - Invalid Access Token
	ErrTokenInvalid = errors.New("Invalid Access Token")
	// ErrEntityAlreadySignedUp - Email or Username is already signed up.
	ErrEntityAlreadySignedUp = errors.New("Email or Username is already signed up")
	// ErrRequest - Bad request.
	ErrRequest = errors.New("Bad request")
	// ErrRequestParsing - Error parsing request data.
	ErrRequestParsing = errors.New("Error parsing request data")
	// ErrImageDecoding - Error decoding image data.
	ErrImageDecoding = errors.New("Error decoding image data")
	// ErrResponseMarshalling - Error marshalling response data.
	ErrResponseMarshalling = errors.New("Error marshalling response data")
	// ErrRegistration - Registration error.
	ErrRegistration = errors.New("Registration error")
	// ErrLogin - Login error.
	ErrLogin = errors.New("Login error")
	// ErrLoginInvalidData - Invalid login data.
	ErrLoginInvalidData = errors.New("Invalid login data")
	// ErrLoginDenied - Login denied.
	ErrLoginDenied = errors.New("Login denied")
	// ErrLoginTokenCreate - Error while generating the access token
	ErrLoginTokenCreate = errors.New("Error while generating the access token")
	// ErrNotLoggedIn - Not logged in.
	ErrNotLoggedIn = errors.New("Not logged in")
	// ErrUnauthorized - Unauthorized
	ErrUnauthorized = errors.New("Unauthorized")
	// ErrOwnerOnlyCanManage - Error while generating the access token
	ErrOwnerOnlyCanManage = errors.New("Only entity owner are allowed to manage entity")
	// ErrEntityInvalidData - Entity invalid data.
	ErrEntityInvalidData = errors.New("Entity invalid data")
	// ErrEntityNotFound - Entity not found.
	ErrEntityNotFound = errors.New("Entity not found")
	// ErrEntitySelect - Cannot select user.
	ErrEntitySelect = errors.New("Cannot select user/s")
	// ErrEntityCreate - Cannot create user.
	ErrEntityCreate = errors.New("Cannot create entity")
	// ErrEntityUpdate - Cannot update model.
	ErrEntityUpdate = errors.New("Cannot update entity")
	// ErrEntityDelete - Cannot delete entity.
	ErrEntityDelete = errors.New("Cannot delete entity")
	// ErrEntitySetProperty - Cannot set property.
	ErrEntitySetProperty = errors.New("Cannot set property")
	// ErrImageProcessing - Error processing image.
	ErrImageProcessing = errors.New("Error processing image")
)
