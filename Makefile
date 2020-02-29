PROTODIR = annotations

proto-build: $(PROTODIR)/annotations.pb.go

%.pb.go: %.proto
	protoc --proto_path=. --proto_path=$(PROTODIR) --go_out=paths=source_relative:. $^