package mongo

import (
	"context"
	"fmt"
	"log"

	"github.com/mitchellh/mapstructure"
	"github.com/pampatzoglou/api/config"
	"github.com/pampatzoglou/api/graph/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ShopRepository interface {
	Save(shop *model.Shop)
	FindAll() []*model.Shop
}

type database struct {
	client *mongo.Client
}

const (
	DATABASE   = "db"
	COLLECTION = "shop"
)

func NewShopRepository() ShopRepository {
	// ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)

	// mongodb+srv://USERNAME:PASSWORD@HOST:PORT
	// MONGODB := os.Getenv("MONGODB")

	// Set client options
	// clientOptions := options.Client().ApplyURI(MONGODB)

	// clientOptions = clientOptions.SetMaxPoolSize(50)

	// Connect to MongoDB
	// userClient, err := mongo.Connect(ctx, clientOptions)

	cfg := config.New()
	mongoClient, ctx, _, err := Connect(cfg.Database.Connector)

	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = mongoClient.Ping(ctx, nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	return &database{
		client: mongoClient,
	}
}

func (db *database) Save(shop *model.Shop) {
	collection := db.client.Database(DATABASE).Collection(COLLECTION)
	_, err := collection.InsertOne(context.TODO(), shop)
	if err != nil {
		log.Fatal(err)
	}
}

func (db *database) FindAll() []*model.Shop {

	ctx := context.Background()

	isCached, productsCache, err := getFromCache(ctx)

	if err != nil {

		fmt.Println(err)

	} else if isCached {

		fmt.Println("is cached", productsCache["data"])

		var s []interface{} = productsCache["data"].([]interface{})
		fmt.Println(s)
		var shop *model.Shop
		var result []*model.Shop
		for _, v := range s { // use type assertion to loop over []interface{}
			err := mapstructure.Decode(v, &shop)
			if err != nil {
				panic(err)
			}
			result = append(result, shop)
		}
		fmt.Println(result)
		return result
		// Convert json string to struct
		//	var v *model.Shop
		//	if err := json.Unmarshal(productsCache, &v); err != nil {
		//		fmt.Println(err)
		//	}
		//var results []*model.Shop
		//	results = append(results, v)
		//return results

	}

	collection := db.client.Database(DATABASE).Collection(COLLECTION)
	cursor, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.TODO())

	var results []*model.Shop
	// var records []bson.M
	for cursor.Next(context.TODO()) {
		var v *model.Shop
		//	var record bson.M
		err := cursor.Decode(&v)
		if err != nil {
			log.Fatal(err)
		}
		results = append(results, v)
		// records = append(records, v)
	}

	res := map[string]interface{}{
		"data": results,
	}
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("not cached", res)

	err = addToCache(ctx, res)
	if err != nil {
		fmt.Println(err)
	}

	return results

}
