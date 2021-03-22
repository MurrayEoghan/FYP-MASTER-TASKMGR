package models

type User struct {
	Username     string `json:"username"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	ProfessionId int    `json:"profession_id"`
	Id           string `json:"id"`
}
