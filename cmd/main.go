package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"

	//"github.com/gin-gonic/gin"
	"os"
	"time"

	//"encoding/json"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/pampatzoglou/api/config"
	"github.com/pampatzoglou/api/graph"
	"github.com/pampatzoglou/api/graph/generated"
	"github.com/pampatzoglou/api/internal/mongo"

	log "github.com/sirupsen/logrus"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/heptiolabs/healthcheck"
)

func recordMetrics() {
	go func() {
		for {
			opsProcessed.Inc()
			time.Sleep(2 * time.Second)
		}
	}()
}

var (
	opsProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "myapp_processed_ops_total",
		Help: "The total number of processed events",
	})
)

func main() {

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

	// Release resource when the main
	// function is returned.
	defer mongo.Close(mongoClient, ctx, cancel)
	/*
		log.Trace("Something very low level.")
		log.Debug("Useful debugging information.")
		log.Info("Something noteworthy happened!")
		log.Warn("You should probably take a look at this.")
		log.Error("Something failed but I'm not quitting.")
	*/
	fs := http.FileServer(http.Dir("./web/static"))
	h := healthcheck.NewHandler()
	// Our app is not happy if we've got more than 100 goroutines running.
	h.AddLivenessCheck("goroutine-threshold", healthcheck.GoroutineCountCheck(100))
	//h.AddReadinessCheck("mongo", mongo.HealthCheck(mongoClient, ctx))
	err = http.ListenAndServe("0.0.0.0:"+cfg.Server.HealthPort, h)
	if err != nil {
		log.Fatal(err)
	}
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))
	http.Handle("/", fs)
	http.Handle("/playground", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)
	recordMetrics()
	http.Handle("/metrics", promhttp.Handler())

	log.Printf("connect to http://localhost:%s/ for GraphQL playground and start queries", cfg.Server.Port)

	example := `{
__schema {
	queryType {
		fields {
			name
		}
	  }
	}
}`
	fmt.Printf("get the schema by:\n%s", example)
	err = http.ListenAndServe(":"+cfg.Server.Port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
