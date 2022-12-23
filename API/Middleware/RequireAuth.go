package Middleware

import (
	"FindMeATutor_User_Service/MongoDB"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"os"
	"time"
)

func RequireAuth(ctx *gin.Context) {
	client, c := MongoDB.ConnectToDatabase()
	var databaseName = os.Getenv("DATABASE_NAME")
	var databaseCollection = os.Getenv("DATABASE_COLLECTION")
	database := client.Database(databaseName)
	collection := database.Collection(databaseCollection)

	tokenString, err := ctx.Cookie("Authorization")
	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("SECRET")), nil
	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if float64(time.Now().Unix()) > claims["foo"].(float64) {
			ctx.AbortWithStatus(http.StatusUnauthorized)
		}
		var user MongoDB.User
		query := bson.D{bson.E{Key: "email", Value: claims["email"]}}
		err := collection.FindOne(c, query).Decode(&user)
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
		}
		ctx.Set("user", user)
		ctx.Next()
	} else {
		ctx.AbortWithStatus(http.StatusUnauthorized)
	}
}
