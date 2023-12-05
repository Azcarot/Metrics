package types

type Gauge float64
type Counter int64
type MemStorage struct {
	Gaugemem   map[string]Gauge
	Countermem map[string]Counter
}

type AllowedMetrics struct {
	Name string
}
type MemInteractions interface {
	GetStoredMetrics() MemStorage
	PostMetrics()
}

var MetricNameTypes = map[string]string{
	"PollCount":    "counter",
	"Frees":        "gauge",
	"GCSys":        "gauge",
	"HeapAlloc":    "gauge",
	"HeapIdle":     "gauge",
	"HeapInuse":    "gauge",
	"HeapObjects":  "gauge",
	"HeapReleased": "gauge",
	"HeapSys":      "gauge",
	"LastGC":       "gauge",
	"Lookups":      "gauge",
	"MCacheInuse":  "gauge",
	"MCacheSys":    "gauge",
	"NextGC":       "gauge",
	"Mallocs":      "gauge",
	"NumForcedGC":  "gauge",
	"NumGC":        "gauge",
	"OtherSys":     "gauge",
	"PauseTotalNs": "gauge",
	"StackInuse":   "gauge",
	"StackSys":     "gauge",
	"Sys":          "gauge",
	"TotalAlloc":   "gauge",
	"RandomValue":  "gauge",
}
