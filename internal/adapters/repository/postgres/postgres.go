package postgres

import (
	"database/sql"
	"fmt"

	"github.com/AntonyIS/notlify-user-svc/config"
	"github.com/AntonyIS/notlify-user-svc/internal/core/domain"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/rds/rdsutils"
	_ "github.com/lib/pq"
)

type PostgresDBClient struct {
	db        *sql.DB
	tablename string
}

func NewPostgresClient(config config.Config) *PostgresDBClient {
	databaseName := config.DatabaseName
	databaseUserTable := config.UserTable
	databaseUser := config.DatabaseUser
	databasePassword := config.DatabasePassword
	databasePort := config.DatabasePort
	databaseHost := config.DatabaseHost
	databaseRegion := config.AWS_DEFAULT_REGION
	var dsn string

	if config.Env == "dev" {
		dsn = fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
			databaseHost,
			databasePort,
			databaseUser,
			databaseName,
			databasePassword,
		)
	} else {
		dbEndpoint := fmt.Sprintf("%s:%s", databaseHost, databasePort)
		creds := credentials.NewEnvCredentials()
		authToken, err := rdsutils.BuildAuthToken(dbEndpoint, databaseRegion, databaseUser, creds)

		if err != nil {
			panic(err)
		}

		dsn = fmt.Sprintf("%s:%s@tcp(%s)/%s?tls=true&allowCleartextPasswords=true",
			databaseUser, authToken, dbEndpoint, databaseName,
		)
	}

	db, err := sql.Open("postgres", dsn)

	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	// Create users table
	migrate(db, databaseUserTable)

	return &PostgresDBClient{
		db:        db,
		tablename: databaseUserTable,
	}
}

func (psql *PostgresDBClient) CreateUser(user *domain.User) (*domain.User, error) {
	var newUser *domain.User
	err := psql.db.QueryRow(
		`INSERT INTO %s (
			id,
			firstname,
			lastname,
			email,
			handle,
			profile_image,
			following,
			followers,
			social_media_links,
			reading_list,
			recommendations
		)
		VALUES 
		(
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
		) RETURNING 
			id,
			firstname,
			lastname,
			email,
			handle,
			profile_image,
			following,
			followers,
			social_media_links,
			reading_list,
			recommendations
		`, user.Id, user.Firstname, user.Lastname, user.Email, user.Handle, user.ProfileImage, user.Following, user.Followers, user.SocialMediaLinks, user.ReadingList, user.Recommendations).
		Scan(&newUser.Id, &newUser.Firstname, &newUser.Lastname, &newUser.Email, &newUser.Handle, &newUser.ProfileImage, &newUser.Following, &newUser.Followers, &newUser.SocialMediaLinks, &newUser.ReadingList, &newUser.Recommendations)

	if err != nil {
		return nil, err
	}
	return newUser, nil
}

func (psql *PostgresDBClient) ReadUser(id string) (*domain.User, error) {
	var user *domain.User
	err := psql.db.QueryRow("SELECT * FROM Users WHERE id=$1", id).Scan(&user.Id, &user.Firstname, &user.Lastname, &user.Email, &user.Handle, &user.ProfileImage, &user.Following, &user.Followers, &user.SocialMediaLinks, &user.ReadingList, &user.Recommendations)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (psql *PostgresDBClient) ReadUsers() ([]*domain.User, error) {
	rows, err := psql.db.Query(fmt.Sprintf("SELECT * FROM %s", psql.tablename))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	users := []*domain.User{}
	for rows.Next() {
		var user *domain.User
		if err := rows.Scan(&user.Id, &user.Firstname, &user.Lastname, &user.Email, &user.Handle, &user.ProfileImage, &user.Following, &user.Followers, &user.SocialMediaLinks, &user.ReadingList, &user.Recommendations); err != nil {
			return nil, err
		}

		users = append(users, user)

	}

	return users, nil
}

func (psql *PostgresDBClient) UpdateUser(user *domain.User) (*domain.User, error) {
	_, err := psql.db.Exec(`
		UPDATE users SET 
			firstname = $1,
			lastname = $2,
			email = $3,
			handle = $4,
			profile_image = $5,
			following = $6,
			followers = $7,
			social_media_links = $8,
			reading_list = $9,
			recommendations = $10
		`, user.Id, user.Firstname, user.Lastname, user.Email, user.Handle, user.ProfileImage, user.Following, user.Followers, user.SocialMediaLinks, user.ReadingList, user.Recommendations)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (psql *PostgresDBClient) DeleteUser(id string) (string, error) {
	_, err := psql.db.Exec("DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return "", err
	}
	return "Entity deleted successfully", nil
}

func migrate(db *sql.DB, userTable string) {
	// Creates new usertable if does not exists
	queryString := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			id VARCHAR(255) PRIMARY KEY UNIQUE,
			firstname VARCHAR(255) NOT NULL,
			lastname VARCHAR(255) NOT NULL,
			email VARCHAR(255) UNIQUE NOT NULL,
			handle VARCHAR(255) UNIQUE NOT NULL,
			about TEXT,
			profile_image varchar(255),
			Following TEXT[],
			Followers TEXT[],
			social_media_links TEXT[],
			reading_list TEXT[],
			recommendations TEXT[],
			blogs TEXT[]
	)
	
	`, userTable)

	_, err := db.Exec(queryString)
	if err != nil {
		panic(err)
	}

}
