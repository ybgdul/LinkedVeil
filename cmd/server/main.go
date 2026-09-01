package main

import (
	"linkedveil/controllers"
	"linkedveil/initializers"
	"linkedveil/middlewares"
	"linkedveil/models"

	"github.com/gin-gonic/gin"
)

func init() { 
	initializers.LoadEnvs()
	initializers.ConnectDB()
}


func main() { 
	ipRateLimiter := models.NewIpLimiter()
	router := gin.Default()

	router.Use(middlewares.RateLimitingMiddleware(ipRateLimiter))
	router.Use(middlewares.StartLogging)
	
	auth := router.Group("/user")
	{
		auth.POST("/signup", controllers.CreateUser)
		auth.POST("/login", controllers.LoginUser)
		auth.GET("/profile", middlewares.CheckAuth, controllers.GetUserProfile)
	}

	protected := router.Group("/api")
	protected.Use(middlewares.CheckAuth)
	{
		auth.GET("/profile", controllers.GetUserProfile)
	}

	router.Run()
}
