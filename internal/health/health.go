package health

import (
	"log"

	"github.com/hootsuite/healthchecks"
	"github.com/pampatzoglou/api/config"
	"github.com/pampatzoglou/api/internal/mongo"
)

func CheckStatus() healthchecks.StatusList {
	cfg := config.New()
	mongoClient, ctx, cancel, err := mongo.Connect(cfg.Database.Connector)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	defer mongo.Close(mongoClient, ctx, cancel)
	pong := mongo.Ping(mongoClient, ctx)

	// Set a default response
	s := healthchecks.Status{
		Description: "mongo",
		Result:      healthchecks.OK,
	}

	// Make sure the pong response is what we expected
	if !pong {
		s = healthchecks.Status{
			Description: "mongo",
			Result:      healthchecks.CRITICAL,
		}
	}

	// Return our response
	return healthchecks.StatusList{StatusList: []healthchecks.Status{s}}
}
