package mongodb_test

import (
	"context"
	"instagram-go/domain"
	"instagram-go/user/repository/mongodb"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestUserRepoSuite(t *testing.T) {
	suite.Run(t, new(UserRepoSuite))
}

type UserRepoSuite struct {
	suite.Suite
	collection *mongo.Collection
}

func (ur *UserRepoSuite) SetupSuite() {
	client, _ := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	ur.collection = client.Database("instagram_test").Collection("users")
}

func (ur *UserRepoSuite) AfterTest(suiteName, testName string) {
	ur.collection.Drop(context.TODO())
}

func (ur *UserRepoSuite) TestInsertDuplicateUser() {
	newUser := domain.NewUser("userid1", "username1", "fullname1", "password1", "email1@gmail.com", nil)
	userRepo := mongodb.NewMongodbUserRepository(ur.collection)

	_ = userRepo.InsertUser(newUser)
	err := userRepo.InsertUser(newUser)

	assert.Error(ur.T(), err, "Should have return error but didn't")
}

func (ur *UserRepoSuite) TestInsertUserSuccessful() {
	newUser := domain.NewUser("userid1", "username1", "fullname1", "password1", "email1@gmail.com", nil)
	userRepo := mongodb.NewMongodbUserRepository(ur.collection)
	err := userRepo.InsertUser(newUser)

	var insertedUser domain.User
	ur.collection.FindOne(context.TODO(), bson.M{"_id": "userid1"}).Decode(&insertedUser)
	assert.Equalf(ur.T(), insertedUser.Id, newUser.Id, "Should have returned the inserted user id %s but got %s", newUser.Id, insertedUser.Id)
	assert.Equalf(ur.T(), insertedUser.Username, newUser.Username, "Should have returned the inserted username %s but got %s", newUser.Username, insertedUser.Username)
	assert.Equalf(ur.T(), insertedUser.Fullname, newUser.Fullname, "Should have returned the inserted fullname %s but got %s", newUser.Fullname, insertedUser.Fullname)
	assert.Equalf(ur.T(), insertedUser.Password, newUser.Password, "Should have returned the inserted password %s but got %s", newUser.Password, insertedUser.Password)
	assert.Equalf(ur.T(), insertedUser.Email, newUser.Email, "Should have returned the inserted email %s but got %s", newUser.Email, insertedUser.Email)
	assert.NoErrorf(ur.T(), err, "Should have not return error but got %s", err)
}

func (ur *UserRepoSuite) TestUpdateNotExistUser() {
	user := bson.M{
		"_id":              "userid1",
		"username":         "username1",
		"full_name":        "fullname1",
		"password":         "password1",
		"email":            "email1",
		"profile_pictures": nil,
	}
	_, _ = ur.collection.InsertOne(context.TODO(), user)

	newUserData := domain.NewUser("notExistUser", "newusername1", "newfullname1", "newpassword1", "email1@gmail.com", nil)
	userRepo := mongodb.NewMongodbUserRepository(ur.collection)
	err := userRepo.UpdateUser(newUserData)

	var notUpdatedUser domain.User
	ur.collection.FindOne(context.TODO(), bson.M{"_id": "userid1"}).Decode(&notUpdatedUser)
	assert.Equalf(ur.T(), "userid1", notUpdatedUser.Id, "Should have returned the old user id %s but got %s", "userid1", notUpdatedUser.Id)
	assert.Equalf(ur.T(), "username1", notUpdatedUser.Username, "Should have returned the old username %s but got %s", "username1", notUpdatedUser.Username)
	assert.Equalf(ur.T(), "fullname1", notUpdatedUser.Fullname, "Should have returned the old fullname %s but got %s", "fullanme1", notUpdatedUser.Fullname)
	assert.Equalf(ur.T(), "password1", notUpdatedUser.Password, "Should have returned the old password %s but got %s", "password1", notUpdatedUser.Password)
	assert.Equalf(ur.T(), "email1", notUpdatedUser.Email, "Should have returned the old email %s but got %s", "email1", notUpdatedUser.Email)
	assert.NoErrorf(ur.T(), err, "Should have not return error but got %s", err)
}

func (ur *UserRepoSuite) TestUpdateUserSuccessful() {
	user := bson.M{
		"_id":              "userid1",
		"username":         "username1",
		"full_name":        "fullname1",
		"password":         "password1",
		"email":            "email1",
		"profile_pictures": nil,
	}
	_, _ = ur.collection.InsertOne(context.TODO(), user)

	newUserData := domain.NewUser("userid1", "new username 1", "new fullname 1", "new password 1", "email1@gmail.com", nil)
	userRepo := mongodb.NewMongodbUserRepository(ur.collection)
	err := userRepo.UpdateUser(newUserData)

	var updatedUser domain.User
	ur.collection.FindOne(context.TODO(), bson.M{"_id": "userid1"}).Decode(&updatedUser)
	assert.Equalf(ur.T(), newUserData.Id, updatedUser.Id, "Should have returned the updated user id %s but got %s", newUserData.Id, updatedUser.Id)
	assert.Equalf(ur.T(), newUserData.Username, updatedUser.Username, "Should have returned the updated username %s but got %s", newUserData.Username, updatedUser.Username)
	assert.Equalf(ur.T(), newUserData.Fullname, updatedUser.Fullname, "Should have returned the updated fullname %s but got %s", newUserData.Fullname, updatedUser.Fullname)
	assert.Equalf(ur.T(), newUserData.Password, updatedUser.Password, "Should have returned the updated password %s but got %s", newUserData.Password, updatedUser.Password)
	assert.Equalf(ur.T(), newUserData.Email, updatedUser.Email, "Should have returned the updated email %s but got %s", newUserData.Email, updatedUser.Email)
	assert.NoErrorf(ur.T(), err, "Should have not return error but got %s", err)
}

func (ur *UserRepoSuite) TestFindNotExistUser() {
	user := bson.M{
		"_id":              "userid1",
		"username":         "username1",
		"full_name":        "fullname1",
		"password":         "password1",
		"email":            "email1",
		"profile_pictures": nil,
	}
	_, _ = ur.collection.InsertOne(context.TODO(), user)

	userRepo := mongodb.NewMongodbUserRepository(ur.collection)
	filter := bson.M{"_id": "nonExistUserId"}
	queryResult, err := userRepo.FindUser(filter)

	assert.Equalf(ur.T(), 0, len(*queryResult), "Should have return the correct amount of user: %v but got %v", 0, len(*queryResult))
	assert.NoErrorf(ur.T(), err, "Should have not return error but got %s", err)
}

func (ur *UserRepoSuite) TestFindUserSuccessful() {
	user := bson.M{
		"_id":              "userid1",
		"username":         "username1",
		"full_name":        "fullname1",
		"password":         "password1",
		"email":            "email1",
		"profile_pictures": nil,
	}
	_, _ = ur.collection.InsertOne(context.TODO(), user)

	userRepo := mongodb.NewMongodbUserRepository(ur.collection)
	filter := bson.M{"_id": "userid1"}
	queryResult, err := userRepo.FindUser(filter)

	assert.Equalf(ur.T(), 1, len(*queryResult), "Should have return the correct amount of user: %v but got %v", 1, len(*queryResult))
	assert.Equalf(ur.T(), "userid1", (*queryResult)[0]["_id"], "Should have received the right id of user: %s but got %s", "userid1", (*queryResult)[0]["_id"])
	assert.NoErrorf(ur.T(), err, "Should have not return error but got %s", err)
}

func (ur *UserRepoSuite) TestFindOneNotExistUser() {
	userRepo := mongodb.NewMongodbUserRepository(ur.collection)
	filter := bson.M{"_id": "notExistUserId"}
	_, err := userRepo.FindOneUser(filter)
	assert.Error(ur.T(), err, "Should have return error but didn't")
}

func (ur *UserRepoSuite) TestFindOneUserSuccessful() {
	user := bson.M{
		"_id":              "userid1",
		"username":         "username1",
		"full_name":        "fullname1",
		"password":         "password1",
		"email":            "email1",
		"profile_pictures": nil,
	}
	_, _ = ur.collection.InsertOne(context.TODO(), user)

	userRepo := mongodb.NewMongodbUserRepository(ur.collection)
	filter := bson.M{"_id": "userid1"}
	foundUser, err := userRepo.FindOneUser(filter)
	assert.Equalf(ur.T(), "userid1", foundUser.Id, "Should have return the correct user id: %s but got %s", "userid1", foundUser.Id)
	assert.Equalf(ur.T(), "username1", foundUser.Username, "Should have return the correct user id: %s but got %s", "username1", foundUser.Username)
	assert.Equalf(ur.T(), "fullname1", foundUser.Fullname, "Should have return the correct user id: %s but got %s", "fullname1", foundUser.Fullname)
	assert.Equalf(ur.T(), "password1", foundUser.Password, "Should have return the correct user id: %s but got %s", "password1", foundUser.Password)
	assert.Equalf(ur.T(), "email1", foundUser.Email, "Should have return the correct user id: %s but got %s", "email1", foundUser.Email)
	assert.NoErrorf(ur.T(), err, "Should have not return error but got %s", err)
}
