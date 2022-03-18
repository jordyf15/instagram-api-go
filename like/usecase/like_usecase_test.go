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

func TestInsertPostLike(t *testing.T) {
	suite.Run(t, new(InsertPostLikeSuite))
}

func TestDeletePostLikeSuite(t *testing.T) {
	suite.Run(t, new(DeletePostLikeSuite))
}

func TestInsertCommentLikeSuite(t *testing.T) {
	suite.Run(t, new(InsertCommentLikeSuite))
}

func TestDeleteCommentLikeSuite(t *testing.T) {
	suite.Run(t, new(DeleteCommentLikeSuite))
}

type InsertPostLikeSuite struct {
	suite.Suite
	headerHelper      *mocks.IHeaderHelper
	postRepository    *mocks.PostRepository
	likeRepository    *mocks.LikeRepository
	commentRepository *mocks.CommentRepository
}

func (ipls *InsertPostLikeSuite) SetupTest() {
	ipls.headerHelper = new(mocks.IHeaderHelper)
	ipls.postRepository = new(mocks.PostRepository)
	ipls.likeRepository = new(mocks.LikeRepository)
	ipls.commentRepository = new(mocks.CommentRepository)
}

func (ipls *InsertPostLikeSuite) TestGetUserIdFromTokenError() {
	ipls.headerHelper.On("GetUserIdFromToken", mock.Anything).Return("", errors.New("GetUserIdFromToken return error"))

	likeUsecase := usecase.NewLikeUsecase(ipls.likeRepository, ipls.postRepository, ipls.commentRepository, ipls.headerHelper)
	err := likeUsecase.InsertPostLike("postid1", "token1")

	expectedError := domain.ErrInternalServerError.Error()
	assert.EqualErrorf(ipls.T(), err, expectedError, "Should have return %s but got %s", expectedError, err.Error())

}

func (ipls *InsertPostLikeSuite) TestFindPostsError() {
	ipls.headerHelper.On("GetUserIdFromToken", mock.Anything).Return("userid1", nil)
	ipls.postRepository.On("FindPosts", mock.Anything).Return(nil, errors.New(""))

	likeUsecase := usecase.NewLikeUsecase(ipls.likeRepository, ipls.postRepository, ipls.commentRepository, ipls.headerHelper)
	err := likeUsecase.InsertPostLike("postid1", "token1")

	expectedError := domain.ErrInternalServerError.Error()
	assert.EqualErrorf(ipls.T(), err, expectedError, "Should have return %s but got %s", expectedError, err.Error())
}

func (ipls *InsertPostLikeSuite) TestPostNotFound() {
	ipls.headerHelper.On("GetUserIdFromToken", mock.Anything).Return("userid1", nil)
	ipls.postRepository.On("FindPosts", mock.Anything).Return(&[]bson.M{}, nil)

	likeUsecase := usecase.NewLikeUsecase(ipls.likeRepository, ipls.postRepository, ipls.commentRepository, ipls.headerHelper)
	err := likeUsecase.InsertPostLike("postid1", "token1")

	expectedError := domain.ErrPostNotFound.Error()
	assert.EqualErrorf(ipls.T(), err, expectedError, "Should have return %s but got %s", expectedError, err.Error())
}

func (ipls *InsertPostLikeSuite) TestFindLikesError() {
	ipls.headerHelper.On("GetUserIdFromToken", mock.Anything).Return("userid1", nil)
	ipls.postRepository.On("FindPosts", mock.Anything).Return(&[]bson.M{
		{"_id": "postid1",
			"user_id":           "userid1",
			"visual_media_urls": []primitive.A{{"jpg.jpg"}, {"png.png"}},
			"caption":           "caption1",
			"created_date":      primitive.NewDateTimeFromTime(time.Now()),
			"updated_date":      primitive.NewDateTimeFromTime(time.Now())},
	}, nil)
	ipls.likeRepository.On("FindLikes", mock.Anything).Return(nil, errors.New("FindLikes return error"))

	likeUsecase := usecase.NewLikeUsecase(ipls.likeRepository, ipls.postRepository, ipls.commentRepository, ipls.headerHelper)
	err := likeUsecase.InsertPostLike("postid1", "token1")

	expectedError := domain.ErrInternalServerError.Error()
	assert.EqualErrorf(ipls.T(), err, expectedError, "Should have return %s but got %s", expectedError, err.Error())
}

