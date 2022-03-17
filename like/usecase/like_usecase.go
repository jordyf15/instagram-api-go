package usecase

import (
	"instagram-go/domain"
	"sync"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

type likeUsecase struct {
	sync.Mutex
	headerHelper      domain.IHeaderHelper
	postRepository    domain.PostRepository
	likeRepository    domain.LikeRepository
	commentRepository domain.CommentRepository
}

func NewLikeUsecase(likeRepository domain.LikeRepository, postRepository domain.PostRepository, commentRepository domain.CommentRepository, headerHelper domain.IHeaderHelper) *likeUsecase {
	return &likeUsecase{
		likeRepository:    likeRepository,
		postRepository:    postRepository,
		commentRepository: commentRepository,
		headerHelper:      headerHelper,
	}
}

func (lu *likeUsecase) InsertPostLike(postId string, tokenString string) error {
	userId, err := lu.headerHelper.GetUserIdFromToken(tokenString)
	if err != nil {
		return domain.ErrInternalServerError
	}
	filter := bson.M{"_id": postId}
	lu.Lock()
	queryResult, err := lu.postRepository.FindPosts(filter)
	lu.Unlock()
	if err != nil {
		return domain.ErrInternalServerError
	}
	if len(*queryResult) == 0 {
		return domain.ErrPostNotFound
	}

	filter = bson.M{"user_id": userId, "resource_id": postId, "resource_type": "post"}
	lu.Lock()
	queryResult, err = lu.likeRepository.FindLikes(filter)
	lu.Unlock()
	if err != nil {
		return domain.ErrInternalServerError
	}
	if len(*queryResult) > 0 {
		return domain.ErrPostLikeConflict
	}
	likeId := "like-" + uuid.NewString()
	like := domain.NewLike(likeId, userId, postId, "post")
	lu.Lock()
	err = lu.likeRepository.InsertLike(like)
	lu.Unlock()
	if err != nil {
		return domain.ErrInternalServerError
	}
	return nil
}

func (lu *likeUsecase) DeletePostLike(likeId string, tokenString string) error {
	userId, err := lu.headerHelper.GetUserIdFromToken(tokenString)
	if err != nil {
		return domain.ErrInternalServerError
	}
	filter := bson.M{"_id": likeId}
	lu.Lock()
	queryResult, err := lu.likeRepository.FindLikes(filter)
	lu.Unlock()
	if err != nil {
		return domain.ErrInternalServerError
	}
	if len(*queryResult) == 0 {
		return domain.ErrLikeNotFound
	}

	lu.Lock()
	like, err := lu.likeRepository.FindOneLike(likeId)
	lu.Unlock()
	if err != nil {
		return domain.ErrInternalServerError
	}
	if like.UserId != userId {
		return domain.ErrUnauthorizedLikeDelete
	}

	lu.Lock()
	err = lu.likeRepository.DeleteLike(likeId)
	lu.Unlock()
	if err != nil {
		return domain.ErrInternalServerError
	}

	return nil
}

func (lu *likeUsecase) InsertCommentLike(commentId string, tokenString string) error {
	userId, err := lu.headerHelper.GetUserIdFromToken(tokenString)
	if err != nil {
		return domain.ErrInternalServerError
	}

	filter := bson.M{"_id": commentId}
	lu.Lock()
	queryResult, err := lu.commentRepository.FindComments(filter)
	lu.Unlock()
	if err != nil {
		return domain.ErrInternalServerError
	}
	if len(*queryResult) == 0 {
		return domain.ErrCommentNotFound
	}

	filter = bson.M{"user_id": userId, "resource_id": commentId, "resource_type": "comment"}
	lu.Lock()
	queryResult, err = lu.likeRepository.FindLikes(filter)
	lu.Unlock()
	if err != nil {
		return domain.ErrInternalServerError
	}
	if len(*queryResult) > 0 {
		return domain.ErrCommentLikeConflict
	}

	likeId := "like-" + uuid.NewString()
	like := domain.NewLike(likeId, userId, commentId, "comment")
	lu.Lock()
	err = lu.likeRepository.InsertLike(like)
	lu.Unlock()
	if err != nil {
		return domain.ErrInternalServerError
	}
	return nil
}

func (lu *likeUsecase) DeleteCommentLike(likeId string, tokenString string) error {
	userId, err := lu.headerHelper.GetUserIdFromToken(tokenString)
	if err != nil {
		return domain.ErrInternalServerError
	}
	filter := bson.M{"_id": likeId}
	lu.Lock()
	queryResult, err := lu.likeRepository.FindLikes(filter)
	lu.Unlock()
	if err != nil {
		return domain.ErrInternalServerError
	}
	if len(*queryResult) == 0 {
		return domain.ErrLikeNotFound
	}

	lu.Lock()
	like, err := lu.likeRepository.FindOneLike(likeId)
	lu.Unlock()

	if err != nil {
		return domain.ErrInternalServerError
	}
	if like.UserId != userId {
		return domain.ErrUnauthorizedLikeDelete
	}

	lu.Lock()
	err = lu.likeRepository.DeleteLike(likeId)
	lu.Unlock()
	if err != nil {
		return domain.ErrInternalServerError
	}
	return nil
}
