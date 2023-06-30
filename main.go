package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"github.com/AntonyIS/car-rent-user-mgmt-svc/config"
)
func init(){
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}

}
func main() {
	appConfig : config.NewConfiguration("dev")
	repo := repository.NewDynamoDBClient()
}
