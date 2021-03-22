package settings

import (
	"log"
	"net/http"

	// "github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	controller "repo/controllers"
)

func Router() {

	router := mux.NewRouter().StrictSlash(true)
	v1Router := router.PathPrefix("/api/v1").Subrouter()
	v1Router.HandleFunc("/login", controller.UserLogIn).Methods(http.MethodPost, http.MethodOptions)
	v1Router.HandleFunc("/signup", controller.UserSignUp).Methods(http.MethodPost, http.MethodOptions)
	v1Router.HandleFunc("/update/profile", controller.UpdateUserProfile).Methods(http.MethodPost, http.MethodOptions)
	v1Router.HandleFunc("/update/account", controller.UpdateUserAccount).Methods(http.MethodPost, http.MethodOptions)
	v1Router.HandleFunc("/user", controller.GetUserById).Methods(http.MethodGet)
	v1Router.HandleFunc("/posts", controller.GetAllPosts).Methods(http.MethodGet)
	v1Router.HandleFunc("/post", controller.GetPost).Methods(http.MethodPost, http.MethodOptions)
	v1Router.HandleFunc("/posts/create", controller.CreatePost).Methods(http.MethodPost, http.MethodOptions)
	v1Router.HandleFunc("/comment", controller.CreateComment).Methods(http.MethodPost, http.MethodOptions)
	v1Router.HandleFunc("/topics", controller.GetTopics).Methods(http.MethodGet)
	v1Router.HandleFunc("/update/post", controller.UpdatePost).Methods(http.MethodPatch, http.MethodOptions)
	v1Router.HandleFunc("/update/comment", controller.UpdateComment).Methods(http.MethodPatch, http.MethodOptions)
	v1Router.HandleFunc("/delete/comment/{id}", controller.DeleteComment).Methods(http.MethodDelete, http.MethodOptions)
	v1Router.HandleFunc("/delete/post/{id}", controller.DeletePost).Methods(http.MethodDelete, http.MethodOptions)
	v1Router.HandleFunc("/update/post/answer", controller.UpdateAnswer).Methods(http.MethodOptions, http.MethodPatch)
	v1Router.HandleFunc("/delete/post/answer/{id}", controller.DeleteAnswer).Methods(http.MethodOptions, http.MethodDelete)
	v1Router.HandleFunc("/post/answer/create", controller.CreateAnswer).Methods(http.MethodOptions, http.MethodPost)
	v1Router.HandleFunc("/posts/recent", controller.GetRecentPosts).Methods(http.MethodOptions, http.MethodPost)
	v1Router.HandleFunc("/user/follow", controller.Follow).Methods(http.MethodOptions, http.MethodPost)
	v1Router.HandleFunc("/user/{id}/followers", controller.GetFollowers).Methods(http.MethodGet)
	v1Router.HandleFunc("/user/{id}/following", controller.GetFollowing).Methods(http.MethodGet)
	v1Router.HandleFunc("/user/unfollow", controller.UnFollow).Methods(http.MethodOptions, http.MethodPost)

	log.Fatal(http.ListenAndServe(":8080", router))
}
