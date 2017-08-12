#!/bin/sh

clear
rm -rf logs/log.txt
go test tests/user_test.go
go test tests/profile_test.go
go test tests/avatar_test.go
go test tests/organization_test.go
go test tests/resource_test.go
go test tests/permission_test.go
go test tests/role_test.go
go test tests/resource_permission_test.go
go test tests/role_permission_test.go
go test tests/user_role_test.go
go test tests/properties_set_test.go
go test tests/plan_test.go
go test tests/plan_subscription_test.go
