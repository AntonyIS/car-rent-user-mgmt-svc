package dynamodb

import "github.com/AntonyIS/car-rent-user-mgmt-svc/internal/core/domain"

type DynamoDBClient struct {
	client string
}

func NewDynamoDBClient() *DynamoDBClient {
	return &DynamoDBClient{
		client: "client",
	}
}

func (db *DynamoDBClient) CreateUser(user *domain.User) (*domain.User, error) {
	return db.CreateUser(user)
}

func (db *DynamoDBClient) ReadUser(id string) (*domain.User, error) {
	return db.ReadUser(id)
}

func (db *DynamoDBClient) ReadUsers() ([]*domain.User, error) {
	return db.ReadUsers()
}

func (db *DynamoDBClient) UpdateUsers(user *domain.User) (*domain.User, error) {
	return db.UpdateUsers(user)
}

func (db *DynamoDBClient) DeleteUser(id string) error {
	return db.DeleteUser(id)
}
