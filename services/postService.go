package services

import (
	"context"
	"fmt"
	"instagram-go/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PostService struct {
	collectionQuerys postCollectionQueryable
}

type postCollectionQueryable interface {
	postInsertOne(context.Context, interface{}) (*mongo.InsertOneResult, error)
	postUpdateOne(context.Context, interface{}, interface{}) (*mongo.UpdateResult, error)
	postFindOne(context.Context, interface{}, *models.Post) error
	postDeleteOne(context.Context, interface{}) (*mongo.DeleteResult, error)
	postFind(context.Context, interface{}) (*mongo.Cursor, error)
	postIterateCursor(*mongo.Cursor, context.Context, *[]bson.M) error
	likeFind(context.Context, interface{}) (*mongo.Cursor, error)
	likeIterateCursor(*mongo.Cursor, context.Context, *[]bson.M) error
}

type postCollectionQuery struct {
	postCollection *mongo.Collection
	likeCollection *mongo.Collection
}

func NewPostCollectionQuery(postCollection *mongo.Collection, likeCollection *mongo.Collection) *postCollectionQuery {
	return &postCollectionQuery{
		postCollection: postCollection,
		likeCollection: likeCollection,
	}
}

func (pcq *postCollectionQuery) postInsertOne(context context.Context, document interface{}) (*mongo.InsertOneResult, error) {
	return pcq.postCollection.InsertOne(context, document)
}

func (pcq *postCollectionQuery) postUpdateOne(context context.Context, filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
	return pcq.postCollection.UpdateOne(context, filter, update)
}

func (pcq *postCollectionQuery) postFindOne(context context.Context, filter interface{}, post *models.Post) error {
	return pcq.postCollection.FindOne(context, filter).Decode(post)
}

func (pcq *postCollectionQuery) postDeleteOne(context context.Context, filter interface{}) (*mongo.DeleteResult, error) {
	return pcq.postCollection.DeleteOne(context, filter)
}

func (pcq *postCollectionQuery) postFind(context context.Context, filter interface{}) (*mongo.Cursor, error) {
	return pcq.postCollection.Find(context, filter)
}

func (pcq *postCollectionQuery) postIterateCursor(cursor *mongo.Cursor, context context.Context, queryResult *[]bson.M) error {
	return cursor.All(context, queryResult)
}

func (pcq *postCollectionQuery) likeFind(context context.Context, filter interface{}) (*mongo.Cursor, error) {
	return pcq.likeCollection.Find(context, filter)
}

func (pcq *postCollectionQuery) likeIterateCursor(cursor *mongo.Cursor, context context.Context, queryResult *[]bson.M) error {
	return cursor.All(context, queryResult)
}

func NewPostService(postCollectionQuery postCollectionQueryable) *PostService {
	return &PostService{
		collectionQuerys: postCollectionQuery,
	}
}

func (ps *PostService) InsertPost(post models.Post) error {
	newPost := bson.D{
		primitive.E{Key: "_id", Value: post.Id},
		primitive.E{Key: "user_id", Value: post.UserId},
		primitive.E{Key: "visual_media_urls", Value: post.VisualMediaUrls},
		primitive.E{Key: "caption", Value: post.Caption},
		primitive.E{Key: "created_date", Value: post.CreatedDate},
		primitive.E{Key: "updated_date", Value: post.UpdatedDate},
	}

	// _, err := ps.postCollection.InsertOne(context.TODO(), newPost)
	_, err := ps.collectionQuerys.postInsertOne(context.TODO(), newPost)
	if err != nil {
		return err
	}
	return nil
}

func (ps *PostService) FindAllPost() ([]models.Post, error) {
	// cursor, err := ps.postCollection.Find(context.TODO(), bson.M{})
	cursor, err := ps.collectionQuerys.postFind(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}
	var queryResult []bson.M
	// if err = cursor.All(context.TODO(), &queryResult);
	if err = ps.collectionQuerys.postIterateCursor(cursor, context.TODO(), &queryResult); err != nil {
		return nil, err
	}
	var posts []models.Post
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
		post := models.NewPost(id, userId, visualMediaUrls, caption, likeCount, createdDate, updatedDate)
		posts = append(posts, *post)
	}
	return posts, nil
}

func (ps *PostService) getPostLikeCount(postId string) (int, error) {
	filter := bson.M{"resource_id": postId, "resource_type": "post"}
	// cursor, err := ps.likeCollection.Find(context.TODO(), filter)
	cursor, err := ps.collectionQuerys.likeFind(context.TODO(), filter)
	if err != nil {
		return 0, err
	}
	var likes []bson.M
	// if err = cursor.All(context.TODO(), &likes);
	if err = ps.collectionQuerys.likeIterateCursor(cursor, context.TODO(), &likes); err != nil {
		return 0, err
	}
	return len(likes), nil
}

func (ps *PostService) UpdatePost(updatedPostId string, newCaption string) error {
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
	// _, err := ps.postCollection.UpdateOne(context.TODO(), filter, update)
	_, err := ps.collectionQuerys.postUpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (ps *PostService) FindPost(postId string) (string, error) {
	var post models.Post
	filter := bson.M{"_id": postId}
	// err := ps.postCollection.FindOne(context.TODO(), filter).Decode(&post)
	err := ps.collectionQuerys.postFindOne(context.TODO(), filter, &post)
	if err != nil {
		return "", err
	}
	return post.UserId, nil
}

func (ps *PostService) DeletePost(postId string) error {
	filter := bson.M{"_id": postId}
	// _, err := ps.postCollection.DeleteOne(context.TODO(), filter)
	_, err := ps.collectionQuerys.postDeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}
	return nil
}
