package MongoDB

import (
	"context"
	"errors"
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

func GetAllUsers() ([]*User, error) {
	client, ctx := ConnectToDatabase()
	var databaseName = os.Getenv("DATABASE_NAME")
	var databaseCollection = os.Getenv("DATABASE_COLLECTION")
	var users []*User

	database := client.Database(databaseName)
	collection := database.Collection(databaseCollection)
	cursor, err := collection.Find(ctx, bson.D{{}})
	if err != nil {
		return nil, err
	}
	for cursor.Next(ctx) {
		var user User
		err := cursor.Decode(&user)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	err = cursor.Close(ctx)
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, errors.New("no users found")
	}
	return users, nil
}

func CreateUser(user *User) error {
	client, ctx := ConnectToDatabase()
	var databaseName = os.Getenv("DATABASE_NAME")
	var databaseCollection = os.Getenv("DATABASE_COLLECTION")

	database := client.Database(databaseName)
	collection := database.Collection(databaseCollection)
	count, err := collection.CountDocuments(ctx, bson.D{bson.E{Key: "email", Value: user.Email}})
	if err != nil {
		return err
	}
	if count == 0 {
		_, err := collection.InsertOne(ctx, user)
		return err
	}
	return err
}

func ReadUser(email *string) (User, error) {
	var user User
	client, ctx := ConnectToDatabase()
	var databaseName = os.Getenv("DATABASE_NAME")
	var databaseCollection = os.Getenv("DATABASE_COLLECTION")

	database := client.Database(databaseName)
	collection := database.Collection(databaseCollection)
	query := bson.D{bson.E{Key: "email", Value: email}}
	err := collection.FindOne(ctx, query).Decode(&user)
	return user, err
}

func UpdateUser(user *User) error {
	client, ctx := ConnectToDatabase()
	var databaseName = os.Getenv("DATABASE_NAME")
	var databaseCollection = os.Getenv("DATABASE_COLLECTION")

	database := client.Database(databaseName)
	collection := database.Collection(databaseCollection)
	filter := bson.D{bson.E{Key: "email", Value: user.Email}}
	update := bson.D{bson.E{Key: "$set", Value: bson.D{bson.E{Key: "firstName", Value: user.FirstName}, bson.E{Key: "password", Value: user.Password}}}}
	result, _ := collection.UpdateOne(ctx, filter, update)
	if result.MatchedCount != 1 {
		return errors.New("no matched user found for update")
	}
	return nil
}

func DeleteUser(email *string) error {
	client, ctx := ConnectToDatabase()
	var databaseName = os.Getenv("DATABASE_NAME")
	var databaseCollection = os.Getenv("DATABASE_COLLECTION")

	database := client.Database(databaseName)
	collection := database.Collection(databaseCollection)
	filter := bson.D{bson.E{Key: "email", Value: email}}
	result, _ := collection.DeleteOne(ctx, filter)
	if result.DeletedCount != 1 {
		return errors.New("no matched user found for deletion")
	}
	return nil
}
