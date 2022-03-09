package models

import "time"

type Comment struct {
	Id          string    `json:"id" bson:"_id"`
	PostId      string    `json:"post_id" bson:"post_id"`
	UserId      string    `json:"user_id" bson:"user_id"`
	Comment     string    `json:"comment" bson:"comment"`
	LikeCount   int       `json:"like_count" bson:"like_count"`
	CreatedDate time.Time `json:"created_date" bson:"created_date"`
	UpdatedDate time.Time `json:"updated_date" bson:"updated_date"`
}
