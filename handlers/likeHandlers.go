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
	service *services.LikeService
	sync.Mutex
	likeHandlerHeader ILikeHandlerHeader
}

func NewLikeHandlers(service *services.LikeService, likeHandlerHeader ILikeHandlerHeader) *LikeHandlers {
	if likeHandlerHeader == nil {
		likeHandlerHeader = newLikeHandlerHeader()
	}
	return &LikeHandlers{
		service:           service,
		likeHandlerHeader: likeHandlerHeader,
	}
}

type ILikeHandlerHeader interface {
	getUserIdFromToken(string) (string, error)
}

type likeHandlerHeader struct {
}

func newLikeHandlerHeader() *likeHandlerHeader {
	return &likeHandlerHeader{}
}

func (lhh *likeHandlerHeader) getUserIdFromToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	if err != nil {
		return "", err
	}
	claims := token.Claims.(jwt.MapClaims)
	userId := fmt.Sprintf("%v", claims["user_id"])
	return userId, nil
}

func (lh *LikeHandlers) PostPostLikeHandler(w http.ResponseWriter, r *http.Request) {
	urlParts := strings.Split(r.URL.String(), "/")
	fmt.Println("asdassd")
	postId := urlParts[2]
	tokenString := r.Header.Get("Authorization")
	// token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
	// 	return []byte("secret"), nil
	// })
	// if err != nil {
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	w.Write([]byte(err.Error()))
	// 	return
	// }
	// claims := token.Claims.(jwt.MapClaims)
	// userIdToken := claims["user_id"]
	// userId := fmt.Sprintf("%v", userIdToken)
	userId, err := lh.likeHandlerHeader.getUserIdFromToken(tokenString)
	if err != nil {
		response := models.NewMessage("An error has occured in our server")
		responseBytes, err := json.Marshal(response)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(responseBytes)
		return
	}

	lh.Lock()
	exist, err := lh.service.IsLikeExist(userId, postId, "post")
	lh.Unlock()
	if err != nil {
		response := models.NewMessage("An error has occured in our server")
		responseBytes, err := json.Marshal(response)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(responseBytes)
		return
	}
	fmt.Println(exist)
	if exist {
		// response := models.Message{"User have already liked this post"}
		response := *models.NewMessage("User have already liked this post")
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
	// newLike := models.Like{likeId, userId, postId, "post"}
	newLike := *models.NewLike(likeId, userId, postId, "post")
	lh.Lock()
	err = lh.service.InsertLike(newLike)
	defer lh.Unlock()
	if err != nil {
		response := models.NewMessage("An error has occured in our server")
		responseBytes, err := json.Marshal(response)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(responseBytes)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (lh *LikeHandlers) DeletePostLikeHandler(w http.ResponseWriter, r *http.Request) {
	urlParts := strings.Split(r.URL.String(), "/")
	likeId := urlParts[4]
	tokenString := r.Header.Get("Authorization")
	// token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
	// 	return []byte("secret"), nil
	// })
	// if err != nil {
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	w.Write([]byte(err.Error()))
	// 	return
	// }
	// claims := token.Claims.(jwt.MapClaims)
	// userIdToken := claims["user_id"]
	// userId := fmt.Sprintf("%v", userIdToken)
	userId, err := lh.likeHandlerHeader.getUserIdFromToken(tokenString)
	if err != nil {
		response := models.NewMessage("An error has occured in our server")
		responseBytes, err := json.Marshal(response)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(responseBytes)
		return
	}
	lh.Lock()
	likeUserId, err := lh.service.GetLikeUserId(likeId)
	lh.Unlock()
	if err != nil {
		response := models.NewMessage("An error has occured in our server")
		responseBytes, err := json.Marshal(response)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(responseBytes)
		return
	}
	if likeUserId != userId {
		// response := models.Message{"You are not authorized"}
		response := *models.NewMessage("You are not authorized")
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
		response := models.NewMessage("An error has occured in our server")
		responseBytes, err := json.Marshal(response)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(responseBytes)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (lh *LikeHandlers) PostCommentLikeHandler(w http.ResponseWriter, r *http.Request) {
	urlParts := strings.Split(r.URL.String(), "/")
	commentId := urlParts[4]
	tokenString := r.Header.Get("Authorization")
	// token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
	// 	return []byte("secret"), nil
	// })
	// if err != nil {
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	w.Write([]byte(err.Error()))
	// 	return
	// }
	// claims := token.Claims.(jwt.MapClaims)
	// userIdToken := claims["user_id"]
	// userId := fmt.Sprintf("%v", userIdToken)
	userId, err := lh.likeHandlerHeader.getUserIdFromToken(tokenString)
	if err != nil {
		response := models.NewMessage("An error has occured in our server")
		responseBytes, err := json.Marshal(response)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(responseBytes)
		return
	}
	lh.Lock()
	exist, err := lh.service.IsLikeExist(userId, commentId, "comment")
	lh.Unlock()
	if err != nil {
		response := models.NewMessage("An error has occured in our server")
		responseBytes, err := json.Marshal(response)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(responseBytes)
		return
	}
	if exist {
		// response := models.Message{"User have already liked this comment"}
		response := *models.NewMessage("User have already liked this comment")
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
	// newLike := models.Like{likeId, userId, commentId, "comment"}
	newLike := *models.NewLike(likeId, userId, commentId, "comment")
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
	// token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
	// 	return []byte("secret"), nil
	// })
	// if err != nil {
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	w.Write([]byte(err.Error()))
	// 	return
	// }
	// claims := token.Claims.(jwt.MapClaims)
	// userIdToken := claims["user_id"]
	// userId := fmt.Sprintf("%v", userIdToken)
	userId, err := lh.likeHandlerHeader.getUserIdFromToken(tokenString)
	if err != nil {
		response := models.NewMessage("An error has occured in our server")
		responseBytes, err := json.Marshal(response)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(responseBytes)
		return
	}
	lh.Lock()
	likeUserId, err := lh.service.GetLikeUserId(likeId)
	lh.Unlock()
	if err != nil {
		response := models.NewMessage("An error has occured in our server")
		responseBytes, err := json.Marshal(response)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(responseBytes)
		return
	}
	if likeUserId != userId {
		// response := models.Message{"You are not authorized"}
		response := *models.NewMessage("You are not authorized")
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
		response := models.NewMessage("An error has occured in our server")
		responseBytes, err := json.Marshal(response)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(responseBytes)
		return
	}
	w.WriteHeader(http.StatusOK)
}
