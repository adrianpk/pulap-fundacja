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

package models

import (
	"github.com/adrianpk/pulap/types"

	sqlxtypes "github.com/jmoiron/sqlx/types"
	"github.com/markbates/pop/nulls"
)

type (
	// User - User model
	User struct {
		IdentifiableModel
		Username             nulls.String       `db:"username" json:"username"`
		Password             string             `json:"password, omitempty" schema:"password"`
		PasswordConfirmation string             `json:"-" schema:"password-confirmation"`
		PasswordHash         string             `db:"password_hash" json:"passwordHash, omitempty"`
		Email                nulls.String       `db:"email" json:"email"`
		FirstName            nulls.String       `db:"first_name" json:"firstName" schema:"first-name"`
		MiddleNames          nulls.String       `db:"middle_names" json:"middleNames" schema:"middle-names"`
		LastName             nulls.String       `db:"last_name" json:"lastName" schema:"last-name"`
		Card                 sqlxtypes.JSONText `db:"card" json:"card"`
		AnnotableModel
		GeolocalizableModel
		AuditableModel
		ValidableDate
	}

	// Profile - Profile model
	Profile struct {
		IdentifiableModel
		Email           nulls.String       `db:"email" json:"email"`
		Bio             nulls.String       `db:"bio" json:"bio"`
		Moto            nulls.String       `db:"moto" json:"moto"`
		Website         nulls.String       `db:"website" json:"website"`
		AnniversaryDate nulls.Time         `db:"anniversary_date" json:"anniversaryDate" schema:"anniversary-date"`
		Avatar          nulls.ByteSlice    `db:"avatar" json:"-"`
		AvatarBase64    string             `json:"avatar, omitempty" schema:"avatar-base-64"`
		AvatarURI       nulls.String       `db:"avatar_uri" json:"avatarURL" schema:"avatar-uri"`
		HeaderURI       nulls.String       `db:"header_uri" json:"headerURL" schema:"header-uri"`
		Card            sqlxtypes.JSONText `db:"card" json:"card"`
		Geolocation     types.NullPoint    `db:"geolocation" json:"geolocation"`
		//UserUsername    nulls.String       `db:"username" json:"username" schema:"user-username"`
		UserID         nulls.String `db:"user_id" json:"userID, omitempty" schema:"user-id"`
		OrganizationID nulls.String `db:"organization_id" json:"organizationID, omitempty" schema:"organization-id"`
		AnnotableModel
		GeolocalizableModel
		AuditableModel
		ValidableDate
	}

	// Organization - Resource model
	Organization struct {
		IdentifiableModel
		UserUsername nulls.String `db:"user_username" json:"userUsername, omitempty" schema:"user-username"`
		UserID       nulls.String `db:"user_id" json:"userID, omitempty" schema:"user-id"`
		AnnotableModel
		GeolocalizableModel
		AuditableModel
		ValidableDate
	}

	// Resource - Resource model
	Resource struct {
		IdentifiableModel
		Tag            nulls.String `db:"tag" json:"tag, omitempty"`
		OrganizationID nulls.String `db:"organization_id" json:"organizationID, omitempty" schema:"organization-id"`
		AuditableModel
		ValidableDate
	}

	// Permission - Resource models
	Permission struct {
		IdentifiableModel
		OrganizationName nulls.String `db:"organization_name" json:"organizationName, omitempty" schema:"organization-name"`
		OrganizationID   nulls.String `db:"organization_id" json:"organizationID, omitempty" schema:"organization-id"`
		AuditableModel
		ValidableDate
	}

	// Role - Resource model
	Role struct {
		IdentifiableModel
		OrganizationID nulls.String `db:"organization_id" json:"organizationID, omitempty" schema:"organization-id"`
		AuditableModel
		ValidableDate
	}

	// ResourcePermission - ResourcePermission model
	ResourcePermission struct {
		IdentifiableModel
		OrganizationID nulls.String `db:"organization_id" json:"organizationID, omitempty" schema:"organization-id"`
		ResourceID     nulls.String `db:"resource_id" json:"resourceID, omitempty" schema:"resource-id"`
		PermissionID   nulls.String `db:"permission_id" json:"permissionID, omitempty" schema:"permission-id"`
		AuditableModel
		ValidableDate
	}

	// RolePermission - RolePermission model
	RolePermission struct {
		IdentifiableModel
		OrganizationID nulls.String `db:"organization_id" json:"organizationID, omitempty" schema:"organization-id"`
		RoleID         nulls.String `db:"role_id" json:"roleID, omitempty" schema:"role-id"`
		PermissionID   nulls.String `db:"permission_id" json:"permissionID, omitempty" schema:"permission-id"`
		AuditableModel
		ValidableDate
	}

	// UserRole - UserRole model
	UserRole struct {
		IdentifiableModel
		OrganizationID nulls.String `db:"organization_id" json:"organizationID, omitempty" schema:"organization-id"`
		UserID         nulls.String `db:"user_id" json:"userID, omitempty" schema:"user-id"`
		RoleID         nulls.String `db:"role_id" json:"roleID, omitempty" schema:"role-id"`
		AuditableModel
		ValidableDate
	}

	// PropertiesSet - PropertiesSet model
	PropertiesSet struct {
		IdentifiableModel
		Position nulls.Int64  `db:"position" json:"position, omitempty" schema:"position-id"`
		HolderID nulls.String `db:"holder_id" json:"holderID, omitempty" schema:"holder-id"`
		AuditableModel
		ValidableDate
	}

	// Property - Property model
	Property struct {
		IdentifiableModel
		StringValue      nulls.String    `db:"string_value" json:"stringValue, omitempty" schema:"string-value"`
		IntValue         nulls.Int64     `db:"int_value" json:"intValue, omitempty" schema:"int-value"`
		FloatValue       nulls.Float64   `db:"float_value" json:"floatValue, omitempty" schema:"float-value"`
		BooleanValue     nulls.Bool      `db:"boolean_value" json:"booleanValue, omitempty" schema:"boolean-value"`
		TimestampValue   nulls.Time      `db:"timestamp_value" json:"timestampValue, omitempty" schema:"timestamp-value"`
		GeolocationValue types.NullPoint `db:"geolocation_value" json:"geolocationValue" schema:"geolocation-value"`
		ValueType        nulls.String    `db:"value_type" json:"valueType, omitempty" schema:"value-type"`
		Position         nulls.Int64     `db:"position" json:"position, omitempty"`
		PropertiesSetID  nulls.String    `db:"properties_set_id" json:"propertiesSetID, omitempty" schema:"properties-set-id"`
		AuditableModel
		ValidableDate
	}
	// Plan - Plan model
	Plan struct {
		IdentifiableModel
		EndsAt             nulls.Time   `db:"ends_at" json:"endsAt, omitempty" schema:"ends-at"`
		PlanSubscriptionID nulls.String `db:"plans_subscripotions_id" json:"planSubscriptionID, omitempty" schema:"plan-subscription-id"`
		AuditableModel
		ValidableDate
	}

	// PlanSubscription - PlanSubscription model
	PlanSubscription struct {
		IdentifiableModel
		EndsAt         nulls.Time   `db:"ends_at" json:"endsAt, omitempty" schema:"ends-at"`
		OrganizationID nulls.String `db:"organization_id" json:"organizationID, omitempty" schema:"organization-id"`
		UserID         nulls.String `db:"user_id" json:"userID, omitempty" schema:"user-id"`
		PlanID         nulls.String `db:"plan_id" json:"planID, omitempty" schema:"plan-id"`
		AuditableModel
		ValidableDate
	}

	// Album - Album model
	Album struct {
		ID          nulls.Int64  `db:"id" json:"id, omitempty"`
		Name        nulls.String `db:"name" json:"name"`
		Description nulls.String `db:"description" json:"description"`
		AuditableModel
	}

	// Image - Image model
	Image struct {
		ID              nulls.Int64  `db:"id" json:"id, omitempty"`
		Name            nulls.String `db:"name" json:"name"`
		Description     nulls.String `db:"description" json:"description"`
		Base64          nulls.String `db:"base-64" json:"base64" schema:"base-64"`
		ThumbnailBase64 nulls.String `db:"thumbnail-base-64" json:"thumbnailBase64" schema:"thumbnail-base-64"`
		Geolocation     `json:"geolocation"`
		AuditableModel
	}

	// Geolocation - Geolocation model
	Geolocation struct {
		Type        nulls.String `db:"type" json:"id, omitempty"`
		Coordinates []float64    `db:"coordinates" json:"coordinates, omitempty"`
	}

	// UserIdentifier - UserIdentifier model
	UserIdentifier struct {
		ID       nulls.Int64  `db:"id" json:"id, omitempty"`
		Username nulls.String `db:"username" json:"username, omitempty"`
	}

	// AgencyIdentifier - AgencyIdentifier model
	AgencyIdentifier struct {
		ID   nulls.Int64  `db:"id" json:"id, omitempty"`
		Name nulls.String `db:"name" json:"name, omitempty"`
	}

	// ImageResource struct {
	// 	Name        nulls.String       `json:"name"`
	// 	Description nulls.String       `json:"description"`
	// 	Encoded     EncodedModel `json:"data"`
	// }
	// // EncodedModel model for encoded file resources
	// EncodedModel struct {
	// 	Base64 nulls.String `json:"base64Data"`
	// }
)
