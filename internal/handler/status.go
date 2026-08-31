package handler

import (
	"fmt"
	"net/http"

	"github.com/PeaceMelodi/PeaceStudio/internal/worker"
)


func StatusHandler(w http.ResponseWriter, r *http.Request) {
	jobID := r.PathValue("id")

	status, exists := worker.GetStatus(jobID)
	if !exists {
		http.Error(w, "job not found", http.StatusNotFound)
		return
	}

	if status.Note != "" {
		fmt.Fprintf(w, "status: %s, note: %s", status.State, status.Note)
		return
	}

	fmt.Fprintf(w, "status: %s", status.State)
}