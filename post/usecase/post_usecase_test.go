package usecase_test

import (
	"errors"
	"instagram-go/domain"
	"instagram-go/domain/mocks"
	"instagram-go/post/usecase"
	"mime/multipart"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestInsertPost(t *testing.T) {
	suite.Run(t, new(InsertPostSuite))
}
func TestFindPosts(t *testing.T) {
	suite.Run(t, new(FindPostsSuite))
}

func TestUpdatePostSuite(t *testing.T) {
	suite.Run(t, new(UpdatePostSuite))
}

func TestDeletePostSuite(t *testing.T) {
	suite.Run(t, new(DeletePostSuite))
}

type InsertPostSuite struct {
	suite.Suite
	mockPostRepository *mocks.PostRepository
	mockLikeRepository *mocks.LikeRepository
	mockFileOsHelper   *mocks.IFileOsHelper
	mockHeaderHelper   *mocks.IHeaderHelper
}

func (ips *InsertPostSuite) SetupTest() {
	ips.mockPostRepository = new(mocks.PostRepository)
	ips.mockLikeRepository = new(mocks.LikeRepository)
	ips.mockFileOsHelper = new(mocks.IFileOsHelper)
	ips.mockHeaderHelper = new(mocks.IHeaderHelper)
}

func (ips *InsertPostSuite) TestGetUserIdTokenError() {
	ips.mockHeaderHelper.On("GetUserIdFromToken", mock.Anything).Return("", errors.New("GetUserIdFromToken return error"))

	postUsecase := usecase.NewPostUseCase(ips.mockPostRepository, ips.mockLikeRepository, ips.mockHeaderHelper, ips.mockFileOsHelper)
	newPost := domain.NewPost("postid1", "userid1", []string{}, "caption1", 0, time.Now(), time.Now())
	err := postUsecase.InsertPost(newPost, "accessToken", []*multipart.FileHeader{})

	expectedError := domain.ErrInternalServerError.Error()
	assert.EqualErrorf(ips.T(), err, expectedError, "Should have return %s but got %s", expectedError, err)
}

func (ips *InsertPostSuite) TestMkDirAllError() {
	ips.mockHeaderHelper.On("GetUserIdFromToken", mock.Anything).Return("userId1", nil)
	ips.mockFileOsHelper.On("MkDirAll", mock.Anything, mock.Anything).Return(errors.New("MkDirAll return error"))

	postUsecase := usecase.NewPostUseCase(ips.mockPostRepository, ips.mockLikeRepository, ips.mockHeaderHelper, ips.mockFileOsHelper)
	newPost := domain.NewPost("postid1", "userid1", []string{}, "caption1", 0, time.Now(), time.Now())
	err := postUsecase.InsertPost(newPost, "accessToken", []*multipart.FileHeader{})

	expectedError := domain.ErrInternalServerError.Error()
	assert.EqualErrorf(ips.T(), err, expectedError, "Should have return %s but got %s", expectedError, err)
}

func (ips *InsertPostSuite) TestInsertPostError() {
	ips.mockHeaderHelper.On("GetUserIdFromToken", mock.Anything).Return("userId1", nil)
	ips.mockFileOsHelper.On("MkDirAll", mock.Anything, mock.Anything).Return(nil)
	ips.mockPostRepository.On("InsertPost", mock.Anything).Return(errors.New("InsertPost return error"))

	postUsecase := usecase.NewPostUseCase(ips.mockPostRepository, ips.mockLikeRepository, ips.mockHeaderHelper, ips.mockFileOsHelper)
	newPost := domain.NewPost("postid1", "userid1", []string{}, "caption1", 0, time.Now(), time.Now())
	err := postUsecase.InsertPost(newPost, "accessToken", []*multipart.FileHeader{})

	expectedError := domain.ErrInternalServerError.Error()
	assert.EqualErrorf(ips.T(), err, expectedError, "Should have return %s but got %s", expectedError, err)
}

func (ips *InsertPostSuite) TestInsertPostSuccessful() {
	ips.mockHeaderHelper.On("GetUserIdFromToken", mock.Anything).Return("userId1", nil)
	ips.mockFileOsHelper.On("MkDirAll", mock.Anything, mock.Anything).Return(nil)
	ips.mockPostRepository.On("InsertPost", mock.Anything).Return(nil)

	postUsecase := usecase.NewPostUseCase(ips.mockPostRepository, ips.mockLikeRepository, ips.mockHeaderHelper, ips.mockFileOsHelper)
	newPost := domain.NewPost("postid1", "userid1", []string{}, "caption1", 0, time.Now(), time.Now())
	err := postUsecase.InsertPost(newPost, "accessToken", []*multipart.FileHeader{})

	assert.NoErrorf(ips.T(), err, "should have not returned error but got %s", err)
}

type FindPostsSuite struct {
	suite.Suite
	mockPostRepository *mocks.PostRepository
	mockLikeRepository *mocks.LikeRepository
	mockFileOsHelper   *mocks.IFileOsHelper
	mockHeaderHelper   *mocks.IHeaderHelper
}

func (fps *FindPostsSuite) SetupTest() {
	fps.mockPostRepository = new(mocks.PostRepository)
	fps.mockFileOsHelper = new(mocks.IFileOsHelper)
	fps.mockHeaderHelper = new(mocks.IHeaderHelper)
	fps.mockLikeRepository = new(mocks.LikeRepository)
}

func (fps *FindPostsSuite) TestFindPostsError() {
	fps.mockPostRepository.On("FindPosts", mock.Anything).Return(nil, errors.New("FindPosts return error"))

	postUsecase := usecase.NewPostUseCase(fps.mockPostRepository, fps.mockLikeRepository, fps.mockHeaderHelper, fps.mockFileOsHelper)
	_, err := postUsecase.FindPosts()

	expectedError := domain.ErrInternalServerError.Error()
	assert.EqualErrorf(fps.T(), err, expectedError, "should have return error %s but got %s", expectedError, err)
}

func (fps *FindPostsSuite) TestFindLikesError() {
	fps.mockPostRepository.On("FindPosts", mock.Anything).Return(&[]bson.M{
		{"_id": "postid1",
			"user_id":           "userid1",
			"visual_media_urls": []primitive.A{{"jpg.jpg"}, {"png.png"}},
			"caption":           "caption1",
			"created_date":      primitive.NewDateTimeFromTime(time.Now()),
			"updated_date":      primitive.NewDateTimeFromTime(time.Now())},
	}, nil)
	fps.mockLikeRepository.On("FindLikes", mock.Anything).Return(nil, errors.New("FindLikes return error"))

	postUsecase := usecase.NewPostUseCase(fps.mockPostRepository, fps.mockLikeRepository, fps.mockHeaderHelper, fps.mockFileOsHelper)
	_, err := postUsecase.FindPosts()

	expectedError := domain.ErrInternalServerError.Error()
	assert.EqualErrorf(fps.T(), err, expectedError, "should have return error %s but got %s", expectedError, err)
}

func (fps *FindPostsSuite) TestFindPostSuccessful() {
	fps.mockPostRepository.On("FindPosts", mock.Anything).Return(&[]bson.M{
		{"_id": "postid1",
			"user_id":           "userid1",
			"visual_media_urls": []primitive.A{{"jpg.jpg"}, {"png.png"}},
			"caption":           "caption1",
			"created_date":      primitive.NewDateTimeFromTime(time.Now()),
			"updated_date":      primitive.NewDateTimeFromTime(time.Now())},
		{"_id": "postid2",
			"user_id":           "userid2",
			"visual_media_urls": []primitive.A{{"jpg.jpg"}, {"png.png"}},
			"caption":           "caption2",
			"created_date":      primitive.NewDateTimeFromTime(time.Now()),
			"updated_date":      primitive.NewDateTimeFromTime(time.Now())},
	}, nil)
	fps.mockLikeRepository.On("FindLikes", mock.Anything).Return(&[]bson.M{
		{"_id": "likeid1", "user_id": "userid1", "resource_id": "postid1", "resource_type": "post"},
	}, nil)

	postUsecase := usecase.NewPostUseCase(fps.mockPostRepository, fps.mockLikeRepository, fps.mockHeaderHelper, fps.mockFileOsHelper)
	result, err := postUsecase.FindPosts()

	assert.NoErrorf(fps.T(), err, "Should have not return error but got %s", err)
	assert.Equal(fps.T(), len(*result), 2, "length of result should be 2")
}

type UpdatePostSuite struct {
	suite.Suite
	mockPostRepository *mocks.PostRepository
	mockLikeRepository *mocks.LikeRepository
	mockFileOsHelper   *mocks.IFileOsHelper
	mockHeaderHelper   *mocks.IHeaderHelper
}

func (ups *UpdatePostSuite) SetupTest() {
	ups.mockPostRepository = new(mocks.PostRepository)
	ups.mockLikeRepository = new(mocks.LikeRepository)
	ups.mockFileOsHelper = new(mocks.IFileOsHelper)
	ups.mockHeaderHelper = new(mocks.IHeaderHelper)
}

func (ups *UpdatePostSuite) TestGetUserIdFromTokenError() {
	ups.mockHeaderHelper.On("GetUserIdFromToken", mock.Anything).Return("", errors.New("GetUserIdFromTokenError return error"))

	postUsecase := usecase.NewPostUseCase(ups.mockPostRepository, ups.mockLikeRepository, ups.mockHeaderHelper, ups.mockFileOsHelper)
	err := postUsecase.UpdatePost("postid1", "updated caption 1", "accessToken")

	expectedError := domain.ErrInternalServerError.Error()
	assert.EqualErrorf(ups.T(), err, expectedError, "Should have returned %s but got %s", expectedError, err.Error())
}

func (ups *UpdatePostSuite) TestFindPostsError() {
	ups.mockHeaderHelper.On("GetUserIdFromToken", mock.Anything).Return("userid1", nil)
	ups.mockPostRepository.On("FindPosts", mock.Anything).Return(nil, errors.New("FindPosts return error"))

	postUsecase := usecase.NewPostUseCase(ups.mockPostRepository, ups.mockLikeRepository, ups.mockHeaderHelper, ups.mockFileOsHelper)
	err := postUsecase.UpdatePost("postid1", "updated caption 1", "accessToken")

	expectedError := domain.ErrInternalServerError.Error()
	assert.EqualErrorf(ups.T(), err, expectedError, "Should have returned %s but got %s", expectedError, err.Error())
}

func (ups *UpdatePostSuite) TestPostNotFound() {
	ups.mockHeaderHelper.On("GetUserIdFromToken", mock.Anything).Return("userid1", nil)
	ups.mockPostRepository.On("FindPosts", mock.Anything).Return(&[]bson.M{}, nil)

	postUsecase := usecase.NewPostUseCase(ups.mockPostRepository, ups.mockLikeRepository, ups.mockHeaderHelper, ups.mockFileOsHelper)
	err := postUsecase.UpdatePost("postid1", "updated caption 1", "accessToken")

	expectedError := domain.ErrPostNotFound.Error()
	assert.EqualErrorf(ups.T(), err, expectedError, "Should have returned %s but got %s", expectedError, err.Error())
}

func (ups *UpdatePostSuite) TestFindOnePostError() {
	foundPost := bson.M{
		"_id":               "postid1",
		"user_id":           "userid1",
		"visual_media_urls": []primitive.A{{"jpg.jpg"}, {"png.png"}},
		"caption":           "a new caption1",
		"created_date":      primitive.NewDateTimeFromTime(time.Now()),
		"updated_date":      primitive.NewDateTimeFromTime(time.Now()),
	}
	ups.mockHeaderHelper.On("GetUserIdFromToken", mock.Anything).Return("userid1", nil)
	ups.mockPostRepository.On("FindPosts", mock.Anything).Return(&[]bson.M{foundPost}, nil)
	ups.mockPostRepository.On("FindOnePost", mock.Anything).Return(nil, errors.New("FindOnePost return error"))

	postUsecase := usecase.NewPostUseCase(ups.mockPostRepository, ups.mockLikeRepository, ups.mockHeaderHelper, ups.mockFileOsHelper)
	err := postUsecase.UpdatePost("postid1", "updated caption 1", "accessToken")

	expectedError := domain.ErrInternalServerError.Error()
	assert.EqualErrorf(ups.T(), err, expectedError, "Should have returned %s but got %s", expectedError, err.Error())
}

func (ups *UpdatePostSuite) TestUnauthorizedPostUpdate() {
	foundPosts := []bson.M{
		{
			"_id":               "postid1",
			"user_id":           "userid1",
			"visual_media_urls": []primitive.A{{"jpg.jpg"}, {"png.png"}},
			"caption":           "a new caption1",
			"created_date":      primitive.NewDateTimeFromTime(time.Now()),
			"updated_date":      primitive.NewDateTimeFromTime(time.Now()),
		},
	}
	foundPost := domain.NewPost("postid1", "userid1", []string{"jpg.jpg", "png.png"}, "caption1", 0, time.Now(), time.Now())
	ups.mockHeaderHelper.On("GetUserIdFromToken", mock.Anything).Return("userid2", nil)
	ups.mockPostRepository.On("FindPosts", mock.Anything).Return(&foundPosts, nil)
	ups.mockPostRepository.On("FindOnePost", mock.Anything).Return(foundPost, nil)

	postUsecase := usecase.NewPostUseCase(ups.mockPostRepository, ups.mockLikeRepository, ups.mockHeaderHelper, ups.mockFileOsHelper)
	err := postUsecase.UpdatePost("postid1", "updated caption 1", "accessToken")

	expectedError := domain.ErrUnauthorizedPostUpdate.Error()
	assert.EqualErrorf(ups.T(), err, expectedError, "Should have returned %s but got %s", expectedError, err.Error())
}

func (ups *UpdatePostSuite) TestUpdatePostError() {
	foundPosts := []bson.M{
		{
			"_id":               "postid1",
			"user_id":           "userid1",
			"visual_media_urls": []primitive.A{{"jpg.jpg"}, {"png.png"}},
			"caption":           "a new caption1",
			"created_date":      primitive.NewDateTimeFromTime(time.Now()),
			"updated_date":      primitive.NewDateTimeFromTime(time.Now()),
		},
	}
	foundPost := domain.NewPost("postid1", "userid1", []string{"jpg.jpg", "png.png"}, "caption1", 0, time.Now(), time.Now())
	ups.mockHeaderHelper.On("GetUserIdFromToken", mock.Anything).Return("userid1", nil)
	ups.mockPostRepository.On("FindPosts", mock.Anything).Return(&foundPosts, nil)
	ups.mockPostRepository.On("FindOnePost", mock.Anything).Return(foundPost, nil)
	ups.mockPostRepository.On("UpdatePost", mock.Anything, mock.Anything).Return(errors.New("UpdatePost return error"))

	postUsecase := usecase.NewPostUseCase(ups.mockPostRepository, ups.mockLikeRepository, ups.mockHeaderHelper, ups.mockFileOsHelper)
	err := postUsecase.UpdatePost("postid1", "updated caption 1", "accessToken")

	expectedError := domain.ErrInternalServerError.Error()
	assert.EqualErrorf(ups.T(), err, expectedError, "Should have returned %s but got %s", expectedError, err.Error())
}

func (ups *UpdatePostSuite) TestUpdatePostSuccessful() {
	foundPosts := []bson.M{
		{
			"_id":               "postid1",
			"user_id":           "userid1",
			"visual_media_urls": []primitive.A{{"jpg.jpg"}, {"png.png"}},
			"caption":           "a new caption1",
			"created_date":      primitive.NewDateTimeFromTime(time.Now()),
			"updated_date":      primitive.NewDateTimeFromTime(time.Now()),
		},
	}
	foundPost := domain.NewPost("postid1", "userid1", []string{"jpg.jpg", "png.png"}, "caption1", 0, time.Now(), time.Now())
	ups.mockHeaderHelper.On("GetUserIdFromToken", mock.Anything).Return("userid1", nil)
	ups.mockPostRepository.On("FindPosts", mock.Anything).Return(&foundPosts, nil)
	ups.mockPostRepository.On("FindOnePost", mock.Anything).Return(foundPost, nil)
	ups.mockPostRepository.On("UpdatePost", mock.Anything, mock.Anything).Return(nil)

}

type DeletePostSuite struct {
	suite.Suite
	mockPostRepository *mocks.PostRepository
	mockLikeRepository *mocks.LikeRepository
	mockFileOsHelper   *mocks.IFileOsHelper
	mockHeaderHelper   *mocks.IHeaderHelper
}

func (dps *DeletePostSuite) SetupTest() {
	dps.mockPostRepository = new(mocks.PostRepository)
	dps.mockLikeRepository = new(mocks.LikeRepository)
	dps.mockFileOsHelper = new(mocks.IFileOsHelper)
	dps.mockHeaderHelper = new(mocks.IHeaderHelper)
}

func (dps *DeletePostSuite) TestGetUserIdFromTokenError() {
	dps.mockHeaderHelper.On("GetUserIdFromToken", mock.Anything).Return("", errors.New("GetUserIdFromToken return error"))

	postUsecase := usecase.NewPostUseCase(dps.mockPostRepository, dps.mockLikeRepository, dps.mockHeaderHelper, dps.mockFileOsHelper)
	err := postUsecase.DeletePost("postid1", "token1")

	expectedError := domain.ErrInternalServerError.Error()
	assert.EqualErrorf(dps.T(), err, expectedError, "Should have return %s but got %s", expectedError, err.Error())
}

func (dps *DeletePostSuite) TestFindPostsError() {
	dps.mockHeaderHelper.On("GetUserIdFromToken", mock.Anything).Return("userid1", nil)
	dps.mockPostRepository.On("FindPosts", mock.Anything).Return(nil, errors.New("FindPosts return error"))

	postUsecase := usecase.NewPostUseCase(dps.mockPostRepository, dps.mockLikeRepository, dps.mockHeaderHelper, dps.mockFileOsHelper)
	err := postUsecase.DeletePost("postid1", "token1")

	expectedError := domain.ErrInternalServerError.Error()
	assert.EqualErrorf(dps.T(), err, expectedError, "Should have return %s but got %s", expectedError, err.Error())
}

func (dps *DeletePostSuite) TestPostNotFound() {
	dps.mockHeaderHelper.On("GetUserIdFromToken", mock.Anything).Return("userid1", nil)
	dps.mockPostRepository.On("FindPosts", mock.Anything).Return(&[]bson.M{}, nil)

	postUsecase := usecase.NewPostUseCase(dps.mockPostRepository, dps.mockLikeRepository, dps.mockHeaderHelper, dps.mockFileOsHelper)
	err := postUsecase.DeletePost("postid1", "token1")

	expectedError := domain.ErrPostNotFound.Error()
	assert.EqualErrorf(dps.T(), err, expectedError, "Should have return %s but got %s", expectedError, err.Error())
}

func (dps *DeletePostSuite) TestFindOnePostError() {
	dps.mockHeaderHelper.On("GetUserIdFromToken", mock.Anything).Return("userid1", nil)
	dps.mockPostRepository.On("FindPosts", mock.Anything).Return(&[]bson.M{
		{
			"_id":               "postid1",
			"user_id":           "userid1",
			"visual_media_urls": []primitive.A{{"jpg.jpg"}, {"png.png"}},
			"caption":           "a new caption1",
			"created_date":      primitive.NewDateTimeFromTime(time.Now()),
			"updated_date":      primitive.NewDateTimeFromTime(time.Now()),
		},
	}, nil)
	dps.mockPostRepository.On("FindOnePost", mock.Anything).Return(nil, errors.New("FindOnePost return error"))

	postUsecase := usecase.NewPostUseCase(dps.mockPostRepository, dps.mockLikeRepository, dps.mockHeaderHelper, dps.mockFileOsHelper)
	err := postUsecase.DeletePost("postid1", "token1")

	expectedError := domain.ErrInternalServerError.Error()
	assert.EqualErrorf(dps.T(), err, expectedError, "Should have return %s but got %s", expectedError, err.Error())
}

func (dps *DeletePostSuite) TestUnauthorizedPostDelete() {
	dps.mockHeaderHelper.On("GetUserIdFromToken", mock.Anything).Return("userid1", nil)
	dps.mockPostRepository.On("FindPosts", mock.Anything).Return(&[]bson.M{
		{
			"_id":               "postid1",
			"user_id":           "userid2",
			"visual_media_urls": []primitive.A{{"jpg.jpg"}, {"png.png"}},
			"caption":           "a new caption1",
			"created_date":      primitive.NewDateTimeFromTime(time.Now()),
			"updated_date":      primitive.NewDateTimeFromTime(time.Now()),
		},
	}, nil)
	dps.mockPostRepository.On("FindOnePost", mock.Anything).Return(domain.NewPost(
		"postid1", "userid2", []string{"jpg.jpg", "png.png"}, "a new caption1", 0, time.Now(), time.Now()), nil)

	postUsecase := usecase.NewPostUseCase(dps.mockPostRepository, dps.mockLikeRepository, dps.mockHeaderHelper, dps.mockFileOsHelper)
	err := postUsecase.DeletePost("postid1", "token1")

	expectedError := domain.ErrUnauthorizedPostDelete.Error()
	assert.EqualErrorf(dps.T(), err, expectedError, "Should have return %s but got %s", expectedError, err.Error())
}

func (dps *DeletePostSuite) TestDeletePostError() {
	dps.mockHeaderHelper.On("GetUserIdFromToken", mock.Anything).Return("userid1", nil)
	dps.mockPostRepository.On("FindPosts", mock.Anything).Return(&[]bson.M{
		{
			"_id":               "postid1",
			"user_id":           "userid1",
			"visual_media_urls": []primitive.A{{"jpg.jpg"}, {"png.png"}},
			"caption":           "a new caption1",
			"created_date":      primitive.NewDateTimeFromTime(time.Now()),
			"updated_date":      primitive.NewDateTimeFromTime(time.Now()),
		},
	}, nil)
	dps.mockPostRepository.On("FindOnePost", mock.Anything).Return(domain.NewPost(
		"postid1", "userid1", []string{"jpg.jpg", "png.png"}, "a new caption1", 0, time.Now(), time.Now()), nil)
	dps.mockPostRepository.On("DeletePost", mock.Anything).Return(errors.New("DeletePost return error"))

	postUsecase := usecase.NewPostUseCase(dps.mockPostRepository, dps.mockLikeRepository, dps.mockHeaderHelper, dps.mockFileOsHelper)
	err := postUsecase.DeletePost("postid1", "token1")

	expectedError := domain.ErrInternalServerError.Error()
	assert.EqualErrorf(dps.T(), err, expectedError, "Should have return %s but got %s", expectedError, err.Error())
}

func (dps *DeletePostSuite) TestDeletePostSuccessful() {
	dps.mockHeaderHelper.On("GetUserIdFromToken", mock.Anything).Return("userid1", nil)
	dps.mockPostRepository.On("FindPosts", mock.Anything).Return(&[]bson.M{
		{
			"_id":               "postid1",
			"user_id":           "userid1",
			"visual_media_urls": []primitive.A{{"jpg.jpg"}, {"png.png"}},
			"caption":           "a new caption1",
			"created_date":      primitive.NewDateTimeFromTime(time.Now()),
			"updated_date":      primitive.NewDateTimeFromTime(time.Now()),
		},
	}, nil)
	dps.mockPostRepository.On("FindOnePost", mock.Anything).Return(domain.NewPost(
		"postid1", "userid1", []string{"jpg.jpg", "png.png"}, "a new caption1", 0, time.Now(), time.Now()), nil)
	dps.mockPostRepository.On("DeletePost", mock.Anything).Return(nil)

	postUsecase := usecase.NewPostUseCase(dps.mockPostRepository, dps.mockLikeRepository, dps.mockHeaderHelper, dps.mockFileOsHelper)
	err := postUsecase.DeletePost("postid1", "token1")

	assert.NoErrorf(dps.T(), err, "should have not return error but got %s", err)
}
