package usecase

import (
	"instagram-go/domain"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

type userUsecase struct {
	userRepository domain.UserRepository
	sync.Mutex
	fileOsHelper         domain.IFileOsHelper
	headerHelper         domain.IHeaderHelper
	authenticationHelper domain.IAuthenticationHelper
}

func NewUserUsecase(userRepository domain.UserRepository, authenticationHelper domain.IAuthenticationHelper, headerHelper domain.IHeaderHelper, fileOsHelper domain.IFileOsHelper) *userUsecase {
	return &userUsecase{
		userRepository:       userRepository,
		headerHelper:         headerHelper,
		fileOsHelper:         fileOsHelper,
		authenticationHelper: authenticationHelper,
	}
}

func (uu *userUsecase) InsertUser(user *domain.User) error {
	filter := bson.M{"username": user.Username}
	uu.Lock()
	findUserQueryResult, err := uu.userRepository.FindUser(filter)
	uu.Unlock()
	if err != nil {
		return domain.ErrInternalServerError
	}
	if len(*findUserQueryResult) > 0 {
		return domain.ErrUsernameConflict
	}
	user.Id = "user-" + uuid.NewString()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		return domain.ErrInternalServerError
	}
	user.Password = string(hashedPassword[:])

	uu.Lock()
	err = uu.userRepository.InsertUser(user)
	uu.Unlock()
	if err != nil {
		return domain.ErrInternalServerError
	}
	return nil
}

func (uu *userUsecase) UpdateUser(user *domain.User, tokenString string, profilePictureFile multipart.File) error {
	var updatedUser *domain.User
	newpath := filepath.Join(".", "profile_pictures")
	err := uu.fileOsHelper.MkDirAll(newpath, os.ModePerm)
	if err != nil {
		return domain.ErrInternalServerError
	}
	userIdToken, err := uu.headerHelper.GetUserIdFromToken(tokenString)
	if err != nil {
		return domain.ErrInternalServerError
	}
	filter := bson.M{"_id": user.Id}

	uu.Lock()
	findUserQueryResult, err := uu.userRepository.FindUser(filter)
	uu.Unlock()
	if err != nil {
		return domain.ErrInternalServerError
	}
	if len(*findUserQueryResult) == 0 {
		return domain.ErrUserNotFound
	}
	if userIdToken != user.Id {
		return domain.ErrUnauthorizedUserUpdate
	}
	if profilePictureFile != nil {
		fileHeader := make([]byte, 512)
		if _, err := profilePictureFile.Read(fileHeader); err != nil {
			return domain.ErrInternalServerError
		}
		if _, err := profilePictureFile.Seek(0, 0); err != nil {
			return domain.ErrInternalServerError
		}
		fileType := http.DetectContentType(fileHeader)
		originalProfilePicture, _, err := uu.fileOsHelper.DecodeImage(profilePictureFile)
		if err != nil {
			return domain.ErrInternalServerError
		}

		smallProfilePictureUrl, err := uu.fileOsHelper.ResizeAndSaveFileToLocale("small", originalProfilePicture, user.Id, fileType)
		if err != nil {
			return domain.ErrInternalServerError
		}

		averageProfilePictureUrl, err := uu.fileOsHelper.ResizeAndSaveFileToLocale("average", originalProfilePicture, user.Id, fileType)
		if err != nil {
			return domain.ErrInternalServerError
		}

		largeProfilePictureUrl, err := uu.fileOsHelper.ResizeAndSaveFileToLocale("large", originalProfilePicture, user.Id, fileType)
		if err != nil {
			return domain.ErrInternalServerError
		}

		smallProfilePicture := domain.NewProfilePicture("small", "150 x 150 px", smallProfilePictureUrl)
		averageProfilePicture := domain.NewProfilePicture("average", "400 x 400 ox", averageProfilePictureUrl)
		largeProfilePicture := domain.NewProfilePicture("large", "800 x 800 px", largeProfilePictureUrl)

		updatedUser = domain.NewUser(user.Id, user.Username, user.Fullname, user.Password, user.Email, []domain.ProfilePicture{*smallProfilePicture, *averageProfilePicture, *largeProfilePicture})
	} else {
		updatedUser = domain.NewUser(user.Id, user.Username, user.Fullname, user.Password, user.Email, nil)
	}

	filter = bson.M{"_id": user.Id}
	uu.Lock()
	oldUser, err := uu.userRepository.FindOneUser(filter)
	uu.Unlock()
	if err != nil {
		return domain.ErrInternalServerError
	}

	if updatedUser.Username == "" {
		updatedUser.Username = oldUser.Username
	}
	if updatedUser.Email == "" {
		updatedUser.Email = oldUser.Email
	}
	if updatedUser.Fullname == "" {
		updatedUser.Fullname = oldUser.Fullname
	}
	if updatedUser.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
		if err != nil {
			return domain.ErrInternalServerError
		}
		updatedUser.Password = string(hashedPassword[:])
	} else {
		updatedUser.Password = oldUser.Password
	}
	if updatedUser.ProfilePictures == nil {
		updatedUser.ProfilePictures = oldUser.ProfilePictures
	}

	uu.Lock()
	err = uu.userRepository.UpdateUser(updatedUser)
	uu.Unlock()
	if err != nil {
		return domain.ErrInternalServerError
	}
	return nil
}

func (uu *userUsecase) VerifyCredential(username string, password string) (string, error) {
	filter := bson.M{"username": username}
	uu.Lock()
	findUserQueryResult, err := uu.userRepository.FindUser(filter)
	uu.Unlock()
	if err != nil {
		return "", domain.ErrInternalServerError
	}
	if len(*findUserQueryResult) == 0 {
		return "", domain.ErrUserNotFound
	}

	uu.Lock()
	user, err := uu.userRepository.FindOneUser(filter)
	uu.Unlock()
	if err != nil {
		return "", domain.ErrInternalServerError
	}

	err = uu.authenticationHelper.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", domain.ErrPasswordWrong
	}

	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_id"] = user.Id
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()
	sign := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := sign.SignedString([]byte("secret"))
	if err != nil {
		return "", domain.ErrInternalServerError
	}
	return token, nil
}
