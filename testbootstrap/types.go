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
	"io"
	"net/http/httptest"

	testfixtures "gopkg.in/testfixtures.v2"
)

type (
	testBootstrap struct {
		Server                                  string
		DBHost, Database, DBUser, DBPass, DBSSL string
		ResourcesDir                            string
		PublicDir                               string
		LogFile                                 string
		LogLevel                                int
		Autoreload                              bool
		BaseDir                                 string
		FixturesDir                             string
		DBConfig                                string
		DBInstance                              *sql.DB
		Fixtures                                *testfixtures.Context
		ServerInstance                          *httptest.Server
		Reader                                  io.Reader //Ignore this for now
		APIPath                                 string
		APIVersion                              string
		APIServerURL                            string
	}
)

func (conf testBootstrap) GetServer() string {
	return conf.Server
}

func (conf testBootstrap) GetDBConnParamenters() map[string]string {
	dbConf := make(map[string]string)
	dbConf["DBHost"] = conf.DBHost
	dbConf["Database"] = conf.Database
	dbConf["DBUser"] = conf.DBUser
	dbConf["DBPass"] = conf.DBPass
	dbConf["DBSSL"] = conf.DBSSL
	return dbConf
}

func (conf testBootstrap) GetBaseDir() string {
	return conf.BaseDir
}

func (conf testBootstrap) GetResourcesDir() string {
	return conf.ResourcesDir
}

func (conf testBootstrap) GetPublicDir() string {
	return conf.PublicDir
}

func (conf testBootstrap) GetLogFile() string {
	return conf.LogFile
}

func (conf testBootstrap) GetLogLevel() int {
	return conf.LogLevel
}

func (conf testBootstrap) IsAutoreloadOn() bool {
	return conf.Autoreload
}
