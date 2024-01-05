package handlers

import (
	"flag"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Azcarot/Metrics/cmd/types"
	"github.com/caarlos0/env/v6"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

var sugar zap.SugaredLogger
var Flag types.Flags

type (
	// берём структуру для хранения сведений об ответе
	responseData struct {
		status int
		size   int
	}

	// добавляем реализацию http.ResponseWriter
	loggingResponseWriter struct {
		http.ResponseWriter // встраиваем оригинальный http.ResponseWriter
		responseData        *responseData
	}
)

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	// записываем ответ, используя оригинальный http.ResponseWriter
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size // захватываем размер
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	// записываем код статуса, используя оригинальный http.ResponseWriter
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode // захватываем код статуса
}

func (st *StorageHandler) HandlePostMetrics() http.Handler {
	postMetric := func(res http.ResponseWriter, req *http.Request) {
		if len(chi.URLParam(req, "name")) == 0 || len(chi.URLParam(req, "value")) == 0 || len(chi.URLParam(req, "type")) == 0 {
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		switch strings.ToLower(chi.URLParam(req, "type")) {
		case types.CounterType, types.GuageType:
			err := st.Storage.StoreMetrics(chi.URLParam(req, "name"), strings.ToLower(chi.URLParam(req, "type")), chi.URLParam(req, "value"))
			if err != nil {
				res.WriteHeader(http.StatusBadRequest)
				return
			}
			res.WriteHeader(http.StatusOK)
		default:
			res.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	return http.HandlerFunc(postMetric)
}

func ParseFlagsAndENV() types.Flags {
	flag.StringVar(&Flag.FlagAddr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&Flag.FlagFileStorage, "f", "/tmp/metrics-db.json", "address of a file-storage")
	flag.IntVar(&Flag.FlagStoreInterval, "i", 300, "interval for storing data")
	flag.BoolVar(&Flag.FlagRestore, "r", true, "reading data from file first")
	flag.Parse()
	var envcfg types.ServerENV
	err := env.Parse(&envcfg)
	if err != nil {
		log.Fatal(err)
	}
	if envcfg.Address != "" {
		Flag.FlagAddr = envcfg.Address
	}
	if envcfg.FileStorage != "" {
		Flag.FlagFileStorage = envcfg.FileStorage
	}
	if (envcfg.Restore) || !envcfg.Restore {
		Flag.FlagRestore = envcfg.Restore
	}
	if len(envcfg.StoreInterval) == 0 {
		storeInterval, err := strconv.Atoi(envcfg.FileStorage)
		if err == nil {
			Flag.FlagStoreInterval = storeInterval
		}
	}
	return Flag
}

func MakeRouter(flag types.Flags) *chi.Mux {
	storagehandler := StorageHandler{
		Storage: &types.MemStorage{
			Gaugemem: make(map[string]types.Gauge), Countermem: make(map[string]types.Counter)},
	}
	if flag.FlagRestore && len(flag.FlagFileStorage) > 0 {
		storagehandler.Storage.ReadMetricsFromFile(flag.FlagFileStorage)
	}
	logger, err := zap.NewDevelopment()
	if err != nil {
		// вызываем панику, если ошибка
		panic(err)
	}
	if flag.FlagStoreInterval > 0 && len(flag.FlagFileStorage) > 0 {
		go func(name string) {
			reporttime := time.Duration(flag.FlagStoreInterval) * time.Second
			reporttimer := time.After(reporttime)
			for {
				<-reporttimer
				fullMetrics := storagehandler.Storage.GetAllMetricsAsMetricType()
				for _, data := range fullMetrics {
					WriteToFile(flag.FlagFileStorage, data)
				}
			}
		}(flag.FlagFileStorage)
	}
	defer logger.Sync()
	// делаем регистратор SugaredLogger
	sugar = *logger.Sugar()
	r := chi.NewRouter()
	r.Use()
	r.Route("/", func(r chi.Router) {
		r.Get("/", WithLogging(GzipHandler(storagehandler.HandleGetAllMetrics(flag))).ServeHTTP)
		r.Post("/update/", WithLogging(GzipHandler(storagehandler.HandleJSONPostMetrics(flag))).ServeHTTP)
		r.Post("/value/", WithLogging(GzipHandler(storagehandler.HandleJSONGetMetrics())).ServeHTTP)
		r.Post("/update/{type}/{name}/{value}", WithLogging(GzipHandler(storagehandler.HandlePostMetrics())).ServeHTTP)
		r.Get("/value/{name}/{type}", WithLogging(GzipHandler(storagehandler.HandleGetMetrics())).ServeHTTP)
	})
	return r
}

type StorageHandler struct {
	Storage types.MemInteractions
}

// WithLogging добавляет дополнительный код для регистрации сведений о запросе
// и возвращает новый http.Handler.
func WithLogging(h http.Handler) http.Handler {
	logFn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		responseData := &responseData{
			status: 0,
			size:   0,
		}
		lw := loggingResponseWriter{
			ResponseWriter: w, // встраиваем оригинальный http.ResponseWriter
			responseData:   responseData,
		}
		h.ServeHTTP(&lw, r) // внедряем реализацию http.ResponseWriter

		duration := time.Since(start)

		sugar.Infoln(
			"uri", r.RequestURI,
			"method", r.Method,
			"body", r.Body,
			"status", responseData.status, // получаем перехваченный код статуса ответа
			"duration", duration,
			"size", responseData.size, // получаем перехваченный размер ответа
		)
	}
	// возвращаем функционально расширенный хендлер
	return http.HandlerFunc(logFn)
}

func (st *StorageHandler) HandleGetMetrics() http.Handler {
	getMetric := func(res http.ResponseWriter, req *http.Request) {

		result, err := st.Storage.GetStoredMetrics(chi.URLParam(req, "type"), chi.URLParam(req, "name"))
		res.Header().Set("Content-Type", "text/html")
		if err != nil {
			res.WriteHeader(http.StatusNotFound)
		} else {
			io.WriteString(res, result)
		}
	}
	return http.HandlerFunc(getMetric)
}

func (st *StorageHandler) HandleGetAllMetrics(flag types.Flags) http.Handler {
	getMetrics := func(res http.ResponseWriter, req *http.Request) {
		result := st.Storage.GetAllMetrics()
		res.Header().Set("Content-Type", "text/html")
		io.WriteString(res, result)
	}
	return http.HandlerFunc(getMetrics)
}
