package app

import (
	"fmt"
	"time"

	"github.com/AntonyIS/notelify-users-service/config"
	"github.com/AntonyIS/notelify-users-service/internal/core/ports"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func InitGinRoutes(svc ports.UserService, logger ports.Logger, conf config.Config) {
	gin.SetMode(gin.DebugMode)

	router := gin.Default()
	router.Use(ginRequestLogger(logger))
	if conf.Env == "prod" {
		router.Use(cors.New(cors.Config{
			AllowOrigins:     []string{"http://notelify-client-service:3000", "http://localhost:3000"},
			AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
			AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
			ExposeHeaders:    []string{"Content-Length"},
			AllowCredentials: true,
		}))

	} else {
		router.Use(cors.New(cors.Config{
			AllowOrigins:     []string{"http://localhost:3000"},
			AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
			AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
			ExposeHeaders:    []string{"Content-Length"},
			AllowCredentials: true,
		}))

	}

	handler := NewGinHandler(svc, conf.SECRET_KEY)

	usersRoutes := router.Group("/v1/users")

	middleware := NewMiddleware(svc,conf.SECRET_KEY)

	usersRoutes.Use(middleware.Authorize)

	{
		usersRoutes.GET("/", handler.ReadUsers)
		usersRoutes.GET("/:user_id", handler.ReadUser)
		usersRoutes.PUT("/:user_id", handler.UpdateUser)
		usersRoutes.DELETE("/:user_id", handler.DeleteUser)
		usersRoutes.DELETE("/delete/all", handler.DeleteAllUsers)
		usersRoutes.POST("/", handler.CreateUser)
		usersRoutes.POST("/login", handler.Login)
		usersRoutes.POST("/logout", handler.Logout)
	}

	logger.Info(fmt.Sprintf("Server running on port :%s", conf.Port))
	router.Run(fmt.Sprintf(":%s", conf.Port))
}

func ginRequestLogger(logger ports.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		end := time.Now()
		latency := end.Sub(start)
		logger.Info(fmt.Sprintf("%s %s %s %d %s %s",
			c.Request.Method,
			c.Request.URL.Path,
			c.Request.Proto,
			c.Writer.Status(),
			latency.String(),
			c.ClientIP(),
		))
	}
}
