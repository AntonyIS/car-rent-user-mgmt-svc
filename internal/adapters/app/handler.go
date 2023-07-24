package app

import (
	"fmt"
	"os"

	"github.com/AntonyIS/notlify-user-svc/config"
	"github.com/AntonyIS/notlify-user-svc/internal/core/ports"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func InitGinRoutes(svc ports.UserService, config config.Config) {
	// Enable detailed error responses
	gin.SetMode(gin.DebugMode)

	// Setup Gin router
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Setup application route handlers
	handler := NewGinHandler(svc)
	

	usersRoutes := router.Group("/api/v1/users")
	{
		usersRoutes.GET("/", handler.ReadUsers)
		usersRoutes.GET("/:id", handler.ReadUser)
		usersRoutes.POST("/", handler.CreateUser)
		usersRoutes.PUT("/:id", handler.UpdateUser)
		usersRoutes.DELETE("/:id", handler.DeleteUser)
	}

	port := fmt.Sprintf(":%s", os.Getenv("SERVER_PORT"))

	router.Run(port)
}