func (ipls *InsertPostLikeSuite) TestPostLikeFound() {
	ipls.headerHelper.On("GetUserIdFromToken", mock.Anything).Return("userid1", nil)
	ipls.postRepository.On("FindPosts", mock.Anything).Return(&[]bson.M{
		{"_id": "postid1",
			"user_id":           "userid1",
			"visual_media_urls": []primitive.A{{"jpg.jpg"}, {"png.png"}},
			"caption":           "caption1",
			"created_date":      primitive.NewDateTimeFromTime(time.Now()),
			"updated_date":      primitive.NewDateTimeFromTime(time.Now())},
	}, nil)
	ipls.likeRepository.On("FindLikes", mock.Anything).Return(&[]bson.M{
		{"_id": "likeid1", "user_id": "userid1", "resource_id": "postid1", "resource_type": "post"},
	}, nil)

	likeUsecase := usecase.NewLikeUsecase(ipls.likeRepository, ipls.postRepository, ipls.commentRepository, ipls.headerHelper)
	err := likeUsecase.InsertPostLike("postid1", "token1")

	expectedError := domain.ErrPostLikeConflict.Error()
	assert.EqualErrorf(ipls.T(), err, expectedError, "Should have return %s but got %s", expectedError, err.Error())
}

func (ipls *InsertPostLikeSuite) TestInsertLikeError() {
	ipls.headerHelper.On("GetUserIdFromToken", mock.Anything).Return("userid1", nil)
	ipls.postRepository.On("FindPosts", mock.Anything).Return(&[]bson.M{
		{"_id": "postid1",
			"user_id":           "userid1",
			"visual_media_urls": []primitive.A{{"jpg.jpg"}, {"png.png"}},
			"caption":           "caption1",
			"created_date":      primitive.NewDateTimeFromTime(time.Now()),
			"updated_date":      primitive.NewDateTimeFromTime(time.Now())},
	}, nil)
	ipls.likeRepository.On("FindLikes", mock.Anything).Return(&[]bson.M{}, nil)
	ipls.likeRepository.On("InsertLike", mock.Anything).Return(errors.New("InsertLike return error"))

	likeUsecase := usecase.NewLikeUsecase(ipls.likeRepository, ipls.postRepository, ipls.commentRepository, ipls.headerHelper)
	err := likeUsecase.InsertPostLike("postid1", "token1")

	expectedError := domain.ErrInternalServerError.Error()
	assert.EqualErrorf(ipls.T(), err, expectedError, "Should have return %s but got %s", expectedError, err.Error())
}

func (ipls *InsertPostLikeSuite) TestInsertLikeSuccessful() {
	ipls.headerHelper.On("GetUserIdFromToken", mock.Anything).Return("userid1", nil)
	ipls.postRepository.On("FindPosts", mock.Anything).Return(&[]bson.M{
		{"_id": "postid1",
			"user_id":           "userid1",
			"visual_media_urls": []primitive.A{{"jpg.jpg"}, {"png.png"}},
			"caption":           "caption1",
			"created_date":      primitive.NewDateTimeFromTime(time.Now()),
			"updated_date":      primitive.NewDateTimeFromTime(time.Now())},
	}, nil)
	ipls.likeRepository.On("FindLikes", mock.Anything).Return(&[]bson.M{}, nil)
	ipls.likeRepository.On("InsertLike", mock.Anything).Return(nil)

	likeUsecase := usecase.NewLikeUsecase(ipls.likeRepository, ipls.postRepository, ipls.commentRepository, ipls.headerHelper)
	err := likeUsecase.InsertPostLike("postid1", "token1")

	assert.NoErrorf(ipls.T(), err, "Should have not return error but got %s", err)
}

