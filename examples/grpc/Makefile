proto-build: proto/hello.pb.go

%.pb.go: %.proto
	protoc --proto_path=proto --proto_path=../../annotations --micro_out=proto --go_out=proto $^