package sqldb

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB
var DB1 *sql.DB
var DB2 *sql.DB
var DB3 *sql.DB

func ConnectDB() {
	var err error
	db, err := sql.Open("mysql", "root:3ManorGrove0@tcp(127.0.0.1:3306)/task_mgr")
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("Welcome back Eoghan\nUser DB connected\n")
	}
	DB = db
	// defer db.Close()
}

func ConnectForumDB() {
	var err error
	db, err := sql.Open("mysql", "root:3ManorGrove0@tcp(127.0.0.1:3306)/forum?parseTime=true")
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("Forum DB connected\n")
	}

	DB1 = db
}

func ConnectFollowerDB() {
	var err error
	db, err := sql.Open("mysql", "root:3ManorGrove0@tcp(127.0.0.1:3306)/follower?parseTime=true")
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("Follower DB Connected\n")
	}
	DB2 = db
}

func ConnectNotificationDB() {
	var err error
	db, err := sql.Open("mysql", "root:3ManorGrove0@tcp(127.0.0.1:3306)/notification_service?parseTime=true")
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("Notification DB Connected\n")
	}
	DB3 = db
}
