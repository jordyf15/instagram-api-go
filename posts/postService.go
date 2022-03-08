package posts

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PostService struct {
	postCollection *mongo.Collection
	likeCollection *mongo.Collection
}

func NewPostService(postCollection *mongo.Collection, likeCollection *mongo.Collection) *PostService {
	return &PostService{
		postCollection: postCollection,
		likeCollection: likeCollection,
	}
}

func (ps *PostService) insertPost(post Post) error {
	newPost := bson.D{
		primitive.E{Key: "_id", Value: post.Id},
		primitive.E{Key: "user_id", Value: post.UserId},
		primitive.E{Key: "visual_media_urls", Value: post.VisualMediaUrls},
		primitive.E{Key: "caption", Value: post.Caption},
		primitive.E{Key: "created_date", Value: post.CreatedDate},
		primitive.E{Key: "updated_date", Value: post.UpdatedDate},
	}

	_, err := ps.postCollection.InsertOne(context.TODO(), newPost)
	if err != nil {
		return err
	}
	return nil
}

func (ps *PostService) findAllPost() ([]Post, error) {
	cursor, err := ps.postCollection.Find(context.TODO(), bson.M{})

	if err != nil {
		return nil, err
	}
	var queryResult []bson.M
	if err = cursor.All(context.TODO(), &queryResult); err != nil {
		return nil, err
	}
	var posts []Post
	for _, v := range queryResult {
		id := fmt.Sprintf("%v", v["_id"])
		userId := fmt.Sprintf("%v", v["user_id"])
		var visualMediaUrls []string

		if visualMediaUrlsPrimitive, ok := v["visual_media_urls"].(primitive.A); ok {
			visualMediaUrlsInterface := []interface{}(visualMediaUrlsPrimitive)
			visualMediaUrls = make([]string, len(visualMediaUrlsInterface))
			for i, url := range visualMediaUrlsInterface {
				visualMediaUrls[i] = url.(string)
			}
		}

		caption := fmt.Sprintf("%v", v["caption"])
		createdDate := v["created_date"].(primitive.DateTime).Time()
		updatedDate := v["updated_date"].(primitive.DateTime).Time()
		likeCount, err := ps.getPostLikeCount(id)
		if err != nil {
			return nil, err
		}
		post := Post{id, userId, visualMediaUrls, caption, likeCount, createdDate, updatedDate}
		posts = append(posts, post)
	}
	return posts, nil
}

func (ps *PostService) getPostLikeCount(postId string) (int, error) {
	filter := bson.M{"resource_id": postId, "resource_type": "post"}
	cursor, err := ps.likeCollection.Find(context.TODO(), filter)
	if err != nil {
		return 0, err
	}
	var likes []bson.M
	if err = cursor.All(context.TODO(), &likes); err != nil {
		return 0, err
	}
	return len(likes), nil
}

func (ps *PostService) updatePost(updatedPostId string, newCaption string) error {
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
	_, err := ps.postCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (ps *PostService) findPost(postId string) (string, error) {
	var post Post
	filter := bson.M{"_id": postId}
	err := ps.postCollection.FindOne(context.TODO(), filter).Decode(&post)
	if err != nil {
		return "", err
	}
	return post.UserId, nil
}

func (ps *PostService) deletePost(postId string) error {
	filter := bson.M{"_id": postId}
	_, err := ps.postCollection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}
	return nil
}