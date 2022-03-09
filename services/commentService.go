package services

import (
	"context"
	"fmt"
	"instagram-go/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CommentService struct {
	commentCollection *mongo.Collection
	likeCollection    *mongo.Collection
}

func NewCommentService(commentCollection *mongo.Collection, likeCollection *mongo.Collection) *CommentService {
	return &CommentService{
		commentCollection: commentCollection,
		likeCollection:    likeCollection,
	}
}

func (cs *CommentService) GetCommentUserId(commentId string) (string, error) {
	var comment models.Comment
	filter := bson.M{"_id": commentId}
	err := cs.commentCollection.FindOne(context.TODO(), filter).Decode(&comment)
	if err != nil {
		return "", err
	}
	return comment.UserId, nil
}

func (cs *CommentService) getCommentLikeCount(commentId string) (int, error) {
	filter := bson.M{"resource_id": commentId, "resource_type": "comment"}
	cursor, err := cs.likeCollection.Find(context.TODO(), filter)
	if err != nil {
		return 0, err
	}
	var likes []bson.M
	if err = cursor.All(context.TODO(), &likes); err != nil {
		return 0, err
	}
	return len(likes), nil
}

func (cs *CommentService) FindAllPostComment(postId string) ([]models.Comment, error) {
	filter := bson.M{"post_id": postId}
	cursor, err := cs.commentCollection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	var queryResult []bson.M
	if err = cursor.All(context.TODO(), &queryResult); err != nil {
		return nil, err
	}
	var comments []models.Comment
	for _, v := range queryResult {
		id := fmt.Sprintf("%v", v["_id"])
		postId := fmt.Sprintf("%v", v["post_id"])
		userId := fmt.Sprintf("%v", v["user_id"])
		commentContent := fmt.Sprintf("%v", v["comment"])
		likeCount, err := cs.getCommentLikeCount(id)
		createdDate := v["created_date"].(primitive.DateTime).Time()
		updatedDate := v["updated_date"].(primitive.DateTime).Time()
		if err != nil {
			return nil, err
		}
		comment := models.Comment{id, postId, userId, commentContent, likeCount, createdDate, updatedDate}
		comments = append(comments, comment)
	}
	return comments, nil
}

func (cs *CommentService) InsertComment(comment models.Comment) error {
	newComment := bson.D{
		primitive.E{Key: "_id", Value: comment.Id},
		primitive.E{Key: "post_id", Value: comment.PostId},
		primitive.E{Key: "user_id", Value: comment.UserId},
		primitive.E{Key: "comment", Value: comment.Comment},
		primitive.E{Key: "created_date", Value: comment.CreatedDate},
		primitive.E{Key: "updated_date", Value: comment.UpdatedDate},
	}

	_, err := cs.commentCollection.InsertOne(context.TODO(), newComment)
	if err != nil {
		return err
	}
	return nil
}

func (cs *CommentService) UpdateComment(updatedCommentId string, newComment string) error {
	filter := bson.M{"_id": updatedCommentId}
	update := bson.D{primitive.E{
		Key: "$set",
		Value: bson.D{primitive.E{
			Key:   "comment",
			Value: newComment},
		},
	}}
	_, err := cs.commentCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (cs *CommentService) DeleteComment(deletedCommentId string) error {
	filter := bson.M{"_id": deletedCommentId}
	_, err := cs.commentCollection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}
	return nil
}
