package repositories

import (
	"log"
	"repo/sqldb"

	"net/http"
	models "repo/models/followerModels"
)

func Follow(ids models.UserIds, w http.ResponseWriter) int {
	var usernames []models.Usernames
	username := &models.Usernames{}
	err := sqldb.DB.QueryRow("SELECT username FROM task_mgr.user WHERE Id = ?", ids.ParentId).Scan(&username.Username)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		log.Print("Error in select")
		return 0
	}
	usernames = append(usernames, *username)

	err2 := sqldb.DB.QueryRow("SELECT username FROM task_mgr.user WHERE Id = ?", ids.SubId).Scan(&username.Username)
	if err2 != nil {
		w.WriteHeader(http.StatusNotFound)
		log.Print("Error in select")
		return 0
	}
	usernames = append(usernames, *username)

	follow := models.Follow{}
	err4 := sqldb.DB2.QueryRow("SELECT * FROM follower.follower_relationships WHERE parent_id = ? AND sub_id = ?", ids.ParentId, ids.SubId).Scan(&follow.RelationshipId, &follow.ParentId, &follow.ParentUser, &follow.SubId, &follow.SubUser, &follow.Following, &follow.Followed)
	if err4 != nil {
		create_primary_relationship, err2 := sqldb.DB2.Prepare(`INSERT INTO follower.follower_relationships (parent_id, parent_user, sub_id, sub_user, following, followed) VALUES (?, ?, ?, ?, ?, ?)`)
		create_primary_relationship.Exec(ids.ParentId, usernames[0].Username, ids.SubId, usernames[1].Username, true, false)
		if err2 != nil {
			w.WriteHeader(http.StatusConflict)
			log.Print("Error inserting")
			return 0
		}
		create_sub_relationship, err3 := sqldb.DB2.Prepare(`INSERT INTO follower.follower_relationships (parent_id, parent_user, sub_id, sub_user, following, followed) VALUES (?, ?, ?, ?, ?, ?)`)
		create_sub_relationship.Exec(ids.SubId, usernames[1].Username, ids.ParentId, usernames[0].Username, false, true)
		if err3 != nil {
			w.WriteHeader(http.StatusConflict)
			log.Print("Error inserting")
			return 0
		}

		return 1
	} else {
		if follow.Following == false || follow.Followed == false {
			stmt, err := sqldb.DB2.Prepare(`UPDATE follower.follower_relationships SET following = ? WHERE parent_id = ? AND sub_id = ?`)
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				log.Print(err)
				return 0
			}
			stmt.Exec(true, ids.ParentId, ids.SubId)

			stmt2, err2 := sqldb.DB2.Prepare(`UPDATE follower.follower_relationships SET followed = ? WHERE parent_id = ? AND sub_id = ?`)
			if err2 != nil {
				w.WriteHeader(http.StatusNotFound)
				log.Print(err)
				return 0
			}

			stmt2.Exec(true, ids.SubId, ids.ParentId)
			return 1

		} else {
			stmt, err := sqldb.DB2.Prepare(`UPDATE follower.follower_relationships SET followed = ? WHERE parent_id = ? AND sub_id = ?`)
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				log.Print(err)
				return 0
			}
			stmt.Exec(true, ids.ParentId, ids.SubId)

			stmt2, err2 := sqldb.DB2.Prepare(`UPDATE follower.follower_relationships SET following = ? WHERE parent_id = ? AND sub_id = ?`)
			if err2 != nil {
				w.WriteHeader(http.StatusNotFound)
				log.Print(err)
				return 0
			}

			stmt2.Exec(true, ids.SubId, ids.ParentId)
			return 1
		}
	}
}

func GetFollowers(userId int) []models.Follower {
	var followers []models.Follower
	rows, err := sqldb.DB2.Query(`SELECT sub_user, sub_id FROM follower.follower_relationships WHERE parent_id = ? AND followed = true`, userId)
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()

	for rows.Next() {
		user := &models.Follower{}
		if err := rows.Scan(&user.Username, &user.Id); err != nil {
			log.Println(err)
		}
		followers = append(followers, *user)

	}
	return followers
}

func GetFollowing(userId int) []models.Follower {
	log.Print(userId)
	var followers []models.Follower
	rows, err := sqldb.DB2.Query(`SELECT sub_user, sub_id FROM follower.follower_relationships WHERE parent_id = ? AND following = true`, userId)
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()
	for rows.Next() {
		user := &models.Follower{}
		if err := rows.Scan(&user.Username, &user.Id); err != nil {
			log.Println(err)
		}
		log.Print(user.Username)
		followers = append(followers, *user)
	}
	return followers
}

func UnFollow(ids models.UserIds, w http.ResponseWriter) int64 {
	var usernames []models.Usernames
	username := &models.Usernames{}
	err := sqldb.DB.QueryRow("SELECT username FROM task_mgr.user WHERE Id = ?", ids.ParentId).Scan(&username.Username)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		log.Print("Error in select")
		return 0
	}
	usernames = append(usernames, *username)

	err2 := sqldb.DB.QueryRow("SELECT username FROM task_mgr.user WHERE Id = ?", ids.SubId).Scan(&username.Username)
	if err2 != nil {
		w.WriteHeader(http.StatusNotFound)
		log.Print("Error in select")
		return 0
	}
	usernames = append(usernames, *username)

	follow := models.Follow{}
	err4 := sqldb.DB2.QueryRow("SELECT * FROM follower.follower_relationships WHERE parent_id = ? AND sub_id = ?", ids.ParentId, ids.SubId).Scan(&follow.RelationshipId, &follow.ParentId, &follow.ParentUser, &follow.SubId, &follow.SubUser, &follow.Following, &follow.Followed)
	if err4 != nil {

		w.WriteHeader(http.StatusNotFound)
		log.Print(err)
		return 0

	} else {
		if follow.Following == true || follow.Followed == true {

			delete_primary_relationship, err := sqldb.DB2.Prepare(`UPDATE follower.follower_relationships SET following = ? WHERE parent_id = ? AND sub_id = ?`)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Print(err)
				return 0
			}
			delete_primary_relationship.Exec(false, ids.ParentId, ids.SubId)

			delete_secondary_relationship, err2 := sqldb.DB2.Prepare(`UPDATE follower.follower_relationships SET followed = ? WHERE parent_id = ? and sub_id = ?`)
			if err2 != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Print(err)
				return 0
			}

			delete_secondary_relationship.Exec(false, ids.SubId, ids.ParentId)
			return 1

		} else {
			log.Print("Relationship doesnt exist")
			return 0
		}
	}
}
