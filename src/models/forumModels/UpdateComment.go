package forumModels

type UpdateComment struct {
	CommentId       int    `json:"comment_id"`
	Comment         string `json:"comment"`
	CommentAuthorId int    `json:"comment_author_id"`
}
