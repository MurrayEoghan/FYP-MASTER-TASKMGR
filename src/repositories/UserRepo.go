package repositories

import (
	"fmt"
	"log"
	"net/http"
	model "repo/models"
	"repo/sqldb"
)

var userLogIn model.LogInUser

type UserRepo interface {
	GetUserByUsernameAndPassword(loggedInUser model.LogInUser) *model.User
	UserExists(username string, email string) bool
	CreateUser(newUser model.NewUser, w http.ResponseWriter)
}

func GetUserByUsernameAndPassword(loggedInUser model.LogInUser) *model.User {
	user := &model.User{}
	row := sqldb.DB.QueryRow(`SELECT * FROM task_mgr.user WHERE username = ? AND password = ?`, loggedInUser.Username, loggedInUser.Password).Scan(&user.Username, &user.Email, &user.Password, &user.Id, &user.Admin)
	if row != nil {
		return &model.User{}
	}
	return user
}

func UserExists(username string, email string) bool {
	user := &model.User{}
	row := sqldb.DB.QueryRow(`SELECT * FROM task_mgr.user WHERE username = ? OR email = ?`, username, email).Scan(&user.Username, &user.Email, &user.Password, &user.Id, &user.Admin)
	if row != nil {
		return false
	}
	return true
}

func CreateUser(newUser model.NewUser, w http.ResponseWriter) {
	var amount int

	count := sqldb.DB.QueryRow("SELECT COUNT(Id) FROM task_mgr.user").Scan(&amount)
	if count != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	stmt, err := sqldb.DB.Prepare(`INSERT INTO task_mgr.user (username, email, password, Id, admin) VALUES (?,?,?,?,?)`)
	row, err := stmt.Exec(newUser.Username, newUser.Email, newUser.Password, amount+1, 0)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Error Inserting Record")
		return
	}
	lastId, err := row.LastInsertId()
	fmt.Printf("Last Insert Id : %d\n", lastId)
	return

}
