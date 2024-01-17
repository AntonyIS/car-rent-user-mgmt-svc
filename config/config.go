package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ENV                 string
	SERVER_PORT         string
	ARTICLE_TABLE       string
	LOGGER_URL          string
	SECRET_KEY          string
	POSTGRES_DB         string
	POSTGRES_USER       string
	POSTGRES_HOST       string
	POSTGRES_PORT       string
	POSTGRES_PASSWORD   string
	ARTICLE_SERVICE_URL string
	DEBUG               bool
	TEST                bool
}

var ENV string

func NewConfig() (*Config, error) {
	ENV := os.Getenv("ENV")

	if ENV == "test" {
		err := godotenv.Load("../../../.env")
		if err != nil {
			return nil, err
		}
	} else {
		err := godotenv.Load(".env")
		if err != nil {
			return nil, err
		}
		ENV = os.Getenv("ENV")
	}

	var (
		SECRET_KEY          = os.Getenv("SECRET_KEY")
		POSTGRES_PASSWORD   = os.Getenv("POSTGRES_PASSWORD")
		POSTGRES_USER       = "postgres"
		POSTGRES_DB         = "postgres"
		POSTGRES_HOST       = "postgres"
		POSTGRES_PORT       = "5432"
		SERVER_PORT         = "8000"
		ARTICLE_TABLE       = "User"
		LOGGER_URL          = "http://localhost:8003/v1/logger"
		ARTICLE_SERVICE_URL = "http://articles"
		DEBUG               = false
		TEST                = false
	)

	switch ENV {
	case "production":
		TEST = false
		DEBUG = false

	case "production_test":
		TEST = true
		DEBUG = true
		ARTICLE_TABLE = "TestUser"

	case "developement_test":
		TEST = true
		DEBUG = true
		ARTICLE_TABLE = "TestUser"

	case "development":
		TEST = true
		DEBUG = true
		POSTGRES_HOST = "localhost"
		ARTICLE_TABLE = "DevUser"

	case "docker":
		TEST = true
		DEBUG = true
		ARTICLE_TABLE = "DockerUser"
	case "docker_test":
		TEST = true
		DEBUG = true
		ARTICLE_TABLE = "DockerUser"
	}

	config := Config{
		ENV:                 ENV,
		SERVER_PORT:         SERVER_PORT,
		ARTICLE_TABLE:       ARTICLE_TABLE,
		SECRET_KEY:          SECRET_KEY,
		LOGGER_URL:          LOGGER_URL,
		DEBUG:               DEBUG,
		TEST:                TEST,
		POSTGRES_DB:         POSTGRES_DB,
		POSTGRES_USER:       POSTGRES_USER,
		POSTGRES_HOST:       POSTGRES_HOST,
		POSTGRES_PORT:       POSTGRES_PORT,
		POSTGRES_PASSWORD:   POSTGRES_PASSWORD,
		ARTICLE_SERVICE_URL: ARTICLE_SERVICE_URL,
	}

	return &config, nil
}
