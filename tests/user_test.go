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

package tests

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"testing"

	"github.com/adrianpk/pulap/bootstrap"
	"github.com/adrianpk/pulap/logger"
	"github.com/adrianpk/pulap/repo"
	"github.com/adrianpk/pulap/testbootstrap"

	_ "github.com/lib/pq"
)

var (
	tbp           = testbootstrap.TestBootstrap
	usersURL      string
	signupURL     string
	loginURL      string
	user1         = "5958b185-8150-4aae-b53f-0c44771ddec5"
	user2         = "3c05e701-b495-4443-b454-2c37e2ecccdf"
	user1Username = "admin"
	user1Role     = "admin"
)

func init() {
	usersURL = fmt.Sprintf("%s/users", tbp.APIServerURL)
	signupURL = fmt.Sprintf("%s/signup", tbp.APIServerURL)
	loginURL = fmt.Sprintf("%s/login", tbp.APIServerURL)
	bootstrap.SetBootParameters(testbootstrap.BootParameters())
	bootstrap.Boot()
}

func TestMain(m *testing.M) {
	tbp.Start(m)
}

func TestSignup(t *testing.T) {
	logger.Debug("TestSignup...")
	tbp.PrepareTestDatabase()
	userJSON := `
	{
		"data": {
			"name": "Arthur",
			"username": "aquaman",
			"password": "sevenseas",
			"email": "arthurcurry@gmail.com"
		}
	}
	`
	tbp.Reader = strings.NewReader(userJSON)
	request, _ := http.NewRequest("POST", signupURL, tbp.Reader)
	res, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Fatal(err)
		t.Error(err)
	}
	if res.StatusCode != http.StatusCreated {
		t.Errorf("Status: %d | Expected: 201-StatusCreated", res.StatusCode)
	}
}

func TestSignupAlreadySignuped(t *testing.T) {
	logger.Debug("TestSignupAlreadySignuped...")
	tbp.PrepareTestDatabase()
	// userJSON := `{"data": {"name": "admin", "username": "Admin", "password": "password", "email": "admin@gmail.com"}}`
	userJSON := `
	{
		"data": {
			"name": "admin",
			"username": "Admin",
			"password": "password",
			"email": "admin@gmail.com"
		}
	}
	`
	tbp.Reader = strings.NewReader(userJSON)
	request, _ := http.NewRequest("POST", signupURL, tbp.Reader)
	res, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Fatal(err)
		t.Errorf("Error executing request: %s", err.Error())
		return
	}
	if res.StatusCode == http.StatusCreated {
		t.Errorf("Status: %d | Expected: 201-StatusCreated", res.StatusCode)
	}
}

func TestLoginByUsername(t *testing.T) {
	logger.Debug("TestLoginByUsername...")
	tbp.PrepareTestDatabase()
	userJSON := `
	{
		"data": {
			"username": "admin",
			"password": "darkknight"
		}
	}
	`
	tbp.Reader = strings.NewReader(userJSON)
	request, _ := http.NewRequest("POST", loginURL, tbp.Reader)
	res, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Fatal(err)
		t.Errorf("Error executing request: %s", err.Error())
		return
	}
	if res.StatusCode != http.StatusOK {
		t.Errorf("Status: %d | Expected: 200-StatusOk", res.StatusCode)
	}
}

func TestLoginByEmail(t *testing.T) {
	logger.Debug("TestLoginByEmail...")
	tbp.PrepareTestDatabase()
	userJSON := `
	{
		"data": {
			"email": "admin@gmail.com",
			"password": "darkknight"
		}
	}
	`
	tbp.Reader = strings.NewReader(userJSON)
	request, _ := http.NewRequest("POST", loginURL, tbp.Reader)
	res, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Fatal(err)
		t.Errorf("Error executing request: %s", err.Error())
		return
	}
	if res.StatusCode != http.StatusOK {
		t.Errorf("Status: %d | Expected: 200-StatusOk", res.StatusCode)
	}
}

func TestLoginBadCredentials(t *testing.T) {
	logger.Debug("TestLoginByEmail...")
	tbp.PrepareTestDatabase()
	userJSON := `
	{
		"data": {
			"username": "Flash",
			"email": "barryallen@gmail.com",
			"password": "speedy"
		}
	}
	`
	tbp.Reader = strings.NewReader(userJSON)
	request, _ := http.NewRequest("POST", loginURL, tbp.Reader)
	res, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Fatal(err)
		t.Errorf("Error executing request: %s", err.Error())
		return
	}
	if res.StatusCode == http.StatusOK {
		t.Errorf("Status: %d | Expected: 401-StatusUnauthorized", res.StatusCode)
	}
}

func TestGetUsers(t *testing.T) {
	logger.Debug("TestGetUsers...")
	tbp.PrepareTestDatabase()
	tbp.Reader = strings.NewReader("")
	request, _ := http.NewRequest("GET", usersURL, tbp.Reader)
	tbp.AuthorizeRequest(request, user1, user1Username, user1Role)
	res, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Fatal(err)
		t.Errorf("Error executing request: %s", err.Error())
		return
	}
	if res.StatusCode != http.StatusOK {
		t.Errorf("Status: %d | Expected: 200-StatusOk", res.StatusCode)
	}
}

func TestCreateUser(t *testing.T) {
	// Not implemented.
	if 1 == 200 {
		t.Errorf("Status: %d | Expected: 200-StatusOk", -1)
	}
}

