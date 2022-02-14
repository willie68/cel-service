## Install gRPC and go-grpc generator

`go install google.golang.org/protobuf/cmd/protoc-gen-go@latest`
`go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest`

Don't forget to set PATH to the exe files. Normally <user home>\go\bin

## Generate with

models
`protoc --go_out=. .\api\cel-service.proto`

gRPC Server
`protoc --go-grpc_out=. .\api\cel-service.proto`