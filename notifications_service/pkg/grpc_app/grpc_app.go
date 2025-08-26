package grpcapp

import (
	"fmt"
	"log"
	"net"
	"notifications_service/pkg/service"

	"google.golang.org/grpc"
)

type GRPCApp struct {
	gRPCServer *grpc.Server
	port       int
}

func NewGRPCApp(port int, service service.GRPC) *GRPCApp {
	gRPCServer := grpc.NewServer()
	NewOrderGRPCService(gRPCServer, service)

	return &GRPCApp{
		gRPCServer: gRPCServer,
		port:       port,
	}
}

func (a *GRPCApp) Run() error {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		log.Fatalf("start grpc server error %v", err)
	}
	log.Println("gRPC servcer start")

	if err = a.gRPCServer.Serve(l); err != nil {
		log.Fatalf("start grpc server error %v", err)
	}

	return nil
}

func (a *GRPCApp) Stop() {
	log.Println("gRPC servcer stop")
	a.gRPCServer.GracefulStop()

}
