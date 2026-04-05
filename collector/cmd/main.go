package main

import(
	"log"
	"net"
	"github-info-system/api/generated"
    grpcAdapter "github-info-system/collector/internal/adapter/grpc"
    "google.golang.org/grpc"

)
func main(){
	lis, err :=net.Listen("tcp", ":50051")
	
	if err!= nil{
		log.Fatalf("Не удалось запустить слушатель: %v", err)
	}

	grpcServer:= grpc.NewServer() 
	
	generated.RegisterCollectorServiceServer(grpcServer, grpcAdapter.NewServer())


	log.Println("Collector сервер запущен на порту 50051")
    if err := grpcServer.Serve(lis); err != nil { 
		
        log.Fatalf("Ошибка при работе сервера: %v", err)
    }
}