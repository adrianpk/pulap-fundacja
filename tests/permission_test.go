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
	tbp             = testbootstrap.TestBootstrap
	permissionsURL  string
	organization1   = "d43809a2-5896-43c4-808e-549f2ee47783"
	organization2   = "b8cef4be-1ec3-44b4-9cbd-551f039f4fc7"
	permission1     = "cf903818-a2c5-46c2-8935-c4fc66fea60f"
	permission2     = "131a9338-39aa-4c33-8661-dcfa21ce726f"
	permission1Name = "Permission1"
	user1           = "5958b185-8150-4aae-b53f-0c44771ddec5"
	user2           = "3c05e701-b495-4443-b454-2c37e2ecccdf"
)

func init() {
	permissionsURL = fmt.Sprintf("%s/organizations", tbp.APIServerURL)
	bootstrap.SetBootParameters(testbootstrap.BootParameters())
	bootstrap.Boot()
}

func TestMain(m *testing.M) {
	tbp.Start(m)
}

func undoPermissionFixture() {
	permissionRepo, err := repo.MakePermissionRepository()
	if err != nil {
		log.Fatal(err)
	}
	err = permissionRepo.Delete(permission1)
	if err != nil {
		log.Fatal(err)
	}
}

func TestGetAllFromOrganization(t *testing.T) {
	logger.Debug("TestGetAll...")
	tbp.PrepareTestDatabase()
	tbp.Reader = strings.NewReader("")
	permissionsOrgURL := fmt.Sprintf("%s/%s/permissions", permissionsURL, organization1)
	request, _ := http.NewRequest("GET", permissionsOrgURL, tbp.Reader)
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

func TestCreatePermission(t *testing.T) {
	logger.Debug("TestCreatePermission...")
	tbp.PrepareTestDatabase()
	undoPermissionFixture()
	permissionJSON := fmt.Sprintf(`
	{
		"data": {
			"name": "Permission",
		  "description": "Permission description.",
			"organizationID": "%s"
		}
	}
	`, organization1)
	tbp.Reader = strings.NewReader(permissionJSON)
	permissionsOrgURL := fmt.Sprintf("%s/%s/permissions", permissionsURL, organization1)
	request, _ := http.NewRequest("POST", permissionsOrgURL, tbp.Reader)
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
	permissionOrgURL := fmt.Sprintf("%s/%s/permissions/%s", permissionsURL, organization1, permission1)
	request, _ := http.NewRequest("GET", permissionOrgURL, tbp.Reader)
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
	permissionsOrgURL := fmt.Sprintf("%s/%s/permissions/%s", permissionsURL, organization1, permission1Name)
	request, _ := http.NewRequest("GET", permissionsOrgURL, tbp.Reader)
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

func TestUpdatePermission(t *testing.T) {
	logger.Debug("TestUpdatePermission...")
	tbp.PrepareTestDatabase()
	permissionJSON := fmt.Sprintf(`
	{
		"data": {
			"id": "%s",
			"name": "Permission new name",
		  "description": "Permission new description.",
			"organizationID": "%s"
		}
	}
	`, permission1, organization1)
	tbp.Reader = strings.NewReader(permissionJSON)
	permissionOrgURL := fmt.Sprintf("%s/%s/permissions/%s", permissionsURL, organization1, permission1)
	request, _ := http.NewRequest("PUT", permissionOrgURL, tbp.Reader)
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

func TestUpdatePermissionWithDatabaseVerify(t *testing.T) {
	logger.Debug("TestUpdatePermissionWithDatabaseVerify...")
	tbp.PrepareTestDatabase()
	newName := "Permission new name"
	newDescription := "Permission new description."
	permissionJSON := fmt.Sprintf(`
	{
		"data": {
			"id": "%s",
			"name": "%s",
		  "description": "%s",
			"organizationID": "%s"
		}
	}
	`, permission1, newName, newDescription, organization1)
	tbp.Reader = strings.NewReader(permissionJSON)
	permissionOrgURL := fmt.Sprintf("%s/%s/permissions/%s", permissionsURL, organization1, permission1)
	request, _ := http.NewRequest("PUT", permissionOrgURL, tbp.Reader)
	tbp.AuthorizeRequest(request, user1, "admin", "admin")
	res, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Fatal(err)
		t.Errorf("Error executing request: %s", err.Error())
		return
	}
	if res.StatusCode == http.StatusNoContent {
		permissionRepo, err := repo.MakePermissionRepository()
		if err != nil {
			log.Fatal(err)
			return
		}
		permission, err := permissionRepo.Get(permission1)
		if err == nil {
			if permission.Name.String == newName && permission.Description.String == newDescription {
				logger.Debug("Permission update: ok.")
			} else {
				error := fmt.Sprintf("Name: '%s' | Expected: '%s' - ", permission.Name.String, newName)
				error += fmt.Sprintf("Description: '%s' | Expected: '%s'", permission.Description.String, newDescription)
				t.Error(error)
			}
		} else {
			t.Error(err.Error())
		}
	} else {
		t.Errorf("Status: %d | Expected: 204-StatusNoContent", res.StatusCode)
	}
}

func TestDeletePermission(t *testing.T) {
	logger.Debug("TestDeletePermission...")
	tbp.PrepareTestDatabase()
	tbp.Reader = strings.NewReader("")
	permissionOrgURL := fmt.Sprintf("%s/%s/permissions/%s", permissionsURL, organization1, permission1)
	request, _ := http.NewRequest("DELETE", permissionOrgURL, tbp.Reader)
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

func TestDeletePermissionWithDatabaseVerify(t *testing.T) {
	logger.Debug("TestDeletePermissionWithDatabaseVerify...")
	tbp.PrepareTestDatabase()
	tbp.Reader = strings.NewReader("")
	permissionOrgURL := fmt.Sprintf("%s/%s/permissions/%s", permissionsURL, organization1, permission1)
	request, _ := http.NewRequest("DELETE", permissionOrgURL, tbp.Reader)
	tbp.AuthorizeRequest(request, user1, "admin", "admin")
	res, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Fatal(err)
		t.Errorf("Error executing request: %s", err.Error())
		return
	}
	if res.StatusCode == http.StatusNoContent {
		permissionRepo, err := repo.MakePermissionRepository()
		if err != nil {
			log.Fatal(err)
			return
		}
		permission, err := permissionRepo.Get(permission1)
		if err != nil {
			logger.Debug("TestDeletePermission: ok")
		} else {
			t.Errorf("Permission: %s | Expected: 'nil'", permission.Name)
		}
	} else {
		t.Errorf("Status: %d | Expected: 204-StatusNoContent", res.StatusCode)
	}
}
