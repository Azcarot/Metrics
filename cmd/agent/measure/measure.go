package measure

import (
	"math/rand"
	"runtime"

	"github.com/Azcarot/Metrics/cmd/types"
)

func CollectMetrics(m types.MemStorage) types.MemStorage {
	var rtm runtime.MemStats
	runtime.ReadMemStats(&rtm)
	m.Gaugemem = make(map[string]types.Gauge)
	m.Countermem = make(map[string]types.Counter)
	m.Countermem["PollCount"]++
	m.Gaugemem["Alloc"] = types.Gauge(rtm.Alloc)
	m.Gaugemem["GCCPUFraction"] = types.Gauge(rtm.GCCPUFraction)
	m.Gaugemem["MSpanInuse"] = types.Gauge(rtm.MSpanInuse)
	m.Gaugemem["MSpanSys"] = types.Gauge(rtm.MSpanSys)
	m.Gaugemem["BuckHashSys"] = types.Gauge(rtm.BuckHashSys)
	m.Gaugemem["Frees"] = types.Gauge(rtm.Frees)
	m.Gaugemem["GCSys"] = types.Gauge(rtm.GCSys)
	m.Gaugemem["HeapAlloc"] = types.Gauge(rtm.HeapAlloc)
	m.Gaugemem["HeapIdle"] = types.Gauge(rtm.HeapIdle)
	m.Gaugemem["HeapInuse"] = types.Gauge(rtm.HeapInuse)
	m.Gaugemem["HeapObjects"] = types.Gauge(rtm.HeapObjects)
	m.Gaugemem["HeapReleased"] = types.Gauge(rtm.HeapReleased)
	m.Gaugemem["HeapSys"] = types.Gauge(rtm.HeapSys)
	m.Gaugemem["LastGC"] = types.Gauge(rtm.LastGC)
	m.Gaugemem["Lookups"] = types.Gauge(rtm.Lookups)
	m.Gaugemem["MCacheInuse"] = types.Gauge(rtm.MCacheInuse)
	m.Gaugemem["MCacheSys"] = types.Gauge(rtm.MCacheSys)
	m.Gaugemem["Mallocs"] = types.Gauge(rtm.Mallocs)
	m.Gaugemem["NextGC"] = types.Gauge(rtm.NextGC)
	m.Gaugemem["NumForcedGC"] = types.Gauge(rtm.NumForcedGC)
	m.Gaugemem["NumGC"] = types.Gauge(rtm.NumGC)
	m.Gaugemem["OtherSys"] = types.Gauge(rtm.OtherSys)
	m.Gaugemem["PauseTotalNs"] = types.Gauge(rtm.PauseTotalNs)
	m.Gaugemem["StackInuse"] = types.Gauge(rtm.StackInuse)
	m.Gaugemem["StackSys"] = types.Gauge(rtm.StackSys)
	m.Gaugemem["Sys"] = types.Gauge(rtm.Sys)
	m.Gaugemem["TotalAlloc"] = types.Gauge(rtm.TotalAlloc)
	m.Gaugemem["RandomValue"] = types.Gauge(rand.Float64())
	return m
}
