// Агент для сбора метрик из системы и отправки их на сервер.
// Собирает данные с указанными флагами интервалами
// Передает данные как в JSON, так и в виде строк через url
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	pb "github.com/Azcarot/Metrics/cmd/proto"
	"github.com/Azcarot/Metrics/internal/agentconfigs"
	"github.com/Azcarot/Metrics/internal/handlers"
	"github.com/Azcarot/Metrics/internal/measure"
	"github.com/Azcarot/Metrics/internal/storage"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func main() {
	fmt.Printf("Build version=%s\nBuild date =%s\nBuild commit =%s\n", buildVersion, buildDate, buildCommit)
	agentflagData := *agentconfigs.SetValues()
	var metric storage.MemStorage
	ctx, cancel := context.WithCancel(context.Background())
	metric.Gaugemem = make(map[string]storage.Gauge)
	metric.Countermem = make(map[string]storage.Counter)
	var workerData handlers.WorkerData
	conn, err := grpc.Dial(agentflagData.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	// получаем переменную интерфейсного типа UsersClient,
	// через которую будем отправлять сообщения
	c := pb.NewMetricsClient(conn)
	shutdown := make(chan bool)
	workerData.Batchrout = agentflagData.Addr + "/updates/"
	workerData.Singlerout = agentflagData.Addr + "/update/"
	workerData.AgentflagData = agentflagData
	sleeptime := time.Duration(agentflagData.Pollint) * time.Second
	reporttime := time.Duration(agentflagData.Reportint) * time.Second
	reporttimer := time.After(reporttime)
	go handlers.GetAgentSignal(workerData, metric, shutdown, agentflagData)
	for {
		select {
		case <-reporttimer:
			body, bodyJSON := agentconfigs.MakeJSON(metric)
			workerData.Body = body
			workerData.BodyJSON = bodyJSON
			for w := 0; w <= agentflagData.RateLimit; w++ {
				go handlers.AgentWorkers(workerData)
				go SendGrpcMetrics(ctx, c, workerData)
			}

			reporttimer = time.After(reporttime)
		case <-shutdown:
			cancel()
			return
		default:
			metrics := make(chan storage.MemStorage)
			additionalMetrics := make(chan storage.MemStorage)
			go measure.CollectMetrics(metrics)
			go measure.CollectPSUtilMetrics(additionalMetrics)
			for i := range metrics {
				for id, value := range i.Gaugemem {
					metric.Gaugemem[id] = value
				}
				for id, value := range i.Countermem {
					metric.Countermem[id] = value
				}
			}
			for i := range additionalMetrics {
				for id, value := range i.Gaugemem {
					metric.Gaugemem[id] = value
				}
			}
			time.Sleep(sleeptime)
		}
	}

}

func SendGrpcMetrics(ctx context.Context, c pb.MetricsClient, workerdata handlers.WorkerData) {
	metrics := []*pb.Metric{}

	err := json.Unmarshal(workerdata.BodyJSON, &metrics)
	if err != nil {
		log.Fatal(err)
	}
	for _, m := range metrics {
		resp, err := c.UpdateMetric(ctx, &pb.UpdateMetricRequest{Metric: m})
		if err != nil {
			log.Print("Ошибка обновление метрики по http2: ", err)
			return
		}
		log := log.Default()
		log.Print(resp)
	}

}