type DeletePostLikeSuite struct {
	suite.Suite
	headerHelper      *mocks.IHeaderHelper
	postRepository    *mocks.PostRepository
	likeRepository    *mocks.LikeRepository
	commentRepository *mocks.CommentRepository
}

func (dpls *DeletePostLikeSuite) SetupTest() {
	dpls.headerHelper = new(mocks.IHeaderHelper)
	dpls.postRepository = new(mocks.PostRepository)
	dpls.likeRepository = new(mocks.LikeRepository)
	dpls.commentRepository = new(mocks.CommentRepository)
}

func (dpls *DeletePostLikeSuite) TestGetUserIdFromTokenError() {
	dpls.headerHelper.On("GetUserIdFromToken", mock.Anything).Return("", errors.New("GetUserIdFromToken return error"))

	likeUsecase := usecase.NewLikeUsecase(dpls.likeRepository, dpls.postRepository, dpls.commentRepository, dpls.headerHelper)
	err := likeUsecase.DeletePostLike("likeid1", "token1")

	expectedError := domain.ErrInternalServerError.Error()
	assert.EqualErrorf(dpls.T(), err, expectedError, "Should have return %s but got %s", expectedError, err.Error())
}
func (dpls *DeletePostLikeSuite) TestFindLikesError() {
	dpls.headerHelper.On("GetUserIdFromToken", mock.Anything).Return("userid1", nil)
	dpls.likeRepository.On("FindLikes", mock.Anything).Return(nil, errors.New("FindLikes return error"))

	likeUsecase := usecase.NewLikeUsecase(dpls.likeRepository, dpls.postRepository, dpls.commentRepository, dpls.headerHelper)
	err := likeUsecase.DeletePostLike("likeid1", "token1")

	expectedError := domain.ErrInternalServerError.Error()
	assert.EqualErrorf(dpls.T(), err, expectedError, "Should have return %s but got %s", expectedError, err.Error())
}

func (dpls *DeletePostLikeSuite) TestLikeNotFound() {
	dpls.headerHelper.On("GetUserIdFromToken", mock.Anything).Return("userid1", nil)
	dpls.likeRepository.On("FindLikes", mock.Anything).Return(&[]bson.M{}, nil)

	likeUsecase := usecase.NewLikeUsecase(dpls.likeRepository, dpls.postRepository, dpls.commentRepository, dpls.headerHelper)
	err := likeUsecase.DeletePostLike("likeid1", "token1")

	expectedError := domain.ErrLikeNotFound.Error()
	assert.EqualErrorf(dpls.T(), err, expectedError, "Should have return %s but got %s", expectedError, err.Error())
}

func (dpls *DeletePostLikeSuite) TestFindOneLikeError() {
	dpls.headerHelper.On("GetUserIdFromToken", mock.Anything).Return("userid1", nil)
	dpls.likeRepository.On("FindLikes", mock.Anything).Return(&[]bson.M{
		{"_id": "likeid1", "user_id": "userid1", "resource_id": "postid1", "resource_type": "post"},
	}, nil)
	dpls.likeRepository.On("FindOneLike", mock.Anything).Return(nil, errors.New("FindOneLike return error"))

	likeUsecase := usecase.NewLikeUsecase(dpls.likeRepository, dpls.postRepository, dpls.commentRepository, dpls.headerHelper)
	err := likeUsecase.DeletePostLike("likeid1", "token1")

	expectedError := domain.ErrInternalServerError.Error()
	assert.EqualErrorf(dpls.T(), err, expectedError, "Should have return %s but got %s", expectedError, err.Error())
}

func (dpls *DeletePostLikeSuite) TestUnauthorizedLikeDelete() {
	dpls.headerHelper.On("GetUserIdFromToken", mock.Anything).Return("userid1", nil)
	dpls.likeRepository.On("FindLikes", mock.Anything).Return(&[]bson.M{
		{"_id": "likeid1", "user_id": "userid2", "resource_id": "postid1", "resource_type": "post"},
	}, nil)
	dpls.likeRepository.On("FindOneLike", mock.Anything).Return(domain.NewLike(
		"likeid1", "userid2", "postid1", "post",
	), nil)

	likeUsecase := usecase.NewLikeUsecase(dpls.likeRepository, dpls.postRepository, dpls.commentRepository, dpls.headerHelper)
	err := likeUsecase.DeletePostLike("likeid1", "token1")

	expectedError := domain.ErrUnauthorizedLikeDelete.Error()
	assert.EqualErrorf(dpls.T(), err, expectedError, "Should have return %s but got %s", expectedError, err.Error())
}

