package main

import (
    "log"
    "net/http"

    "github-info-system/gateway/internal/adapter/rest"
)

func main(){
	collectorAddr:= "localhost:50051"
	handler, err := rest.NewHandler(collectorAddr)
	if err != nil {
        log.Fatalf("Не удалось создать обработчик: %v", err)
    }

	http.HandleFunc("/api/v1/repos/", handler.GetRepoInfo)
    http.HandleFunc("/health", handler.HealthCheck)
	

    http.Handle("/swagger/", http.StripPrefix("/swagger/", http.FileServer(http.Dir("./docs"))))
    
    log.Println("Gateway сервер запущен на порту 8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
	
	
}
