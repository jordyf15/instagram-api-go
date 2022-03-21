package domain

import (
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
)

type Like struct {
	Id           string `json:"id" bson:"_id"`
	UserId       string `json:"user_id" bson:"user_id"`
	ResourceId   string `json:"resource_id" bson:"resource_id"`
	ResourceType string `json:"resource_type" bson:"resource_type"`
}

func NewLike(id string, userId string, resourceId string, resourceType string) *Like {
	return &Like{
		Id:           id,
		UserId:       userId,
		ResourceId:   resourceId,
		ResourceType: resourceType,
	}
}

type LikeUsecase interface {
	InsertPostLike(string, string) error
	DeletePostLike(string, string) error
	InsertCommentLike(string, string) error
	DeleteCommentLike(string, string) error
}

type LikeRepository interface {
	InsertLike(*Like) error
	FindLikes(interface{}) (*[]bson.M, error)
	FindOneLike(string) (*Like, error)
	DeleteLike(string) error
}

type LikeHandler interface {
	PostLikePost(http.ResponseWriter, *http.Request)
	DeleteLikePost(http.ResponseWriter, *http.Request)
	PostCommentLike(http.ResponseWriter, *http.Request)
	DeleteCommentLike(http.ResponseWriter, *http.Request)
}
