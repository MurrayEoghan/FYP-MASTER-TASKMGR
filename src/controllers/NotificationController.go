package controllers

import (
	"encoding/json"
	"net/http"

	// models "repo/models/notificationModels"
	repo "repo/repositories"
	"strconv"

	"github.com/gorilla/mux"
)

func GetNotifications(w http.ResponseWriter, r *http.Request) {
	addCorsHeader(w)
	if r.Method == "GET" {
		params := mux.Vars(r)
		id, err := strconv.Atoi(params["id"])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			createErrorMsg("Error parsing url value", w)

			return
		}
		res := repo.GetNotifications(id, w)
		if res == nil {
			w.WriteHeader(http.StatusNoContent)
			createErrorMsg("No Followers", w)
			return
		} else {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(res)

		}
		return
	}
}

func DeleteNotifications(w http.ResponseWriter, r *http.Request) {
	addCorsHeader(w)
	if r.Method == "DELETE" {

		params := mux.Vars(r)
		id, err := strconv.Atoi(params["id"])
		if err != nil {
			createErrorMsg("Error parsing url value", w)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		value, err := repo.DeleteNotifications(id, w)
		if value < 0 {
			createErrorMsg("No notifications exist", w)
			w.WriteHeader(http.StatusNoContent)
			return
		} else {
			w.WriteHeader(http.StatusAccepted)
			return
		}
	}
}
