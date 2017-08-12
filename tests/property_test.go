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
	tbp               = testbootstrap.TestBootstrap
	propertiesSetsURL string
	propertiesSet1    = "b2672e94-e1ba-4c23-a428-b91528d06d1f"
	propertiesSet2    = "d4824e55-6305-4898-ab3a-7023568a1d27"
	property1         = "a3066a59-5e2e-4305-9844-d9e2c2317924"
	property2         = "8b348adf-01b9-45d1-8d11-e9d4e0090e20"
	property1Name     = "Property1"
	user1             = "5958b185-8150-4aae-b53f-0c44771ddec5"
	user2             = "3c05e701-b495-4443-b454-2c37e2ecccdf"
)

func init() {
	propertiesSetsURL = fmt.Sprintf("%s/properties-set", tbp.APIServerURL)
	bootstrap.SetBootParameters(testbootstrap.BootParameters())
	bootstrap.Boot()
}

func TestMain(m *testing.M) {
	tbp.Start(m)
}

func undoPropertyFixture() {
	propertyRepo, err := repo.MakePropertyRepository()
	if err != nil {
		log.Fatal(err)
	}
	err = propertyRepo.Delete(property1)
	if err != nil {
		log.Fatal(err)
	}
}

func TestGetAllFromOrganization(t *testing.T) {
	logger.Debug("TestGetAll...")
	tbp.PrepareTestDatabase()
	tbp.Reader = strings.NewReader("")
	propertiesURL := fmt.Sprintf("%s/%s/properties", propertiesSetsURL, propertiesSet1)
	request, _ := http.NewRequest("GET", propertiesURL, tbp.Reader)
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

func TestCreateProperty(t *testing.T) {
	logger.Debug("TestCreateProperty...")
	tbp.PrepareTestDatabase()
	undoPropertyFixture()
	newName := "Property name"
	newDescription := "Property description."
	newValue := "Property1-Value"
	newValueType := "s"
	propertyJSON := fmt.Sprintf(`
	{
		"data": {
			"id": "%s",
		 	"name": "%s",
		 	"description": "%s",
  		"string_value": "%s",
		 	"valueType": "%s",
		 	"position": 0,
		 	"properties_set_id": "%s"
		}
	}
	`, property1, newName, newDescription, newValue, newValueType, propertiesSet1)
	tbp.Reader = strings.NewReader(propertyJSON)
	propertiesURL := fmt.Sprintf("%s/%s/properties", propertiesSetsURL, propertiesSet1)
	request, _ := http.NewRequest("POST", propertiesURL, tbp.Reader)
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
	propertyURL := fmt.Sprintf("%s/%s/properties/%s", propertiesSetsURL, propertiesSet1, property1)
	request, _ := http.NewRequest("GET", propertyURL, tbp.Reader)
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
	propertyURL := fmt.Sprintf("%s/%s/properties/%s", propertiesSetsURL, propertiesSet1, property1Name)
	request, _ := http.NewRequest("GET", propertyURL, tbp.Reader)
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

func TestUpdateProperty(t *testing.T) {
	logger.Debug("TestUpdateProperty...")
	tbp.PrepareTestDatabase()
	newName := "Property new name"
	newDescription := "Property new description."
	newValue := 512
	newValueType := "i"
	propertyJSON := fmt.Sprintf(`
	{
		"data": {
			"id": "%s",
		 	"name": "%s",
		 	"description": "%s",
		 	"stringValue": "",
		 	"intValue": %d,
		 	"valueType": "%s",
		 	"position": 0,
		 	"propertiesSetId": "%s"
		}
	}
	`, property1, newName, newDescription, newValue, newValueType, propertiesSet1)
	tbp.Reader = strings.NewReader(propertyJSON)
	propertyURL := fmt.Sprintf("%s/%s/properties/%s", propertiesSetsURL, propertiesSet1, property1)
	request, _ := http.NewRequest("PUT", propertyURL, tbp.Reader)
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

func TestDoUpdatePropertyWithDatabaseVerify(t *testing.T) {
	logger.Debug("TestUpdatePropertyWithDatabaseVerify...")
	tbp.PrepareTestDatabase()
	newName := "Property new name"
	newDescription := "Property new description."
	newValue := 512
	newValueType := "i"
	propertyJSON := fmt.Sprintf(`
	{
		"data": {
			"id": "%s",
		 	"name": "%s",
		 	"description": "%s",
		 	"stringValue": "",
		 	"intValue": %d,
		 	"valueType": "%s",
		 	"position": 0,
		 	"propertiesSetId": "%s"
		}
	}
	`, property1, newName, newDescription, newValue, newValueType, propertiesSet1)
	tbp.Reader = strings.NewReader(propertyJSON)
	propertyURL := fmt.Sprintf("%s/%s/properties/%s", propertiesSetsURL, propertiesSet1, property1)
	request, _ := http.NewRequest("PUT", propertyURL, tbp.Reader)
	tbp.AuthorizeRequest(request, user1, "admin", "admin")
	res, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Fatal(err)
		t.Errorf("Error executing request: %s", err.Error())
		return
	}
	if res.StatusCode == http.StatusNoContent {
		propertyRepo, err := repo.MakePropertyRepository()
		if err != nil {
			log.Fatal(err)
			return
		}
		property, err := propertyRepo.Get(property1)
		if err == nil {
			if property.Name.String == newName && property.Description.String == newDescription {
				logger.Debug("Property update: ok.")
			} else {
				error := fmt.Sprintf("Name: '%s' | Expected: '%s' - ", property.Name.String, newName)
				error += fmt.Sprintf("Description: '%s' | Expected: '%s'", property.Description.String, newDescription)
				t.Error(error)
			}
		} else {
			t.Error(err.Error())
		}
	} else {
		t.Errorf("Status: %d | Expected: 204-StatusNoContent", res.StatusCode)
	}
}

func TestDeleteProperty(t *testing.T) {
	logger.Debug("TestDeleteProperty...")
	tbp.PrepareTestDatabase()
	tbp.Reader = strings.NewReader("")
	propertyURL := fmt.Sprintf("%s/%s/properties/%s", propertiesSetsURL, propertiesSet1, property1)
	request, _ := http.NewRequest("DELETE", propertyURL, tbp.Reader)
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

func TestDeletePropertyWithDatabaseVerify(t *testing.T) {
	logger.Debug("TestDeletePropertyWithDatabaseVerify...")
	tbp.PrepareTestDatabase()
	tbp.Reader = strings.NewReader("")
	propertyURL := fmt.Sprintf("%s/%s/properties/%s", propertiesSetsURL, propertiesSet1, property1)
	request, _ := http.NewRequest("DELETE", propertyURL, tbp.Reader)
	tbp.AuthorizeRequest(request, user1, "admin", "admin")
	res, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Fatal(err)
		t.Errorf("Error executing request: %s", err.Error())
		return
	}
	if res.StatusCode == http.StatusNoContent {
		propertyRepo, err := repo.MakePropertyRepository()
		if err != nil {
			log.Fatal(err)
			return
		}
		property, err := propertyRepo.Get(property1)
		if err != nil {
			logger.Debug("TestDeleteProperty: ok")
		} else {
			t.Errorf("Property: %s | Expected: 'nil'", property.Name)
		}
	} else {
		t.Errorf("Status: %d | Expected: 204-StatusNoContent", res.StatusCode)
	}
}