func (dpls *DeletePostLikeSuite) TestDeleteLikeError() {
	dpls.headerHelper.On("GetUserIdFromToken", mock.Anything).Return("userid1", nil)
	dpls.likeRepository.On("FindLikes", mock.Anything).Return(&[]bson.M{
		{"_id": "likeid1", "user_id": "userid1", "resource_id": "postid1", "resource_type": "post"},
	}, nil)
	dpls.likeRepository.On("FindOneLike", mock.Anything).Return(domain.NewLike(
		"likeid1", "userid1", "postid1", "post",
	), nil)
	dpls.likeRepository.On("DeleteLike", mock.Anything).Return(errors.New("Delete like return error"))

	likeUsecase := usecase.NewLikeUsecase(dpls.likeRepository, dpls.postRepository, dpls.commentRepository, dpls.headerHelper)
	err := likeUsecase.DeletePostLike("likeid1", "token1")

	expectedError := domain.ErrInternalServerError.Error()
	assert.EqualErrorf(dpls.T(), err, expectedError, "Should have return %s but got %s", expectedError, err.Error())
}

func (dpls *DeletePostLikeSuite) TestDeleteLikeSuccessful() {
	dpls.headerHelper.On("GetUserIdFromToken", mock.Anything).Return("userid1", nil)
	dpls.likeRepository.On("FindLikes", mock.Anything).Return(&[]bson.M{
		{"_id": "likeid1", "user_id": "userid1", "resource_id": "postid1", "resource_type": "post"},
	}, nil)
	dpls.likeRepository.On("FindOneLike", mock.Anything).Return(domain.NewLike(
		"likeid1", "userid1", "postid1", "post",
	), nil)
	dpls.likeRepository.On("DeleteLike", mock.Anything).Return(nil)

	likeUsecase := usecase.NewLikeUsecase(dpls.likeRepository, dpls.postRepository, dpls.commentRepository, dpls.headerHelper)
	err := likeUsecase.DeletePostLike("likeid1", "token1")

	assert.NoErrorf(dpls.T(), err, "should have not return error but got %s", err)
}

type InsertCommentLikeSuite struct {
	suite.Suite
	commentRepository *mocks.CommentRepository
	likeRepository    *mocks.LikeRepository
	postRepository    *mocks.PostRepository
	headerHelper      *mocks.IHeaderHelper
}

func (icls *InsertCommentLikeSuite) SetupTest() {
	icls.commentRepository = new(mocks.CommentRepository)
	icls.likeRepository = new(mocks.LikeRepository)
	icls.headerHelper = new(mocks.IHeaderHelper)
}

func (icls *InsertCommentLikeSuite) TestGetUserIdFromTokenError() {
	icls.headerHelper.On("GetUserIdFromToken", mock.Anything).Return("", errors.New("GetUserIdFromToken return error"))

	likeUsecase := usecase.NewLikeUsecase(icls.likeRepository, icls.postRepository, icls.commentRepository, icls.headerHelper)
	err := likeUsecase.InsertCommentLike("likeid1", "token1")

	expectedError := domain.ErrInternalServerError.Error()
	assert.EqualErrorf(icls.T(), err, expectedError, "Should have return %s but got %s", expectedError, err.Error())
}

func (icls *InsertCommentLikeSuite) TestFindCommentError() {
	icls.headerHelper.On("GetUserIdFromToken", mock.Anything).Return("userid1", nil)
	icls.commentRepository.On("FindComments", mock.Anything).Return(nil, errors.New("FindComments return error"))

	likeUsecase := usecase.NewLikeUsecase(icls.likeRepository, icls.postRepository, icls.commentRepository, icls.headerHelper)
	err := likeUsecase.InsertCommentLike("likeid1", "token1")

	expectedError := domain.ErrInternalServerError.Error()
	assert.EqualErrorf(icls.T(), err, expectedError, "Should have return %s but got %s", expectedError, err.Error())
}

