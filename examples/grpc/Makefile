PROTODIR = proto

install-dev:
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc
	go install google.golang.org/protobuf/cmd/protoc-gen-go
	@echo "Do not forget to install protoc C++ libraries manually"

proto-build: $(PROTODIR)/greeter.pb.go

%.pb.go: %.proto
	protoc --proto_path=. --proto_path=$(PROTODIR) --go-grpc_out=paths=source_relative:. --go_out=paths=source_relative:. $^