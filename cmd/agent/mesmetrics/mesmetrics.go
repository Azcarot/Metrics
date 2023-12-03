package mesmetrics

import (
	"fmt"
	"runtime/metrics"
)

type gauge float64
type counter int64
type MemStorage struct {
	gaugemem   map[string]gauge
	countermem map[string]counter
}

type AllowedMetrics struct {
	name  string
	valid bool
}

func GetMetrics(m MemStorage) MemStorage {
	met := []AllowedMetrics{{
		name:  "BuckHashSys",
		valid: true},
		{
			name:  "Frees",
			valid: true},
		{name: "GCCPUFraction",
			valid: true},
		{name: "GCSys",
			valid: true},
		{name: "HeapAlloc",
			valid: true},
		{name: "HeapIdle",
			valid: true},
		{name: "HeapInuse",
			valid: true},
		{name: "HeapObjects",
			valid: true},
		{name: "HeapReleased",
			valid: true},
		{name: "HeapSys",
			valid: true},
		{name: "LastGC",
			valid: true},
		{name: "Lookups",
			valid: true},
		{name: "MCacheInuse",
			valid: true},
		{name: "MCacheSys",
			valid: true},
		{name: "MSpanInuse",
			valid: true},
		{name: "MSpanSys",
			valid: true},
		{name: "Mallocs",
			valid: true},
		{name: "NextGC",
			valid: true},
		{name: "NumForcedGC",
			valid: true},
		{name: "NumGC",
			valid: true},
		{name: "OtherSys",
			valid: true},
		{name: "PauseTotalNs",
			valid: true},
		{name: "StackInuse",
			valid: true},
		{name: "StackSys",
			valid: true},

		{name: "Sys",
			valid: true},
		{name: "TotalAlloc",
			valid: true},
	}
	for _, metric := range met {
		if metric.valid {
			//Аллоцируем память
			sample := make([]metrics.Sample, 1)
			sample[0].Name = metric.name
			// Читаем метрику
			metrics.Read(sample)
			// Проверяем, поддерживается ли метрика
			if sample[0].Value.Kind() == metrics.KindBad {
				panic(fmt.Sprintf("metric %q no longer supported", metric.name))
			}
			newmem := gauge(sample[0].Value.Uint64())
			m.gaugemem[metric.name] = newmem
			m.countermem[metric.name]++

		}
	}
	return m
}
