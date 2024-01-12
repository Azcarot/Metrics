package handlers

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/Azcarot/Metrics/internal/storage"
)

func (st *StorageHandler) ShutdownSave(s *http.Server, flag storage.Flags) {
	terminateSignals := make(chan os.Signal, 1)
	signal.Notify(terminateSignals, syscall.SIGINT, syscall.SIGTERM) //NOTE:: syscall.SIGKILL we cannot catch kill -9 as its force kill signal.

	_, ok := <-terminateSignals
	if ok && len(flag.FlagFileStorage) != 0 {
		FinalData := Storagehandler.Storage.GetAllMetricsAsMetricType()

		for _, metric := range FinalData {
			storage.WriteToFile(flag.FlagFileStorage, metric)
		}
		s.Shutdown(context.Background())
	}

}

func GetSignal(s *http.Server, f storage.Flags) {
	Storagehandler = StorageHandler{
		Storage: &storage.MemStorage{
			Gaugemem: make(map[string]storage.Gauge), Countermem: make(map[string]storage.Counter)},
	}
	Storagehandler.ShutdownSave(s, f)
}
