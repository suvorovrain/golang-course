package main

import (
	"log"
	"os"

	"github.com/artem-smola/golang-course/task2/gateway/docs"
	"github.com/artem-smola/golang-course/task2/gateway/internal/adapter"
	"github.com/artem-smola/golang-course/task2/gateway/internal/controller"
	"github.com/artem-smola/golang-course/task2/gateway/internal/usecase"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

const (
	defaultGatewayServerAddr = ":8081"
	defaultCollectorGRPCAddr = "localhost:50051"
	defaultSwaggerHost       = "localhost:8081"
	defaultSwaggerScheme     = "http"
)

// @title Gateway API
// @version 1.0
// @description REST gateway for retrieving GitHub repository information.
// @host localhost:8081
// @BasePath /
func main() {
	gatewayServerAddr := getEnv("GATEWAY_ADDR", defaultGatewayServerAddr)
	collectorGRPCAddr := getEnv("COLLECTOR_ADDR", defaultCollectorGRPCAddr)
	docs.SwaggerInfo.Host = getEnv("SWAGGER_HOST", defaultSwaggerHost)
	docs.SwaggerInfo.Schemes = []string{getEnv("SWAGGER_SCHEME", defaultSwaggerScheme)}

	gatewayAdapter, err := adapter.NewGatewayGRPCClientAdapter(collectorGRPCAddr)
	if err != nil {
		log.Fatalf("failed to create gateway grpc adapter: %v", err)
	}

	gatewayUsecase := usecase.NewGatewayUsecase(gatewayAdapter)
	gatewayController := controller.NewGatewayRESTServerController(gatewayUsecase)

	router := gin.Default()
	gatewayController.RegisterRoutes(router)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	log.Printf("gateway REST server is listening on %s", gatewayServerAddr)
	if err := router.Run(gatewayServerAddr); err != nil {
		log.Fatalf("failed to run gateway server: %v", err)
	}
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