func (icls *InsertCommentLikeSuite) TestCommentNotFound() {
	icls.headerHelper.On("GetUserIdFromToken", mock.Anything).Return("userid1", nil)
	icls.commentRepository.On("FindComments", mock.Anything).Return(&[]bson.M{}, nil)

	likeUsecase := usecase.NewLikeUsecase(icls.likeRepository, icls.postRepository, icls.commentRepository, icls.headerHelper)
	err := likeUsecase.InsertCommentLike("likeid1", "token1")

	expectedError := domain.ErrCommentNotFound.Error()
	assert.EqualErrorf(icls.T(), err, expectedError, "Should have return %s but got %s", expectedError, err.Error())
}

func (icls *InsertCommentLikeSuite) TestFindLikesError() {
	icls.headerHelper.On("GetUserIdFromToken", mock.Anything).Return("userid1", nil)
	icls.commentRepository.On("FindComments", mock.Anything).Return(&[]bson.M{
		{"_id": "commentid1", "post_id": "postid1", "user_id": "userid1", "comment": "comment1",
			"created_date": primitive.NewDateTimeFromTime(time.Now()), "updated_date": primitive.NewDateTimeFromTime(time.Now())},
	}, nil)
	icls.likeRepository.On("FindLikes", mock.Anything).Return(nil, errors.New("FindLikes return error"))

	likeUsecase := usecase.NewLikeUsecase(icls.likeRepository, icls.postRepository, icls.commentRepository, icls.headerHelper)
	err := likeUsecase.InsertCommentLike("likeid1", "token1")

	expectedError := domain.ErrInternalServerError.Error()
	assert.EqualErrorf(icls.T(), err, expectedError, "Should have return %s but got %s", expectedError, err.Error())
}

func (icls *InsertCommentLikeSuite) TestCommentLikeFound() {
	icls.headerHelper.On("GetUserIdFromToken", mock.Anything).Return("userid1", nil)
	icls.commentRepository.On("FindComments", mock.Anything).Return(&[]bson.M{
		{"_id": "commentid1", "post_id": "postid1", "user_id": "userid1", "comment": "comment1",
			"created_date": primitive.NewDateTimeFromTime(time.Now()), "updated_date": primitive.NewDateTimeFromTime(time.Now())},
	}, nil)
	icls.likeRepository.On("FindLikes", mock.Anything).Return(&[]bson.M{
		{"_id": "likeid1", "user_id": "userid1", "resource_id": "commentid1", "resource_type": "comment"},
	}, nil)

	likeUsecase := usecase.NewLikeUsecase(icls.likeRepository, icls.postRepository, icls.commentRepository, icls.headerHelper)
	err := likeUsecase.InsertCommentLike("likeid1", "token1")

	expectedError := domain.ErrCommentLikeConflict.Error()
	assert.EqualErrorf(icls.T(), err, expectedError, "Should have return %s but got %s", expectedError, err.Error())
}

func (icls *InsertCommentLikeSuite) TestInsertLikeError() {
	icls.headerHelper.On("GetUserIdFromToken", mock.Anything).Return("userid1", nil)
	icls.commentRepository.On("FindComments", mock.Anything).Return(&[]bson.M{
		{"_id": "commentid1", "post_id": "postid1", "user_id": "userid1", "comment": "comment1",
			"created_date": primitive.NewDateTimeFromTime(time.Now()), "updated_date": primitive.NewDateTimeFromTime(time.Now())},
	}, nil)
	icls.likeRepository.On("FindLikes", mock.Anything).Return(&[]bson.M{}, nil)
	icls.likeRepository.On("InsertLike", mock.Anything).Return(errors.New("InsertLike return error"))

	likeUsecase := usecase.NewLikeUsecase(icls.likeRepository, icls.postRepository, icls.commentRepository, icls.headerHelper)
	err := likeUsecase.InsertCommentLike("likeid1", "token1")

	expectedError := domain.ErrInternalServerError.Error()
	assert.EqualErrorf(icls.T(), err, expectedError, "Should have return %s but got %s", expectedError, err.Error())
}

