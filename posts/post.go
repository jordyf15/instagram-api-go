package posts

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
