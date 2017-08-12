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
	tbp                 = testbootstrap.TestBootstrap
	rolePermissionsURL  string
	user1               = "5958b185-8150-4aae-b53f-0c44771ddec5"
	user2               = "3c05e701-b495-4443-b454-2c37e2ecccdf"
	organization1       = "d43809a2-5896-43c4-808e-549f2ee47783"
	organization2       = "b8cef4be-1ec3-44b4-9cbd-551f039f4fc7"
	rolePermission1     = "092ace1d-6c78-499a-a5e9-125e5f65ce2a"
	rolePermission2     = "f25a4af8-661f-4197-99b8-cdee96bc4e6f"
	role1               = "9b6869e4-f51a-4197-9608-f2898bd764d8"
	role2               = "1a40baee-e968-4fb0-8cc9-ecb62e2f2a76"
	permission1         = "cf903818-a2c5-46c2-8935-c4fc66fea60f"
	permission2         = "131a9338-39aa-4c33-8661-dcfa21ce726f"
	rolePermission1Name = "RolePermission1"
)

func init() {
	rolePermissionsURL = fmt.Sprintf("%s/organizations", tbp.APIServerURL)
	bootstrap.SetBootParameters(testbootstrap.BootParameters())
	bootstrap.Boot()
}

func TestMain(m *testing.M) {
	tbp.Start(m)
}

func undoRolePermissionFixture() {
	rolePermRepo, err := repo.MakeRolePermissionRepository()
	if err != nil {
		log.Fatal(err)
	}
	err = rolePermRepo.Delete(rolePermission1)
	if err != nil {
		log.Fatal(err)
	}
}

