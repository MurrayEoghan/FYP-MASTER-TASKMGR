package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	models "repo/models/forumModels"
	repo "repo/repositories"
	"strconv"

	"github.com/gorilla/mux"
)

func GetAllPosts(w http.ResponseWriter, r *http.Request) {
	addCorsHeader(w)

	if r.Method == "GET" {
		posts := repo.GetAllPosts()
		if posts == nil {
			w.WriteHeader(http.StatusInternalServerError)
			createErrorMsg("No Posts Found", w)
			return
		} else {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(posts)
		}
		return
	}
}

func CreatePost(w http.ResponseWriter, r *http.Request) {
	addCorsHeader(w)
	if r.Method == "POST" {
		newPost := &models.CreatePostModel{}
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&newPost)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			createErrorMsg("Incorrect payload", w)
			log.Printf(err.Error())
			return
		}
		success := repo.CreatePost(*newPost, w)
		if success == 0 {
			w.WriteHeader(http.StatusConflict)
			createErrorMsg("An Unexpected error has occured", w)
			return
		} else {
			w.WriteHeader(http.StatusCreated)

		}
		return

	}
}

func CreateComment(w http.ResponseWriter, r *http.Request) {
	addCorsHeader(w)
	if r.Method == "POST" {
		newComment := &models.CreateComment{}
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&newComment)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			createErrorMsg("Incorrect Payload", w)
			log.Print(err)
			return
		}
		success := repo.CreateComment(*newComment, w)
		if success == 0 {
			w.WriteHeader(http.StatusInternalServerError)
			createErrorMsg("An Unexpected error has occured", w)
			return
		} else {
			w.WriteHeader(http.StatusCreated)
		}
		return
	}
}

func GetTopics(w http.ResponseWriter, r *http.Request) {
	addCorsHeader(w)
	if r.Method == "GET" {
		topics := repo.GetTopics()
		if topics == nil {
			w.WriteHeader(http.StatusInternalServerError)
			createErrorMsg("No Topics Found. Check Server", w)
			return
		} else {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(topics)
		}
		return
	}
}

func GetPost(w http.ResponseWriter, r *http.Request) {
	addCorsHeader(w)
	if r.Method == "POST" {
		postId := &models.PostId{}
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&postId)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			createErrorMsg("Incorrect Body", w)
			return
		}

		post := repo.GetPost(postId.PostId, w)
		if post == nil {
			w.WriteHeader(http.StatusNotFound)
			createErrorMsg("No Post with that ID found.", w)
			return
		} else {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(post)
		}
		return

	}
}

func GetRecentPosts(w http.ResponseWriter, r *http.Request) {
	addCorsHeader(w)
	if r.Method == "POST" {
		authorId := &models.AuthorId{}
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&authorId)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			createErrorMsg("Incorrect Payload", w)
			return
		}
		posts := repo.GetRecentPosts(authorId.AuthorId, w)
		if posts == nil {
			w.WriteHeader(http.StatusNotFound)
			createErrorMsg("Check Repo. No Posts Found", w)
			return
		} else {

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(posts)
		}
		return

	}
}

func UpdatePost(w http.ResponseWriter, r *http.Request) {
	addCorsHeader(w)
	if r.Method == "PATCH" {
		newCommentBody := &models.UpdatePost{}
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&newCommentBody)
		if err != nil {
			createErrorMsg("Bad Payload", w)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		value := repo.UpdatePost(*newCommentBody, w)
		if value == 0 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		} else {
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
}

func UpdateComment(w http.ResponseWriter, r *http.Request) {
	addCorsHeader(w)
	if r.Method == "PATCH" {
		newCommentBody := &models.UpdateComment{}
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&newCommentBody)
		if err != nil {
			createErrorMsg("Bad Payload", w)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		value := repo.UpdateComment(*newCommentBody, w)
		if value == 0 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		} else {
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
}

func DeleteComment(w http.ResponseWriter, r *http.Request) {
	addCorsHeader(w)
	if r.Method == "DELETE" {

		params := mux.Vars(r)
		id, err := strconv.Atoi(params["id"])
		if err != nil {
			createErrorMsg("Error parsing url value", w)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		value, err := repo.DeleteComment(id, w)
		if value < 0 {
			createErrorMsg("Unknown Error Occurred", w)
			w.WriteHeader(http.StatusInternalServerError)
			return
		} else {
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
}

func DeletePost(w http.ResponseWriter, r *http.Request) {
	addCorsHeader(w)
	if r.Method == "DELETE" {

		params := mux.Vars(r)
		id, err := strconv.Atoi(params["id"])
		if err != nil {
			createErrorMsg("Error parsing url value", w)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		post, err := repo.DeletePost(id, w)
		comments, err := repo.DeleteComments(id, w)
		answer, err := repo.DeleteAnswer(id, w)
		if post < 0 && comments < 0 && answer < 0 {
			createErrorMsg("Unknown Error Occurred", w)
			w.WriteHeader(http.StatusInternalServerError)
			return
		} else {
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
}

func UpdateAnswer(w http.ResponseWriter, r *http.Request) {
	addCorsHeader(w)
	if r.Method == "PATCH" {
		newAnswerBody := &models.UpdateAnswer{}
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&newAnswerBody)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			createErrorMsg("Bad Payload", w)

			return
		}
		value := repo.UpdateAnswer(*newAnswerBody, w)
		if value == 0 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		} else {
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
}

func DeleteAnswer(w http.ResponseWriter, r *http.Request) {
	addCorsHeader(w)
	if r.Method == "DELETE" {

		params := mux.Vars(r)
		id, err := strconv.Atoi(params["id"])
		if err != nil {
			createErrorMsg("Error parsing url value", w)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		answer, err := repo.DeleteAnswer(id, w)

		if answer < 0 {
			createErrorMsg("Unknown Error Occurred", w)
			w.WriteHeader(http.StatusInternalServerError)
			return
		} else {
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
}

func CreateAnswer(w http.ResponseWriter, r *http.Request) {
	addCorsHeader(w)
	if r.Method == "POST" {
		newAnswer := &models.CreateAnswer{}
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&newAnswer)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			createErrorMsg("Incorrect Payload", w)
			log.Print(err)
			return
		}
		success := repo.CreateAnswer(*newAnswer, w)
		if success == 0 {
			w.WriteHeader(http.StatusInternalServerError)
			createErrorMsg("An Unexpected error has occured", w)
			return
		} else {
			w.WriteHeader(http.StatusCreated)
		}
		return
	}
}
