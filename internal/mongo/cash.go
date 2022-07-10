package mongo

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

func addToCache(ctx context.Context, data map[string]interface{}) error {

	fmt.Println(data)

	redisClient := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})

	jsonString, err := json.Marshal(data)

	if err != nil {
		return err
	}

	//err = redisClient.Set(ctx, "products_cache", jsonString, 30*time.Second).Err()
	err = redisClient.Set("products_cache", jsonString, 30*time.Second).Err()
	if err != nil {
		return nil
	}

	return nil
}

func getFromCache(ctx context.Context) (bool, map[string]interface{}, error) {

	redisClient := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})

	productsCache, err := redisClient.Get("products_cache").Bytes()

	if err != nil {
		return false, nil, nil
	}

	res := map[string]interface{}{}

	err = json.Unmarshal(productsCache, &res)

	if err != nil {
		return false, nil, nil
	}

	return true, res, nil
}
