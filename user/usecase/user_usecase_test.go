package usecase_test

import (
	"errors"
	"instagram-go/domain"
	"instagram-go/domain/mocks"
	"instagram-go/user/usecase"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
)

func TestInsertUser(t *testing.T) {
	suite.Run(t, new(InsertUserSuite))
}
func TestUpdateUser(t *testing.T) {
	suite.Run(t, new(UpdateUserSuite))
}

func TestVerifyCredential(t *testing.T) {
	suite.Run(t, new(VerifyCredentialSuite))
}

type InsertUserSuite struct {
	suite.Suite
	mockUserRepo             *mocks.UserRepository
	mockFileOsHelper         *mocks.IFileOsHelper
	mockHeaderHelper         *mocks.IHeaderHelper
	mockAuthenticationHelper *mocks.IAuthenticationHelper
}

func (ius *InsertUserSuite) SetupTest() {
	ius.mockUserRepo = new(mocks.UserRepository)
	ius.mockFileOsHelper = new(mocks.IFileOsHelper)
	ius.mockHeaderHelper = new(mocks.IHeaderHelper)
	ius.mockAuthenticationHelper = new(mocks.IAuthenticationHelper)
}

func (ius *InsertUserSuite) TestFindUserError() {
	mockUser := domain.NewUser("userid1", "username1", "fullname1", "password1", "email1@gmail.com", nil)
	ius.mockUserRepo.On("FindUser", mock.Anything).Return(nil, errors.New("FindUser return error"))

	userUsecase := usecase.NewUserUsecase(ius.mockUserRepo, ius.mockAuthenticationHelper, ius.mockHeaderHelper, ius.mockFileOsHelper)
	result := userUsecase.InsertUser(mockUser)

	expectedError := domain.ErrInternalServerError.Error()
	assert.EqualErrorf(ius.T(), result, expectedError, "Should have return %s but got %s", expectedError, result.Error())
}

func (ius *InsertUserSuite) TestUserNameConflict() {
	mockUser := domain.NewUser("userid1", "username1", "fullname1", "password1", "email1@gmail.com", nil)
	ius.mockUserRepo.On("FindUser", mock.Anything).Return(&[]bson.M{{
		"_id":                  "user1",
		"username":             "username1",
		"fullname":             "fullname1",
		"password":             "password1",
		"email":                "email1@gmail.com",
		"profile_pictures_url": []string{},
	}}, nil)

	userUsecase := usecase.NewUserUsecase(ius.mockUserRepo, ius.mockAuthenticationHelper, ius.mockHeaderHelper, ius.mockFileOsHelper)
	result := userUsecase.InsertUser(mockUser)

	expectedError := domain.ErrUsernameConflict.Error()
	assert.EqualErrorf(ius.T(), result, expectedError, "should have return %s but got %s", expectedError, result.Error())
}

func (ius *InsertUserSuite) TestInsertUserError() {
	mockUser := domain.NewUser("userid1", "username1", "fullname1", "password1", "email1@gmail.com", nil)
	ius.mockUserRepo.On("FindUser", mock.Anything).Return(&[]bson.M{}, nil)
	ius.mockUserRepo.On("InsertUser", mock.Anything).Return(errors.New("InsertUser return error"))

	userUsecase := usecase.NewUserUsecase(ius.mockUserRepo, ius.mockAuthenticationHelper, ius.mockHeaderHelper, ius.mockFileOsHelper)
	result := userUsecase.InsertUser(mockUser)

	expectedError := domain.ErrInternalServerError.Error()
	assert.EqualErrorf(ius.T(), result, expectedError, "should have return %s but got %s", expectedError, result.Error())
}
func (ius *InsertUserSuite) TestSuccessfulInsert() {
	mockUser := domain.NewUser("userid1", "username1", "fullname1", "password1", "email1@gmail.com", nil)
	ius.mockUserRepo.On("FindUser", mock.Anything).Return(&[]bson.M{}, nil)
	ius.mockUserRepo.On("InsertUser", mock.Anything).Return(nil)

	userUsecase := usecase.NewUserUsecase(ius.mockUserRepo, ius.mockAuthenticationHelper, ius.mockHeaderHelper, ius.mockFileOsHelper)
	result := userUsecase.InsertUser(mockUser)

	assert.NoErrorf(ius.T(), result, "should have not returned error but got %s", result)
}

