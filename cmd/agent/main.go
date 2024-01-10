package main

import (
	"time"

	"github.com/Azcarot/Metrics/internal/agentconfigs"
	"github.com/Azcarot/Metrics/internal/handlers"
	"github.com/Azcarot/Metrics/internal/measure"
	"github.com/Azcarot/Metrics/internal/storage"
)

func main() {
	flagData := *agentconfigs.SetValues()
	var metric storage.MemStorage
	sleeptime := time.Duration(flagData.Pollint) * time.Second
	reporttime := time.Duration(flagData.Reportint) * time.Second
	reporttimer := time.After(reporttime)
	for {
		select {
		case <-reporttimer:
			body := handlers.MakeJSON(metric)
			for _, buf := range body {
				handlers.PostJSONMetrics(buf, flagData.Addr)
			}
			reporttimer = time.After(reporttime)
		default:
			metric = measure.CollectMetrics(metric)
			time.Sleep(sleeptime)
		}
	}
}
