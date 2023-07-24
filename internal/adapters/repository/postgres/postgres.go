package postgres

import (
	"database/sql"
	"fmt"

	"github.com/AntonyIS/notlify-user-svc/config"
	"github.com/AntonyIS/notlify-user-svc/internal/core/domain"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/rds/rdsutils"
)

type PostgresDBClient struct {
	client *sql.DB
}

func NewPostgresClient(config *config.Config) *PostgresDBClient {
	dbEndpoint := fmt.Sprintf("%s:%d", config.DatabaseHost, config.DatabasePort)
	creds := credentials.NewEnvCredentials()
	authToken, err := rdsutils.BuildAuthToken(dbEndpoint, config.AWS_DEFAULT_REGION, config.DatabaseUser, creds)

	if err != nil {
		panic(err)
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?tls=true&allowCleartextPasswords=true",
		config.DatabaseUser, authToken, dbEndpoint, config.DatabaseName,
	)

	client, err := sql.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}

	err = client.Ping()
	if err != nil {
		panic(err)
	}

	return &PostgresDBClient{
		client: client,
	}
}

func (psql *PostgresDBClient) CreateUser(user *domain.User) (*domain.User, error) {
	return nil, nil
}

func (psql *PostgresDBClient) ReadUser(id string) (*domain.User, error) {
	return nil, nil
}

func (psql *PostgresDBClient) ReadUsers() ([]*domain.User, error) {
	return nil, nil
}

func (psql *PostgresDBClient) UpdateUser(user *domain.User) (*domain.User, error) {
	return nil, nil
}

func (psql *PostgresDBClient) DeleteUser(id string) (string, error) {
	return "", nil
}
