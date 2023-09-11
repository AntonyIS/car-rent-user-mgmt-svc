package app

import (
	"fmt"

	"github.com/AntonyIS/notelify-users-service/config"
	"github.com/AntonyIS/notelify-users-service/internal/adapters/logger"
	"github.com/AntonyIS/notelify-users-service/internal/core/ports"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func InitGinRoutes(svc ports.UserService, logger logger.LoggerType, conf config.Config) {
	gin.SetMode(gin.DebugMode)
	router := gin.Default()
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	router.Use(cors.New(config))

	handler := NewGinHandler(svc, conf.SECRET_KEY)

	usersRoutes := router.Group("/v1/users")

	{
		usersRoutes.GET("/", handler.ReadUsers)
		usersRoutes.GET("/:id", handler.ReadUser)
		usersRoutes.PUT("/:id", handler.UpdateUser)
		usersRoutes.DELETE("/:id", handler.DeleteUser)
		usersRoutes.DELETE("/delete/all", handler.DeleteAllUsers)
		usersRoutes.POST("/", handler.CreateUser)
		usersRoutes.POST("/login", handler.Login)
		usersRoutes.POST("/logout", handler.Logout)
	}
	logger.PostLogMessage(fmt.Sprintf("Server running on port :%s", conf.Port))
	router.Run(fmt.Sprintf(":%s", conf.Port))
}
