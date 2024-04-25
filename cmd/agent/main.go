// Агент для сбора метрик из системы и отпраки их на сервер.
// Собирает данные с указанными флагами интервалами
// Передает данные как в JSON, так и в виде строк через url
package main

import (
	"fmt"
	"time"

	"github.com/Azcarot/Metrics/internal/agentconfigs"
	"github.com/Azcarot/Metrics/internal/handlers"
	"github.com/Azcarot/Metrics/internal/measure"
	"github.com/Azcarot/Metrics/internal/storage"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func main() {
	fmt.Printf("Build version=%s\nBuild date =%s\nBuild commit =%s\n", buildVersion, buildDate, buildCommit)
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
	go handlers.GetAgentSignal(workerData, metric, agentflagData)
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
