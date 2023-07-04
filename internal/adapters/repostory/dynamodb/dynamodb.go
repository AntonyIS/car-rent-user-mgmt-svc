package dynamodb


import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"

	"github.com/AntonyIS/car-rent-user-mgmt-svc/internal/core/domain"
	"github.com/AntonyIS/car-rent-user-mgmt-svc/config"

	"errors"
)



type DynamoDBClient struct {
	client *dynamodb.DynamoDB
	tableName string
}

func NewDynamoDBClient(c *config.AppConfig) *DynamoDBClient {
	creds := credentials.NewStaticCredentials(
		c.AWSAccessKeyID, 
		c.AWSSecretAccessKey, 
		"",
	)

	sess := session.Must(session.NewSession(&aws.Config{
		Region:      aws.String(c.Region),
		Credentials: creds,
	}))

	// Create DynamoDB client
	client := dynamodb.New(sess)
	return &DynamoDBClient{
		client: client,
		tableName : c.UserTablename,
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

	newUser, err := db.ReadUser(user.ID)
	if err != nil {
		return nil,err
	}

	return newUser, nil
}

func (db *DynamoDBClient) ReadUser(id string) (*domain.User, error) {
	result, err := db.client.GetItem(&dynamodb.GetItemInput{
        TableName: aws.String(db.tableName),
        Key: map[string]*dynamodb.AttributeValue{
            "ID": {
                N: aws.String(id),
            },
        },
    })
	if err != nil {
       return nil , err
    }

	user := domain.User{}

	err = dynamodbattribute.UnmarshalMap(result.Item, &user)

    if err != nil {
		return nil , err
	}

    if user.ID == "" {
        return nil , errors.New("Could not find user")
    }
	return &user, nil
}

func (db *DynamoDBClient) ReadUsers() ([]*domain.User, error) {
	users := []*domain.User{}
	filt := expression.Name("Id").AttributeNotExists()
	proj := expression.NamesList(
		expression.Name("id"),
		expression.Name("firstname"),
		expression.Name("lastname"),
		expression.Name("email"),
		expression.Name("projects"),
		expression.Name("payment"),
		expression.Name("vehicle"),
	
	)

	expr, err := expression.NewBuilder().WithFilter(filt).WithProjection(proj).Build()

	if err != nil {
		return nil, err
	}
	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(db.tableName),
	}

	result, err := db.client.Scan(params)

	if err != nil {
		return nil, err
	}

	for _, item := range result.Items {
		var user domain.User

		err = dynamodbattribute.UnmarshalMap(item, &user)
		if err != nil {
			return nil, err
		}

		users = append(users, &user)

	}

	return users, nil
}

func (db *DynamoDBClient) UpdateUsers(user *domain.User) (*domain.User, error) {
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

	return user, nil
}

func (db *DynamoDBClient) DeleteUser(id string) error {
	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"ID": {
				S: aws.String(id),
			},
		},
		TableName: aws.String(db.tableName),
	}

	res, err := db.client.DeleteItem(input)
	if res == nil {
		return err
	}
	if err != nil {
		return err
	}
	return nil
}
