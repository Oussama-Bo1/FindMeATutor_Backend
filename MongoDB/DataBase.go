package MongoDB

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
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
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		return err
	}
	user.Password = string(hash)
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
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		return err
	}
	user.Password = string(hash)
	filter := bson.D{bson.E{Key: "email", Value: user.Email}}
	update := bson.D{bson.E{Key: "$set", Value: bson.D{bson.E{Key: "aboutMe", Value: user.AboutMe}, bson.E{Key: "address", Value: user.Address}, bson.E{Key: "birthDate", Value: user.BirthDate}, bson.E{Key: "currentStudents", Value: user.CurrentStudents}, bson.E{Key: "currentTeachers", Value: user.CurrentTeachers}, bson.E{Key: "email", Value: user.Email}, bson.E{Key: "firstName", Value: user.FirstName}, bson.E{Key: "identity_verified", Value: user.IdentityVerified}, bson.E{Key: "lastName", Value: user.LastName}, bson.E{Key: "memberSince", Value: user.MemberSince}, bson.E{Key: "password", Value: user.Password}, bson.E{Key: "phone", Value: user.Phone}, bson.E{Key: "picture_url", Value: user.PictureUrl}, bson.E{Key: "price_per_hour", Value: user.PricePerHour}, bson.E{Key: "reviews", Value: user.Reviews}, bson.E{Key: "role", Value: user.Role}, bson.E{Key: "searching", Value: user.Searching}, bson.E{Key: "searchingForSubjects", Value: user.SearchingForSubjects}, bson.E{Key: "skills", Value: user.Skills}, bson.E{Key: "username", Value: user.Username}, bson.E{Key: "requests", Value: user.Requests}, bson.E{Key: "appointments", Value: user.Appointments}}}}
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

func LoginUser(email *string, password string) (error, string) {
	var user User
	client, ctx := ConnectToDatabase()
	var databaseName = os.Getenv("DATABASE_NAME")
	var databaseCollection = os.Getenv("DATABASE_COLLECTION")

	database := client.Database(databaseName)
	collection := database.Collection(databaseCollection)
	query := bson.D{bson.E{Key: "email", Value: email}}
	err := collection.FindOne(ctx, query).Decode(&user)
	if err != nil {
		return err, ""
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return err, ""
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": user.Email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		return err, ""
	}
	return nil, tokenString
}