func (icls *InsertCommentLikeSuite) TestInsertLikeSuccessful() {
	icls.headerHelper.On("GetUserIdFromToken", mock.Anything).Return("userid1", nil)
	icls.commentRepository.On("FindComments", mock.Anything).Return(&[]bson.M{
		{"_id": "commentid1", "post_id": "postid1", "user_id": "userid1", "comment": "comment1",
			"created_date": primitive.NewDateTimeFromTime(time.Now()), "updated_date": primitive.NewDateTimeFromTime(time.Now())},
	}, nil)
	icls.likeRepository.On("FindLikes", mock.Anything).Return(&[]bson.M{}, nil)
	icls.likeRepository.On("InsertLike", mock.Anything).Return(nil)

	likeUsecase := usecase.NewLikeUsecase(icls.likeRepository, icls.postRepository, icls.commentRepository, icls.headerHelper)
	err := likeUsecase.InsertCommentLike("likeid1", "token1")

	assert.NoErrorf(icls.T(), err, "Should have not return error but got %s", err)
}

type DeleteCommentLikeSuite struct {
	suite.Suite
	commentRepository *mocks.CommentRepository
	likeRepository    *mocks.LikeRepository
	postRepository    *mocks.PostRepository
	headerHelper      *mocks.IHeaderHelper
}

func (dcls *DeleteCommentLikeSuite) SetupTest() {
	dcls.commentRepository = new(mocks.CommentRepository)
	dcls.likeRepository = new(mocks.LikeRepository)
	dcls.postRepository = new(mocks.PostRepository)
	dcls.headerHelper = new(mocks.IHeaderHelper)
}

func (dcls *DeleteCommentLikeSuite) TestGetUserIdFromTokenError() {
	dcls.headerHelper.On("GetUserIdFromToken", mock.Anything).Return("", errors.New("GetUserIdFromToken return error"))

	likeUsecase := usecase.NewLikeUsecase(dcls.likeRepository, dcls.postRepository, dcls.commentRepository, dcls.headerHelper)
	err := likeUsecase.DeleteCommentLike("likeid1", "token1")

	expectedError := domain.ErrInternalServerError.Error()
	assert.EqualErrorf(dcls.T(), err, expectedError, "Should have return %s but got %s", expectedError, err.Error())
}

func (dcls *DeleteCommentLikeSuite) TestFindLikesError() {
	dcls.headerHelper.On("GetUserIdFromToken", mock.Anything).Return("userid1", nil)
	dcls.likeRepository.On("FindLikes", mock.Anything).Return(nil, errors.New("FindLikes return error"))

	likeUsecase := usecase.NewLikeUsecase(dcls.likeRepository, dcls.postRepository, dcls.commentRepository, dcls.headerHelper)
	err := likeUsecase.DeleteCommentLike("likeid1", "token1")

	expectedError := domain.ErrInternalServerError.Error()
	assert.EqualErrorf(dcls.T(), err, expectedError, "Should have return %s but got %s", expectedError, err.Error())
}

func (dcls *DeleteCommentLikeSuite) TestLikeNotFound() {
	dcls.headerHelper.On("GetUserIdFromToken", mock.Anything).Return("userid1", nil)
	dcls.likeRepository.On("FindLikes", mock.Anything).Return(&[]bson.M{}, nil)

	likeUsecase := usecase.NewLikeUsecase(dcls.likeRepository, dcls.postRepository, dcls.commentRepository, dcls.headerHelper)
	err := likeUsecase.DeleteCommentLike("likeid1", "token1")

	expectedError := domain.ErrLikeNotFound.Error()
	assert.EqualErrorf(dcls.T(), err, expectedError, "Should have return %s but got %s", expectedError, err.Error())
}

