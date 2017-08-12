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

	"github.com/adrianpk/fundacja/bootstrap"
	"github.com/adrianpk/fundacja/logger"
	"github.com/adrianpk/fundacja/repo"
	"github.com/adrianpk/fundacja/testbootstrap"

	_ "github.com/lib/pq"
)

var (
	tbp           = testbootstrap.TestBootstrap
	resourcesURL  string
	organization1 = "d43809a2-5896-43c4-808e-549f2ee47783"
	organization2 = "b8cef4be-1ec3-44b4-9cbd-551f039f4fc7"
	resource1     = "656fe625-1f1e-4528-bec3-c11cf254cfe4"
	resource2     = "55a192d4-5d1a-4b1a-9fd0-f644394e9457"
	resource1Name = "Controllers1"
	user1         = "5958b185-8150-4aae-b53f-0c44771ddec5"
	user2         = "3c05e701-b495-4443-b454-2c37e2ecccdf"
)

func init() {
	resourcesURL = fmt.Sprintf("%s/organizations", tbp.APIServerURL)
	bootstrap.SetBootParameters(testbootstrap.BootParameters())
	bootstrap.Boot()
}

func TestMain(m *testing.M) {
	tbp.Start(m)
}

func undoResourceFixture() {
	resourceRepo, err := repo.MakeResourceRepository()
	if err != nil {
		log.Fatal(err)
	}
	err = resourceRepo.Delete(resource1)
	if err != nil {
		log.Fatal(err)
	}
}

func TestGetAllFromOrganization(t *testing.T) {
	logger.Debug("TestGetAll...")
	tbp.PrepareTestDatabase()
	tbp.Reader = strings.NewReader("")
	resourcesOrgURL := fmt.Sprintf("%s/%s/resources", resourcesURL, organization1)
	request, _ := http.NewRequest("GET", resourcesOrgURL, tbp.Reader)
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

func TestCreateResource(t *testing.T) {
	logger.Debug("TestCreateResource...")
	tbp.PrepareTestDatabase()
	undoResourceFixture()
	resourceJSON := fmt.Sprintf(`
	{
		"data": {
			"name": "Resource",
		  "description": "Resource description.",
			"organizationID": "%s"
		}
	}
	`, organization1)
	tbp.Reader = strings.NewReader(resourceJSON)
	resourcesOrgURL := fmt.Sprintf("%s/%s/resources", resourcesURL, organization1)
	request, _ := http.NewRequest("POST", resourcesOrgURL, tbp.Reader)
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

func TestGetFromOrganization(t *testing.T) {
	logger.Debug("TestGet...")
	tbp.PrepareTestDatabase()
	tbp.Reader = strings.NewReader("")
	resourceOrgURL := fmt.Sprintf("%s/%s/resources/%s", resourcesURL, organization1, resource1)
	request, _ := http.NewRequest("GET", resourceOrgURL, tbp.Reader)
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
	resourceOrgURL := fmt.Sprintf("%s/%s/resources/%s", resourcesURL, organization1, resource1Name)
	request, _ := http.NewRequest("GET", resourceOrgURL, tbp.Reader)
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

func TestUpdateResource(t *testing.T) {
	logger.Debug("TestUpdateResource...")
	tbp.PrepareTestDatabase()
	resourceJSON := fmt.Sprintf(`
	{
		"data": {
			"id": "%s",
			"name": "Resource new name",
		  "description": "Resource new description.",
			"organizationID": "%s"
		}
	}
	`, resource1, organization1)
	tbp.Reader = strings.NewReader(resourceJSON)
	resourceOrgURL := fmt.Sprintf("%s/%s/resources/%s", resourcesURL, organization1, resource1)
	request, _ := http.NewRequest("PUT", resourceOrgURL, tbp.Reader)
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

func TestUpdateResourceWithDatabaseVerify(t *testing.T) {
	logger.Debug("TestUpdateResourceWithDatabaseVerify...")
	tbp.PrepareTestDatabase()
	newName := "Resourcy1"
	newDescription := "Something new about resource number 1"
	resourceJSON := fmt.Sprintf(`
	{
		"data": {
			"id": "%s",
			"name": "%s",
		  "description": "%s",
			"organizationID": "%s"
		}
	}
	`, resource1, newName, newDescription, organization1)
	tbp.Reader = strings.NewReader(resourceJSON)
	resourceOrgURL := fmt.Sprintf("%s/%s/resources/%s", resourcesURL, organization1, resource1)
	request, _ := http.NewRequest("PUT", resourceOrgURL, tbp.Reader)
	tbp.AuthorizeRequest(request, user1, "admin", "admin")
	res, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Fatal(err)
		t.Errorf("Error executing request: %s", err.Error())
		return
	}
	if res.StatusCode == http.StatusNoContent {
		resourceRepo, err := repo.MakeResourceRepository()
		if err != nil {
			log.Fatal(err)
			return
		}
		resource, err := resourceRepo.Get(resource1)
		if err == nil {
			if resource.Name.String == newName && resource.Description.String == newDescription {
				logger.Debug("Resource update: ok.")
			} else {
				error := fmt.Sprintf("Name: '%s' | Expected: '%s' - ", resource.Name.String, newName)
				error += fmt.Sprintf("Description: '%s' | Expected: '%s'", resource.Description.String, newDescription)
				t.Error(error)
			}
		} else {
			t.Error(err.Error())
		}
	} else {
		t.Errorf("Status: %d | Expected: 204-StatusNoContent", res.StatusCode)
	}
}

func TestDeleteResource(t *testing.T) {
	logger.Debug("TestDeleteResource...")
	tbp.PrepareTestDatabase()
	tbp.Reader = strings.NewReader("")
	resourceOrgURL := fmt.Sprintf("%s/%s/resources/%s", resourcesURL, organization1, resource1)
	request, _ := http.NewRequest("DELETE", resourceOrgURL, tbp.Reader)
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

func TestDeleteResourceWithDatabaseVerify(t *testing.T) {
	logger.Debug("TestDeleteResourceWithDatabaseVerify...")
	tbp.PrepareTestDatabase()
	tbp.Reader = strings.NewReader("")
	resourceOrgURL := fmt.Sprintf("%s/%s/resources/%s", resourcesURL, organization1, resource1)
	request, _ := http.NewRequest("DELETE", resourceOrgURL, tbp.Reader)
	tbp.AuthorizeRequest(request, user1, "admin", "admin")
	res, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Fatal(err)
		t.Errorf("Error executing request: %s", err.Error())
		return
	}
	if res.StatusCode == http.StatusNoContent {
		resourceRepo, err := repo.MakeResourceRepository()
		if err != nil {
			log.Fatal(err)
			return
		}
		resource, err := resourceRepo.Get(resource1)
		if err != nil {
			logger.Debug("TestDeleteResource: ok")
		} else {
			t.Errorf("Resource: %s | Expected: 'nil'", resource.Name)
		}
	} else {
		t.Errorf("Status: %d | Expected: 204-StatusNoContent", res.StatusCode)
	}
}
