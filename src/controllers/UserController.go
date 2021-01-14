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
	UpdateUserProfile(w http.ResponseWriter, r *http.Request)
	UpdateUserAccount(w http.ResponseWriter, r *http.Request)
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
			log.Printf(err.Error())
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
	addCorsHeader(w)
	if r.Method == "POST" {

		newUser := &model.NewUser{}
		userId := &model.UserId{}
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&newUser)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			createErrorMsg("Internal Server Error", w)

			log.Printf(err.Error())
			return
		}

		if newUser.Email == "" || newUser.Username == "" || newUser.Password == "" {
			w.WriteHeader(http.StatusBadRequest)
			createErrorMsg("Missing Field", w)
			log.Printf("User left something blank")
			return
		}
		exists := *rep.UserExists(newUser.Username, newUser.Email)
		if exists.Username == "" && exists.Email == "" {
			userId.UserId = rep.CreateUser(*newUser, w)
			json.NewEncoder(w).Encode(userId)
			w.WriteHeader(http.StatusCreated)
			return

		} else {

			w.WriteHeader(http.StatusConflict)
			createErrorMsg("User Already Exists", w)
			log.Printf("User Already Exists")
		}

		return
	}
}

func UpdateUserProfile(w http.ResponseWriter, r *http.Request) {
	addCorsHeader(w)
	if r.Method == "POST" {

		profileModel := &model.NewUserProfile{}
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&profileModel)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			createErrorMsg("Missing Field", w)
			log.Printf("Incorrect param or missing field in request body")
			return
		}
		profileUpdate := rep.UpdateProfile(*profileModel)
		if profileUpdate <= 0 {
			w.WriteHeader(http.StatusInternalServerError)
			createErrorMsg("Internal Server Error", w)
			log.Printf("Error inserting record")
			return
		} else {
			w.WriteHeader(http.StatusNoContent)
			return
		}

	}
}

func UpdateUserAccount(w http.ResponseWriter, r *http.Request) {
	addCorsHeader(w)
	if r.Method == "POST" {
		accountModel := &model.UpdateUserAccount{}
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&accountModel)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Printf(err.Error())
			createErrorMsg("Incorrect Request Body", w)

			return
		}
		if accountModel.Email == "" || accountModel.Password == "" || accountModel.Username == "" {
			w.WriteHeader(http.StatusBadRequest)
			log.Printf("Missing Field")
			createErrorMsg("Missing or Null Field", w)

			return
		}
		update := rep.UpdateAccount(*accountModel)
		if update <= 0 {
			w.WriteHeader(http.StatusInternalServerError)
			createErrorMsg("Internal Server Error", w)
			log.Printf("Error inserting record")
			return
		} else {
			w.WriteHeader(http.StatusNoContent)
			return
		}

	}
}
