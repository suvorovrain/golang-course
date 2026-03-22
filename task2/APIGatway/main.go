package main

import (
	"log"
	"net/http"
	"os"

	"APIGatway/internal/client/grpc"
	http_handler "APIGatway/internal/handler/http"
)

func main() {
	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "8080"
	}

	collectorAddr := os.Getenv("COLLECTOR_ADDR")
	if collectorAddr == "" {
		collectorAddr = "localhost:50051"
	}

	collectorClient, err := grpc.NewCollectorClient(collectorAddr)
	if err != nil {
		log.Fatalf("Failed to create collector client: %v", err)
	}
	defer collectorClient.Close()

	handler := http_handler.NewHandler(collectorClient)

	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	log.Printf("API Gateway started on port %s", httpPort)
	log.Printf("Try: http://localhost:%s/api/repositories?url=https://github.com/torvalds/linux", httpPort)

	if err := http.ListenAndServe(":"+httpPort, mux); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
