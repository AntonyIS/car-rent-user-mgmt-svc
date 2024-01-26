package postgres

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/AntonyIS/notelify-users-service/config"
	"github.com/AntonyIS/notelify-users-service/internal/core/domain"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

type PostgresDBClient struct {
	db                 *sql.DB
	tablename          string
	articlesServiceURL string
}

func NewPostgresClient(appConfig config.Config) (*PostgresDBClient, error) {
	dbname := appConfig.POSTGRES_DB
	tablename := appConfig.USER_TABLE
	user := appConfig.POSTGRES_USER
	password := appConfig.POSTGRES_PASSWORD
	port := appConfig.POSTGRES_PORT
	host := appConfig.POSTGRES_HOST
	articlesServiceURL := appConfig.ARTICLE_SERVICE_URL

	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", host, port, user, dbname, password)

	db, err := sql.Open("postgres", dsn)

	if err != nil {
		return nil, err
	}

	err = db.Ping()

	if err != nil {
		return nil, err
	}

	err = migrateDB(db, tablename)
	if err != nil {
		return nil, err
	}

	return &PostgresDBClient{db: db, tablename: tablename, articlesServiceURL: articlesServiceURL}, nil
}

func (psql *PostgresDBClient) CreateUser(user *domain.User) (*domain.User, error) {
	query := fmt.Sprintf(
		`INSERT INTO %s 
			(user_id,firstname,lastname,email,password,handle,about,articles,profile_image,following,followers) 
			VALUES 
			($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
		psql.tablename)
	_, err := psql.db.Exec(
		query,
		user.UserId,
		user.Firstname,
		user.Lastname,
		user.Email,
		user.Password,
		user.Handle,
		user.About,
		pq.Array(user.Articles),
		user.ProfileImage,
		user.Following,
		user.Followers,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (psql *PostgresDBClient) ReadUserWithId(user_id string) (*domain.User, error) {
	var user domain.User
	queryString := fmt.Sprintf(`SELECT user_id,firstname,lastname,email,handle,about,articles,profile_image,following,followers FROM %s WHERE user_id=$1`, psql.tablename)
	err := psql.db.QueryRow(queryString, user_id).Scan(&user.UserId, &user.Firstname, &user.Lastname, &user.Email, &user.Handle, &user.About, pq.Array(&user.Articles), &user.ProfileImage, &user.Following, &user.Followers)
	if err != nil {
		return nil, err
	}
	articleSvcURL := fmt.Sprintf("%s/author/%s", psql.articlesServiceURL, user_id)
	var articles []domain.Article
	articles, err = getUserArticles(articleSvcURL)
	if err != nil {
		return nil, err
	}
	user.Articles = articles
	return &user, nil
}

func (psql *PostgresDBClient) ReadUsers() ([]domain.User, error) {
	rows, err := psql.db.Query(fmt.Sprintf("SELECT user_id,firstname,lastname,email,handle,about,articles,profile_image,following,followers FROM %s", psql.tablename))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []domain.User{}
	for rows.Next() {
		var user domain.User

		if err := rows.Scan(&user.UserId, &user.Firstname, &user.Lastname, &user.Email, &user.Handle, &user.About, pq.Array(&user.Articles), &user.ProfileImage, &user.Following, &user.Followers); err != nil {

			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (psql *PostgresDBClient) ReadUserWithEmail(email string) (*domain.User, error) {
	var user domain.User
	queryString := fmt.Sprintf(`SELECT user_id,firstname,lastname,email,handle,about,articles,profile_image,following,followers FROM %s WHERE email=$1`, psql.tablename)
	err := psql.db.QueryRow(queryString, email).Scan(&user.UserId, &user.Firstname, &user.Lastname, &user.Email, &user.Handle, &user.About, pq.Array(&user.Articles), &user.ProfileImage, &user.Following, &user.Followers)
	if err != nil {
		return nil, err
	}
	
	return &user, nil
}

func (psql *PostgresDBClient) UpdateUser(user *domain.User) (*domain.User, error) {
	queryString := fmt.Sprintf(`UPDATE %s SET 
		firstname = $2,
		lastname = $3,
		handle = $4,
		about = $5,
		articles = $6,
		profile_image = $7,
		following = $8,
		followers = $9
	`, psql.tablename)

	_, err := psql.db.Exec(queryString, user.Firstname, user.Lastname, user.Handle, user.About, user.Articles, user.ProfileImage, user.Following, user.Followers)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (psql *PostgresDBClient) DeleteUser(user_id string) (string, error) {
	queryString := fmt.Sprintf(`DELETE FROM %s WHERE user_id = $1`, psql.tablename)
	_, err := psql.db.Exec(queryString, user_id)
	if err != nil {
		return "", err
	}
	return "Entity deleted successfully", nil
}

func (psql *PostgresDBClient) DeleteAllUsers() (string, error) {
	queryString := fmt.Sprintf(`DELETE FROM %s`, psql.tablename)
	_, err := psql.db.Exec(queryString)
	if err != nil {
		return "", err
	}
	return "All items deletes successfully", nil
}

func migrateDB(db *sql.DB, userTable string) error {
	queryString := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			user_id VARCHAR(255) PRIMARY KEY UNIQUE,
			firstname VARCHAR(255) NOT NULL,
			lastname VARCHAR(255) NOT NULL,
			email VARCHAR(255) UNIQUE NOT NULL,
			password VARCHAR(255) UNIQUE NOT NULL,
			handle VARCHAR(255),
			about TEXT,
			articles TEXT [],
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

func getUserArticles(url string) ([]domain.Article, error) {
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

	var articles []domain.Article

	err = json.Unmarshal([]byte(res), &articles)
	if err != nil {
		return nil, err
	}

	return articles, nil
}
