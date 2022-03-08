package likes

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type LikeService struct {
	collection *mongo.Collection
}

func NewLikeService(collection *mongo.Collection) *LikeService {
	return &LikeService{collection: collection}
}

func (ls *LikeService) insertLike(like Like) error {
	newLike := bson.D{
		primitive.E{Key: "_id", Value: like.Id},
		primitive.E{Key: "user_id", Value: like.UserId},
		primitive.E{Key: "resource_id", Value: like.ResourceId},
		primitive.E{Key: "resource_type", Value: like.ResourceType},
	}
	_, err := ls.collection.InsertOne(context.TODO(), newLike)
	if err != nil {
		return err
	}
	return nil
}

func (ls *LikeService) deleteLike(likeId string) error {
	filter := bson.M{"_id": likeId}
	_, err := ls.collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}
	return nil
}

func (ls *LikeService) isLikeExist(userId string, resourceId string, resourceType string) (bool, error) {
	filter := bson.M{"resource_id": resourceId, "user_id": userId, "resource_type": resourceType}
	cursor, err := ls.collection.Find(context.TODO(), filter, options.Find().SetLimit(1))
	if err != nil {
		return true, err
	}
	var queryResult []bson.M
	if err = cursor.All(context.TODO(), &queryResult); err != nil {
		return true, err
	}
	if len(queryResult) == 0 {
		return false, nil
	} else {
		return true, nil
	}
}

func (ls *LikeService) getLikeUserId(likeId string) (string, error) {
	var like Like
	filter := bson.M{"_id": likeId}
	err := ls.collection.FindOne(context.TODO(), filter).Decode(&like)
	if err != nil {
		return "", err
	}
	return like.UserId, nil
}
