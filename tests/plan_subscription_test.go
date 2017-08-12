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
	subscriptorURL    string
	user1             = "5958b185-8150-4aae-b53f-0c44771ddec5"
	user2             = "3c05e701-b495-4443-b454-2c37e2ecccdf"
	plan1             = "50273641-8fa0-4275-af3c-a7d2e3025c61"
	plan2             = "51dbccff-7503-462a-b92e-ed4b852b6bab"
	planSubscription1 = "e82ae81a-aadf-43ff-990f-2d28689f0384"
	planSubscription2 = "05da878c-03bb-4d55-8cae-2c5bd149def9"
	username          = "admin"
)

func init() {
	subscriptorURL = fmt.Sprintf("%s/plan-subscriptor", tbp.APIServerURL)
	bootstrap.SetBootParameters(testbootstrap.BootParameters())
	bootstrap.Boot()
}

func TestMain(m *testing.M) {
	tbp.Start(m)
}

func undoResourceFixture() {
	resourceRepo, err := repo.MakePlanSubscriptionRepository()
	if err != nil {
		log.Fatal(err)
	}
	err = resourceRepo.Delete(planSubscription1)
	if err != nil {
		log.Fatal(err)
	}
}

func TestCreatePlanSubscription(t *testing.T) {
	logger.Debug("TestCreatePlanSubscription...")
	tbp.PrepareTestDatabase()
	undoResourceFixture()
	planSubscription1JSON := fmt.Sprintf(`
	{
		"data": {
		  "name": "Admin::Organization",
		  "description": "Admin::Organization subscription.",
		  "startedAt": 1483268400,
		  "endsAt": 1514804400,
		  "userId": "%s",
		  "planId": "%s"
		}
	}
	`, user1, plan1)
	tbp.Reader = strings.NewReader(planSubscription1JSON)
	planSubURL := fmt.Sprintf("%s/%s/plan-subscription", subscriptorURL, user1)
	request, _ := http.NewRequest("POST", planSubURL, tbp.Reader)
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

func TestGetFromSubscriptor(t *testing.T) {
	logger.Debug("TestGetFrom...")
	tbp.PrepareTestDatabase()
	tbp.Reader = strings.NewReader("")
	planSubURL := fmt.Sprintf("%s/%s/plan-subscription", subscriptorURL, user1)
	request, _ := http.NewRequest("GET", planSubURL, tbp.Reader)
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
	planSubURL := fmt.Sprintf("%s/%s/plan-subscription", subscriptorURL, user1)
	request, _ := http.NewRequest("GET", planSubURL, tbp.Reader)
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

func TestUpdatePlanSubscription(t *testing.T) {
	logger.Debug("TestUpdateResource...")
	tbp.PrepareTestDatabase()
	newName := "Admin::Organization new name"
	newDescription := "Admin::Organization new subscription.."
	planSubscription1JSON := fmt.Sprintf(`
	{
		"data": {
      "id": "%s",
      "name": "%s",
      "description": "%s",
      "userID": "%s",
      "planID": "%s"
		}
	}
	`, planSubscription1, newName, newDescription, user1, plan2)
	tbp.Reader = strings.NewReader(planSubscription1JSON)
	planSubURL := fmt.Sprintf("%s/%s/plan-subscription", subscriptorURL, user1)
	request, _ := http.NewRequest("PUT", planSubURL, tbp.Reader)
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
	newName := "Admin::Organization new name"
	newDescription := "Admin::Organization new subscription.."
	planSubscription1JSON := fmt.Sprintf(`
  {
		"data": {
      "id": "%s",
      "name": "%s",
      "description": "%s",
      "userID": "%s",
      "planID": "%s"
		}
	}
	`, planSubscription1, newName, newDescription, user1, plan2)
	tbp.Reader = strings.NewReader(planSubscription1JSON)
	planSubURL := fmt.Sprintf("%s/%s/plan-subscription", subscriptorURL, user1)
	request, _ := http.NewRequest("PUT", planSubURL, tbp.Reader)
	tbp.AuthorizeRequest(request, user1, "admin", "admin")
	res, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Fatal(err)
		t.Errorf("Error executing request: %s", err.Error())
		return
	}
	if res.StatusCode == http.StatusNoContent {
		planSubRepo, err := repo.MakePlanSubscriptionRepository()
		if err != nil {
			log.Fatal(err)
			return
		}
		planSub, err := planSubRepo.Get(planSubscription1)
		if err == nil {
			if planSub.Name.String == newName && planSub.Description.String == newDescription && planSub.PlanID.String == plan2 {
				logger.Debug("PlanSubscriptions update: ok.")
			} else {
				error := fmt.Sprintf("Name: '%s' | Expected: '%s' - ", planSub.Name.String, newName)
				error += fmt.Sprintf("Description: '%s' | Expected: '%s'", planSub.Description.String, newDescription)
				error += fmt.Sprintf("Plan: '%s' | Expected: '%s'", planSub.PlanID.String, plan2)
				t.Error(error)
			}
		} else {
			t.Error(err.Error())
		}
	} else {
		t.Errorf("Status: %d | Expected: 204-StatusNoContent", res.StatusCode)
	}
}

func TestDeletePlanSubscription(t *testing.T) {
	logger.Debug("TestDeletePlanSubscription...")
	tbp.PrepareTestDatabase()
	tbp.Reader = strings.NewReader("")
	planSubURL := fmt.Sprintf("%s/%s/plan-subscription", subscriptorURL, user1)
	request, _ := http.NewRequest("DELETE", planSubURL, tbp.Reader)
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

func TestDeletePlanSubscriptionWithDatabaseVerify(t *testing.T) {
	logger.Debug("TestDeletePlanSubscriptionWithDatabaseVerify...")
	tbp.PrepareTestDatabase()
	tbp.Reader = strings.NewReader("")
	planSubURL := fmt.Sprintf("%s/%s/plan-subscription", subscriptorURL, user1)
	request, _ := http.NewRequest("DELETE", planSubURL, tbp.Reader)
	tbp.AuthorizeRequest(request, user1, "admin", "admin")
	res, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Fatal(err)
		t.Errorf("Error executing request: %s", err.Error())
		return
	}
	if res.StatusCode == http.StatusNoContent {
		planSubRepoRepo, err := repo.MakePlanSubscriptionRepository()
		if err != nil {
			log.Fatal(err)
			return
		}
		planSub, err := planSubRepoRepo.Get(planSubscription1)
		if err != nil {
			logger.Debug("TestDeletePlanSubscriptionWithDatabaseVerify: ok")
		} else {
			t.Errorf("Resource: %s | Expected: 'nil'", planSub.Name)
		}
	} else {
		t.Errorf("Status: %d | Expected: 204-StatusNoContent", res.StatusCode)
	}
}
