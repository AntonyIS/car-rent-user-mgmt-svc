/*
Package name : repository
File name : dynamodb.go
Author : Antony Injila
Description :
	- Host dynamoDb database specific methods
*/

package dynamodb

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"


    "fmt"

	"github.com/AntonyIS/car-rent-user-mgmt-svc/internal/core/domain"
)



type DynamoDBClient struct {
	client *dynamodb.DynamoDB
	tableName string
}

func NewDynamoDBClient(tableName string) *DynamoDBClient {
	creds := credentials.NewStaticCredentials(os.Getenv("AWS_ACCESS_KEY_ID"), os.Getenv("AWS_SECRET_ACCESS_KEY"), "")

	sess := session.Must(session.NewSession(&aws.Config{
		Region:      aws.String(c.Region),
		Credentials: creds,
	}))
	// Create DynamoDB client
	client := dynamodb.New(sess)
	return &DynamoDBClient{
		client: client,
		tableName : tableName,
	}
}

func (db *DynamoDBClient) CreateUser(user *domain.User) (*domain.User, error) {
	entityParsed, err := dynamodbattribute.MarshalMap(user)
	if err != nil {
		return nil, err
	}

	input := &dynamodb.PutItemInput{
		Item:      entityParsed,
		TableName: aws.String(db.tableName),
	}

	_, err = db.client.PutItem(input)

	if err != nil {
		return nil, err
	}

	user, err = db.ReadUser(user.id)
	if err != nil {
		return nil, errs.Wrap(errors.New(fmt.Sprintf("%s: %s", internalServerError, err)), "adapters.repository.dynamodb.CreateUser")
	}

	return user, nil
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
