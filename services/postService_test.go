package services

import (
	"context"
	"errors"
	"instagram-go/models"
	"testing"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var postInsertOneMock func(context.Context, interface{}) (*mongo.InsertOneResult, error)
var postUpdateOneMock func(context.Context, interface{}, interface{}) (*mongo.UpdateResult, error)
var postFindOneMock func(context.Context, interface{}, *models.Post) error
var postDeleteOneMock func(context.Context, interface{}) (*mongo.DeleteResult, error)
var postFindMock func(context.Context, interface{}) (*mongo.Cursor, error)
var postIterateCursorMock func(*mongo.Cursor, context.Context, *[]bson.M) error
var likeFindMock func(context.Context, interface{}) (*mongo.Cursor, error)
var likeIterateCursorMock func(*mongo.Cursor, context.Context, *[]bson.M) error

type postCollectionQueryMock struct {
}

func newPostCollectionQueryMock() *postCollectionQueryMock {
	return &postCollectionQueryMock{}
}

func (pcqm *postCollectionQueryMock) postInsertOne(context context.Context, document interface{}) (*mongo.InsertOneResult, error) {
	return postInsertOneMock(context, document)
}

func (pcqm *postCollectionQueryMock) postUpdateOne(context context.Context, filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
	return postUpdateOneMock(context, filter, update)
}

func (pcqm *postCollectionQueryMock) postFindOne(context context.Context, filter interface{}, post *models.Post) error {
	return postFindOneMock(context, filter, post)
}

func (pcqm *postCollectionQueryMock) postDeleteOne(context context.Context, filter interface{}) (*mongo.DeleteResult, error) {
	return postDeleteOneMock(context, filter)
}

func (pcqm *postCollectionQueryMock) postFind(context context.Context, filter interface{}) (*mongo.Cursor, error) {
	return postFindMock(context, filter)
}

func (pcqm *postCollectionQueryMock) postIterateCursor(cursor *mongo.Cursor, context context.Context, queryResult *[]bson.M) error {
	return postIterateCursorMock(cursor, context, queryResult)
}

func (pcqm *postCollectionQueryMock) likeFind(context context.Context, filter interface{}) (*mongo.Cursor, error) {
	return likeFindMock(context, filter)
}

func (pcqm *postCollectionQueryMock) likeIterateCursor(cursor *mongo.Cursor, context context.Context, queryResult *[]bson.M) error {
	return likeIterateCursorMock(cursor, context, queryResult)
}

func TestInsertPost(t *testing.T) {
	postColQueryMock := newPostCollectionQueryMock()
	postService := NewPostService(postColQueryMock)
	newPost := models.Post{
		Id:              "post-" + uuid.NewString(),
		UserId:          "user-" + uuid.NewString(),
		VisualMediaUrls: []string{},
		Caption:         "a new post",
		LikeCount:       0,
		CreatedDate:     time.Now(),
		UpdatedDate:     time.Now(),
	}

	postInsertOneMock = func(context context.Context, document interface{}) (*mongo.InsertOneResult, error) {
		return nil, errors.New("insert one throws error")
	}
	err := postService.InsertPost(newPost)
	if err == nil {
		t.Error("If InsertOne to db returns an error then InsertPost should return an error")
	}

	err = nil
	postInsertOneMock = func(context context.Context, document interface{}) (*mongo.InsertOneResult, error) {
		return nil, nil
	}
	err = postService.InsertPost(newPost)
	if err != nil {
		t.Error("If InsertOne to db does not returns an error then InsertPost should also not return an error")
	}
}

func TestFindAllPost(t *testing.T) {
	postColQueryMock := newPostCollectionQueryMock()
	postService := NewPostService(postColQueryMock)

	postFindMock = func(context context.Context, filter interface{}) (*mongo.Cursor, error) {
		return nil, errors.New("Find on db returns an error")
	}
	postIterateCursorMock = func(c *mongo.Cursor, ctx context.Context, m *[]bson.M) error {
		return nil
	}
	_, err := postService.FindAllPost()
	if err == nil {
		t.Error("If Find on db returns an error than FindAllPost should also return an error")
	}

	postFindMock = func(context context.Context, filter interface{}) (*mongo.Cursor, error) {
		return nil, nil
	}
	postIterateCursorMock = func(c *mongo.Cursor, ctx context.Context, m *[]bson.M) error {
		return errors.New("Iterate cursor returns an error")
	}
	_, err = postService.FindAllPost()
	if err == nil {
		t.Error("If Iterate cursor returns an error than FindAllPost should also returns an error")
	}

	postFindMock = func(context context.Context, filter interface{}) (*mongo.Cursor, error) {
		return nil, nil
	}
	postIterateCursorMock = func(c *mongo.Cursor, ctx context.Context, m *[]bson.M) error {
		return nil
	}
	_, err = postService.FindAllPost()
	if err != nil {
		t.Error("If Find on db and Iterate cursor does not return an error than FindAllPost should also not returns an error")
	}
}

func TestGetPostLikeCount(t *testing.T) {
	postColQueryMock := newPostCollectionQueryMock()
	postService := NewPostService(postColQueryMock)
	findPostId := "post-8b652b94-ff70-444c-a044-266ec7779c45"

	likeFindMock = func(context context.Context, filter interface{}) (*mongo.Cursor, error) {
		return nil, errors.New("Find likes on db returns an error")
	}
	likeIterateCursorMock = func(cursor *mongo.Cursor, context context.Context, queryResults *[]bson.M) error {
		return nil
	}
	_, err := postService.getPostLikeCount(findPostId)
	if err == nil {
		t.Error("If Find likes on db returns an error than getPostLikeCount should also return an error")
	}

	err = nil
	likeFindMock = func(context context.Context, filter interface{}) (*mongo.Cursor, error) {
		return nil, nil
	}
	likeIterateCursorMock = func(cursor *mongo.Cursor, context context.Context, queryResults *[]bson.M) error {
		return errors.New("Iterate cursors returns an error")
	}
	_, err = postService.getPostLikeCount(findPostId)
	if err == nil {
		t.Error("If Iterate cursor returns an error than getPostLikeCount should also return an error")
	}

	err = nil
	likeFindMock = func(context context.Context, filter interface{}) (*mongo.Cursor, error) {
		return nil, nil
	}
	likeIterateCursorMock = func(cursor *mongo.Cursor, context context.Context, queryResults *[]bson.M) error {
		return nil
	}
	_, err = postService.getPostLikeCount(findPostId)
	if err != nil {
		t.Error("If Find Like on db and Iterate cursor does not returns an error than getPostLikeCount should also not return an error")
	}
}

func TestUpdatePost(t *testing.T) {
	postColQueryMock := newPostCollectionQueryMock()
	postService := NewPostService(postColQueryMock)
	updatedPostId := "post-8b652b94-ff70-444c-a044-266ec7779c45"
	newCaption := "an updated post caption"

	postUpdateOneMock = func(context context.Context, filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
		return nil, errors.New("UpdateOne returns an error")
	}
	err := postService.UpdatePost(updatedPostId, newCaption)
	if err == nil {
		t.Error("If UpdateOne to db returns an error then UpdatePost should also return an error")
	}

	postUpdateOneMock = func(context context.Context, filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
		return nil, nil
	}
	err = nil
	if err != nil {
		t.Error("If UpdateOne to db does not return an error then UpdatePost should also not return an error")
	}
}

func TestFindPost(t *testing.T) {
	postColQueryMock := newPostCollectionQueryMock()
	postService := NewPostService(postColQueryMock)
	findPostId := "post-8b652b94-ff70-444c-a044-266ec7779c45"

	postFindOneMock = func(context context.Context, filter interface{}, post *models.Post) error {
		return errors.New("FindOne returns an error")
	}
	_, err := postService.FindPost(findPostId)
	if err == nil {
		t.Error("If FindOne to db returns an error than FindPost should also return an error")
	}

	err = nil
	postFindOneMock = func(context context.Context, filter interface{}, post *models.Post) error {
		return nil
	}
	_, err = postService.FindPost(findPostId)
	if err != nil {
		t.Error("If FindOne to db does not returns an error than FindPost should also not return an error")
	}
}

func TestDeletePost(t *testing.T) {
	postColQueryMock := newPostCollectionQueryMock()
	postService := NewPostService(postColQueryMock)
	deletedPostId := "post-8b652b94-ff70-444c-a044-266ec7779c45"

	postDeleteOneMock = func(context context.Context, filter interface{}) (*mongo.DeleteResult, error) {
		return nil, errors.New("DeleteOne query return errors")
	}
	err := postService.DeletePost(deletedPostId)
	if err == nil {
		t.Error("If DeleteOne to db returns an error then DeletePost should also returns an error")
	}

	err = nil
	postDeleteOneMock = func(context context.Context, filter interface{}) (*mongo.DeleteResult, error) {
		return nil, nil
	}
	err = postService.DeletePost(deletedPostId)
	if err != nil {
		t.Error("If DeleteOne to db does not return an error then DeletePost should also not return an error")
	}
}
