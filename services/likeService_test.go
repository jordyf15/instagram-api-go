package services

import (
	"context"
	"errors"
	"instagram-go/models"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var insertOneLikeMock func(context.Context, interface{}) (*mongo.InsertOneResult, error)
var deleteOneLikeMock func(context.Context, interface{}) (*mongo.DeleteResult, error)
var findOneLikeMock func(context.Context, interface{}) (*models.Like, error)
var findLikeMock func(context.Context, interface{}) (*[]bson.M, error)

type likeCollectionQueryMock struct {
}

func newLikeCollectionQueryMock() *likeCollectionQueryMock {
	return &likeCollectionQueryMock{}
}

func (lcqm *likeCollectionQueryMock) insertOne(context context.Context, document interface{}) (*mongo.InsertOneResult, error) {
	return insertOneLikeMock(context, document)
}

func (lcqm *likeCollectionQueryMock) deleteOne(context context.Context, filter interface{}) (*mongo.DeleteResult, error) {
	return deleteOneLikeMock(context, filter)
}

func (lcqm *likeCollectionQueryMock) findOne(context context.Context, filter interface{}) (*models.Like, error) {
	return findOneLikeMock(context, filter)
}

func (lcqm *likeCollectionQueryMock) find(context context.Context, filter interface{}) (*[]bson.M, error) {
	return findLikeMock(context, filter)
}

func TestInsertLike(t *testing.T) {
	likeColQuerymock := newLikeCollectionQueryMock()
	likeService := NewLikeService(likeColQuerymock)
	newLike := models.Like{
		Id:           "like-5c533843-2138-4b8b-87d2-170d34626975",
		UserId:       "user-1cab874c-d558-45a2-aaef-6487a415c261",
		ResourceId:   "comment-bfe6f419-0553-46dd-a58d-9dd186ebd2e4",
		ResourceType: "comment",
	}

	insertOneLikeMock = func(ctx context.Context, i interface{}) (*mongo.InsertOneResult, error) {
		return nil, errors.New("insertOne to db return error")
	}
	err := likeService.InsertLike(newLike)
	if err == nil {
		t.Error("If insertOne to db return error than InsertLike should also return error")
	}

	err = nil
	insertOneLikeMock = func(ctx context.Context, i interface{}) (*mongo.InsertOneResult, error) {
		return nil, nil
	}
	err = likeService.InsertLike(newLike)
	if err != nil {
		t.Error("If insertOne to db does not return error than InsertLike should also not return error")
	}
}

func TestDeleteLike(t *testing.T) {
	likeColQuerymock := newLikeCollectionQueryMock()
	likeService := NewLikeService(likeColQuerymock)
	deleteLikeId := "like-5c533843-2138-4b8b-87d2-170d34626975"

	deleteOneLikeMock = func(ctx context.Context, i interface{}) (*mongo.DeleteResult, error) {
		return nil, errors.New("deleteOne to db return error")
	}
	err := likeService.DeleteLike(deleteLikeId)
	if err == nil {
		t.Error("If DeleteOne to db return error than DeleteLike should also return error")
	}

	err = nil
	deleteOneLikeMock = func(ctx context.Context, i interface{}) (*mongo.DeleteResult, error) {
		return nil, nil
	}
	err = likeService.DeleteLike(deleteLikeId)
	if err != nil {
		t.Error("If DeleteOne to db does not return error than DeleteLike should also not return error")
	}
}

func TestIsLikeExist(t *testing.T) {
	likeColQuerymock := newLikeCollectionQueryMock()
	likeService := NewLikeService(likeColQuerymock)
	userId := "user-1cab874c-d558-45a2-aaef-6487a415c261"
	resourceId := "comment-bfe6f419-0553-46dd-a58d-9dd186ebd2e4"
	resourceType := "comment"

	findLikeMock = func(ctx context.Context, i interface{}) (*[]bson.M, error) {
		return &[]bson.M{}, errors.New("find to db return error")
	}
	_, err := likeService.IsLikeExist(userId, resourceId, resourceType)
	if err == nil {
		t.Error("If find to db return error than IsLikeExist should also return error")
	}

	err = nil
	findLikeMock = func(ctx context.Context, i interface{}) (*[]bson.M, error) {
		return &[]bson.M{}, nil
	}
	_, err = likeService.IsLikeExist(userId, resourceId, resourceType)
	if err != nil {
		t.Error("if find to db does not return error than is Like exist should also not return error")
	}

	findLikeMock = func(ctx context.Context, i interface{}) (*[]bson.M, error) {
		return &[]bson.M{
			{"_id": "like-5c533843-2138-4b8b-87d2-170d34626975",
				"user_id":       "user-1cab874c-d558-45a2-aaef-6487a415c261",
				"resource_id":   "comment-bfe6f419-0553-46dd-a58d-9dd186ebd2e4",
				"resource_type": "comment"},
		}, nil
	}
	isExist, _ := likeService.IsLikeExist(userId, resourceId, resourceType)
	if !isExist {
		t.Error("If a like is found than IsLikeExist should return true")
	}

	findLikeMock = func(ctx context.Context, i interface{}) (*[]bson.M, error) {
		return &[]bson.M{}, nil
	}
	isExist, _ = likeService.IsLikeExist(userId, resourceId, resourceType)
	if isExist {
		t.Error("If a like is not found than IsLikeExist should return false")
	}
}

func TestGetLikeUserId(t *testing.T) {
	likeColQuerymock := newLikeCollectionQueryMock()
	likeService := NewLikeService(likeColQuerymock)
	likeId := "like-5c533843-2138-4b8b-87d2-170d34626975"
	returnedLike := models.NewLike("like-5c533843-2138-4b8b-87d2-170d34626975", "user-1cab874c-d558-45a2-aaef-6487a415c261",
		"comment-bfe6f419-0553-46dd-a58d-9dd186ebd2e4", "comment")

	findOneLikeMock = func(ctx context.Context, i interface{}) (*models.Like, error) {
		return nil, errors.New("FindOne to db return error")
	}
	_, err := likeService.GetLikeUserId(likeId)
	if err == nil {
		t.Error("if FindOne to db return error than GetLikeUserId should also return error")
	}

	err = nil
	findOneLikeMock = func(ctx context.Context, i interface{}) (*models.Like, error) {
		return returnedLike, nil
	}
	_, err = likeService.GetLikeUserId(likeId)
	if err != nil {
		t.Error("if FindOne to db does not return error than GetLikeUserId should also not return error")
	}

	findOneLikeMock = func(ctx context.Context, i interface{}) (*models.Like, error) {
		return returnedLike, nil
	}
	if returnedLike.Id != likeId {
		t.Error("GetLikeUserId should return the id of the found like")
	}
}
