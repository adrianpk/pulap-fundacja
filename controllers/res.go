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

package controllers

import (
	"github.com/adrianpk/fundacja/models"
)

type (
	// UserResource For Post - /users/signup
	UserResource struct {
		Data models.User `json:"data"`
	}

	// UsersResource For Get - /users/r
	UsersResource struct {
		Data []models.User `json:"data"`
	}

	// LoginResource For Post - /users/login
	LoginResource struct {
		Data LoginModel `json:"data"`
	}

	// LoginModel for authentication
	LoginModel struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// AuthUserResource for authorized user Post - /users/login
	AuthUserResource struct {
		Data AuthUserModel `json:"data"`
	}

	// AuthUserModel for authorized user with access token
	AuthUserModel struct {
		User  models.User `json:"user"`
		Token string      `json:"token"`
	}

	// AvatarResource for authorized user with access token
	AvatarResource struct {
		Data AvatarModel `json:"data"`
	}

	// AvatarModel - Resource
	AvatarModel struct {
		ProfileID string `json:"profileID"`
		Base64    string `json:"base64"`
	}

	// // ImageResource model for Image Upload
	// ImageResource struct {
	// 	Name        string       `json:"name"`
	// 	Description string       `json:"description"`
	// 	Encoded     EncodedModel `json:"data"`
	// }
	// // EncodedModel model for encoded file resources
	// EncodedModel struct {
	// 	Base64 string `json:"base64Data"`
	// }

	// ProfileResource - Resource
	ProfileResource struct {
		Data models.Profile `json:"data"`
	}

	// OrganizationResource - Resource
	OrganizationResource struct {
		Data models.Organization `json:"data"`
	}

	// OrganizationsResource - Resource
	OrganizationsResource struct {
		Data []models.Organization `json:"data"`
	}

	// ResourceResource - Resource
	ResourceResource struct {
		Data models.Resource `json:"data"`
	}

	// ResourcesResource - Resource
	ResourcesResource struct {
		Data []models.Resource `json:"data"`
	}

	// PermissionResource  - Resource
	PermissionResource struct {
		Data models.Permission `json:"data"`
	}

	// PermissionsResource - Resource
	PermissionsResource struct {
		Data []models.Permission `json:"data"`
	}

	// RoleResource - Resource
	RoleResource struct {
		Data models.Role `json:"data"`
	}

	// RolesResource For Get - /roles
	RolesResource struct {
		Data []models.Role `json:"data"`
	}

	// ResourcePermissionResource - Resource
	ResourcePermissionResource struct {
		Data models.ResourcePermission `json:"data"`
	}

	// ResourcePermissionsResource  - Resource
	ResourcePermissionsResource struct {
		Data []models.ResourcePermission `json:"data"`
	}

	// RolePermissionResource - Resource
	RolePermissionResource struct {
		Data models.RolePermission `json:"data"`
	}

	// RolePermissionsResource - Resource
	RolePermissionsResource struct {
		Data []models.RolePermission `json:"data"`
	}

	// UserRoleResource  - Resource
	UserRoleResource struct {
		Data models.UserRole `json:"data"`
	}

	// UserRolesResource - Resource
	UserRolesResource struct {
		Data []models.UserRole `json:"data"`
	}

	// PropertiesSetsResource - Resource
	PropertiesSetsResource struct {
		Data []models.PropertiesSet `json:"data"`
	}

	// PropertiesSetResource  - Resource
	PropertiesSetResource struct {
		Data models.PropertiesSet `json:"data"`
	}

	// PropertiesResource - Resource
	PropertiesResource struct {
		Data []models.Property `json:"data"`
	}

	// PropertyResource - Resource
	PropertyResource struct {
		Data models.Property `json:"data"`
	}

	// PlanSubscriptionsResource  - Resource
	PlanSubscriptionsResource struct {
		Data []models.PlanSubscription `json:"data"`
	}

	// PlanSubscriptionResource - Resource
	PlanSubscriptionResource struct {
		Data models.PlanSubscription `json:"data"`
	}

	// PlansResource - Resource
	PlansResource struct {
		Data []models.Plan `json:"data"`
	}

	// PlanResource - Resource
	PlanResource struct {
		Data models.Plan `json:"data"`
	}
)
