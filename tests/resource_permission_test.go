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
	tbp                     = testbootstrap.TestBootstrap
	organizationsURL        string
	user1                   = "5958b185-8150-4aae-b53f-0c44771ddec5"
	user2                   = "3c05e701-b495-4443-b454-2c37e2ecccdf"
	organization1           = "d43809a2-5896-43c4-808e-549f2ee47783"
	organization2           = "b8cef4be-1ec3-44b4-9cbd-551f039f4fc7"
	resourcePermission1     = "e42213a8-cbd9-4957-b82c-6805ef59d123"
	resourcePermission2     = "cb1e77a8-38ed-49ea-ad47-2ba880fbe57c"
	resource1               = "656fe625-1f1e-4528-bec3-c11cf254cfe4"
	resource2               = "55a192d4-5d1a-4b1a-9fd0-f644394e9457"
	permission1             = "cf903818-a2c5-46c2-8935-c4fc66fea60f"
	permission2             = "131a9338-39aa-4c33-8661-dcfa21ce726f"
	resourcePermission1Name = "ResourcePermission1"
)

func init() {
	organizationsURL = fmt.Sprintf("%s/organizations", tbp.APIServerURL)
	bootstrap.SetBootParameters(testbootstrap.BootParameters())
	bootstrap.Boot()
}

func TestMain(m *testing.M) {
	tbp.Start(m)
}

func undoResourcePermissionFixture() {
	resourcePermRepo, err := repo.MakeResourcePermissionRepository()
	if err != nil {
		log.Fatal(err)
	}
	err = resourcePermRepo.Delete(resourcePermission1)
	if err != nil {
		log.Fatal(err)
	}
}

func TestGetAllFromOrganization(t *testing.T) {
	logger.Debug("TestGetAllFromOrganization...")
	tbp.PrepareTestDatabase()
	tbp.Reader = strings.NewReader("")
	resourcePermissionsOrgURL := fmt.Sprintf("%s/%s/resource-permissions", organizationsURL, organization1)
	request, _ := http.NewRequest("GET", resourcePermissionsOrgURL, tbp.Reader)
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

func TestCreateResourcePermissionInOrganization(t *testing.T) {
	logger.Debug("TestCreateresource...")
	tbp.PrepareTestDatabase()
	undoResourcePermissionFixture()
	roleJSON := fmt.Sprintf(`
	{
		"data": {
				"organizationID": "%s",
				"resourceID": "%s",
				"permissionID": "%s"
		}
	}
	`, organization1, resource1, permission1)
	tbp.Reader = strings.NewReader(roleJSON)
	resourcePermissionsOrgURL := fmt.Sprintf("%s/%s/resource-permissions", organizationsURL, organization1)
	request, _ := http.NewRequest("POST", resourcePermissionsOrgURL, tbp.Reader)
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
	resourcePermissionOrgURL := fmt.Sprintf("%s/%s/resource-permissions/%s", organizationsURL, organization1, resourcePermission1)
	request, _ := http.NewRequest("GET", resourcePermissionOrgURL, tbp.Reader)
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

func TestUpdateResourcePermissionInOrganization(t *testing.T) {
	logger.Debug("TestUpdateResourcePermissionInOrganization...")
	tbp.PrepareTestDatabase()
	newDescription := "RolePermission new description."
	roleJSON := fmt.Sprintf(`
	{
		"data": {
				"id": "%s",
				"description": "%s",
				"organizationID": "%s",
				"resourceID": "%s",
				"permissionID": "%s"
		}
	}
	`, resourcePermission1, newDescription, organization1, resource1, permission2)
	tbp.Reader = strings.NewReader(roleJSON)
	resourcePermissionsOrgURL := fmt.Sprintf("%s/%s/resource-permissions/%s", organizationsURL, organization1, resourcePermission1)
	request, _ := http.NewRequest("PUT", resourcePermissionsOrgURL, tbp.Reader)
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

func TestUpdateResourcePermissionWithDatabaseVerify(t *testing.T) {
	logger.Debug("TestUpdateResourcePermissionWithDatabaseVerify...")
	tbp.PrepareTestDatabase()
	newDescription := "RolePermission new description."
	roleJSON := fmt.Sprintf(`
	{
		"data": {
				"id": "%s",
				"description": "%s",
				"organizationID": "%s",
				"resourceID": "%s",
				"permissionID": "%s"
		}
	}
	`, resourcePermission1, newDescription, organization1, resource1, permission2)
	tbp.Reader = strings.NewReader(roleJSON)
	resourcePermissionsOrgURL := fmt.Sprintf("%s/%s/resource-permissions/%s", organizationsURL, organization1, resourcePermission1)
	request, _ := http.NewRequest("PUT", resourcePermissionsOrgURL, tbp.Reader)
	tbp.AuthorizeRequest(request, user1, "admin", "admin")
	res, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Fatal(err)
		t.Errorf("Error executing request: %s", err.Error())
		return
	}
	if res.StatusCode == http.StatusNoContent {
		resourcePermRepo, err := repo.MakeResourcePermissionRepository()
		if err != nil {
			log.Fatal(err)
			return
		}
		resourcePerm, err := resourcePermRepo.Get(resourcePermission1)
		if err == nil {
			if resourcePerm.Description.String == newDescription && resourcePerm.PermissionID.String == permission2 {
				logger.Debug("ResourcePermission update: ok.")
			} else {
				error := fmt.Sprintf("Description: '%s' | Expected: '%s'", resourcePerm.Description.String, newDescription)
				error += fmt.Sprintf("Permission ID: '%s' | Expected: '%s' - ", resourcePerm.PermissionID, permission2)
				t.Error(error)
			}
		} else {
			t.Error(err.Error())
		}
	} else {
		t.Errorf("Status: %d | Expected: 204-StatusNoContent", res.StatusCode)
	}
}

func TestDeleteResourcePermission(t *testing.T) {
	logger.Debug("TestDeleteResourcePermission...")
	tbp.PrepareTestDatabase()
	tbp.Reader = strings.NewReader("")
	roleURL := fmt.Sprintf("%s/%s/resource-permissions/%s", organizationsURL, organization1, resourcePermission1)
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

func TestDeleteResourcePermissionWithDatabaseVerify(t *testing.T) {
	logger.Debug("TestDeleteResourcePermissionWithDatabaseVerify...")
	tbp.PrepareTestDatabase()
	tbp.Reader = strings.NewReader("")
	roleURL := fmt.Sprintf("%s/%s/resource-permissions/%s", organizationsURL, organization1, resourcePermission1)
	request, _ := http.NewRequest("DELETE", roleURL, tbp.Reader)
	tbp.AuthorizeRequest(request, user1, "admin", "admin")
	res, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Fatal(err)
		t.Errorf("Error executing request: %s", err.Error())
		return
	}
	if res.StatusCode == http.StatusNoContent {
		resourcePermRepo, err := repo.MakeResourcePermissionRepository()
		if err != nil {
			log.Fatal(err)
			return
		}
		resourcePerm, err := resourcePermRepo.Get(resourcePermission1)
		if err != nil {
			logger.Debug("TestDeleteRole: ok")
		} else {
			t.Errorf("Role: %s | Expected: 'nil'", resourcePerm.Name)
		}
	} else {
		t.Errorf("Status: %d | Expected: 204-StatusNoContent", res.StatusCode)
	}
}
