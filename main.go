package main

import (
	"fmt"
	"os"

	"github.com/AntonyIS/notlify-user-svc/config"
	"github.com/joho/godotenv"
)

func main() {
	conf, err := NewConfig("dev")
	if err != nil {
		panic(err)
	}
	fmt.Println(*conf)
}

func NewConfig(Env string) (*config.Config, error) {
	err := godotenv.Load(".env")

	if err != nil {
		return nil, err
	}
	var (
		AWS_ACCESS_KEY     = os.Getenv("AWS_ACCESS_KEY")
		AWS_SECRET_KEY_ID  = os.Getenv("AWS_SECRET_KEY_ID")
		AWS_DEFAULT_REGION = os.Getenv("AWS_DEFAULT_REGION")
		Port               = "8080"
		Testing            = false
		Debugging          = false
		UserTable          = "DevUserTable"
		DatabaseName       = "Notlify"
		DatabaseUser       = "Antony"
		DatabaseHost       = os.Getenv("DatabaseHost")
		DatabasePort       = 3306
		DatabaseRegion     = os.Getenv("AWS_DEFAULT_REGION")
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

	config := config.Config{
		Env:                Env,
		Port:               Port,
		UserTable:          UserTable,
		AWS_ACCESS_KEY:     AWS_ACCESS_KEY,
		AWS_SECRET_KEY_ID:  AWS_SECRET_KEY_ID,
		AWS_DEFAULT_REGION: AWS_DEFAULT_REGION,
		Debugging:          Debugging,
		Testing:            Testing,
		DatabaseName:       DatabaseName,
		DatabaseUser:       DatabaseUser,
		DatabaseHost:       DatabaseHost,
		DatabasePort:       DatabasePort,
		DatabaseRegion:     DatabaseRegion,
	}

	return &config, nil
}
