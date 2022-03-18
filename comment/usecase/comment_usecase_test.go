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

func TestFindCommentSuite(t *testing.T) {
	suite.Run(t, new(FindCommentSuite))
}

func TestPostCommentSuite(t *testing.T) {
	suite.Run(t, new(PostCommentSuite))
}

func TestPutCommentSuite(t *testing.T) {
	suite.Run(t, new(PutCommentSuite))
}

func TestDeleteCommentSuite(t *testing.T) {
	suite.Run(t, new(DeleteCommentSuite))
}

type FindCommentSuite struct {
	suite.Suite
	commentRepository *mocks.CommentRepository
	postRepository    *mocks.PostRepository
	likeRepository    *mocks.LikeRepository
	headerHelper      *mocks.IHeaderHelper
}

func (fcs *FindCommentSuite) SetupTest() {
	fcs.commentRepository = new(mocks.CommentRepository)
	fcs.postRepository = new(mocks.PostRepository)
	fcs.likeRepository = new(mocks.LikeRepository)
	fcs.headerHelper = new(mocks.IHeaderHelper)
}

func (fcs *FindCommentSuite) TestFindPostError() {
	fcs.postRepository.On("FindPosts", mock.Anything).Return(nil, errors.New("FindPosts return error"))

	commentUsecase := usecase.NewCommentUsecase(fcs.commentRepository, fcs.postRepository, fcs.likeRepository, fcs.headerHelper)
	_, err := commentUsecase.FindComments("postid1")

	expectedError := domain.ErrInternalServerError.Error()
	assert.EqualErrorf(fcs.T(), err, expectedError, "Should have return %s but got %s", expectedError, err.Error())
}

func (fcs *FindCommentSuite) TestPostNotFound() {
	fcs.postRepository.On("FindPosts", mock.Anything).Return(&[]bson.M{}, nil)

	commentUsecase := usecase.NewCommentUsecase(fcs.commentRepository, fcs.postRepository, fcs.likeRepository, fcs.headerHelper)
	_, err := commentUsecase.FindComments("postid1")

	expectedError := domain.ErrPostNotFound.Error()
	assert.EqualErrorf(fcs.T(), err, expectedError, "Should have return %s but got %s", expectedError, err.Error())
}

func (fcs *FindCommentSuite) TestFindCommentError() {
	fcs.postRepository.On("FindPosts", mock.Anything).Return(&[]bson.M{
		{"_id": "postid1",
			"user_id":           "userid1",
			"visual_media_urls": []primitive.A{{"jpg.jpg"}, {"png.png"}},
			"caption":           "caption1",
			"created_date":      primitive.NewDateTimeFromTime(time.Now()),
			"updated_date":      primitive.NewDateTimeFromTime(time.Now())},
	}, nil)
	fcs.commentRepository.On("FindComments", mock.Anything).Return(nil, errors.New("FindComments return error"))

	commentUsecase := usecase.NewCommentUsecase(fcs.commentRepository, fcs.postRepository, fcs.likeRepository, fcs.headerHelper)
	_, err := commentUsecase.FindComments("postid1")

	expectedError := domain.ErrInternalServerError.Error()
	assert.EqualErrorf(fcs.T(), err, expectedError, "Should have return %s but got %s", expectedError, err.Error())
}

