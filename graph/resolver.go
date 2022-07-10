package graph

import (
	"github.com/go-redis/redis/v8"
	"github.com/pampatzoglou/api/graph/model"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	shops []*model.Shop
	Redis redis.Client
}
