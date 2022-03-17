package http

import (
	"encoding/json"
	"instagram-go/domain"
	"net/http"
	"strings"
)

type LikeHandler struct {
	likeUsecase domain.LikeUsecase
}

func NewLikeHandler(likeUsecase domain.LikeUsecase) *LikeHandler {
	return &LikeHandler{
		likeUsecase: likeUsecase,
	}
}

func (lh *LikeHandler) PostLikePost(w http.ResponseWriter, r *http.Request) {
	urlParts := strings.Split(r.URL.String(), "/")
	postId := urlParts[2]
	tokenString := r.Header.Get("Authorization")
	err := lh.likeUsecase.InsertPostLike(postId, tokenString)
	if err != nil {
		response := domain.NewMessage(err.Error())
		responseBytes, errMarshal := json.Marshal(response)
		if errMarshal != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errMarshal.Error()))
			return
		}
		w.WriteHeader(likeGetStatusCode(err))
		w.Write(responseBytes)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (lh *LikeHandler) DeleteLikePost(w http.ResponseWriter, r *http.Request) {
	urlParts := strings.Split(r.URL.String(), "/")
	likeId := urlParts[4]
	tokenString := r.Header.Get("Authorization")
	err := lh.likeUsecase.DeletePostLike(likeId, tokenString)
	if err != nil {
		response := domain.NewMessage(err.Error())
		responseBytes, errMarshal := json.Marshal(response)
		if errMarshal != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errMarshal.Error()))
			return
		}
		w.WriteHeader(likeGetStatusCode(err))
		w.Write(responseBytes)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (lh *LikeHandler) PostCommentLike(w http.ResponseWriter, r *http.Request) {
	urlParts := strings.Split(r.URL.String(), "/")
	commentId := urlParts[4]
	tokenString := r.Header.Get("Authorization")

	err := lh.likeUsecase.InsertCommentLike(commentId, tokenString)
	if err != nil {
		response := domain.NewMessage(err.Error())
		responseBytes, errMarshal := json.Marshal(response)
		if errMarshal != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errMarshal.Error()))
			return
		}
		w.WriteHeader(likeGetStatusCode(err))
		w.Write(responseBytes)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (lh *LikeHandler) DeleteCommentLike(w http.ResponseWriter, r *http.Request) {
	urlParts := strings.Split(r.URL.String(), "/")
	likeId := urlParts[6]
	tokenString := r.Header.Get("Authorization")

	err := lh.likeUsecase.DeleteCommentLike(likeId, tokenString)
	if err != nil {
		response := domain.NewMessage(err.Error())
		responseBytes, errMarshal := json.Marshal(response)
		if errMarshal != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errMarshal.Error()))
			return
		}
		w.WriteHeader(likeGetStatusCode(err))
		w.Write(responseBytes)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func likeGetStatusCode(err error) int {
	switch err {
	case domain.ErrInternalServerError:
		return http.StatusInternalServerError
	case domain.ErrPostNotFound, domain.ErrLikeNotFound, domain.ErrCommentNotFound:
		return http.StatusNotFound
	case domain.ErrPostLikeConflict, domain.ErrCommentLikeConflict:
		return http.StatusConflict
	case domain.ErrUnauthorizedLikeDelete:
		return http.StatusUnauthorized
	}
	return http.StatusOK
}
