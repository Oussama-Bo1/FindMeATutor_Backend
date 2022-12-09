package main

import (
	"FindMeATutor_User_Service/API"
	"context"
	"log"

	"FindMeATutor_User_Service/MongoDB"
	"go.mongodb.org/mongo-driver/mongo"
)

func main() {
	client, ctx := MongoDB.ConnectToDatabase()
	data := MongoDB.GetData(client, ctx)
	API.CreateAPI(data)
	defer func(client *mongo.Client, ctx context.Context) {
		var err = client.Disconnect(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}(client, ctx)
}
