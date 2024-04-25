package handlers

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/Azcarot/Metrics/internal/agentconfigs"
	"github.com/Azcarot/Metrics/internal/storage"
)

func GetAgentSignal(workerData WorkerData, metric storage.MemStorage, agentflagData agentconfigs.AgentData) {
	terminateSignals := make(chan os.Signal, 1)
	signal.Notify(terminateSignals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT) //NOTE:: syscall.SIGKILL we cannot catch kill -9 as its force kill signal.

	_, ok := <-terminateSignals
	if ok {
		body, bodyJSON := agentconfigs.MakeJSON(metric)
		workerData.Body = body
		workerData.BodyJSON = bodyJSON
		for w := 0; w <= agentflagData.RateLimit; w++ {
			go AgentWorkers(workerData)
		}

	}

}
