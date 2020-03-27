package main

import (
	"context"
	"log"
	"net"

	"github.com/jexia/maestro/examples/grpc/proto"
	"google.golang.org/grpc"
)

type Say struct{}

func (s *Say) Hello(ctx context.Context, req *proto.Request) (*proto.Response, error) {
	log.Print("Received Say.Hello request")
	res := new(proto.Response)
	res.Msg = "Hello " + req.Name

	return res, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}

	server := grpc.NewServer()
	proto.RegisterSayServer(server, &Say{})

	log.Println("server up and running on port :50051")
	server.Serve(lis)
}
