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

func TestUserusecase(t *testing.T) {
	suite.Run(t, new(UserUsecaseSuite))
}

type UserUsecaseSuite struct {
	suite.Suite
	mockUserRepo             *mocks.UserRepository
	mockFileOsHelper         *mocks.IFileOsHelper
	mockHeaderHelper         *mocks.IHeaderHelper
	mockAuthenticationHelper *mocks.IAuthenticationHelper
}

func (us *UserUsecaseSuite) SetupTest() {
	us.mockUserRepo = new(mocks.UserRepository)
	us.mockFileOsHelper = new(mocks.IFileOsHelper)
	us.mockHeaderHelper = new(mocks.IHeaderHelper)
	us.mockAuthenticationHelper = new(mocks.IAuthenticationHelper)
}

func (us *UserUsecaseSuite) TestInsertUserFindUserError() {
	mockUser := domain.NewUser("userid1", "username1", "fullname1", "password1", "email1@gmail.com", nil)
	us.mockUserRepo.On("FindUser", mock.AnythingOfType("M")).Return(nil, errors.New("FindUser return error"))

	userUsecase := usecase.NewUserUsecase(us.mockUserRepo, us.mockAuthenticationHelper, us.mockHeaderHelper, us.mockFileOsHelper)
	result := userUsecase.InsertUser(mockUser)

	expectedError := domain.ErrInternalServerError.Error()
	assert.EqualErrorf(us.T(), result, expectedError, "Should have return %s but got %s", expectedError, result.Error())
}

func (us *UserUsecaseSuite) TestInsertUserUserNameConflict() {
	mockUser := domain.NewUser("userid1", "username1", "fullname1", "password1", "email1@gmail.com", nil)
	us.mockUserRepo.On("FindUser", mock.AnythingOfType("M")).Return(&[]bson.M{{
		"_id":                  "user1",
		"username":             "username1",
		"fullname":             "fullname1",
		"password":             "password1",
		"email":                "email1@gmail.com",
		"profile_pictures_url": []string{},
	}}, nil)

	userUsecase := usecase.NewUserUsecase(us.mockUserRepo, us.mockAuthenticationHelper, us.mockHeaderHelper, us.mockFileOsHelper)
	result := userUsecase.InsertUser(mockUser)

	expectedError := domain.ErrUsernameConflict.Error()
	assert.EqualErrorf(us.T(), result, expectedError, "should have return %s but got %s", expectedError, result.Error())
}

func (us *UserUsecaseSuite) TestInsertUserError() {
	mockUser := domain.NewUser("userid1", "username1", "fullname1", "password1", "email1@gmail.com", nil)
	us.mockUserRepo.On("FindUser", mock.AnythingOfType("M")).Return(&[]bson.M{}, nil)
	us.mockUserRepo.On("InsertUser", mock.AnythingOfType("*domain.User")).Return(errors.New("InsertUser return error"))

	userUsecase := usecase.NewUserUsecase(us.mockUserRepo, us.mockAuthenticationHelper, us.mockHeaderHelper, us.mockFileOsHelper)
	result := userUsecase.InsertUser(mockUser)

	expectedError := domain.ErrInternalServerError.Error()
	assert.EqualErrorf(us.T(), result, expectedError, "should have return %s but got %s", expectedError, result.Error())
}

func (us *UserUsecaseSuite) TestInsertUserSuccessful() {
	mockUser := domain.NewUser("userid1", "username1", "fullname1", "password1", "email1@gmail.com", nil)
	us.mockUserRepo.On("FindUser", mock.AnythingOfType("M")).Return(&[]bson.M{}, nil)
	us.mockUserRepo.On("InsertUser", mock.AnythingOfType("*domain.User")).Return(nil)

	userUsecase := usecase.NewUserUsecase(us.mockUserRepo, us.mockAuthenticationHelper, us.mockHeaderHelper, us.mockFileOsHelper)
	result := userUsecase.InsertUser(mockUser)

	assert.NoErrorf(us.T(), result, "should have not returned error but got %s", result)
}

