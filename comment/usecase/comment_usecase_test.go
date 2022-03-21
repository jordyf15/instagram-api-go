package usecase_test

import (
	"errors"
	"instagram-go/comment/usecase"
	"instagram-go/domain"
	"instagram-go/domain/mocks"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestCommentUsecaseSuite(t *testing.T) {
	suite.Run(t, new(CommentUsecaseSuite))
}

type CommentUsecaseSuite struct {
	suite.Suite
	commentRepository *mocks.CommentRepository
	postRepository    *mocks.PostRepository
	likeRepository    *mocks.LikeRepository
	headerHelper      *mocks.IHeaderHelper
}

func (cu *CommentUsecaseSuite) SetupTest() {
	cu.commentRepository = new(mocks.CommentRepository)
	cu.postRepository = new(mocks.PostRepository)
	cu.likeRepository = new(mocks.LikeRepository)
	cu.headerHelper = new(mocks.IHeaderHelper)
}

func (cu *CommentUsecaseSuite) TestFindCommentFindPostError() {
	cu.postRepository.On("FindPosts", mock.AnythingOfType("M")).Return(nil, errors.New("FindPosts return error"))

	commentUsecase := usecase.NewCommentUsecase(cu.commentRepository, cu.postRepository, cu.likeRepository, cu.headerHelper)
	_, err := commentUsecase.FindComments("postid1")

	expectedError := domain.ErrInternalServerError.Error()
	assert.EqualErrorf(cu.T(), err, expectedError, "Should have return %s but got %s", expectedError, err.Error())
}

func (cu *CommentUsecaseSuite) TestFindCommentPostNotFound() {
	cu.postRepository.On("FindPosts", mock.AnythingOfType("M")).Return(&[]bson.M{}, nil)

	commentUsecase := usecase.NewCommentUsecase(cu.commentRepository, cu.postRepository, cu.likeRepository, cu.headerHelper)
	_, err := commentUsecase.FindComments("postid1")

	expectedError := domain.ErrPostNotFound.Error()
	assert.EqualErrorf(cu.T(), err, expectedError, "Should have return %s but got %s", expectedError, err.Error())
}

func (cu *CommentUsecaseSuite) TestFindCommentFindCommentsError() {
	cu.postRepository.On("FindPosts", mock.AnythingOfType("M")).Return(&[]bson.M{
		{"_id": "postid1",
			"user_id":           "userid1",
			"visual_media_urls": []primitive.A{{"jpg.jpg"}, {"png.png"}},
			"caption":           "caption1",
			"created_date":      primitive.NewDateTimeFromTime(time.Now()),
			"updated_date":      primitive.NewDateTimeFromTime(time.Now())},
	}, nil)
	cu.commentRepository.On("FindComments", mock.AnythingOfType("M")).Return(nil, errors.New("FindComments return error"))

	commentUsecase := usecase.NewCommentUsecase(cu.commentRepository, cu.postRepository, cu.likeRepository, cu.headerHelper)
	_, err := commentUsecase.FindComments("postid1")

	expectedError := domain.ErrInternalServerError.Error()
	assert.EqualErrorf(cu.T(), err, expectedError, "Should have return %s but got %s", expectedError, err.Error())
}

func (cu *CommentUsecaseSuite) TestFindCommentFindLikesError() {
	cu.postRepository.On("FindPosts", mock.AnythingOfType("M")).Return(&[]bson.M{
		{"_id": "commentid1", "post_id": "postid1", "user_id": "userid1", "comment": "comment1",
			"created_date": primitive.NewDateTimeFromTime(time.Now()), "updated_date": primitive.NewDateTimeFromTime(time.Now())},
	}, nil)
	cu.commentRepository.On("FindComments", mock.AnythingOfType("M")).Return(&[]bson.M{
		{"_id": "commentid1", "post_id": "postid1", "user_id": "userid1", "comment": "comment1",
			"created_date": primitive.NewDateTimeFromTime(time.Now()), "updated_date": primitive.NewDateTimeFromTime(time.Now())},
		{"_id": "commentid2", "post_id": "postid1", "user_id": "userid1", "comment": "comment2",
			"created_date": primitive.NewDateTimeFromTime(time.Now()), "updated_date": primitive.NewDateTimeFromTime(time.Now())},
	}, nil)
	cu.likeRepository.On("FindLikes", mock.AnythingOfType("M")).Return(nil, errors.New("FindLikes return error"))

	commentUsecase := usecase.NewCommentUsecase(cu.commentRepository, cu.postRepository, cu.likeRepository, cu.headerHelper)
	_, err := commentUsecase.FindComments("postid1")

	expectedError := domain.ErrInternalServerError.Error()
	assert.EqualErrorf(cu.T(), err, expectedError, "Should have return error %s but got %s", expectedError, err.Error())
}

func (cu *CommentUsecaseSuite) TestFindCommentSuccessful() {
	cu.postRepository.On("FindPosts", mock.AnythingOfType("M")).Return(&[]bson.M{
		{"_id": "commentid1", "post_id": "postid1", "user_id": "userid1", "comment": "comment1",
			"created_date": primitive.NewDateTimeFromTime(time.Now()), "updated_date": primitive.NewDateTimeFromTime(time.Now())},
	}, nil)
	cu.commentRepository.On("FindComments", mock.AnythingOfType("M")).Return(&[]bson.M{
		{"_id": "commentid1", "post_id": "postid1", "user_id": "userid1", "comment": "comment1",
			"created_date": primitive.NewDateTimeFromTime(time.Now()), "updated_date": primitive.NewDateTimeFromTime(time.Now())},
		{"_id": "commentid2", "post_id": "postid1", "user_id": "userid1", "comment": "comment2",
			"created_date": primitive.NewDateTimeFromTime(time.Now()), "updated_date": primitive.NewDateTimeFromTime(time.Now())},
	}, nil)
	cu.likeRepository.On("FindLikes", mock.AnythingOfType("M")).Return(&[]bson.M{}, nil)

	commentUsecase := usecase.NewCommentUsecase(cu.commentRepository, cu.postRepository, cu.likeRepository, cu.headerHelper)
	comments, err := commentUsecase.FindComments("postid1")

	assert.NoErrorf(cu.T(), err, "Should not have return error but got %s", err)
	assert.Equal(cu.T(), 2, len(*comments), "Should have return 2 comments")
	assert.Equal(cu.T(), "commentid1", (*comments)[0].Id, "The first id should be correct")
	assert.Equal(cu.T(), "commentid2", (*comments)[1].Id, "The first id should be correct")
}

func (cu *CommentUsecaseSuite) TestPostCommentGetUserIdFromTokenError() {
	cu.headerHelper.On("GetUserIdFromToken", mock.AnythingOfType("string")).Return("", errors.New("GetUserIdFromToken return error"))
	comment := domain.NewComment("commentid1", "postid1", "userid1", "comment1", 0, time.Now(), time.Now())

	commentUsecase := usecase.NewCommentUsecase(cu.commentRepository, cu.postRepository, cu.likeRepository, cu.headerHelper)
	err := commentUsecase.PostComment(comment, "token1")

	expectedError := domain.ErrInternalServerError.Error()
	assert.EqualErrorf(cu.T(), err, expectedError, "Should have return error %s but got %s", expectedError, err.Error())
}

func (cu *CommentUsecaseSuite) TestPostCommentFindPostsError() {
	comment := domain.NewComment("commentid1", "postid1", "userid1", "comment1", 0, time.Now(), time.Now())
	cu.headerHelper.On("GetUserIdFromToken", mock.AnythingOfType("string")).Return("userid1", nil)
	cu.postRepository.On("FindPosts", mock.AnythingOfType("M")).Return(nil, errors.New("FindPosts return error"))

	commentUsecase := usecase.NewCommentUsecase(cu.commentRepository, cu.postRepository, cu.likeRepository, cu.headerHelper)
	err := commentUsecase.PostComment(comment, "token1")

	expectedError := domain.ErrInternalServerError.Error()
	assert.EqualErrorf(cu.T(), err, expectedError, "Should have return error %s but got %s", expectedError, err.Error())
}
func (cu *CommentUsecaseSuite) TestPostCommentPostNotFound() {
	comment := domain.NewComment("commentid1", "postid1", "userid1", "comment1", 0, time.Now(), time.Now())
	cu.headerHelper.On("GetUserIdFromToken", mock.AnythingOfType("string")).Return("userid1", nil)
	cu.postRepository.On("FindPosts", mock.AnythingOfType("M")).Return(&[]bson.M{}, nil)

	commentUsecase := usecase.NewCommentUsecase(cu.commentRepository, cu.postRepository, cu.likeRepository, cu.headerHelper)
	err := commentUsecase.PostComment(comment, "token1")

	expectedError := domain.ErrPostNotFound.Error()
	assert.EqualErrorf(cu.T(), err, expectedError, "Should have return error %s but got %s", expectedError, err.Error())
}

func (cu *CommentUsecaseSuite) TestPostCommentInsertCommentError() {
	comment := domain.NewComment("commentid1", "postid1", "userid1", "comment1", 0, time.Now(), time.Now())
	cu.headerHelper.On("GetUserIdFromToken", mock.AnythingOfType("string")).Return("userid1", nil)
	cu.postRepository.On("FindPosts", mock.AnythingOfType("M")).Return(&[]bson.M{
		{"_id": "postid1",
			"user_id":           "userid1",
			"visual_media_urls": []primitive.A{{"jpg.jpg"}, {"png.png"}},
			"caption":           "caption1",
			"created_date":      primitive.NewDateTimeFromTime(time.Now()),
			"updated_date":      primitive.NewDateTimeFromTime(time.Now())},
	}, nil)
	cu.commentRepository.On("InsertComment", mock.AnythingOfType("*domain.Comment")).Return(errors.New("InsertComment return error"))

	commentUsecase := usecase.NewCommentUsecase(cu.commentRepository, cu.postRepository, cu.likeRepository, cu.headerHelper)
	err := commentUsecase.PostComment(comment, "token1")

	expectedError := domain.ErrInternalServerError.Error()
	assert.EqualErrorf(cu.T(), err, expectedError, "Should have return error %s but got %s", expectedError, err.Error())
}

func (cu *CommentUsecaseSuite) TestPostCommentInsertCommentSuccessful() {
	comment := domain.NewComment("commentid1", "postid1", "userid1", "comment1", 0, time.Now(), time.Now())
	cu.headerHelper.On("GetUserIdFromToken", mock.AnythingOfType("string")).Return("userid1", nil)
	cu.postRepository.On("FindPosts", mock.AnythingOfType("M")).Return(&[]bson.M{
		{"_id": "postid1",
			"user_id":           "userid1",
			"visual_media_urls": []primitive.A{{"jpg.jpg"}, {"png.png"}},
			"caption":           "caption1",
			"created_date":      primitive.NewDateTimeFromTime(time.Now()),
			"updated_date":      primitive.NewDateTimeFromTime(time.Now())},
	}, nil)
	cu.commentRepository.On("InsertComment", mock.AnythingOfType("*domain.Comment")).Return(nil)

	commentUsecase := usecase.NewCommentUsecase(cu.commentRepository, cu.postRepository, cu.likeRepository, cu.headerHelper)
	err := commentUsecase.PostComment(comment, "token1")

	assert.NoErrorf(cu.T(), err, "Should have not return error but got %s", err)
}

func (cu *CommentUsecaseSuite) TestPutCommentGetUserIdFromTokenError() {
	comment := domain.NewComment("commentid1", "postid1", "userid1", "comment1", 0, time.Now(), time.Now())
	cu.headerHelper.On("GetUserIdFromToken", mock.AnythingOfType("string")).Return("", errors.New("GetUserIdFromToken return error"))

	commentUsecase := usecase.NewCommentUsecase(cu.commentRepository, cu.postRepository, cu.likeRepository, cu.headerHelper)
	err := commentUsecase.PutComment(comment, "token1")

	expectedError := domain.ErrInternalServerError.Error()
	assert.EqualErrorf(cu.T(), err, expectedError, "Should have return error %s but got %s", expectedError, err.Error())
}

func (cu *CommentUsecaseSuite) TestPutCommentFindCommentError() {
	comment := domain.NewComment("commentid1", "postid1", "userid1", "comment1", 0, time.Now(), time.Now())
	cu.headerHelper.On("GetUserIdFromToken", mock.AnythingOfType("string")).Return("userid1", nil)
	cu.commentRepository.On("FindComments", mock.AnythingOfType("M")).Return(nil, errors.New("FindComments return error"))

	commentUsecase := usecase.NewCommentUsecase(cu.commentRepository, cu.postRepository, cu.likeRepository, cu.headerHelper)
	err := commentUsecase.PutComment(comment, "token1")

	expectedError := domain.ErrInternalServerError.Error()
	assert.EqualErrorf(cu.T(), err, expectedError, "Should have return error %s but got %s", expectedError, err.Error())
}

func (cu *CommentUsecaseSuite) TestPutCommentCommentNotFound() {
	comment := domain.NewComment("commentid1", "postid1", "userid1", "comment1", 0, time.Now(), time.Now())
	cu.headerHelper.On("GetUserIdFromToken", mock.AnythingOfType("string")).Return("userid1", nil)
	cu.commentRepository.On("FindComments", mock.AnythingOfType("M")).Return(&[]bson.M{}, nil)

	commentUsecase := usecase.NewCommentUsecase(cu.commentRepository, cu.postRepository, cu.likeRepository, cu.headerHelper)
	err := commentUsecase.PutComment(comment, "token1")

	expectedError := domain.ErrCommentNotFound.Error()
	assert.EqualErrorf(cu.T(), err, expectedError, "Should have return error %s but got %s", expectedError, err.Error())
}
func (cu *CommentUsecaseSuite) TestPutCommentFindOneCommentError() {
	comment := domain.NewComment("commentid1", "postid1", "userid1", "comment1", 0, time.Now(), time.Now())
	cu.headerHelper.On("GetUserIdFromToken", mock.AnythingOfType("string")).Return("userid1", nil)
	cu.commentRepository.On("FindComments", mock.AnythingOfType("M")).Return(&[]bson.M{
		{"_id": "commentid1", "post_id": "postid1", "user_id": "userid1", "comment": "comment1",
			"created_date": primitive.NewDateTimeFromTime(time.Now()), "updated_date": primitive.NewDateTimeFromTime(time.Now())},
	}, nil)
	cu.commentRepository.On("FindOneComment", mock.AnythingOfType("string")).Return(nil, errors.New("FindOneComment return error"))

	commentUsecase := usecase.NewCommentUsecase(cu.commentRepository, cu.postRepository, cu.likeRepository, cu.headerHelper)
	err := commentUsecase.PutComment(comment, "token1")

	expectedError := domain.ErrInternalServerError.Error()
	assert.EqualErrorf(cu.T(), err, expectedError, "Should have return error %s but got %s", expectedError, err.Error())
}

func (cu *CommentUsecaseSuite) TestPutCommentUnauthorizedCommentUpdate() {
	comment := domain.NewComment("commentid1", "postid1", "userid1", "comment1", 0, time.Now(), time.Now())
	cu.headerHelper.On("GetUserIdFromToken", mock.AnythingOfType("string")).Return("userid1", nil)
	cu.commentRepository.On("FindComments", mock.AnythingOfType("M")).Return(&[]bson.M{
		{"_id": "commentid1", "post_id": "postid1", "user_id": "userid1", "comment": "comment1",
			"created_date": primitive.NewDateTimeFromTime(time.Now()), "updated_date": primitive.NewDateTimeFromTime(time.Now())},
	}, nil)
	cu.commentRepository.On("FindOneComment", mock.AnythingOfType("string")).Return(domain.NewComment(
		"commentid1", "postid1", "userid2", "comment1", 0, time.Now(), time.Now(),
	), nil)

	commentUsecase := usecase.NewCommentUsecase(cu.commentRepository, cu.postRepository, cu.likeRepository, cu.headerHelper)
	err := commentUsecase.PutComment(comment, "token1")

	expectedError := domain.ErrUnauthorizedCommentUpdate.Error()
	assert.EqualErrorf(cu.T(), err, expectedError, "Should have return error %s but got %s", expectedError, err.Error())
}

func (cu *CommentUsecaseSuite) TestPutCommentUpdateCommentError() {
	comment := domain.NewComment("commentid1", "postid1", "userid1", "comment1", 0, time.Now(), time.Now())
	cu.headerHelper.On("GetUserIdFromToken", mock.AnythingOfType("string")).Return("userid1", nil)
	cu.commentRepository.On("FindComments", mock.AnythingOfType("M")).Return(&[]bson.M{
		{"_id": "commentid1", "post_id": "postid1", "user_id": "userid1", "comment": "comment1",
			"created_date": primitive.NewDateTimeFromTime(time.Now()), "updated_date": primitive.NewDateTimeFromTime(time.Now())},
	}, nil)
	cu.commentRepository.On("FindOneComment", mock.AnythingOfType("string")).Return(domain.NewComment(
		"commentid1", "postid1", "userid1", "comment1", 0, time.Now(), time.Now(),
	), nil)
	cu.commentRepository.On("UpdateComment", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(errors.New("UpdateComment return error"))

	commentUsecase := usecase.NewCommentUsecase(cu.commentRepository, cu.postRepository, cu.likeRepository, cu.headerHelper)
	err := commentUsecase.PutComment(comment, "token1")

	expectedError := domain.ErrInternalServerError.Error()
	assert.EqualErrorf(cu.T(), err, expectedError, "Should have return error %s but got %s", expectedError, err.Error())
}

func (cu *CommentUsecaseSuite) TestPutCommentUpdateCommentSuccessful() {
	comment := domain.NewComment("commentid1", "postid1", "userid1", "comment1", 0, time.Now(), time.Now())
	cu.headerHelper.On("GetUserIdFromToken", mock.AnythingOfType("string")).Return("userid1", nil)
	cu.commentRepository.On("FindComments", mock.AnythingOfType("M")).Return(&[]bson.M{
		{"_id": "commentid1", "post_id": "postid1", "user_id": "userid1", "comment": "comment1",
			"created_date": primitive.NewDateTimeFromTime(time.Now()), "updated_date": primitive.NewDateTimeFromTime(time.Now())},
	}, nil)
	cu.commentRepository.On("FindOneComment", mock.AnythingOfType("string")).Return(domain.NewComment(
		"commentid1", "postid1", "userid1", "comment1", 0, time.Now(), time.Now(),
	), nil)
	cu.commentRepository.On("UpdateComment", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil)

	commentUsecase := usecase.NewCommentUsecase(cu.commentRepository, cu.postRepository, cu.likeRepository, cu.headerHelper)
	err := commentUsecase.PutComment(comment, "token1")

	assert.NoErrorf(cu.T(), err, "Should have not return error but got %s", err)
}

func (cu *CommentUsecaseSuite) TestDeleteCommentGetUserIdFromTokenError() {
	cu.headerHelper.On("GetUserIdFromToken", mock.AnythingOfType("string")).Return("", errors.New("GetUserIdFromToken return error"))

	commentUsecase := usecase.NewCommentUsecase(cu.commentRepository, cu.postRepository, cu.likeRepository, cu.headerHelper)
	err := commentUsecase.DeleteComment("commentid1", "token1")

	expectedError := domain.ErrInternalServerError.Error()
	assert.EqualErrorf(cu.T(), err, expectedError, "Should have return error %s but got %s", expectedError, err.Error())
}

func (cu *CommentUsecaseSuite) TestDeleteCommentFindCommentsError() {
	cu.headerHelper.On("GetUserIdFromToken", mock.AnythingOfType("string")).Return("userid1", nil)
	cu.commentRepository.On("FindComments", mock.AnythingOfType("M")).Return(nil, errors.New("FindComments return error"))

	commentUsecase := usecase.NewCommentUsecase(cu.commentRepository, cu.postRepository, cu.likeRepository, cu.headerHelper)
	err := commentUsecase.DeleteComment("commentid1", "token1")

	expectedError := domain.ErrInternalServerError.Error()
	assert.EqualErrorf(cu.T(), err, expectedError, "Should have return error %s but got %s", expectedError, err.Error())
}

func (cu *CommentUsecaseSuite) TestDeleteCommentCommentNotFound() {
	cu.headerHelper.On("GetUserIdFromToken", mock.AnythingOfType("string")).Return("userid1", nil)
	cu.commentRepository.On("FindComments", mock.AnythingOfType("M")).Return(&[]bson.M{}, nil)

	commentUsecase := usecase.NewCommentUsecase(cu.commentRepository, cu.postRepository, cu.likeRepository, cu.headerHelper)
	err := commentUsecase.DeleteComment("commentid1", "token1")

	expectedError := domain.ErrCommentNotFound.Error()
	assert.EqualErrorf(cu.T(), err, expectedError, "Should have return error %s but got %s", expectedError, err.Error())
}

func (cu *CommentUsecaseSuite) TestDeleteCommentFindOneCommentError() {
	cu.headerHelper.On("GetUserIdFromToken", mock.AnythingOfType("string")).Return("userid1", nil)
	cu.commentRepository.On("FindComments", mock.AnythingOfType("M")).Return(&[]bson.M{
		{"_id": "commentid1", "post_id": "postid1", "user_id": "userid1", "comment": "comment1",
			"created_date": primitive.NewDateTimeFromTime(time.Now()), "updated_date": primitive.NewDateTimeFromTime(time.Now())},
	}, nil)
	cu.commentRepository.On("FindOneComment", mock.AnythingOfType("string")).Return(nil, errors.New("FindOneComment return error"))

	commentUsecase := usecase.NewCommentUsecase(cu.commentRepository, cu.postRepository, cu.likeRepository, cu.headerHelper)
	err := commentUsecase.DeleteComment("commentid1", "token1")

	expectedError := domain.ErrInternalServerError.Error()
	assert.EqualErrorf(cu.T(), err, expectedError, "Should have return error %s but got %s", expectedError, err.Error())
}

func (cu *CommentUsecaseSuite) TestDeleteCommentUnauthorizedCommentDelete() {
	cu.headerHelper.On("GetUserIdFromToken", mock.AnythingOfType("string")).Return("userid1", nil)
	cu.commentRepository.On("FindComments", mock.AnythingOfType("M")).Return(&[]bson.M{
		{"_id": "commentid1", "post_id": "postid1", "user_id": "userid2", "comment": "comment1",
			"created_date": primitive.NewDateTimeFromTime(time.Now()), "updated_date": primitive.NewDateTimeFromTime(time.Now())},
	}, nil)
	cu.commentRepository.On("FindOneComment", mock.AnythingOfType("string")).Return(domain.NewComment(
		"commentid1", "postid1", "userid2", "comment1", 0, time.Now(), time.Now(),
	), nil)

	commentUsecase := usecase.NewCommentUsecase(cu.commentRepository, cu.postRepository, cu.likeRepository, cu.headerHelper)
	err := commentUsecase.DeleteComment("commentid1", "token1")

	expectedError := domain.ErrUnauthorizedCommentDelete.Error()
	assert.EqualErrorf(cu.T(), err, expectedError, "Should have return error %s but got %s", expectedError, err.Error())
}

func (cu *CommentUsecaseSuite) TestDeleteCommentDeleteCommentError() {
	cu.headerHelper.On("GetUserIdFromToken", mock.AnythingOfType("string")).Return("userid1", nil)
	cu.commentRepository.On("FindComments", mock.AnythingOfType("M")).Return(&[]bson.M{
		{"_id": "commentid1", "post_id": "postid1", "user_id": "userid1", "comment": "comment1",
			"created_date": primitive.NewDateTimeFromTime(time.Now()), "updated_date": primitive.NewDateTimeFromTime(time.Now())},
	}, nil)
	cu.commentRepository.On("FindOneComment", mock.AnythingOfType("string")).Return(domain.NewComment(
		"commentid1", "postid1", "userid1", "comment1", 0, time.Now(), time.Now(),
	), nil)
	cu.commentRepository.On("DeleteComment", mock.AnythingOfType("string")).Return(errors.New("DeleteComment return error"))

	commentUsecase := usecase.NewCommentUsecase(cu.commentRepository, cu.postRepository, cu.likeRepository, cu.headerHelper)
	err := commentUsecase.DeleteComment("commentid1", "token1")

	expectedError := domain.ErrInternalServerError.Error()
	assert.EqualErrorf(cu.T(), err, expectedError, "Should have return error %s but got %s", expectedError, err.Error())
}

func (cu *CommentUsecaseSuite) TestDeleteCommentSuccessful() {
	cu.headerHelper.On("GetUserIdFromToken", mock.AnythingOfType("string")).Return("userid1", nil)
	cu.commentRepository.On("FindComments", mock.AnythingOfType("M")).Return(&[]bson.M{
		{"_id": "commentid1", "post_id": "postid1", "user_id": "userid1", "comment": "comment1",
			"created_date": primitive.NewDateTimeFromTime(time.Now()), "updated_date": primitive.NewDateTimeFromTime(time.Now())},
	}, nil)
	cu.commentRepository.On("FindOneComment", mock.AnythingOfType("string")).Return(domain.NewComment(
		"commentid1", "postid1", "userid1", "comment1", 0, time.Now(), time.Now(),
	), nil)
	cu.commentRepository.On("DeleteComment", mock.AnythingOfType("string")).Return(nil)

	commentUsecase := usecase.NewCommentUsecase(cu.commentRepository, cu.postRepository, cu.likeRepository, cu.headerHelper)
	err := commentUsecase.DeleteComment("commentid1", "token1")

	assert.NoErrorf(cu.T(), err, "Should have not return error but got %s", err)
}
