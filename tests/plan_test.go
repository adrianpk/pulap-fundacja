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
	tbp               = testbootstrap.TestBootstrap
	plansURL          string
	user1             = "5958b185-8150-4aae-b53f-0c44771ddec5"
	user2             = "3c05e701-b495-4443-b454-2c37e2ecccdf"
	plan1             = "50273641-8fa0-4275-af3c-a7d2e3025c61"
	plan2             = "51dbccff-7503-462a-b92e-ed4b852b6bab"
	planSubscription1 = "6966bdba-6227-4a99-b098-1da083f363dd"
	planSubscription2 = "8f6342b8-4817-4a4a-8780-f3a43fec6c6e"
	plan1Name         = "Plan"
)

func init() {
	plansURL = fmt.Sprintf("%s/plans", tbp.APIServerURL)
	bootstrap.SetBootParameters(testbootstrap.BootParameters())
	bootstrap.Boot()
}

func TestMain(m *testing.M) {
	tbp.Start(m)
}

func undoPlanFixture() {
	planRepo, err := repo.MakePlanRepository()
	if err != nil {
		log.Fatal(err)
	}
	err = planRepo.Delete(plan1)
	if err != nil {
		log.Fatal(err)
	}
}

func TestGetAll(t *testing.T) {
	logger.Debug("TestGetAll...")
	tbp.PrepareTestDatabase()
	tbp.Reader = strings.NewReader("")
	request, _ := http.NewRequest("GET", plansURL, tbp.Reader)
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

func TestCreatePlan(t *testing.T) {
	logger.Debug("TestCreatePlan...")
	tbp.PrepareTestDatabase()
	undoPlanFixture()
	planJSON := fmt.Sprintf(`
	{
		"data": {
			"name": "Plan",
		  "description": "Plan description.",
			"startedAt": 1483268400,
			"endsAt": 1514804400,
			"createdBy": "%s"
		}
	}
	`, user1)
	tbp.Reader = strings.NewReader(planJSON)
	planURL := fmt.Sprintf("%s", plansURL)
	request, _ := http.NewRequest("POST", planURL, tbp.Reader)
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
	planURL := fmt.Sprintf("%s/%s", plansURL, plan1)
	request, _ := http.NewRequest("GET", planURL, tbp.Reader)
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
	planURL := fmt.Sprintf("%s/%s", plansURL, plan1Name)
	request, _ := http.NewRequest("GET", planURL, tbp.Reader)
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

func TestUpdatePlan(t *testing.T) {
	logger.Debug("TestUpdatePlan...")
	tbp.PrepareTestDatabase()
	newName := "Plan new name"
	newDescription := "Plan new description."
	planJSON := fmt.Sprintf(`
	{
		"data": {
			"name": "%s",
		  "description": "%s",
			"startedAt": 1483268400,
			"endsAt": 1514804400,
			"createdBy": "%s"
		}
	}
	`, newName, newDescription)
	tbp.Reader = strings.NewReader(planJSON)
	planURL := fmt.Sprintf("%s/%s", plansURL, plan1)
	request, _ := http.NewRequest("PUT", planURL, tbp.Reader)
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

func TestDoUpdatePlanWithDatabaseVerify(t *testing.T) {
	logger.Debug("TestUpdateUserWithDatabaseVerify...")
	tbp.PrepareTestDatabase()
	newName := "Plan new name"
	newDescription := "Plan new description."
	planJSON := fmt.Sprintf(`
	{
		"data": {
			"name": "%s",
		  "description": "%s",
			"startedAt": 1483268400,
			"endsAt": 1514804400,
			"createdBy": "%s"
		}
	}
	`, newName, newDescription)
	tbp.Reader = strings.NewReader(planJSON)
	planURL := fmt.Sprintf("%s/%s", plansURL, plan1)
	request, _ := http.NewRequest("PUT", planURL, tbp.Reader)
	tbp.AuthorizeRequest(request, user1, "admin", "admin")
	res, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Fatal(err)
		t.Errorf("Error executing request: %s", err.Error())
		return
	}
	if res.StatusCode == http.StatusNoContent {
		planRepo, err := repo.MakePlanRepository()
		if err != nil {
			log.Fatal(err)
			return
		}
		plan, err := planRepo.Get(plan1)
		if err == nil {
			if plan.Name.String == newName && plan.Description.String == newDescription {
				logger.Debug("Plan update: ok.")
			} else {
				error := fmt.Sprintf("Name: '%s' | Expected: '%s' - ", plan.Name.String, newName)
				error += fmt.Sprintf("Description: '%s' | Expected: '%s'", plan.Description.String, newDescription)
				t.Error(error)
			}
		} else {
			t.Error(err.Error())
		}
	} else {
		t.Errorf("Status: %d | Expected: 204-StatusNoContent", res.StatusCode)
	}
}

func TestDeletePlan(t *testing.T) {
	logger.Debug("TestDeletePlan...")
	tbp.PrepareTestDatabase()
	tbp.Reader = strings.NewReader("")
	planURL := fmt.Sprintf("%s/%s", plansURL, plan1)
	request, _ := http.NewRequest("DELETE", planURL, tbp.Reader)
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

func TestDeletePlanWithDatabaseVerify(t *testing.T) {
	logger.Debug("TestDeletePlanWithDatabaseVerify...")
	tbp.PrepareTestDatabase()
	tbp.Reader = strings.NewReader("")
	planURL := fmt.Sprintf("%s/%s", plansURL, plan1)
	request, _ := http.NewRequest("DELETE", planURL, tbp.Reader)
	tbp.AuthorizeRequest(request, user1, "admin", "admin")
	res, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Fatal(err)
		t.Errorf("Error executing request: %s", err.Error())
		return
	}
	if res.StatusCode == http.StatusNoContent {
		planRepo, err := repo.MakePlanRepository()
		if err != nil {
			log.Fatal(err)
			return
		}
		plan, err := planRepo.Get(plan1)
		if err != nil {
			logger.Debug("TestDeletePlan: ok")
		} else {
			t.Errorf("Plan: %s | Expected: 'nil'", plan.Name)
		}
	} else {
		t.Errorf("Status: %d | Expected: 204-StatusNoContent", res.StatusCode)
	}
}
