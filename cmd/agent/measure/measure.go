package measure

import (
	"fmt"
	"runtime/metrics"
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

func GetMetrics(m MemStorage, a []AllowedMetrics) MemStorage {

	for _, metric := range a {
		//Аллоцируем память
		sample := make([]metrics.Sample, 1)
		sample[0].Name = metric.Name
		// Читаем метрику
		metrics.Read(sample)
		// Проверяем, поддерживается ли метрика
		if sample[0].Value.Kind() == metrics.KindBad {
			panic(fmt.Sprintf("metric %q no longer supported", metric.Name))
		}
		newmem := Gauge(sample[0].Value.Uint64())
		m.Gaugemem[metric.Name] = newmem
		m.Countermem[metric.Name]++

	}
	return m
}
