package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	conf, err := NewConfig("dev")
	if err != nil {
		panic(err)
	}
	fmt.Println(*conf)
}

type Config struct {
	Env               string
	Port              string
	UserTable         string
	AWS_ACCESS_KEY    string
	AWS_SECRET_KEY_ID string
	Debugging         bool
	Testing           bool
}

func NewConfig(Env string) (*Config, error) {
	err := godotenv.Load(".env")

	if err != nil {
		return nil, err
	}
	var (
		AWS_ACCESS_KEY    = os.Getenv("AWS_ACCESS_KEY")
		AWS_SECRET_KEY_ID = os.Getenv("AWS_SECRET_KEY_ID")
		Port              = "8080"
		Testing           = false
		Debugging         = false
		UserTable         = "DevUserTable"
	)

	switch Env {
	case "testing":
		Testing = true
		Debugging = true
		UserTable = "DevUserTable"
	case "dev":
		Testing = true
		Debugging = true
		UserTable = "DevUserTable"
	case "prod":
		Testing = false
		Debugging = false
		UserTable = "UserTable"
	}

	config := Config{
		Env:               Env,
		Port:              Port,
		UserTable:         UserTable,
		AWS_ACCESS_KEY:    AWS_ACCESS_KEY,
		AWS_SECRET_KEY_ID: AWS_SECRET_KEY_ID,
		Debugging:         Debugging,
		Testing:           Testing,
	}

	return &config, nil
}