func (fcs *FindCommentSuite) TestFindLikesError() {
	fcs.postRepository.On("FindPosts", mock.Anything).Return(&[]bson.M{
		{"_id": "commentid1", "post_id": "postid1", "user_id": "userid1", "comment": "comment1",
			"created_date": primitive.NewDateTimeFromTime(time.Now()), "updated_date": primitive.NewDateTimeFromTime(time.Now())},
	}, nil)
	fcs.commentRepository.On("FindComments", mock.Anything).Return(&[]bson.M{
		{"_id": "commentid1", "post_id": "postid1", "user_id": "userid1", "comment": "comment1",
			"created_date": primitive.NewDateTimeFromTime(time.Now()), "updated_date": primitive.NewDateTimeFromTime(time.Now())},
		{"_id": "commentid2", "post_id": "postid1", "user_id": "userid1", "comment": "comment2",
			"created_date": primitive.NewDateTimeFromTime(time.Now()), "updated_date": primitive.NewDateTimeFromTime(time.Now())},
	}, nil)
	fcs.likeRepository.On("FindLikes", mock.Anything).Return(nil, errors.New("FindLikes return error"))

	commentUsecase := usecase.NewCommentUsecase(fcs.commentRepository, fcs.postRepository, fcs.likeRepository, fcs.headerHelper)
	_, err := commentUsecase.FindComments("postid1")

	expectedError := domain.ErrInternalServerError.Error()
	assert.EqualErrorf(fcs.T(), err, expectedError, "Should have return error %s but got %s", expectedError, err.Error())
}

func (fcs *FindCommentSuite) TestFindCommentSuccessful() {
	fcs.postRepository.On("FindPosts", mock.Anything).Return(&[]bson.M{
		{"_id": "commentid1", "post_id": "postid1", "user_id": "userid1", "comment": "comment1",
			"created_date": primitive.NewDateTimeFromTime(time.Now()), "updated_date": primitive.NewDateTimeFromTime(time.Now())},
	}, nil)
	fcs.commentRepository.On("FindComments", mock.Anything).Return(&[]bson.M{
		{"_id": "commentid1", "post_id": "postid1", "user_id": "userid1", "comment": "comment1",
			"created_date": primitive.NewDateTimeFromTime(time.Now()), "updated_date": primitive.NewDateTimeFromTime(time.Now())},
		{"_id": "commentid2", "post_id": "postid1", "user_id": "userid1", "comment": "comment2",
			"created_date": primitive.NewDateTimeFromTime(time.Now()), "updated_date": primitive.NewDateTimeFromTime(time.Now())},
	}, nil)
	fcs.likeRepository.On("FindLikes", mock.Anything).Return(&[]bson.M{}, nil)

	commentUsecase := usecase.NewCommentUsecase(fcs.commentRepository, fcs.postRepository, fcs.likeRepository, fcs.headerHelper)
	comments, err := commentUsecase.FindComments("postid1")

	assert.NoErrorf(fcs.T(), err, "Should not have return error but got %s", err)
	assert.Equal(fcs.T(), 2, len(*comments), "Should have return 2 comments")
	assert.Equal(fcs.T(), "commentid1", (*comments)[0].Id, "The first id should be correct")
	assert.Equal(fcs.T(), "commentid2", (*comments)[1].Id, "The first id should be correct")
}

type PostCommentSuite struct {
	suite.Suite
	commentRepository *mocks.CommentRepository
	postRepository    *mocks.PostRepository
	likeRepository    *mocks.LikeRepository
	headerHelper      *mocks.IHeaderHelper
}

func (pcs *PostCommentSuite) SetupTest() {
	pcs.commentRepository = new(mocks.CommentRepository)
	pcs.postRepository = new(mocks.PostRepository)
	pcs.likeRepository = new(mocks.LikeRepository)
	pcs.headerHelper = new(mocks.IHeaderHelper)
}

func (pcs *PostCommentSuite) TestGetUserIdFromTokenError() {
	pcs.headerHelper.On("GetUserIdFromToken", mock.Anything).Return("", errors.New("GetUserIdFromToken return error"))
	comment := domain.NewComment("commentid1", "postid1", "userid1", "comment1", 0, time.Now(), time.Now())

	commentUsecase := usecase.NewCommentUsecase(pcs.commentRepository, pcs.postRepository, pcs.likeRepository, pcs.headerHelper)
	err := commentUsecase.PostComment(comment, "token1")

	expectedError := domain.ErrInternalServerError.Error()
	assert.EqualErrorf(pcs.T(), err, expectedError, "Should have return error %s but got %s", expectedError, err.Error())
}

