package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ENV                   string
	SERVER_PORT           string
	USER_TABLE            string
	LOGGER_URL            string
	SECRET_KEY            string
	POSTGRES_DB           string
	POSTGRES_USER         string
	POSTGRES_HOST         string
	POSTGRES_PORT         string
	POSTGRES_PASSWORD     string
	ARTICLE_SERVICE_URL   string
	GITHUB_CLIENT_ID      string
	GITHUB_CLIENT_SECRET  string
	GITHUB_REDIRECT_URL   string
	LINKEDIN_REDIRECT_URL string
	DEBUG                 bool
	TEST                  bool
}

func NewConfig() (*Config, error) {
	ENV := os.Getenv("ENV")
	switch ENV {
	case "development":
		err := godotenv.Load(".env")
		if err != nil {
			return nil, err
		}
	}

	var (
		SECRET_KEY            = os.Getenv("SECRET_KEY")
		POSTGRES_PASSWORD     = os.Getenv("POSTGRES_PASSWORD")
		GITHUB_CLIENT_ID      = os.Getenv("GITHUB_CLIENT_ID")
		GITHUB_CLIENT_SECRET  = os.Getenv("GITHUB_CLIENT_SECRET")
		POSTGRES_USER         = "postgres"
		POSTGRES_DB           = "postgres"
		POSTGRES_HOST         = "postgres"
		POSTGRES_PORT         = "5432"
		SERVER_PORT           = "8000"
		USER_TABLE            = "Users"
		LOGGER_URL            = "http://logger:8002/logger/v1/users"
		ARTICLE_SERVICE_URL   = "http://articles:8001/posts/v1"
		GITHUB_REDIRECT_URL   = "http://users:3000/github/oauth2/callback"
		LINKEDIN_REDIRECT_URL = "http://users:3000/linkedin/oauth2/callback"
		DEBUG                 = false
		TEST                  = false
	)

	switch ENV {
	case "production":
		TEST = false
		DEBUG = false
		GITHUB_REDIRECT_URL = "http://users:3000/github/oauth2/callback"
		LINKEDIN_REDIRECT_URL = "http://users:3000/linkedin/oauth2/callback"

	case "production_test":
		TEST = true
		DEBUG = true
		USER_TABLE = "ProductionTestUsers"

	case "development":
		TEST = true
		DEBUG = true
		POSTGRES_HOST = "localhost"
		USER_TABLE = "DevUsers"
		LOGGER_URL = "http://localhost:8002/logger/v1/users"
		ARTICLE_SERVICE_URL = "http://localhost:8001/posts/v1"
		GITHUB_REDIRECT_URL = "http://localhost:3000/github/oauth2/callback"
		LINKEDIN_REDIRECT_URL = "http://localhost:3000/linkedin/oauth2/callback"

	case "development_test":
		TEST = true
		DEBUG = true
		SECRET_KEY = "testsecret"
		POSTGRES_PASSWORD = "pass1234"
		POSTGRES_HOST = "localhost"
		USER_TABLE = "TestUsers"
		LOGGER_URL = "http://localhost:8002/logger/v1/users"
		ARTICLE_SERVICE_URL = "http://localhost:8001/posts/v1"
		GITHUB_REDIRECT_URL = "http://localhost:3000/github/oauth2/callback"
		LINKEDIN_REDIRECT_URL = "http://localhost:3000/linkedin/oauth2/callback"

	case "docker":
		TEST = true
		DEBUG = true
		USER_TABLE = "DockerUsers"
		LOGGER_URL = "http://logger:8002/logger/v1/users"

	case "docker_test":
		TEST = true
		DEBUG = true
		USER_TABLE = "DockerUsers"
		LOGGER_URL = "http://logger:8002/logger/v1/users"
	}

	config := Config{
		ENV:                   ENV,
		SERVER_PORT:           SERVER_PORT,
		USER_TABLE:            USER_TABLE,
		SECRET_KEY:            SECRET_KEY,
		LOGGER_URL:            LOGGER_URL,
		DEBUG:                 DEBUG,
		TEST:                  TEST,
		POSTGRES_DB:           POSTGRES_DB,
		POSTGRES_USER:         POSTGRES_USER,
		POSTGRES_HOST:         POSTGRES_HOST,
		POSTGRES_PORT:         POSTGRES_PORT,
		POSTGRES_PASSWORD:     POSTGRES_PASSWORD,
		ARTICLE_SERVICE_URL:   ARTICLE_SERVICE_URL,
		GITHUB_CLIENT_ID:      GITHUB_CLIENT_ID,
		GITHUB_CLIENT_SECRET:  GITHUB_CLIENT_SECRET,
		GITHUB_REDIRECT_URL:   GITHUB_REDIRECT_URL,
		LINKEDIN_REDIRECT_URL: LINKEDIN_REDIRECT_URL,
	}

	return &config, nil
}
