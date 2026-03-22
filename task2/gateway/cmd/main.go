package main

import (
	"log"
	"os"

	"task2/gateway/internal/adapter"
	"task2/gateway/internal/handler"

	_ "task2/gateway/docs" //

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// @title           My GitHub Repository API
// @version         1.0
// @description     This is a server for getting information about github repository
// @host            localhost:8080
// @BasePath        /
func main() {
	collectorAddr := os.Getenv("COLLECTOR_ADDR")
	if collectorAddr == "" {
		collectorAddr = "localhost:50051"
	}

	conn, err := grpc.Dial(collectorAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect to collector: %v", err)
	}
	defer conn.Close()

	collectorAdapter := adapter.NewGRPCCollectorAdapter(conn)
	httpHandler := handler.NewHTTPHandler(collectorAdapter)

	r := gin.Default()
	r.GET("/api/repo/:owner/:repo", httpHandler.GetRepo)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	log.Println("Gateway REST server started on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("failed to run gateway: %v", err)
	}
}
