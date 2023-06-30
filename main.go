package main

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}

}
func main() {
	fmt.Println("Car rent service")
}
