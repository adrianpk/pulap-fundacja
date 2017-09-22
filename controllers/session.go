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
//

package controllers

import (
	"net/http"

	"github.com/adrianpk/fundacja/logger"
	"github.com/adrianpk/fundacja/models"
	"github.com/gorilla/sessions"
)

const (
	sessionSecret = "fundacja-secret"
	sessionName   = "fundacja-session"
	cookieName    = "fundacja-session"
)

var (
	sessionStore *sessions.CookieStore
)

func init() {
	sessionStore = sessions.NewCookieStore([]byte(sessionSecret))
}

func setSession(w http.ResponseWriter, r *http.Request, user models.User, remember bool) error {
	// Gen and configure session
	session, err := sessionStore.Get(r, sessionName)
	if err != nil {
		return err
	}
	configureSession(session, remember)
	// Set values
	session.Values["userID"] = user.ID.String
	err = session.Save(r, w)
	if err != nil {
		return err
	}
	return nil
}

func clearSession(w http.ResponseWriter, r *http.Request) error {
	session, err := sessionStore.Get(r, sessionName)
	if err != nil {
		return err
	}
	// Invalidate session values
	session.Values["userID"] = nil
	// Save it before we write to the response/return from the handler.
	err = session.Save(r, w)
	if err != nil {
		return err
	}
	return nil
}

func configureSession(session *sessions.Session, remember bool) {
	mins := 20
	if remember {
		logger.Debug("Remember is true")
		mins = 10080 // 1 Week
	}
	session.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   mins * 60,
		HttpOnly: true,
	}
}
