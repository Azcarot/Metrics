package main

import (
	"agent/measure"
	"agent/postmetrics"
	"fmt"
	"net/http"
)

func main() {
	mets := []measure.AllowedMetrics{
		{Name: "Frees"},
		{Name: "GCCPUFraction"},
		{Name: "GCSys"},
		{Name: "HeapAlloc"},
		{Name: "HeapIdle"},
		{Name: "HeapInuse"},
		{Name: "HeapObjects"},
		{Name: "HeapReleased"},
		{Name: "HeapSys"},
		{Name: "LastGC"},
		{Name: "Lookups"},
		{Name: "MCacheInuse"},
		{Name: "MCacheSys"},
		{Name: "MSpanInuse"},
		{Name: "MSpanSys"},
		{Name: "Mallocs"},
		{Name: "NextGC"},
		{Name: "NumForcedGC"},
		{Name: "NumGC"},
		{Name: "OtherSys"},
		{Name: "PauseTotalNs"},
		{Name: "StackInuse"},
		{Name: "StackSys"},

		{Name: "Sys"},
		{Name: "TotalAlloc"},
	}
	var metric measure.MemStorage
	//Так как у нас метрики пока нигде не хранятся, аллоцируем память
	for range mets {
		metric.Gaugemem = make(map[string]measure.Gauge)
	}
	metric = measure.GetMetrics(metric, mets)
	urls := postmetrics.Makepath(metric)
	var resp *http.Response
	for _, url := range urls {
		resp = postmetrics.PostMetrics(url)
	}
	fmt.Sprintf("metric no longer supported", resp)

}
