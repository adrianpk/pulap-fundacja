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
	"fmt"
	"net/http"

	"github.com/adrianpk/pulap/bootstrap"
	"github.com/adrianpk/pulap/logger"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
)

var (
	appRouter     *mux.Router
	apiV1Router   *mux.Router
	apiLoginPath  string
	apiSignupPath string
	loginPath     string
	signupPath    string
)

// GetRouter - Returns app main router
func GetRouter(config bootstrap.Configuration) *mux.Router {
	InitAppRouter()
	InitSignupAndLoginRouter()
	InitSubRouters()
	InitAPIV1Router()
	InitAPIV1SubRouters()
	InitPublicFilesystem(config.GetPublicDir())
	return appRouter
}

// InitAppRouter - Initialize app routes.
func InitAppRouter() {
	appRouter = NewRouter()
	//appRouter.HandleFunc("/", home())
}

// InitAPIV1Router - Get a router for API calls.
func InitAPIV1Router() {
	apiV1Path := "/api/v1"
	apiLoginPath = "/api/v1/login"
	apiSignupPath = "/api/v1/signup"
	apiV1Router = NewRouter()
	InitAPILoginRouter()
	InitAPISignUpRouter()
	// Middleware
	appRouter.PathPrefix(apiV1Path).Handler(
		negroni.New(
			negroni.HandlerFunc(bootstrap.Authorize),
			negroni.Wrap(apiV1Router),
		))
}

// InitAPIV1SubRouters - Initialize API subrouters.
func InitAPIV1SubRouters() {
	InitAPIUserRouter()
	InitAPIOrganizationRouter()
	InitAPIPropertiesSetRouter()
	InitAPIPropertyRouter()
	InitAPIPlanSubscriptionRouter()
	InitAPIPlanRouter()
}

// InitSignupAndLoginRouter - Get a router for API calls.
func InitSignupAndLoginRouter() {
	loginPath = "/login"
	signupPath = "/signup"
	InitLoginRouter()
	InitSignUpRouter()
	// Middleware
	// appRouter.PathPrefix(apiV1Path).Handler(
	// 	negroni.New(
	// 		//negroni.HandlerFunc(bootstrap.Authorize),
	// 		negroni.Wrap(apiV1Router),
	// 	))
}

// InitSubRouters - Initialize subrouters.
func InitSubRouters() {
	InitUserRouter()
	InitOrganizationRouter()
	// InitAPIPropertiesSetRouter()
	// InitAPIPropertyRouter()
	// InitAPIPlanSubscriptionRouter()
	// InitAPIPlanRouter()
}

// InitPublicFilesystem - Initialize public filesystem.
func InitPublicFilesystem(publicDir string) {
	// Paths
	//publicFsPath := "/Users/adrian/go/src/github.com/adrianpk/pulap/resources/public"
	publicPath := "/{rest:.*}"
	// Fileserver
	fsServer := http.FileServer(http.Dir(publicDir))
	// Router
	fsRouter := appRouter.PathPrefix(publicPath).Subrouter()
	// Resources
	fsRouter.Handle(publicPath, fsServer)
}

// InitPublicFilesystems - Initialize public filesystem.
func InitPublicFilesystems(publicDir string) {
	// Paths
	//fsPath := "/Users/adrian/go/src/github.com/adrianpk/pulap/resources/public"
	publicPath := "/"
	// Fileserver
	fsServer := http.FileServer(http.Dir(publicDir))
	// Resources
	appRouter.Handle(publicPath, fsServer)
}

// NewRouter - Creates a new mux.router.
func NewRouter() *mux.Router {
	r := mux.NewRouter()
	r.StrictSlash(false)
	r.KeepContext = true
	return r
}

func home() func(http.ResponseWriter, *http.Request) {
	logger.Debug("Home.....")
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/html")
		head := "<head><title>Fundacja</title></head>"
		body := "<body><div>Fundacja is working!</div></body>"
		w.Write([]byte(fmt.Sprintf("<html>%s%s</html>", head, body)))
	}
}
