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
	"github.com/adrianpk/pulap/controllers"
	"github.com/gorilla/mux"
)

// InitUserRouter - Initialize router for users.
func InitUserRouter() *mux.Router {
	// Paths
	usersPath := "/users" ///{rest:.*}
	// Router
	userRouter := appRouter.PathPrefix(usersPath).Subrouter()
	userRouter.StrictSlash(true)
	// Resource
	userRouter.HandleFunc("/", controllers.IndexUsers).Methods("GET")
	userRouter.HandleFunc("/new", controllers.NewUser).Methods("GET")
	userRouter.HandleFunc("", controllers.CreateUser).Methods("POST")
	userRouter.HandleFunc("/{user}", controllers.ShowUser).Methods("GET")
	userRouter.HandleFunc("/edit/{user}", controllers.EditUser).Methods("GET")
	userRouter.HandleFunc("/{user}", controllers.UpdateUser).Methods("POST")
	userRouter.HandleFunc("/init-delete/{user}", controllers.InitDeleteUser).Methods("POST")
	userRouter.HandleFunc("/delete/{user}", controllers.DeleteUser).Methods("POST")
	// Resource
	userRouter.HandleFunc("/{user}/profile", controllers.GetUserProfile).Methods("GET")
	userRouter.HandleFunc("/{user}/profile", controllers.CreateUserProfile).Methods("POST")
	userRouter.HandleFunc("/{user}/profile", controllers.UpdateUserProfile).Methods("PUT")
	userRouter.HandleFunc("/{user}/profile", controllers.DeleteUserProfile).Methods("DELETE")
	// Resource
	userRouter.HandleFunc("/{user}/profile/avatar", controllers.HandleAvatar)
	// Middleware
	// appRouter.PathPrefix(usersPath).Handler(
	// 	negroni.New(
	// 		//negroni.HandlerFunc(bootstrap.Authorize),
	// 		negroni.Wrap(userRouter),
	// 	))
	return userRouter
}
