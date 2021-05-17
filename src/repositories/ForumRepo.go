package repositories

import (
	"context"
	"fmt"
	"log"
	"net/http"
	model "repo/models/forumModels"
	"repo/sqldb"
	"time"
)

type ForumRepo interface {
	GetAllPosts()
}

var (
	ctx context.Context
)

func GetAllPosts() []model.ListPostModel {
	var listings []model.ListPostModel
	rows, err := sqldb.DB1.Query("SELECT * FROM forum.listings")
	if err != nil {
		log.Println(err)

	}
	defer rows.Close()

	for rows.Next() {
		listing := &model.ListPostModel{}
		if err := rows.Scan(&listing.Id, &listing.Votes, &listing.Title, &listing.Date, &listing.Author, &listing.AuthorId, &listing.TopicId, &listing.SubTopicId, &listing.Content); err != nil {
			log.Println(err)
		}
		listings = append(listings, *listing)
	}
	return listings

}

func GetRecentPosts(id int, w http.ResponseWriter) []model.ListPostModel {
	var listings []model.ListPostModel
	rows, err := sqldb.DB1.Query("SELECT  * FROM forum.listings WHERE authorId = ? ORDER BY post_id DESC LIMIT 5;", id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return listings
	}
	for rows.Next() {
		listing := &model.ListPostModel{}
		if err := rows.Scan(&listing.Id, &listing.Votes, &listing.Title, &listing.Date, &listing.Author, &listing.AuthorId, &listing.TopicId, &listing.SubTopicId, &listing.Content); err != nil {
			log.Println(err)
		}
		listings = append(listings, *listing)
	}
	return listings
}

func CreatePost(newPost model.CreatePostModel, w http.ResponseWriter) int64 {
	currentTime := time.Now()
	stmt, err := sqldb.DB1.Prepare(`INSERT INTO forum.listings (votes, title, date, author, authorId, topic_id, sub_topic_id, content) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`)
	postRow, err := stmt.Exec(0, newPost.Title, currentTime.Format("2006-01-02"), newPost.Author, newPost.AuthorId, newPost.TopicId, newPost.SubTopicId, newPost.Content)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error Creating Post")
		return 0
	}

	newPostId, err := postRow.LastInsertId()
	fmt.Printf("New Post ID : %d", newPostId)
	return newPostId

}

func GetTopics() []model.Topics {
	var topics []model.Topics
	rows, err := sqldb.DB1.Query("SELECT * FROM forum.topics JOIN  forum.sub_topics WHERE forum.topics.topic_id = forum.sub_topics.parent_topic_id")
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()

	for rows.Next() {
		topic := &model.Topics{}
		if err := rows.Scan(&topic.TopicId, &topic.Topic, &topic.SubTopicId, &topic.SubTopic, &topic.ParentTopicId); err != nil {
			log.Println(err)
		}
		topics = append(topics, *topic)
	}
	return topics
}

func GetPost(postId int, w http.ResponseWriter) *model.Post {

	post := &model.Post{}

	err := sqldb.DB1.QueryRow("SELECT post_id, title, date, author, authorId, content FROM forum.listings WHERE post_id = ?", postId).Scan(&post.PostId, &post.Title, &post.Date, &post.Author, &post.AuthorId, &post.Content)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		log.Println(err)
		return &model.Post{}
	}

	rows, err := sqldb.DB1.Query("SELECT * FROM forum.comments WHERE parent_post_id = ?", postId)
	if err != nil {
		log.Println(err)
		post.Comments = make([]model.Comment, 0)
		return post
	}
	defer rows.Close()

	for rows.Next() {
		comment := &model.Comment{}
		if err := rows.Scan(&comment.Id, &comment.Comment, &comment.ParentPostId, &comment.CommentAuthor, &comment.CommentAuthorId, &comment.Date); err != nil {
			log.Println(err)

		}
		if err2 := sqldb.DB.QueryRow("SELECT user.profession_id FROM task_mgr.user WHERE Id = ?", comment.CommentAuthorId).Scan(&comment.CommentAuthorProfessionId); err2 != nil {
			log.Println(err2)
		}

		post.Comments = append(post.Comments, *comment)
	}

	err2 := sqldb.DB.QueryRow("SELECT user.profession_id FROM task_mgr.user WHERE Id = ?", post.AuthorId).Scan(&post.ProfessionId)
	if err2 != nil {
		w.WriteHeader(http.StatusNotFound)
		log.Println(err2)
		return &model.Post{}
	}

	err3 := sqldb.DB1.QueryRow("SELECT answer_id, answer_author_id, answer_author_profession, answer, answer_author FROM forum.answers WHERE parent_post_id = ?", postId).Scan(&post.AnswerId, &post.AnswerAuthorId, &post.AnswerAuthorProfession, &post.Answer, &post.AnswerAuthor)
	if err3 != nil {
		log.Println(err)
		return post
	}
	return post
}

