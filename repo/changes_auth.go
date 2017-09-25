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

package repo

import (
	"log"
	"reflect"

	"github.com/adrianpk/pulap/models"
)

// UserChanges - Creates a map ([string]interface{}) including al changing field.
func UserChanges(user *models.User, reference models.User) map[string]string {
	changes := make(map[string]string)
	if reference.Username.String != user.Username.String {
		changes["username"] = ":username"
	}
	if reference.PasswordHash != user.PasswordHash {
		changes["password_hash"] = ":password_hash"
	}
	if reference.Email.String != user.Email.String {
		changes["email"] = ":email"
	}
	if reference.FirstName.String != user.FirstName.String {
		changes["first_name"] = ":first_name"
	}
	if reference.MiddleNames.String != user.MiddleNames.String {
		changes["middle_names"] = ":middle_names"
	}
	if reference.LastName.String != user.LastName.String {
		changes["last_name"] = ":last_name"
	}
	if !reflect.DeepEqual(reference.Card, user.Card) {
		if isJSON(user.Card.String()) {
			changes["card"] = ":card"
		}
	}
	if !reflect.DeepEqual(reference.Annotations, user.Annotations) {
		if isJSON(user.Annotations.String()) {
			changes["annotations"] = ":annotations"
		}
	}
	if reference.StartedAt.Time != user.StartedAt.Time {
		if true {
			changes["started_at"] = ":started_at"
		}
	}
	if reference.Geolocation.Point.String() != user.Geolocation.Point.String() {
		changes["geolocation"] = ":geolocation"
	}
	if reference.IsActive.Bool != user.IsActive.Bool {
		changes["is_active"] = ":is_active"
	}
	if reference.IsLogicalDeleted.Bool != user.IsLogicalDeleted.Bool {
		changes["is_logical_deleted"] = ":is_logical_deleted"
	}
	return changes
}

// ProfileChanges - Creates a map ([string]interface{}) including al changing field.
func ProfileChanges(profile *models.Profile, reference models.Profile) map[string]string {
	log.Print(profile.Avatar)
	changes := make(map[string]string)
	if reference.Name.String != profile.Name.String {
		changes["name"] = ":name"
	}
	if reference.Email != profile.Email {
		changes["email"] = ":email"
	}
	if reference.Description.String != profile.Description.String {
		changes["description"] = ":description"
	}
	if reference.Bio.String != profile.Bio.String {
		changes["bio"] = ":bio"
	}
	if reference.Moto.String != profile.Moto.String {
		changes["moto"] = ":moto"
	}
	if reference.Website.String != profile.Website.String {
		changes["website"] = ":website"
	}
	if reference.AnniversaryDate.Time != profile.AnniversaryDate.Time {
		if true {
			changes["anniversary_date"] = ":anniversary_date"
		}
	}
	if len(profile.Avatar.ByteSlice) > 0 && !reflect.DeepEqual(reference.Avatar.ByteSlice, profile.Avatar.ByteSlice) {
		changes["avatar"] = ":avatar"
	}
	if reference.AvatarURI.String != profile.AvatarURI.String {
		changes["avatar_uri"] = ":avatar_uri"
	}
	if reference.HeaderURI.String != profile.HeaderURI.String {
		changes["header_uri"] = ":header_uri"
	}
	if !reflect.DeepEqual(reference.Annotations, profile.Annotations) {
		if isJSON(profile.Annotations.String()) {
			changes["annotations"] = ":annotations"
		}
	}
	if reference.StartedAt.Time != profile.StartedAt.Time {
		if true {
			changes["started_at"] = ":started_at"
		}
	}
	if reference.Geolocation.Point.String() != profile.Geolocation.Point.String() {
		changes["geolocation"] = ":geolocation"
	}
	if reference.IsActive.Bool != profile.IsActive.Bool {
		changes["is_active"] = ":is_active"
	}
	if reference.IsLogicalDeleted.Bool != profile.IsLogicalDeleted.Bool {
		changes["is_logical_deleted"] = ":is_logical_deleted"
	}
	return changes
}

