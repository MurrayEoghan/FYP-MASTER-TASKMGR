package controllers

import (
	"encoding/json"
	"net/http"
	models "repo/models/followerModels"
	repo "repo/repositories"
	"strconv"

	"github.com/gorilla/mux"
)

func Follow(w http.ResponseWriter, r *http.Request) {
	addCorsHeader(w)
	if r.Method == "POST" {
		ids := &models.UserIds{}
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&ids)
		if err != nil {
			createErrorMsg("Incorrect Payload", w)
			w.WriteHeader(http.StatusBadRequest)
		}
		follow := repo.Follow(*ids, w)
		if follow == 0 {
			w.WriteHeader(http.StatusInternalServerError)
			createErrorMsg("Error creating relationship", w)

		} else {
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

}
func UnFollow(w http.ResponseWriter, r *http.Request) {
	addCorsHeader(w)
	if r.Method == "POST" {
		ids := &models.UserIds{}
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&ids)
		if err != nil {
			createErrorMsg("Incorrect Payload", w)
			w.WriteHeader(http.StatusBadRequest)
		}
		follow := repo.UnFollow(*ids, w)
		if follow == 0 {
			createErrorMsg("Error deleting relationship", w)
			w.WriteHeader(http.StatusInternalServerError)

		} else {
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

}

func GetFollowers(w http.ResponseWriter, r *http.Request) {
	addCorsHeader(w)
	if r.Method == "GET" {
		params := mux.Vars(r)
		id, err := strconv.Atoi(params["id"])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			createErrorMsg("Error parsing url value", w)

			return
		}
		res := repo.GetFollowers(id)
		if res == nil {
			w.WriteHeader(http.StatusNotFound)
			createErrorMsg("No Followers", w)
			return
		} else {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(res)

		}
		return
	}
}

func GetFollowing(w http.ResponseWriter, r *http.Request) {
	addCorsHeader(w)
	if r.Method == "GET" {
		params := mux.Vars(r)
		id, err := strconv.Atoi(params["id"])

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			createErrorMsg("Error parsing url value", w)

			return
		}
		res := repo.GetFollowing(id)

		if res == nil {
			w.WriteHeader(http.StatusNotFound)
			createErrorMsg("No Following", w)
			return
		} else {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(res)
		}

	}
}

func GetFollowingPosts(w http.ResponseWriter, r *http.Request) {
	addCorsHeader(w)
	if r.Method == "GET" {
		params := mux.Vars(r)
		id, err := strconv.Atoi(params["id"])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			createErrorMsg("Error parsing url value", w)
			return
		}
		values, err := repo.GetFollowingPosts(id, w)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		} else {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(values)

		}

	}
}
