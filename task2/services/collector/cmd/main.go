package main

import (
	"golang-course/task2/pkg/api"
	"golang-course/task2/services/collector/internal/adapter/github"
	"golang-course/task2/services/collector/internal/controller/grpcController"
	"golang-course/task2/services/collector/internal/usecase"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
)

func main() {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		log.Fatalf("github token not found")
	}

	githubClient := github.NewClient(token)

	uc := usecase.New(githubClient)

	srv := grpc_server.New(uc)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	api.RegisterRepositoryServiceServer(grpcServer, srv)

	grpcServer.Serve(lis)

}
