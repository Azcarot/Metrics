package measure

import (
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
		m.Gaugemem = make(map[string]Gauge)
		m.Countermem = make(map[string]Counter)
		// Читаем метрику
		metrics.Read(sample)
		// Проверяем, поддерживается ли метрика
		var newmem Gauge
		if sample[0].Value.Kind() != metrics.KindBad && sample[0].Value.Kind() == metrics.KindFloat64 {
			newmem = Gauge(sample[0].Value.Float64())
			m.Countermem[metric.Name]++
			m.Gaugemem[metric.Name] = newmem
		}

	}
	return m
}
