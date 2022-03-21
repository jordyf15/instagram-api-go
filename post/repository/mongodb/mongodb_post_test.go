package mongodb_test

import (
	"context"
	"instagram-go/domain"
	"instagram-go/post/repository/mongodb"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestPostRepoSuite(t *testing.T) {
	suite.Run(t, new(PostRepoSuite))
}

type PostRepoSuite struct {
	suite.Suite
	collection *mongo.Collection
}

func (pr *PostRepoSuite) SetupSuite() {
	client, _ := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	pr.collection = client.Database("instagram_test").Collection("posts")
}

func (pr *PostRepoSuite) AfterTest(suiteName, testName string) {
	pr.collection.Drop(context.TODO())
}

func (pr *PostRepoSuite) TestInsertPostSuccessful() {
	newPost := domain.NewPost("postid1", "userid1", []string{}, "caption1", 0, time.Now(), time.Now())
	postRepo := mongodb.NewMongodbPostRepository(pr.collection)
	err := postRepo.InsertPost(newPost)

	var insertedPost domain.Post
	pr.collection.FindOne(context.TODO(), bson.M{"_id": "postid1"}).Decode(&insertedPost)
	assert.Equalf(pr.T(), newPost.Id, insertedPost.Id, "Should have return the correct id %s but got %s", newPost.Id, insertedPost.Id)
	assert.Equalf(pr.T(), newPost.UserId, insertedPost.UserId, "Should have return the correct userid %s but got %s", newPost.UserId, insertedPost.UserId)
	assert.Equalf(pr.T(), newPost.Caption, insertedPost.Caption, "Should have returned the correct caption %s but got %s", newPost.Caption, insertedPost.Caption)
	assert.NoError(pr.T(), err, "Should have not return error")
}

func (pr *PostRepoSuite) TestInsertDuplicatePost() {
	newPost := domain.NewPost("postid1", "userid1", []string{}, "caption1", 0, time.Now(), time.Now())
	postRepo := mongodb.NewMongodbPostRepository(pr.collection)
	_ = postRepo.InsertPost(newPost)
	err := postRepo.InsertPost(newPost)

	assert.Error(pr.T(), err, "Should have return error but didn't")
}

func (pr *PostRepoSuite) TestFindPostsSuccessful() {
	post := bson.M{
		"_id":               "postid1",
		"user_id":           "userid1",
		"visual_media_urls": nil,
		"caption":           "caption1",
		"created_date":      primitive.NewDateTimeFromTime(time.Now()),
		"updated_date":      primitive.NewDateTimeFromTime(time.Now()),
	}
	_, _ = pr.collection.InsertOne(context.TODO(), post)

	postRepo := mongodb.NewMongodbPostRepository(pr.collection)
	filter := bson.M{"_id": "postid1"}
	queryResult, err := postRepo.FindPosts(filter)

	assert.Equalf(pr.T(), 1, len(*queryResult), "Should have return the correct amount of post %v but got %v", 1, len(*queryResult))
	assert.Equalf(pr.T(), "postid1", (*queryResult)[0]["_id"], "Should have received the correct postid %s but got %s", "postid1", (*queryResult)[0]["_id"])
	assert.Equalf(pr.T(), "userid1", (*queryResult)[0]["user_id"], "Should have received the correct userid %s but got %s", "userid1", (*queryResult)[0]["user_id"])
	assert.Equalf(pr.T(), "caption1", (*queryResult)[0]["caption"], "Should have received the correct caption %s but got %s", "caption1", (*queryResult)[0]["caption"])
	assert.NoError(pr.T(), err, "Should have not returned error")
}

func (pr *PostRepoSuite) TestFindNotExistPosts() {
	post := bson.M{
		"_id":               "postid1",
		"user_id":           "userid1",
		"visual_media_urls": nil,
		"caption":           "caption1",
		"created_date":      primitive.NewDateTimeFromTime(time.Now()),
		"updated_date":      primitive.NewDateTimeFromTime(time.Now()),
	}
	_, _ = pr.collection.InsertOne(context.TODO(), post)

	postRepo := mongodb.NewMongodbPostRepository(pr.collection)
	filter := bson.M{"_id": "nonExistPost"}
	queryResult, err := postRepo.FindPosts(filter)

	assert.Equalf(pr.T(), 0, len(*queryResult), "Should have return the correct amount of post %v but got %v", 0, len(*queryResult))
	assert.NoError(pr.T(), err, "Should have not returned error")
}

func (pr *PostRepoSuite) TestFindOnePostSuccessful() {
	post := bson.M{
		"_id":               "postid1",
		"user_id":           "userid1",
		"visual_media_urls": nil,
		"caption":           "caption1",
		"created_date":      primitive.NewDateTimeFromTime(time.Now()),
		"updated_date":      primitive.NewDateTimeFromTime(time.Now()),
	}
	_, _ = pr.collection.InsertOne(context.TODO(), post)

	postRepo := mongodb.NewMongodbPostRepository(pr.collection)
	foundPost, err := postRepo.FindOnePost("postid1")
	assert.Equalf(pr.T(), "postid1", foundPost.Id, "Should have return the correct _id %s but got %s", "postid1", foundPost.Id)
	assert.Equalf(pr.T(), "userid1", foundPost.UserId, "Should have return the correct user_id %s but got %s", "userid1", foundPost.UserId)
	assert.Equalf(pr.T(), "caption1", foundPost.Caption, "Should have return the correct caption %s but got %s", "caption1", foundPost.Caption)
	assert.NoError(pr.T(), err, "Should have not return error")
}

func (pr *PostRepoSuite) TestFindOneNotExistPost() {
	postRepo := mongodb.NewMongodbPostRepository(pr.collection)
	_, err := postRepo.FindOnePost("notExistPostId")
	assert.Error(pr.T(), err, "Should return error but didn't")
}

func (pr *PostRepoSuite) TestUpdateNotExistPost() {
	post := bson.M{
		"_id":               "postid1",
		"user_id":           "userid1",
		"visual_media_urls": nil,
		"caption":           "caption1",
		"created_date":      primitive.NewDateTimeFromTime(time.Now()),
		"updated_date":      primitive.NewDateTimeFromTime(time.Now()),
	}
	_, _ = pr.collection.InsertOne(context.TODO(), post)

	postRepo := mongodb.NewMongodbPostRepository(pr.collection)
	err := postRepo.UpdatePost("notExistPostId", "newcaption1")

	var notUpdatedPost domain.Post
	pr.collection.FindOne(context.TODO(), bson.M{"_id": "postid1"}).Decode(&notUpdatedPost)
	assert.Equalf(pr.T(), "postid1", notUpdatedPost.Id, "Should have received the old _id %s but got %s", "postid1", notUpdatedPost.Id)
	assert.Equalf(pr.T(), "userid1", notUpdatedPost.UserId, "Should have received the old user_id %s but got %s", "userid1", notUpdatedPost.UserId)
	assert.Equalf(pr.T(), "caption1", notUpdatedPost.Caption, "Should have received the old caption %s but got %s", "newcaption1", notUpdatedPost.Caption)
	assert.NoErrorf(pr.T(), err, "Should have not return error but got %s", err)
}

func (pr *PostRepoSuite) TestUpdatePostSuccessful() {
	post := bson.M{
		"_id":               "postid1",
		"user_id":           "userid1",
		"visual_media_urls": nil,
		"caption":           "caption1",
		"created_date":      primitive.NewDateTimeFromTime(time.Now()),
		"updated_date":      primitive.NewDateTimeFromTime(time.Now()),
	}
	_, _ = pr.collection.InsertOne(context.TODO(), post)

	postRepo := mongodb.NewMongodbPostRepository(pr.collection)
	err := postRepo.UpdatePost("postid1", "newcaption1")

	var updatedPost domain.Post
	pr.collection.FindOne(context.TODO(), bson.M{"_id": "postid1"}).Decode(&updatedPost)
	assert.Equalf(pr.T(), "postid1", updatedPost.Id, "Should have received the correct _id %s but got %s", "postid1", updatedPost.Id)
	assert.Equalf(pr.T(), "userid1", updatedPost.UserId, "Should have received the correct user_id %s but got %s", "userid1", updatedPost.UserId)
	assert.Equalf(pr.T(), "newcaption1", updatedPost.Caption, "Should have received the correct caption %s but got %s", "newcaption1", updatedPost.Caption)
	assert.NoErrorf(pr.T(), err, "Should have not return error but got %s", err)
}

func (pr *PostRepoSuite) TestDeletePostSuccessful() {
	post := bson.M{
		"_id":               "postid1",
		"user_id":           "userid1",
		"visual_media_urls": nil,
		"caption":           "caption1",
		"created_date":      primitive.NewDateTimeFromTime(time.Now()),
		"updated_date":      primitive.NewDateTimeFromTime(time.Now()),
	}
	_, _ = pr.collection.InsertOne(context.TODO(), post)

	postRepo := mongodb.NewMongodbPostRepository(pr.collection)
	err := postRepo.DeletePost("postid1")

	filter := bson.M{"_id": "postid1"}
	cursor, _ := pr.collection.Find(context.TODO(), filter)
	var queryResult []bson.M
	cursor.All(context.TODO(), &queryResult)

	assert.Equalf(pr.T(), 0, len(queryResult), "Should have return the correct amount of post %s but got %s", 0, len(queryResult))
	assert.NoError(pr.T(), err, "Should have not return error")
}

func (pr *PostRepoSuite) TestDeleteNotExistPost() {
	post := bson.M{
		"_id":               "postid1",
		"user_id":           "userid1",
		"visual_media_urls": nil,
		"caption":           "caption1",
		"created_date":      primitive.NewDateTimeFromTime(time.Now()),
		"updated_date":      primitive.NewDateTimeFromTime(time.Now()),
	}
	_, _ = pr.collection.InsertOne(context.TODO(), post)
	postRepo := mongodb.NewMongodbPostRepository(pr.collection)
	err := postRepo.DeletePost("notExistPostId")

	filter := bson.M{}
	cursor, _ := pr.collection.Find(context.TODO(), filter)
	var queryResult []bson.M
	cursor.All(context.TODO(), &queryResult)

	assert.Equalf(pr.T(), 1, len(queryResult), "Should have return the correct amount of post %v but got %v", 1, len(queryResult))
	assert.NoError(pr.T(), err, "Should have not return error")
}
