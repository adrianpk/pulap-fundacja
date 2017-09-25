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
	tbp      = testbootstrap.TestBootstrap
	usersURL string
	user1    = "5958b185-8150-4aae-b53f-0c44771ddec5"
	user2    = "3c05e701-b495-4443-b454-2c37e2ecccdf"
	profile1 = "28bb0dad-ece8-44a1-8c45-c4898968bee5"
	profile2 = "1b57cb73-7f61-4323-ae87-86b4d0569178"
)

func init() {
	usersURL = fmt.Sprintf("%s/users", tbp.APIServerURL)
	bootstrap.SetBootParameters(testbootstrap.BootParameters())
	bootstrap.Boot()
}

func TestMain(m *testing.M) {
	tbp.Start(m)
}

func undoProfileFixture() {
	profileRepo, err := repo.MakeProfileRepository()
	if err != nil {
		log.Fatal(err)
	}
	err = profileRepo.Delete(user1)
	if err != nil {
		log.Fatal(err)
	}
}

func TestCreateProfile(t *testing.T) {
	logger.Debug("TestCreateProfile...")
	tbp.PrepareTestDatabase()
	undoProfileFixture()
	profileJSON := `
	{
		"data": {
			"name": "Administrator",
			"email": "admin@gmail.com",
		  "description": "Administrator's profile",
		  "bio": "Administrator's bio",
		  "moto": "Administrator's moto",
		  "website": "admin.com",
		  "anniversaryDate": 1420070400,
			"userID": "5958b185-8150-4aae-b53f-0c44771ddec5"
		}
	}
	`
	tbp.Reader = strings.NewReader(profileJSON)
	profileURL := fmt.Sprintf("%s/%s/profile", usersURL, user1)
	request, _ := http.NewRequest("POST", profileURL, tbp.Reader)
	tbp.AuthorizeRequest(request, user1, "admin", "admin")
	res, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Fatal(err)
		t.Error(err)
	}
	if res.StatusCode != http.StatusCreated {
		t.Errorf("Status: %d | Expected: 201-StatusCreated", res.StatusCode)
	}
}

func TestGetProfileByUserID(t *testing.T) {
	logger.Debug("TestGetProfileByUserID...")
	tbp.PrepareTestDatabase()
	tbp.Reader = strings.NewReader("")
	profileURL := fmt.Sprintf("%s/%s/profile", usersURL, user1)
	request, _ := http.NewRequest("GET", profileURL, tbp.Reader)
	tbp.AuthorizeRequest(request, user1, "admin", "admin")
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

func TestGetProfileByUsername(t *testing.T) {
	logger.Debug("TestGetByUsername...")
	tbp.PrepareTestDatabase()
	tbp.Reader = strings.NewReader("")
	profileURL := fmt.Sprintf("%s/%s/profile", usersURL, "admin")
	request, _ := http.NewRequest("GET", profileURL, tbp.Reader)
	tbp.AuthorizeRequest(request, user1, "admin", "admin")
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

func TestDoUpdateProfile(t *testing.T) {
	logger.Debug("TestUpdateProfile...")
	tbp.PrepareTestDatabase()
	newName := "Admin"
	newEmail := "admin@gmail.com"
	profileJSON := fmt.Sprintf(`
	{
		"data": {
			"id": "%s",
			"name": "%s",
			"email": "%s",
		  "description": "Administrator profile updated",
		  "bio": "Administrator bio updated",
		  "moto": "Administrator moto updated",
		  "website": "admin.com",
		  "anniversaryDate": 1420070400,
			"userID": "5958b185-8150-4aae-b53f-0c44771ddec5"
		}
	}
	`, profile1, newName, newEmail)
	tbp.Reader = strings.NewReader(profileJSON)
	profileURL := fmt.Sprintf("%s/%s/profile", usersURL, user1)
	request, _ := http.NewRequest("PUT", profileURL, tbp.Reader)
	tbp.AuthorizeRequest(request, user1, "admin", "admin")
	res, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Fatal(err)
		t.Errorf("Error executing request: %s", err.Error())
		return
	}
	if res.StatusCode != http.StatusNoContent {
		t.Errorf("Status: %d | Expected: 204-StatusNoContent", res.StatusCode)
	}
}

func TestUpdateProfileWithDatabaseVerify(t *testing.T) {
	logger.Debug("TestUpdateUserWithDatabaseVerify...")
	tbp.PrepareTestDatabase()
	newName := "Admin"
	newEmail := "admin@gmail.com"
	profileJSON := fmt.Sprintf(`
	{
		"data": {
			"id": "%s",
			"name": "%s",
			"email": "%s",
		  "description": "Administrator profile updated",
		  "bio": "Administrator bio updated",
		  "moto": "Administrator moto updated",
		  "website": "admin.com",
		  "anniversaryDate": 1420070400,
			"userID": "5958b185-8150-4aae-b53f-0c44771ddec5"
		}
	}
	`, profile1, newName, newEmail)
	tbp.Reader = strings.NewReader(profileJSON)
	profileURL := fmt.Sprintf("%s/%s/profile", usersURL, user1)
	request, _ := http.NewRequest("PUT", profileURL, tbp.Reader)
	tbp.AuthorizeRequest(request, user1, "admin", "admin")
	res, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Fatal(err)
		t.Errorf("Error executing request: %s", err.Error())
		return
	}
	if res.StatusCode == http.StatusNoContent {
		profileRepo, err := repo.MakeProfileRepository()
		if err != nil {
			log.Fatal(err)
			return
		}
		profile, err := profileRepo.GetByUserID(user1)
		if err == nil {
			if profile.Name.String == newName && profile.Email.String == newEmail {
				logger.Debug("Profile update: ok.")
			} else {
				error := fmt.Sprintf("Name: '%s' | Expected: '%s' - ", profile.Name.String, newName)
				error += fmt.Sprintf("Email: '%s' | Expected: '%s'", profile.Email, newEmail)
				t.Error(error)
			}
		} else {
			t.Error(err.Error())
		}
	} else {
		t.Errorf("Status: %d | Expected: 204-StatusNoContent", res.StatusCode)
	}
}

func TestDeleteProfile(t *testing.T) {
	logger.Debug("TestDeleteProfile...")
	tbp.PrepareTestDatabase()
	tbp.Reader = strings.NewReader("")
	profileURL := fmt.Sprintf("%s/%s/profile", usersURL, user1)
	request, _ := http.NewRequest("DELETE", profileURL, tbp.Reader)
	tbp.AuthorizeRequest(request, user1, "admin", "admin")
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

func TestDeleteProfileWithDatabaseVerify(t *testing.T) {
	logger.Debug("TestDeleteProfileWithDatabaseVerify...")
	tbp.PrepareTestDatabase()
	tbp.Reader = strings.NewReader("")
	profileURL := fmt.Sprintf("%s/%s/profile", usersURL, user1)
	request, _ := http.NewRequest("DELETE", profileURL, tbp.Reader)
	tbp.AuthorizeRequest(request, user1, "admin", "admin")
	res, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Fatal(err)
		t.Errorf("Error executing request: %s", err.Error())
		return
	}
	if res.StatusCode == http.StatusNoContent {
		profileRepo, err := repo.MakeProfileRepository()
		if err != nil {
			log.Fatal(err)
			return
		}
		profile, err := profileRepo.GetByUserID(user1)
		if err != nil {
			logger.Debug("TestDeleteProfile: ok")
		} else {
			t.Errorf("Profile: %s | Expected: 'nil'", profile.Name)
		}
	} else {
		t.Errorf("Status: %d | Expected: 204-StatusNoContent", res.StatusCode)
	}
}
