package MongoDB

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"time"
)

// ConnectToDatabase
// Connects to a given database
func ConnectToDatabase() (*mongo.Client, context.Context) {
	var databaseUri = os.Getenv("DATABASE_URI")

	client, err := mongo.NewClient(options.Client().ApplyURI(databaseUri))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	return client, ctx
}

func GetData(client *mongo.Client, ctx context.Context) []bson.M {
	var databaseName = os.Getenv("DATABASE_NAME")
	var databaseCollection = os.Getenv("DATABASE_COLLECTION")

	database := client.Database(databaseName)
	collection := database.Collection(databaseCollection)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	var data []bson.M
	if err = cursor.All(ctx, &data); err != nil {
		log.Fatal(err)
	}
	fmt.Println(data)
	return data
}
