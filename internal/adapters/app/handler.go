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

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	handler := NewGinHandler(svc, logger, conf.SECRET_KEY)

	usersRoutes := router.Group("/v1/users")

	// middleware := NewMiddleware(svc, logger, conf.SECRET_KEY)

	// usersRoutes.Use(middleware.Authorize)

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

	logger.Info(fmt.Sprintf("Server running on port 0.0.0.0:%s", conf.SERVER_PORT))
	router.Run(fmt.Sprintf("0.0.0.0:%s", conf.SERVER_PORT))
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
