package storage_test

import "github.com/Azcarot/Metrics/internal/storage"

func Example() {
	type StorageHandler struct {
		Storage storage.MemInteractions
	}
	Storagehandler := StorageHandler{
		Storage: &storage.MemStorage{
			Gaugemem: make(map[string]storage.Gauge), Countermem: make(map[string]storage.Counter)},
	}
	value := float64(333)
	data := storage.Metrics{ID: "222", MType: storage.GuageType, Value: &value}
	storage.MemInteractions.StoreMetrics(Storagehandler.Storage, data)
	//Output:
	//
}
