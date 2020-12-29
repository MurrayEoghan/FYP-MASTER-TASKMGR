package sqldb

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func ConnectDB() {
	var err error
	db, err := sql.Open("mysql", "root:3ManorGrove0@tcp(127.0.0.1:3306)/task_mgr")
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("Welcome back Eoghan\n")
	}
	DB = db
	// defer db.Close()
}
