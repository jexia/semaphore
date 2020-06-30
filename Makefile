PROTODIR = api

install-dev:
	go get -u google.golang.org/grpc
	go get -u github.com/golang/protobuf/{proto,protoc-gen-go}
	@echo "Do not forget to install protoc C++ libraries manually"

proto-build: $(PROTODIR)/annotations.pb.go

%.pb.go: %.proto
	protoc --proto_path=. --proto_path=$(PROTODIR) --go_out=paths=source_relative:. $^