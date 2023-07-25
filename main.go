package main

import (
	"github.com/AntonyIS/notlify-user-svc/config"
	"github.com/AntonyIS/notlify-user-svc/internal/adapters/app"
	"github.com/AntonyIS/notlify-user-svc/internal/adapters/repository/postgres"
	"github.com/AntonyIS/notlify-user-svc/internal/core/services"
)

func main() {
	conf, err := config.NewConfig("dev")
	if err != nil {
		panic(err)
	}
	// // Postgres Client
	postgresDBrepo := postgres.NewPostgresClient(*conf)
	// // User service
	userSVC := services.NewUserManagementService(postgresDBrepo)
	// // Initialize HTTP server
	app.InitGinRoutes(userSVC, conf.Port)
}
