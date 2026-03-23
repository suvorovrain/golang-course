package main

import (
	"golang-course/task2/services/api-gateway/internal/adapter"
	"golang-course/task2/services/api-gateway/internal/controller/httpController"
	"golang-course/task2/services/api-gateway/internal/usecase"
	"log"
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	_ "golang-course/task2/services/api-gateway/docs"

	"github.com/swaggo/http-swagger"
)

// @title           GitHub Repo Info API
// @version         1.0
// @description     This is a sample server for getting GitHub repo info.
// @host            localhost:8080
// @BasePath        /
func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	adapter := adapter.New(conn)
	usecase := usecase.New(adapter)
	handler := httpController.New(usecase)

	http.HandleFunc("/swagger/", httpSwagger.WrapHandler)
	http.HandleFunc("/repo", handler.GetRepo)
	log.Println("Starting server")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
