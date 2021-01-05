package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	model "repo/models"
	rep "repo/repositories"
)

type UserController interface {
	UserLogIn(w http.ResponseWriter, r *http.Request)
	UserSignUp(w http.ResponseWriter, r *http.Request)
	createErrorMsg(message string, w http.ResponseWriter)
}

func createErrorMsg(message string, w http.ResponseWriter) {
	errorMsg := &model.ErrorMessage{}
	errorMsg.Message = message
	json.NewEncoder(w).Encode(&errorMsg)

}

func addCorsHeader(res http.ResponseWriter) {
	headers := res.Header()
	headers.Add("Access-Control-Allow-Origin", "*")
	headers.Add("Vary", "Accept-Encoding, Origin")
	headers.Add("Vary", "Access-Control-Request-Method")
	headers.Add("Vary", "Access-Control-Request-Headers")
	headers.Add("Access-Control-Allow-Headers", "Content-Type, Origin, Accept, token, Access-Control-Allow-Headers")
	headers.Add("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
}

func UserLogIn(w http.ResponseWriter, r *http.Request) {

	addCorsHeader(w)
	if r.Method == "POST" {
		var u model.User
		loggedInUser := &model.LogInUser{}

		decoder := json.NewDecoder(r.Body)

		err := decoder.Decode(&loggedInUser)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Fatal(err)
		}
		if loggedInUser.Username == "" || loggedInUser.Password == "" {
			w.WriteHeader(http.StatusBadRequest)
		}
		u = *rep.GetUserByUsernameAndPassword(*loggedInUser)
		if u.Username == "" {
			w.WriteHeader(http.StatusNotFound)
		} else {

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(u)
		}
	}

}

func UserSignUp(w http.ResponseWriter, r *http.Request) {
	newUser := &model.NewUser{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&newUser)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		createErrorMsg("Internal Server Error", w)

		log.Fatal(err)
		return
	}
	if newUser.Email == "" || newUser.Username == "" || newUser.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		createErrorMsg("Missing Field", w)
		log.Printf("User left something blank")
		return
	}
	exists := rep.UserExists(newUser.Username, newUser.Email)
	if exists {
		w.WriteHeader(http.StatusConflict)
		createErrorMsg("User Already Exists", w)
		log.Printf("User Already Exists")
		return
	} else {
		rep.CreateUser(*newUser, w)
	}

	return

}
