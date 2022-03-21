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

func TestPostUsecaseSuite(t *testing.T) {
	suite.Run(t, new(PostUsecaseSuite))
}

type PostUsecaseSuite struct {
	suite.Suite
	mockPostRepository *mocks.PostRepository
	mockLikeRepository *mocks.LikeRepository
	mockFileOsHelper   *mocks.IFileOsHelper
	mockHeaderHelper   *mocks.IHeaderHelper
}

func (pu *PostUsecaseSuite) SetupTest() {
	pu.mockPostRepository = new(mocks.PostRepository)
	pu.mockLikeRepository = new(mocks.LikeRepository)
	pu.mockFileOsHelper = new(mocks.IFileOsHelper)
	pu.mockHeaderHelper = new(mocks.IHeaderHelper)
}

func (pu *PostUsecaseSuite) TestInsertPostGetUserIdTokenError() {
	pu.mockHeaderHelper.On("GetUserIdFromToken", mock.AnythingOfType("string")).Return("", errors.New("GetUserIdFromToken return error"))

	postUsecase := usecase.NewPostUseCase(pu.mockPostRepository, pu.mockLikeRepository, pu.mockHeaderHelper, pu.mockFileOsHelper)
	newPost := domain.NewPost("postid1", "userid1", []string{}, "caption1", 0, time.Now(), time.Now())
	err := postUsecase.InsertPost(newPost, "accessToken", []*multipart.FileHeader{})

	expectedError := domain.ErrInternalServerError.Error()
	assert.EqualErrorf(pu.T(), err, expectedError, "Should have return %s but got %s", expectedError, err)
}

func (pu *PostUsecaseSuite) TestInsertPostMkDirAllError() {
	pu.mockHeaderHelper.On("GetUserIdFromToken", mock.AnythingOfType("string")).Return("userId1", nil)
	pu.mockFileOsHelper.On("MkDirAll", mock.AnythingOfType("string"), mock.AnythingOfType("fs.FileMode")).Return(errors.New("MkDirAll return error"))

	postUsecase := usecase.NewPostUseCase(pu.mockPostRepository, pu.mockLikeRepository, pu.mockHeaderHelper, pu.mockFileOsHelper)
	newPost := domain.NewPost("postid1", "userid1", []string{}, "caption1", 0, time.Now(), time.Now())
	err := postUsecase.InsertPost(newPost, "accessToken", []*multipart.FileHeader{})

	expectedError := domain.ErrInternalServerError.Error()
	assert.EqualErrorf(pu.T(), err, expectedError, "Should have return %s but got %s", expectedError, err)
}

func (pu *PostUsecaseSuite) TestInsertPostInsertPostError() {
	pu.mockHeaderHelper.On("GetUserIdFromToken", mock.AnythingOfType("string")).Return("userId1", nil)
	pu.mockFileOsHelper.On("MkDirAll", mock.AnythingOfType("string"), mock.AnythingOfType("fs.FileMode")).Return(nil)
	pu.mockPostRepository.On("InsertPost", mock.AnythingOfType("*domain.Post")).Return(errors.New("InsertPost return error"))

	postUsecase := usecase.NewPostUseCase(pu.mockPostRepository, pu.mockLikeRepository, pu.mockHeaderHelper, pu.mockFileOsHelper)
	newPost := domain.NewPost("postid1", "userid1", []string{}, "caption1", 0, time.Now(), time.Now())
	err := postUsecase.InsertPost(newPost, "accessToken", []*multipart.FileHeader{})

	expectedError := domain.ErrInternalServerError.Error()
	assert.EqualErrorf(pu.T(), err, expectedError, "Should have return %s but got %s", expectedError, err)
}

func (pu *PostUsecaseSuite) TestInsertPostSuccessful() {
	pu.mockHeaderHelper.On("GetUserIdFromToken", mock.AnythingOfType("string")).Return("userId1", nil)
	pu.mockFileOsHelper.On("MkDirAll", mock.AnythingOfType("string"), mock.AnythingOfType("fs.FileMode")).Return(nil)
	pu.mockPostRepository.On("InsertPost", mock.AnythingOfType("*domain.Post")).Return(nil)

	postUsecase := usecase.NewPostUseCase(pu.mockPostRepository, pu.mockLikeRepository, pu.mockHeaderHelper, pu.mockFileOsHelper)
	newPost := domain.NewPost("postid1", "userid1", []string{}, "caption1", 0, time.Now(), time.Now())
	err := postUsecase.InsertPost(newPost, "accessToken", []*multipart.FileHeader{})

	assert.NoErrorf(pu.T(), err, "should have not returned error but got %s", err)
}

