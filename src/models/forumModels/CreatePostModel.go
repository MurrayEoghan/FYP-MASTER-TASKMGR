package forumModels

type CreatePostModel struct {
	Title      string `json:"title"`
	Author     string `json:"author"`
	AuthorId   int    `json:"authorId"`
	TopicId    int    `json:"topic"`
	SubTopicId int    `json:"subTopic"`
	Content    string `json:"content"`
}
