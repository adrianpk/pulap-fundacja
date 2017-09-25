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

package testbootstrap

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/adrianpk/pulap/bootstrap"
	"github.com/adrianpk/pulap/handler"

	testfixtures "gopkg.in/testfixtures.v2"
)

// 	TestBootstrap testBootstrap
//  holds the configuration values from config.json file
var (
	env           = "test"
	baseTest      = ""
	fixturesDir   = "resources/fixtures"
	TestBootstrap testBootstrap
	apiPath       = "api"
	apiVersion    = "v1"
)

// BootParameters - Default boot parameters for tests
func BootParameters() map[string]string {
	params := make(map[string]string)
	// Envs: "dev", "test", "prod"
	// Migrations: "m", "r", "n" - migrate, rollback, none
	params["env"] = "test"
	params["app_home"] = os.Getenv("PULAP_HOME")
	params["migration"] = "m"
	return params
}

// Reads config.json and decode into TestBootstrap
func init() {
	configFile := fmt.Sprintf("config/config_%s.json", env)
	baseTest = getBaseDir()
	fullConfigPath := path.Join(baseTest, configFile)
	file, err := os.Open(fullConfigPath)
	defer file.Close()
	TestBootstrap.BaseDir = fullConfigPath
	if err != nil {
		log.Fatalf("[Fundacja-Test]: %s\n", err)
	}
	decoder := json.NewDecoder(file)
	TestBootstrap = testBootstrap{}
	err = decoder.Decode(&TestBootstrap)
	//log.Printf("Test config: %v", TestBootstrap)
	if err != nil {
		log.Fatalf("[Fundacja-Test]: %s\n", err)
	}
	TestBootstrap.BaseDir = baseTest
	TestBootstrap.FixturesDir = path.Join(baseTest, fixturesDir)
	TestBootstrap.DBConfig = fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s", TestBootstrap.DBUser, TestBootstrap.DBPass, TestBootstrap.Database, TestBootstrap.DBSSL)
	TestBootstrap.ServerInstance = httptest.NewServer(handler.AppHandler(TestBootstrap))
	TestBootstrap.APIPath = apiPath
	TestBootstrap.APIVersion = apiVersion
	TestBootstrap.APIServerURL = fmt.Sprintf("%s/%s/%s", TestBootstrap.ServerInstance.URL, TestBootstrap.APIPath, TestBootstrap.APIVersion)
}

func (configurator *testBootstrap) Start(m *testing.M) {
	var err error
	// Open connection with the test database.
	// Do NOT import fixtures in a production database!
	// Existing data would be deleted
	// connStr := fmt.Sprintf("user=%s dbname=%s sslmode=%s", cfg.User, cfg.DB, cfg.SSL)
	TestBootstrap.DBInstance, err = sql.Open("postgres", TestBootstrap.DBConfig)
	if err != nil {
		log.Fatal(err)
	}
	// creating the context that hold the fixtures
	TestBootstrap.Fixtures, err = testfixtures.NewFolder(TestBootstrap.DBInstance, &testfixtures.PostgreSQL{}, TestBootstrap.FixturesDir)
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(m.Run())
}

func (configurator *testBootstrap) PrepareTestDatabase() {
	if err := TestBootstrap.Fixtures.Load(); err != nil {
		log.Fatal(err)
	}
}

func (configurator *testBootstrap) AuthorizeRequest(req *http.Request, user, username, role string) {
	token, _ := bootstrap.GenerateJWT(user, username, role)
	var bearer = "Bearer " + token
	req.Header.Add("authorization", bearer)
}

func getBaseDir() string {
	exPath, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	parentPath := filepath.Dir(exPath)
	return parentPath
}
