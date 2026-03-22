package main

import (
	"log"
	"net"

	"github.com/artem-smola/golang-course/task2/collector/internal/adapter"
	"github.com/artem-smola/golang-course/task2/collector/internal/controller"
	"github.com/artem-smola/golang-course/task2/collector/internal/usecase"
	"github.com/artem-smola/golang-course/task2/proto/gen"
	"google.golang.org/grpc"
)

const grpcServerPort = ":50051"

func main() {
	collectorAdapter := &adapter.CollectorAdapter{}
	collectorUsecase := usecase.NewCollectorUsecase(collectorAdapter)
	collectorController := controller.NewCollectorGRPCServerController(collectorUsecase)

	listener, err := net.Listen("tcp", grpcServerPort)
	if err != nil {
		log.Fatalf("failed to listen on %s: %v", grpcServerPort, err)
	}

	grpcServer := grpc.NewServer()
	gen.RegisterGRPCServiceServer(grpcServer, collectorController)

	log.Printf("collector gRPC server is listening on %s", grpcServerPort)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve gRPC server: %v", err)
	}
}
