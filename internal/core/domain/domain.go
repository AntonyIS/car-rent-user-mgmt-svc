package domain

import "golang.org/x/crypto/bcrypt"

type User struct {
	Id           string `json:"id"`
	Firstname    string `json:"firstname"`
	Lastname     string `json:"lastname"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	Handle       string `json:"handle"`
	About        string `json:"about"`
	ProfileImage string `json:"profile_image"`
	Following    int    `json:"following"`
	Followers    int    `json:"followers"`
}

func (u User) CheckPasswordHarsh(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return false
	}
	return true
}

