package types

type Metrics struct {
	ID    string   `json:"id"`                     // имя метрики
	MType string   `json:"type"`                   // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,string,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,string,omitempty"` // значение метрики в случае передачи gauge
}

const JSONContentType = "application/json"
