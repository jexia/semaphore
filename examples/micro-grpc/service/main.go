package main

import (
	"log"
	"time"

	"github.com/jexia/maestro/examples/micro-grpc/proto"
	"github.com/micro/go-micro/v2/service"
	"github.com/micro/go-micro/v2/service/grpc"

	"context"
)

type Say struct{}

func (s *Say) Hello(ctx context.Context, req *proto.Request, rsp *proto.Response) error {
	log.Print("Received Say.Hello request")
	rsp.Msg = "Hello " + req.Name
	return nil
}

func main() {
	service := grpc.NewService(
		service.Name("go.micro.srv.greeter"),
		service.RegisterTTL(time.Second*30),
		service.RegisterInterval(time.Second*10),
	)

	// optionally setup command line usage
	service.Init()

	// Register Handlers
	proto.RegisterSayHandler(service.Server(), new(Say))

	// Run server
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
