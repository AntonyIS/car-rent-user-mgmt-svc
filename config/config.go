package config

import (
	"errors"
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	ErrNotFound       = errors.New("item not found")
	ErrInvalidItem    = errors.New("invalid item")
	ErrInternalServer = errors.New("internal server error")
)

type AppConfig struct {
	Env          string
	Port         string
	Region       string
	AWSSecretAccessKey string
	AWSAccessKeyID string
	UserTablename string
	Testing      bool
	Debugging      bool
}

func NewConfiguration(Env string) *AppConfig {

	var (
		awsAccessKeyID    = os.Getenv("AWS_ACCESS_KEY_ID")
		awsSecretAccessKey   = os.Getenv("AWS_SECRET_ACCESS_KEY")
		region           = os.Getenv("AWS_DEFAULT_REGION")
		serverPort       = os.Getenv("SERVER_PORT")
		devUserTablename = os.Getenv("DEV_USERS_TABLE")
		prodUserTablename = os.Getenv("PROD_USERS_TABLE")
		testing          = false
		debugging = false	
		userTablename = ""
	)

	switch Env {

	case "dev":
		userTablename = devUserTablename
		testing          = true
		debugging = true	

	case "pro":
		userTablename = prodUserTablename
		testing          = true
		debugging = true
	}	

	return &AppConfig{
		Env:          Env,
		AWSAccessKeyID : awsAccessKeyID,
		AWSSecretAccessKey : awsSecretAccessKey,
		Port:         serverPort,
		UserTablename:   userTablename,
		Region:       region,
		Testing:      testing,
		Debugging:      debugging,
	}
}

func LoadEnv() error {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}
	return nil
}