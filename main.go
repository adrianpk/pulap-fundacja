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

package main

//go:generate rm -rf views/views.go
//go:generate go-bindata -pkg views -o views/views.go resources/templates/...
//g.:g....... go-bindata -pkg views -o views/public.go resources/public/...

import (
	"net/http"
	"os"

	"github.com/adrianpk/fundacja/bootstrap"
	"github.com/adrianpk/fundacja/controllers"
	"github.com/adrianpk/fundacja/handler"
	"github.com/adrianpk/fundacja/logger"
)

func main() {
	bootstrap.SetBootParameters(mockBootParameters())
	bootstrap.Boot()
	controllers.Initialize()
	handler := handler.AppHandler(bootstrap.AppConfig)
	server := &http.Server{
		Addr:    bootstrap.AppConfig.GetServer(),
		Handler: handler,
	}
	err := server.ListenAndServe()
	if err != nil {
		logger.Debugf("Error: %s", err)
		panic(err)
	}
	logger.Debug("Listening...")
}

// FIX: Just for framework testing
func mockBootParameters() map[string]string {
	params := make(map[string]string)
	params["env"] = "test"
	params["app_home"] = os.Getenv("FUNDACJA_HOME")
	params["migration"] = "n"
	return params
}
