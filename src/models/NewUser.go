package models

type NewUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Fname    string `json:"fname"`
	Lname    string `json:"lname"`
	Age      int    `json:"age"`
	Sex      string `json:"sex"`
	Address1 string `json:"addr1"`
	Address2 string `json:"addr2"`
	Address3 string `json:"addr3"`
	County   string `json:"county"`
	Country  string `json:"country"`
}
