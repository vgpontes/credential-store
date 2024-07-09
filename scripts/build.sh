#/bin/bash

set -e -x

cd lambda
go env -w GOOS=linux GOARCH=arm64
go build -tags lambda.norpc -o ../build/auth_service/bootstrap ./auth_service/auth_service.go
go build -tags lambda.norpc -o ../build/users/bootstrap ./users/users.go ./users/users_table.go ./users/users_types.go
cd ..