package worker

import (
	"fmt"
)

type Job struct {
	ID       string
	FilePath string
	Action   string
}

var JobQueue = make(chan Job, 100)

func StartWorkers(numWorkers int) {
	for i := 1; i <= numWorkers; i++ {
		go worker(i)
	}
}

func worker(id int) {
	for job := range JobQueue {
		fmt.Printf("worker %d picked up: %s (action: %s)\n", id, job.FilePath, job.Action)
		SetStatus(job.ID, "Getting worked on", "")

		var outputPath string
		var err error

		switch job.Action {
		case "grayscale":
			outputPath, err = ToGrayscale(job.ID, job.FilePath)
		case "blur":
			outputPath, err = ToBlur(job.ID, job.FilePath)
		case "resize":
			outputPath, err = ToResize(job.ID, job.FilePath)
		case "pixelate":
			outputPath, err = ToPixelate(job.ID, job.FilePath)
		default:
			note := fmt.Sprintf("unknown action: %s", job.Action)
			fmt.Printf("worker %d %s\n", id, note)
			SetStatus(job.ID, "Failed", note)
			continue
		}

		if err != nil {
			note := err.Error()
			fmt.Printf("worker %d failed to process %s: %v\n", id, job.FilePath, err)
			SetStatus(job.ID, "Failed", note)
			continue
		}

		fmt.Printf("worker %d finished: %s\n", id, outputPath)
		SetFinished(job.ID, outputPath)
	}
}