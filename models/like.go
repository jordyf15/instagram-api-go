package models

type Like struct {
	Id           string `json:"id" bson:"_id"`
	UserId       string `json:"user_id" bson:"user_id"`
	ResourceId   string `json:"resource_id" bson:"resource_id"`
	ResourceType string `json:"resource_type" bson:"resource_type"`
}
