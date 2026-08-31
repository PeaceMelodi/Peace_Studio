package handler

import (
	"net/http"

	"github.com/PeaceMelodi/PeaceStudio/internal/worker"
)

func ResultHandler(w http.ResponseWriter, r *http.Request) {
	jobID := r.PathValue("id")

	status, exists := worker.GetStatus(jobID)
	if !exists {
		http.Error(w, "job not found", http.StatusNotFound)
		return
	}

	if status.State != "Finished" || status.OutputPath == "" {
		http.Error(w, "job is not finished yet", http.StatusNotFound)
		return
	}

	http.ServeFile(w, r, status.OutputPath)
}