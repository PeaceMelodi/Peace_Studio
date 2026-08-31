package worker

import "sync"

type Status struct {
	State      string 
	Note       string 
	OutputPath string 
}


var statusStore = make(map[string]Status)


var statusMutex sync.Mutex

func SetStatus(jobID string, state string, note string) {
	statusMutex.Lock()
	defer statusMutex.Unlock()
	existing := statusStore[jobID]
	existing.State = state
	existing.Note = note
	statusStore[jobID] = existing
}

func SetFinished(jobID string, outputPath string) {
	statusMutex.Lock()
	defer statusMutex.Unlock()
	statusStore[jobID] = Status{State: "Finished", OutputPath: outputPath}
}

func GetStatus(jobID string) (Status, bool) {
	statusMutex.Lock()
	defer statusMutex.Unlock()
	status, exists := statusStore[jobID]
	return status, exists
}