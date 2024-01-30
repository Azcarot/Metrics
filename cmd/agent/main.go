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
	var metric storage.MemStorage
	metric.Gaugemem = make(map[string]storage.Gauge)
	metric.Countermem = make(map[string]storage.Counter)
	var workerData handlers.WorkerData

	workerData.Batchrout = agentflagData.Addr + "/updates/"
	workerData.Singlerout = agentflagData.Addr + "/update/"
	sleeptime := time.Duration(agentflagData.Pollint) * time.Second
	reporttime := time.Duration(agentflagData.Reportint) * time.Second
	reporttimer := time.After(reporttime)
	for {
		select {
		case <-reporttimer:
			body, bodyJSON := agentconfigs.MakeJSON(metric)
			workerData.Body = body
			workerData.BodyJSON = bodyJSON
			for w := 0; w <= agentflagData.RateLimit; w++ {
				go handlers.AgentWorkers(workerData)
			}
			reporttimer = time.After(reporttime)
		default:
			metrics := make(chan storage.MemStorage)
			additionalMetrics := make(chan storage.MemStorage)
			go measure.CollectMetrics(metrics)
			go measure.CollectPSUtilMetrics(additionalMetrics)
			for i := range metrics {
				for id, value := range i.Gaugemem {
					metric.Gaugemem[id] = value
				}
				for id, value := range i.Countermem {
					metric.Countermem[id] = value
				}
			}
			for i := range additionalMetrics {
				for id, value := range i.Gaugemem {
					metric.Gaugemem[id] = value
				}
			}
			time.Sleep(sleeptime)
		}
	}
}
