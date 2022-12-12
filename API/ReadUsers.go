package API

import (
	"FindMeATutor_User_Service/MongoDB"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetAllUsers(ctx *gin.Context) {
	users, err := MongoDB.GetAllUsers()
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, users)
}

func CreateUser(ctx *gin.Context) {
	var user MongoDB.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := MongoDB.CreateUser(&user)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "User created"})
}

func ReadUser(ctx *gin.Context) {
	email := ctx.Param("email")
	user, err := MongoDB.ReadUser(&email)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, user)
}

func UpdateUser(ctx *gin.Context) {
	var user MongoDB.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := MongoDB.UpdateUser(&user)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "User updated"})
}

func DeleteUser(ctx *gin.Context) {
	email := ctx.Param("email")
	err := MongoDB.DeleteUser(&email)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "User deleted"})
}

func RegisterUserRoutes(router *gin.RouterGroup) {
	router.POST("/createUser", CreateUser)
	router.GET("/readUser/:email", ReadUser)
	router.GET("/getAllUsers", GetAllUsers)
	router.PATCH("/updateUser", UpdateUser)
	router.DELETE("/deleteUser/:email", DeleteUser)
}