func (pcs *PostCommentSuite) TestFindPostsError() {
	comment := domain.NewComment("commentid1", "postid1", "userid1", "comment1", 0, time.Now(), time.Now())
	pcs.headerHelper.On("GetUserIdFromToken", mock.Anything).Return("userid1", nil)
	pcs.postRepository.On("FindPosts", mock.Anything).Return(nil, errors.New("FindPosts return error"))

	commentUsecase := usecase.NewCommentUsecase(pcs.commentRepository, pcs.postRepository, pcs.likeRepository, pcs.headerHelper)
	err := commentUsecase.PostComment(comment, "token1")

	expectedError := domain.ErrInternalServerError.Error()
	assert.EqualErrorf(pcs.T(), err, expectedError, "Should have return error %s but got %s", expectedError, err.Error())
}
func (pcs *PostCommentSuite) TestPostNotFound() {
	comment := domain.NewComment("commentid1", "postid1", "userid1", "comment1", 0, time.Now(), time.Now())
	pcs.headerHelper.On("GetUserIdFromToken", mock.Anything).Return("userid1", nil)
	pcs.postRepository.On("FindPosts", mock.Anything).Return(&[]bson.M{}, nil)

	commentUsecase := usecase.NewCommentUsecase(pcs.commentRepository, pcs.postRepository, pcs.likeRepository, pcs.headerHelper)
	err := commentUsecase.PostComment(comment, "token1")

	expectedError := domain.ErrPostNotFound.Error()
	assert.EqualErrorf(pcs.T(), err, expectedError, "Should have return error %s but got %s", expectedError, err.Error())
}

func (pcs *PostCommentSuite) TestInsertCommentError() {
	comment := domain.NewComment("commentid1", "postid1", "userid1", "comment1", 0, time.Now(), time.Now())
	pcs.headerHelper.On("GetUserIdFromToken", mock.Anything).Return("userid1", nil)
	pcs.postRepository.On("FindPosts", mock.Anything).Return(&[]bson.M{
		{"_id": "postid1",
			"user_id":           "userid1",
			"visual_media_urls": []primitive.A{{"jpg.jpg"}, {"png.png"}},
			"caption":           "caption1",
			"created_date":      primitive.NewDateTimeFromTime(time.Now()),
			"updated_date":      primitive.NewDateTimeFromTime(time.Now())},
	}, nil)
	pcs.commentRepository.On("InsertComment", mock.Anything).Return(errors.New("InsertComment return error"))

	commentUsecase := usecase.NewCommentUsecase(pcs.commentRepository, pcs.postRepository, pcs.likeRepository, pcs.headerHelper)
	err := commentUsecase.PostComment(comment, "token1")

	expectedError := domain.ErrInternalServerError.Error()
	assert.EqualErrorf(pcs.T(), err, expectedError, "Should have return error %s but got %s", expectedError, err.Error())
}

func (pcs *PostCommentSuite) TestInsertCommentSuccessful() {
	comment := domain.NewComment("commentid1", "postid1", "userid1", "comment1", 0, time.Now(), time.Now())
	pcs.headerHelper.On("GetUserIdFromToken", mock.Anything).Return("userid1", nil)
	pcs.postRepository.On("FindPosts", mock.Anything).Return(&[]bson.M{
		{"_id": "postid1",
			"user_id":           "userid1",
			"visual_media_urls": []primitive.A{{"jpg.jpg"}, {"png.png"}},
			"caption":           "caption1",
			"created_date":      primitive.NewDateTimeFromTime(time.Now()),
			"updated_date":      primitive.NewDateTimeFromTime(time.Now())},
	}, nil)
	pcs.commentRepository.On("InsertComment", mock.Anything).Return(nil)

	commentUsecase := usecase.NewCommentUsecase(pcs.commentRepository, pcs.postRepository, pcs.likeRepository, pcs.headerHelper)
	err := commentUsecase.PostComment(comment, "token1")

	assert.NoErrorf(pcs.T(), err, "Should have not return error but got %s", err)
}

