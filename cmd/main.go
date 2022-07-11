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

	"github.com/go-redis/redis"
	// https://betterprogramming.pub/graphql-subscriptions-with-go-6eb25dec5cd1
	//https://github.com/scorpionknifes/gqlmanage
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
	err := prometheus.Register(totalRequests)
	if err != nil {
		panic(err)
	}
	err = prometheus.Register(responseStatus)
	if err != nil {
		panic(err)
	}
	err = prometheus.Register(httpDuration)
	if err != nil {
		panic(err)
	}
}

func main() {

	//Redis
	//https://tutorialedge.net/golang/go-redis-tutorial/

	fmt.Println("Testing Golang Redis")
	client := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})

	pong, err := client.Ping().Result()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(pong)
	//End Redis

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
	routerMetrics := mux.NewRouter()
	routerMetrics.Use(prometheusMiddleware)
	routerMetrics.Handle("/metrics", promhttp.Handler())

	router := mux.NewRouter()
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))
	router.Handle("/", fs)
	router.Handle("/playground", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", srv)
	router.Handle("/health", h.Handler())

	//log.Printf("connect to http://localhost:%s/ for GraphQL playground and start queries", cfg.Server.Port)

	// go http.ListenAndServe(":9000", routerMetrics)
	go log.Fatal(http.ListenAndServe(":8000", routerMetrics))
	fmt.Println("Serving requests on port 9000")
	fmt.Println("Serving requests on port 8000")
	// http.ListenAndServe(":8000", router)
	log.Fatal(http.ListenAndServe(":8000", router))
	//select {} // block forever to prevent exiting
}
