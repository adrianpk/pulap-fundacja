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

// InitAPIPlanSubscriptionRouter - Initialize API router for properties sets.
func InitAPIPlanSubscriptionRouter() *mux.Router {
	// Paths
	planSubscriptionPath := "/api/v1/plan-subscriptor/{plan-subscriptor}/plan-subscription"
	// Router
	planSubscriptionRouter := apiV1Router.PathPrefix(planSubscriptionPath).Subrouter()
	// Resource
	planSubscriptionRouter.HandleFunc("", api.GetPlanSubscription).Methods("GET")
	planSubscriptionRouter.HandleFunc("", api.CreatePlanSubscription).Methods("POST")
	planSubscriptionRouter.HandleFunc("", api.UpdatePlanSubscription).Methods("PUT")
	planSubscriptionRouter.HandleFunc("", api.DeletePlanSubscription).Methods("DELETE")
	return planSubscriptionRouter
}
