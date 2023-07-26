package main

import (
	"flag"

	"github.com/AntonyIS/notlify-user-svc/config"
	"github.com/AntonyIS/notlify-user-svc/internal/adapters/app"
	"github.com/AntonyIS/notlify-user-svc/internal/adapters/repository/postgres"
	"github.com/AntonyIS/notlify-user-svc/internal/core/services"
)

var env string

func init() {
	flag.StringVar(&env, "env", "dev", "The environment the application is running")
	flag.Parse()
}

func main() {
	conf, err := config.NewConfig(env)
	if err != nil {
		panic(err)
	}
	// // Postgres Client
	postgresDBRepo, err := postgres.NewPostgresClient(*conf)
	if err != nil {
		panic(err)
	} else {
		// // User service
		userSVC := services.NewUserManagementService(postgresDBRepo)
		// // Initialize HTTP server
		app.InitGinRoutes(userSVC, *conf)
	}

}
