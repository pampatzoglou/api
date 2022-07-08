package mongo

import (
	"context"
	"fmt"
	"log"

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
	collection := db.client.Database(DATABASE).Collection(COLLECTION)
	cursor, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.TODO())
	var results []*model.Shop
	for cursor.Next(context.TODO()) {
		var v *model.Shop
		err := cursor.Decode(&v)
		if err != nil {
			log.Fatal(err)
		}
		results = append(results, v)
	}
	return results
}
