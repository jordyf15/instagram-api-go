package http

import (
	"encoding/json"
	"instagram-go/domain"
	"io/ioutil"
	"net/http"
	"strings"
)

type PostHandler struct {
	postUsecase domain.PostUsecase
}

func NewPostHandler(postUsecase domain.PostUsecase) domain.PostHandler {
	return &PostHandler{
		postUsecase: postUsecase,
	}
}

func (ph *PostHandler) Posts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		ph.postPost(w, r)
		return
	case "GET":
		ph.getPosts(w, r)
		return
	}
}

func (ph *PostHandler) Post(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "PUT":
		ph.putPost(w, r)
		return
	case "DELETE":
		ph.deletePost(w, r)
		return
	}
}

func (ph *PostHandler) postPost(w http.ResponseWriter, r *http.Request) {
	tokenString := r.Header.Get("Authorization")

	r.ParseMultipartForm(10 << 20)
	formData := r.MultipartForm
	visualMedias := formData.File["visual_medias"]
	if len(visualMedias) == 0 {
		response := domain.NewMessage(domain.ErrMissingVisualMediasInput.Error())
		responseBytes, errMarshal := json.Marshal(response)
		if errMarshal != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errMarshal.Error()))
			return
		}
		w.WriteHeader(postGetStatusCode(domain.ErrMissingVisualMediasInput))
		w.Write(responseBytes)
		return
	}
	for _, v := range visualMedias {
		visualMedia, err := v.Open()
		if err != nil {
			response := domain.NewMessage(domain.ErrInternalServerError.Error())
			responseBytes, errMarshal := json.Marshal(response)
			if errMarshal != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return
			}
			w.WriteHeader(postGetStatusCode(domain.ErrInternalServerError))
			w.Write(responseBytes)
			return
		}
		defer visualMedia.Close()
		fileNameParts := strings.Split(v.Filename, ".")
		extension := fileNameParts[len(fileNameParts)-1]
		if extension != "gif" && extension != "jpg" && extension != "png" &&
			extension != "tiff" && extension != "webm" && extension != "mp4" {
			response := domain.NewMessage(domain.ErrUnsupportedVisualMediaType.Error())
			responseBytes, errMarshal := json.Marshal(response)
			if errMarshal != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return
			}
			w.WriteHeader(postGetStatusCode(domain.ErrUnsupportedVisualMediaType))
			w.Write(responseBytes)
			return
		}
	}

	caption := r.FormValue("caption")
	if caption == "" {
		response := domain.NewMessage(domain.ErrMissingCaptionInput.Error())
		responseBytes, errMarshal := json.Marshal(response)
		if errMarshal != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errMarshal.Error()))
			return
		}
		w.WriteHeader(postGetStatusCode(domain.ErrMissingCaptionInput))
		w.Write(responseBytes)
		return
	}
	var post domain.Post
	post.Caption = caption

	err := ph.postUsecase.InsertPost(&post, tokenString, visualMedias)
	if err != nil {
		response := domain.NewMessage(err.Error())
		responseBytes, errMarshal := json.Marshal(response)
		if errMarshal != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errMarshal.Error()))
			return
		}
		w.WriteHeader(postGetStatusCode(err))
		w.Write(responseBytes)
		return
	}
	response := domain.NewMessage("Post successfully Created")
	responseBytes, errMarshal := json.Marshal(response)
	if errMarshal != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(responseBytes)
}

func (ph *PostHandler) getPosts(w http.ResponseWriter, r *http.Request) {
	posts, err := ph.postUsecase.FindPosts()
	if err != nil {
		response := domain.NewMessage(err.Error())
		responseBytes, errMarshal := json.Marshal(response)
		if errMarshal != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errMarshal.Error()))
			return
		}
		w.WriteHeader(postGetStatusCode(err))
		w.Write(responseBytes)
		return
	}
	dataPosts := domain.NewDataPosts(*posts)
	response := domain.NewDataResponsePosts(*dataPosts)
	responseBytes, errMarshal := json.Marshal(response)
	if errMarshal != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(errMarshal.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(responseBytes)
}

func (ph *PostHandler) putPost(w http.ResponseWriter, r *http.Request) {
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
		w.WriteHeader(postGetStatusCode(domain.ErrInternalServerError))
		w.Write(responseBytes)
		return
	}
	var post domain.Post
	err = json.Unmarshal(bodyBytes, &post)
	if err != nil {
		response := domain.NewMessage(domain.ErrInternalServerError.Error())
		responseBytes, errMarshal := json.Marshal(response)
		if errMarshal != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errMarshal.Error()))
			return
		}
		w.WriteHeader(postGetStatusCode(domain.ErrInternalServerError))
		w.Write(responseBytes)
		return
	}

	if post.Caption == "" {
		response := domain.NewMessage(domain.ErrMissingCaptionInput.Error())
		responseBytes, errMarshal := json.Marshal(response)
		if errMarshal != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errMarshal.Error()))
			return
		}
		w.WriteHeader(postGetStatusCode(domain.ErrMissingCaptionInput))
		w.Write(responseBytes)
		return
	}

	urlParts := strings.Split(r.URL.String(), "/")
	postId := urlParts[2]
	tokenString := r.Header.Get("Authorization")

	err = ph.postUsecase.UpdatePost(postId, post.Caption, tokenString)
	if err != nil {
		response := domain.NewMessage(err.Error())
		responseBytes, errMarshal := json.Marshal(response)
		if errMarshal != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errMarshal.Error()))
			return
		}
		w.WriteHeader(postGetStatusCode(err))
		w.Write(responseBytes)
		return
	}

	response := domain.NewMessage("Post successfully Updated")
	responseBytes, errMarshal := json.Marshal(response)
	if errMarshal != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(errMarshal.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(responseBytes)
}

func (ph *PostHandler) deletePost(w http.ResponseWriter, r *http.Request) {
	urlParts := strings.Split(r.URL.String(), "/")
	postId := urlParts[2]
	tokenString := r.Header.Get("Authorization")

	err := ph.postUsecase.DeletePost(postId, tokenString)
	if err != nil {
		response := domain.NewMessage(err.Error())
		responseBytes, errMarshal := json.Marshal(response)
		if errMarshal != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errMarshal.Error()))
			return
		}
		w.WriteHeader(postGetStatusCode(err))
		w.Write(responseBytes)
		return
	}

	response := domain.NewMessage("Post successfully Deleted")
	responseBytes, errMarshal := json.Marshal(response)
	if errMarshal != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(responseBytes)
}

func postGetStatusCode(err error) int {
	switch err {
	case domain.ErrMissingVisualMediasInput, domain.ErrUnsupportedVisualMediaType, domain.ErrMissingCaptionInput:
		return http.StatusBadRequest
	case domain.ErrInternalServerError:
		return http.StatusInternalServerError
	case domain.ErrPostNotFound:
		return http.StatusNotFound
	case domain.ErrUnauthorizedPostUpdate, domain.ErrUnauthorizedPostDelete:
		return http.StatusUnauthorized
	}
	return http.StatusOK
}