type UpdateUserSuite struct {
	suite.Suite
	mockUserRepo             *mocks.UserRepository
	mockFileOsHelper         *mocks.IFileOsHelper
	mockHeaderHelper         *mocks.IHeaderHelper
	mockAuthenticationHelper *mocks.IAuthenticationHelper
}

func (uus *UpdateUserSuite) SetupTest() {
	uus.mockUserRepo = new(mocks.UserRepository)
	uus.mockFileOsHelper = new(mocks.IFileOsHelper)
	uus.mockHeaderHelper = new(mocks.IHeaderHelper)
	uus.mockAuthenticationHelper = new(mocks.IAuthenticationHelper)
}

func (uus *UpdateUserSuite) TestMkDirAllError() {
	mockUser := domain.NewUser("userid1", "username1", "fullname1", "password1", "email1@gmail.com", nil)
	uus.mockFileOsHelper.On("MkDirAll", mock.Anything, mock.Anything).Return(errors.New("MkDirAll return error"))

	userUsecase := usecase.NewUserUsecase(uus.mockUserRepo, uus.mockAuthenticationHelper, uus.mockHeaderHelper, uus.mockFileOsHelper)
	result := userUsecase.UpdateUser(mockUser, mock.Anything, nil)

	expectedError := domain.ErrInternalServerError.Error()
	assert.EqualErrorf(uus.T(), result, expectedError, "Should have return %s but got %s", expectedError, result.Error())
}

func (uus *UpdateUserSuite) TestGetUserTokenError() {
	mockUser := domain.NewUser("userid1", "username1", "fullname1", "password1", "email1@gmail.com", nil)
	uus.mockFileOsHelper.On("MkDirAll", mock.Anything, mock.Anything).Return(nil)
	uus.mockHeaderHelper.On("GetUserIdFromToken", mock.Anything).Return("user1id", errors.New("GetUserIdFromToken return error"))
	userUsecase := usecase.NewUserUsecase(uus.mockUserRepo, uus.mockAuthenticationHelper, uus.mockHeaderHelper, uus.mockFileOsHelper)
	result := userUsecase.UpdateUser(mockUser, mock.Anything, nil)

	expectedError := domain.ErrInternalServerError.Error()
	assert.EqualErrorf(uus.T(), result, expectedError, "Should have return %s but got %s", expectedError, result.Error())
}

func (uus *UpdateUserSuite) TestFindUserError() {
	mockUser := domain.NewUser("userid1", "username1", "fullname1", "password1", "email1@gmail.com", nil)
	uus.mockFileOsHelper.On("MkDirAll", mock.Anything, mock.Anything).Return(nil)
	uus.mockHeaderHelper.On("GetUserIdFromToken", mock.Anything).Return("user1id", nil)
	uus.mockUserRepo.On("FindUser", mock.Anything).Return(nil, errors.New("FindUser return error"))

	userUsecase := usecase.NewUserUsecase(uus.mockUserRepo, uus.mockAuthenticationHelper, uus.mockHeaderHelper, uus.mockFileOsHelper)
	result := userUsecase.UpdateUser(mockUser, mock.Anything, nil)

	expectedError := domain.ErrInternalServerError.Error()
	assert.EqualErrorf(uus.T(), result, expectedError, "Should have return %s but got %s", expectedError, result.Error())
}

func (uus *UpdateUserSuite) TestUserNotFoundError() {
	mockUser := domain.NewUser("userid1", "username1", "fullname1", "password1", "email1@gmail.com", nil)
	uus.mockFileOsHelper.On("MkDirAll", mock.Anything, mock.Anything).Return(nil)
	uus.mockHeaderHelper.On("GetUserIdFromToken", mock.Anything).Return("user1id", nil)
	uus.mockUserRepo.On("FindUser", mock.Anything).Return(&[]bson.M{}, nil)

	userUsecase := usecase.NewUserUsecase(uus.mockUserRepo, uus.mockAuthenticationHelper, uus.mockHeaderHelper, uus.mockFileOsHelper)
	result := userUsecase.UpdateUser(mockUser, mock.Anything, nil)

	expectedError := domain.ErrUserNotFound.Error()
	assert.EqualErrorf(uus.T(), result, expectedError, "Should have return %s but got %s", expectedError, result.Error())
}

