package handlers

import (
	"encoding/json"
	"fmt"
	"instagram-go/models"
	"instagram-go/services"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type CommentHandlers struct {
	service services.CommentService
	sync.Mutex
}

func NewCommentHandlers(service services.CommentService) *CommentHandlers {
	return &CommentHandlers{
		service: service,
	}
}

func (ch *CommentHandlers) Comments(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		ch.getComments(w, r)
		return
	case "POST":
		ch.postComment(w, r)
		return
	}
}

func (ch *CommentHandlers) Comment(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "PUT":
		ch.putComment(w, r)
		return
	case "DELETE":
		ch.deleteComment(w, r)
		return
	}
}

func (ch *CommentHandlers) getComments(w http.ResponseWriter, r *http.Request) {
	urlParts := strings.Split(r.URL.String(), "/")
	postId := urlParts[2]

	ch.Lock()
	comments, err := ch.service.FindAllPostComment(postId)
	defer ch.Unlock()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	// response := models.DataResponseComments{models.DataComments{comments}}
	dataComments := *models.NewDataComments(comments)
	response := *models.NewDataResponseComments(dataComments)
	responseBytes, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(responseBytes)
}
func (ch *CommentHandlers) postComment(w http.ResponseWriter, r *http.Request) {
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
	newCommentId := "comment-" + uuid.NewString()

	urlParts := strings.Split(r.URL.String(), "/")
	postId := urlParts[2]

	bodyBytes, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	ct := r.Header.Get("content-type")
	if ct != r.Header.Get("content-type") {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte(fmt.Sprintf("need content-type 'application/json', but got '%s'", ct)))
		return
	}

	var comment models.Comment
	err = json.Unmarshal(bodyBytes, &comment)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	comment.Id = newCommentId
	comment.PostId = postId
	comment.UserId = userId
	comment.CreatedDate = time.Now()
	comment.UpdatedDate = comment.CreatedDate

	ch.Lock()
	err = ch.service.InsertComment(comment)
	defer ch.Unlock()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	// response := models.Message{"Comment successfully Created"}
	response := *models.NewMessage("Comment successfully Created")
	responseBytes, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(responseBytes)
}

func (ch *CommentHandlers) putComment(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	ct := r.Header.Get("content-type")
	if ct != r.Header.Get("content-type") {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte(fmt.Sprintf("need content-type 'application/json', but got '%s'", ct)))
		return
	}

	var comment models.Comment
	err = json.Unmarshal(bodyBytes, &comment)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

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

	ch.Lock()
	commentUserId, err := ch.service.GetCommentUserId(commentId)
	ch.Unlock()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	if commentUserId != userId {
		// response := models.Message{"You are not authorized to update this comment"}
		response := *models.NewMessage("You are not authorized to update this comment")
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
	ch.Lock()
	err = ch.service.UpdateComment(commentId, comment.Comment)
	defer ch.Unlock()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	// response := models.Message{"Comment successfully Updated"}
	response := *models.NewMessage("Comment successfully Updated")
	responseBytes, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(responseBytes)
}

func (ch *CommentHandlers) deleteComment(w http.ResponseWriter, r *http.Request) {
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

	ch.Lock()
	commentUserId, err := ch.service.GetCommentUserId(commentId)
	ch.Unlock()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	if commentUserId != userId {
		// response := models.Message{"You are not authorized to delete this comment"}
		response := *models.NewMessage("You are not authorized to delete this comment")
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

	ch.Lock()
	err = ch.service.DeleteComment(commentId)
	defer ch.Unlock()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	// response := models.Message{"Comment successfully Deleted"}
	response := *models.NewMessage("Comment successfully Deleted")
	responseBytes, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(responseBytes)
}
