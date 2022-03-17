package main

import (
	"context"
	"fmt"
	commentHttp "instagram-go/comment/delivery/http"
	commentRepo "instagram-go/comment/repository/mongodb"
	commentUsecase "instagram-go/comment/usecase"
	"instagram-go/domain"
	likeHttp "instagram-go/like/delivery/http"
	likeRepo "instagram-go/like/repository/mongodb"
	likeUsecase "instagram-go/like/usecase"
	"instagram-go/middlewares"
	postHttp "instagram-go/post/delivery/http"
	postRepo "instagram-go/post/repository/mongodb"
	postUsecase "instagram-go/post/usecase"
	userHttp "instagram-go/user/delivery/http"
	userRepo "instagram-go/user/repository/mongodb"
	userUsecase "instagram-go/user/usecase"
	"net/http"
	"strings"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	client, dbErr := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if dbErr != nil {
		panic(dbErr)
	}
	if client != nil {
		fmt.Println("database connected")
	}
	usersCollection := client.Database("instagram").Collection("users")
	postsCollection := client.Database("instagram").Collection("posts")
	likesCollection := client.Database("instagram").Collection("likes")
	commentsCollection := client.Database("instagram").Collection("comments")

	userRepository := userRepo.NewMongodbUserRepository(usersCollection)
	postRepository := postRepo.NewMongodbPostRepository(postsCollection)
	likeRepository := likeRepo.NewMongodbLikeRepository(likesCollection)
	commentRepository := commentRepo.NewMongodbCommentRepository(commentsCollection)

	authenticationHelper := domain.NewAuthenticationHelper()
	fileOsHelper := domain.NewFileOsHelper()
	headerHelper := domain.NewHeaderHelper()

	userUseCase := userUsecase.NewUserUsecase(userRepository, authenticationHelper, headerHelper, fileOsHelper)
	postUsecase := postUsecase.NewPostUseCase(postRepository, likeRepository, headerHelper, fileOsHelper)
	likeUsecase := likeUsecase.NewLikeUsecase(likeRepository, postRepository, commentRepository, headerHelper)
	commentUsecase := commentUsecase.NewCommentUsecase(commentRepository, postRepository, likeRepository, headerHelper)

	postHandler := postHttp.NewPostHandler(postUsecase)
	userHandler := userHttp.NewUserHandler(userUseCase)
	likeHandler := likeHttp.NewLikeHandler(likeUsecase)
	commentHandler := commentHttp.NewCommentHandler(commentUsecase)

	mux := http.NewServeMux()
	mux.HandleFunc("/users", userHandler.PostUser)
	mux.HandleFunc("/authentications", userHandler.AuthenticateUser)
	mux.HandleFunc("/users/", userHandler.PutUser)
	mux.HandleFunc("/posts", postHandler.Posts)
	mux.HandleFunc("/posts/", func(w http.ResponseWriter, r *http.Request) {
		urlParts := strings.Split(r.URL.String(), "/")
		if len(urlParts) == 3 {
			postHandler.Post(w, r)
		} else if len(urlParts) == 4 {
			if urlParts[3] == "likes" && r.Method == "POST" {
				likeHandler.PostLikePost(w, r)
			} else if urlParts[3] == "comments" {
				commentHandler.Comments(w, r)
			}
		} else if len(urlParts) == 5 {
			if urlParts[3] == "likes" && r.Method == "DELETE" {
				likeHandler.DeleteLikePost(w, r)
			} else if urlParts[3] == "comments" {
				commentHandler.Comment(w, r)
			}
		} else if len(urlParts) == 6 && r.Method == "POST" && urlParts[3] == "comments" {
			likeHandler.PostCommentLike(w, r)
		} else if len(urlParts) == 7 && r.Method == "DELETE" {
			likeHandler.DeleteCommentLike(w, r)
		}
	})

	wrappedMux := middlewares.NewAuthenticateMiddleware(mux)
	err := http.ListenAndServe(":8000", wrappedMux)
	if err != nil {
		panic(err)
	}
}
