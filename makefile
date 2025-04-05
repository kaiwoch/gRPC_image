LOCAL_BIN := $(CURDIR)/bin

install-deps-windows:
	go env -w GOBIN=$(LOCAL_BIN)
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
	go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

install-deps-macOS:
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
	GOBIN=$(LOCAL_BIN) go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

get-deps:
	go get -u google.golang.org/protobuf/cmd/protoc-gen-go
	go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc
	go get -u google.golang.org/grpc/status
	go get -u google.golang.org/protobuf/reflect/protoreflect
	go get -u google.golang.org/protobuf/runtime/protoimpl
	go get -u google.golang.org/protobuf/types/known/emptypb
	go get -u google.golang.org/protobuf/types/known/timestamppb
	go get -u google.golang.org/protobuf/types/known/wrapperspb
	go get -u google.golang.org/grpc

generate-server_api_v1-windows:
	protoc --proto_path api/server_api_v1 \
	--go_out=pkg/server_api_v1 --go_opt=paths=source_relative \
	--plugin=protoc-gen-go=bin/protoc-gen-go.exe \
	--go-grpc_out=pkg/server_api_v1 --go-grpc_opt=paths=source_relative \
	--plugin=protoc-gen-go-grpc=bin/protoc-gen-go-grpc.exe \
	api/server_api_v1/server_api_v1.proto