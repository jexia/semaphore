package main

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/jexia/semaphore/examples/grpc/proto"
	"google.golang.org/grpc"
)

// Say represents a simple gRPC service
type Say struct{}

// Hello returns a message
func (s *Say) Hello(ctx context.Context, req *proto.Request) (*proto.Response, error) {
	log.Print("Received Say.Hello request")

	res := new(proto.Response)
	res.Msg = "Hello " + req.Name
	res.Meta = &proto.Meta{
		Session: time.Now().Unix(),
	}

	return res, nil
}

func main() {
	lis, err := net.Listen("tcp", ":5050")
	if err != nil {
		log.Fatal(err)
	}

	server := grpc.NewServer()
	proto.RegisterSayServer(server, &Say{})

	log.Println("server up and running on port :5050")
	server.Serve(lis)
}
