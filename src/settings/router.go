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
	log.Fatal(http.ListenAndServe(":8080", router))
}
