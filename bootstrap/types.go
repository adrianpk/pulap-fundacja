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

type (
	// Configuration - Configuration interface
	Configuration interface {
		GetServer() string
		GetDBConnParamenters() map[string]string // DBHost, Database, DBUser, DBPass, DBSSL string
		GetBaseDir() string
		GetResourcesDir() string
		GetPublicDir() string
		GetLogFile() string
		GetLogLevel() int
		IsAutoreloadOn() bool
	}
)

type (
	configuration struct {
		Server                                  string
		DBHost, Database, DBUser, DBPass, DBSSL string
		BaseDir                                 string
		ResourcesDir                            string
		PublicDir                               string
		LogLevel                                int
		LogFile                                 string
		Autoreload                              bool
	}
)

func (conf configuration) GetServer() string {
	return conf.Server
}

func (conf configuration) GetDBConnParamenters() map[string]string {
	dbConf := make(map[string]string)
	dbConf["DBHost"] = conf.DBHost
	dbConf["Database"] = conf.Database
	dbConf["DBUser"] = conf.DBUser
	dbConf["DBPass"] = conf.DBPass
	dbConf["DBSSL"] = conf.DBSSL
	return dbConf
}

func (conf configuration) GetBaseDir() string {
	return conf.BaseDir
}

func (conf configuration) GetResourcesDir() string {
	return conf.ResourcesDir
}

func (conf configuration) GetPublicDir() string {
	return conf.PublicDir
}

func (conf configuration) GetLogFile() string {
	return conf.LogFile
}

func (conf configuration) GetLogLevel() int {
	return conf.LogLevel
}

func (conf configuration) IsAutoreloadOn() bool {
	return conf.Autoreload
}
