package mongodb_test

import (
	"context"
	"instagram-go/comment/repository/mongodb"
	"instagram-go/domain"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestCommentRepoSuite(t *testing.T) {
	suite.Run(t, new(CommentRepoSuite))
}

type CommentRepoSuite struct {
	suite.Suite
	collection *mongo.Collection
}

func (cr *CommentRepoSuite) SetupSuite() {
	client, _ := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	cr.collection = client.Database("instagram_test").Collection("comments")
}

func (cr *CommentRepoSuite) AfterTest(suiteName, testName string) {
	cr.collection.Drop(context.TODO())
}

func (cr *CommentRepoSuite) TestFindNotExistComment() {
	comment := bson.M{
		"_id":          "commentid1",
		"post_id":      "postid1",
		"user_id":      "userid1",
		"comment":      "comment1",
		"created_date": primitive.NewDateTimeFromTime(time.Now()),
		"updated_date": primitive.NewDateTimeFromTime(time.Now()),
	}
	_, _ = cr.collection.InsertOne(context.TODO(), comment)

	commentRepo := mongodb.NewMongodbCommentRepository(cr.collection)
	filter := bson.M{"_id": "notExistCommentId"}
	queryResult, err := commentRepo.FindComments(filter)

	assert.Equalf(cr.T(), 0, len(*queryResult), "Should have return the correct amount of comments %v but got %v", 0, len(*queryResult))
	assert.NoError(cr.T(), err, "Should have not return error")
}

func (cr *CommentRepoSuite) TestFindCommentsSuccessful() {
	comment := bson.M{
		"_id":          "commentid1",
		"post_id":      "postid1",
		"user_id":      "userid1",
		"comment":      "comment1",
		"created_date": primitive.NewDateTimeFromTime(time.Now()),
		"updated_date": primitive.NewDateTimeFromTime(time.Now()),
	}
	_, _ = cr.collection.InsertOne(context.TODO(), comment)

	commentRepo := mongodb.NewMongodbCommentRepository(cr.collection)
	filter := bson.M{"_id": "commentid1"}
	queryResult, err := commentRepo.FindComments(filter)

	assert.Equalf(cr.T(), 1, len(*queryResult), "Should have return the correct amount of comments %v but got %v", 1, len(*queryResult))
	assert.Equalf(cr.T(), "commentid1", (*queryResult)[0]["_id"], "Should have return the correct comment id %s but got %s", "commentid1", (*queryResult)[0]["_id"])
	assert.Equalf(cr.T(), "postid1", (*queryResult)[0]["post_id"], "Should have return the correct post id %s but got %s", "postid1", (*queryResult)[0]["post_id"])
	assert.Equalf(cr.T(), "userid1", (*queryResult)[0]["user_id"], "Should have return the correct user id %s but got %s", "userid1", (*queryResult)[0]["user_id"])
	assert.Equalf(cr.T(), "comment1", (*queryResult)[0]["comment"], "Should have return the correct comment %s but got %s", "comment1", (*queryResult)[0]["comment"])
	assert.NoError(cr.T(), err, "Should have not return error")
}

func (cr *CommentRepoSuite) TestInsertDuplicateComment() {
	newComment := domain.NewComment("commentid1", "postid1", "userid1", "comment1", 0, time.Now(), time.Now())
	commentRepo := mongodb.NewMongodbCommentRepository(cr.collection)
	_ = commentRepo.InsertComment(newComment)
	err := commentRepo.InsertComment(newComment)

	assert.Error(cr.T(), err, "Should have return an error but didn't")
}

func (cr *CommentRepoSuite) TestInsertCommentSuccessful() {
	newComment := domain.NewComment("commentid1", "postid1", "userid1", "comment1", 0, time.Now(), time.Now())
	commentRepo := mongodb.NewMongodbCommentRepository(cr.collection)
	err := commentRepo.InsertComment(newComment)

	var insertedComment domain.Comment
	cr.collection.FindOne(context.TODO(), bson.M{"_id": "commentid1"}).Decode(&insertedComment)
	assert.Equalf(cr.T(), "commentid1", insertedComment.Id, "Should have received the correct comment id %s but got %s", "commentid1", insertedComment.Id)
	assert.Equalf(cr.T(), "userid1", insertedComment.UserId, "Should have received the correct user id %s but got %s", "userid1", insertedComment.UserId)
	assert.Equalf(cr.T(), "postid1", insertedComment.PostId, "Should have received the correct post id %s but got %s", "postid1", insertedComment.PostId)
	assert.Equalf(cr.T(), "comment1", insertedComment.Comment, "Should have received the correct comment %s but got %s", "comment1", insertedComment.Comment)
	assert.NoError(cr.T(), err, "Should have not return an error")
}

func (cr *CommentRepoSuite) TestFindOneNotExistComment() {
	commentRepo := mongodb.NewMongodbCommentRepository(cr.collection)
	_, err := commentRepo.FindOneComment("notExistCommentId")

	assert.Error(cr.T(), err, "Should have return an error but didn't")
}

func (cr *CommentRepoSuite) TestFindOneCommentSuccessful() {
	comment := bson.M{
		"_id":          "commentid1",
		"post_id":      "postid1",
		"user_id":      "userid1",
		"comment":      "comment1",
		"created_date": primitive.NewDateTimeFromTime(time.Now()),
		"updated_date": primitive.NewDateTimeFromTime(time.Now()),
	}
	_, _ = cr.collection.InsertOne(context.TODO(), comment)

	commentRepo := mongodb.NewMongodbCommentRepository(cr.collection)
	foundComment, err := commentRepo.FindOneComment("commentid1")

	assert.Equalf(cr.T(), "commentid1", foundComment.Id, "Should have return the correct comment id %s but got %s", "commentid1", foundComment.Id)
	assert.Equalf(cr.T(), "postid1", foundComment.PostId, "Should have return the correct post id %s but got %s", "postid1", foundComment.PostId)
	assert.Equalf(cr.T(), "userid1", foundComment.UserId, "Should have return the correct user id %s but got %s", "userid1", foundComment.UserId)
	assert.Equalf(cr.T(), "comment1", foundComment.Comment, "Should have return the correct comment %s but got %s", "comment1", foundComment.Comment)
	assert.NoError(cr.T(), err, "Should have not return an error")
}

func (cr *CommentRepoSuite) TestUpdateNotExistComment() {
	comment := bson.M{
		"_id":          "commentid1",
		"post_id":      "postid1",
		"user_id":      "userid1",
		"comment":      "comment1",
		"created_date": primitive.NewDateTimeFromTime(time.Now()),
		"updated_date": primitive.NewDateTimeFromTime(time.Now()),
	}
	_, _ = cr.collection.InsertOne(context.TODO(), comment)

	commentRepo := mongodb.NewMongodbCommentRepository(cr.collection)
	err := commentRepo.UpdateComment("notExistCommentId", "newcomment1")

	var notUpdatedComment domain.Comment
	cr.collection.FindOne(context.TODO(), bson.M{"_id": "commentid1"}).Decode(&notUpdatedComment)

	assert.Equalf(cr.T(), "commentid1", notUpdatedComment.Id, "Should have received the old comment id %s but got %s", "commentid1", notUpdatedComment.Id)
	assert.Equalf(cr.T(), "postid1", notUpdatedComment.PostId, "Should have received the old post id %s but got %s", "postid1", notUpdatedComment.PostId)
	assert.Equalf(cr.T(), "userid1", notUpdatedComment.UserId, "Should have received the old user id %s but got %s", "userid1", notUpdatedComment.UserId)
	assert.Equalf(cr.T(), "comment1", notUpdatedComment.Comment, "Should have received the old comment %s but got %s", "comment1", notUpdatedComment.Comment)
	assert.NoError(cr.T(), err, "Should have not return error")
}

func (cr *CommentRepoSuite) TestUpdateCommentSuccessful() {
	comment := bson.M{
		"_id":          "commentid1",
		"post_id":      "postid1",
		"user_id":      "userid1",
		"comment":      "comment1",
		"created_date": primitive.NewDateTimeFromTime(time.Now()),
		"updated_date": primitive.NewDateTimeFromTime(time.Now()),
	}
	_, _ = cr.collection.InsertOne(context.TODO(), comment)

	commentRepo := mongodb.NewMongodbCommentRepository(cr.collection)
	err := commentRepo.UpdateComment("commentid1", "newcomment1")

	var updatedComment domain.Comment
	cr.collection.FindOne(context.TODO(), bson.M{"_id": "commentid1"}).Decode(&updatedComment)

	assert.Equalf(cr.T(), "commentid1", updatedComment.Id, "Should have received the correct comment id %s but got %s", "commentid1", updatedComment.Id)
	assert.Equalf(cr.T(), "postid1", updatedComment.PostId, "Should have received the correct post id %s but got %s", "postid1", updatedComment.PostId)
	assert.Equalf(cr.T(), "userid1", updatedComment.UserId, "Should have received the correct user id %s but got %s", "userid1", updatedComment.UserId)
	assert.Equalf(cr.T(), "newcomment1", updatedComment.Comment, "Should have received the correct comment %s but got %s", "newcomment1", updatedComment.Comment)
	assert.NoError(cr.T(), err, "Should have not return error")
}

func (cr *CommentRepoSuite) TestDeleteNotExistComment() {
	comment := bson.M{
		"_id":          "commentid1",
		"post_id":      "postid1",
		"user_id":      "userid1",
		"comment":      "comment1",
		"created_date": primitive.NewDateTimeFromTime(time.Now()),
		"updated_date": primitive.NewDateTimeFromTime(time.Now()),
	}
	_, _ = cr.collection.InsertOne(context.TODO(), comment)

	commentRepo := mongodb.NewMongodbCommentRepository(cr.collection)
	err := commentRepo.DeleteComment("notExistCommentId")

	filter := bson.M{}
	cursor, _ := cr.collection.Find(context.TODO(), filter)
	var queryResult []bson.M
	cursor.All(context.TODO(), &queryResult)

	assert.Equalf(cr.T(), 1, len(queryResult), "Should have return the correct amount of comments %v but got %v", 1, len(queryResult))
	assert.NoError(cr.T(), err, "Should have not return error")
}

func (cr *CommentRepoSuite) TestDeleteCommentSuccessful() {
	comment := bson.M{
		"_id":          "commentid1",
		"post_id":      "postid1",
		"user_id":      "userid1",
		"comment":      "comment1",
		"created_date": primitive.NewDateTimeFromTime(time.Now()),
		"updated_date": primitive.NewDateTimeFromTime(time.Now()),
	}
	_, _ = cr.collection.InsertOne(context.TODO(), comment)

	commentRepo := mongodb.NewMongodbCommentRepository(cr.collection)
	err := commentRepo.DeleteComment("commentid1")

	filter := bson.M{"_id": "commentid1"}
	cursor, _ := cr.collection.Find(context.TODO(), filter)
	var queryResult []bson.M
	cursor.All(context.TODO(), &queryResult)

	assert.Equalf(cr.T(), 0, len(queryResult), "Should have return the correct amount of comments %v but got %v", 0, len(queryResult))
	assert.NoError(cr.T(), err, "Should have not return error")
}
