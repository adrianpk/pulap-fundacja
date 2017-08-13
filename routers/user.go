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
	"github.com/adrianpk/fundacja/api"
	"github.com/gorilla/mux"
)

// InitAPIUserRouter - Initialize API router for users.
func InitAPIUserRouter() *mux.Router {
	// Paths
	usersPath := "/api/v1/users"
	// Router
	userRouter := apiV1Router.PathPrefix(usersPath).Subrouter()
	// Resource
	userRouter.HandleFunc("", api.GetUsers).Methods("GET")
	userRouter.HandleFunc("", api.CreateUser).Methods("POST")
	userRouter.HandleFunc("/{user}", api.GetUser).Methods("GET")
	userRouter.HandleFunc("/{user}", api.UpdateUser).Methods("PUT")
	userRouter.HandleFunc("/{user}", api.DeleteUser).Methods("DELETE")
	// Resource
	userRouter.HandleFunc("/{user}/profile", api.GetUserProfile).Methods("GET")
	userRouter.HandleFunc("/{user}/profile", api.CreateUserProfile).Methods("POST")
	userRouter.HandleFunc("/{user}/profile", api.UpdateUserProfile).Methods("PUT")
	userRouter.HandleFunc("/{user}/profile", api.DeleteUserProfile).Methods("DELETE")
	// Resource
	userRouter.HandleFunc("/{user}/profile/avatar", api.HandleAvatar)
	return userRouter
}
