package services

import (
	"context"
	"instagram-go/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	collectionQuerys userCollectionQueryable
}

type IUserService interface {
	InsertUser(models.User) error
	UpdateUser(models.User) error
	CheckIfUsernameExist(string) (bool, error)
	CheckIfUserExist(string) (bool, error)
}

type userCollectionQueryable interface {
	insertOne(context.Context, interface{}) (*mongo.InsertOneResult, error)
	findOne(context.Context, interface{}, *models.User) error
	updateOne(context.Context, interface{}, interface{}) (*mongo.UpdateResult, error)
	find(context.Context, interface{}) (*[]bson.M, error)
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

func (ucq *userCollectionQuery) find(context context.Context, filter interface{}) (*[]bson.M, error) {
	cursor, err := ucq.collection.Find(context, filter, options.Find().SetLimit(1))
	if err != nil {
		return nil, err
	}
	var queryResult []bson.M
	if err = cursor.All(context, &queryResult); err != nil {
		return nil, err
	}
	return &queryResult, nil
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

func (us *UserService) CheckIfUsernameExist(username string) (bool, error) {
	filter := bson.M{"username": username}
	queryResult, err := us.collectionQuerys.find(context.TODO(), filter)
	if err != nil {
		return true, err
	}
	if len(*queryResult) == 0 {
		return false, nil
	} else {
		return true, nil
	}
}

func (us *UserService) CheckIfUserExist(id string) (bool, error) {
	filter := bson.M{"_id": id}
	queryResult, err := us.collectionQuerys.find(context.TODO(), filter)
	if err != nil {
		return true, err
	}
	if len(*queryResult) == 0 {
		return false, nil
	} else {
		return true, nil
	}
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
