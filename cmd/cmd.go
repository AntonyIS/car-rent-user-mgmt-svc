package cmd

import (
	"github.com/AntonyIS/notelify-users-service/config"
	"github.com/AntonyIS/notelify-users-service/internal/adapters/app"
	"github.com/AntonyIS/notelify-users-service/internal/adapters/repository/postgres"
	"github.com/AntonyIS/notelify-users-service/internal/core/domain"
	"github.com/AntonyIS/notelify-users-service/internal/core/services"
)

func RunService() {
	// Read application environment and load configurations
	conf, err := config.NewConfig()
	if err != nil {
		panic(err)
	}
	newLoggerService := services.NewLoggingManagementService(conf.LOGGER_URL)

	databaseRepo, err := postgres.NewPostgresClient(*conf)
	if err != nil {
		logEntry := domain.LogMessage{
			LogLevel: "ERROR",
			Service:  "users",
			Message:  err.Error(),
		}
		newLoggerService.LogError(logEntry)
		panic(err)
	}

	// Initialize the article service
	articleService := services.NewUserManagementService(databaseRepo, newLoggerService)
	// Run HTTP Server
	app.InitGinRoutes(articleService, newLoggerService, *conf)

}
