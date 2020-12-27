proto-buffers:
	protoc -I protos/src --go_out=./protos protos/src/**/*.proto

proto-go-plugin:
	protoc -I protos/src --go-grpc_out=./protos protos/src/**/*.proto

protoc: proto-buffers proto-go-plugin

run-authentication:
	go run services/authentication/main.go

.PHONY: proto-buffers proto-go-plugin protoc