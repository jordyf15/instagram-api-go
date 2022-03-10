package services

import (
	"context"
	"instagram-go/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	collectionQuerys userCollectionQueryable
}

type userCollectionQueryable interface {
	insertOne(context.Context, interface{}) (*mongo.InsertOneResult, error)
	findOne(context.Context, interface{}, *models.User) error
	updateOne(context.Context, interface{}, interface{}) (*mongo.UpdateResult, error)
}

type userCollectionQuery struct {
	collection *mongo.Collection
}

func NewUserCollectionQuery(collection *mongo.Collection) *userCollectionQuery {
	return &userCollectionQuery{
		collection: collection,
	}
}

func (ucq *userCollectionQuery) insertOne(context context.Context, document interface{}) (*mongo.InsertOneResult, error) {
	return ucq.collection.InsertOne(context, document)
}

func (ucq *userCollectionQuery) findOne(context context.Context, filter interface{}, willBeUpdatedUser *models.User) error {
	return ucq.collection.FindOne(context, filter).Decode(willBeUpdatedUser)
}

func (ucq *userCollectionQuery) updateOne(context context.Context, filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
	return ucq.collection.UpdateOne(context, filter, update)
}

func NewUserService(userCollectionQuery userCollectionQueryable) *UserService {
	return &UserService{
		collectionQuerys: userCollectionQuery,
	}
}

func (us *UserService) InsertUser(user models.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		return err
	}
	newUser := bson.D{
		primitive.E{Key: "_id", Value: user.Id},
		primitive.E{Key: "username", Value: user.Username},
		primitive.E{Key: "full_name", Value: user.Fullname},
		primitive.E{Key: "password", Value: string(hashedPassword[:])},
		primitive.E{Key: "email", Value: user.Email},
		primitive.E{Key: "profile_pictures", Value: nil},
	}

	result, err := us.collectionQuerys.insertOne(context.TODO(), newUser)
	if err != nil {
		return err
	}
	if result != nil {
		return nil
	}
	return nil
}

func (us *UserService) UpdateUser(newUserData models.User) error {
	var willBeUpdatedUser models.User
	filter := bson.M{"_id": newUserData.Id}
	err := us.collectionQuerys.findOne(context.TODO(), filter, &willBeUpdatedUser)
	if err != nil {
		return err
	}
	if newUserData.Email != "" {
		willBeUpdatedUser.Email = newUserData.Email
	}
	if newUserData.Username != "" {
		willBeUpdatedUser.Username = newUserData.Username
	}
	if newUserData.Fullname != "" {
		willBeUpdatedUser.Fullname = newUserData.Fullname
	}
	if newUserData.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUserData.Password), 10)
		if err != nil {
			return err
		}
		willBeUpdatedUser.Password = string(hashedPassword)
	}
	if newUserData.ProfilePictures != nil {
		willBeUpdatedUser.ProfilePictures = newUserData.ProfilePictures
	}
	update := bson.D{primitive.E{Key: "$set", Value: bson.D{
		primitive.E{Key: "username", Value: willBeUpdatedUser.Username},
		primitive.E{Key: "full_name", Value: willBeUpdatedUser.Fullname},
		primitive.E{Key: "password", Value: willBeUpdatedUser.Password},
		primitive.E{Key: "email", Value: willBeUpdatedUser.Email},
		primitive.E{Key: "profile_pictures", Value: willBeUpdatedUser.ProfilePictures},
	},
	},
	}
	_, err = us.collectionQuerys.updateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}
	return nil
}
