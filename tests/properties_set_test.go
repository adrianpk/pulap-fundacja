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
	tbp            = testbootstrap.TestBootstrap
	holdersURL     string
	user1          = "5958b185-8150-4aae-b53f-0c44771ddec5"
	user2          = "3c05e701-b495-4443-b454-2c37e2ecccdf"
	propertiesSet1 = "b2672e94-e1ba-4c23-a428-b91528d06d1f"
	propertiesSet2 = "d4824e55-6305-4898-ab3a-7023568a1d27"
	resource1Name  = "Controllers1"
)

func init() {
	holdersURL = fmt.Sprintf("%s/holders", tbp.APIServerURL)
	bootstrap.SetBootParameters(testbootstrap.BootParameters())
	bootstrap.Boot()
}

func TestMain(m *testing.M) {
	tbp.Start(m)
}

func undoResourceFixture() {
	resourceRepo, err := repo.MakePropertiesSetRepository()
	if err != nil {
		log.Fatal(err)
	}
	err = resourceRepo.Delete(propertiesSet1)
	if err != nil {
		log.Fatal(err)
	}
}

func TestGetAllForHolder(t *testing.T) {
	logger.Debug("TestGetAllForHolder...")
	tbp.PrepareTestDatabase()
	tbp.Reader = strings.NewReader("")
	propSetHolderURL := fmt.Sprintf("%s/%s/properties-sets", holdersURL, user1)
	request, _ := http.NewRequest("GET", propSetHolderURL, tbp.Reader)
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

func TestCreatePropertiesSet(t *testing.T) {
	logger.Debug("TestCreatePropertiesSet...")
	tbp.PrepareTestDatabase()
	undoResourceFixture()
	propertiesSet1JSON := fmt.Sprintf(`
	{
		"data": {
			"name": "PropertiesSet",
		  "description": "PropertiesSet description.",
			"position": 0,
			"holderID": "%s"
		}
	}
	`, user1)
	tbp.Reader = strings.NewReader(propertiesSet1JSON)
	propSetHolderURL := fmt.Sprintf("%s/%s/properties-sets", holdersURL, user1)
	request, _ := http.NewRequest("POST", propSetHolderURL, tbp.Reader)
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

func TestGetFromHolder(t *testing.T) {
	logger.Debug("TestGetFromHolder...")
	tbp.PrepareTestDatabase()
	tbp.Reader = strings.NewReader("")
	propSetHolderURL := fmt.Sprintf("%s/%s/properties-sets", holdersURL, user1)
	request, _ := http.NewRequest("GET", propSetHolderURL, tbp.Reader)
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
	propSetHolderURL := fmt.Sprintf("%s/%s/properties-sets", holdersURL, user1)
	request, _ := http.NewRequest("GET", propSetHolderURL, tbp.Reader)
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

func TestUpdatePropertiesSet(t *testing.T) {
	logger.Debug("TestUpdateResource...")
	tbp.PrepareTestDatabase()
	newName := "PropertiesSet new name"
	newDescription := "PropertiesSet new description."
	newPosition := int64(1)
	propertiesSet1JSON := fmt.Sprintf(`
	{
		"data": {
			"name": "%s",
			"description": "%s",
			"position": %d,
			"holderID": "%s"
		}
	}
	`, newName, newDescription, newPosition, user1)
	tbp.Reader = strings.NewReader(propertiesSet1JSON)
	propSetHolderURL := fmt.Sprintf("%s/%s/properties-sets/%s", holdersURL, user1, propertiesSet1)
	request, _ := http.NewRequest("PUT", propSetHolderURL, tbp.Reader)
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

func TestUpdatePropSetWithDatabaseVerify(t *testing.T) {
	logger.Debug("TestUpdatePropSetWithDatabaseVerify...")
	tbp.PrepareTestDatabase()
	tbp.PrepareTestDatabase()
	newName := "PropertiesSet new name"
	newDescription := "PropertiesSet new description."
	newPosition := int64(1)
	propertiesSet1JSON := fmt.Sprintf(`
	{
		"data": {
			"name": "%s",
			"description": "%s",
			"position": %d,
			"holderID": "%s"
		}
	}
	`, newName, newDescription, newPosition, user1)
	tbp.Reader = strings.NewReader(propertiesSet1JSON)
	propSetHolderURL := fmt.Sprintf("%s/%s/properties-sets/%s", holdersURL, user1, propertiesSet1)
	request, _ := http.NewRequest("PUT", propSetHolderURL, tbp.Reader)
	tbp.AuthorizeRequest(request, user1, "admin", "admin")
	res, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Fatal(err)
		t.Errorf("Error executing request: %s", err.Error())
		return
	}
	if res.StatusCode == http.StatusNoContent {
		propSetRepo, err := repo.MakePropertiesSetRepository()
		if err != nil {
			log.Fatal(err)
			return
		}
		propSet, err := propSetRepo.Get(propertiesSet1)
		if err == nil {
			if propSet.Name.String == newName && propSet.Description.String == newDescription && propSet.Position.Int64 == newPosition {
				logger.Debug("PropertiesSet update: ok.")
			} else {
				error := fmt.Sprintf("Name: '%s' | Expected: '%s' - ", propSet.Name.String, newName)
				error += fmt.Sprintf("Description: '%s' | Expected: '%s'", propSet.Description.String, newDescription)
				error += fmt.Sprintf("Position: '%s' | Expected: '%s'", propSet.Description.String, newPosition)
				t.Error(error)
			}
		} else {
			t.Error(err.Error())
		}
	} else {
		t.Errorf("Status: %d | Expected: 204-StatusNoContent", res.StatusCode)
	}
}

func TestDeletePropertiesSet(t *testing.T) {
	logger.Debug("TestDeletePropertiesSet...")
	tbp.PrepareTestDatabase()
	tbp.Reader = strings.NewReader("")
	propSetHolderURL := fmt.Sprintf("%s/%s/properties-sets/%s", holdersURL, user1, propertiesSet1)
	request, _ := http.NewRequest("DELETE", propSetHolderURL, tbp.Reader)
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

func TestDeletePropertiesSetWithDatabaseVerify(t *testing.T) {
	logger.Debug("TestDeletePropertiesSetWithDatabaseVerify...")
	tbp.PrepareTestDatabase()
	tbp.Reader = strings.NewReader("")
	propSetHolderURL := fmt.Sprintf("%s/%s/properties-sets/%s", holdersURL, user1, propertiesSet1)
	request, _ := http.NewRequest("DELETE", propSetHolderURL, tbp.Reader)
	tbp.AuthorizeRequest(request, user1, "admin", "admin")
	res, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Fatal(err)
		t.Errorf("Error executing request: %s", err.Error())
		return
	}
	if res.StatusCode == http.StatusNoContent {
		propSetRepoRepo, err := repo.MakePropertiesSetRepository()
		if err != nil {
			log.Fatal(err)
			return
		}
		propSet, err := propSetRepoRepo.Get(propertiesSet1)
		if err != nil {
			logger.Debug("TestDeletePropertiesSetWithDatabaseVerify: ok")
		} else {
			t.Errorf("Resource: %s | Expected: 'nil'", propSet.Name)
		}
	} else {
		t.Errorf("Status: %d | Expected: 204-StatusNoContent", res.StatusCode)
	}
}
