package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/hellofresh/health-go/v4"
	"github.com/pampatzoglou/api/config"
	"github.com/pampatzoglou/api/graph"
	"github.com/pampatzoglou/api/graph/generated"
	"github.com/pampatzoglou/api/internal/mongo"

	log "github.com/sirupsen/logrus"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, http.StatusOK}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

var totalRequests = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Number of get requests.",
	},
	[]string{"path"},
)

var responseStatus = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "response_status",
		Help: "Status of HTTP response",
	},
	[]string{"status"},
)

var httpDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
	Name: "http_response_time_seconds",
	Help: "Duration of HTTP requests.",
}, []string{"path"})

func prometheusMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		route := mux.CurrentRoute(r)
		path, _ := route.GetPathTemplate()

		timer := prometheus.NewTimer(httpDuration.WithLabelValues(path))
		rw := NewResponseWriter(w)
		next.ServeHTTP(rw, r)

		statusCode := rw.statusCode

		responseStatus.WithLabelValues(strconv.Itoa(statusCode)).Inc()
		totalRequests.WithLabelValues(path).Inc()

		timer.ObserveDuration()
	})
}

func init() {
	prometheus.Register(totalRequests)
	prometheus.Register(responseStatus)
	prometheus.Register(httpDuration)
}

func main() {

	//
	log.Println("os.Args", os.Args)
	cfg := config.New()

	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	ll, err := log.ParseLevel(cfg.Server.LogLevel)
	if err != nil {
		ll = log.DebugLevel
	}
	// set global log level
	log.SetLevel(ll)

	mongoClient, ctx, cancel, err := mongo.Connect(cfg.Database.Connector)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	defer mongo.Close(mongoClient, ctx, cancel)
	fs := http.FileServer(http.Dir("./web/static"))
	h, _ := health.New()
	err = h.Register(health.Config{
		Name:      "mongo-check",
		Timeout:   time.Second * 5,
		SkipOnErr: true,
		Check: func(ctx context.Context) error {
			mongo.Ping(mongoClient, ctx)
			return nil
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	router := mux.NewRouter()
	router.Use(prometheusMiddleware)

	// Prometheus endpoint
	router.Handle("/prometheus", promhttp.Handler())

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))

	router.Handle("/", fs)
	router.Handle("/playground", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", srv)
	router.Handle("/metrics", promhttp.Handler())
	// http.Handle("/health", h.Handler())

	//log.Printf("connect to http://localhost:%s/ for GraphQL playground and start queries", cfg.Server.Port)

	fmt.Println("Serving requests on port 8000")
	err = http.ListenAndServe(":8000", router)

	//	err = http.ListenAndServe(":"+cfg.Server.Port, nil)
	if err != nil {
		log.Fatal(err)
	}

}
