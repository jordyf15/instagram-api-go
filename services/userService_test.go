package services

import (
	"context"
	"errors"
	"instagram-go/models"
	"testing"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var insertOneMock func(context.Context, interface{}) (*mongo.InsertOneResult, error)
var findOneMock func(context.Context, interface{}, *models.User) error
var updateOneMock func(context.Context, interface{}, interface{}) (*mongo.UpdateResult, error)
var findMock func(context.Context, interface{}) (*[]bson.M, error)

type userCollectionQueryMock struct {
}

func (ucqm *userCollectionQueryMock) insertOne(context context.Context, document interface{}) (*mongo.InsertOneResult, error) {
	return insertOneMock(context, document)
}

func (ucqm *userCollectionQueryMock) findOne(context context.Context, filter interface{}, willBeUpdatedUser *models.User) error {
	return findOneMock(context, filter, willBeUpdatedUser)
}

func (ucqm *userCollectionQueryMock) updateOne(context context.Context, filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
	return updateOneMock(context, filter, update)
}

func (ucqm *userCollectionQueryMock) find(context context.Context, filter interface{}) (*[]bson.M, error) {
	return findMock(context, filter)
}

func newUserCollectionQueryMock() *userCollectionQueryMock {
	return &userCollectionQueryMock{}
}

func TestInsertUser(t *testing.T) {
	usercolQueryMock := newUserCollectionQueryMock()
	insertOneMock = func(context context.Context, document interface{}) (*mongo.InsertOneResult, error) {
		return nil, errors.New("insert one query error")
	}
	userService := NewUserService(usercolQueryMock)

	newUser := models.User{
		Id:              "user-" + uuid.NewString(),
		Username:        "jordyf15",
		Fullname:        "jordy ferdian",
		Password:        "jordyjordy",
		Email:           "jordyferdian@gmail.com",
		ProfilePictures: nil,
	}
	err := userService.InsertUser(newUser)
	if err == nil {
		t.Error("If insertOne to db throws an error then the insertUser should also throw an error")
	}

	insertOneMock = func(context context.Context, document interface{}) (*mongo.InsertOneResult, error) {
		return nil, nil
	}
	err = userService.InsertUser(newUser)
	if err != nil {
		t.Error("if insertOne to db does not throw error than the intertUser function should also not throw an error")
	}
}

func TestUpdateUser(t *testing.T) {
	usercolQueryMock := newUserCollectionQueryMock()
	findOneMock = func(context context.Context, filter interface{}, willBeUpdatedUser *models.User) error {
		return errors.New("error on FindOne query")
	}
	updateOneMock = func(context context.Context, filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
		return nil, nil
	}
	userService := NewUserService(usercolQueryMock)
	updatedUser := models.User{
		Id:              "user-" + uuid.NewString(),
		Username:        "jordyf15",
		Fullname:        "jordy ferdian",
		Password:        "jordyjordy",
		Email:           "jordy@gmail.com",
		ProfilePictures: nil,
	}
	err := userService.UpdateUser(updatedUser)
	if err == nil {
		t.Error("if FindOne query throws an error then UpdateUser should also throw error")
	}
	err = nil

	findOneMock = func(context context.Context, filter interface{}, willBeUpdatedUser *models.User) error {
		return nil
	}
	updateOneMock = func(context context.Context, filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
		return nil, errors.New("error on UpdateOne query")
	}
	err = userService.UpdateUser(updatedUser)
	if err == nil {
		t.Error("if UpdateOne query throws an error then UpdateUser should also throw error")
	}

	findOneMock = func(context context.Context, filter interface{}, willBeUpdatedUser *models.User) error {
		return nil
	}
	updateOneMock = func(context context.Context, filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
		return nil, nil
	}
	err = nil
	err = userService.UpdateUser(updatedUser)
	if err != nil {
		t.Error("if FindOne and UpdateOne query does not throws an error then UpdateUser should also not throw error")
	}

}

func TestCheckIfUsernameExist(t *testing.T) {
	usercolQueryMock := newUserCollectionQueryMock()
	userService := NewUserService(usercolQueryMock)
	username := "jordyf15"
	findMock = func(ctx context.Context, i interface{}) (*[]bson.M, error) {
		return &[]bson.M{}, errors.New("Find on db returns error")
	}
	if _, err := userService.CheckIfUsernameExist(username); err == nil {
		t.Error("If Find on db returns error then CheckIfUsernameExist should also return error")
	}

	findMock = func(ctx context.Context, i interface{}) (*[]bson.M, error) {
		return &[]bson.M{{"_id": "user-45d6cd8e-795a-4710-9e81-5332d57e819b", "username": "jordyf15", "full_name": "jordy ferdian", "password": "jordyjordy", "email": "jordyferdian88@gmail.com", "profile_pictures": nil}}, nil
	}
	isUsernameExist, err := userService.CheckIfUsernameExist(username)
	if err != nil {
		t.Error("If Find on db does not returns error then CheckIfUsernameExist should also not return error")
	}
	if !isUsernameExist {
		t.Error("If Find on db does return a user then CheckIfUsername exist should return true")
	}

	findMock = func(ctx context.Context, i interface{}) (*[]bson.M, error) {
		return &[]bson.M{}, nil
	}
	if isUsernameExist, _ := userService.CheckIfUsernameExist(username); isUsernameExist {
		t.Error("If Find on db return nothing then CheckIfUsernameExist should return false")
	}
}

func TestCheckIfUserExist(t *testing.T) {
	usercolQueryMock := newUserCollectionQueryMock()
	userService := NewUserService(usercolQueryMock)
	userId := "user-45d6cd8e-795a-4710-9e81-5332d57e819b"
	findMock = func(ctx context.Context, i interface{}) (*[]bson.M, error) {
		return &[]bson.M{}, errors.New("Find on db returns error")
	}
	if _, err := userService.CheckIfUserExist(userId); err == nil {
		t.Error("If Find on db returns error then CheckIfUserExist should also return error")
	}

	findMock = func(ctx context.Context, i interface{}) (*[]bson.M, error) {
		return &[]bson.M{{"_id": "user-45d6cd8e-795a-4710-9e81-5332d57e819b", "username": "jordyf15", "full_name": "jordy ferdian", "password": "jordyjordy", "email": "jordyferdian88@gmail.com", "profile_pictures": nil}}, nil
	}
	isUserExist, err := userService.CheckIfUserExist(userId)
	if err != nil {
		t.Error("If Find on db does not return error then CheckIfUserExist should also not return error")
	}
	if !isUserExist {
		t.Error("if FInd on db return a user then CheckIfUserExist should return true")
	}

	findMock = func(ctx context.Context, i interface{}) (*[]bson.M, error) {
		return &[]bson.M{}, nil
	}
	isUserExist, _ = userService.CheckIfUserExist(userId)
	if isUserExist {
		t.Error("If Find on db does not return a user then CheckIfUserExist should return false")
	}
}
