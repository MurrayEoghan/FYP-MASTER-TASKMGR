package repositories

import (
	"log"
	"net/http"
	models "repo/models/followerModels"
	forumModels "repo/models/forumModels"
	"repo/sqldb"
	"strings"
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
	create_notification, err5 := sqldb.DB3.Prepare(`INSERT INTO notification_service.notifications ( notification_type, viewed, initiated_by_id, initiated_by_name, recipient_name, recipient_id, cause_entity) VALUES (?,?,?,?,?,?,?)`)
	create_notification.Exec(3, 1, ids.ParentId, usernames[0].Username, usernames[1].Username, ids.SubId, 0)
	if err5 != nil {
		log.Print("Unable to create follow notification")
	}
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

func GetFollowingPosts(id int, w http.ResponseWriter) ([]forumModels.ListPostModel, error) {

	var ids []int
	var posts []forumModels.ListPostModel
	rows, err := sqldb.DB2.Query(`SELECT sub_id FROM follower.follower_relationships WHERE parent_id = ? AND following = true`, id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		log.Println(err)

		return []forumModels.ListPostModel{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var user int
		if err := rows.Scan(&user); err != nil {
			log.Println(err)
		}
		ids = append(ids, user)

	}
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		args[i] = id
	}

	stmt := `SELECT * FROM forum.listings WHERE authorId IN (?` + strings.Repeat(",?", len(args)-1) + `) ORDER BY date DESC LIMIT 5`

	rows2, err2 := sqldb.DB1.Query(stmt, args...)
	if err2 != nil {
		log.Print("Made it here")
		return []forumModels.ListPostModel{}, err2
	}

	for rows2.Next() {
		post := &forumModels.ListPostModel{}
		if err2 := rows2.Scan(&post.Id, &post.Votes, &post.Title, &post.Date, &post.Author, &post.AuthorId, &post.TopicId, &post.SubTopicId, &post.Content); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err2)
			return []forumModels.ListPostModel{}, err2
		}
		posts = append(posts, *post)
	}
	return posts, nil

}
