package models

type User struct {
	Id              string           `json:"id"`
	Username        string           `json:"username"`
	Fullname        string           `json:"fullname"`
	Password        string           `json:"password"`
	Email           string           `json:"email"`
	ProfilePictures []ProfilePicture `json:"profile_pictures"`
}

type ProfilePicture struct {
	Type string `json:"type"`
	Size string `json:"size"`
	Url  string `json:"url"`
}
