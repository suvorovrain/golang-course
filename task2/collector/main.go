package main

import (
	"collector/internal/usecase/client"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"

	pb "collector/api/proto"
	grpc_handler "collector/internal/handler/grpc"
	_ "collector/internal/usecase/client"
)

func main() {
	port := os.Getenv("GRPC_PORT")
	if port == "" {
		port = "50051"
	}

	githubToken := os.Getenv("GITHUB_TOKEN")

	useCase := client.NewUseCase(githubToken)

	handler := grpc_handler.NewHandler(useCase)

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterCollectorServiceServer(grpcServer, handler)

	log.Printf("Collector service started on port %s", port)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
