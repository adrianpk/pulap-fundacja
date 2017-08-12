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
	tbp              = testbootstrap.TestBootstrap
	organizationsURL string
	organization1    = "d43809a2-5896-43c4-808e-549f2ee47783"
	organization2    = "b8cef4be-1ec3-44b4-9cbd-551f039f4fc7"
	userRole1        = "c4bd3a52-e7c0-4219-80b2-d50556ea6b8e"
	userRole2        = "fd1a1483-8c41-4400-9e07-713f6a350c92"
	user1            = "5958b185-8150-4aae-b53f-0c44771ddec5"
	user2            = "3c05e701-b495-4443-b454-2c37e2ecccdf"
	role1            = "9b6869e4-f51a-4197-9608-f2898bd764d8"
	role2            = "1a40baee-e968-4fb0-8cc9-ecb62e2f2a76"
	userRole1Name    = "UserRole1"
)

func init() {
	organizationsURL = fmt.Sprintf("%s/organizations", tbp.APIServerURL)
	bootstrap.SetBootParameters(testbootstrap.BootParameters())
	bootstrap.Boot()
}

func TestMain(m *testing.M) {
	tbp.Start(m)
}

func undoUserRoleFixture() {
	userPermRepo, err := repo.MakeUserRoleRepository()
	if err != nil {
		log.Fatal(err)
	}
	err = userPermRepo.Delete(userRole1)
	if err != nil {
		log.Fatal(err)
	}
}

func TestGetAllFromOrganization(t *testing.T) {
	logger.Debug("TestGetAllFromOrganization...")
	tbp.PrepareTestDatabase()
	tbp.Reader = strings.NewReader("")
	userRolesURL := fmt.Sprintf("%s/%s/user-roles", organizationsURL, organization1)
	request, _ := http.NewRequest("GET", userRolesURL, tbp.Reader)
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

func TestCreateUserRoleInOrganization(t *testing.T) {
	logger.Debug("TestCreateuser...")
	tbp.PrepareTestDatabase()
	undoUserRoleFixture()
	roleJSON := fmt.Sprintf(`
	{
		"data": {
				"organizationID": "%s",
				"userID": "%s",
				"roleID": "%s"
		}
	}
	`, organization1, user1, role1)
	tbp.Reader = strings.NewReader(roleJSON)
	userRolesURL := fmt.Sprintf("%s/%s/user-roles", organizationsURL, organization1)
	request, _ := http.NewRequest("POST", userRolesURL, tbp.Reader)
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
	logger.Debug("TestGetFromOrganization...")
	tbp.PrepareTestDatabase()
	tbp.Reader = strings.NewReader("")
	userRoleURL := fmt.Sprintf("%s/%s/user-roles/%s", organizationsURL, organization1, userRole1)
	request, _ := http.NewRequest("GET", userRoleURL, tbp.Reader)
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

func TestUpdateUserRoleInOrganization(t *testing.T) {
	logger.Debug("TestUpdateUserRoleInOrganization...")
	tbp.PrepareTestDatabase()
	newDescription := "UserRole new description."
	roleJSON := fmt.Sprintf(`
	{
		"data": {
				"id": "%s",
				"description": "%s",
				"organizationID": "%s",
				"userID": "%s",
				"roleID": "%s"
		}
	}
	`, userRole1, newDescription, organization1, user1, role2)
	tbp.Reader = strings.NewReader(roleJSON)
	userRoleURL := fmt.Sprintf("%s/%s/user-roles/%s", organizationsURL, organization1, userRole1)
	request, _ := http.NewRequest("PUT", userRoleURL, tbp.Reader)
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

func TestUpdateUserRoleWithDatabaseVerify(t *testing.T) {
	logger.Debug("TestUpdateUserRoleWithDatabaseVerify...")
	tbp.PrepareTestDatabase()
	newDescription := "UserRole new description."
	roleJSON := fmt.Sprintf(`
	{
		"data": {
				"id": "%s",
				"description": "%s",
				"organizationID": "%s",
				"userID": "%s",
				"roleID": "%s"
		}
	}
	`, userRole1, newDescription, organization1, user1, role2)
	tbp.Reader = strings.NewReader(roleJSON)
	userRoleURL := fmt.Sprintf("%s/%s/user-roles/%s", organizationsURL, organization1, userRole1)
	request, _ := http.NewRequest("PUT", userRoleURL, tbp.Reader)
	tbp.AuthorizeRequest(request, user1, "admin", "admin")
	res, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Fatal(err)
		t.Errorf("Error executing request: %s", err.Error())
		return
	}
	if res.StatusCode == http.StatusNoContent {
		userPermRepo, err := repo.MakeUserRoleRepository()
		if err != nil {
			log.Fatal(err)
			return
		}
		userPerm, err := userPermRepo.Get(userRole1)
		if err == nil {
			if userPerm.Description.String == newDescription && userPerm.RoleID.String == role2 {
				logger.Debug("UserRole update: ok.")
			} else {
				error := fmt.Sprintf("Description: '%s' | Expected: '%s'", userPerm.Description.String, newDescription)
				error += fmt.Sprintf("Role ID: '%s' | Expected: '%s' - ", userPerm.RoleID.String, role2)
				t.Error(error)
			}
		} else {
			t.Error(err.Error())
		}
	} else {
		t.Errorf("Status: %d | Expected: 204-StatusNoContent", res.StatusCode)
	}
}

func TestDeleteUserRole(t *testing.T) {
	logger.Debug("TestDeleteUserRole...")
	tbp.PrepareTestDatabase()
	tbp.Reader = strings.NewReader("")
	userRoleURL := fmt.Sprintf("%s/%s/user-roles/%s", organizationsURL, organization1, userRole1)
	request, _ := http.NewRequest("DELETE", userRoleURL, tbp.Reader)
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

func TestDeleteUserRoleWithDatabaseVerify(t *testing.T) {
	logger.Debug("TestDeleteUserRoleWithDatabaseVerify...")
	tbp.PrepareTestDatabase()
	tbp.Reader = strings.NewReader("")
	userRoleURL := fmt.Sprintf("%s/%s/user-roles/%s", organizationsURL, organization1, userRole1)
	request, _ := http.NewRequest("DELETE", userRoleURL, tbp.Reader)
	tbp.AuthorizeRequest(request, user1, "admin", "admin")
	res, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Fatal(err)
		t.Errorf("Error executing request: %s", err.Error())
		return
	}
	if res.StatusCode == http.StatusNoContent {
		userPermRepo, err := repo.MakeUserRoleRepository()
		if err != nil {
			log.Fatal(err)
			return
		}
		userPerm, err := userPermRepo.Get(userRole1)
		if err != nil {
			logger.Debug("TestDeleteRole: ok")
		} else {
			t.Errorf("Role: %s | Expected: 'nil'", userPerm.Name)
		}
	} else {
		t.Errorf("Status: %d | Expected: 204-StatusNoContent", res.StatusCode)
	}
}
