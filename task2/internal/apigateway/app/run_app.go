package apigatewayapp

import (
	"log"

	_ "github.com/Chaice1/golang-course/task2/docs"
	collectorpb "github.com/Chaice1/golang-course/task2/gen"
	apigatewayclient "github.com/Chaice1/golang-course/task2/internal/apigateway/adapters/client"
	apigatewayhandler "github.com/Chaice1/golang-course/task2/internal/apigateway/adapters/handler"
	apigatewayusecase "github.com/Chaice1/golang-course/task2/internal/apigateway/usecase"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func RunApp(config Config) {

	conn, err := grpc.NewClient(config.GRPCport, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatal(err)
	}

	client := collectorpb.NewCollectorClient(conn)

	CollectorClient := apigatewayclient.NewCollectorClient(client)

	ApigatewayUsecase := apigatewayusecase.NewUsecaseApiGateway(CollectorClient)

	apigatewayHandler := apigatewayhandler.NewApiGatewayHandler(ApigatewayUsecase)

	server := gin.Default()

	server.GET("/get_repo_info/:owner/:repo", apigatewayHandler.GetRepoInfo)
	server.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	if err := server.Run(config.HTTPport); err != nil {
		log.Fatal(err)
	}
}
