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
	userRow, err := stmt.Exec(newUser.Username, newUser.Email, newUser.Password, amount+1, 0)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Error Inserting Record")
		return
	}
	stmt2, err := sqldb.DB.Prepare(`INSERT INTO task_mgr.user_profile (Id, Fname, Lname, Age, Sex, Add1, Add2, Add3, County, Country) VALUES (?,?,?,?,?,?,?,?,?,?)`)
	userProfileRow, err := stmt2.Exec(amount+1, newUser.Fname, newUser.Lname, newUser.Age, newUser.Sex, newUser.Address1, newUser.Address2, newUser.Address3, newUser.County, newUser.Country)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Error Inserting Profile Record")
	}
	lastId, err := userRow.LastInsertId()
	lastProfileId, err := userProfileRow.LastInsertId()
	fmt.Printf("Last Insert Id User Table: %d\nLast Insert Id Profile Table: %d\n", lastId, lastProfileId)
	return

}