func (uus *UpdateUserSuite) TestUnauthorizedUserUpdateError() {
	mockUser := domain.NewUser("userid1", "username1", "fullname1", "password1", "email1@gmail.com", nil)
	uus.mockFileOsHelper.On("MkDirAll", mock.Anything, mock.Anything).Return(nil)
	uus.mockHeaderHelper.On("GetUserIdFromToken", mock.Anything).Return("userid3", nil)
	uus.mockUserRepo.On("FindUser", mock.Anything).Return(&[]bson.M{
		{
			"_id":              "userid1",
			"username":         "user1",
			"full_name":        "user1",
			"password":         "password1",
			"email":            "email1@gmail.com",
			"profile_pictures": nil,
		},
	}, nil)

	userUsecase := usecase.NewUserUsecase(uus.mockUserRepo, uus.mockAuthenticationHelper, uus.mockHeaderHelper, uus.mockFileOsHelper)
	result := userUsecase.UpdateUser(mockUser, mock.Anything, nil)

	expectedError := domain.ErrUnauthorizedUserUpdate.Error()
	assert.EqualErrorf(uus.T(), result, expectedError, "Should have return %s but got %s", expectedError, result.Error())
}

func (uus *UpdateUserSuite) TestResizeAndSaveFileToLocaleError() {
	mockUser := domain.NewUser("userid1", "username1", "fullname1", "password1", "email1@gmail.com", nil)
	uus.mockFileOsHelper.On("MkDirAll", mock.Anything, mock.Anything).Return(nil)
	uus.mockHeaderHelper.On("GetUserIdFromToken", mock.Anything).Return("userid1", nil)
	uus.mockUserRepo.On("FindUser", mock.Anything).Return(&[]bson.M{
		{
			"_id":              "userid1",
			"username":         "user1",
			"full_name":        "user1",
			"password":         "password1",
			"email":            "email1@gmail.com",
			"profile_pictures": nil,
		},
	}, nil)
	file, _ := os.Open("./test_profile_pictures/jpg.jpg")
	uus.mockFileOsHelper.On("DecodeImage", file).Return(nil, "jpg.jpg", nil)
	uus.mockFileOsHelper.On("ResizeAndSaveFileToLocale", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return("jpg.jpg", errors.New("ResizeAndSaveFIleToLocale return error"))

	userUsecase := usecase.NewUserUsecase(uus.mockUserRepo, uus.mockAuthenticationHelper, uus.mockHeaderHelper, uus.mockFileOsHelper)
	result := userUsecase.UpdateUser(mockUser, mock.Anything, file)

	expectedError := domain.ErrInternalServerError.Error()
	assert.EqualErrorf(uus.T(), result, expectedError, "Should have return %s but got %s", expectedError, result.Error())
}

func (uus *UpdateUserSuite) TestDecodeImageError() {
	mockUser := domain.NewUser("userid1", "username1", "fullname1", "password1", "email1@gmail.com", nil)
	uus.mockFileOsHelper.On("MkDirAll", mock.Anything, mock.Anything).Return(nil)
	uus.mockHeaderHelper.On("GetUserIdFromToken", mock.Anything).Return("userid1", nil)
	uus.mockUserRepo.On("FindUser", mock.Anything).Return(&[]bson.M{
		{
			"_id":              "userid1",
			"username":         "user1",
			"full_name":        "user1",
			"password":         "password1",
			"email":            "email1@gmail.com",
			"profile_pictures": nil,
		},
	}, nil)
	file, _ := os.Open("./test_profile_pictures/jpg.jpg")
	uus.mockFileOsHelper.On("DecodeImage", file).Return(nil, "jpg.jpg", errors.New("DecodeImage return error"))

	userUsecase := usecase.NewUserUsecase(uus.mockUserRepo, uus.mockAuthenticationHelper, uus.mockHeaderHelper, uus.mockFileOsHelper)
	result := userUsecase.UpdateUser(mockUser, mock.Anything, file)

	expectedError := domain.ErrInternalServerError.Error()
	assert.EqualErrorf(uus.T(), result, expectedError, "Should have return %s but got %s", expectedError, result.Error())
}

func (uus *UpdateUserSuite) TestFindOneUserError() {
	mockUser := domain.NewUser("userid1", "username1", "fullname1", "password1", "email1@gmail.com", nil)
	uus.mockFileOsHelper.On("MkDirAll", mock.Anything, mock.Anything).Return(nil)
	uus.mockHeaderHelper.On("GetUserIdFromToken", mock.Anything).Return("userid1", nil)
	uus.mockUserRepo.On("FindUser", mock.Anything).Return(&[]bson.M{
		{
			"_id":              "userid1",
			"username":         "username1",
			"full_name":        "fullname1",
			"password":         "password1",
			"email":            "email1@gmail.com",
			"profile_pictures": nil,
		},
	}, nil)
	uus.mockUserRepo.On("FindOneUser", mock.Anything).Return(nil, errors.New("FindOneUser return error"))

	userUsecase := usecase.NewUserUsecase(uus.mockUserRepo, uus.mockAuthenticationHelper, uus.mockHeaderHelper, uus.mockFileOsHelper)
	result := userUsecase.UpdateUser(mockUser, mock.Anything, nil)

	expectedError := domain.ErrInternalServerError.Error()
	assert.EqualErrorf(uus.T(), result, expectedError, "Should have return %s but got %s", expectedError, result.Error())
}
func (uus *UpdateUserSuite) TestUpdateUserError() {
	mockUser := domain.NewUser("userid1", "username1", "fullname1", "password1", "email1@gmail.com", nil)
	uus.mockFileOsHelper.On("MkDirAll", mock.Anything, mock.Anything).Return(nil)
	uus.mockHeaderHelper.On("GetUserIdFromToken", mock.Anything).Return("userid1", nil)
	uus.mockUserRepo.On("FindUser", mock.Anything).Return(&[]bson.M{
		{
			"_id":              "userid1",
			"username":         "username1",
			"full_name":        "fullname1",
			"password":         "password1",
			"email":            "email1@gmail.com",
			"profile_pictures": nil,
		},
	}, nil)
	foundUser := domain.NewUser("userid1", "username1", "fullname1", "password1", "email1@gmail.com", nil)
	uus.mockUserRepo.On("FindOneUser", mock.Anything).Return(foundUser, nil)
	uus.mockUserRepo.On("UpdateUser", mock.Anything).Return(errors.New("UpdateUser return error"))

	userUsecase := usecase.NewUserUsecase(uus.mockUserRepo, uus.mockAuthenticationHelper, uus.mockHeaderHelper, uus.mockFileOsHelper)
	result := userUsecase.UpdateUser(mockUser, mock.Anything, nil)

	expectedError := domain.ErrInternalServerError.Error()
	assert.EqualErrorf(uus.T(), result, expectedError, "Should have return %s but got %s", expectedError, result.Error())
}
func (uus *UpdateUserSuite) TestUpdateSuccessful() {
	mockUser := domain.NewUser("userid1", "username1", "fullname1", "password1", "email1@gmail.com", nil)
	uus.mockFileOsHelper.On("MkDirAll", mock.Anything, mock.Anything).Return(nil)
	uus.mockHeaderHelper.On("GetUserIdFromToken", mock.Anything).Return("userid1", nil)
	uus.mockUserRepo.On("FindUser", mock.Anything).Return(&[]bson.M{
		{
			"_id":              "userid1",
			"username":         "username1",
			"full_name":        "fullname1",
			"password":         "password1",
			"email":            "email1@gmail.com",
			"profile_pictures": nil,
		},
	}, nil)
	foundUser := domain.NewUser("userid1", "username1", "fullname1", "password1", "email1@gmail.com", nil)
	uus.mockUserRepo.On("FindOneUser", mock.Anything).Return(foundUser, nil)
	uus.mockUserRepo.On("UpdateUser", mock.Anything).Return(nil)

	userUsecase := usecase.NewUserUsecase(uus.mockUserRepo, uus.mockAuthenticationHelper, uus.mockHeaderHelper, uus.mockFileOsHelper)
	result := userUsecase.UpdateUser(mockUser, mock.Anything, nil)

	assert.NoErrorf(uus.T(), result, "should have not returned error but got %s", result)
}

type VerifyCredentialSuite struct {
	suite.Suite
	mockUserRepo             *mocks.UserRepository
	mockFileOsHelper         *mocks.IFileOsHelper
	mockHeaderHelper         *mocks.IHeaderHelper
	mockAuthenticationHelper *mocks.IAuthenticationHelper
}

func (vcs *VerifyCredentialSuite) SetupTest() {
	vcs.mockUserRepo = new(mocks.UserRepository)
	vcs.mockFileOsHelper = new(mocks.IFileOsHelper)
	vcs.mockHeaderHelper = new(mocks.IHeaderHelper)
	vcs.mockAuthenticationHelper = new(mocks.IAuthenticationHelper)
}
func (vcs *VerifyCredentialSuite) TestFindUserError() {
	vcs.mockUserRepo.On("FindUser", mock.Anything).Return(nil, errors.New("FindUser return error"))

	userUsecase := usecase.NewUserUsecase(vcs.mockUserRepo, vcs.mockAuthenticationHelper, vcs.mockHeaderHelper, vcs.mockFileOsHelper)
	_, err := userUsecase.VerifyCredential("username1", "password1")
	expectedError := domain.ErrInternalServerError.Error()

	assert.EqualErrorf(vcs.T(), err, expectedError, "Should have return %s but got %s", expectedError, err.Error())
}

func (vcs *VerifyCredentialSuite) TestUserNotFoundError() {
	vcs.mockUserRepo.On("FindUser", mock.Anything).Return(&[]bson.M{}, nil)

	userUsecase := usecase.NewUserUsecase(vcs.mockUserRepo, vcs.mockAuthenticationHelper, vcs.mockHeaderHelper, vcs.mockFileOsHelper)
	_, err := userUsecase.VerifyCredential("username1", "password1")
	expectedError := domain.ErrUserNotFound.Error()

	assert.EqualErrorf(vcs.T(), err, expectedError, "Should have return %s but got %s", expectedError, err.Error())
}

func (vcs *VerifyCredentialSuite) TestFindOneUserError() {
	foundUsers := []bson.M{
		{"_id": "userid1", "username": "username1", "full_name": "fullname1", "password": "password1", "email": "email1@gmail.com", "profile_pictures": nil},
	}
	vcs.mockUserRepo.On("FindUser", mock.Anything).Return(&foundUsers, nil)
	vcs.mockUserRepo.On("FindOneUser", mock.Anything).Return(nil, errors.New("FindOneUser return error"))

	userUsecase := usecase.NewUserUsecase(vcs.mockUserRepo, vcs.mockAuthenticationHelper, vcs.mockHeaderHelper, vcs.mockFileOsHelper)
	_, err := userUsecase.VerifyCredential("username1", "password1")
	expectedError := domain.ErrInternalServerError.Error()

	assert.EqualErrorf(vcs.T(), err, expectedError, "Should have return %s but got %s", expectedError, err.Error())
}

func (vcs *VerifyCredentialSuite) TestCompareHashAndPasswordError() {
	foundUsers := []bson.M{
		{"_id": "userid1", "username": "username1", "full_name": "fullname1", "password": "password1", "email": "email1@gmail.com", "profile_pictures": nil},
	}
	foundUser := domain.NewUser("userid1", "username1", "fullname1", "password1", "email1@gmail.com", nil)
	vcs.mockUserRepo.On("FindUser", mock.Anything).Return(&foundUsers, nil)
	vcs.mockUserRepo.On("FindOneUser", mock.Anything).Return(foundUser, nil)
	vcs.mockAuthenticationHelper.On("CompareHashAndPassword", mock.Anything, mock.Anything).Return(errors.New("CompareHashAndPassword return error"))

	userUsecase := usecase.NewUserUsecase(vcs.mockUserRepo, vcs.mockAuthenticationHelper, vcs.mockHeaderHelper, vcs.mockFileOsHelper)
	_, err := userUsecase.VerifyCredential("username1", "password1")
	expectedError := domain.ErrPasswordWrong.Error()

	assert.EqualErrorf(vcs.T(), err, expectedError, "Should have return %s but got %s", expectedError, err.Error())
}

func (vcs *VerifyCredentialSuite) TestVerifyCredentialSuccessful() {
	foundUsers := []bson.M{
		{"_id": "userid1", "username": "username1", "full_name": "fullname1", "password": "password1", "email": "email1@gmail.com", "profile_pictures": nil},
	}
	foundUser := domain.NewUser("userid1", "username1", "fullname1", "password1", "email1@gmail.com", nil)
	vcs.mockUserRepo.On("FindUser", mock.Anything).Return(&foundUsers, nil)
	vcs.mockUserRepo.On("FindOneUser", mock.Anything).Return(foundUser, nil)
	vcs.mockAuthenticationHelper.On("CompareHashAndPassword", mock.Anything, mock.Anything).Return(nil)

	userUsecase := usecase.NewUserUsecase(vcs.mockUserRepo, vcs.mockAuthenticationHelper, vcs.mockHeaderHelper, vcs.mockFileOsHelper)
	_, err := userUsecase.VerifyCredential("username1", "password1")
	assert.NoErrorf(vcs.T(), err, "should have not returned error but got %s", err)
}
