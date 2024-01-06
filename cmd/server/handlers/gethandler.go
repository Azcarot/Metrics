package handlers

import (
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Azcarot/Metrics/cmd/storage"
	"github.com/caarlos0/env/v6"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

var sugar zap.SugaredLogger
var Flag storage.Flags
var storagehandler StorageHandler

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
		case storage.CounterType, storage.GuageType:
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

func ParseFlagsAndENV() storage.Flags {
	flag.StringVar(&Flag.FlagAddr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&Flag.FlagFileStorage, "f", "/tmp/metrics-db.json", "address of a file-storage")
	flag.IntVar(&Flag.FlagStoreInterval, "i", 300, "interval for storing data")
	flag.BoolVar(&Flag.FlagRestore, "r", true, "reading data from file first")
	flag.Parse()
	var envcfg storage.ServerENV
	err := env.Parse(&envcfg)
	if err != nil {
		log.Fatal(err)
	}

	if len(envcfg.Address) > 0 {
		Flag.FlagAddr = envcfg.Address
	}

	if len(envcfg.FileStorage) > 0 {
		Flag.FlagFileStorage = envcfg.FileStorage
	}
	restore := os.Getenv("RESTORE")
	if len(restore) > 0 {
		envrestore, err := strconv.ParseBool(restore)
		if err != nil {
			Flag.FlagRestore = envrestore
		}
	}
	if len(envcfg.StoreInterval) == 0 {
		storeInterval, err := strconv.Atoi(envcfg.StoreInterval)
		if err == nil {
			Flag.FlagStoreInterval = storeInterval
		}
	}
	return Flag
}
func GetSignal(s *http.Server, f storage.Flags) {
	storagehandler = StorageHandler{
		Storage: &storage.MemStorage{
			Gaugemem: make(map[string]storage.Gauge), Countermem: make(map[string]storage.Counter)},
	}

	storagehandler.Storage.ShutdownSave(s, f)
}

func MakeRouter(flag storage.Flags) *chi.Mux {
	storagehandler = StorageHandler{
		Storage: &storage.MemStorage{
			Gaugemem: make(map[string]storage.Gauge), Countermem: make(map[string]storage.Counter)},
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
					storage.WriteToFile(name, data)
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
		r.Get("/", WithLogging(GzipHandler(storagehandler.HandleGetAllMetrics())).ServeHTTP)
		r.Post("/update/", WithLogging(GzipHandler(storagehandler.HandleJSONPostMetrics(flag))).ServeHTTP)
		r.Post("/value/", WithLogging(GzipHandler(storagehandler.HandleJSONGetMetrics())).ServeHTTP)
		r.Post("/update/{type}/{name}/{value}", WithLogging(GzipHandler(storagehandler.HandlePostMetrics())).ServeHTTP)
		r.Get("/value/{name}/{type}", WithLogging(GzipHandler(storagehandler.HandleGetMetrics())).ServeHTTP)
	})
	return r
}

type StorageHandler struct {
	Storage storage.MemInteractions
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

func (st *StorageHandler) HandleGetAllMetrics() http.Handler {
	getMetrics := func(res http.ResponseWriter, req *http.Request) {
		result := st.Storage.GetAllMetrics()
		res.Header().Set("Content-Type", "text/html")
		io.WriteString(res, result)
	}
	return http.HandlerFunc(getMetrics)
}