// OrganizationChanges - Creates a map ([string]interface{}) including al changing field.
func OrganizationChanges(organization *models.Organization, reference models.Organization) map[string]interface{} {
	changes := make(map[string]interface{})
	if reference.Name.String != organization.Name.String {
		changes["name"] = ":name"
	}
	if reference.Description.String != organization.Description.String {
		changes["description"] = ":description"
	}
	if !reflect.DeepEqual(reference.Annotations, organization.Annotations) {
		if isJSON(organization.Annotations.String()) {
			changes["annotations"] = ":annotations"
		}
	}
	if reference.StartedAt.Time != organization.StartedAt.Time {
		if true {
			changes["started_at"] = ":started_at"
		}
	}
	if reference.Geolocation.Point.String() != organization.Geolocation.Point.String() {
		changes["geolocation"] = ":geolocation"
	}
	if reference.CreatedBy.String != organization.CreatedBy.String {
		if organization.CreatedBy.String != "" {
			changes["created_by"] = ":created_by"
		}
	}
	if reference.IsActive.Bool != organization.IsActive.Bool {
		changes["is_active"] = ":is_active"
	}
	if reference.IsLogicalDeleted.Bool != organization.IsLogicalDeleted.Bool {
		changes["is_logical_deleted"] = ":is_logical_deleted"
	}
	if reference.UpdatedAt.Time != organization.UpdatedAt.Time {
		if true {
			changes["updated_at"] = ":updated_at"
		}
	}
	if reference.UserID != organization.UserID {
		changes["user_id"] = ":user_id"
	}
	return changes
}

// ResourceChanges - Creates a map ([string]interface{}) including al changing field.
func ResourceChanges(resource *models.Resource, reference models.Resource) map[string]string {
	changes := make(map[string]string)
	if reference.Name.String != resource.Name.String {
		changes["name"] = ":name"
	}
	if reference.Description.String != resource.Description.String {
		changes["description"] = ":description"
	}
	if reference.IsActive.Bool != resource.IsActive.Bool {
		changes["is_active"] = ":is_active"
	}
	if reference.IsLogicalDeleted.Bool != resource.IsLogicalDeleted.Bool {
		changes["is_logical_deleted"] = ":is_logical_deleted"
	}
	if reference.UpdatedAt.Time != resource.UpdatedAt.Time {
		if true {
			changes["updated_at"] = ":updated_at"
		}
	}
	if reference.OrganizationID != resource.OrganizationID {
		changes["organization_id"] = ":organization_id"
	}
	return changes
}

// PermissionChanges - Creates a map ([string]interface{}) including al changing field.
func PermissionChanges(permission *models.Permission, reference models.Permission) map[string]string {
	changes := make(map[string]string)
	if reference.Name.String != permission.Name.String {
		changes["name"] = ":name"
	}
	if reference.Description.String != permission.Description.String {
		changes["description"] = ":description"
	}
	if reference.IsActive.Bool != permission.IsActive.Bool {
		changes["is_active"] = ":is_active"
	}
	if reference.IsLogicalDeleted.Bool != permission.IsLogicalDeleted.Bool {
		changes["is_logical_deleted"] = "is_logical_deleted"
	}
	if reference.UpdatedAt.Time != permission.UpdatedAt.Time {
		if true {
			changes["updated_at"] = ":updated_at"
		}
	}
	if reference.OrganizationID != permission.OrganizationID {
		changes["organization_id"] = ":organization_id"
	}
	return changes
}

// RoleChanges - Creates a map ([string]interface{}) including al changing field.
func RoleChanges(role *models.Role, reference models.Role) map[string]string {
	changes := make(map[string]string)
	if reference.Name.String != role.Name.String {
		changes["name"] = ":name"
	}
	if reference.Description.String != role.Description.String {
		changes["description"] = ":description"
	}
	if reference.IsActive.Bool != role.IsActive.Bool {
		changes["is_active"] = ":is_active"
	}
	if reference.IsLogicalDeleted.Bool != role.IsLogicalDeleted.Bool {
		changes["is_logical_deleted"] = ":is_logical_deleted"
	}
	if reference.UpdatedAt.Time != role.UpdatedAt.Time {
		if true {
			changes["updated_at"] = ":updated_at"
		}
	}
	if reference.OrganizationID != role.OrganizationID {
		changes["organization_id"] = ":organization_id"
	}
	return changes
}

