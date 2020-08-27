PROTODIR = api

default: dev

dev:
	sh ./scripts/build.sh

bootstrap:
	go install -mod mod google.golang.org/grpc/cmd/protoc-gen-go-grpc
	go install -mod mod google.golang.org/protobuf/cmd/protoc-gen-go
	@echo "Do not forget to install protoc C++ libraries manually"

proto-build: $(PROTODIR)/annotations.pb.go

test:
	go test -mod vendor -cover -race ./...

bench:
	go test -mod vendor -benchmem -run=^$ -bench=. ./...

%.pb.go: %.proto
	protoc --proto_path=. --proto_path=$(PROTODIR) --go_out=paths=source_relative:. $^
