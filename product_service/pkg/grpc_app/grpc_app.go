package grpcapp

import (
	"product_service/pkg/service"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
)

type GRPCApp struct {
	gRPCServerver *grpc.Server
	port          int
}

func NewGRPCApp(port int, service service.GRPC) *GRPCApp {
	gRPCServer := grpc.NewServer()
	NewProductGRPCServer(gRPCServer, service)

	return &GRPCApp{
		gRPCServerver: gRPCServer,
		port:          port,
	}
}

func (a *GRPCApp) Run() error {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		log.Fatalf("start grpc server error %v", err)
	}
	log.Println("gRPC servcer start")

	if err = a.gRPCServerver.Serve(l); err != nil {
		log.Fatalf("start grpc server error %v", err)
	}

	return nil
}

func (a *GRPCApp) Stop() {
	log.Println("gRPC servcer stop")
	a.gRPCServerver.GracefulStop()

}
