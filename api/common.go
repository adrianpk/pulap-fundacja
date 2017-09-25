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

package api

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"

	"github.com/adrianpk/pulap/app"
	"github.com/adrianpk/pulap/bootstrap"
	"github.com/adrianpk/pulap/logger"
	"github.com/adrianpk/pulap/models"
	"github.com/adrianpk/pulap/repo"

	"github.com/twinj/uuid"
)

func isUUID(id string) bool {
	_, err := uuid.Parse(id)
	if err != nil {
		return false
	}
	return true
}

func loggedInUserID(r *http.Request) string {
	userID, err := sessionUserID(r)
	if err != nil {
		logger.Dump(err)
		return ""
	}
	return userID
}

func sessionUserID(r *http.Request) (string, error) {
	claims, ok := r.Context().Value(bootstrap.UserCtxKey).(bootstrap.AppClaims)
	if ok {
		userID := claims.UserID
		return userID, nil
	}
	return "", app.ErrNotLoggedIn
}

func sessionUser(r *http.Request) (user models.User, err error) {
	user = models.User{}
	userRepo, err := repo.MakeUserRepository()
	if err != nil {
		return user, err
	}
	userID, err := sessionUserID(r)
	if err != nil {
		return user, err
	}
	return userRepo.Get(userID)
}

func userFromSessionID(r *http.Request) (models.User, error) {
	claims, ok := r.Context().Value(bootstrap.UserCtxKey).(bootstrap.AppClaims)
	if ok {
		userID := claims.UserID
		userRepo, err := repo.MakeUserRepository()
		if err != nil {
			return models.User{}, err
		}
		return userRepo.Get(userID)
	}
	return models.User{}, nil
}

func profileFromSessionID(r *http.Request) (profile models.Profile, err error) {
	user, err := userFromSessionID(r)
	if err != nil {
		logger.Dump(err)
		return models.Profile{}, err
	}
	profileRepo, err := repo.MakeProfileRepository()
	if err != nil {
		logger.Dump(err)
		return models.Profile{}, err
	}
	profile, err = profileRepo.GetByUserID(user.ID.String)
	if err != nil {
		logger.Dump(err)
		return models.Profile{}, err
	}
	return profile, nil
}

func verifyOwnership(ownerID string, r *http.Request) error {
	userID, err := sessionUserID(r)
	if err != nil {
		logger.Debugf("[temp][verifyOwnership] Error: %v", err)
		return err
	}
	if ownerID != userID {
		return app.ErrOwnerOnlyCanManage
	}
	return nil

}

func verifyID(update models.IdentifiableModel, current models.IdentifiableModel) error {
	if update.ID.String != current.ID.String {
		return app.ErrEntityInvalidData
	}
	return nil
}

// Get models by ID
func getUser(userID string) (models.User, error) {
	userRepo, err := repo.MakeUserRepository()
	if err != nil {
		return models.User{}, err
	}
	return userRepo.Get(userID)
}

func getOrganization(orgID string) (models.Organization, error) {
	orgRepo, err := repo.MakeOrganizationRepository()
	if err != nil {
		return models.Organization{}, err
	}
	return orgRepo.Get(orgID)
}

func getResource(resID string) (models.Resource, error) {
	resRepo, err := repo.MakeResourceRepository()
	if err != nil {
		return models.Resource{}, err
	}
	return resRepo.Get(resID)
}

func getRole(roleID string) (models.Role, error) {
	roleRepo, err := repo.MakeRoleRepository()
	if err != nil {
		return models.Role{}, err
	}
	return roleRepo.Get(roleID)
}

func getPermission(permissionID string) (models.Permission, error) {
	permRepo, err := repo.MakePermissionRepository()
	if err != nil {
		return models.Permission{}, err
	}
	return permRepo.Get(permissionID)
}

func getPlan(planID string) (models.Plan, error) {
	planRepo, err := repo.MakePlanRepository()
	if err != nil {
		return models.Plan{}, err
	}
	return planRepo.Get(planID)
}

func debugRequestHeader(r *http.Request) {
	logger.Debugf("Request header: %v", r.Header)
}

func debugRequestBody(r *http.Request) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	str := buf.String()
	logger.Debugf("Request body", str)
}

func debugRequest(r *http.Request) {
	logger.Debug(formatRequest(r))
}

// formatRequest generates ascii representation of a request
func formatRequest(r *http.Request) string {
	// Create return string
	var request []string
	// Add the request string
	url := fmt.Sprintf("%v %v %v", r.Method, r.URL, r.Proto)
	request = append(request, url)
	// Add the host
	request = append(request, fmt.Sprintf("Host: %v", r.Host))
	// Loop through headers
	for name, headers := range r.Header {
		name = strings.ToLower(name)
		for _, h := range headers {
			request = append(request, fmt.Sprintf("%v: %v", name, h))
		}
	}
	// If this is a POST, add post data
	if r.Method == "POST" || r.Method == "PUT" || r.Method == "PATCH" {
		r.ParseForm()
		request = append(request, "\n")
		request = append(request, r.Form.Encode())
	}
	// Return the request as a string
	return strings.Join(request, "\n")
}
