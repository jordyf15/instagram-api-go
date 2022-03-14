package services

import (
	"context"
	"instagram-go/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type LikeService struct {
	collectionQuerys likeCollectionQueryable
}

type ILikeService interface {
	InsertLike(models.Like) error
	DeleteLike(string) error
	IsLikeExist(string, string, string) (bool, error)
	GetLikeUserId(string) (string, error)
	IsLikeExistById(string) (bool, error)
}

func NewLikeService(likeCollectionQuery likeCollectionQueryable) *LikeService {
	return &LikeService{collectionQuerys: likeCollectionQuery}
}

type likeCollectionQueryable interface {
	insertOne(context.Context, interface{}) (*mongo.InsertOneResult, error)
	deleteOne(context.Context, interface{}) (*mongo.DeleteResult, error)
	findOne(context.Context, interface{}) (*models.Like, error)
	find(context.Context, interface{}) (*[]bson.M, error)
}

type likeCollectionQuery struct {
	collection *mongo.Collection
}

func NewLikeCollectionQuery(collection *mongo.Collection) *likeCollectionQuery {
	return &likeCollectionQuery{
		collection: collection,
	}
}
func (lcq *likeCollectionQuery) insertOne(context context.Context, document interface{}) (*mongo.InsertOneResult, error) {
	return lcq.collection.InsertOne(context, document)
}
func (lcq *likeCollectionQuery) deleteOne(context context.Context, filter interface{}) (*mongo.DeleteResult, error) {
	return lcq.collection.DeleteOne(context, filter)
}
func (lcq *likeCollectionQuery) findOne(context context.Context, filter interface{}) (*models.Like, error) {
	var like models.Like
	err := lcq.collection.FindOne(context, filter).Decode(&like)
	return &like, err
}
func (lcq *likeCollectionQuery) find(context context.Context, filter interface{}) (*[]bson.M, error) {
	cursor, err := lcq.collection.Find(context, filter, options.Find().SetLimit(1))
	if err != nil {
		return nil, err
	}
	var queryResult []bson.M
	if err = cursor.All(context, &queryResult); err != nil {
		return nil, err
	}
	return &queryResult, nil
}

func (ls *LikeService) InsertLike(like models.Like) error {
	newLike := bson.D{
		primitive.E{Key: "_id", Value: like.Id},
		primitive.E{Key: "user_id", Value: like.UserId},
		primitive.E{Key: "resource_id", Value: like.ResourceId},
		primitive.E{Key: "resource_type", Value: like.ResourceType},
	}
	_, err := ls.collectionQuerys.insertOne(context.TODO(), newLike)
	if err != nil {
		return err
	}
	return nil
}

func (ls *LikeService) DeleteLike(likeId string) error {
	filter := bson.M{"_id": likeId}
	_, err := ls.collectionQuerys.deleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}
	return nil
}

func (ls *LikeService) IsLikeExist(userId string, resourceId string, resourceType string) (bool, error) {
	filter := bson.M{"resource_id": resourceId, "user_id": userId, "resource_type": resourceType}
	queryResult, err := ls.collectionQuerys.find(context.TODO(), filter)
	if err != nil {
		return true, err
	}
	if len(*queryResult) == 0 {
		return false, nil
	} else {
		return true, nil
	}
}

func (ls *LikeService) IsLikeExistById(likeId string) (bool, error) {
	filter := bson.M{"_id": likeId}
	queryResult, err := ls.collectionQuerys.find(context.TODO(), filter)
	if err != nil {
		return true, err
	}
	if len(*queryResult) == 0 {
		return false, nil
	} else {
		return true, nil
	}
}

func (ls *LikeService) GetLikeUserId(likeId string) (string, error) {
	filter := bson.M{"_id": likeId}
	like, err := ls.collectionQuerys.findOne(context.TODO(), filter)
	if err != nil {
		return "", err
	}
	return like.UserId, nil
}