func CreateComment(newComment model.CreateComment, w http.ResponseWriter) int64 {

	post := &model.Post{}

	err := sqldb.DB1.QueryRow("SELECT post_id, title, date, author, authorId, content FROM forum.listings WHERE post_id = ?", newComment.ParentPostId).Scan(&post.PostId, &post.Title, &post.Date, &post.Author, &post.AuthorId, &post.Content)
	if err != nil {
		log.Println(err)
	}
	if post.AuthorId != newComment.CommentAuthorId {
		insert, err := sqldb.DB3.Prepare("INSERT INTO notification_service.notifications (notification_type, viewed, initiated_by_id, initiated_by_name, recipient_name, recipient_id, cause_entity) VALUES (?,?,?,?,?,?,?)")
		insert.Exec(1, 1, newComment.CommentAuthorId, newComment.CommentAuthor, post.Author, post.AuthorId, post.PostId)
		if err != nil {
			log.Print("Error creating comment notification")
		}
	}

	currentTime := time.Now()
	stmt, err := sqldb.DB1.Prepare(`INSERT INTO forum.comments (comment, parent_post_id, comment_author, comment_author_id, post_date) VALUES (?, ?, ?, ?, ?)`)
	postRow, err := stmt.Exec(newComment.Comment, newComment.ParentPostId, newComment.CommentAuthor, newComment.CommentAuthorId, currentTime.Format("2006-01-02"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error Posting Comment")
		return 0
	}

	newPostId, err := postRow.LastInsertId()
	fmt.Printf("New Comment ID : %d", newPostId)
	return newPostId

}

func UpdatePost(newPost model.UpdatePost, w http.ResponseWriter) int64 {
	check, err := sqldb.DB.Query(`SELECT * FROM forum.listings WHERE post_id = ? AND authorId = ?`, newPost.PostId, newPost.PostAuthorId)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return 0
	}
	if check.Next() == false {
		w.WriteHeader(http.StatusNotFound)
		return 0
	}

	stmt, err := sqldb.DB.Prepare(`UPDATE forum.listings SET content = ? WHERE post_id = ? AND authorId = ?`)
	if err != nil {

		w.WriteHeader(http.StatusNotFound)
		log.Fatal(err)
		return 0
	}
	stmt.Exec(newPost.Content, newPost.PostId, newPost.PostAuthorId)
	return 1
}

func UpdateComment(newComment model.UpdateComment, w http.ResponseWriter) int64 {

	check, err := sqldb.DB.Query(`SELECT * FROM forum.comments WHERE id = ? AND comment_author_id = ?`, newComment.CommentId, newComment.CommentAuthorId)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return 0
	}
	if check.Next() == false {
		w.WriteHeader(http.StatusNotFound)
		return 0
	}

	stmt, err := sqldb.DB.Prepare(`UPDATE forum.comments SET comment = ? WHERE id = ? AND comment_author_id = ?`)
	if err != nil {

		w.WriteHeader(http.StatusInternalServerError)
		log.Fatal(err)
		return 0
	}
	stmt.Exec(newComment.Comment, newComment.CommentId, newComment.CommentAuthorId)
	return 1
}

func DeleteComment(deleteCommentId int, w http.ResponseWriter) (int64, error) {
	result, err := sqldb.DB.Exec(`DELETE FROM forum.comments WHERE id = ?`, deleteCommentId)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return 0, err
	} else {
		return result.LastInsertId()
	}
}

func DeletePost(deleteCommentId int, w http.ResponseWriter) (int64, error) {
	post, err := sqldb.DB.Exec(`DELETE FROM forum.listings WHERE post_id = ?`, deleteCommentId)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return 0, err
	} else {
		return post.LastInsertId()
	}
}

func DeleteComments(deleteCommentId int, w http.ResponseWriter) (int64, error) {
	result, err := sqldb.DB.Exec(`DELETE FROM forum.comments WHERE parent_post_id = ?`, deleteCommentId)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return 0, err
	} else {
		return result.LastInsertId()
	}
}

func UpdateAnswer(newAnswer model.UpdateAnswer, w http.ResponseWriter) int64 {

	check, err := sqldb.DB.Query(`SELECT * FROM forum.answers WHERE parent_post_id = ? AND answer_id = ?`, newAnswer.PostId, newAnswer.AnswerId)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return 0
	}
	if check.Next() == false {
		w.WriteHeader(http.StatusNotFound)
		return 0
	}

	stmt, err := sqldb.DB.Prepare(`UPDATE forum.answers SET answer = ? WHERE parent_post_id = ? AND answer_id = ?`)
	if err != nil {

		w.WriteHeader(http.StatusInternalServerError)
		log.Fatal(err)
		return 0
	}
	stmt.Exec(newAnswer.NewAnswer, newAnswer.PostId, newAnswer.AnswerId)
	return 1
}

func DeleteAnswer(deleteAnswerId int, w http.ResponseWriter) (int64, error) {
	post, err := sqldb.DB.Exec(`DELETE FROM forum.answers WHERE parent_post_id = ?`, deleteAnswerId)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return 0, err
	} else {
		return post.LastInsertId()
	}
}

func CreateAnswer(newAnswer model.CreateAnswer, w http.ResponseWriter) int64 {

	post := &model.Post{}

	err := sqldb.DB1.QueryRow("SELECT post_id, title, date, author, authorId, content FROM forum.listings WHERE post_id = ?", newAnswer.ParentPostId).Scan(&post.PostId, &post.Title, &post.Date, &post.Author, &post.AuthorId, &post.Content)
	if err != nil {
		log.Println(err)
	}
	insert, err := sqldb.DB3.Prepare("INSERT INTO notification_service.notifications (notification_type, viewed, initiated_by_id, initiated_by_name, recipient_name, recipient_id, cause_entity) VALUES (?,?,?,?,?,?,?)")
	insert.Exec(2, 1, newAnswer.AuthorId, newAnswer.Author, post.Author, post.AuthorId, post.PostId)
	if err != nil {
		log.Print("Error creating answer notification")
	}

	stmt, err := sqldb.DB1.Prepare(`INSERT INTO forum.answers (parent_post_id, answer_author_id, answer_author_profession, answer, answer_author) VALUES (?, ?, ?, ?, ?)`)
	postRow, err := stmt.Exec(newAnswer.ParentPostId, newAnswer.AuthorId, newAnswer.AuthorProfession, newAnswer.Answer, newAnswer.Author)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error Posting Comment")
		return 0
	}

	newPostId, err := postRow.LastInsertId()

	return newPostId

}