func TestGetAllFromOrganization(t *testing.T) {
	logger.Debug("TestGetAllFromOrganization...")
	tbp.PrepareTestDatabase()
	tbp.Reader = strings.NewReader("")
	rolePermissionsOrgURL := fmt.Sprintf("%s/%s/role-permissions", rolePermissionsURL, organization1)
	request, _ := http.NewRequest("GET", rolePermissionsOrgURL, tbp.Reader)
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

func TestCreateRolePermissionInOrganization(t *testing.T) {
	logger.Debug("TestCreateRole...")
	tbp.PrepareTestDatabase()
	undoRolePermissionFixture()
	roleJSON := fmt.Sprintf(`
	{
		"data": {
			"organizationID": "d43809a2-5896-43c4-808e-549f2ee47783",
		  "roleID": "9b6869e4-f51a-4197-9608-f2898bd764d8",
		  "permissionID": "cf903818-a2c5-46c2-8935-c4fc66fea60f"
		}
	}
	`, organization1)
	tbp.Reader = strings.NewReader(roleJSON)
	rolePermissionsOrgURL := fmt.Sprintf("%s/%s/role-permissions", rolePermissionsURL, organization1)
	request, _ := http.NewRequest("POST", rolePermissionsOrgURL, tbp.Reader)
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
	rolePermissionOrgURL := fmt.Sprintf("%s/%s/role-permissions/%s", rolePermissionsURL, organization1, rolePermission1)
	request, _ := http.NewRequest("GET", rolePermissionOrgURL, tbp.Reader)
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

func TestUpdateRolePermissionInOrganization(t *testing.T) {
	logger.Debug("TestUpdateRolePermissionInOrganization...")
	tbp.PrepareTestDatabase()
	newDescription := "RolePermission new description."
	roleJSON := fmt.Sprintf(`
	{
		"data": {
				"id": "%s",
				"description": "%s",
				"organizationID": "d43809a2-5896-43c4-808e-549f2ee47783",
				"roleID": "9b6869e4-f51a-4197-9608-f2898bd764d8",
				"permissionID": "%s"
		}
	}
	`, rolePermission1, newDescription, permission2)
	tbp.Reader = strings.NewReader(roleJSON)
	rolePermissionsOrgURL := fmt.Sprintf("%s/%s/role-permissions/%s", rolePermissionsURL, organization1, rolePermission1)
	request, _ := http.NewRequest("PUT", rolePermissionsOrgURL, tbp.Reader)
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

func TestUpdateRolePermissionWithDatabaseVerify(t *testing.T) {
	logger.Debug("TestUpdateRolePermissionWithDatabaseVerify...")
	tbp.PrepareTestDatabase()
	newDescription := "RolePermission new description."
	roleJSON := fmt.Sprintf(`
	{
		"data": {
			"id": "%s",
			"description": "%s",
			"organizationID": "d43809a2-5896-43c4-808e-549f2ee47783",
			"roleID": "9b6869e4-f51a-4197-9608-f2898bd764d8",
			"permissionID": "%s"
		}
	}
	`, rolePermission1, newDescription, permission2)
	tbp.Reader = strings.NewReader(roleJSON)
	rolePermissionsOrgURL := fmt.Sprintf("%s/%s/role-permissions/%s", rolePermissionsURL, organization1, rolePermission1)
	request, _ := http.NewRequest("PUT", rolePermissionsOrgURL, tbp.Reader)
	tbp.AuthorizeRequest(request, user1, "admin", "admin")
	res, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Fatal(err)
		t.Errorf("Error executing request: %s", err.Error())
		return
	}
	if res.StatusCode == http.StatusNoContent {
		rolePermRepo, err := repo.MakeRolePermissionRepository()
		if err != nil {
			log.Fatal(err)
			return
		}
		rolePerm, err := rolePermRepo.Get(rolePermission1)
		if err == nil {
			if rolePerm.Description.String == newDescription && rolePerm.PermissionID.String == permission2 {
				logger.Debug("RolePermission update: ok.")
			} else {
				error := fmt.Sprintf("Description: '%s' | Expected: '%s'", rolePerm.Description.String, newDescription)
				error += fmt.Sprintf("Permission ID: '%s' | Expected: '%s' - ", rolePerm.PermissionID, permission2)
				t.Error(error)
			}
		} else {
			t.Error(err.Error())
		}
	} else {
		t.Errorf("Status: %d | Expected: 204-StatusNoContent", res.StatusCode)
	}
}

func TestDeleteRolePermission(t *testing.T) {
	logger.Debug("TestDeleteRolePermission...")
	tbp.PrepareTestDatabase()
	tbp.Reader = strings.NewReader("")
	rolePermissionURL := fmt.Sprintf("%s/%s/role-permissions/%s", rolePermissionsURL, organization1, rolePermission1)
	request, _ := http.NewRequest("DELETE", rolePermissionURL, tbp.Reader)
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

func TestDeleteRolePermissionWithDatabaseVerify(t *testing.T) {
	logger.Debug("TestDeleteRolePermissionWithDatabaseVerify...")
	tbp.PrepareTestDatabase()
	tbp.Reader = strings.NewReader("")
	rolePermissionURL := fmt.Sprintf("%s/%s/role-permissions/%s", rolePermissionsURL, organization1, rolePermission1)
	request, _ := http.NewRequest("DELETE", rolePermissionURL, tbp.Reader)
	tbp.AuthorizeRequest(request, user1, "admin", "admin")
	res, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Fatal(err)
		t.Errorf("Error executing request: %s", err.Error())
		return
	}
	if res.StatusCode == http.StatusNoContent {
		rolePermRepo, err := repo.MakeRolePermissionRepository()
		if err != nil {
			log.Fatal(err)
			return
		}
		rolePerm, err := rolePermRepo.Get(rolePermission1)
		if err != nil {
			logger.Debug("TestDeleteRole: ok")
		} else {
			t.Errorf("Role: %s | Expected: 'nil'", rolePerm.Name)
		}
	} else {
		t.Errorf("Status: %d | Expected: 204-StatusNoContent", res.StatusCode)
	}
}
