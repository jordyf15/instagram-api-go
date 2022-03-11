package models

type User struct {
	Id              string           `json:"id" bson:"_id"`
	Username        string           `json:"username" bson:"username"`
	Fullname        string           `json:"fullname" bson:"full_name"`
	Password        string           `json:"password" bson:"password"`
	Email           string           `json:"email" bson:"email"`
	ProfilePictures []ProfilePicture `json:"profile_pictures" bson:"profile_pictures"`
}

func NewUser(id string, username string, fullname string, password string, email string, profilePictures []ProfilePicture) *User {
	return &User{
		Id:              id,
		Username:        username,
		Fullname:        fullname,
		Password:        password,
		Email:           email,
		ProfilePictures: profilePictures,
	}
}

type ProfilePicture struct {
	Type string `json:"type"`
	Size string `json:"size"`
	Url  string `json:"url"`
}

func NewProfilePicture(profilePictureType string, size string, url string) *ProfilePicture {
	return &ProfilePicture{
		Type: profilePictureType,
		Size: size,
		Url:  url,
	}
}
