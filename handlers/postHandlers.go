package handlers

import (
	"encoding/json"
	"fmt"
	"instagram-go/models"
	"instagram-go/services"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type PostHandlers struct {
	sync.Mutex
	service *services.PostService
}

func NewPostHandlers(service *services.PostService) *PostHandlers {
	return &PostHandlers{
		service: service,
	}
}

func (ph *PostHandlers) Posts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		ph.postPostHandler(w, r)
		return
	case "GET":
		ph.getPostsHandler(w, r)
		return
	}
}

func (ph *PostHandlers) Post(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "PUT":
		ph.putPostHandler(w, r)
		return
	case "DELETE":
		ph.deletePostHandler(w, r)
		return
	}
}

func (ph *PostHandlers) postPostHandler(w http.ResponseWriter, r *http.Request) {
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
	newPostId := "post-" + uuid.NewString()
	r.ParseMultipartForm(10 << 20)
	formData := r.MultipartForm
	visualMedias := formData.File["visual_medias"]
	var visualMediaUrls []string
	newpath := filepath.Join(".", "visual_medias")
	err = os.MkdirAll(newpath, os.ModePerm)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	for k := range visualMedias {
		visualMedia, err := visualMedias[k].Open()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		defer visualMedia.Close()
		fileNameParts := strings.Split(visualMedias[k].Filename, ".")
		extension := fileNameParts[len(fileNameParts)-1]
		visualMediaUrl := "./visual_medias/" + newPostId + "-" + strconv.Itoa(k) + "." + extension
		out, err := os.Create(visualMediaUrl)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		defer out.Close()

		_, err = io.Copy(out, visualMedia)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		visualMediaUrls = append(visualMediaUrls, visualMediaUrl)
	}
	caption := r.FormValue("caption")
	createdTime := time.Now()
	newPost := *models.NewPost(newPostId, userId, visualMediaUrls, caption, 0, createdTime, createdTime)

	ph.Lock()
	err = ph.service.InsertPost(newPost)
	defer ph.Unlock()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	} else {
		response := *models.NewMessage("Post successfully Created")
		responseBytes, err := json.Marshal(response)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusCreated)
		w.Write(responseBytes)
	}
}

func (ph *PostHandlers) getPostsHandler(w http.ResponseWriter, r *http.Request) {
	ph.Lock()
	posts, err := ph.service.FindAllPost()
	defer ph.Unlock()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	dataPosts := *models.NewDataPosts(posts)
	response := *models.NewDataResponsePosts(dataPosts)
	responseBytes, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(responseBytes)
}

func (ph *PostHandlers) putPostHandler(w http.ResponseWriter, r *http.Request) {
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

	var post models.Post
	err = json.Unmarshal(bodyBytes, &post)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

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
	ph.Lock()
	postUserId, err := ph.service.GetPostUserId(postId)
	ph.Unlock()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	if postUserId != userId {
		response := *models.NewMessage("You are not authorized to update this post")
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
	ph.Lock()
	err = ph.service.UpdatePost(postId, post.Caption)
	defer ph.Unlock()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	response := *models.NewMessage("Post successfully Updated")
	responseBytes, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(responseBytes)

}

func (ph *PostHandlers) deletePostHandler(w http.ResponseWriter, r *http.Request) {
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

	ph.Lock()
	postUserId, err := ph.service.GetPostUserId(postId)
	ph.Unlock()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	if postUserId != userId {
		response := *models.NewMessage("You are not authorized to delete this post")
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

	ph.Lock()
	err = ph.service.DeletePost(postId)
	defer ph.Unlock()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	response := models.NewMessage("Post successfully Deleted")
	responseBytes, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(responseBytes)
}
