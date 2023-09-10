package postgres

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	appConfig "github.com/AntonyIS/notelify-user-service/config"
	"github.com/AntonyIS/notelify-user-service/internal/adapters/logger"
	"github.com/AntonyIS/notelify-user-service/internal/core/domain"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

type PostgresDBClient struct {
	db            *sql.DB
	tablename     string
	loggerService logger.LoggerType
}

func NewPostgresClient(appConfig appConfig.Config, logger logger.LoggerType) (*PostgresDBClient, error) {
	dbname := appConfig.DatabaseName
	tablename := appConfig.UserTable
	user := appConfig.DatabaseUser
	password := appConfig.DatabasePassword
	port := appConfig.DatabasePort
	host := appConfig.DatabaseHost
	region := appConfig.AWS_DEFAULT_REGION
	rdsInstanceIdentifier := appConfig.RDSInstanceIdentifier

	var dsn string

	if appConfig.Env == "dev" {
		dsn = fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable", host, port, user, dbname, password)
	} else {

		// Create a new AWS session
		awsSession := session.Must(session.NewSession(&aws.Config{
			Region: aws.String(region),
		}))

		// Create an RDS client
		rdsClient := rds.New(awsSession)

		// Describe the DB instance to get its endpoint
		describeInput := &rds.DescribeDBInstancesInput{
			DBInstanceIdentifier: &rdsInstanceIdentifier,
		}

		describeOutput, err := rdsClient.DescribeDBInstances(describeInput)

		if err != nil {
			logger.PostLogMessage(fmt.Sprintf("Failed to describe DB instance: %s", err.Error()))
		}

		if len(describeOutput.DBInstances) == 0 {
			logger.PostLogMessage("DB instance not found")
		}

		dsn = fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=require", host, port, dbname, user, password)
	}

	db, err := sql.Open("postgres", dsn)

	if err != nil {
		logger.PostLogMessage(err.Error())
		return nil, err
	}

	err = db.Ping()

	if err != nil {
		logger.PostLogMessage(err.Error())
		return nil, err
	}

	err = migrateDB(db, tablename)
	if err != nil {
		logger.PostLogMessage(err.Error())
		return nil, err

	}

	return &PostgresDBClient{db: db, tablename: tablename, loggerService: logger}, nil
}

func (psql *PostgresDBClient) CreateUser(user *domain.User) (*domain.User, error) {
	queryString := fmt.Sprintf(
		`INSERT INTO %s 
			(id,firstname,lastname,email,password,handle,about,contents,profile_image,following,followers) 
			VALUES 
			($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
		psql.tablename)
	_, err := psql.db.Exec(
		queryString,
		user.Id,
		user.Firstname,
		user.Lastname,
		user.Email,
		user.Password,
		user.Handle,
		user.About,
		pq.Array(user.Contents),
		user.ProfileImage,
		user.Following,
		user.Followers,
	)

	if err != nil {
		psql.loggerService.PostLogMessage(err.Error())
		return nil, err
	}

	return user, nil
}

func (psql *PostgresDBClient) ReadUser(id string) (*domain.User, error) {
	var user domain.User
	queryString := fmt.Sprintf(`SELECT id,firstname,lastname,email,handle,about,contents,profile_image,following,followers FROM %s WHERE id=$1`, psql.tablename)
	err := psql.db.QueryRow(queryString, id).Scan(&user.Id, &user.Firstname, &user.Lastname, &user.Email, &user.Handle, &user.About, pq.Array(&user.Contents), &user.ProfileImage, &user.Following, &user.Followers)

	if err != nil {
		psql.loggerService.PostLogMessage(fmt.Sprintf("user with id [%s] not found: %s", id, err.Error()))
		return nil, errors.New(fmt.Sprintf("user with id [%s] not found", id))
	}
	contentSvcURL := fmt.Sprintf("http://127.0.0.1:8081/v1/contents/author/%s", id)
	var contents []domain.Content
	contents, err = getUserContent(contentSvcURL)
	if err != nil {
		psql.loggerService.PostLogMessage("unable to read user content")
	}
	user.Contents = contents
	return &user, nil
}

func (psql *PostgresDBClient) ReadUsers() ([]domain.User, error) {

	rows, err := psql.db.Query(fmt.Sprintf("SELECT id,firstname,lastname,email,handle,about,contents,profile_image,following,followers FROM %s", psql.tablename))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []domain.User{}
	for rows.Next() {
		var user domain.User

		if err := rows.Scan(&user.Id, &user.Firstname, &user.Lastname, &user.Email, &user.Handle, &user.About, pq.Array(&user.Contents), &user.ProfileImage, &user.Following, &user.Followers); err != nil {
			psql.loggerService.PostLogMessage(err.Error())

			return nil, err
		}

		users = append(users, user)

	}
	return users, nil
}

func (psql *PostgresDBClient) ReadUserWithEmail(email string) (*domain.User, error) {
	var user domain.User
	queryString := fmt.Sprintf(`SELECT id,firstname,lastname,email,handle,about,contents,profile_image,following,followers FROM %s WHERE email=$1`, psql.tablename)
	err := psql.db.QueryRow(queryString, email).Scan(&user.Id, &user.Firstname, &user.Lastname, &user.Email, &user.Handle, &user.About, &user.Contents, &user.ProfileImage, &user.Following, &user.Followers)
	if err != nil {
		psql.loggerService.PostLogMessage(fmt.Sprintf("user with id [%s] not found: %s", email, err.Error()))
		return nil, errors.New(fmt.Sprintf("user with id [%s] not found", email))
	}

	return &user, nil
}

func (psql *PostgresDBClient) UpdateUser(user *domain.User) (*domain.User, error) {
	queryString := fmt.Sprintf(`UPDATE %s SET 
		firstname = $2,
		lastname = $3,
		handle = $4,
		about = $5,
		contents = $6,
		profile_image = $7,
		following = $8,
		followers = $9
	`, psql.tablename)

	_, err := psql.db.Exec(queryString, user.Firstname, user.Lastname, user.Handle, user.About, user.Contents, user.ProfileImage, user.Following, user.Followers)
	if err != nil {
		psql.loggerService.PostLogMessage(err.Error())
		return nil, err
	}
	return user, nil
}

func (psql *PostgresDBClient) DeleteUser(id string) (string, error) {

	queryString := fmt.Sprintf(`DELETE FROM %s WHERE id = $1`, psql.tablename)
	_, err := psql.db.Exec(queryString, id)
	if err != nil {
		psql.loggerService.PostLogMessage(err.Error())
		return "", err
	}
	return "Entity deleted successfully", nil
}

func migrateDB(db *sql.DB, userTable string) error {
	queryString := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			id VARCHAR(255) PRIMARY KEY UNIQUE,
			firstname VARCHAR(255) NOT NULL,
			lastname VARCHAR(255) NOT NULL,
			email VARCHAR(255) UNIQUE NOT NULL,
			password VARCHAR(255) UNIQUE NOT NULL,
			handle VARCHAR(255),
			about TEXT,
			contents TEXT [],
			profile_image varchar(255),
			Following int,
			Followers int
	)
	`, userTable)

	_, err := db.Exec(queryString)
	if err != nil {
		return err
	}

	return nil

}

func getUserContent(url string) ([]domain.Content, error) {
	// Makes a get HTTP request to the content service

	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return nil, err
	}

	res := string(body)

	var content []domain.Content

	err = json.Unmarshal([]byte(res), &content)
	if err != nil {
		return nil, err
	}

	return content, nil
}
