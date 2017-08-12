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
	"github.com/adrianpk/fundacja/controllers"

	"github.com/gorilla/mux"
)

// InitAPIPropertiesSetRouter - Initialize API router for properties sets.
func InitAPIPropertiesSetRouter() *mux.Router {
	// Paths
	propertiesSetsPath := "/api/v1/holders/{holder}/properties-sets"
	// Router
	propertiesSetRouter := apiV1Router.PathPrefix(propertiesSetsPath).Subrouter()
	// Resource
	propertiesSetRouter.HandleFunc("", controllers.GetPropertiesSets).Methods("GET")
	propertiesSetRouter.HandleFunc("", controllers.CreatePropertiesSet).Methods("POST")
	propertiesSetRouter.HandleFunc("/{properties-set}", controllers.GetPropertiesSet).Methods("GET")
	propertiesSetRouter.HandleFunc("/{properties-set}", controllers.UpdatePropertiesSet).Methods("PUT")
	propertiesSetRouter.HandleFunc("/{properties-set}", controllers.DeletePropertiesSet).Methods("DELETE")
	// // Resource
	// propertiesSetRouter.HandleFunc("/{properties-set}/properties", controllers.GetProperties).Methods("GET")
	// propertiesSetRouter.HandleFunc("/{properties-set}/properties", controllers.CreateProperty).Methods("POST")
	// propertiesSetRouter.HandleFunc("/{properties-set}/properties/{property}", controllers.GetProperty).Methods("GET")
	// propertiesSetRouter.HandleFunc("/{properties-set}/properties/{property}", controllers.UpdateProperty).Methods("PUT")
	// propertiesSetRouter.HandleFunc("/{properties-set}/properties/{property}", controllers.DeleteProperty).Methods("DELETE")
	return propertiesSetRouter
}
