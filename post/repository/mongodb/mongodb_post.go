package mongodb

import (
	"context"
	"instagram-go/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongodbPostRepository struct {
	collection *mongo.Collection
}

func NewMongodbPostRepository(collection *mongo.Collection) *mongodbPostRepository {
	return &mongodbPostRepository{
		collection: collection,
	}
}

func (pr *mongodbPostRepository) InsertPost(post *domain.Post) error {
	newPost := bson.D{
		primitive.E{Key: "_id", Value: post.Id},
		primitive.E{Key: "user_id", Value: post.UserId},
		primitive.E{Key: "visual_media_urls", Value: post.VisualMediaUrls},
		primitive.E{Key: "caption", Value: post.Caption},
		primitive.E{Key: "created_date", Value: post.CreatedDate},
		primitive.E{Key: "updated_date", Value: post.UpdatedDate},
	}
	_, err := pr.collection.InsertOne(context.TODO(), newPost)
	if err != nil {
		return err
	}
	return nil
}

func (pr *mongodbPostRepository) FindPosts(filter interface{}) (*[]bson.M, error) {
	cursor, err := pr.collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	var queryResult []bson.M
	if err = cursor.All(context.TODO(), &queryResult); err != nil {
		return nil, err
	}
	return &queryResult, nil
}

func (pr *mongodbPostRepository) FindOnePost(searchedPostId string) (*domain.Post, error) {
	var post domain.Post
	filter := bson.M{"_id": searchedPostId}
	err := pr.collection.FindOne(context.TODO(), filter).Decode(&post)
	return &post, err
}

func (pr *mongodbPostRepository) UpdatePost(updatedPostId string, newCaption string) error {
	filter := bson.M{"_id": updatedPostId}
	update := bson.D{primitive.E{
		Key: "$set",
		Value: bson.D{primitive.E{
			Key:   "caption",
			Value: newCaption,
		},
		},
	},
	}
	_, err := pr.collection.UpdateOne(context.TODO(), filter, update)
	return err
}

func (pr *mongodbPostRepository) DeletePost(deletedPostId string) error {
	filter := bson.M{"_id": deletedPostId}
	_, err := pr.collection.DeleteOne(context.TODO(), filter)
	return err
}
