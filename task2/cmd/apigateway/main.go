package main

import apigatewayapp "github.com/Chaice1/golang-course/task2/internal/apigateway/app"

func main() {
	config := apigatewayapp.Config{
		GRPCport: "collector:8082",
		HTTPport: ":8080",
	}

	apigatewayapp.RunApp(config)
}