type PutCommentSuite struct {
	suite.Suite
	commentRepository *mocks.CommentRepository
	postRepository    *mocks.PostRepository
	likeRepository    *mocks.LikeRepository
	headerHelper      *mocks.IHeaderHelper
}

func (pcs *PutCommentSuite) SetupTest() {
	pcs.commentRepository = new(mocks.CommentRepository)
	pcs.postRepository = new(mocks.PostRepository)
	pcs.likeRepository = new(mocks.LikeRepository)
	pcs.headerHelper = new(mocks.IHeaderHelper)
}

func (pcs *PutCommentSuite) TestGetUserIdFromTokenError() {
	comment := domain.NewComment("commentid1", "postid1", "userid1", "comment1", 0, time.Now(), time.Now())
	pcs.headerHelper.On("GetUserIdFromToken", mock.Anything).Return("", errors.New("GetUserIdFromToken return error"))

	commentUsecase := usecase.NewCommentUsecase(pcs.commentRepository, pcs.postRepository, pcs.likeRepository, pcs.headerHelper)
	err := commentUsecase.PutComment(comment, "token1")

	expectedError := domain.ErrInternalServerError.Error()
	assert.EqualErrorf(pcs.T(), err, expectedError, "Should have return error %s but got %s", expectedError, err.Error())
}

func (pcs *PutCommentSuite) TestFindCommentError() {
	comment := domain.NewComment("commentid1", "postid1", "userid1", "comment1", 0, time.Now(), time.Now())
	pcs.headerHelper.On("GetUserIdFromToken", mock.Anything).Return("userid1", nil)
	pcs.commentRepository.On("FindComments", mock.Anything).Return(nil, errors.New("FindComments return error"))

	commentUsecase := usecase.NewCommentUsecase(pcs.commentRepository, pcs.postRepository, pcs.likeRepository, pcs.headerHelper)
	err := commentUsecase.PutComment(comment, "token1")

	expectedError := domain.ErrInternalServerError.Error()
	assert.EqualErrorf(pcs.T(), err, expectedError, "Should have return error %s but got %s", expectedError, err.Error())
}

func (pcs *PutCommentSuite) TestCommentNotFound() {
	comment := domain.NewComment("commentid1", "postid1", "userid1", "comment1", 0, time.Now(), time.Now())
	pcs.headerHelper.On("GetUserIdFromToken", mock.Anything).Return("userid1", nil)
	pcs.commentRepository.On("FindComments", mock.Anything).Return(&[]bson.M{}, nil)

	commentUsecase := usecase.NewCommentUsecase(pcs.commentRepository, pcs.postRepository, pcs.likeRepository, pcs.headerHelper)
	err := commentUsecase.PutComment(comment, "token1")

	expectedError := domain.ErrCommentNotFound.Error()
	assert.EqualErrorf(pcs.T(), err, expectedError, "Should have return error %s but got %s", expectedError, err.Error())
}
func (pcs *PutCommentSuite) TestFindOneCommentError() {
	comment := domain.NewComment("commentid1", "postid1", "userid1", "comment1", 0, time.Now(), time.Now())
	pcs.headerHelper.On("GetUserIdFromToken", mock.Anything).Return("userid1", nil)
	pcs.commentRepository.On("FindComments", mock.Anything).Return(&[]bson.M{
		{"_id": "commentid1", "post_id": "postid1", "user_id": "userid1", "comment": "comment1",
			"created_date": primitive.NewDateTimeFromTime(time.Now()), "updated_date": primitive.NewDateTimeFromTime(time.Now())},
	}, nil)
	pcs.commentRepository.On("FindOneComment", mock.Anything).Return(nil, errors.New("FindOneComment return error"))

	commentUsecase := usecase.NewCommentUsecase(pcs.commentRepository, pcs.postRepository, pcs.likeRepository, pcs.headerHelper)
	err := commentUsecase.PutComment(comment, "token1")

	expectedError := domain.ErrInternalServerError.Error()
	assert.EqualErrorf(pcs.T(), err, expectedError, "Should have return error %s but got %s", expectedError, err.Error())
}

