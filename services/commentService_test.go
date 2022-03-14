package services

import (
	"context"
	"errors"
	"instagram-go/models"
	"testing"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var findOneCommentMock func(context.Context, interface{}) (*models.Comment, error)
var findCommentLikeMock func(context.Context, interface{}) (*[]bson.M, error)
var findPostCommentMock func(context.Context, interface{}) (*[]bson.M, error)
var insertOneCommentMock func(context.Context, interface{}) (*mongo.InsertOneResult, error)
var updateOneCommentMock func(context.Context, interface{}, interface{}) (*mongo.UpdateResult, error)
var deleteOneCommentMock func(context.Context, interface{}) (*mongo.DeleteResult, error)
var findCommentMock func(context.Context, interface{}) (*[]bson.M, error)

type commentCollectionQueryMock struct {
}

func newCommentCollectionQueryMock() *commentCollectionQueryMock {
	return &commentCollectionQueryMock{}
}

func (ccqm *commentCollectionQueryMock) findComment(context context.Context, filter interface{}) (*[]bson.M, error) {
	return findCommentMock(context, filter)
}

func (ccqm *commentCollectionQueryMock) findOneComment(context context.Context, filter interface{}) (*models.Comment, error) {
	return findOneCommentMock(context, filter)
}

func (ccqm *commentCollectionQueryMock) findCommentLike(context context.Context, filter interface{}) (*[]bson.M, error) {
	return findCommentLikeMock(context, filter)
}

func (ccqm *commentCollectionQueryMock) findPostComment(context context.Context, filter interface{}) (*[]bson.M, error) {
	return findPostCommentMock(context, filter)
}

func (ccqm *commentCollectionQueryMock) insertOneComment(context context.Context, document interface{}) (*mongo.InsertOneResult, error) {
	return insertOneCommentMock(context, document)
}

func (ccqm *commentCollectionQueryMock) updateOneComment(context context.Context, filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
	return updateOneCommentMock(context, filter, update)
}

func (ccqm *commentCollectionQueryMock) deleteOneComment(context context.Context, filter interface{}) (*mongo.DeleteResult, error) {
	return deleteOneCommentMock(context, filter)
}

func TestCheckIfCommentExist(t *testing.T) {
	commentColQueryMock := newCommentCollectionQueryMock()
	commentService := NewCommentService(commentColQueryMock)
	commentId := "comment-1cab874c-d558-45a2-aaef-6487a415c261"
	comments := []bson.M{
		{"_id": "comment-" + uuid.NewString(),
			"post_id":      "post-" + uuid.NewString(),
			"user_id":      "user-" + uuid.NewString(),
			"comment":      "a comment",
			"like_count":   0,
			"created_date": primitive.NewDateTimeFromTime(time.Now()),
			"updated_date": primitive.NewDateTimeFromTime(time.Now())},
	}
	findCommentMock = func(ctx context.Context, i interface{}) (*[]bson.M, error) {
		return nil, errors.New("find to db return error")
	}
	if _, err := commentService.CheckIfCommentExist(commentId); err == nil {
		t.Error("If find to db return error than CheckIfCommentExist should also return error")
	}

	findCommentMock = func(ctx context.Context, i interface{}) (*[]bson.M, error) {
		return &[]bson.M{}, nil
	}
	isExist, err := commentService.CheckIfCommentExist(commentId)
	if err != nil {
		t.Error("If find to db does not return error than CheckIfCommentExist should also not return error")
	}
	if isExist {
		t.Error("If find to db does not return result than CheckIfCommentExist should not return true")
	}

	findCommentMock = func(ctx context.Context, i interface{}) (*[]bson.M, error) {
		return &comments, nil
	}
	isExist, _ = commentService.CheckIfCommentExist(commentId)
	if !isExist {
		t.Error("If find to db does return result than CheckIfCommentExist should return true")
	}
}

func TestGetCommentUserId(t *testing.T) {
	commentColQueryMock := newCommentCollectionQueryMock()
	commentService := NewCommentService(commentColQueryMock)
	foundComment := models.NewComment("comment-22e9eab8-0fe7-4a09-81a9-d0482ede1cf9", "post-4ca9098f-f646-46ff-8fbd-4aa14681fbfa", "user-2af4bb67-e06e-4c85-917e-80f95c140afe", "a new comment", 0, time.Now(), time.Now())

	findOneCommentMock = func(ctx context.Context, i interface{}) (*models.Comment, error) {
		return nil, errors.New("FindOne to db returns error")
	}
	_, err := commentService.GetCommentUserId(foundComment.Id)
	if err == nil {
		t.Error("if FindOne to db returns error than GetCommentUserId should also return error")
	}

	err = nil
	findOneCommentMock = func(ctx context.Context, i interface{}) (*models.Comment, error) {
		return foundComment, nil
	}
	foundCommentUserId, err := commentService.GetCommentUserId(foundComment.Id)
	if err != nil {
		t.Error("If FindOne to db does not return error than GetCommentUserId should also not return error")
	}
	if foundCommentUserId != foundComment.UserId {
		t.Error("GetCommentUserId should return the userid of the foundComment")
	}
}

func TestGetCommentLikeCount(t *testing.T) {
	commentColQueryMock := newCommentCollectionQueryMock()
	commentService := NewCommentService(commentColQueryMock)
	commentId := "comment-22e9eab8-0fe7-4a09-81a9-d0482ede1cf9"
	likes := []primitive.M{
		{"_id": "like-a0fd5081-f3ec-45d3-89c5-dd9873aca3fe", "user_id": "user-2af4bb67-e06e-4c85-917e-80f95c140afe", "resource_id": "comment-22e9eab8-0fe7-4a09-81a9-d0482ede1cf9", "resource_type": "comment"},
		{"_id": "like-a0fd5081-f3ec-45d3-89c5-dd9873aca3fe", "user_id": "user-2af4bb67-e06e-4c85-917e-80f95c140afe", "resource_id": "comment-22e9eab8-0fe7-4a09-81a9-d0482ede1cf9", "resource_type": "comment"},
		{"_id": "like-a0fd5081-f3ec-45d3-89c5-dd9873aca3fe", "user_id": "user-2af4bb67-e06e-4c85-917e-80f95c140afe", "resource_id": "comment-22e9eab8-0fe7-4a09-81a9-d0482ede1cf9", "resource_type": "comment"},
	}

	findCommentLikeMock = func(ctx context.Context, i interface{}) (*[]bson.M, error) {
		return nil, errors.New("find on db returns error")
	}
	_, err := commentService.getCommentLikeCount(commentId)
	if err == nil {
		t.Error("If Find on db returns error then getCommentLikeCount should also return error")
	}

	err = nil
	findCommentLikeMock = func(ctx context.Context, i interface{}) (*[]bson.M, error) {
		return &likes, nil
	}
	likeCount, err := commentService.getCommentLikeCount(commentId)
	if err != nil {
		t.Error("if Find on db does not return error then getCommentLiekCount should also not return error")
	}
	if likeCount != 3 {
		t.Error("getCommentLikeCount should return the correct amount of likes")
	}
}

func TestFindAllPostComment(t *testing.T) {
	commentColQueryMock := newCommentCollectionQueryMock()
	commentService := NewCommentService(commentColQueryMock)
	postId := "post-4ca9098f-f646-46ff-8fbd-4aa14681fbfa"
	comments := []bson.M{
		{"_id": "comment-" + uuid.NewString(),
			"post_id":      "post-4ca9098f-f646-46ff-8fbd-4aa14681fbfa",
			"user_id":      "user-" + uuid.NewString(),
			"comment":      "a new comment",
			"like_count":   0,
			"created_date": primitive.NewDateTimeFromTime(time.Now()),
			"updated_date": primitive.NewDateTimeFromTime(time.Now())},
		{"_id": "comment-" + uuid.NewString(),
			"post_id":      "post-4ca9098f-f646-46ff-8fbd-4aa14681fbfa",
			"user_id":      "user-" + uuid.NewString(),
			"comment":      "a new comment",
			"like_count":   0,
			"created_date": primitive.NewDateTimeFromTime(time.Now()),
			"updated_date": primitive.NewDateTimeFromTime(time.Now())},
		{"_id": "comment-" + uuid.NewString(),
			"post_id":      "post-4ca9098f-f646-46ff-8fbd-4aa14681fbfa",
			"user_id":      "user-" + uuid.NewString(),
			"comment":      "a new comment",
			"like_count":   0,
			"created_date": primitive.NewDateTimeFromTime(time.Now()),
			"updated_date": primitive.NewDateTimeFromTime(time.Now())},
	}

	findPostCommentMock = func(ctx context.Context, i interface{}) (*[]bson.M, error) {
		return nil, errors.New("Find on db return error")
	}
	_, err := commentService.FindAllPostComment(postId)
	if err == nil {
		t.Error("if Find on db return error then FindAllPostComment should also return error")
	}

	err = nil
	findPostCommentMock = func(ctx context.Context, i interface{}) (*[]bson.M, error) {
		return &comments, nil
	}
	findCommentLikeMock = func(ctx context.Context, i interface{}) (*[]bson.M, error) {
		return nil, errors.New("Find likes on db returns error")
	}
	_, err = commentService.FindAllPostComment(postId)
	if err == nil {
		t.Error("if Find Likes on db returns error then FindAllPostComment should also return error")
	}

	err = nil
	findPostCommentMock = func(ctx context.Context, i interface{}) (*[]bson.M, error) {
		return &comments, nil
	}
	findCommentLikeMock = func(ctx context.Context, i interface{}) (*[]bson.M, error) {
		return &[]bson.M{}, nil
	}
	foundComments, err := commentService.FindAllPostComment("post-4ca9098f-f646-46ff-8fbd-4aa14681fbfa")
	if err != nil {
		t.Error("if Find on db does not return error then FindAllPostComment should also not return error")
	}
	if len(foundComments) != 3 {
		t.Error("FindAllPostComment should return all found comments")
	}
}

func TestInsertComment(t *testing.T) {
	commentColQueryMock := newCommentCollectionQueryMock()
	commentService := NewCommentService(commentColQueryMock)
	insertedComment := models.NewComment("comment-"+uuid.NewString(),
		"post-4ca9098f-f646-46ff-8fbd-4aa14681fbfa",
		"user-"+uuid.NewString(),
		"a new comment", 0, time.Now(), time.Now())
	insertOneCommentMock = func(ctx context.Context, i interface{}) (*mongo.InsertOneResult, error) {
		return nil, errors.New("InsertOne on db return error")
	}
	if err := commentService.InsertComment(*insertedComment); err == nil {
		t.Error("if InsertOne on db return error than InsertComment should also return error")
	}

	insertOneCommentMock = func(ctx context.Context, i interface{}) (*mongo.InsertOneResult, error) {
		return nil, nil
	}
	if err := commentService.InsertComment(*insertedComment); err != nil {
		t.Error("If InsertOne on db does not return error than InsertComment should also not return error")
	}
}

func TestUpdateComment(t *testing.T) {
	commentColQueryMock := newCommentCollectionQueryMock()
	commentService := NewCommentService(commentColQueryMock)

	updateOneCommentMock = func(ctx context.Context, i1, i2 interface{}) (*mongo.UpdateResult, error) {
		return nil, errors.New("UpdateOne on db returns error")
	}
	if err := commentService.UpdateComment("comment-4ca9098f-f646-46ff-8fbd-4aa14681fbfa", "an updated comment"); err == nil {
		t.Error("if UpdateOne on db returns error then UpdateComment should also return error")
	}

	updateOneCommentMock = func(ctx context.Context, i1, i2 interface{}) (*mongo.UpdateResult, error) {
		return nil, nil
	}
	if err := commentService.UpdateComment("comment-4ca9098f-f646-46ff-8fbd-4aa14681fbfa", "an updated comment"); err != nil {
		t.Error("if UpdateOne on db does not returns error then UpdateComment should also not return error")
	}
}

func TestDeleteComment(t *testing.T) {
	commentColQueryMock := newCommentCollectionQueryMock()
	commentService := NewCommentService(commentColQueryMock)

	deleteOneCommentMock = func(ctx context.Context, i interface{}) (*mongo.DeleteResult, error) {
		return nil, errors.New("DeleteOne on db return error")
	}
	if err := commentService.DeleteComment("comment-4ca9098f-f646-46ff-8fbd-4aa14681fbfa"); err == nil {
		t.Error("If DeleteOne on db returns error then DeleteComment should also return error")
	}

	deleteOneCommentMock = func(ctx context.Context, i interface{}) (*mongo.DeleteResult, error) {
		return nil, nil
	}
	if err := commentService.DeleteComment("comment-4ca9098f-f646-46ff-8fbd-4aa14681fbfa"); err != nil {
		t.Error("If DeleteOne on db does not return error then DeleteComment should also not return error")
	}
}
