package users

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	collection *mongo.Collection
}

func NewUserService(collection *mongo.Collection) *UserService {
	return &UserService{
		collection: collection,
	}
}

func (us *UserService) insertUser(user User) error {
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

	result, err := us.collection.InsertOne(context.TODO(), newUser)
	if err != nil {
		return err
	}
	if result != nil {
		return nil
	}
	return nil
}

func (us *UserService) updateUser(newUserData User) error {
	var willBeUpdatedUser User
	filter := bson.M{"_id": newUserData.Id}
	err := us.collection.FindOne(context.TODO(), filter).Decode(&willBeUpdatedUser)
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
	_, err = us.collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}
	return nil
}
