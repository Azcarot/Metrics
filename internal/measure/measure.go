// Сбор системных метрик
package measure

import (
	"math/rand"
	"runtime"
	"strconv"

	"github.com/Azcarot/Metrics/internal/storage"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

// CollectMetrics собирает системные метрики методом runtime.ReadMemStats
func CollectMetrics(results chan<- storage.MemStorage) {
	var rtm runtime.MemStats
	var m storage.MemStorage

	runtime.ReadMemStats(&rtm)
	m.Gaugemem = make(map[string]storage.Gauge)
	m.Countermem = make(map[string]storage.Counter)
	m.Countermem["PollCount"]++
	m.Gaugemem["Alloc"] = storage.Gauge(rtm.Alloc)
	m.Gaugemem["GCCPUFraction"] = storage.Gauge(rtm.GCCPUFraction)
	m.Gaugemem["MSpanInuse"] = storage.Gauge(rtm.MSpanInuse)
	m.Gaugemem["MSpanSys"] = storage.Gauge(rtm.MSpanSys)
	m.Gaugemem["BuckHashSys"] = storage.Gauge(rtm.BuckHashSys)
	m.Gaugemem["Frees"] = storage.Gauge(rtm.Frees)
	m.Gaugemem["GCSys"] = storage.Gauge(rtm.GCSys)
	m.Gaugemem["HeapAlloc"] = storage.Gauge(rtm.HeapAlloc)
	m.Gaugemem["HeapIdle"] = storage.Gauge(rtm.HeapIdle)
	m.Gaugemem["HeapInuse"] = storage.Gauge(rtm.HeapInuse)
	m.Gaugemem["HeapObjects"] = storage.Gauge(rtm.HeapObjects)
	m.Gaugemem["HeapReleased"] = storage.Gauge(rtm.HeapReleased)
	m.Gaugemem["HeapSys"] = storage.Gauge(rtm.HeapSys)
	m.Gaugemem["LastGC"] = storage.Gauge(rtm.LastGC)
	m.Gaugemem["Lookups"] = storage.Gauge(rtm.Lookups)
	m.Gaugemem["MCacheInuse"] = storage.Gauge(rtm.MCacheInuse)
	m.Gaugemem["MCacheSys"] = storage.Gauge(rtm.MCacheSys)
	m.Gaugemem["Mallocs"] = storage.Gauge(rtm.Mallocs)
	m.Gaugemem["NextGC"] = storage.Gauge(rtm.NextGC)
	m.Gaugemem["NumForcedGC"] = storage.Gauge(rtm.NumForcedGC)
	m.Gaugemem["NumGC"] = storage.Gauge(rtm.NumGC)
	m.Gaugemem["OtherSys"] = storage.Gauge(rtm.OtherSys)
	m.Gaugemem["PauseTotalNs"] = storage.Gauge(rtm.PauseTotalNs)
	m.Gaugemem["StackInuse"] = storage.Gauge(rtm.StackInuse)
	m.Gaugemem["StackSys"] = storage.Gauge(rtm.StackSys)
	m.Gaugemem["Sys"] = storage.Gauge(rtm.Sys)
	m.Gaugemem["TotalAlloc"] = storage.Gauge(rtm.TotalAlloc)
	m.Gaugemem["RandomValue"] = storage.Gauge(rand.Float64())
	results <- m
	defer close(results)
}

// CollectPSUtilMetrics собирает метрики пакета gopsUtil
func CollectPSUtilMetrics(results chan<- storage.MemStorage) {
	var m storage.MemStorage
	v, _ := mem.VirtualMemory()
	m.Gaugemem = make(map[string]storage.Gauge)
	m.Gaugemem["TotalMemory"] = storage.Gauge(v.Total)
	m.Gaugemem["FreeMemory"] = storage.Gauge(v.Free)
	result, err := cpu.Percent(0, true)
	if err == nil {
		for i := range result {
			metricName := "CPUutilization" + strconv.Itoa(i)
			m.Gaugemem[metricName] = storage.Gauge(result[i])
		}
	}
	results <- m
	defer close(results)
}
