package main

import (
	"fmt"
	"net/http"

	"github.com/PeaceMelodi/PeaceStudio/internal/handler"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "github.com/PeaceMelodi/PeaceStudio/docs"
	"github.com/PeaceMelodi/PeaceStudio/internal/worker"
)

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	worker.StartWorkers(3) // 3 background workers picking up pictures

	mux := http.NewServeMux()
	mux.HandleFunc("/upload", handler.UploadHandler)
	mux.HandleFunc("/status/{id}", handler.StatusHandler)
	mux.HandleFunc("/result/{id}", handler.ResultHandler)
	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	fmt.Println("PeaceStudio server running on http://localhost:8080")
	fmt.Println("Swagger UI available at http://localhost:8080/swagger/index.html")

	err := http.ListenAndServe(":8080", corsMiddleware(mux))
	if err != nil {
		fmt.Println("server error:", err)
	}
}