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
	collectionQuerys commentCollectionQueryable
}

func NewCommentService(commentCollectionQuery commentCollectionQueryable) *CommentService {
	return &CommentService{
		collectionQuerys: commentCollectionQuery,
	}
}

type commentCollectionQueryable interface {
	findOneComment(context.Context, interface{}) (*models.Comment, error)
	findCommentLike(context.Context, interface{}) (*[]bson.M, error)
	findPostComment(context.Context, interface{}) (*[]bson.M, error)
	insertOneComment(context.Context, interface{}) (*mongo.InsertOneResult, error)
	updateOneComment(context.Context, interface{}, interface{}) (*mongo.UpdateResult, error)
	deleteOneComment(context.Context, interface{}) (*mongo.DeleteResult, error)
}

type commentCollectionQuery struct {
	commentCollection *mongo.Collection
	likeCollection    *mongo.Collection
}

func NewCommentCollectionQuery(commentCollection *mongo.Collection, likeCollection *mongo.Collection) *commentCollectionQuery {
	return &commentCollectionQuery{
		commentCollection: commentCollection,
		likeCollection:    likeCollection,
	}
}

func (ccq *commentCollectionQuery) findOneComment(context context.Context, filter interface{}) (*models.Comment, error) {
	var comment models.Comment
	err := ccq.commentCollection.FindOne(context, filter).Decode(&comment)
	return &comment, err
}

func (ccq *commentCollectionQuery) findCommentLike(context context.Context, filter interface{}) (*[]bson.M, error) {
	cursor, err := ccq.likeCollection.Find(context, filter)
	if err != nil {
		return nil, err
	}
	var queryResult []bson.M
	if err = cursor.All(context, &queryResult); err != nil {
		return nil, err
	}
	return &queryResult, nil
}

func (ccq *commentCollectionQuery) findPostComment(context context.Context, filter interface{}) (*[]bson.M, error) {
	cursor, err := ccq.commentCollection.Find(context, filter)
	if err != nil {
		return nil, err
	}
	var queryResult []bson.M
	if err = cursor.All(context, &queryResult); err != nil {
		return nil, err
	}
	return &queryResult, nil
}

func (ccq *commentCollectionQuery) insertOneComment(context context.Context, document interface{}) (*mongo.InsertOneResult, error) {
	return ccq.commentCollection.InsertOne(context, document)
}

func (ccq *commentCollectionQuery) updateOneComment(context context.Context, filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
	return ccq.commentCollection.UpdateOne(context, filter, update)
}

func (ccq *commentCollectionQuery) deleteOneComment(context context.Context, filter interface{}) (*mongo.DeleteResult, error) {
	return ccq.commentCollection.DeleteOne(context, filter)
}

func (cs *CommentService) GetCommentUserId(commentId string) (string, error) {
	filter := bson.M{"_id": commentId}
	comment, err := cs.collectionQuerys.findOneComment(context.TODO(), filter)
	if err != nil {
		return "", err
	}
	return comment.UserId, nil
}

func (cs *CommentService) getCommentLikeCount(commentId string) (int, error) {
	filter := bson.M{"resource_id": commentId, "resource_type": "comment"}
	likes, err := cs.collectionQuerys.findCommentLike(context.TODO(), filter)
	if err != nil {
		return 0, err
	}
	return len(*likes), nil
}

func (cs *CommentService) FindAllPostComment(postId string) ([]models.Comment, error) {
	filter := bson.M{"post_id": postId}
	queryResult, err := cs.collectionQuerys.findPostComment(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	var comments []models.Comment
	for _, v := range *queryResult {
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
		comment := models.NewComment(id, postId, userId, commentContent, likeCount, createdDate, updatedDate)
		comments = append(comments, *comment)
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

	_, err := cs.collectionQuerys.insertOneComment(context.TODO(), newComment)
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
	_, err := cs.collectionQuerys.updateOneComment(context.TODO(), filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (cs *CommentService) DeleteComment(deletedCommentId string) error {
	filter := bson.M{"_id": deletedCommentId}
	_, err := cs.collectionQuerys.deleteOneComment(context.TODO(), filter)
	if err != nil {
		return err
	}
	return nil
}
