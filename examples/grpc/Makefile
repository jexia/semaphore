proto-build: proto/hello.pb.go

%.pb.go: %.proto
	protoc --proto_path=proto --proto_path=../../../ --go_out=plugins=grpc:proto $^