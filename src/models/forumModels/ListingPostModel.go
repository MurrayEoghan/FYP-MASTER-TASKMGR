package forumModels

import "time"

type ListPostModel struct {
	Id         int       `json:"id"`
	Votes      int       `json:"votes"`
	Title      string    `json:"title"`
	Date       time.Time `json:"date"`
	Author     string    `json:"author"`
	AuthorId   int       `json:"authorId"`
	TopicId    int       `json:"topic"`
	SubTopicId int       `json:"subTopic"`
	Content    string    `json:"content"`
}
