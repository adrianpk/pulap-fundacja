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
	"github.com/adrianpk/pulap/api"

	"github.com/gorilla/mux"
)

// InitAPIPropertyRouter - Initialize API router for roles.
func InitAPIPropertyRouter() *mux.Router {
	// Paths
	propertiesPath := "/api/v1/properties-set/{properties-set}/properties"
	// Router
	propertyRouter := apiV1Router.PathPrefix(propertiesPath).Subrouter()
	// Resource
	propertyRouter.HandleFunc("", api.GetProperties).Methods("GET")
	propertyRouter.HandleFunc("", api.CreateProperty).Methods("POST")
	propertyRouter.HandleFunc("/{property}", api.GetProperty).Methods("GET")
	propertyRouter.HandleFunc("/{property}", api.UpdateProperty).Methods("PUT")
	propertyRouter.HandleFunc("/{property}", api.DeleteProperty).Methods("DELETE")
	return propertyRouter
}
