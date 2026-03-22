package main

import (
	"log"
	"net"

	"task2/collector/internal/adapter"
	grpc_handler "task2/collector/internal/handler"
	"task2/collector/internal/usecase"
	"task2/proto"

	"google.golang.org/grpc"
)

func main() {
	githubAdapter := adapter.NewGitHubAdapter()
	repoUsecase := usecase.NewRepoUsecase(githubAdapter)
	grpcHandler := grpc_handler.NewServer(repoUsecase)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	proto.RegisterRepositoryServiceServer(s, grpcHandler)

	log.Printf("Collector gRPC server started on :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
