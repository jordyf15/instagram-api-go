package mongodb

import (
	"context"
	"instagram-go/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongodbUserRepository struct {
	collection *mongo.Collection
}

func NewMongodbUserRepository(collection *mongo.Collection) domain.UserRepository {
	return &mongodbUserRepository{
		collection: collection,
	}
}

func (mur *mongodbUserRepository) InsertUser(user *domain.User) error {
	newUser := bson.D{
		primitive.E{Key: "_id", Value: user.Id},
		primitive.E{Key: "username", Value: user.Username},
		primitive.E{Key: "full_name", Value: user.Fullname},
		primitive.E{Key: "password", Value: user.Password},
		primitive.E{Key: "email", Value: user.Email},
		primitive.E{Key: "profile_pictures", Value: nil},
	}
	_, err := mur.collection.InsertOne(context.TODO(), newUser)
	if err != nil {
		return err
	}
	return nil
}

func (mur *mongodbUserRepository) UpdateUser(newUserData *domain.User) error {
	filter := bson.M{"_id": newUserData.Id}
	update := bson.D{primitive.E{Key: "$set", Value: bson.D{
		primitive.E{Key: "username", Value: newUserData.Username},
		primitive.E{Key: "full_name", Value: newUserData.Fullname},
		primitive.E{Key: "password", Value: newUserData.Password},
		primitive.E{Key: "email", Value: newUserData.Email},
		primitive.E{Key: "profile_pictures", Value: newUserData.ProfilePictures},
	},
	},
	}
	_, err := mur.collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (mur *mongodbUserRepository) FindUser(filter interface{}) (*[]bson.M, error) {
	cursor, err := mur.collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	var queryResult []bson.M
	if err = cursor.All(context.TODO(), &queryResult); err != nil {
		return nil, err
	}
	return &queryResult, nil
}

func (mur *mongodbUserRepository) FindOneUser(filter interface{}) (*domain.User, error) {
	var user domain.User
	err := mur.collection.FindOne(context.TODO(), filter).Decode(&user)
	return &user, err
}
