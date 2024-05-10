package protofuncs

import (
	// ...
	"context"
	"fmt"
	"sync"

	_ "google.golang.org/grpc"

	// импортируем пакет со сгенерированными protobuf-файлами
	pb "github.com/Azcarot/Metrics/cmd/proto"
	"github.com/Azcarot/Metrics/internal/storage"
)

// MetricsServer поддерживает все необходимые методы сервера.
type MetricsServer struct {
	// нужно встраивать тип pb.Unimplemented<TypeName>
	// для совместимости с будущими версиями
	pb.UnimplementedMetricsServer

	// используем sync.Map для хранения метрик
	metrics sync.Map
}

// UpdateMetrics реализует интерфейс добавления метрики.
func (s *MetricsServer) UpdateMetrics(ctx context.Context, in *pb.UpdateMetricsRequest) (*pb.UpdateMetricsResponse, error) {
	var response pb.UpdateMetricsResponse
	for _, metric := range in.Metrics {
		if metric.Mtype == storage.GuageType {
			if value, ok := s.metrics.LoadOrStore(metric.Id, metric.Delta); ok {
				newDelta, ok := value.(int64)
				if ok {
					metric.Delta = newDelta + metric.Delta
					s.metrics.Store(metric.Id, metric)
					if _, ok := s.metrics.Load(metric); !ok {

						response.Error = fmt.Sprintf("Не удалось сохранить метрику %s ", metric.Id)
					}
				}
			} else {
				if _, ok := s.metrics.Load(metric.Id); !ok {
					response.Error = fmt.Sprintf("Не удалось сохранить метрику %s ", metric.Id)
				}
			}
		} else {
			s.metrics.Store(metric.Id, metric)
			if _, ok := s.metrics.Load(metric.Id); !ok {
				response.Error = fmt.Sprintf("Не удалось сохранить метрику %s отсутствует", metric.Id)
			}

		}

	}
	return &response, nil
}