func (pcs *PutCommentSuite) TestUnauthorizedCommentUpdate() {
	comment := domain.NewComment("commentid1", "postid1", "userid1", "comment1", 0, time.Now(), time.Now())
	pcs.headerHelper.On("GetUserIdFromToken", mock.Anything).Return("userid1", nil)
	pcs.commentRepository.On("FindComments", mock.Anything).Return(&[]bson.M{
		{"_id": "commentid1", "post_id": "postid1", "user_id": "userid1", "comment": "comment1",
			"created_date": primitive.NewDateTimeFromTime(time.Now()), "updated_date": primitive.NewDateTimeFromTime(time.Now())},
	}, nil)
	pcs.commentRepository.On("FindOneComment", mock.Anything).Return(domain.NewComment(
		"commentid1", "postid1", "userid2", "comment1", 0, time.Now(), time.Now(),
	), nil)

	commentUsecase := usecase.NewCommentUsecase(pcs.commentRepository, pcs.postRepository, pcs.likeRepository, pcs.headerHelper)
	err := commentUsecase.PutComment(comment, "token1")

	expectedError := domain.ErrUnauthorizedCommentUpdate.Error()
	assert.EqualErrorf(pcs.T(), err, expectedError, "Should have return error %s but got %s", expectedError, err.Error())
}

func (pcs *PutCommentSuite) TestUpdateCommentError() {
	comment := domain.NewComment("commentid1", "postid1", "userid1", "comment1", 0, time.Now(), time.Now())
	pcs.headerHelper.On("GetUserIdFromToken", mock.Anything).Return("userid1", nil)
	pcs.commentRepository.On("FindComments", mock.Anything).Return(&[]bson.M{
		{"_id": "commentid1", "post_id": "postid1", "user_id": "userid1", "comment": "comment1",
			"created_date": primitive.NewDateTimeFromTime(time.Now()), "updated_date": primitive.NewDateTimeFromTime(time.Now())},
	}, nil)
	pcs.commentRepository.On("FindOneComment", mock.Anything).Return(domain.NewComment(
		"commentid1", "postid1", "userid1", "comment1", 0, time.Now(), time.Now(),
	), nil)
	pcs.commentRepository.On("UpdateComment", mock.Anything, mock.Anything).Return(errors.New("UpdateComment return error"))

	commentUsecase := usecase.NewCommentUsecase(pcs.commentRepository, pcs.postRepository, pcs.likeRepository, pcs.headerHelper)
	err := commentUsecase.PutComment(comment, "token1")

	expectedError := domain.ErrInternalServerError.Error()
	assert.EqualErrorf(pcs.T(), err, expectedError, "Should have return error %s but got %s", expectedError, err.Error())
}

func (pcs *PutCommentSuite) TestUpdateCommentSuccessful() {
	comment := domain.NewComment("commentid1", "postid1", "userid1", "comment1", 0, time.Now(), time.Now())
	pcs.headerHelper.On("GetUserIdFromToken", mock.Anything).Return("userid1", nil)
	pcs.commentRepository.On("FindComments", mock.Anything).Return(&[]bson.M{
		{"_id": "commentid1", "post_id": "postid1", "user_id": "userid1", "comment": "comment1",
			"created_date": primitive.NewDateTimeFromTime(time.Now()), "updated_date": primitive.NewDateTimeFromTime(time.Now())},
	}, nil)
	pcs.commentRepository.On("FindOneComment", mock.Anything).Return(domain.NewComment(
		"commentid1", "postid1", "userid1", "comment1", 0, time.Now(), time.Now(),
	), nil)
	pcs.commentRepository.On("UpdateComment", mock.Anything, mock.Anything).Return(nil)

	commentUsecase := usecase.NewCommentUsecase(pcs.commentRepository, pcs.postRepository, pcs.likeRepository, pcs.headerHelper)
	err := commentUsecase.PutComment(comment, "token1")

	assert.NoErrorf(pcs.T(), err, "Should have not return error but got %s", err)
}

