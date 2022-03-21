package usecase_test

import (
	"errors"
	"instagram-go/domain"
	"instagram-go/domain/mocks"
	"instagram-go/like/usecase"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestLikeUsecaseSuite(t *testing.T) {
	suite.Run(t, new(LikeUsecaseSuite))
}

type LikeUsecaseSuite struct {
	suite.Suite
	headerHelper      *mocks.IHeaderHelper
	postRepository    *mocks.PostRepository
	likeRepository    *mocks.LikeRepository
	commentRepository *mocks.CommentRepository
}

func (lu *LikeUsecaseSuite) SetupTest() {
	lu.headerHelper = new(mocks.IHeaderHelper)
	lu.postRepository = new(mocks.PostRepository)
	lu.likeRepository = new(mocks.LikeRepository)
	lu.commentRepository = new(mocks.CommentRepository)
}

func (lu *LikeUsecaseSuite) TestInsertPostLikeGetUserIdFromTokenError() {
	lu.headerHelper.On("GetUserIdFromToken", mock.AnythingOfType("string")).Return("", errors.New("GetUserIdFromToken return error"))

	likeUsecase := usecase.NewLikeUsecase(lu.likeRepository, lu.postRepository, lu.commentRepository, lu.headerHelper)
	err := likeUsecase.InsertPostLike("postid1", "token1")

	expectedError := domain.ErrInternalServerError.Error()
	assert.EqualErrorf(lu.T(), err, expectedError, "Should have return %s but got %s", expectedError, err.Error())

}

func (lu *LikeUsecaseSuite) TestInsertPostLikeFindPostsError() {
	lu.headerHelper.On("GetUserIdFromToken", mock.AnythingOfType("string")).Return("userid1", nil)
	lu.postRepository.On("FindPosts", mock.AnythingOfType("M")).Return(nil, errors.New(""))

	likeUsecase := usecase.NewLikeUsecase(lu.likeRepository, lu.postRepository, lu.commentRepository, lu.headerHelper)
	err := likeUsecase.InsertPostLike("postid1", "token1")

	expectedError := domain.ErrInternalServerError.Error()
	assert.EqualErrorf(lu.T(), err, expectedError, "Should have return %s but got %s", expectedError, err.Error())
}

func (lu *LikeUsecaseSuite) TestInsertPostLikePostNotFound() {
	lu.headerHelper.On("GetUserIdFromToken", mock.AnythingOfType("string")).Return("userid1", nil)
	lu.postRepository.On("FindPosts", mock.AnythingOfType("M")).Return(&[]bson.M{}, nil)

	likeUsecase := usecase.NewLikeUsecase(lu.likeRepository, lu.postRepository, lu.commentRepository, lu.headerHelper)
	err := likeUsecase.InsertPostLike("postid1", "token1")

	expectedError := domain.ErrPostNotFound.Error()
	assert.EqualErrorf(lu.T(), err, expectedError, "Should have return %s but got %s", expectedError, err.Error())
}

func (lu *LikeUsecaseSuite) TestInsertPostLikeFindLikesError() {
	lu.headerHelper.On("GetUserIdFromToken", mock.AnythingOfType("string")).Return("userid1", nil)
	lu.postRepository.On("FindPosts", mock.AnythingOfType("M")).Return(&[]bson.M{
		{"_id": "postid1",
			"user_id":           "userid1",
			"visual_media_urls": []primitive.A{{"jpg.jpg"}, {"png.png"}},
			"caption":           "caption1",
			"created_date":      primitive.NewDateTimeFromTime(time.Now()),
			"updated_date":      primitive.NewDateTimeFromTime(time.Now())},
	}, nil)
	lu.likeRepository.On("FindLikes", mock.AnythingOfType("M")).Return(nil, errors.New("FindLikes return error"))

	likeUsecase := usecase.NewLikeUsecase(lu.likeRepository, lu.postRepository, lu.commentRepository, lu.headerHelper)
	err := likeUsecase.InsertPostLike("postid1", "token1")

	expectedError := domain.ErrInternalServerError.Error()
	assert.EqualErrorf(lu.T(), err, expectedError, "Should have return %s but got %s", expectedError, err.Error())
}

func (lu *LikeUsecaseSuite) TestInsertPostLikePostLikeFound() {
	lu.headerHelper.On("GetUserIdFromToken", mock.AnythingOfType("string")).Return("userid1", nil)
	lu.postRepository.On("FindPosts", mock.AnythingOfType("M")).Return(&[]bson.M{
		{"_id": "postid1",
			"user_id":           "userid1",
			"visual_media_urls": []primitive.A{{"jpg.jpg"}, {"png.png"}},
			"caption":           "caption1",
			"created_date":      primitive.NewDateTimeFromTime(time.Now()),
			"updated_date":      primitive.NewDateTimeFromTime(time.Now())},
	}, nil)
	lu.likeRepository.On("FindLikes", mock.AnythingOfType("M")).Return(&[]bson.M{
		{"_id": "likeid1", "user_id": "userid1", "resource_id": "postid1", "resource_type": "post"},
	}, nil)

	likeUsecase := usecase.NewLikeUsecase(lu.likeRepository, lu.postRepository, lu.commentRepository, lu.headerHelper)
	err := likeUsecase.InsertPostLike("postid1", "token1")

	expectedError := domain.ErrPostLikeConflict.Error()
	assert.EqualErrorf(lu.T(), err, expectedError, "Should have return %s but got %s", expectedError, err.Error())
}

func (lu *LikeUsecaseSuite) TestInsertPostLikeInsertLikeError() {
	lu.headerHelper.On("GetUserIdFromToken", mock.AnythingOfType("string")).Return("userid1", nil)
	lu.postRepository.On("FindPosts", mock.AnythingOfType("M")).Return(&[]bson.M{
		{"_id": "postid1",
			"user_id":           "userid1",
			"visual_media_urls": []primitive.A{{"jpg.jpg"}, {"png.png"}},
			"caption":           "caption1",
			"created_date":      primitive.NewDateTimeFromTime(time.Now()),
			"updated_date":      primitive.NewDateTimeFromTime(time.Now())},
	}, nil)
	lu.likeRepository.On("FindLikes", mock.AnythingOfType("M")).Return(&[]bson.M{}, nil)
	lu.likeRepository.On("InsertLike", mock.AnythingOfType("*domain.Like")).Return(errors.New("InsertLike return error"))

	likeUsecase := usecase.NewLikeUsecase(lu.likeRepository, lu.postRepository, lu.commentRepository, lu.headerHelper)
	err := likeUsecase.InsertPostLike("postid1", "token1")

	expectedError := domain.ErrInternalServerError.Error()
	assert.EqualErrorf(lu.T(), err, expectedError, "Should have return %s but got %s", expectedError, err.Error())
}

func (lu *LikeUsecaseSuite) TestInsertPostLikeInsertLikeSuccessful() {
	lu.headerHelper.On("GetUserIdFromToken", mock.AnythingOfType("string")).Return("userid1", nil)
	lu.postRepository.On("FindPosts", mock.AnythingOfType("M")).Return(&[]bson.M{
		{"_id": "postid1",
			"user_id":           "userid1",
			"visual_media_urls": []primitive.A{{"jpg.jpg"}, {"png.png"}},
			"caption":           "caption1",
			"created_date":      primitive.NewDateTimeFromTime(time.Now()),
			"updated_date":      primitive.NewDateTimeFromTime(time.Now())},
	}, nil)
	lu.likeRepository.On("FindLikes", mock.AnythingOfType("M")).Return(&[]bson.M{}, nil)
	lu.likeRepository.On("InsertLike", mock.AnythingOfType("*domain.Like")).Return(nil)

	likeUsecase := usecase.NewLikeUsecase(lu.likeRepository, lu.postRepository, lu.commentRepository, lu.headerHelper)
	err := likeUsecase.InsertPostLike("postid1", "token1")

	assert.NoErrorf(lu.T(), err, "Should have not return error but got %s", err)
}

func (lu *LikeUsecaseSuite) TestDeletePostLikeGetUserIdFromTokenError() {
	lu.headerHelper.On("GetUserIdFromToken", mock.AnythingOfType("string")).Return("", errors.New("GetUserIdFromToken return error"))

	likeUsecase := usecase.NewLikeUsecase(lu.likeRepository, lu.postRepository, lu.commentRepository, lu.headerHelper)
	err := likeUsecase.DeletePostLike("likeid1", "token1")

	expectedError := domain.ErrInternalServerError.Error()
	assert.EqualErrorf(lu.T(), err, expectedError, "Should have return %s but got %s", expectedError, err.Error())
}
func (lu *LikeUsecaseSuite) TestDeletePostLikeFindLikesError() {
	lu.headerHelper.On("GetUserIdFromToken", mock.AnythingOfType("string")).Return("userid1", nil)
	lu.likeRepository.On("FindLikes", mock.AnythingOfType("M")).Return(nil, errors.New("FindLikes return error"))

	likeUsecase := usecase.NewLikeUsecase(lu.likeRepository, lu.postRepository, lu.commentRepository, lu.headerHelper)
	err := likeUsecase.DeletePostLike("likeid1", "token1")

	expectedError := domain.ErrInternalServerError.Error()
	assert.EqualErrorf(lu.T(), err, expectedError, "Should have return %s but got %s", expectedError, err.Error())
}

func (lu *LikeUsecaseSuite) TestDeletePostLikeLikeNotFound() {
	lu.headerHelper.On("GetUserIdFromToken", mock.AnythingOfType("string")).Return("userid1", nil)
	lu.likeRepository.On("FindLikes", mock.AnythingOfType("M")).Return(&[]bson.M{}, nil)

	likeUsecase := usecase.NewLikeUsecase(lu.likeRepository, lu.postRepository, lu.commentRepository, lu.headerHelper)
	err := likeUsecase.DeletePostLike("likeid1", "token1")

	expectedError := domain.ErrLikeNotFound.Error()
	assert.EqualErrorf(lu.T(), err, expectedError, "Should have return %s but got %s", expectedError, err.Error())
}

func (lu *LikeUsecaseSuite) TestDeletePostLikeFindOneLikeError() {
	lu.headerHelper.On("GetUserIdFromToken", mock.AnythingOfType("string")).Return("userid1", nil)
	lu.likeRepository.On("FindLikes", mock.AnythingOfType("M")).Return(&[]bson.M{
		{"_id": "likeid1", "user_id": "userid1", "resource_id": "postid1", "resource_type": "post"},
	}, nil)
	lu.likeRepository.On("FindOneLike", mock.AnythingOfType("string")).Return(nil, errors.New("FindOneLike return error"))

	likeUsecase := usecase.NewLikeUsecase(lu.likeRepository, lu.postRepository, lu.commentRepository, lu.headerHelper)
	err := likeUsecase.DeletePostLike("likeid1", "token1")

	expectedError := domain.ErrInternalServerError.Error()
	assert.EqualErrorf(lu.T(), err, expectedError, "Should have return %s but got %s", expectedError, err.Error())
}

func (lu *LikeUsecaseSuite) TestDeletePostLikeUnauthorizedLikeDelete() {
	lu.headerHelper.On("GetUserIdFromToken", mock.AnythingOfType("string")).Return("userid1", nil)
	lu.likeRepository.On("FindLikes", mock.AnythingOfType("M")).Return(&[]bson.M{
		{"_id": "likeid1", "user_id": "userid2", "resource_id": "postid1", "resource_type": "post"},
	}, nil)
	lu.likeRepository.On("FindOneLike", mock.AnythingOfType("string")).Return(domain.NewLike(
		"likeid1", "userid2", "postid1", "post",
	), nil)

	likeUsecase := usecase.NewLikeUsecase(lu.likeRepository, lu.postRepository, lu.commentRepository, lu.headerHelper)
	err := likeUsecase.DeletePostLike("likeid1", "token1")

	expectedError := domain.ErrUnauthorizedLikeDelete.Error()
	assert.EqualErrorf(lu.T(), err, expectedError, "Should have return %s but got %s", expectedError, err.Error())
}

func (lu *LikeUsecaseSuite) TestDeletePostLikeDeleteLikeError() {
	lu.headerHelper.On("GetUserIdFromToken", mock.AnythingOfType("string")).Return("userid1", nil)
	lu.likeRepository.On("FindLikes", mock.AnythingOfType("M")).Return(&[]bson.M{
		{"_id": "likeid1", "user_id": "userid1", "resource_id": "postid1", "resource_type": "post"},
	}, nil)
	lu.likeRepository.On("FindOneLike", mock.AnythingOfType("string")).Return(domain.NewLike(
		"likeid1", "userid1", "postid1", "post",
	), nil)
	lu.likeRepository.On("DeleteLike", mock.AnythingOfType("string")).Return(errors.New("Delete like return error"))

	likeUsecase := usecase.NewLikeUsecase(lu.likeRepository, lu.postRepository, lu.commentRepository, lu.headerHelper)
	err := likeUsecase.DeletePostLike("likeid1", "token1")

	expectedError := domain.ErrInternalServerError.Error()
	assert.EqualErrorf(lu.T(), err, expectedError, "Should have return %s but got %s", expectedError, err.Error())
}

func (lu *LikeUsecaseSuite) TestDeletePostLikeDeleteLikeSuccessful() {
	lu.headerHelper.On("GetUserIdFromToken", mock.AnythingOfType("string")).Return("userid1", nil)
	lu.likeRepository.On("FindLikes", mock.AnythingOfType("M")).Return(&[]bson.M{
		{"_id": "likeid1", "user_id": "userid1", "resource_id": "postid1", "resource_type": "post"},
	}, nil)
	lu.likeRepository.On("FindOneLike", mock.AnythingOfType("string")).Return(domain.NewLike(
		"likeid1", "userid1", "postid1", "post",
	), nil)
	lu.likeRepository.On("DeleteLike", mock.AnythingOfType("string")).Return(nil)

	likeUsecase := usecase.NewLikeUsecase(lu.likeRepository, lu.postRepository, lu.commentRepository, lu.headerHelper)
	err := likeUsecase.DeletePostLike("likeid1", "token1")

	assert.NoErrorf(lu.T(), err, "should have not return error but got %s", err)
}

func (lu *LikeUsecaseSuite) TestInsertCommentLikeGetUserIdFromTokenError() {
	lu.headerHelper.On("GetUserIdFromToken", mock.AnythingOfType("string")).Return("", errors.New("GetUserIdFromToken return error"))

	likeUsecase := usecase.NewLikeUsecase(lu.likeRepository, lu.postRepository, lu.commentRepository, lu.headerHelper)
	err := likeUsecase.InsertCommentLike("likeid1", "token1")

	expectedError := domain.ErrInternalServerError.Error()
	assert.EqualErrorf(lu.T(), err, expectedError, "Should have return %s but got %s", expectedError, err.Error())
}

func (lu *LikeUsecaseSuite) TestInsertCommentLikeFindCommentError() {
	lu.headerHelper.On("GetUserIdFromToken", mock.AnythingOfType("string")).Return("userid1", nil)
	lu.commentRepository.On("FindComments", mock.AnythingOfType("M")).Return(nil, errors.New("FindComments return error"))

	likeUsecase := usecase.NewLikeUsecase(lu.likeRepository, lu.postRepository, lu.commentRepository, lu.headerHelper)
	err := likeUsecase.InsertCommentLike("likeid1", "token1")

	expectedError := domain.ErrInternalServerError.Error()
	assert.EqualErrorf(lu.T(), err, expectedError, "Should have return %s but got %s", expectedError, err.Error())
}

func (lu *LikeUsecaseSuite) TestInsertCommentLikeCommentNotFound() {
	lu.headerHelper.On("GetUserIdFromToken", mock.AnythingOfType("string")).Return("userid1", nil)
	lu.commentRepository.On("FindComments", mock.AnythingOfType("M")).Return(&[]bson.M{}, nil)

	likeUsecase := usecase.NewLikeUsecase(lu.likeRepository, lu.postRepository, lu.commentRepository, lu.headerHelper)
	err := likeUsecase.InsertCommentLike("likeid1", "token1")

	expectedError := domain.ErrCommentNotFound.Error()
	assert.EqualErrorf(lu.T(), err, expectedError, "Should have return %s but got %s", expectedError, err.Error())
}

func (lu *LikeUsecaseSuite) TestInsertCommentLikeFindLikesError() {
	lu.headerHelper.On("GetUserIdFromToken", mock.AnythingOfType("string")).Return("userid1", nil)
	lu.commentRepository.On("FindComments", mock.AnythingOfType("M")).Return(&[]bson.M{
		{"_id": "commentid1", "post_id": "postid1", "user_id": "userid1", "comment": "comment1",
			"created_date": primitive.NewDateTimeFromTime(time.Now()), "updated_date": primitive.NewDateTimeFromTime(time.Now())},
	}, nil)
	lu.likeRepository.On("FindLikes", mock.AnythingOfType("M")).Return(nil, errors.New("FindLikes return error"))

	likeUsecase := usecase.NewLikeUsecase(lu.likeRepository, lu.postRepository, lu.commentRepository, lu.headerHelper)
	err := likeUsecase.InsertCommentLike("likeid1", "token1")

	expectedError := domain.ErrInternalServerError.Error()
	assert.EqualErrorf(lu.T(), err, expectedError, "Should have return %s but got %s", expectedError, err.Error())
}

func (lu *LikeUsecaseSuite) TestInsertCommentLikeCommentLikeFound() {
	lu.headerHelper.On("GetUserIdFromToken", mock.AnythingOfType("string")).Return("userid1", nil)
	lu.commentRepository.On("FindComments", mock.AnythingOfType("M")).Return(&[]bson.M{
		{"_id": "commentid1", "post_id": "postid1", "user_id": "userid1", "comment": "comment1",
			"created_date": primitive.NewDateTimeFromTime(time.Now()), "updated_date": primitive.NewDateTimeFromTime(time.Now())},
	}, nil)
	lu.likeRepository.On("FindLikes", mock.AnythingOfType("M")).Return(&[]bson.M{
		{"_id": "likeid1", "user_id": "userid1", "resource_id": "commentid1", "resource_type": "comment"},
	}, nil)

	likeUsecase := usecase.NewLikeUsecase(lu.likeRepository, lu.postRepository, lu.commentRepository, lu.headerHelper)
	err := likeUsecase.InsertCommentLike("likeid1", "token1")

	expectedError := domain.ErrCommentLikeConflict.Error()
	assert.EqualErrorf(lu.T(), err, expectedError, "Should have return %s but got %s", expectedError, err.Error())
}

func (lu *LikeUsecaseSuite) TestInsertCommentLikeInsertLikeError() {
	lu.headerHelper.On("GetUserIdFromToken", mock.AnythingOfType("string")).Return("userid1", nil)
	lu.commentRepository.On("FindComments", mock.AnythingOfType("M")).Return(&[]bson.M{
		{"_id": "commentid1", "post_id": "postid1", "user_id": "userid1", "comment": "comment1",
			"created_date": primitive.NewDateTimeFromTime(time.Now()), "updated_date": primitive.NewDateTimeFromTime(time.Now())},
	}, nil)
	lu.likeRepository.On("FindLikes", mock.AnythingOfType("M")).Return(&[]bson.M{}, nil)
	lu.likeRepository.On("InsertLike", mock.AnythingOfType("*domain.Like")).Return(errors.New("InsertLike return error"))

	likeUsecase := usecase.NewLikeUsecase(lu.likeRepository, lu.postRepository, lu.commentRepository, lu.headerHelper)
	err := likeUsecase.InsertCommentLike("likeid1", "token1")

	expectedError := domain.ErrInternalServerError.Error()
	assert.EqualErrorf(lu.T(), err, expectedError, "Should have return %s but got %s", expectedError, err.Error())
}

func (lu *LikeUsecaseSuite) TestInsertCommentLikeInsertLikeSuccessful() {
	lu.headerHelper.On("GetUserIdFromToken", mock.AnythingOfType("string")).Return("userid1", nil)
	lu.commentRepository.On("FindComments", mock.AnythingOfType("M")).Return(&[]bson.M{
		{"_id": "commentid1", "post_id": "postid1", "user_id": "userid1", "comment": "comment1",
			"created_date": primitive.NewDateTimeFromTime(time.Now()), "updated_date": primitive.NewDateTimeFromTime(time.Now())},
	}, nil)
	lu.likeRepository.On("FindLikes", mock.AnythingOfType("M")).Return(&[]bson.M{}, nil)
	lu.likeRepository.On("InsertLike", mock.AnythingOfType("*domain.Like")).Return(nil)

	likeUsecase := usecase.NewLikeUsecase(lu.likeRepository, lu.postRepository, lu.commentRepository, lu.headerHelper)
	err := likeUsecase.InsertCommentLike("likeid1", "token1")

	assert.NoErrorf(lu.T(), err, "Should have not return error but got %s", err)
}

func (lu *LikeUsecaseSuite) TestDeleteCommentLikeGetUserIdFromTokenError() {
	lu.headerHelper.On("GetUserIdFromToken", mock.AnythingOfType("string")).Return("", errors.New("GetUserIdFromToken return error"))

	likeUsecase := usecase.NewLikeUsecase(lu.likeRepository, lu.postRepository, lu.commentRepository, lu.headerHelper)
	err := likeUsecase.DeleteCommentLike("likeid1", "token1")

	expectedError := domain.ErrInternalServerError.Error()
	assert.EqualErrorf(lu.T(), err, expectedError, "Should have return %s but got %s", expectedError, err.Error())
}

func (lu *LikeUsecaseSuite) TestDeleteCommentLikeFindLikesError() {
	lu.headerHelper.On("GetUserIdFromToken", mock.AnythingOfType("string")).Return("userid1", nil)
	lu.likeRepository.On("FindLikes", mock.AnythingOfType("M")).Return(nil, errors.New("FindLikes return error"))

	likeUsecase := usecase.NewLikeUsecase(lu.likeRepository, lu.postRepository, lu.commentRepository, lu.headerHelper)
	err := likeUsecase.DeleteCommentLike("likeid1", "token1")

	expectedError := domain.ErrInternalServerError.Error()
	assert.EqualErrorf(lu.T(), err, expectedError, "Should have return %s but got %s", expectedError, err.Error())
}

func (lu *LikeUsecaseSuite) TestDeleteCommentLikeLikeNotFound() {
	lu.headerHelper.On("GetUserIdFromToken", mock.AnythingOfType("string")).Return("userid1", nil)
	lu.likeRepository.On("FindLikes", mock.AnythingOfType("M")).Return(&[]bson.M{}, nil)

	likeUsecase := usecase.NewLikeUsecase(lu.likeRepository, lu.postRepository, lu.commentRepository, lu.headerHelper)
	err := likeUsecase.DeleteCommentLike("likeid1", "token1")

	expectedError := domain.ErrLikeNotFound.Error()
	assert.EqualErrorf(lu.T(), err, expectedError, "Should have return %s but got %s", expectedError, err.Error())
}

func (lu *LikeUsecaseSuite) TestDeleteCommentLikeFindOneLikeError() {
	lu.headerHelper.On("GetUserIdFromToken", mock.AnythingOfType("string")).Return("userid1", nil)
	lu.likeRepository.On("FindLikes", mock.AnythingOfType("M")).Return(&[]bson.M{
		{"_id": "likeid1", "user_id": "userid1", "resource_id": "commentid1", "resource_type": "comment"},
	}, nil)
	lu.likeRepository.On("FindOneLike", mock.AnythingOfType("string")).Return(nil, errors.New("FindOneLike return error"))

	likeUsecase := usecase.NewLikeUsecase(lu.likeRepository, lu.postRepository, lu.commentRepository, lu.headerHelper)
	err := likeUsecase.DeleteCommentLike("likeid1", "token1")

	expectedError := domain.ErrInternalServerError.Error()
	assert.EqualErrorf(lu.T(), err, expectedError, "Should have return %s but got %s", expectedError, err.Error())
}

func (lu *LikeUsecaseSuite) TestDeleteCommentLikeUnauthorizedLikeDelete() {
	lu.headerHelper.On("GetUserIdFromToken", mock.AnythingOfType("string")).Return("userid1", nil)
	lu.likeRepository.On("FindLikes", mock.AnythingOfType("M")).Return(&[]bson.M{
		{"_id": "likeid1", "user_id": "userid2", "resource_id": "commentid1", "resource_type": "comment"},
	}, nil)
	lu.likeRepository.On("FindOneLike", mock.AnythingOfType("string")).Return(domain.NewLike(
		"likeid1", "userid2", "commentid1", "comment",
	), nil)

	likeUsecase := usecase.NewLikeUsecase(lu.likeRepository, lu.postRepository, lu.commentRepository, lu.headerHelper)
	err := likeUsecase.DeleteCommentLike("likeid1", "token1")

	expectedError := domain.ErrUnauthorizedLikeDelete.Error()
	assert.EqualErrorf(lu.T(), err, expectedError, "Should have return %s but got %s", expectedError, err.Error())
}

func (lu *LikeUsecaseSuite) TestDeleteCommentLikeDeleteLikeError() {
	lu.headerHelper.On("GetUserIdFromToken", mock.AnythingOfType("string")).Return("userid1", nil)
	lu.likeRepository.On("FindLikes", mock.AnythingOfType("M")).Return(&[]bson.M{
		{"_id": "likeid1", "user_id": "userid1", "resource_id": "commentid1", "resource_type": "comment"},
	}, nil)
	lu.likeRepository.On("FindOneLike", mock.AnythingOfType("string")).Return(domain.NewLike(
		"likeid1", "userid1", "commentid1", "comment",
	), nil)
	lu.likeRepository.On("DeleteLike", mock.AnythingOfType("string")).Return(errors.New("DeleteLike return error"))

	likeUsecase := usecase.NewLikeUsecase(lu.likeRepository, lu.postRepository, lu.commentRepository, lu.headerHelper)
	err := likeUsecase.DeleteCommentLike("likeid1", "token1")

	expectedError := domain.ErrInternalServerError.Error()
	assert.EqualErrorf(lu.T(), err, expectedError, "Should have return %s but got %s", expectedError, err.Error())
}

func (lu *LikeUsecaseSuite) TestDeleteCommentLikeDeleteLikeSuccessful() {
	lu.headerHelper.On("GetUserIdFromToken", mock.AnythingOfType("string")).Return("userid1", nil)
	lu.likeRepository.On("FindLikes", mock.AnythingOfType("M")).Return(&[]bson.M{
		{"_id": "likeid1", "user_id": "userid1", "resource_id": "commentid1", "resource_type": "comment"},
	}, nil)
	lu.likeRepository.On("FindOneLike", mock.AnythingOfType("string")).Return(domain.NewLike(
		"likeid1", "userid1", "commentid1", "comment",
	), nil)
	lu.likeRepository.On("DeleteLike", mock.AnythingOfType("string")).Return(nil)

	likeUsecase := usecase.NewLikeUsecase(lu.likeRepository, lu.postRepository, lu.commentRepository, lu.headerHelper)
	err := likeUsecase.DeleteCommentLike("likeid1", "token1")

	assert.NoErrorf(lu.T(), err, "Should have not return error but got %s", err)
}
