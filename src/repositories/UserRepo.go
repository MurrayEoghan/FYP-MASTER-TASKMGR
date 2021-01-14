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
	CreateUser(newUser model.NewUser, w http.ResponseWriter) int64
	UpdateProfile(profile model.NewUserProfile) int64
	UpdateAccount(account model.UpdateUserAccount) int64
}

func GetUserByUsernameAndPassword(loggedInUser model.LogInUser) *model.User {
	user := &model.User{}
	row := sqldb.DB.QueryRow(`SELECT * FROM task_mgr.user WHERE username = ? AND password = ?`, loggedInUser.Username, loggedInUser.Password).Scan(&user.Username, &user.Email, &user.Password, &user.Id, &user.Admin)
	if row != nil {
		return &model.User{}
	}
	return user
}

func UserExists(username string, email string) *model.User {
	user := &model.User{}
	row := sqldb.DB.QueryRow(`SELECT * FROM task_mgr.user WHERE username = ? OR email = ?`, username, email).Scan(&user.Username, &user.Email, &user.Password, &user.Id, &user.Admin)
	if row != nil {
		return &model.User{}
	}
	return user
}

func CreateUser(newUser model.NewUser, w http.ResponseWriter) int64 {
	var amount int

	count := sqldb.DB.QueryRow("SELECT COUNT(Id) FROM task_mgr.user").Scan(&amount)
	if count != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return 0
	}

	stmt, err := sqldb.DB.Prepare(`INSERT INTO task_mgr.user (username, email, password, Id, admin) VALUES (?,?,?,?,?)`)
	userRow, err := stmt.Exec(newUser.Username, newUser.Email, newUser.Password, amount+1, 0)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Error Inserting Record")
		return 0
	}
	stmt2, err := sqldb.DB.Prepare(`INSERT INTO task_mgr.user_profile (Id, Fname, Lname, Age, Sex, Add1, Add2, Add3, County, Country) VALUES (?,?,?,?,?,?,?,?,?,?)`)
	userProfileRow, err := stmt2.Exec(amount+1, "", "", 0, "", "", "", "", "", "")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Error Inserting Profile Record")
	}
	lastId, err := userRow.LastInsertId()
	lastProfileId, err := userProfileRow.LastInsertId()
	fmt.Printf("Last Insert Id User Table: %d\nLast Insert Id Profile Table: %d\n", lastId, lastProfileId)
	return lastId

}

func UpdateProfile(profile model.NewUserProfile) int64 {

	stmt, err := sqldb.DB.Prepare(`UPDATE task_mgr.user_profile SET Fname = ?, Lname = ?, Age = ?, Sex = ?, Add1 = ?, Add2 = ?, Add3 = ?, County = ?, Country = ? WHERE Id = ?`)
	if err != nil {
		log.Fatal(err)
		return 0
	}
	stmt.Exec(profile.Fname, profile.Lname, profile.Age, profile.Gender, profile.Address1, profile.Address2, profile.Address3, profile.County, profile.Country, profile.Id)
	return 1
}

func UpdateAccount(account model.UpdateUserAccount) int64 {
	stmt, err := sqldb.DB.Prepare(`UPDATE task_mgr.user SET username = ?, email = ?, password = ? WHERE Id = ?`)
	if err != nil {
		log.Fatal(err)
		return 0
	}
	stmt.Exec(account.Username, account.Email, account.Password, account.Id)
	return 1
}
