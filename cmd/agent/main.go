package main

import (
	"time"

	"github.com/Azcarot/Metrics/internal/agentconfigs"
	"github.com/Azcarot/Metrics/internal/handlers"
	"github.com/Azcarot/Metrics/internal/measure"
	"github.com/Azcarot/Metrics/internal/storage"
)

func main() {
	agentflagData := *agentconfigs.SetValues()
	batchrout := agentflagData.Addr + "/updates/"
	singlerout := agentflagData.Addr + "/update/"
	var metric storage.MemStorage
	sleeptime := time.Duration(agentflagData.Pollint) * time.Second
	reporttime := time.Duration(agentflagData.Reportint) * time.Second
	reporttimer := time.After(reporttime)
	for {
		select {
		case <-reporttimer:
			body, bodyJSON := agentconfigs.MakeJSON(metric)
			handlers.PostJSONMetrics(bodyJSON, batchrout)
			for _, buf := range body {
				handlers.PostJSONMetrics(buf, singlerout)
			}
			reporttimer = time.After(reporttime)
		default:

			metric = measure.CollectMetrics(metric)
			time.Sleep(sleeptime)
		}
	}
}
