package domain

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	UserId       string    `json:"user_id"`
	Firstname    string    `json:"firstname"`
	Lastname     string    `json:"lastname"`
	Email        string    `json:"email"`
	Password     string    `json:"password"`
	Handle       string    `json:"handle"`
	About        string    `json:"about"`
	Articles     []Article `json:"articles"`
	ProfileImage string    `json:"profile_image"`
	Following    int       `json:"following"`
	Followers    int       `json:"followers"`
}

type Article struct {
	ArticleID    string    `json:"article_id"`
	Title        string    `json:"title"`
	Subtitle     string    `json:"subtitle"`
	Introduction string    `json:"introduction"`
	Body         string    `json:"body"`
	Tags         []string  `json:"tags"`
	PublishDate  time.Time `json:"publish_date"`
	AuthorID     string    `json:"author_id"`
}

func (u User) CheckPasswordHarsh(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}
