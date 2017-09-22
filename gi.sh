#!/bin/bash
rm -rf views/views.go
go generate
go install
