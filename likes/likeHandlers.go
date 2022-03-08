package likes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type LikeHandlers struct {
	service LikeService
	sync.Mutex
}

func NewLikeHandlers(service LikeService) *LikeHandlers {
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
	exist, err := lh.service.isLikeExist(userId, postId, "post")
	lh.Unlock()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	if exist {
		response := message{"User have already liked this post"}
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
	newLike := Like{likeId, userId, postId, "post"}
	lh.Lock()
	err = lh.service.insertLike(newLike)
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
	likeUserId, err := lh.service.getLikeUserId(likeId)
	lh.Unlock()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	if likeUserId != userId {
		response := message{"You are not authorized"}
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
	err = lh.service.deleteLike(likeId)
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
	exist, err := lh.service.isLikeExist(userId, commentId, "comment")
	lh.Unlock()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	if exist {
		response := message{"User have already liked this comment"}
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
	newLike := Like{likeId, userId, commentId, "comment"}
	lh.Lock()
	err = lh.service.insertLike(newLike)
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
	likeUserId, err := lh.service.getLikeUserId(likeId)
	lh.Unlock()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	if likeUserId != userId {
		response := message{"You are not authorized"}
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
	err = lh.service.deleteLike(likeId)
	defer lh.Unlock()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
}

type message struct {
	Message string `json:"message"`
}
