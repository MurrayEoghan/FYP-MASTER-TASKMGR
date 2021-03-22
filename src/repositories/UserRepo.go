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
	UserExistsWithDifferentId(account model.UpdateUserAccount) *model.User
	GetUserById(id int) *model.WholeUser
}

func GetUserByUsernameAndPassword(loggedInUser model.LogInUser) *model.User {
	user := &model.User{}
	row := sqldb.DB.QueryRow(`SELECT * FROM task_mgr.user WHERE username = ? AND password = ?`, loggedInUser.Username, loggedInUser.Password).Scan(&user.Username, &user.Email, &user.Password, &user.Id, &user.ProfessionId)
	if row != nil {
		return &model.User{}
	}
	return user
}

func UserExists(username string, email string) *model.User {
	user := &model.User{}
	row := sqldb.DB.QueryRow(`SELECT * FROM task_mgr.user WHERE username = ? OR email = ?`, username, email).Scan(&user.Username, &user.Email, &user.Password, &user.Id, &user.ProfessionId)
	if row != nil {
		return &model.User{}
	}
	return user
}

func CreateUser(newUser model.NewUser, w http.ResponseWriter) int64 {

	stmt, err := sqldb.DB.Prepare(`INSERT INTO task_mgr.user (username, email, password,  admin) VALUES (?,?,?,?)`)
	userRow, err := stmt.Exec(newUser.Username, newUser.Email, newUser.Password, 0)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Error Inserting Record")
		return 0
	}
	stmt2, err := sqldb.DB.Prepare(`INSERT INTO task_mgr.user_profile (Fname, Lname, Age, Sex, Add1, Add2, Add3, County, Country) VALUES (?,?,?,?,?,?,?,?,?)`)
	userProfileRow, err := stmt2.Exec("", "", 0, "", "", "", "", "", "")
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

func UserExistsWithDifferentId(account model.UpdateUserAccount) *model.User {
	user := &model.User{}
	stmt, err := sqldb.DB.Query("SELECT * FROM task_mgr.user WHERE Id != ? AND (email = ? OR username = ?)", account.Id, account.Email, account.Username)
	if err != nil {
		log.Printf(err.Error())
		return &model.User{}
	}
	defer stmt.Close()
	for stmt.Next() {
		err := stmt.Scan(&user.Username, &user.Email, &user.Password, &user.Id, &user.ProfessionId)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("%s", user.Username)
	}
	return user
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

func GetUserById(id int) *model.WholeUser {
	wholeUser := &model.WholeUser{}
	row := sqldb.DB.QueryRow(`SELECT task_mgr.user.username, task_mgr.user.email, task_mgr.user.profession_id, task_mgr.user.Id, task_mgr.user_profile.Fname, task_mgr.user_profile.Lname, task_mgr.user_profile.Age, task_mgr.user_profile.Sex, task_mgr.user_profile.Add1, task_mgr.user_profile.Add2, task_mgr.user_profile.Add3, task_mgr.user_profile.County, task_mgr.user_profile.Country FROM task_mgr.user, task_mgr.user_profile WHERE task_mgr.user.Id = ? AND task_mgr.user_profile.Id = ?;`, id, id).Scan(&wholeUser.Username, &wholeUser.Email, &wholeUser.ProfessionId, &wholeUser.Id, &wholeUser.Fname, &wholeUser.Lname, &wholeUser.Age, &wholeUser.Gender, &wholeUser.Address1, &wholeUser.Address2, &wholeUser.Address3, &wholeUser.County, &wholeUser.Country)
	if row != nil {
		return &model.WholeUser{}
	}
	return wholeUser
}