type DeleteCommentSuite struct {
	suite.Suite
	commentRepository *mocks.CommentRepository
	postRepository    *mocks.PostRepository
	likeRepository    *mocks.LikeRepository
	headerHelper      *mocks.IHeaderHelper
}

func (dcs *DeleteCommentSuite) SetupTest() {
	dcs.commentRepository = new(mocks.CommentRepository)
	dcs.postRepository = new(mocks.PostRepository)
	dcs.likeRepository = new(mocks.LikeRepository)
	dcs.headerHelper = new(mocks.IHeaderHelper)
}

func (dcs *DeleteCommentSuite) TestGetUserIdFromTokenError() {
	dcs.headerHelper.On("GetUserIdFromToken", mock.Anything).Return("", errors.New("GetUserIdFromToken return error"))

	commentUsecase := usecase.NewCommentUsecase(dcs.commentRepository, dcs.postRepository, dcs.likeRepository, dcs.headerHelper)
	err := commentUsecase.DeleteComment("commentid1", "token1")

	expectedError := domain.ErrInternalServerError.Error()
	assert.EqualErrorf(dcs.T(), err, expectedError, "Should have return error %s but got %s", expectedError, err.Error())
}

func (dcs *DeleteCommentSuite) TestFindCommentsError() {
	dcs.headerHelper.On("GetUserIdFromToken", mock.Anything).Return("userid1", nil)
	dcs.commentRepository.On("FindComments", mock.Anything).Return(nil, errors.New("FindComments return error"))

	commentUsecase := usecase.NewCommentUsecase(dcs.commentRepository, dcs.postRepository, dcs.likeRepository, dcs.headerHelper)
	err := commentUsecase.DeleteComment("commentid1", "token1")

	expectedError := domain.ErrInternalServerError.Error()
	assert.EqualErrorf(dcs.T(), err, expectedError, "Should have return error %s but got %s", expectedError, err.Error())
}

func (dcs *DeleteCommentSuite) TestCommentNotFound() {
	dcs.headerHelper.On("GetUserIdFromToken", mock.Anything).Return("userid1", nil)
	dcs.commentRepository.On("FindComments", mock.Anything).Return(&[]bson.M{}, nil)

	commentUsecase := usecase.NewCommentUsecase(dcs.commentRepository, dcs.postRepository, dcs.likeRepository, dcs.headerHelper)
	err := commentUsecase.DeleteComment("commentid1", "token1")

	expectedError := domain.ErrCommentNotFound.Error()
	assert.EqualErrorf(dcs.T(), err, expectedError, "Should have return error %s but got %s", expectedError, err.Error())
}

func (dcs *DeleteCommentSuite) TestFindOneCommentError() {
	dcs.headerHelper.On("GetUserIdFromToken", mock.Anything).Return("userid1", nil)
	dcs.commentRepository.On("FindComments", mock.Anything).Return(&[]bson.M{
		{"_id": "commentid1", "post_id": "postid1", "user_id": "userid1", "comment": "comment1",
			"created_date": primitive.NewDateTimeFromTime(time.Now()), "updated_date": primitive.NewDateTimeFromTime(time.Now())},
	}, nil)
	dcs.commentRepository.On("FindOneComment", mock.Anything).Return(nil, errors.New("FindOneComment return error"))

	commentUsecase := usecase.NewCommentUsecase(dcs.commentRepository, dcs.postRepository, dcs.likeRepository, dcs.headerHelper)
	err := commentUsecase.DeleteComment("commentid1", "token1")

	expectedError := domain.ErrInternalServerError.Error()
	assert.EqualErrorf(dcs.T(), err, expectedError, "Should have return error %s but got %s", expectedError, err.Error())
}

