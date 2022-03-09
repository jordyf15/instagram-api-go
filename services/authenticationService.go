package services

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type AuthenticationService struct {
	collection *mongo.Collection
}

func NewAuthenticationService(collection *mongo.Collection) *AuthenticationService {
	return &AuthenticationService{
		collection: collection,
	}
}

func (as *AuthenticationService) VerifyCredential(username string, password string) (string, error) {
	var user user
	filter := bson.M{"username": username}
	err := as.collection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		return "", err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", err
	}

	return user.ID, err
}

type user struct {
	ID       string `bson:"_id"`
	Username string
	Password string
}
