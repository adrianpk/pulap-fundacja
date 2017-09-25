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
	user1            = "5958b185-8150-4aae-b53f-0c44771ddec5"
	user2            = "3c05e701-b495-4443-b454-2c37e2ecccdf"
	organizationsURL string
	organization1    = "d43809a2-5896-43c4-808e-549f2ee47783"
	organization2    = "b8cef4be-1ec3-44b4-9cbd-551f039f4fc7"
	role1            = "9b6869e4-f51a-4197-9608-f2898bd764d8"
	role2            = "1a40baee-e968-4fb0-8cc9-ecb62e2f2a76"
	role1Name        = "Role1"
)

func init() {
	organizationsURL = fmt.Sprintf("%s/organizations", tbp.APIServerURL)
	bootstrap.SetBootParameters(testbootstrap.BootParameters())
	bootstrap.Boot()
}

func TestMain(m *testing.M) {
	tbp.Start(m)
}

func undoRoleFixture() {
	roleRepo, err := repo.MakeRoleRepository()
	if err != nil {
		log.Fatal(err)
	}
	err = roleRepo.Delete(role1)
	if err != nil {
		log.Fatal(err)
	}
}

func TestGetAllFromOrganization(t *testing.T) {
	logger.Debug("TestGetAll...")
	tbp.PrepareTestDatabase()
	tbp.Reader = strings.NewReader("")
	rolesOrgURL := fmt.Sprintf("%s/%s/roles", organizationsURL, organization1)
	request, _ := http.NewRequest("GET", rolesOrgURL, tbp.Reader)
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

func TestCreateRole(t *testing.T) {
	logger.Debug("TestCreateRole...")
	tbp.PrepareTestDatabase()
	undoRoleFixture()
	roleJSON := fmt.Sprintf(`
	{
		"data": {
			"name": "Role",
		  "description": "Role description.",
			"organizationID": "%s"
		}
	}
	`, organization1)
	tbp.Reader = strings.NewReader(roleJSON)
	rolesURL := fmt.Sprintf("%s/%s/roles", organizationsURL, organization1)
	request, _ := http.NewRequest("POST", rolesURL, tbp.Reader)
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
	roleOrgURL := fmt.Sprintf("%s/%s/roles/%s", organizationsURL, organization1, role1)
	request, _ := http.NewRequest("GET", roleOrgURL, tbp.Reader)
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
	roleURL := fmt.Sprintf("%s/%s/roles/%s", organizationsURL, organization1, role1Name)
	request, _ := http.NewRequest("GET", roleURL, tbp.Reader)
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

func TestUpdateRole(t *testing.T) {
	logger.Debug("TestUpdateRole...")
	tbp.PrepareTestDatabase()
	roleJSON := fmt.Sprintf(`
	{
		"data": {
			"id": "%s",
			"name": "Role new name",
		  "description": "Role new description",
			"organizationID": "%s"
		}
	}
	`, role1, organization1)
	tbp.Reader = strings.NewReader(roleJSON)
	roleURL := fmt.Sprintf("%s/%s/roles/%s", organizationsURL, organization1, role1)
	request, _ := http.NewRequest("PUT", roleURL, tbp.Reader)
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

func TestUpdateRoleWithDatabaseVerify(t *testing.T) {
	logger.Debug("TestUpdateRoleWithDatabaseVerify...")
	tbp.PrepareTestDatabase()
	newName := "Role new name"
	newDescription := "Role new description"
	roleJSON := fmt.Sprintf(`
	{
		"data": {
			"id": "%s",
			"name": "%s",
		  "description": "%s",
			"organizationID": "%s"
		}
	}
	`, role1, newName, newDescription, organization1)
	tbp.Reader = strings.NewReader(roleJSON)
	roleURL := fmt.Sprintf("%s/%s/roles/%s", organizationsURL, organization1, role1)
	request, _ := http.NewRequest("PUT", roleURL, tbp.Reader)
	tbp.AuthorizeRequest(request, user1, "admin", "admin")
	res, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Fatal(err)
		t.Errorf("Error executing request: %s", err.Error())
		return
	}
	if res.StatusCode == http.StatusNoContent {
		roleRepo, err := repo.MakeRoleRepository()
		if err != nil {
			log.Fatal(err)
			return
		}
		role, err := roleRepo.Get(role1)
		if err == nil {
			if role.Name.String == newName && role.Description.String == newDescription {
				logger.Debug("Role update: ok.")
			} else {
				error := fmt.Sprintf("Name: '%s' | Expected: '%s' - ", role.Name.String, newName)
				error += fmt.Sprintf("Description: '%s' | Expected: '%s'", role.Description.String, newDescription)
				t.Error(error)
			}
		} else {
			t.Error(err.Error())
		}
	} else {
		t.Errorf("Status: %d | Expected: 204-StatusNoContent", res.StatusCode)
	}
}

func TestDeleteRole(t *testing.T) {
	logger.Debug("TestDeleteRole...")
	tbp.PrepareTestDatabase()
	tbp.Reader = strings.NewReader("")
	roleURL := fmt.Sprintf("%s/%s/roles/%s", organizationsURL, organization1, role1)
	request, _ := http.NewRequest("DELETE", roleURL, tbp.Reader)
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

func TestDeleteRoleWithDatabaseVerify(t *testing.T) {
	logger.Debug("TestDeleteRoleWithDatabaseVerify...")
	tbp.PrepareTestDatabase()
	tbp.Reader = strings.NewReader("")
	roleURL := fmt.Sprintf("%s/%s/roles/%s", organizationsURL, organization1, role1)
	request, _ := http.NewRequest("DELETE", roleURL, tbp.Reader)
	tbp.AuthorizeRequest(request, user1, "admin", "admin")
	res, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Fatal(err)
		t.Errorf("Error executing request: %s", err.Error())
		return
	}
	if res.StatusCode == http.StatusNoContent {
		roleRepo, err := repo.MakeRoleRepository()
		if err != nil {
			log.Fatal(err)
			return
		}
		role, err := roleRepo.Get(role1)
		if err != nil {
			logger.Debug("TestDeleteRole: ok")
		} else {
			t.Errorf("Role: %s | Expected: 'nil'", role.Name)
		}
	} else {
		t.Errorf("Status: %d | Expected: 204-StatusNoContent", res.StatusCode)
	}
}