func TestGet(t *testing.T) {
	logger.Debug("TestGet...")
	tbp.PrepareTestDatabase()
	tbp.Reader = strings.NewReader("")
	userURL := fmt.Sprintf("%s/%s", usersURL, user1)
	request, _ := http.NewRequest("GET", userURL, tbp.Reader)
	tbp.AuthorizeRequest(request, user1, user1Username, user1Role)
	res, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Fatal(err)
		t.Errorf("Error executing request: %s", err.Error())
		return
	}
	if res.StatusCode != http.StatusOK {
		t.Errorf("Status: %d | Expected: 200-StatusOk", res.StatusCode)
	}
}

func TestGetByUsername(t *testing.T) {
	logger.Debug("TestGetByUsername...")
	tbp.PrepareTestDatabase()
	tbp.Reader = strings.NewReader("")
	userURL := fmt.Sprintf("%s/%s", usersURL, "admin")
	request, _ := http.NewRequest("GET", userURL, tbp.Reader)
	tbp.AuthorizeRequest(request, user1, user1Username, user1Role)
	res, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Fatal(err)
		t.Errorf("Error executing request: %s", err.Error())
		return
	}
	if res.StatusCode != http.StatusOK {
		t.Errorf("Status: %d | Expected: 200-StatusOk", res.StatusCode)
	}
}

func TestUpdateUser(t *testing.T) {
	logger.Debug("TestUpdateUser...")
	tbp.PrepareTestDatabase()
	userJSON := `
	{
		"data": {
			"username": "administrator",
			"password": "password",
			"email": "administrator@gmail.com",
			"firstName": "Tim",
			"middleNames": "",
			"lastName": "Drake",
			"startedAt": 1735693261
		}
	}
	`
	tbp.Reader = strings.NewReader(userJSON)
	userURL := fmt.Sprintf("%s/%s", usersURL, user1)
	request, _ := http.NewRequest("PUT", userURL, tbp.Reader)
	tbp.AuthorizeRequest(request, user1, user1Username, user1Role)
	res, err := http.DefaultClient.Do(request)
	if err != nil {
		t.Errorf("Error executing request: %s", err.Error())
		return
	}
	if res.StatusCode != http.StatusNoContent {
		t.Errorf("Status: %d | Expected: 204-StatusNoContent", res.StatusCode)
	}
}

func TestUpdateUserWithDatabaseVerify(t *testing.T) {
	logger.Debug("TestUpdateUserWithDatabaseVerify...")
	tbp.PrepareTestDatabase()
	newUsername := "administrator"
	newEmail := "administrator@gmail.com"
	userJSON := fmt.Sprintf(`
	{
		"data": {
			"username": "%s",
			"password": "%s",
			"email": "administrator@gmail.com",
			"firstName": "Tim",
			"middleNames": "",
			"lastName": "Drake",
			"startedAt": 1735693261
		}
	}
	`, newUsername, newEmail)
	tbp.Reader = strings.NewReader(userJSON)
	userURL := fmt.Sprintf("%s/%s", usersURL, user1)
	request, _ := http.NewRequest("PUT", userURL, tbp.Reader)
	tbp.AuthorizeRequest(request, user1, user1Username, user1Role)
	res, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Fatal(err)
		return
	}
	if res.StatusCode == http.StatusNoContent {
		userRepo, err := repo.MakeUserRepository()
		if err != nil {
			log.Fatal(err)
			return
		}
		user, err := userRepo.Get(user1)
		if err == nil {
			if user.Username.String == newUsername && user.Email.String == newEmail {
				logger.Debug("User update: ok.")
			} else {
				error := fmt.Sprintf("Username: '%s' | Expected: '%s' - ", user.Username.String, newUsername)
				error += fmt.Sprintf("Email: '%s' | Expected: '%s'", user.Email.String, newEmail)
				t.Error(error)
			}
		} else {
			t.Error(err.Error())
		}
	} else {
		t.Errorf("Status: %d | Expected: 204-StatusNoContent", res.StatusCode)
	}
}

func TestDeleteUser(t *testing.T) {
	logger.Debug("TestDeleteUserWithDatabaseVerify...")
	tbp.PrepareTestDatabase()
	tbp.Reader = strings.NewReader("")
	userURL := fmt.Sprintf("%s/%s", usersURL, user1)
	request, _ := http.NewRequest("DELETE", userURL, tbp.Reader)
	tbp.AuthorizeRequest(request, user1, user1Username, user1Role)
	res, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Fatal(err)
		t.Errorf("Error executing request: %s", err.Error())
		return
	}
	if res.StatusCode != http.StatusNoContent {
		t.Errorf("Status: %d | Expected: 200-StatusNoContent", res.StatusCode)
	}
}

func TestDeleteUserWithDatabaseVerify(t *testing.T) {
	logger.Debug("TestDeleteUser...")
	tbp.PrepareTestDatabase()
	tbp.Reader = strings.NewReader("")
	userURL := fmt.Sprintf("%s/%s", usersURL, user1)
	request, _ := http.NewRequest("DELETE", userURL, tbp.Reader)
	tbp.AuthorizeRequest(request, user1, user1Username, user1Role)
	res, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Fatal(err)
		t.Errorf("Error executing request: %s", err.Error())
		return
	}
	if res.StatusCode == http.StatusNoContent {
		userRepo, err := repo.MakeUserRepository()
		if err != nil {
			log.Fatal(err)
			return
		}
		user, err := userRepo.Get(user1)
		if err != nil {
			logger.Debug("TestDeleteUser: ok")
		} else {
			t.Errorf("User: %s | Expected: 'nil'", user.Username)
		}
	} else {
		t.Errorf("Status: %d | Expected: 204-StatusNoContent", res.StatusCode)
	}
}
