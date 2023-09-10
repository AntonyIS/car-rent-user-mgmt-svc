package app

import (
	"fmt"

	"github.com/AntonyIS/notelify-user-service/config"
	"github.com/AntonyIS/notelify-user-service/internal/adapters/logger"
	"github.com/AntonyIS/notelify-user-service/internal/core/ports"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func InitGinRoutes(svc ports.UserService, logger logger.LoggerType, conf config.Config) {
	gin.SetMode(gin.DebugMode)
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	handler := NewGinHandler(svc, conf.SECRET_KEY)

	usersRoutes := router.Group("/v1/users")

	// if conf.Env == "prod" {
	// 	middleware := NewMiddleware(svc, conf.SECRET_KEY)
	// 	usersRoutes.Use(middleware.Authorize)
	// }

	{
		usersRoutes.GET("/", handler.ReadUsers)
		usersRoutes.GET("/:id", handler.ReadUser)
		usersRoutes.PUT("/:id", handler.UpdateUser)
		usersRoutes.DELETE("/:id", handler.DeleteUser)
		usersRoutes.POST("/", handler.CreateUser)
		usersRoutes.POST("/login", handler.Login)
		usersRoutes.POST("/logout", handler.Logout)
	}
	logger.PostLogMessage(fmt.Sprintf("Server running on port :%s", conf.Port))
	router.Run(fmt.Sprintf(":%s", conf.Port))
}
