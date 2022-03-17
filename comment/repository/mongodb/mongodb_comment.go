package mongodb

import (
	"context"
	"instagram-go/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongodbCommentRepository struct {
	collection *mongo.Collection
}

func NewMongodbCommentRepository(collection *mongo.Collection) *mongodbCommentRepository {
	return &mongodbCommentRepository{
		collection: collection,
	}
}

func (mcr *mongodbCommentRepository) FindComments(filter interface{}) (*[]bson.M, error) {
	cursor, err := mcr.collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	var queryResult []bson.M
	if err = cursor.All(context.TODO(), &queryResult); err != nil {
		return nil, err
	}
	return &queryResult, nil
}

func (mcr *mongodbCommentRepository) InsertComment(comment *domain.Comment) error {
	newComment := bson.D{
		primitive.E{Key: "_id", Value: comment.Id},
		primitive.E{Key: "post_id", Value: comment.PostId},
		primitive.E{Key: "user_id", Value: comment.UserId},
		primitive.E{Key: "comment", Value: comment.Comment},
		primitive.E{Key: "created_date", Value: comment.CreatedDate},
		primitive.E{Key: "updated_date", Value: comment.UpdatedDate},
	}

	_, err := mcr.collection.InsertOne(context.TODO(), newComment)
	return err
}

func (mcr *mongodbCommentRepository) FindOneComment(commentId string) (*domain.Comment, error) {
	var comment domain.Comment
	filter := bson.M{"_id": commentId}
	err := mcr.collection.FindOne(context.TODO(), filter).Decode(&comment)
	return &comment, err
}

func (mcr *mongodbCommentRepository) UpdateComment(commentId string, commentContent string) error {
	filter := bson.M{"_id": commentId}
	update := bson.D{primitive.E{
		Key: "$set",
		Value: bson.D{primitive.E{
			Key:   "comment",
			Value: commentContent},
		},
	}}
	_, err := mcr.collection.UpdateOne(context.TODO(), filter, update)
	return err
}

func (mcr *mongodbCommentRepository) DeleteComment(commentId string) error {
	filter := bson.M{"_id": commentId}
	_, err := mcr.collection.DeleteOne(context.TODO(), filter)
	return err
}
