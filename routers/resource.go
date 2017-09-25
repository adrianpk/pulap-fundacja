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

package routers

import (
	"github.com/adrianpk/pulap/bootstrap"
	"github.com/adrianpk/pulap/api"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
)

// InitAPIResourceRoutes - Initialize API router for resources.
func InitAPIResourceRoutes() *mux.Router {
	// Paths
	resourcesPath := "/organizations/{organization}/resources"
	// Router
	resourceRouter := apiV1Router.PathPrefix(resourcesPath).Subrouter()
	// Resource
	resourceRouter.HandleFunc("", api.GetResources).Methods("GET")
	resourceRouter.HandleFunc("", api.CreateResource).Methods("POST")
	resourceRouter.HandleFunc("/{resource}", api.GetResource).Methods("GET")
	resourceRouter.HandleFunc("/{resource}", api.UpdateResource).Methods("PUT")
	resourceRouter.HandleFunc("/{resource}", api.DeleteResource).Methods("DELETE")
	// Middleware
	apiV1Router.Handle(resourcesPath, negroni.New(
		negroni.HandlerFunc(bootstrap.Authorize),
		negroni.Wrap(resourceRouter),
	))
	return resourceRouter
}
