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
	tbp              = testbootstrap.TestBootstrap
	organizationsURL string
	user1            = "5958b185-8150-4aae-b53f-0c44771ddec5"
	user2            = "3c05e701-b495-4443-b454-2c37e2ecccdf"
	organization1    = "d43809a2-5896-43c4-808e-549f2ee47783"
	organization2    = "b8cef4be-1ec3-44b4-9cbd-551f039f4fc7"
)

func init() {
	organizationsURL = fmt.Sprintf("%s/organizations", tbp.APIServerURL)
	bootstrap.SetBootParameters(testbootstrap.BootParameters())
	bootstrap.Boot()
}

func TestMain(m *testing.M) {
	tbp.Start(m)
}

func undoOrganizationFixture() {
	organizationRepo, err := repo.MakeOrganizationRepository()
	if err != nil {
		log.Fatal(err)
	}
	err = organizationRepo.Delete(organization1)
	if err != nil {
		log.Fatal(err)
	}
}

func TestGetAll(t *testing.T) {
	logger.Debug("TestGetAll...")
	tbp.PrepareTestDatabase()
	tbp.Reader = strings.NewReader("")
	request, _ := http.NewRequest("GET", organizationsURL, tbp.Reader)
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

func TestCreateOrganization(t *testing.T) {
	logger.Debug("TestCreateOrganization...")
	tbp.PrepareTestDatabase()
	undoOrganizationFixture()
	organizationJSON := fmt.Sprintf(`
	{
		"data": {
			"name": "Organization",
		  "description": "Organization description",
			"userID": "%s"
		}
	}
	`, user1)
	tbp.Reader = strings.NewReader(organizationJSON)
	organizationURL := fmt.Sprintf("%s", organizationsURL)
	request, _ := http.NewRequest("POST", organizationURL, tbp.Reader)
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

func TestGet(t *testing.T) {
	logger.Debug("TestGet...")
	tbp.PrepareTestDatabase()
	tbp.Reader = strings.NewReader("")
	organizationURL := fmt.Sprintf("%s/%s", organizationsURL, organization1)
	request, _ := http.NewRequest("GET", organizationURL, tbp.Reader)
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

func TestGetByName(t *testing.T) {
	logger.Debug("TestGetByName...")
	tbp.PrepareTestDatabase()
	tbp.Reader = strings.NewReader("")
	organizationURL := fmt.Sprintf("%s/%s", organizationsURL, "Organization")
	request, _ := http.NewRequest("GET", organizationURL, tbp.Reader)
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

func TestUpdateOrganization(t *testing.T) {
	logger.Debug("TestUpdateOrganization...")
	tbp.PrepareTestDatabase()
	organizationJSON := fmt.Sprintf(`
	{
		"data": {
			"id": "%s",
			"name": "Organization new name",
		  "description": "Organization new description.",
			"userID": "%s"
		}
	}
	`, organization1, user1)
	tbp.Reader = strings.NewReader(organizationJSON)
	organizationURL := fmt.Sprintf("%s/%s", organizationsURL, organization1)
	request, _ := http.NewRequest("PUT", organizationURL, tbp.Reader)
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

func TestUpdateOrganizationWithDatabaseVerify(t *testing.T) {
	logger.Debug("TestUpdateUserWithDatabaseVerify...")
	tbp.PrepareTestDatabase()
	newName := "Organization new name"
	newDescription := "Organization new description."
	organizationJSON := fmt.Sprintf(`
	{
		"data": {
			"id": "%s",
			"name": "%s",
		  "description": "%s",
			"userID": "%s"
		}
	}
	`, organization1, newName, newDescription, user1)
	tbp.Reader = strings.NewReader(organizationJSON)
	organizationURL := fmt.Sprintf("%s/%s", organizationsURL, organization1)
	request, _ := http.NewRequest("PUT", organizationURL, tbp.Reader)
	tbp.AuthorizeRequest(request, user1, "admin", "admin")
	res, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Fatal(err)
		t.Errorf("Error executing request: %s", err.Error())
		return
	}
	if res.StatusCode == http.StatusNoContent {
		organizationRepo, err := repo.MakeOrganizationRepository()
		if err != nil {
			log.Fatal(err)
			return
		}
		organization, err := organizationRepo.Get(organization1)
		if err == nil {
			if organization.Name.String == newName && organization.Description.String == newDescription {
				logger.Debug("Organization update: ok.")
			} else {
				error := fmt.Sprintf("Name: '%s' | Expected: '%s' - ", organization.Name.String, newName)
				error += fmt.Sprintf("Description: '%s' | Expected: '%s'", organization.Description.String, newDescription)
				t.Error(error)
			}
		} else {
			t.Error(err.Error())
		}
	} else {
		t.Errorf("Status: %d | Expected: 204-StatusNoContent", res.StatusCode)
	}
}

func TestDeleteOrganization(t *testing.T) {
	logger.Debug("TestDeleteOrganization...")
	tbp.PrepareTestDatabase()
	tbp.Reader = strings.NewReader("")
	organizationURL := fmt.Sprintf("%s/%s", organizationsURL, organization1)
	request, _ := http.NewRequest("DELETE", organizationURL, tbp.Reader)
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

func TestDeleteOrganizationWithDatabaseVerify(t *testing.T) {
	logger.Debug("TestDeleteOrganizationWithDatabaseVerify...")
	tbp.PrepareTestDatabase()
	tbp.Reader = strings.NewReader("")
	organizationURL := fmt.Sprintf("%s/%s", organizationsURL, organization1)
	request, _ := http.NewRequest("DELETE", organizationURL, tbp.Reader)
	tbp.AuthorizeRequest(request, user1, "admin", "admin")
	res, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Fatal(err)
		t.Errorf("Error executing request: %s", err.Error())
		return
	}
	if res.StatusCode == http.StatusNoContent {
		organizationRepo, err := repo.MakeOrganizationRepository()
		if err != nil {
			log.Fatal(err)
			return
		}
		organization, err := organizationRepo.Get(organization1)
		if err != nil {
			logger.Debug("TestDeleteOrganization: ok")
		} else {
			t.Errorf("Organization: %s | Expected: 'nil'", organization.Name)
		}
	} else {
		t.Errorf("Status: %d | Expected: 204-StatusNoContent", res.StatusCode)
	}
}
