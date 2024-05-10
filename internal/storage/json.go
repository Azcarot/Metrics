package storage

type Metrics struct {
	ID    string   `json:"id" mapstructure:"Id"`                 // имя метрики
	MType string   `json:"type" mapstructure:"Mtype"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty" mapstructure:"Delta"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty" mapstructure:"Value"` // значение метрики в случае передачи gauge
}

const JSONContentType = "application/json"
