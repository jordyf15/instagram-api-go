package http

import (
	"encoding/json"
	"instagram-go/domain"
	"io/ioutil"
	"net/http"
	"strings"
)

type CommentHandler struct {
	commentUsecase domain.CommentUsecase
}

func NewCommentHandler(commentUsecase domain.CommentUsecase) *CommentHandler {
	return &CommentHandler{
		commentUsecase: commentUsecase,
	}
}

func (ch *CommentHandler) Comments(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		ch.getComments(w, r)
		return
	case "POST":
		ch.postComment(w, r)
		return
	}
}

func (ch *CommentHandler) Comment(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "PUT":
		ch.putComment(w, r)
		return
	case "DELETE":
		ch.deleteComment(w, r)
		return
	}
}

func (ch *CommentHandler) getComments(w http.ResponseWriter, r *http.Request) {
	urlParts := strings.Split(r.URL.String(), "/")
	postId := urlParts[2]
	comments, err := ch.commentUsecase.FindComments(postId)
	if err != nil {
		response := domain.NewMessage(err.Error())
		responseBytes, errMarshal := json.Marshal(response)
		if errMarshal != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errMarshal.Error()))
			return
		}
		w.WriteHeader(commentGetStatusCode(err))
		w.Write(responseBytes)
		return
	}
	dataComments := domain.NewDataComments(comments)
	response := domain.NewDataResponseComments(dataComments)
	responseBytes, errMarshal := json.Marshal(response)
	if errMarshal != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(errMarshal.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(responseBytes)
}

func (ch *CommentHandler) postComment(w http.ResponseWriter, r *http.Request) {
	tokenString := r.Header.Get("Authorization")
	urlParts := strings.Split(r.URL.String(), "/")
	postId := urlParts[2]
	bodyBytes, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	var comment domain.Comment
	err = json.Unmarshal(bodyBytes, &comment)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	if comment.Comment == "" {
		response := domain.NewMessage(domain.ErrMissingCommentInput.Error())
		responseBytes, errMarshal := json.Marshal(response)
		if errMarshal != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errMarshal.Error()))
			return
		}
		w.WriteHeader(commentGetStatusCode(domain.ErrMissingCommentInput))
		w.Write(responseBytes)
		return
	}
	comment.PostId = postId
	err = ch.commentUsecase.PostComment(&comment, tokenString)
	if err != nil {
		response := domain.NewMessage(err.Error())
		responseBytes, errMarshal := json.Marshal(response)
		if errMarshal != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errMarshal.Error()))
			return
		}
		w.WriteHeader(commentGetStatusCode(err))
		w.Write(responseBytes)
		return
	}
	response := domain.NewMessage("Comment successfully Created")
	responseBytes, errMarshal := json.Marshal(response)
	if errMarshal != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(errMarshal.Error()))
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(responseBytes)
}

func (ch *CommentHandler) putComment(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		response := domain.NewMessage(domain.ErrInternalServerError.Error())
		responseBytes, errMarshal := json.Marshal(response)
		if errMarshal != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errMarshal.Error()))
			return
		}
		w.WriteHeader(commentGetStatusCode(err))
		w.Write(responseBytes)
		return
	}
	var comment domain.Comment
	err = json.Unmarshal(bodyBytes, &comment)
	if err != nil {
		response := domain.NewMessage(domain.ErrInternalServerError.Error())
		responseBytes, errMarshal := json.Marshal(response)
		if errMarshal != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errMarshal.Error()))
			return
		}
		w.WriteHeader(commentGetStatusCode(err))
		w.Write(responseBytes)
		return
	}
	urlParts := strings.Split(r.URL.String(), "/")
	commentId := urlParts[4]
	comment.Id = commentId
	tokenString := r.Header.Get("Authorization")
	if comment.Comment == "" {
		response := domain.NewMessage(domain.ErrMissingCommentInput.Error())
		responseBytes, errMarshal := json.Marshal(response)
		if errMarshal != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errMarshal.Error()))
			return
		}
		w.WriteHeader(commentGetStatusCode(domain.ErrMissingCommentInput))
		w.Write(responseBytes)
		return
	}
	err = ch.commentUsecase.PutComment(&comment, tokenString)
	if err != nil {
		response := domain.NewMessage(err.Error())
		responseBytes, errMarshal := json.Marshal(response)
		if errMarshal != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errMarshal.Error()))
			return
		}
		w.WriteHeader(commentGetStatusCode(err))
		w.Write(responseBytes)
		return
	}
	response := domain.NewMessage("Comment successfully Updated")
	responseBytes, errMarshal := json.Marshal(response)
	if errMarshal != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(errMarshal.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(responseBytes)
}

func (ch *CommentHandler) deleteComment(w http.ResponseWriter, r *http.Request) {
	urlParts := strings.Split(r.URL.String(), "/")
	commentId := urlParts[4]
	tokenString := r.Header.Get("Authorization")
	err := ch.commentUsecase.DeleteComment(commentId, tokenString)
	if err != nil {
		response := domain.NewMessage(err.Error())
		responseBytes, errMarshal := json.Marshal(response)
		if errMarshal != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errMarshal.Error()))
			return
		}
		w.WriteHeader(commentGetStatusCode(err))
		w.Write(responseBytes)
		return
	}

	response := domain.NewMessage("Comment successfully Deleted")
	responseBytes, errMarshal := json.Marshal(response)
	if errMarshal != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(errMarshal.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(responseBytes)
}

func commentGetStatusCode(err error) int {
	switch err {
	case domain.ErrInternalServerError:
		return http.StatusInternalServerError
	case domain.ErrPostNotFound, domain.ErrCommentNotFound:
		return http.StatusNotFound
	case domain.ErrUnauthorizedCommentUpdate, domain.ErrUnauthorizedCommentDelete:
		return http.StatusUnauthorized
	case domain.ErrMissingCommentInput:
		return http.StatusBadRequest
	}
	return http.StatusOK
}
