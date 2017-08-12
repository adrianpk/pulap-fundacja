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
//

package bootstrap

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
)

type (
	configuration struct {
		Server                                  string
		DBHost, Database, DBUser, DBPass, DBSSL string
		BaseDir                                 string
		LogLevel                                int
		LogFile                                 string
	}
)

// AppConfig holds the configuration values from config.json file
var AppConfig configuration

// Initialize AppConfig
func initConfig(env string) {
	loadAppConfig(env)
}

// Reads config.json and decode into AppConfig
func loadAppConfig(env string) {
	configFile := fmt.Sprintf("config/config_%s.json", env)
	fullConfigPath := path.Join(BaseDir, configFile)
	file, err := os.Open(fullConfigPath)
	defer file.Close()
	AppConfig.BaseDir = fullConfigPath
	if err != nil {
		log.Fatalf("[Fundacja]: %s\n", err)
	}
	decoder := json.NewDecoder(file)
	AppConfig = configuration{}
	err = decoder.Decode(&AppConfig)
	//log.Printf("Config: %v", AppConfig)
	if err != nil {
		log.Fatalf("[Fundacja]: %s\n", err)
	}
}