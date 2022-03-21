package usecase

import (
	"fmt"
	"instagram-go/domain"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type postUsecase struct {
	postRepository domain.PostRepository
	likeRepository domain.LikeRepository
	headerHelper   domain.IHeaderHelper
	fileOsHelper   domain.IFileOsHelper
	sync.Mutex
}

func NewPostUseCase(postRepository domain.PostRepository, likeRepository domain.LikeRepository, headerHelper domain.IHeaderHelper, fileOsHelper domain.IFileOsHelper) domain.PostUsecase {
	return &postUsecase{
		postRepository: postRepository,
		likeRepository: likeRepository,
		fileOsHelper:   fileOsHelper,
		headerHelper:   headerHelper,
	}
}

func (pu *postUsecase) InsertPost(post *domain.Post, tokenString string, visualMedias []*multipart.FileHeader) error {
	userId, err := pu.headerHelper.GetUserIdFromToken(tokenString)
	if err != nil {
		return domain.ErrInternalServerError
	}
	post.UserId = userId
	post.Id = "post-" + uuid.NewString()
	newPath := filepath.Join(".", "visual_medias")
	err = pu.fileOsHelper.MkDirAll(newPath, os.ModePerm)
	if err != nil {
		return domain.ErrInternalServerError
	}
	for k, v := range visualMedias {
		visualMedia, err := v.Open()
		if err != nil {
			return domain.ErrInternalServerError
		}
		defer visualMedia.Close()
		fileNameParts := strings.Split(v.Filename, ".")
		extension := fileNameParts[len(fileNameParts)-1]
		visualMediaUrl := "./visual_medias/" + post.Id + strconv.Itoa(k) + "." + extension
		out, err := pu.fileOsHelper.Create(visualMediaUrl)
		if err != nil {
			return domain.ErrInternalServerError
		}
		defer out.Close()
		_, err = pu.fileOsHelper.Copy(out, visualMedia)
		if err != nil {
			return domain.ErrInternalServerError
		}
		post.VisualMediaUrls = append(post.VisualMediaUrls, visualMediaUrl)
	}
	post.CreatedDate = time.Now()
	post.UpdatedDate = post.CreatedDate

	pu.Lock()
	err = pu.postRepository.InsertPost(post)
	pu.Unlock()

	if err != nil {
		return domain.ErrInternalServerError
	}
	return nil
}

func (pu *postUsecase) FindPosts() (*[]domain.Post, error) {
	filter := bson.M{}
	pu.Lock()
	queryResult, err := pu.postRepository.FindPosts(filter)
	pu.Unlock()
	if err != nil {
		return nil, domain.ErrInternalServerError
	}
	var posts []domain.Post
	for _, v := range *queryResult {
		id := fmt.Sprintf("%v", v["_id"])
		userId := fmt.Sprintf("%v", v["user_id"])
		var visualMediaUrls []string

		if visualMediaUrlsPrimitive, ok := v["visual_media_urls"].(primitive.A); ok {
			visualMediaUrlsInterface := []interface{}(visualMediaUrlsPrimitive)
			visualMediaUrls = make([]string, len(visualMediaUrlsInterface))
			for i, url := range visualMediaUrlsInterface {
				visualMediaUrls[i] = url.(string)
			}
		}
		caption := fmt.Sprintf("%v", v["caption"])
		createdDate := v["created_date"].(primitive.DateTime).Time()
		updatedDate := v["updated_date"].(primitive.DateTime).Time()

		filter := bson.M{"resource_id": id, "resource_type": "post"}
		pu.Lock()
		likes, err := pu.likeRepository.FindLikes(filter)
		pu.Unlock()
		if err != nil {
			return nil, domain.ErrInternalServerError
		}
		likeCount := len(*likes)
		post := domain.NewPost(id, userId, visualMediaUrls, caption, likeCount, createdDate, updatedDate)
		posts = append(posts, *post)
	}
	return &posts, nil
}

func (pu *postUsecase) UpdatePost(updatedPostId string, newCaption string, tokenString string) error {
	userId, err := pu.headerHelper.GetUserIdFromToken(tokenString)
	if err != nil {
		return domain.ErrInternalServerError
	}

	filter := bson.M{"_id": updatedPostId}
	pu.Lock()
	queryResult, err := pu.postRepository.FindPosts(filter)
	pu.Unlock()
	if err != nil {
		return domain.ErrInternalServerError
	}
	if len(*queryResult) == 0 {
		return domain.ErrPostNotFound
	}

	pu.Lock()
	post, err := pu.postRepository.FindOnePost(updatedPostId)
	pu.Unlock()
	if err != nil {
		return domain.ErrInternalServerError
	}
	if post.UserId != userId {
		return domain.ErrUnauthorizedPostUpdate
	}

	pu.Lock()
	err = pu.postRepository.UpdatePost(updatedPostId, newCaption)
	pu.Unlock()
	if err != nil {
		return domain.ErrInternalServerError
	}
	return nil
}

func (pu *postUsecase) DeletePost(deletedPostId string, tokenString string) error {
	userId, err := pu.headerHelper.GetUserIdFromToken(tokenString)
	if err != nil {
		return domain.ErrInternalServerError
	}
	filter := bson.M{"_id": deletedPostId}
	pu.Lock()
	queryResult, err := pu.postRepository.FindPosts(filter)
	pu.Unlock()
	if err != nil {
		return domain.ErrInternalServerError
	}
	if len(*queryResult) == 0 {
		return domain.ErrPostNotFound
	}

	pu.Lock()
	post, err := pu.postRepository.FindOnePost(deletedPostId)
	pu.Unlock()
	if err != nil {
		return domain.ErrInternalServerError
	}
	if post.UserId != userId {
		return domain.ErrUnauthorizedPostDelete
	}

	pu.Lock()
	err = pu.postRepository.DeletePost(deletedPostId)
	pu.Unlock()
	if err != nil {
		return domain.ErrInternalServerError
	}
	return nil
}
