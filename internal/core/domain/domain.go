package domain

import (
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type FollowUser struct {
	UserId       string `json:"user_id"`
	Firstname    string `json:"firstname"`
	Lastname     string `json:"lastname"`
	Email        string `json:"email"`
	Handle       string `json:"handle"`
	About        string `json:"about"`
	ProfileImage string `json:"profile_image"`
}

type User struct {
	UserId       string       `json:"user_id"`
	GitHubId     string       `json:"github_id"`
	LinkedInId   string       `json:"linkedin_id"`
	Firstname    string       `json:"firstname"`
	Lastname     string       `json:"lastname"`
	Email        string       `json:"email"`
	Password     string       `json:"password"`
	Handle       string       `json:"handle"`
	About        string       `json:"about"`
	Articles     []Article    `json:"articles"`
	ProfileImage string       `json:"profile_image"`
	Following    []FollowUser `json:"following"`
	Followers    []FollowUser `json:"followers"`
	AccessToken  string       `json:"access_token"`
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

type LogMessage struct {
	LogLevel string `json:"log_level"`
	Message  string `json:"message"`
	Service  string `json:"service"`
}

type GithubUser struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Firstname   string `json:"firstname"`
	Lastname    string `json:"lastname"`
	AvatarURL   string `json:"avatar_url"`
	AccessToken string `json:"access_token"`
	Email       string `json:"email"`
	Handle      string `json:"handle"`
}

func (g *GithubUser) InitGithubUser() User {
	nameParts := strings.Split(g.Name, " ")
	if len(nameParts) >= 2 {
		g.Firstname = nameParts[0]
		g.Lastname = strings.Join(nameParts[1:], " ")
	}

	user := User{
		UserId:       "",
		GitHubId:     strconv.Itoa(g.ID),
		LinkedInId:   "",
		Firstname:    g.Firstname,
		Lastname:     g.Lastname,
		Email:        "",
		Password:     "",
		Handle:       g.Handle,
		About:        "",
		Articles:     []Article{},
		ProfileImage: g.AvatarURL,
		Following:    []FollowUser{},
		Followers:    []FollowUser{},
		AccessToken:  g.AccessToken,
	}

	return user
}
