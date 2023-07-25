package domain

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
