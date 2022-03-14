package handlers

import (
	"encoding/json"
	"fmt"
	"instagram-go/models"
	"instagram-go/services"
	"io"
	"io/fs"
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
	service           services.IPostService
	postHandlerHeader IPostHandlerHeader
	postFileOsHandler IPostFileOsHandler
}

func NewPostHandlers(service services.IPostService, postHandlerHeader IPostHandlerHeader, postFileOsHandler IPostFileOsHandler) *PostHandlers {
	if postHandlerHeader == nil {
		postHandlerHeader = newPostHandlerHeader()
	}
	if postFileOsHandler == nil {
		postFileOsHandler = newPostFileOsHandler()
	}
	return &PostHandlers{
		service:           service,
		postHandlerHeader: postHandlerHeader,
		postFileOsHandler: postFileOsHandler,
	}
}

type IPostHandlerHeader interface {
	getUserIdFromToken(string) (string, error)
}
type postHandlerHeader struct {
}

func newPostHandlerHeader() *postHandlerHeader {
	return &postHandlerHeader{}
}

func (phh *postHandlerHeader) getUserIdFromToken(tokenString string) (string, error) {
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

type postFileOsHandler struct {
}

type IPostFileOsHandler interface {
	create(string) (*os.File, error)
	copy(io.Writer, io.Reader) (int64, error)
	mkDirAll(string, fs.FileMode) error
}

func newPostFileOsHandler() *postFileOsHandler {
	return &postFileOsHandler{}
}

func (pfoh *postFileOsHandler) create(name string) (*os.File, error) {
	return os.Create(name)
}

func (pfoh *postFileOsHandler) copy(dst io.Writer, src io.Reader) (int64, error) {
	return io.Copy(dst, src)
}

func (pfoh *postFileOsHandler) mkDirAll(path string, perm fs.FileMode) error {
	return os.MkdirAll(path, perm)
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

	userId, err := ph.postHandlerHeader.getUserIdFromToken(tokenString)
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

	newPostId := "post-" + uuid.NewString()
	r.ParseMultipartForm(10 << 20)
	formData := r.MultipartForm
	visualMedias := formData.File["visual_medias"]
	var visualMediaUrls []string
	newpath := filepath.Join(".", "visual_medias")

	err = ph.postFileOsHandler.mkDirAll(newpath, os.ModePerm)
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
	if len(visualMedias) == 0 {
		response := models.NewMessage("Visual Medias must not be empty")
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
		if extension != "gif" && extension != "jpg" && extension != "png" &&
			extension != "tiff" && extension != "webm" && extension != "mp4" {
			response := models.NewMessage("Uploaded Visual Medias type is not supported")
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
		visualMediaUrl := "./visual_medias/" + newPostId + "-" + strconv.Itoa(k) + "." + extension
		out, err := ph.postFileOsHandler.create(visualMediaUrl)
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
		defer out.Close()

		_, err = ph.postFileOsHandler.copy(out, visualMedia)
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
		visualMediaUrls = append(visualMediaUrls, visualMediaUrl)
	}

	caption := r.FormValue("caption")
	if caption == "" {
		response := models.NewMessage("Caption must not be empty")
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
	createdTime := time.Now()
	newPost := *models.NewPost(newPostId, userId, visualMediaUrls, caption, 0, createdTime, createdTime)

	ph.Lock()
	err = ph.service.InsertPost(newPost)
	defer ph.Unlock()

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
	dataPosts := *models.NewDataPosts(posts)
	response := *models.NewDataResponsePosts(dataPosts)
	responseBytes, err := json.Marshal(response)
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
	w.Write(responseBytes)
}

func (ph *PostHandlers) putPostHandler(w http.ResponseWriter, r *http.Request) {
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
		response := models.NewMessage(fmt.Sprintf("need content-type 'application/json', but got '%s'", ct))
		responseBytes, err := json.Marshal(response)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write(responseBytes)
		return
	}

	var post models.Post
	err = json.Unmarshal(bodyBytes, &post)
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
	postId := urlParts[2]
	tokenString := r.Header.Get("Authorization")

	userId, err := ph.postHandlerHeader.getUserIdFromToken(tokenString)
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

	ph.Lock()
	isPostExist, err := ph.service.CheckIfPostExist(postId)
	ph.Unlock()
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
		response := *models.NewMessage("Post does not exist")
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

	ph.Lock()
	postUserId, err := ph.service.GetPostUserId(postId)
	ph.Unlock()

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

	if postUserId != userId {
		response := *models.NewMessage("User is not authorized to update this post")
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
	if post.Caption == "" {
		response := *models.NewMessage("Caption must not be empty")
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
	ph.Lock()
	err = ph.service.UpdatePost(postId, post.Caption)
	defer ph.Unlock()
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

	userId, err := ph.postHandlerHeader.getUserIdFromToken(tokenString)
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
	ph.Lock()
	isPostExist, err := ph.service.CheckIfPostExist(postId)
	ph.Unlock()
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
		response := *models.NewMessage("Post does not exist")
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
	ph.Lock()
	postUserId, err := ph.service.GetPostUserId(postId)
	ph.Unlock()

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

	if postUserId != userId {
		response := *models.NewMessage("User is not authorized to delete this post")
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