func (us *UserUsecaseSuite) TestUpdateUserMkDirAllError() {
	mockUser := domain.NewUser("userid1", "username1", "fullname1", "password1", "email1@gmail.com", nil)
	us.mockFileOsHelper.On("MkDirAll", mock.AnythingOfType("string"), mock.AnythingOfType("fs.FileMode")).Return(errors.New("MkDirAll return error"))

	userUsecase := usecase.NewUserUsecase(us.mockUserRepo, us.mockAuthenticationHelper, us.mockHeaderHelper, us.mockFileOsHelper)
	result := userUsecase.UpdateUser(mockUser, "test", nil)

	expectedError := domain.ErrInternalServerError.Error()
	assert.EqualErrorf(us.T(), result, expectedError, "Should have return %s but got %s", expectedError, result.Error())
}

func (us *UserUsecaseSuite) TestUpdateUserGetUserTokenError() {
	mockUser := domain.NewUser("userid1", "username1", "fullname1", "password1", "email1@gmail.com", nil)
	us.mockFileOsHelper.On("MkDirAll", mock.AnythingOfType("string"), mock.AnythingOfType("fs.FileMode")).Return(nil)
	us.mockHeaderHelper.On("GetUserIdFromToken", mock.AnythingOfType("string")).Return("user1id", errors.New("GetUserIdFromToken return error"))
	userUsecase := usecase.NewUserUsecase(us.mockUserRepo, us.mockAuthenticationHelper, us.mockHeaderHelper, us.mockFileOsHelper)
	result := userUsecase.UpdateUser(mockUser, "test", nil)

	expectedError := domain.ErrInternalServerError.Error()
	assert.EqualErrorf(us.T(), result, expectedError, "Should have return %s but got %s", expectedError, result.Error())
}

func (us *UserUsecaseSuite) TestUpdateUserFindUserError() {
	mockUser := domain.NewUser("userid1", "username1", "fullname1", "password1", "email1@gmail.com", nil)
	us.mockFileOsHelper.On("MkDirAll", mock.AnythingOfType("string"), mock.AnythingOfType("fs.FileMode")).Return(nil)
	us.mockHeaderHelper.On("GetUserIdFromToken", mock.AnythingOfType("string")).Return("user1id", nil)
	us.mockUserRepo.On("FindUser", mock.AnythingOfType("M")).Return(nil, errors.New("FindUser return error"))

	userUsecase := usecase.NewUserUsecase(us.mockUserRepo, us.mockAuthenticationHelper, us.mockHeaderHelper, us.mockFileOsHelper)
	result := userUsecase.UpdateUser(mockUser, "test", nil)

	expectedError := domain.ErrInternalServerError.Error()
	assert.EqualErrorf(us.T(), result, expectedError, "Should have return %s but got %s", expectedError, result.Error())
}

func (us *UserUsecaseSuite) TestUpdateUserUserNotFoundError() {
	mockUser := domain.NewUser("userid1", "username1", "fullname1", "password1", "email1@gmail.com", nil)
	us.mockFileOsHelper.On("MkDirAll", mock.AnythingOfType("string"), mock.AnythingOfType("fs.FileMode")).Return(nil)
	us.mockHeaderHelper.On("GetUserIdFromToken", mock.AnythingOfType("string")).Return("user1id", nil)
	us.mockUserRepo.On("FindUser", mock.AnythingOfType("M")).Return(&[]bson.M{}, nil)

	userUsecase := usecase.NewUserUsecase(us.mockUserRepo, us.mockAuthenticationHelper, us.mockHeaderHelper, us.mockFileOsHelper)
	result := userUsecase.UpdateUser(mockUser, "test", nil)

	expectedError := domain.ErrUserNotFound.Error()
	assert.EqualErrorf(us.T(), result, expectedError, "Should have return %s but got %s", expectedError, result.Error())
}

