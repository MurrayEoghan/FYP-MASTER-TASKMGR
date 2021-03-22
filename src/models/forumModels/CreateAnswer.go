package forumModels

type CreateAnswer struct {
	ParentPostId     int    `json:"parent_post_id"`
	Answer           string `json:"answer"`
	AuthorId         int    `json:"author_id"`
	Author           string `json:"author"`
	AuthorProfession int    `json:"author_profession"`
}
