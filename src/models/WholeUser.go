package models

type WholeUser struct {
	Id           int    `json:"id"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	ProfessionId int    `json:"profession_id"`
	Fname        string `json:"fname"`
	Lname        string `json:"lname"`
	Age          int    `json:"age"`
	Gender       string `json:"gender"`
	Address1     string `json:"addr1"`
	Address2     string `json:"addr2"`
	Address3     string `json:"addr3"`
	County       string `json:"county"`
	Country      string `json:"country"`
}
