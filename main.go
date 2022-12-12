package main

import (
	"FindMeATutor_User_Service/API"
	"github.com/gin-gonic/gin"
)

func main() {
	server := gin.Default()
	basePath := server.Group("/v1")
	API.RegisterUserRoutes(basePath)
	err := server.Run()
	if err != nil {
		return
	}
}
