proto-buffers:
	protoc -I protos/src --go_out=./protos protos/src/**/*.proto

proto-go-plugin:
	protoc -I protos/src --go-grpc_out=./protos protos/src/**/*.proto

protoc: proto-buffers proto-go-plugin

run-authentication:
	cd services/authentication/;go run main.go
run-registration:
	cd services/registration/;go run main.go
run-user:
	cd services/user/;go run main.go
authentication-migration:
	@echo "Enter migration name";
	@read MIGRATION; migrate create -ext sql -dir services/authentication/db/migrations -seq $$MIGRATION

.PHONY: proto-buffers proto-go-plugin protoc