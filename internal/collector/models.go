package collector

type MetricJSON struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

// collector представляет структуру коллектора метрик
type collector struct {
	storage *memStorage
}

// Структура memStorage представляет собой хранилище данных в памяти для коллектора метрик.
// counters - это мапа, которая хранит значения счетчиков метрик.
// gauges - это мапа, которая хранит значения показателей метрик.
type memStorage struct {
	counters map[string]int
	gauges   map[string]string
}
