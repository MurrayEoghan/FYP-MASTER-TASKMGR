package forumModels

type UpdateAnswer struct {
	PostId    int    `json:"post_id"`
	AnswerId  int    `json:"answer_id"`
	NewAnswer string `json:"answer"`
}
