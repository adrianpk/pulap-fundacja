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

import "github.com/adrianpk/pulap/models"

// PropertiesSetChanges - Creates a map ([string]interface{}) including al changing field.
func PropertiesSetChanges(propertiesSet *models.PropertiesSet, reference models.PropertiesSet) map[string]string {
	changes := make(map[string]string)
	if reference.Name.String != propertiesSet.Name.String {
		changes["name"] = ":name"
	}
	if reference.Description.String != propertiesSet.Description.String {
		changes["description"] = ":description"
	}
	if reference.Position.Int64 != propertiesSet.Position.Int64 {
		changes["position"] = ":position"
	}
	if reference.HolderID.String != propertiesSet.HolderID.String {
		changes["holder_id"] = ":holder_id"
	}
	if reference.IsActive.Bool != propertiesSet.IsActive.Bool {
		changes["is_active"] = ":is_active"
	}
	if reference.IsLogicalDeleted.Bool != propertiesSet.IsLogicalDeleted.Bool {
		changes["is_logical_deleted"] = ":is_logical_deleted"
	}
	if reference.UpdatedAt.Time != propertiesSet.UpdatedAt.Time {
		if true {
			changes["updated_at"] = ":updated_at"
		}
	}
	return changes
}

// PropertyChanges - Creates a map ([string]interface{}) including al changing field.
func PropertyChanges(property *models.Property, reference models.Property) map[string]string {
	changes := make(map[string]string)
	if reference.Name.String != property.Name.String {
		changes["name"] = ":name"
	}
	if reference.Description.String != property.Description.String {
		changes["description"] = ":description"
	}
	if reference.StringValue.String != property.StringValue.String {
		changes["string_value"] = ":string_value"
	}
	if reference.IntValue.Int64 != property.IntValue.Int64 {
		changes["int_value"] = ":int_value"
	}
	if reference.FloatValue.Float64 != property.FloatValue.Float64 {
		changes["float_value"] = ":float_value"
	}
	if reference.BooleanValue.Bool != property.BooleanValue.Bool {
		changes["boolean_value"] = ":boolean_value"
	}
	if reference.TimestampValue.Time != property.TimestampValue.Time {
		changes["timestamp_value"] = ":timestamp_value"
	}
	if reference.ValueType.String != property.ValueType.String {
		changes["value_type"] = ":value_type"
	}
	if reference.Position.Int64 != property.Position.Int64 {
		changes["position"] = ":position"
	}
	if reference.IsLogicalDeleted.Bool != property.IsLogicalDeleted.Bool {
		changes["is_logical_deleted"] = ":is_logical_deleted"
	}
	if reference.UpdatedAt.Time != property.UpdatedAt.Time {
		if true {
			changes["updated_at"] = ":updated_at"
		}
	}
	return changes
}
