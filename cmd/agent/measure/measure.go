package measure

import (
	"math/rand"
	"runtime"
)

type Gauge float64
type Counter int64
type MemStorage struct {
	Gaugemem   map[string]Gauge
	Countermem map[string]Counter
}

type AllowedMetrics struct {
	Name string
}

func GetMetrics(m MemStorage) MemStorage {
	var rtm runtime.MemStats
	runtime.ReadMemStats(&rtm)
	m.Gaugemem = make(map[string]Gauge)
	m.Countermem = make(map[string]Counter)
	m.Countermem["PollCount"]++
	m.Gaugemem["BuckHashSys"] = Gauge(rtm.BuckHashSys)
	m.Gaugemem["Frees"] = Gauge(rtm.Frees)
	m.Gaugemem["GCSys"] = Gauge(rtm.GCSys)
	m.Gaugemem["HeapAlloc"] = Gauge(rtm.HeapAlloc)
	m.Gaugemem["HeapIdle"] = Gauge(rtm.HeapIdle)
	m.Gaugemem["HeapInuse"] = Gauge(rtm.HeapInuse)
	m.Gaugemem["HeapObjects"] = Gauge(rtm.HeapObjects)
	m.Gaugemem["HeapReleased"] = Gauge(rtm.HeapReleased)
	m.Gaugemem["HeapSys"] = Gauge(rtm.HeapSys)
	m.Gaugemem["LastGC"] = Gauge(rtm.LastGC)
	m.Gaugemem["Lookups"] = Gauge(rtm.Lookups)
	m.Gaugemem["MCacheInuse"] = Gauge(rtm.MCacheInuse)
	m.Gaugemem["MCacheSys"] = Gauge(rtm.MCacheSys)
	m.Gaugemem["Mallocs"] = Gauge(rtm.Mallocs)
	m.Gaugemem["NextGC"] = Gauge(rtm.NextGC)
	m.Gaugemem["NumForcedGC"] = Gauge(rtm.NumForcedGC)
	m.Gaugemem["NumGC"] = Gauge(rtm.NumGC)
	m.Gaugemem["OtherSys"] = Gauge(rtm.OtherSys)
	m.Gaugemem["PauseTotalNs"] = Gauge(rtm.PauseTotalNs)
	m.Gaugemem["StackInuse"] = Gauge(rtm.StackInuse)
	m.Gaugemem["StackSys"] = Gauge(rtm.StackSys)
	m.Gaugemem["Sys"] = Gauge(rtm.Sys)
	m.Gaugemem["TotalAlloc"] = Gauge(rtm.TotalAlloc)
	m.Gaugemem["RandomValue"] = Gauge(rand.Float64())
	return m
}
