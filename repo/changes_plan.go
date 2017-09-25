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

// PlanSubscriptionChanges - Creates a map ([string]interface{}) including al changing field.
func PlanSubscriptionChanges(planSubscription *models.PlanSubscription, reference models.PlanSubscription) map[string]string {
	changes := make(map[string]string)
	if reference.Name.String != planSubscription.Name.String {
		changes["name"] = ":name"
	}
	if reference.Description.String != planSubscription.Description.String {
		changes["description"] = ":description"
	}
	if reference.UserID.String != planSubscription.UserID.String {
		changes["user_id"] = ":user_id"
	}
	if reference.PlanID.String != planSubscription.PlanID.String {
		changes["plan_id"] = ":plan_id"
	}
	if reference.StartedAt.Time != planSubscription.StartedAt.Time {
		if true {
			changes["started_at"] = ":started_at"
		}
	}
	if reference.EndsAt.Time != planSubscription.EndsAt.Time {
		if true {
			changes["ends_at"] = ":ends_at"
		}
	}
	if reference.IsActive.Bool != planSubscription.IsActive.Bool {
		changes["is_active"] = ":is_active"
	}
	if reference.IsLogicalDeleted.Bool != planSubscription.IsLogicalDeleted.Bool {
		changes["is_logical_deleted"] = ":is_logical_deleted"
	}
	if reference.UpdatedAt.Time != planSubscription.UpdatedAt.Time {
		if true {
			changes["updated_at"] = ":updated_at"
		}
	}
	return changes
}

// PlanChanges - Creates a map ([string]interface{}) including al changing field.
func PlanChanges(plan *models.Plan, reference models.Plan) map[string]string {
	changes := make(map[string]string)
	if reference.Name.String != plan.Name.String {
		changes["name"] = ":name"
	}
	if reference.Description.String != plan.Description.String {
		changes["description"] = ":description"
	}
	if reference.StartedAt.Time != plan.StartedAt.Time {
		if true {
			changes["started_at"] = ":started_at"
		}
	}
	if reference.EndsAt.Time != plan.EndsAt.Time {
		if true {
			changes["ends_at"] = ":ends_at"
		}
	}
	if reference.IsActive.Bool != plan.IsActive.Bool {
		changes["is_active"] = ":is_active"
	}
	if reference.IsLogicalDeleted.Bool != plan.IsLogicalDeleted.Bool {
		changes["is_logical_deleted"] = ":is_logical_deleted"
	}
	if reference.UpdatedAt.Time != plan.UpdatedAt.Time {
		if true {
			changes["updated_at"] = ":updated_at"
		}
	}
	return changes
}
