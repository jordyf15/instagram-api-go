package usecase

import (
	"fmt"
	"instagram-go/domain"
	"sync"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type commentUsecase struct {
	sync.Mutex
	commentRepository domain.CommentRepository
	postRepository    domain.PostRepository
	likeRepository    domain.LikeRepository
	headerHelper      domain.IHeaderHelper
}

func NewCommentUsecase(commentRepository domain.CommentRepository, postRepository domain.PostRepository, likeRepository domain.LikeRepository, headerHelper domain.IHeaderHelper) domain.CommentUsecase {
	return &commentUsecase{
		commentRepository: commentRepository,
		postRepository:    postRepository,
		likeRepository:    likeRepository,
		headerHelper:      headerHelper,
	}
}

func (cu *commentUsecase) FindComments(postId string) (*[]domain.Comment, error) {
	filter := bson.M{"_id": postId}
	cu.Lock()
	queryResult, err := cu.postRepository.FindPosts(filter)
	cu.Unlock()
	if err != nil {
		return nil, domain.ErrInternalServerError
	}
	if len(*queryResult) == 0 {
		return nil, domain.ErrPostNotFound
	}
	filter = bson.M{"post_id": postId}
	cu.Lock()
	queryResult, err = cu.commentRepository.FindComments(filter)
	cu.Unlock()
	if err != nil {
		return nil, domain.ErrInternalServerError
	}
	var comments []domain.Comment
	for _, v := range *queryResult {
		id := fmt.Sprintf("%v", v["_id"])
		postId := fmt.Sprintf("%v", v["post_id"])
		userId := fmt.Sprintf("%v", v["user_id"])
		commentContent := fmt.Sprintf("%v", v["comment"])
		filter = bson.M{"resource_id": id, "resource_type": "comment"}
		cu.Lock()
		likeQueryResult, err := cu.likeRepository.FindLikes(filter)
		cu.Unlock()
		if err != nil {
			return nil, domain.ErrInternalServerError
		}
		likeCount := len(*likeQueryResult)
		createdDate := v["created_date"].(primitive.DateTime).Time()
		updatedDate := v["updated_date"].(primitive.DateTime).Time()
		if err != nil {
			return nil, err
		}
		comment := domain.NewComment(id, postId, userId, commentContent, likeCount, createdDate, updatedDate)
		comments = append(comments, *comment)
	}
	return &comments, nil
}

func (cu *commentUsecase) PostComment(comment *domain.Comment, tokenString string) error {
	userId, err := cu.headerHelper.GetUserIdFromToken(tokenString)
	if err != nil {
		return domain.ErrInternalServerError
	}
	newCommentId := "comment-" + uuid.NewString()
	filter := bson.M{"_id": comment.PostId}
	cu.Lock()
	queryResult, err := cu.postRepository.FindPosts(filter)
	cu.Unlock()
	if err != nil {
		return domain.ErrInternalServerError
	}
	if len(*queryResult) == 0 {
		return domain.ErrPostNotFound
	}
	comment.Id = newCommentId
	comment.UserId = userId
	comment.CreatedDate = time.Now()
	comment.UpdatedDate = comment.CreatedDate

	cu.Lock()
	err = cu.commentRepository.InsertComment(comment)
	cu.Unlock()
	if err != nil {
		return domain.ErrInternalServerError
	}
	return nil
}

func (cu *commentUsecase) PutComment(comment *domain.Comment, tokenString string) error {
	userId, err := cu.headerHelper.GetUserIdFromToken(tokenString)
	if err != nil {
		return domain.ErrInternalServerError
	}
	filter := bson.M{"_id": comment.Id}
	cu.Lock()
	queryResult, err := cu.commentRepository.FindComments(filter)
	cu.Unlock()
	if err != nil {
		return domain.ErrInternalServerError
	}
	if len(*queryResult) == 0 {
		return domain.ErrCommentNotFound
	}
	cu.Lock()
	willBeUpdatedComment, err := cu.commentRepository.FindOneComment(comment.Id)
	cu.Unlock()
	if err != nil {
		return domain.ErrInternalServerError
	}
	if willBeUpdatedComment.UserId != userId {
		return domain.ErrUnauthorizedCommentUpdate
	}

	cu.Lock()
	err = cu.commentRepository.UpdateComment(comment.Id, comment.Comment)
	cu.Unlock()
	if err != nil {
		return domain.ErrInternalServerError
	}
	return nil
}

func (cu *commentUsecase) DeleteComment(commentId string, tokenString string) error {
	userId, err := cu.headerHelper.GetUserIdFromToken(tokenString)
	if err != nil {
		return domain.ErrInternalServerError
	}
	filter := bson.M{"_id": commentId}
	cu.Lock()
	queryResult, err := cu.commentRepository.FindComments(filter)
	cu.Unlock()
	if err != nil {
		return domain.ErrInternalServerError
	}
	if len(*queryResult) == 0 {
		return domain.ErrCommentNotFound
	}

	cu.Lock()
	willBeDeletedComment, err := cu.commentRepository.FindOneComment(commentId)
	cu.Unlock()
	if err != nil {
		return domain.ErrInternalServerError
	}
	if willBeDeletedComment.UserId != userId {
		return domain.ErrUnauthorizedCommentDelete
	}

	cu.Lock()
	err = cu.commentRepository.DeleteComment(commentId)
	cu.Unlock()
	if err != nil {
		return domain.ErrInternalServerError
	}
	return nil
}
