package mongodb

import (
	"context"
	"instagram-go/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongodbLikeRepository struct {
	collection *mongo.Collection
}

func NewMongodbLikeRepository(collection *mongo.Collection) *mongodbLikeRepository {
	return &mongodbLikeRepository{
		collection: collection,
	}
}

func (mlr *mongodbLikeRepository) InsertLike(like *domain.Like) error {
	newLike := bson.D{
		primitive.E{Key: "_id", Value: like.Id},
		primitive.E{Key: "user_id", Value: like.UserId},
		primitive.E{Key: "resource_id", Value: like.ResourceId},
		primitive.E{Key: "resource_type", Value: like.ResourceType},
	}
	_, err := mlr.collection.InsertOne(context.TODO(), newLike)
	if err != nil {
		return err
	}
	return nil
}

func (mlr *mongodbLikeRepository) FindLikes(filter interface{}) (*[]bson.M, error) {
	cursor, err := mlr.collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	var queryResult []bson.M
	if err = cursor.All(context.TODO(), &queryResult); err != nil {
		return nil, err
	}
	return &queryResult, nil
}

func (mlr *mongodbLikeRepository) FindOneLike(likeId string) (*domain.Like, error) {
	filter := bson.M{"_id": likeId}
	var like domain.Like
	err := mlr.collection.FindOne(context.TODO(), filter).Decode(&like)
	return &like, err
}

func (mlr *mongodbLikeRepository) DeleteLike(likeId string) error {
	filter := bson.M{"_id": likeId}
	_, err := mlr.collection.DeleteOne(context.TODO(), filter)
	return err
}
