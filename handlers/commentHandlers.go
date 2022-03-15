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
	commentService services.ICommentService
	postService    services.IPostService
	sync.Mutex
	commentHandlerHeader ICommentHandlerHeader
}

func NewCommentHandlers(commentService services.ICommentService, postService services.IPostService, commentHandlerHeader ICommentHandlerHeader) *CommentHandlers {
	if commentHandlerHeader == nil {
		commentHandlerHeader = newCommentHandlerHeader()
	}
	return &CommentHandlers{
		commentService:       commentService,
		postService:          postService,
		commentHandlerHeader: commentHandlerHeader,
	}
}

type ICommentHandlerHeader interface {
	getUserIdFromToken(string) (string, error)
}

type commentHandlerHeader struct {
}

func newCommentHandlerHeader() *commentHandlerHeader {
	return &commentHandlerHeader{}
}

func (chh *commentHandlerHeader) getUserIdFromToken(tokenString string) (string, error) {
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
	isPostExist, err := ch.postService.CheckIfPostExist(postId)
	ch.Unlock()
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
	if !isPostExist {
		response := models.NewMessage("Post does not exist")
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

	ch.Lock()
	comments, err := ch.commentService.FindAllPostComment(postId)
	defer ch.Unlock()
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

	userId, err := ch.commentHandlerHeader.getUserIdFromToken(tokenString)
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
	newCommentId := "comment-" + uuid.NewString()

	urlParts := strings.Split(r.URL.String(), "/")
	postId := urlParts[2]

	ch.Lock()
	isPostExist, err := ch.postService.CheckIfPostExist(postId)
	ch.Unlock()
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
	if !isPostExist {
		response := models.NewMessage("Post does not exist")
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

	if comment.Comment == "" {
		response := models.NewMessage("Comment must not be empty")
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
	comment.Id = newCommentId
	comment.PostId = postId
	comment.UserId = userId
	comment.CreatedDate = time.Now()
	comment.UpdatedDate = comment.CreatedDate

	ch.Lock()
	err = ch.commentService.InsertComment(comment)
	defer ch.Unlock()
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

	ct := r.Header.Get("content-type")
	if ct != r.Header.Get("content-type") {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte(fmt.Sprintf("need content-type 'application/json', but got '%s'", ct)))
		return
	}

	var comment models.Comment
	err = json.Unmarshal(bodyBytes, &comment)
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

	urlParts := strings.Split(r.URL.String(), "/")
	commentId := urlParts[4]

	tokenString := r.Header.Get("Authorization")
	userId, err := ch.commentHandlerHeader.getUserIdFromToken(tokenString)
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
	ch.Lock()
	isCommentExist, err := ch.commentService.CheckIfCommentExist(commentId)
	ch.Unlock()
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
	if !isCommentExist {
		response := models.NewMessage("Comment does not exist")
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

	ch.Lock()
	commentUserId, err := ch.commentService.GetCommentUserId(commentId)
	ch.Unlock()

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
	if commentUserId != userId {
		response := *models.NewMessage("User is not authorized to update this comment")
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
	if comment.Comment == "" {
		response := models.NewMessage("Comment must not be empty")
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
	ch.Lock()
	err = ch.commentService.UpdateComment(commentId, comment.Comment)
	defer ch.Unlock()
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
	userId, err := ch.commentHandlerHeader.getUserIdFromToken(tokenString)
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
	ch.Lock()
	isCommentExist, err := ch.commentService.CheckIfCommentExist(commentId)
	ch.Unlock()
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
	if !isCommentExist {
		response := models.NewMessage("Comment does not exist")
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

	ch.Lock()
	commentUserId, err := ch.commentService.GetCommentUserId(commentId)
	ch.Unlock()

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

	if commentUserId != userId {
		response := *models.NewMessage("User is not authorized to delete this comment")
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
	err = ch.commentService.DeleteComment(commentId)
	defer ch.Unlock()
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
