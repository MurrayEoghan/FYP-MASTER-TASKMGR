package models

type User struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Admin    bool   `json:"admin"`
	Id       string `json:"id"`
}
