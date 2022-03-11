package services

import (
	"context"
	"instagram-go/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type AuthenticationService struct {
	authenticationQuerys       authenticationQueryable
	authenticationVerification authenticationVerifyable
}

func NewAuthenticationService(authenticationQuerys authenticationQueryable, authentificationVerification authenticationVerifyable) *AuthenticationService {
	return &AuthenticationService{
		authenticationQuerys:       authenticationQuerys,
		authenticationVerification: authentificationVerification,
	}
}

type authenticationQueryable interface {
	findOneUser(context.Context, interface{}) (*models.User, error)
}

type authenticationQuery struct {
	collection *mongo.Collection
}

type authenticationVerifyable interface {
	compareHashAndPassword([]byte, []byte) error
}

type authenticationVerification struct {
}

func (av *authenticationVerification) compareHashAndPassword(hashedPassword []byte, password []byte) error {
	return bcrypt.CompareHashAndPassword(hashedPassword, password)
}

func NewAuthenticationVerification() *authenticationVerification {
	return &authenticationVerification{}
}

func NewAuthenticationQuery(collection *mongo.Collection) *authenticationQuery {
	return &authenticationQuery{
		collection: collection,
	}
}

func (aq *authenticationQuery) findOneUser(context context.Context, filter interface{}) (*models.User, error) {
	var user models.User
	err := aq.collection.FindOne(context, filter).Decode(&user)
	return &user, err
}

func (as *AuthenticationService) VerifyCredential(username string, password string) (string, error) {
	filter := bson.M{"username": username}
	user, err := as.authenticationQuerys.findOneUser(context.TODO(), filter)
	if err != nil {
		return "", err
	}
	err = as.authenticationVerification.compareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", err
	}

	return user.Id, err
}