func (pu *PostUsecaseSuite) TestFindPostFindPostsError() {
	pu.mockPostRepository.On("FindPosts", mock.AnythingOfType("M")).Return(nil, errors.New("FindPosts return error"))

	postUsecase := usecase.NewPostUseCase(pu.mockPostRepository, pu.mockLikeRepository, pu.mockHeaderHelper, pu.mockFileOsHelper)
	_, err := postUsecase.FindPosts()

	expectedError := domain.ErrInternalServerError.Error()
	assert.EqualErrorf(pu.T(), err, expectedError, "should have return error %s but got %s", expectedError, err)
}

func (pu *PostUsecaseSuite) TestFindPostFindLikesError() {
	pu.mockPostRepository.On("FindPosts", mock.AnythingOfType("M")).Return(&[]bson.M{
		{"_id": "postid1",
			"user_id":           "userid1",
			"visual_media_urls": []primitive.A{{"jpg.jpg"}, {"png.png"}},
			"caption":           "caption1",
			"created_date":      primitive.NewDateTimeFromTime(time.Now()),
			"updated_date":      primitive.NewDateTimeFromTime(time.Now())},
	}, nil)
	pu.mockLikeRepository.On("FindLikes", mock.AnythingOfType("M")).Return(nil, errors.New("FindLikes return error"))

	postUsecase := usecase.NewPostUseCase(pu.mockPostRepository, pu.mockLikeRepository, pu.mockHeaderHelper, pu.mockFileOsHelper)
	_, err := postUsecase.FindPosts()

	expectedError := domain.ErrInternalServerError.Error()
	assert.EqualErrorf(pu.T(), err, expectedError, "should have return error %s but got %s", expectedError, err)
}

