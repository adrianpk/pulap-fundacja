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
	"log"
)

const (
	rollbackAll   = true
	migrationsNum = 13
)

var (
	env = "test"
	// BaseDir - App base directory
	BaseDir     = ""
	mgrSw       = "n" // m=migrate, r=rollback, default=none
	doMigration = true
	doRollback  = true
)

// SetBootParameters - Initialize boot
func SetBootParameters(parameters map[string]string) {
	env = parameters["env"]
	BaseDir = parameters["app_home"]
	mgrSw = parameters["migration"]
}

// Boot - Boot app
func Boot() {
	// Initialize AppConfig variable
	initConfig(env)
	// Initialize Logger objects with Log Level
	initLogger()
	// Initialize migrations
	initMigrationOrRollback()
	// Initialize private/public keys for JWT authentication
	initKeys()
	// Start a MongoDB session
	//createDbSession()
	// Add indexes into MongoDB
	//addIndexes()
	//logBootParameters()
}

func logBootParameters() {
	log.Printf("env: %s\n", env)
	log.Printf("BaseDir: %s\n", BaseDir)
	log.Printf("mgrSw: %s\n", mgrSw)
}