func (dcls *DeleteCommentLikeSuite) TestFindOneLikeError() {
	dcls.headerHelper.On("GetUserIdFromToken", mock.Anything).Return("userid1", nil)
	dcls.likeRepository.On("FindLikes", mock.Anything).Return(&[]bson.M{
		{"_id": "likeid1", "user_id": "userid1", "resource_id": "commentid1", "resource_type": "comment"},
	}, nil)
	dcls.likeRepository.On("FindOneLike", mock.Anything).Return(nil, errors.New("FindOneLike return error"))

	likeUsecase := usecase.NewLikeUsecase(dcls.likeRepository, dcls.postRepository, dcls.commentRepository, dcls.headerHelper)
	err := likeUsecase.DeleteCommentLike("likeid1", "token1")

	expectedError := domain.ErrInternalServerError.Error()
	assert.EqualErrorf(dcls.T(), err, expectedError, "Should have return %s but got %s", expectedError, err.Error())
}

func (dcls *DeleteCommentLikeSuite) TestUnauthorizedLikeDelete() {
	dcls.headerHelper.On("GetUserIdFromToken", mock.Anything).Return("userid1", nil)
	dcls.likeRepository.On("FindLikes", mock.Anything).Return(&[]bson.M{
		{"_id": "likeid1", "user_id": "userid2", "resource_id": "commentid1", "resource_type": "comment"},
	}, nil)
	dcls.likeRepository.On("FindOneLike", mock.Anything).Return(domain.NewLike(
		"likeid1", "userid2", "commentid1", "comment",
	), nil)

	likeUsecase := usecase.NewLikeUsecase(dcls.likeRepository, dcls.postRepository, dcls.commentRepository, dcls.headerHelper)
	err := likeUsecase.DeleteCommentLike("likeid1", "token1")

	expectedError := domain.ErrUnauthorizedLikeDelete.Error()
	assert.EqualErrorf(dcls.T(), err, expectedError, "Should have return %s but got %s", expectedError, err.Error())
}

func (dcls *DeleteCommentLikeSuite) TestDeleteLikeError() {
	dcls.headerHelper.On("GetUserIdFromToken", mock.Anything).Return("userid1", nil)
	dcls.likeRepository.On("FindLikes", mock.Anything).Return(&[]bson.M{
		{"_id": "likeid1", "user_id": "userid1", "resource_id": "commentid1", "resource_type": "comment"},
	}, nil)
	dcls.likeRepository.On("FindOneLike", mock.Anything).Return(domain.NewLike(
		"likeid1", "userid1", "commentid1", "comment",
	), nil)
	dcls.likeRepository.On("DeleteLike", mock.Anything).Return(errors.New("DeleteLike return error"))

	likeUsecase := usecase.NewLikeUsecase(dcls.likeRepository, dcls.postRepository, dcls.commentRepository, dcls.headerHelper)
	err := likeUsecase.DeleteCommentLike("likeid1", "token1")

	expectedError := domain.ErrInternalServerError.Error()
	assert.EqualErrorf(dcls.T(), err, expectedError, "Should have return %s but got %s", expectedError, err.Error())
}

func (dcls *DeleteCommentLikeSuite) TestDeleteLikeSuccessful() {
	dcls.headerHelper.On("GetUserIdFromToken", mock.Anything).Return("userid1", nil)
	dcls.likeRepository.On("FindLikes", mock.Anything).Return(&[]bson.M{
		{"_id": "likeid1", "user_id": "userid1", "resource_id": "commentid1", "resource_type": "comment"},
	}, nil)
	dcls.likeRepository.On("FindOneLike", mock.Anything).Return(domain.NewLike(
		"likeid1", "userid1", "commentid1", "comment",
	), nil)
	dcls.likeRepository.On("DeleteLike", mock.Anything).Return(nil)

	likeUsecase := usecase.NewLikeUsecase(dcls.likeRepository, dcls.postRepository, dcls.commentRepository, dcls.headerHelper)
	err := likeUsecase.DeleteCommentLike("likeid1", "token1")

	assert.NoErrorf(dcls.T(), err, "Should have not return error but got %s", err)
}