func (us *UserUsecaseSuite) TestUpdateUserUnauthorizedUserUpdateError() {
	mockUser := domain.NewUser("userid1", "username1", "fullname1", "password1", "email1@gmail.com", nil)
	us.mockFileOsHelper.On("MkDirAll", mock.AnythingOfType("string"), mock.AnythingOfType("fs.FileMode")).Return(nil)
	us.mockHeaderHelper.On("GetUserIdFromToken", mock.AnythingOfType("string")).Return("userid3", nil)
	us.mockUserRepo.On("FindUser", mock.AnythingOfType("M")).Return(&[]bson.M{
		{
			"_id":              "userid1",
			"username":         "user1",
			"full_name":        "user1",
			"password":         "password1",
			"email":            "email1@gmail.com",
			"profile_pictures": nil,
		},
	}, nil)

	userUsecase := usecase.NewUserUsecase(us.mockUserRepo, us.mockAuthenticationHelper, us.mockHeaderHelper, us.mockFileOsHelper)
	result := userUsecase.UpdateUser(mockUser, "test", nil)

	expectedError := domain.ErrUnauthorizedUserUpdate.Error()
	assert.EqualErrorf(us.T(), result, expectedError, "Should have return %s but got %s", expectedError, result.Error())
}

func (us *UserUsecaseSuite) TestUpdateUserResizeAndSaveFileToLocaleError() {
	mockUser := domain.NewUser("userid1", "username1", "fullname1", "password1", "email1@gmail.com", nil)
	us.mockFileOsHelper.On("MkDirAll", mock.AnythingOfType("string"), mock.AnythingOfType("fs.FileMode")).Return(nil)
	us.mockHeaderHelper.On("GetUserIdFromToken", mock.AnythingOfType("string")).Return("userid1", nil)
	us.mockUserRepo.On("FindUser", mock.AnythingOfType("M")).Return(&[]bson.M{
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
	us.mockFileOsHelper.On("DecodeImage", file).Return(nil, "jpg.jpg", nil)
	us.mockFileOsHelper.On("ResizeAndSaveFileToLocale", mock.AnythingOfType("string"), mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return("jpg.jpg", errors.New("ResizeAndSaveFIleToLocale return error"))

	userUsecase := usecase.NewUserUsecase(us.mockUserRepo, us.mockAuthenticationHelper, us.mockHeaderHelper, us.mockFileOsHelper)
	result := userUsecase.UpdateUser(mockUser, "token", file)

	expectedError := domain.ErrInternalServerError.Error()
	assert.EqualErrorf(us.T(), result, expectedError, "Should have return %s but got %s", expectedError, result.Error())
}

func (us *UserUsecaseSuite) TestUpdateUserDecodeImageError() {
	mockUser := domain.NewUser("userid1", "username1", "fullname1", "password1", "email1@gmail.com", nil)
	us.mockFileOsHelper.On("MkDirAll", mock.AnythingOfType("string"), mock.AnythingOfType("fs.FileMode")).Return(nil)
	us.mockHeaderHelper.On("GetUserIdFromToken", mock.AnythingOfType("string")).Return("userid1", nil)
	us.mockUserRepo.On("FindUser", mock.AnythingOfType("M")).Return(&[]bson.M{
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
	us.mockFileOsHelper.On("DecodeImage", file).Return(nil, "jpg.jpg", errors.New("DecodeImage return error"))

	userUsecase := usecase.NewUserUsecase(us.mockUserRepo, us.mockAuthenticationHelper, us.mockHeaderHelper, us.mockFileOsHelper)
	result := userUsecase.UpdateUser(mockUser, "token", file)

	expectedError := domain.ErrInternalServerError.Error()
	assert.EqualErrorf(us.T(), result, expectedError, "Should have return %s but got %s", expectedError, result.Error())
}

func (us *UserUsecaseSuite) TestUpdateUserFindOneUserError() {
	mockUser := domain.NewUser("userid1", "username1", "fullname1", "password1", "email1@gmail.com", nil)
	us.mockFileOsHelper.On("MkDirAll", mock.AnythingOfType("string"), mock.AnythingOfType("fs.FileMode")).Return(nil)
	us.mockHeaderHelper.On("GetUserIdFromToken", mock.AnythingOfType("string")).Return("userid1", nil)
	us.mockUserRepo.On("FindUser", mock.AnythingOfType("M")).Return(&[]bson.M{
		{
			"_id":              "userid1",
			"username":         "username1",
			"full_name":        "fullname1",
			"password":         "password1",
			"email":            "email1@gmail.com",
			"profile_pictures": nil,
		},
	}, nil)
	us.mockUserRepo.On("FindOneUser", mock.AnythingOfType("M")).Return(nil, errors.New("FindOneUser return error"))

	userUsecase := usecase.NewUserUsecase(us.mockUserRepo, us.mockAuthenticationHelper, us.mockHeaderHelper, us.mockFileOsHelper)
	result := userUsecase.UpdateUser(mockUser, "token", nil)

	expectedError := domain.ErrInternalServerError.Error()
	assert.EqualErrorf(us.T(), result, expectedError, "Should have return %s but got %s", expectedError, result.Error())
}

func (us *UserUsecaseSuite) TestUpdateUserError() {
	mockUser := domain.NewUser("userid1", "username1", "fullname1", "password1", "email1@gmail.com", nil)
	us.mockFileOsHelper.On("MkDirAll", mock.AnythingOfType("string"), mock.AnythingOfType("fs.FileMode")).Return(nil)
	us.mockHeaderHelper.On("GetUserIdFromToken", mock.AnythingOfType("string")).Return("userid1", nil)
	us.mockUserRepo.On("FindUser", mock.AnythingOfType("M")).Return(&[]bson.M{
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
	us.mockUserRepo.On("FindOneUser", mock.AnythingOfType("M")).Return(foundUser, nil)
	us.mockUserRepo.On("UpdateUser", mock.AnythingOfType("*domain.User")).Return(errors.New("UpdateUser return error"))

	userUsecase := usecase.NewUserUsecase(us.mockUserRepo, us.mockAuthenticationHelper, us.mockHeaderHelper, us.mockFileOsHelper)
	result := userUsecase.UpdateUser(mockUser, "token", nil)

	expectedError := domain.ErrInternalServerError.Error()
	assert.EqualErrorf(us.T(), result, expectedError, "Should have return %s but got %s", expectedError, result.Error())
}

func (us *UserUsecaseSuite) TestUpdateUserSuccessful() {
	mockUser := domain.NewUser("userid1", "username1", "fullname1", "password1", "email1@gmail.com", nil)
	us.mockFileOsHelper.On("MkDirAll", mock.AnythingOfType("string"), mock.AnythingOfType("fs.FileMode")).Return(nil)
	us.mockHeaderHelper.On("GetUserIdFromToken", mock.AnythingOfType("string")).Return("userid1", nil)
	us.mockUserRepo.On("FindUser", mock.AnythingOfType("M")).Return(&[]bson.M{
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
	us.mockUserRepo.On("FindOneUser", mock.AnythingOfType("M")).Return(foundUser, nil)
	us.mockUserRepo.On("UpdateUser", mock.AnythingOfType("*domain.User")).Return(nil)

	userUsecase := usecase.NewUserUsecase(us.mockUserRepo, us.mockAuthenticationHelper, us.mockHeaderHelper, us.mockFileOsHelper)
	result := userUsecase.UpdateUser(mockUser, "token", nil)

	assert.NoErrorf(us.T(), result, "should have not returned error but got %s", result)
}

func (us *UserUsecaseSuite) TestVerfiyCredentialFindUserError() {
	us.mockUserRepo.On("FindUser", mock.AnythingOfType("M")).Return(nil, errors.New("FindUser return error"))

	userUsecase := usecase.NewUserUsecase(us.mockUserRepo, us.mockAuthenticationHelper, us.mockHeaderHelper, us.mockFileOsHelper)
	_, err := userUsecase.VerifyCredential("username1", "password1")
	expectedError := domain.ErrInternalServerError.Error()

	assert.EqualErrorf(us.T(), err, expectedError, "Should have return %s but got %s", expectedError, err.Error())
}

func (us *UserUsecaseSuite) TestVerifyCredentialUserNotFoundError() {
	us.mockUserRepo.On("FindUser", mock.AnythingOfType("M")).Return(&[]bson.M{}, nil)

	userUsecase := usecase.NewUserUsecase(us.mockUserRepo, us.mockAuthenticationHelper, us.mockHeaderHelper, us.mockFileOsHelper)
	_, err := userUsecase.VerifyCredential("username1", "password1")
	expectedError := domain.ErrUserNotFound.Error()

	assert.EqualErrorf(us.T(), err, expectedError, "Should have return %s but got %s", expectedError, err.Error())
}

func (us *UserUsecaseSuite) TestVerifyCredentialFindOneUserError() {
	foundUsers := []bson.M{
		{"_id": "userid1", "username": "username1", "full_name": "fullname1", "password": "password1", "email": "email1@gmail.com", "profile_pictures": nil},
	}
	us.mockUserRepo.On("FindUser", mock.AnythingOfType("M")).Return(&foundUsers, nil)
	us.mockUserRepo.On("FindOneUser", mock.AnythingOfType("M")).Return(nil, errors.New("FindOneUser return error"))

	userUsecase := usecase.NewUserUsecase(us.mockUserRepo, us.mockAuthenticationHelper, us.mockHeaderHelper, us.mockFileOsHelper)
	_, err := userUsecase.VerifyCredential("username1", "password1")
	expectedError := domain.ErrInternalServerError.Error()

	assert.EqualErrorf(us.T(), err, expectedError, "Should have return %s but got %s", expectedError, err.Error())
}

func (us *UserUsecaseSuite) TestVerifyCredentialCompareHashAndPasswordError() {
	foundUsers := []bson.M{
		{"_id": "userid1", "username": "username1", "full_name": "fullname1", "password": "password1", "email": "email1@gmail.com", "profile_pictures": nil},
	}
	foundUser := domain.NewUser("userid1", "username1", "fullname1", "password1", "email1@gmail.com", nil)
	us.mockUserRepo.On("FindUser", mock.AnythingOfType("M")).Return(&foundUsers, nil)
	us.mockUserRepo.On("FindOneUser", mock.AnythingOfType("M")).Return(foundUser, nil)
	us.mockAuthenticationHelper.On("CompareHashAndPassword", mock.AnythingOfType("[]uint8"), mock.AnythingOfType("[]uint8")).Return(errors.New("CompareHashAndPassword return error"))

	userUsecase := usecase.NewUserUsecase(us.mockUserRepo, us.mockAuthenticationHelper, us.mockHeaderHelper, us.mockFileOsHelper)
	_, err := userUsecase.VerifyCredential("username1", "password1")
	expectedError := domain.ErrPasswordWrong.Error()

	assert.EqualErrorf(us.T(), err, expectedError, "Should have return %s but got %s", expectedError, err.Error())
}

func (us *UserUsecaseSuite) TestVerifyCredentialSuccessful() {
	foundUsers := []bson.M{
		{"_id": "userid1", "username": "username1", "full_name": "fullname1", "password": "password1", "email": "email1@gmail.com", "profile_pictures": nil},
	}
	foundUser := domain.NewUser("userid1", "username1", "fullname1", "password1", "email1@gmail.com", nil)
	us.mockUserRepo.On("FindUser", mock.AnythingOfType("M")).Return(&foundUsers, nil)
	us.mockUserRepo.On("FindOneUser", mock.AnythingOfType("M")).Return(foundUser, nil)
	us.mockAuthenticationHelper.On("CompareHashAndPassword", mock.AnythingOfType("[]uint8"), mock.AnythingOfType("[]uint8")).Return(nil)

	userUsecase := usecase.NewUserUsecase(us.mockUserRepo, us.mockAuthenticationHelper, us.mockHeaderHelper, us.mockFileOsHelper)
	_, err := userUsecase.VerifyCredential("username1", "password1")
	assert.NoErrorf(us.T(), err, "should have not returned error but got %s", err)
}
