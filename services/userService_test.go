package services

import (
	"context"
	"errors"
	"instagram-go/models"
	"testing"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
)

var insertOneMock func(context.Context, interface{}) (*mongo.InsertOneResult, error)
var findOneMock func(context.Context, interface{}, *models.User) error
var updateOneMock func(context.Context, interface{}, interface{}) (*mongo.UpdateResult, error)

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
