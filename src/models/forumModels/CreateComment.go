package forumModels

type CreateComment struct {
	Comment         string `json:"comment"`
	ParentPostId    int    `json:"parent_post_id"`
	CommentAuthor   string `json:"comment_author"`
	CommentAuthorId int    `json:"comment_author_id"`
}