func (pu *PostUsecaseSuite) TestFindPostSuccessful() {
	pu.mockPostRepository.On("FindPosts", mock.AnythingOfType("M")).Return(&[]bson.M{
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
	pu.mockLikeRepository.On("FindLikes", mock.AnythingOfType("M")).Return(&[]bson.M{
		{"_id": "likeid1", "user_id": "userid1", "resource_id": "postid1", "resource_type": "post"},
	}, nil)

	postUsecase := usecase.NewPostUseCase(pu.mockPostRepository, pu.mockLikeRepository, pu.mockHeaderHelper, pu.mockFileOsHelper)
	result, err := postUsecase.FindPosts()

	assert.NoErrorf(pu.T(), err, "Should have not return error but got %s", err)
	assert.Equal(pu.T(), len(*result), 2, "length of result should be 2")
}

func (pu *PostUsecaseSuite) TestUpdatePostGetUserIdFromTokenError() {
	pu.mockHeaderHelper.On("GetUserIdFromToken", mock.AnythingOfType("string")).Return("", errors.New("GetUserIdFromTokenError return error"))

	postUsecase := usecase.NewPostUseCase(pu.mockPostRepository, pu.mockLikeRepository, pu.mockHeaderHelper, pu.mockFileOsHelper)
	err := postUsecase.UpdatePost("postid1", "updated caption 1", "accessToken")

	expectedError := domain.ErrInternalServerError.Error()
	assert.EqualErrorf(pu.T(), err, expectedError, "Should have returned %s but got %s", expectedError, err.Error())
}

func (pu *PostUsecaseSuite) TestUpdatePostFindPostsError() {
	pu.mockHeaderHelper.On("GetUserIdFromToken", mock.AnythingOfType("string")).Return("userid1", nil)
	pu.mockPostRepository.On("FindPosts", mock.AnythingOfType("M")).Return(nil, errors.New("FindPosts return error"))

	postUsecase := usecase.NewPostUseCase(pu.mockPostRepository, pu.mockLikeRepository, pu.mockHeaderHelper, pu.mockFileOsHelper)
	err := postUsecase.UpdatePost("postid1", "updated caption 1", "accessToken")

	expectedError := domain.ErrInternalServerError.Error()
	assert.EqualErrorf(pu.T(), err, expectedError, "Should have returned %s but got %s", expectedError, err.Error())
}

func (pu *PostUsecaseSuite) TestUpdatePostPostNotFound() {
	pu.mockHeaderHelper.On("GetUserIdFromToken", mock.AnythingOfType("string")).Return("userid1", nil)
	pu.mockPostRepository.On("FindPosts", mock.AnythingOfType("M")).Return(&[]bson.M{}, nil)

	postUsecase := usecase.NewPostUseCase(pu.mockPostRepository, pu.mockLikeRepository, pu.mockHeaderHelper, pu.mockFileOsHelper)
	err := postUsecase.UpdatePost("postid1", "updated caption 1", "accessToken")

	expectedError := domain.ErrPostNotFound.Error()
	assert.EqualErrorf(pu.T(), err, expectedError, "Should have returned %s but got %s", expectedError, err.Error())
}

func (pu *PostUsecaseSuite) TestUpdatePostFindOnePostError() {
	foundPost := bson.M{
		"_id":               "postid1",
		"user_id":           "userid1",
		"visual_media_urls": []primitive.A{{"jpg.jpg"}, {"png.png"}},
		"caption":           "a new caption1",
		"created_date":      primitive.NewDateTimeFromTime(time.Now()),
		"updated_date":      primitive.NewDateTimeFromTime(time.Now()),
	}
	pu.mockHeaderHelper.On("GetUserIdFromToken", mock.AnythingOfType("string")).Return("userid1", nil)
	pu.mockPostRepository.On("FindPosts", mock.AnythingOfType("M")).Return(&[]bson.M{foundPost}, nil)
	pu.mockPostRepository.On("FindOnePost", mock.AnythingOfType("string")).Return(nil, errors.New("FindOnePost return error"))

	postUsecase := usecase.NewPostUseCase(pu.mockPostRepository, pu.mockLikeRepository, pu.mockHeaderHelper, pu.mockFileOsHelper)
	err := postUsecase.UpdatePost("postid1", "updated caption 1", "accessToken")

	expectedError := domain.ErrInternalServerError.Error()
	assert.EqualErrorf(pu.T(), err, expectedError, "Should have returned %s but got %s", expectedError, err.Error())
}

func (pu *PostUsecaseSuite) TestUpdatePostUnauthorizedPostUpdate() {
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
	pu.mockHeaderHelper.On("GetUserIdFromToken", mock.AnythingOfType("string")).Return("userid2", nil)
	pu.mockPostRepository.On("FindPosts", mock.AnythingOfType("M")).Return(&foundPosts, nil)
	pu.mockPostRepository.On("FindOnePost", mock.AnythingOfType("string")).Return(foundPost, nil)

	postUsecase := usecase.NewPostUseCase(pu.mockPostRepository, pu.mockLikeRepository, pu.mockHeaderHelper, pu.mockFileOsHelper)
	err := postUsecase.UpdatePost("postid1", "updated caption 1", "accessToken")

	expectedError := domain.ErrUnauthorizedPostUpdate.Error()
	assert.EqualErrorf(pu.T(), err, expectedError, "Should have returned %s but got %s", expectedError, err.Error())
}

func (pu *PostUsecaseSuite) TestUpdatePostUpdatePostError() {
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
	pu.mockHeaderHelper.On("GetUserIdFromToken", mock.AnythingOfType("string")).Return("userid1", nil)
	pu.mockPostRepository.On("FindPosts", mock.AnythingOfType("M")).Return(&foundPosts, nil)
	pu.mockPostRepository.On("FindOnePost", mock.AnythingOfType("string")).Return(foundPost, nil)
	pu.mockPostRepository.On("UpdatePost", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(errors.New("UpdatePost return error"))

	postUsecase := usecase.NewPostUseCase(pu.mockPostRepository, pu.mockLikeRepository, pu.mockHeaderHelper, pu.mockFileOsHelper)
	err := postUsecase.UpdatePost("postid1", "updated caption 1", "accessToken")

	expectedError := domain.ErrInternalServerError.Error()
	assert.EqualErrorf(pu.T(), err, expectedError, "Should have returned %s but got %s", expectedError, err.Error())
}

func (pu *PostUsecaseSuite) TestUpdatePostSuccessful() {
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
	pu.mockHeaderHelper.On("GetUserIdFromToken", mock.AnythingOfType("string")).Return("userid1", nil)
	pu.mockPostRepository.On("FindPosts", mock.AnythingOfType("M")).Return(&foundPosts, nil)
	pu.mockPostRepository.On("FindOnePost", mock.AnythingOfType("string")).Return(foundPost, nil)
	pu.mockPostRepository.On("UpdatePost", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil)

}

func (pu *PostUsecaseSuite) TestDeletePostGetUserIdFromTokenError() {
	pu.mockHeaderHelper.On("GetUserIdFromToken", mock.AnythingOfType("string")).Return("", errors.New("GetUserIdFromToken return error"))

	postUsecase := usecase.NewPostUseCase(pu.mockPostRepository, pu.mockLikeRepository, pu.mockHeaderHelper, pu.mockFileOsHelper)
	err := postUsecase.DeletePost("postid1", "token1")

	expectedError := domain.ErrInternalServerError.Error()
	assert.EqualErrorf(pu.T(), err, expectedError, "Should have return %s but got %s", expectedError, err.Error())
}

func (pu *PostUsecaseSuite) TestDeletePostFindPostsError() {
	pu.mockHeaderHelper.On("GetUserIdFromToken", mock.AnythingOfType("string")).Return("userid1", nil)
	pu.mockPostRepository.On("FindPosts", mock.AnythingOfType("M")).Return(nil, errors.New("FindPosts return error"))

	postUsecase := usecase.NewPostUseCase(pu.mockPostRepository, pu.mockLikeRepository, pu.mockHeaderHelper, pu.mockFileOsHelper)
	err := postUsecase.DeletePost("postid1", "token1")

	expectedError := domain.ErrInternalServerError.Error()
	assert.EqualErrorf(pu.T(), err, expectedError, "Should have return %s but got %s", expectedError, err.Error())
}

func (pu *PostUsecaseSuite) TestDeletePostPostNotFound() {
	pu.mockHeaderHelper.On("GetUserIdFromToken", mock.AnythingOfType("string")).Return("userid1", nil)
	pu.mockPostRepository.On("FindPosts", mock.AnythingOfType("M")).Return(&[]bson.M{}, nil)

	postUsecase := usecase.NewPostUseCase(pu.mockPostRepository, pu.mockLikeRepository, pu.mockHeaderHelper, pu.mockFileOsHelper)
	err := postUsecase.DeletePost("postid1", "token1")

	expectedError := domain.ErrPostNotFound.Error()
	assert.EqualErrorf(pu.T(), err, expectedError, "Should have return %s but got %s", expectedError, err.Error())
}

func (pu *PostUsecaseSuite) TestDeletePostFindOnePostError() {
	pu.mockHeaderHelper.On("GetUserIdFromToken", mock.AnythingOfType("string")).Return("userid1", nil)
	pu.mockPostRepository.On("FindPosts", mock.AnythingOfType("M")).Return(&[]bson.M{
		{
			"_id":               "postid1",
			"user_id":           "userid1",
			"visual_media_urls": []primitive.A{{"jpg.jpg"}, {"png.png"}},
			"caption":           "a new caption1",
			"created_date":      primitive.NewDateTimeFromTime(time.Now()),
			"updated_date":      primitive.NewDateTimeFromTime(time.Now()),
		},
	}, nil)
	pu.mockPostRepository.On("FindOnePost", mock.AnythingOfType("string")).Return(nil, errors.New("FindOnePost return error"))

	postUsecase := usecase.NewPostUseCase(pu.mockPostRepository, pu.mockLikeRepository, pu.mockHeaderHelper, pu.mockFileOsHelper)
	err := postUsecase.DeletePost("postid1", "token1")

	expectedError := domain.ErrInternalServerError.Error()
	assert.EqualErrorf(pu.T(), err, expectedError, "Should have return %s but got %s", expectedError, err.Error())
}

func (pu *PostUsecaseSuite) TestDeletePostUnauthorizedPostDelete() {
	pu.mockHeaderHelper.On("GetUserIdFromToken", mock.AnythingOfType("string")).Return("userid1", nil)
	pu.mockPostRepository.On("FindPosts", mock.AnythingOfType("M")).Return(&[]bson.M{
		{
			"_id":               "postid1",
			"user_id":           "userid2",
			"visual_media_urls": []primitive.A{{"jpg.jpg"}, {"png.png"}},
			"caption":           "a new caption1",
			"created_date":      primitive.NewDateTimeFromTime(time.Now()),
			"updated_date":      primitive.NewDateTimeFromTime(time.Now()),
		},
	}, nil)
	pu.mockPostRepository.On("FindOnePost", mock.AnythingOfType("string")).Return(domain.NewPost(
		"postid1", "userid2", []string{"jpg.jpg", "png.png"}, "a new caption1", 0, time.Now(), time.Now()), nil)

	postUsecase := usecase.NewPostUseCase(pu.mockPostRepository, pu.mockLikeRepository, pu.mockHeaderHelper, pu.mockFileOsHelper)
	err := postUsecase.DeletePost("postid1", "token1")

	expectedError := domain.ErrUnauthorizedPostDelete.Error()
	assert.EqualErrorf(pu.T(), err, expectedError, "Should have return %s but got %s", expectedError, err.Error())
}

func (pu *PostUsecaseSuite) TestDeletePostDeletePostError() {
	pu.mockHeaderHelper.On("GetUserIdFromToken", mock.AnythingOfType("string")).Return("userid1", nil)
	pu.mockPostRepository.On("FindPosts", mock.AnythingOfType("M")).Return(&[]bson.M{
		{
			"_id":               "postid1",
			"user_id":           "userid1",
			"visual_media_urls": []primitive.A{{"jpg.jpg"}, {"png.png"}},
			"caption":           "a new caption1",
			"created_date":      primitive.NewDateTimeFromTime(time.Now()),
			"updated_date":      primitive.NewDateTimeFromTime(time.Now()),
		},
	}, nil)
	pu.mockPostRepository.On("FindOnePost", mock.AnythingOfType("string")).Return(domain.NewPost(
		"postid1", "userid1", []string{"jpg.jpg", "png.png"}, "a new caption1", 0, time.Now(), time.Now()), nil)
	pu.mockPostRepository.On("DeletePost", mock.AnythingOfType("string")).Return(errors.New("DeletePost return error"))

	postUsecase := usecase.NewPostUseCase(pu.mockPostRepository, pu.mockLikeRepository, pu.mockHeaderHelper, pu.mockFileOsHelper)
	err := postUsecase.DeletePost("postid1", "token1")

	expectedError := domain.ErrInternalServerError.Error()
	assert.EqualErrorf(pu.T(), err, expectedError, "Should have return %s but got %s", expectedError, err.Error())
}

func (pu *PostUsecaseSuite) TestDeletePostSuccessful() {
	pu.mockHeaderHelper.On("GetUserIdFromToken", mock.AnythingOfType("string")).Return("userid1", nil)
	pu.mockPostRepository.On("FindPosts", mock.AnythingOfType("M")).Return(&[]bson.M{
		{
			"_id":               "postid1",
			"user_id":           "userid1",
			"visual_media_urls": []primitive.A{{"jpg.jpg"}, {"png.png"}},
			"caption":           "a new caption1",
			"created_date":      primitive.NewDateTimeFromTime(time.Now()),
			"updated_date":      primitive.NewDateTimeFromTime(time.Now()),
		},
	}, nil)
	pu.mockPostRepository.On("FindOnePost", mock.AnythingOfType("string")).Return(domain.NewPost(
		"postid1", "userid1", []string{"jpg.jpg", "png.png"}, "a new caption1", 0, time.Now(), time.Now()), nil)
	pu.mockPostRepository.On("DeletePost", mock.AnythingOfType("string")).Return(nil)

	postUsecase := usecase.NewPostUseCase(pu.mockPostRepository, pu.mockLikeRepository, pu.mockHeaderHelper, pu.mockFileOsHelper)
	err := postUsecase.DeletePost("postid1", "token1")

	assert.NoErrorf(pu.T(), err, "should have not return error but got %s", err)
}
