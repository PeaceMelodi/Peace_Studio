package main

import (
	"fmt"
	"net/http"

	"github.com/PeaceMelodi/PeaceStudio/internal/handler"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "github.com/PeaceMelodi/PeaceStudio/docs"
	"github.com/PeaceMelodi/PeaceStudio/internal/worker"
)


func main() {
	worker.StartWorkers(3)

	http.HandleFunc("/upload", handler.UploadHandler)
	http.HandleFunc("/status/{id}", handler.StatusHandler)
	http.Handle("/swagger/", httpSwagger.WrapHandler)

	fmt.Println("PeaceStudio server running on http://localhost:8080")
	fmt.Println("Swagger UI available at http://localhost:8080/swagger/index.html")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("server error:", err)
	}
}