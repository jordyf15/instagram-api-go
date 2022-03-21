package mongodb_test

import (
	"context"
	"instagram-go/domain"
	"instagram-go/like/repository/mongodb"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestLikeRepoSuite(t *testing.T) {
	suite.Run(t, new(LikeRepoSuite))
}

type LikeRepoSuite struct {
	suite.Suite
	collection *mongo.Collection
}

func (lr *LikeRepoSuite) SetupSuite() {
	client, _ := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	lr.collection = client.Database("instagram_test").Collection("likes")
}

func (lr *LikeRepoSuite) AfterTest(suiteName, testName string) {
	lr.collection.Drop(context.TODO())
}

func (lr *LikeRepoSuite) TestInsertDuplicateLike() {
	likeRepo := mongodb.NewMongodbLikeRepository(lr.collection)
	newLike := domain.NewLike("likeid1", "userid1", "postid1", "post")

	_ = likeRepo.InsertLike(newLike)
	err := likeRepo.InsertLike(newLike)

	assert.Error(lr.T(), err, "Should have return an error but didn't")
}

func (lr *LikeRepoSuite) TestInsertLikeSuccessful() {
	likeRepo := mongodb.NewMongodbLikeRepository(lr.collection)
	newLike := domain.NewLike("likeid1", "userid1", "postid1", "post")

	err := likeRepo.InsertLike(newLike)

	var insertedLike domain.Like
	lr.collection.FindOne(context.TODO(), bson.M{"_id": "likeid1"}).Decode(&insertedLike)
	assert.Equalf(lr.T(), newLike.Id, insertedLike.Id, "Should have return the correct like id %s but got %s", newLike.Id, insertedLike.Id)
	assert.Equalf(lr.T(), newLike.UserId, insertedLike.UserId, "Should have return the correct user id %s but got %s", newLike.UserId, insertedLike.UserId)
	assert.Equalf(lr.T(), newLike.ResourceId, insertedLike.ResourceId, "Should have return the correct resource id %s but got %s", newLike.ResourceId, insertedLike.ResourceId)
	assert.Equalf(lr.T(), newLike.ResourceType, insertedLike.ResourceType, "Should have return the correct resource type %s but got %s", newLike.ResourceType, insertedLike.ResourceType)
	assert.NoError(lr.T(), err, "Should have not return an error")
}

func (lr *LikeRepoSuite) TestFindNotExistLikes() {
	like := bson.M{
		"_id":           "likeid1",
		"user_id":       "userid1",
		"resource_id":   "postid1",
		"resource_type": "post",
	}
	_, _ = lr.collection.InsertOne(context.TODO(), like)

	likeRepo := mongodb.NewMongodbLikeRepository(lr.collection)
	filter := bson.M{"_id": "notExistLikeId"}
	queryResult, err := likeRepo.FindLikes(filter)

	assert.Equalf(lr.T(), 0, len(*queryResult), "Should have return the correct amount of like %v but got %v", 0, len(*queryResult))
	assert.NoError(lr.T(), err, "Should have not return error")
}

func (lr *LikeRepoSuite) TestFindLikesSuccessful() {
	like := bson.M{
		"_id":           "likeid1",
		"user_id":       "userid1",
		"resource_id":   "postid1",
		"resource_type": "post",
	}
	_, _ = lr.collection.InsertOne(context.TODO(), like)

	likeRepo := mongodb.NewMongodbLikeRepository(lr.collection)
	filter := bson.M{"_id": "likeid1"}
	queryResult, err := likeRepo.FindLikes(filter)

	assert.Equalf(lr.T(), 1, len(*queryResult), "Should have return the correct amount of like %v but got %v", 1, len(*queryResult))
	assert.Equalf(lr.T(), "likeid1", (*queryResult)[0]["_id"], "Should have return the correct like id %s but got %s", "likeid1", (*queryResult)[0]["_id"])
	assert.Equalf(lr.T(), "userid1", (*queryResult)[0]["user_id"], "Should have return the correct user id %s but got %s", "userid1", (*queryResult)[0]["user_id"])
	assert.Equalf(lr.T(), "postid1", (*queryResult)[0]["resource_id"], "Should have return the correct resource id %s but got %s", "postid1", (*queryResult)[0]["resource_id"])
	assert.Equalf(lr.T(), "post", (*queryResult)[0]["resource_type"], "Should have return the correct resource type %s but got %s", "post", (*queryResult)[0]["resource_type"])
	assert.NoError(lr.T(), err, "Should have not return error")
}

func (lr *LikeRepoSuite) TestFindOneNotExistLike() {
	likeRepo := mongodb.NewMongodbLikeRepository(lr.collection)
	_, err := likeRepo.FindOneLike("notExistLikeId")
	assert.Error(lr.T(), err, "Should return error but didn't")
}

func (lr *LikeRepoSuite) TestFindOneLikeSuccessful() {
	like := bson.M{
		"_id":           "likeid1",
		"user_id":       "userid1",
		"resource_id":   "postid1",
		"resource_type": "post",
	}
	_, _ = lr.collection.InsertOne(context.TODO(), like)

	likeRepo := mongodb.NewMongodbLikeRepository(lr.collection)
	foundLike, err := likeRepo.FindOneLike("likeid1")

	assert.Equalf(lr.T(), "likeid1", foundLike.Id, "Should have return the correct like id %s but got %s", "likeid1", foundLike.Id)
	assert.Equalf(lr.T(), "userid1", foundLike.UserId, "Should have return the correct user id %s but got %s", "userid1", foundLike.UserId)
	assert.Equalf(lr.T(), "postid1", foundLike.ResourceId, "Should have return the correct resource id %s but got %s", "postid1", foundLike.ResourceId)
	assert.Equalf(lr.T(), "post", foundLike.ResourceType, "Should have return the correct resource type %s but got %s", "post", foundLike.ResourceType)
	assert.NoError(lr.T(), err, "Should have not return error")
}

func (lr *LikeRepoSuite) TestDeleteNotExistLike() {
	like := bson.M{
		"_id":           "likeid1",
		"user_id":       "userid1",
		"resource_id":   "postid1",
		"resource_type": "post",
	}
	_, _ = lr.collection.InsertOne(context.TODO(), like)

	likeRepo := mongodb.NewMongodbLikeRepository(lr.collection)
	err := likeRepo.DeleteLike("notExistLike")

	filter := bson.M{}
	cursor, _ := lr.collection.Find(context.TODO(), filter)
	var queryResult []bson.M
	cursor.All(context.TODO(), &queryResult)

	assert.Equalf(lr.T(), 1, len(queryResult), "Should have return the correct amount of likes %v but got %v", 1, len(queryResult))
	assert.NoError(lr.T(), err, "Should have not return error")
}

func (lr *LikeRepoSuite) TestDeleteLikeSuccessful() {
	like := bson.M{
		"_id":           "likeid1",
		"user_id":       "userid1",
		"resource_id":   "postid1",
		"resource_type": "post",
	}
	_, _ = lr.collection.InsertOne(context.TODO(), like)

	likeRepo := mongodb.NewMongodbLikeRepository(lr.collection)
	err := likeRepo.DeleteLike("likeid1")

	filter := bson.M{"_id": "likeid1"}
	cursor, _ := lr.collection.Find(context.TODO(), filter)
	var queryResult []bson.M
	cursor.All(context.TODO(), &queryResult)

	assert.Equalf(lr.T(), 0, len(queryResult), "Should have return the correct amount of likes %v but got %v", 0, len(queryResult))
	assert.NoError(lr.T(), err, "Should have not return error")
}
