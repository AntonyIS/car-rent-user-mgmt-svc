package domain

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id           string    `json:"id"`
	Firstname    string    `json:"firstname"`
	Lastname     string    `json:"lastname"`
	Email        string    `json:"email"`
	Password     string    `json:"password"`
	Handle       string    `json:"handle"`
	About        string    `json:"about"`
	Contents     []Content `json:"contents"`
	ProfileImage string    `json:"profile_image"`
	Following    int       `json:"following"`
	Followers    int       `json:"followers"`
}

type Content struct {
	ContentId       string            `json:"content_id"`
	CreatorId       string            `json:"creator_id"`
	Title           string            `json:"title"`
	Body            string            `json:"body"`
	Images          map[string]string `json:"images"`
	Videos          map[string]string `json:"videos"`
	PublicationDate time.Time         `json:"publication_date"`
}

func (u User) CheckPasswordHarsh(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return false
	}
	return true
}
