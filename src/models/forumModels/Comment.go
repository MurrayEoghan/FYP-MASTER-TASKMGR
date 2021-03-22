package forumModels

import "time"

type Comment struct {
	Id                        int       `json:"comment_id"`
	Comment                   string    `json:"comment"`
	ParentPostId              int       `json:"parent_post_id"`
	CommentAuthor             string    `json:"comment_author"`
	CommentAuthorId           int       `json:"comment_author_id"`
	CommentAuthorProfessionId int       `json:"comment_author_profession_id,omitempty"`
	Date                      time.Time `json:"date"`
}
