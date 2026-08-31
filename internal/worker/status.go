package worker

import "sync"

type Status struct {
	State string 
	Note  string 
}

var statusStore = make(map[string]Status)

var statusMutex sync.Mutex

func SetStatus(jobID string, state string, note string) {
	statusMutex.Lock()
	defer statusMutex.Unlock()
	statusStore[jobID] = Status{State: state, Note: note}
}

func GetStatus(jobID string) (Status, bool) {
	statusMutex.Lock()
	defer statusMutex.Unlock()
	status, exists := statusStore[jobID]
	return status, exists
}