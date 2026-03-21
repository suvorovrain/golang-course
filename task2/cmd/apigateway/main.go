package main

import apigatewayapp "github.com/Chaice1/golang-course/task2/internal/apigateway/app"

// @title           GitHub Collector API
// @version         1.0
// @description     API Gateway for GitHub Repository Information.
// @host            localhost:8080
// @BasePath        /
func main() {
	config := apigatewayapp.Config{
		GRPCport: "collector:8082",
		HTTPport: ":8080",
	}

	apigatewayapp.RunApp(config)
}
