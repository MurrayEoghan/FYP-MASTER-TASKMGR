package main

import (
	r "repo/settings"
	db "repo/sqldb"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db.ConnectDB()
	db.ConnectForumDB()
	db.ConnectFollowerDB()
	db.ConnectNotificationDB()
	r.Router()
}
