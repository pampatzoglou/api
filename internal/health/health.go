package health

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/hellofresh/health-go/v4"
	"github.com/pampatzoglou/api/config"
	"github.com/pampatzoglou/api/internal/mongo"
)

func main() {
	cfg := config.New()
	mongoClient, ctx, cancel, err := mongo.Connect(cfg.Database.Connector)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	defer mongo.Close(mongoClient, ctx, cancel)

	h, _ := health.New()
	h.Register(health.Config{
		Name:      "mongo-check",
		Timeout:   time.Second * 5,
		SkipOnErr: true,
		Check: func(ctx context.Context) error {
			mongo.Ping(mongoClient, ctx)
			return nil
		},
	})

	r := chi.NewRouter()
	r.Get("/status", h.HandlerFunc)
	http.ListenAndServe(":"+cfg.Server.HealthPort, nil)
}
