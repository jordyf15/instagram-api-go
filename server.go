package main

import (
	"context"
	"fmt"
	"instagram-go/authentications"
	"instagram-go/comments"
	"instagram-go/likes"
	"instagram-go/middlewares"
	"instagram-go/posts"
	"instagram-go/users"
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

	userService := users.NewUserService(usersCollection)
	userHandlers := users.NewUserHandlers(*userService)

	authenticationService := authentications.NewAuthenticationService(usersCollection)
	authenticationHandlers := authentications.NewAuthenticationHandler(*authenticationService)

	postService := posts.NewPostService(postsCollection, likesCollection)
	postHandlers := posts.NewPostHandlers(*postService)

	likeService := likes.NewLikeService(likesCollection)
	likeHandlers := likes.NewLikeHandlers(*likeService)

	commentService := comments.NewCommentService(commentsCollection, likesCollection)
	commentHandlers := comments.NewCommentHandlers(*commentService)

	mux := http.NewServeMux()
	mux.HandleFunc("/users", userHandlers.PostUserHandler)
	mux.HandleFunc("/authentications", authenticationHandlers.PostAuthenticationHandler)
	mux.HandleFunc("/users/", userHandlers.PutUserHandler)
	mux.HandleFunc("/posts", postHandlers.Posts)
	mux.HandleFunc("/posts/", func(w http.ResponseWriter, r *http.Request) {
		urlParts := strings.Split(r.URL.String(), "/")
		if len(urlParts) == 3 {
			postHandlers.Post(w, r)
		} else if len(urlParts) == 4 {
			if urlParts[3] == "likes" && r.Method == "POST" {
				likeHandlers.PostPostLikeHandler(w, r)
			} else if urlParts[3] == "comments" {
				commentHandlers.Comments(w, r)
			}
		} else if len(urlParts) == 5 {
			if urlParts[3] == "likes" && r.Method == "DELETE" {
				likeHandlers.DeletePostLikeHandler(w, r)
			} else if urlParts[3] == "comments" {
				commentHandlers.Comment(w, r)
			}
		} else if len(urlParts) == 6 && r.Method == "POST" && urlParts[3] == "comments" {
			likeHandlers.PostCommentLikeHandler(w, r)
		} else if len(urlParts) == 7 && r.Method == "DELETE" {
			likeHandlers.DeleteCommentLikeHandler(w, r)
		}
	})
	wrappedMux := middlewares.NewAuthenticateMiddleware(mux)
	err := http.ListenAndServe(":8000", wrappedMux)
	if err != nil {
		panic(err)
	}
}