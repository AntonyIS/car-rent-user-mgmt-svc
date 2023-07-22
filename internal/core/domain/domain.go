package domain

type User struct {
	Id               string   `json:"id"`
	Firstname        string   `json:"firstname"`
	Lastname         string   `json:"lastname"`
	Email            string   `json:"email"`
	Handle           string   `json:"handle"`
	ProfileImage     string   `json:"profile_image"`
	Following        int      `json:"following"`
	Followers        int      `json:"followers"`
	SocialMediaLinks []string `json:"social_media_links"`
	ReadingList      []string `json:"reading_list"`
	Recommendations  []string `json:"recommendations"`
}
