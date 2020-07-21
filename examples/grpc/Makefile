install-dev:
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc
	go install google.golang.org/protobuf/cmd/protoc-gen-go
	@echo "Do not forget to install protoc C++ libraries manually"

proto-build: proto/hello.pb.go

%.pb.go: %.proto
	protoc --proto_path=proto --proto_path=../../ --go-grpc_out=proto --go_out=proto $^