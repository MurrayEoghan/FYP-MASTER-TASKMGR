package forumModels

type UpdatePost struct {
	PostId       int    `json:"post_id"`
	PostAuthorId int    `json:"author_id"`
	Content      string `json:"content"`
}
