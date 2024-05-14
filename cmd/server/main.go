// Основной серверный пакет. Ицидиирует связь с бд, создает роутер и слушает назначенный порт
package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"

	_ "net/http/pprof"

	pb "github.com/Azcarot/Metrics/cmd/proto"
	"github.com/Azcarot/Metrics/internal/handlers"
	"github.com/Azcarot/Metrics/internal/serverconfigs"
	"github.com/Azcarot/Metrics/internal/storage"
	"github.com/mitchellh/mapstructure"
	"github.com/soheilhy/cmux"
	"google.golang.org/grpc"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

// MetricsServer поддерживает все необходимые методы сервера.
type MetricsServer struct {
	// нужно встраивать тип pb.Unimplemented<TypeName>
	// для совместимости с будущими версиями
	pb.UnimplementedMetricsServer

	// используем sync.Map для хранения метрик
	metrics sync.Map
}

var Flags storage.Flags

func main() {
	fmt.Printf("Build version=%s\nBuild date =%s\nBuild commit =%s\n", buildVersion, buildDate, buildCommit)
	Flags := serverconfigs.ParseFlagsAndENV()
	if Flags.FlagCrypto != "" {
		var err error
		serverconfigs.PrivateKey, err = serverconfigs.GetPrivateKey(Flags.FlagCrypto)
		if err != nil {
			panic(err)
		}
	}
	if Flags.FlagDBAddr != "" {
		err := storage.NewConn(Flags)
		if err != nil {
			panic(err)
		}
		storage.ST.CheckDBConnection()
		storage.ST.CreateTablesForMetrics()
		defer storage.DB.Close(context.Background())
	}
	r := handlers.MakeRouter(Flags)
	l, err := net.Listen("tcp", Flags.FlagAddr)
	if err != nil {
		log.Fatal(err)
	}

	// Create a cmux.
	m := cmux.New(l)

	// Match connections in order:
	// First grpc, then HTTP, and otherwise Go RPC/TCP.

	grpcL := m.Match(cmux.HTTP2())
	httpL := m.Match(cmux.HTTP1Fast())

	// Create your protocol servers.
	grpcS := grpc.NewServer()
	pb.RegisterMetricsServer(grpcS, &MetricsServer{})

	httpS := &http.Server{
		Handler: r,
	}

	// Use the muxed listeners for your servers.
	go grpcS.Serve(grpcL)
	go httpS.Serve(httpL)
	go handlers.GetSignal(httpS, Flags)
	//Сервер для pprof
	go func() {
		log.Println(http.ListenAndServe(":6060", nil))
	}()

	m.Serve()

}

// UpdateMetric реализует интерфейс добавления метрики.
func (s *MetricsServer) UpdateMetric(ctx context.Context, in *pb.UpdateMetricRequest) (*pb.UpdateMetricResponse, error) {
	var response pb.UpdateMetricResponse
	if in.Metric.Mtype == pb.Mtype_MTYPE_DELTA {
		if value, ok := s.metrics.LoadOrStore(in.Metric.Id, in.Metric.Delta); ok {
			newDelta, ok := value.(int64)
			if ok {
				in.Metric.Delta = newDelta + in.Metric.Delta
				s.metrics.Store(in.Metric.Id, in.Metric)
				if _, ok := s.metrics.Load(in.Metric.Id); ok {
					response.Metric = in.Metric
				} else {
					return &response, fmt.Errorf("не удалось обновить Delta метрику %s", in.Metric.Id)
				}
			}
		} else {
			if _, ok := s.metrics.Load(in.Metric.Id); ok {
				response.Metric = in.Metric
			} else {
				return &response, fmt.Errorf("не удалось сохранить Delta метрику %s", in.Metric.Id)
			}
		}
	} else {
		s.metrics.Store(in.Metric.Id, in.Metric)
		if _, ok := s.metrics.Load(in.Metric.Id); ok {
			response.Metric = in.Metric
		} else {
			return &response, fmt.Errorf("не удалось сохранить Gauge метрику %s", in.Metric.Id)
		}

	}
	var metricData storage.Metrics
	err := mapstructure.Decode(in.Metric, &metricData)

	if err != nil {
		fmt.Println(in.Metric)
		return &response, fmt.Errorf("не правильный формат метрики %s", in.Metric.Id)
	}
	if len(Flags.FlagDBAddr) != 0 {
		storage.PgxStorage.WriteMetricsToPstgrs(storage.ST, metricData)
	}

	if len(Flags.FlagFileStorage) != 0 && Flags.FlagStoreInterval == 0 {
		fileName := Flags.FlagFileStorage
		storage.WriteToFile(fileName, metricData)
	}
	return &response, nil
}
