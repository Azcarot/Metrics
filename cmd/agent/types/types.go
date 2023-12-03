package types

type MemStorage struct {
	gaugemem   map[string]gauge
	countermem map[string]counter
}

type AllowedMetrics struct {
	name  string
	valid bool
}

type gauge float64
type counter int64