// ResourcePermissionChanges - Creates a map ([string]interface{}) including al changing field.
func ResourcePermissionChanges(resourcePermission *models.ResourcePermission, reference models.ResourcePermission) map[string]string {
	changes := make(map[string]string)
	if reference.Name.String != resourcePermission.Name.String {
		changes["name"] = ":name"
	}
	if reference.Description.String != resourcePermission.Description.String {
		changes["description"] = ":description"
	}
	if reference.IsActive.Bool != resourcePermission.IsActive.Bool {
		changes["is_active"] = ":is_active"
	}
	if reference.IsLogicalDeleted.Bool != resourcePermission.IsLogicalDeleted.Bool {
		changes["is_logical_deleted"] = ":is_logical_deleted"
	}
	if reference.UpdatedAt.Time != resourcePermission.UpdatedAt.Time {
		if true {
			changes["updated_at"] = ":updated_at"
		}
	}
	if resourcePermission.OrganizationID.String != "" && reference.OrganizationID != resourcePermission.OrganizationID {
		changes["organization_id"] = ":organization_id"
	}
	if resourcePermission.ResourceID.String != "" && reference.ResourceID != resourcePermission.ResourceID {
		changes["resource_id"] = ":resource_id"
	}
	if resourcePermission.PermissionID.String != "" && reference.PermissionID != resourcePermission.PermissionID {
		changes["permission_id"] = ":permission_id"
	}
	return changes
}

// RolePermissionChanges - Creates a map ([string]interface{}) including al changing field.
func RolePermissionChanges(rolePermission *models.RolePermission, reference models.RolePermission) map[string]string {
	changes := make(map[string]string)
	if reference.Name.String != rolePermission.Name.String {
		changes["name"] = ":name"
	}
	if reference.Description.String != rolePermission.Description.String {
		changes["description"] = ":description"
	}
	if reference.IsActive.Bool != rolePermission.IsActive.Bool {
		changes["is_active"] = ":is_active"
	}
	if reference.IsLogicalDeleted.Bool != rolePermission.IsLogicalDeleted.Bool {
		changes["is_logical_deleted"] = ":is_logical_deleted"
	}
	if reference.UpdatedAt.Time != rolePermission.UpdatedAt.Time {
		if true {
			changes["updated_at"] = "updated_at"
		}
	}
	if rolePermission.OrganizationID.String != "" && reference.OrganizationID != rolePermission.OrganizationID {
		changes["organization_id"] = ":organization_id"
	}
	if rolePermission.RoleID.String != "" && reference.RoleID != rolePermission.RoleID {
		changes["role_id"] = ":role_id"
	}
	if rolePermission.PermissionID.String != "" && reference.PermissionID != rolePermission.PermissionID {
		changes["permission_id"] = ":permission_id"
	}
	return changes
}

// UserRoleChanges - Creates a map ([string]interface{}) including al changing field.
func UserRoleChanges(userRole *models.UserRole, reference models.UserRole) map[string]string {
	changes := make(map[string]string)
	if reference.Name.String != userRole.Name.String {
		changes["name"] = ":name"
	}
	if reference.Description.String != userRole.Description.String {
		changes["description"] = ":description"
	}
	if reference.IsActive.Bool != userRole.IsActive.Bool {
		changes["is_active"] = ":is_active"
	}
	if reference.IsLogicalDeleted.Bool != userRole.IsLogicalDeleted.Bool {
		changes["is_logical_deleted"] = ":is_logical_deleted"
	}
	if reference.UpdatedAt.Time != userRole.UpdatedAt.Time {
		if true {
			changes["updated_at"] = ":updated_at"
		}
	}
	if userRole.OrganizationID.String != "" && reference.OrganizationID != userRole.OrganizationID {
		changes["organization_id"] = ":organization_id"
	}
	if userRole.UserID.String != "" && reference.UserID != userRole.UserID {
		changes["user_id"] = ":user_id"
	}
	if userRole.RoleID.String != "" && reference.RoleID != userRole.RoleID {
		changes["role_id"] = ":role_id"
	}
	return changes
}