func (dcs *DeleteCommentSuite) TestUnauthorizedCommentDelete() {
	dcs.headerHelper.On("GetUserIdFromToken", mock.Anything).Return("userid1", nil)
	dcs.commentRepository.On("FindComments", mock.Anything).Return(&[]bson.M{
		{"_id": "commentid1", "post_id": "postid1", "user_id": "userid2", "comment": "comment1",
			"created_date": primitive.NewDateTimeFromTime(time.Now()), "updated_date": primitive.NewDateTimeFromTime(time.Now())},
	}, nil)
	dcs.commentRepository.On("FindOneComment", mock.Anything).Return(domain.NewComment(
		"commentid1", "postid1", "userid2", "comment1", 0, time.Now(), time.Now(),
	), nil)

	commentUsecase := usecase.NewCommentUsecase(dcs.commentRepository, dcs.postRepository, dcs.likeRepository, dcs.headerHelper)
	err := commentUsecase.DeleteComment("commentid1", "token1")

	expectedError := domain.ErrUnauthorizedCommentDelete.Error()
	assert.EqualErrorf(dcs.T(), err, expectedError, "Should have return error %s but got %s", expectedError, err.Error())
}

func (dcs *DeleteCommentSuite) TestDeleteCommentError() {
	dcs.headerHelper.On("GetUserIdFromToken", mock.Anything).Return("userid1", nil)
	dcs.commentRepository.On("FindComments", mock.Anything).Return(&[]bson.M{
		{"_id": "commentid1", "post_id": "postid1", "user_id": "userid1", "comment": "comment1",
			"created_date": primitive.NewDateTimeFromTime(time.Now()), "updated_date": primitive.NewDateTimeFromTime(time.Now())},
	}, nil)
	dcs.commentRepository.On("FindOneComment", mock.Anything).Return(domain.NewComment(
		"commentid1", "postid1", "userid1", "comment1", 0, time.Now(), time.Now(),
	), nil)
	dcs.commentRepository.On("DeleteComment", mock.Anything).Return(errors.New("DeleteComment return error"))

	commentUsecase := usecase.NewCommentUsecase(dcs.commentRepository, dcs.postRepository, dcs.likeRepository, dcs.headerHelper)
	err := commentUsecase.DeleteComment("commentid1", "token1")

	expectedError := domain.ErrInternalServerError.Error()
	assert.EqualErrorf(dcs.T(), err, expectedError, "Should have return error %s but got %s", expectedError, err.Error())
}

func (dcs *DeleteCommentSuite) TestDeleteCommentSuccessful() {
	dcs.headerHelper.On("GetUserIdFromToken", mock.Anything).Return("userid1", nil)
	dcs.commentRepository.On("FindComments", mock.Anything).Return(&[]bson.M{
		{"_id": "commentid1", "post_id": "postid1", "user_id": "userid1", "comment": "comment1",
			"created_date": primitive.NewDateTimeFromTime(time.Now()), "updated_date": primitive.NewDateTimeFromTime(time.Now())},
	}, nil)
	dcs.commentRepository.On("FindOneComment", mock.Anything).Return(domain.NewComment(
		"commentid1", "postid1", "userid1", "comment1", 0, time.Now(), time.Now(),
	), nil)
	dcs.commentRepository.On("DeleteComment", mock.Anything).Return(nil)

	commentUsecase := usecase.NewCommentUsecase(dcs.commentRepository, dcs.postRepository, dcs.likeRepository, dcs.headerHelper)
	err := commentUsecase.DeleteComment("commentid1", "token1")

	assert.NoErrorf(dcs.T(), err, "Should have not return error but got %s", err)
}
