package handler

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/PeaceMelodi/PeaceStudio/internal/worker"
)

const maxUploadSize = 5 << 20 


var validActions = map[string]bool{
	"grayscale": true,
	"blur":      true,
	"resize":    true,
	"pixelate":  true,
}


func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "file too large", http.StatusBadRequest)
		return
	}

	
	action := r.FormValue("action")
	if !validActions[action] {
		http.Error(w, "invalid or missing action: must be grayscale, blur, resize, or pixelate", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("picture")
	if err != nil {
		http.Error(w, "no picture found in request", http.StatusBadRequest)
		return
	}
	defer file.Close()

	if header.Size > maxUploadSize {
		http.Error(w, "file exceeds 5MB limit", http.StatusBadRequest)
		return
	}

	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		http.Error(w, "unable to read file", http.StatusInternalServerError)
		return
	}
	contentType := http.DetectContentType(buffer)

	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/gif":  true,
		"image/webp": true,
	}

	if !allowedTypes[contentType] {
		http.Error(w, "file is not a valid image", http.StatusBadRequest)
		return
	}

	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		http.Error(w, "unable to process file", http.StatusInternalServerError)
		return
	}

	dstPath := filepath.Join("uploads", "temp", header.Filename)
	dst, err := os.Create(dstPath)
	if err != nil {
		http.Error(w, "unable to save file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, "error saving file", http.StatusInternalServerError)
		return
	}

	jobID := header.Filename + "-" + strconv.FormatInt(time.Now().UnixNano(), 10)

	worker.SetStatus(jobID, "Waiting", "")

	worker.JobQueue <- worker.Job{ID: jobID, FilePath: dstPath, Action: action}

	fmt.Fprintf(w, "picture %s uploaded successfully, action: %s, job id: %s", header.Filename, action, jobID)
}