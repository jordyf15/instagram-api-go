package handlers

import (
	"encoding/json"
	"fmt"
	"instagram-go/models"
	"instagram-go/services"
	"net/http"
	"strings"
	"sync"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type LikeHandlers struct {
	service services.LikeService
	sync.Mutex
}

func NewLikeHandlers(service services.LikeService) *LikeHandlers {
	return &LikeHandlers{
		service: service,
	}
}

func (lh *LikeHandlers) PostPostLikeHandler(w http.ResponseWriter, r *http.Request) {
	urlParts := strings.Split(r.URL.String(), "/")
	postId := urlParts[2]
	tokenString := r.Header.Get("Authorization")
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	claims := token.Claims.(jwt.MapClaims)
	userIdToken := claims["user_id"]
	userId := fmt.Sprintf("%v", userIdToken)

	lh.Lock()
	exist, err := lh.service.IsLikeExist(userId, postId, "post")
	lh.Unlock()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	if exist {
		response := models.Message{"User have already liked this post"}
		responseBytes, err := json.Marshal(response)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write(responseBytes)
		return
	}

	likeId := "like-" + uuid.NewString()
	newLike := models.Like{likeId, userId, postId, "post"}
	lh.Lock()
	err = lh.service.InsertLike(newLike)
	defer lh.Unlock()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (lh *LikeHandlers) DeletePostLikeHandler(w http.ResponseWriter, r *http.Request) {
	urlParts := strings.Split(r.URL.String(), "/")
	likeId := urlParts[4]
	tokenString := r.Header.Get("Authorization")
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	claims := token.Claims.(jwt.MapClaims)
	userIdToken := claims["user_id"]
	userId := fmt.Sprintf("%v", userIdToken)

	lh.Lock()
	likeUserId, err := lh.service.GetLikeUserId(likeId)
	lh.Unlock()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	if likeUserId != userId {
		response := models.Message{"You are not authorized"}
		responseBytes, err := json.Marshal(response)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(responseBytes)
		return
	}
	lh.Lock()
	err = lh.service.DeleteLike(likeId)
	defer lh.Unlock()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (lh *LikeHandlers) PostCommentLikeHandler(w http.ResponseWriter, r *http.Request) {
	urlParts := strings.Split(r.URL.String(), "/")
	commentId := urlParts[4]
	tokenString := r.Header.Get("Authorization")
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	claims := token.Claims.(jwt.MapClaims)
	userIdToken := claims["user_id"]
	userId := fmt.Sprintf("%v", userIdToken)

	lh.Lock()
	exist, err := lh.service.IsLikeExist(userId, commentId, "comment")
	lh.Unlock()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	if exist {
		response := models.Message{"User have already liked this comment"}
		responseBytes, err := json.Marshal(response)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write(responseBytes)
		return
	}

	likeId := "like-" + uuid.NewString()
	newLike := models.Like{likeId, userId, commentId, "comment"}
	lh.Lock()
	err = lh.service.InsertLike(newLike)
	defer lh.Unlock()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (lh *LikeHandlers) DeleteCommentLikeHandler(w http.ResponseWriter, r *http.Request) {
	urlParts := strings.Split(r.URL.String(), "/")
	likeId := urlParts[6]
	tokenString := r.Header.Get("Authorization")
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	claims := token.Claims.(jwt.MapClaims)
	userIdToken := claims["user_id"]
	userId := fmt.Sprintf("%v", userIdToken)

	lh.Lock()
	likeUserId, err := lh.service.GetLikeUserId(likeId)
	lh.Unlock()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	if likeUserId != userId {
		response := models.Message{"You are not authorized"}
		responseBytes, err := json.Marshal(response)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(responseBytes)
		return
	}
	lh.Lock()
	err = lh.service.DeleteLike(likeId)
	defer lh.Unlock()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
}
