package services

import (
	"context"
	"errors"
	"instagram-go/models"
	"testing"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var postInsertOneMock func(context.Context, interface{}) (*mongo.InsertOneResult, error)
var postUpdateOneMock func(context.Context, interface{}, interface{}) (*mongo.UpdateResult, error)
var postFindOneMock func(context.Context, interface{}) (*models.Post, error)
var postDeleteOneMock func(context.Context, interface{}) (*mongo.DeleteResult, error)
var postFindMock func(context.Context, interface{}) (*[]bson.M, error)
var likeFindMock func(context.Context, interface{}) (*[]bson.M, error)

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

func (pcqm *postCollectionQueryMock) postFindOne(context context.Context, filter interface{}) (*models.Post, error) {
	return postFindOneMock(context, filter)
}

func (pcqm *postCollectionQueryMock) postDeleteOne(context context.Context, filter interface{}) (*mongo.DeleteResult, error) {
	return postDeleteOneMock(context, filter)
}

func (pcqm *postCollectionQueryMock) postFind(context context.Context, filter interface{}) (*[]bson.M, error) {
	return postFindMock(context, filter)
}

func (pcqm *postCollectionQueryMock) likeFind(context context.Context, filter interface{}) (*[]bson.M, error) {
	return likeFindMock(context, filter)
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

	likeFindMock = func(ctx context.Context, i interface{}) (*[]bson.M, error) {
		return &[]bson.M{}, nil
	}
	postFindMock = func(ctx context.Context, i interface{}) (*[]bson.M, error) {
		return &[]bson.M{}, errors.New("Find to db returns an error")
	}
	_, err := postService.FindAllPost()
	if err == nil {
		t.Error("If Find on db returns an error than FindAllPost should also return an error")
	}

	err = nil
	postFindMock = func(ctx context.Context, i interface{}) (*[]bson.M, error) {
		return &[]bson.M{
			{"_id": "post-" + uuid.NewString(),
				"user_id":           "user-" + uuid.NewString(),
				"visual_media_urls": []string{},
				"caption":           "a new caption",
				"like_count":        0,
				"created_date":      primitive.NewDateTimeFromTime(time.Now()),
				"updated_date":      primitive.NewDateTimeFromTime(time.Now()),
			}, {"_id": "post-" + uuid.NewString(),
				"user_id":           "user-" + uuid.NewString(),
				"visual_media_urls": []string{},
				"caption":           "a new caption",
				"like_count":        0,
				"created_date":      primitive.NewDateTimeFromTime(time.Now()),
				"updated_date":      primitive.NewDateTimeFromTime(time.Now()),
			}, {"_id": "post-" + uuid.NewString(),
				"user_id":           "user-" + uuid.NewString(),
				"visual_media_urls": []string{},
				"caption":           "a new caption",
				"like_count":        0,
				"created_date":      primitive.NewDateTimeFromTime(time.Now()),
				"updated_date":      primitive.NewDateTimeFromTime(time.Now()),
			},
		}, nil
	}
	allPost, err := postService.FindAllPost()
	if err != nil {
		t.Error("If Iterate cursor returns an error than FindAllPost should also returns an error")
	}
	if len(allPost) != 3 {
		t.Error("FindAllPost shall return also founded Post from db")
	}
}

func TestGetPostLikeCount(t *testing.T) {
	postColQueryMock := newPostCollectionQueryMock()
	postService := NewPostService(postColQueryMock)
	findPostId := "post-8b652b94-ff70-444c-a044-266ec7779c45"

	likeFindMock = func(ctx context.Context, i interface{}) (*[]bson.M, error) {
		return nil, errors.New("Find likes on db returns an error")
	}
	_, err := postService.getPostLikeCount(findPostId)
	if err == nil {
		t.Error("Find likes on db returns an error")
	}
	err = nil
	likeFindMock = func(ctx context.Context, i interface{}) (*[]bson.M, error) {
		return &[]bson.M{
			{"_id": "like-5c533843-2138-4b8b-87d2-170d34626975",
				"user_id":       "user-1cab874c-d558-45a2-aaef-6487a415c261",
				"resource_id":   "post-8b652b94-ff70-444c-a044-266ec7779c45",
				"resource_type": "post"},
		}, nil
	}
	likeCount, err := postService.getPostLikeCount(findPostId)
	if err != nil {
		t.Error("If Find likes on db does not return error than getPostLikeCount should also not return error")
	}
	if likeCount != 1 {
		t.Error("If likes is found than getPostLikeCount should return the correct quantity of the found likes")
	}

	likeFindMock = func(ctx context.Context, i interface{}) (*[]bson.M, error) {
		return &[]bson.M{}, nil
	}
	likeCount, _ = postService.getPostLikeCount(findPostId)
	if likeCount != 0 {
		t.Error("If likes is not found than getPostLikeCount should return 0")
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

func TestGetPostUserId(t *testing.T) {
	postColQueryMock := newPostCollectionQueryMock()
	postService := NewPostService(postColQueryMock)
	findPostId := "post-8b652b94-ff70-444c-a044-266ec7779c45"
	foundUserId := "user-8b652b94-ff70-444c-a044-266ec7779c45"

	postFindOneMock = func(ctx context.Context, i interface{}) (*models.Post, error) {
		return nil, errors.New("FindOne to db returns an error")
	}
	_, err := postService.GetPostUserId(findPostId)
	if err == nil {
		t.Error("if FindOne to db returns an error than GetPostUserId should also return an error")
	}

	postFindOneMock = func(ctx context.Context, i interface{}) (*models.Post, error) {
		return models.NewPost("post-8b652b94-ff70-444c-a044-266ec7779c45", foundUserId, []string{}, "a caption", 0, time.Now(), time.Now()), nil
	}
	fetchedUserId, err := postService.GetPostUserId(findPostId)
	if err != nil {
		t.Error("If FindOne to db does not returns an error than GetPostUserId should also not return an error")
	}
	if foundUserId != fetchedUserId {
		t.Error("the GetPostUserId should return the proper userid of the post")
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

func TestCheckIfPostExist(t *testing.T) {
	postColQueryMock := newPostCollectionQueryMock()
	postService := NewPostService(postColQueryMock)
	searchedPostId := "post-8b652b94-ff70-444c-a044-266ec7779c45"

	postFindMock = func(ctx context.Context, i interface{}) (*[]bson.M, error) {
		return &[]bson.M{}, errors.New("Find on db returns error")
	}
	if _, err := postService.CheckIfPostExist(searchedPostId); err == nil {
		t.Error("If Find on db returns error than CheckIfPostExist should also return error")
	}

	postFindMock = func(ctx context.Context, i interface{}) (*[]bson.M, error) {
		return &[]bson.M{}, nil
	}
	isExist, err := postService.CheckIfPostExist(searchedPostId)
	if err != nil {
		t.Error("If Find on db does not return error than CheckIfPostExist should also not return error")
	}
	if isExist {
		t.Error("If Find on db does not return any post than CheckIfPostExist should return false")
	}

	postFindMock = func(ctx context.Context, i interface{}) (*[]bson.M, error) {
		return &[]bson.M{{"_id": "post-" + uuid.NewString(),
			"user_id":           "user-" + uuid.NewString(),
			"visual_media_urls": []string{},
			"caption":           "a new caption",
			"like_count":        0,
			"created_date":      primitive.NewDateTimeFromTime(time.Now()),
			"updated_date":      primitive.NewDateTimeFromTime(time.Now()),
		}}, nil
	}
	isExist, _ = postService.CheckIfPostExist(searchedPostId)
	if !isExist {
		t.Error("If Find on db does return a post than CheckIfPostExist should return true")
	}
}
