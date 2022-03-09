package models

import "time"

type Post struct {
	Id              string    `json:"id" bson:"_id"`
	UserId          string    `json:"user_id" bson:"user_id"`
	VisualMediaUrls []string  `json:"visual_media_urls" bson:"visual_media_urls"`
	Caption         string    `json:"caption" bson:"caption"`
	LikeCount       int       `json:"like_count" bson:"like_count"`
	CreatedDate     time.Time `json:"created_date" bson:"created_date"`
	UpdatedDate     time.Time `json:"updated_date" bson:"updated_date"`
}

func NewPost(id string, userId string, visualMediaUrls []string, caption string, likeCount int, createdDate time.Time, updatedDate time.Time) *Post {
	return &Post{
		Id:              id,
		UserId:          userId,
		VisualMediaUrls: visualMediaUrls,
		Caption:         caption,
		LikeCount:       likeCount,
		CreatedDate:     createdDate,
		UpdatedDate:     updatedDate,
	}
}
