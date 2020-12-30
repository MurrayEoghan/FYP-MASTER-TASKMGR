package settings

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	controller "repo/controllers"
)

func Router() {

	router := mux.NewRouter().StrictSlash(true)
	v1Router := router.PathPrefix("/api/v1").Subrouter()
	v1Router.HandleFunc("/login", controller.UserLogIn).Methods(http.MethodGet)
	v1Router.HandleFunc("/signup", controller.UserSignUp).Methods(http.MethodPost)
	log.Fatal(http.ListenAndServe(":8080", router))
}
