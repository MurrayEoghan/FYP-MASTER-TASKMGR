package forumModels

import "time"

type Post struct {
	PostId                 int       `json:"post_id"`
	Title                  string    `json:"title"`
	Date                   time.Time `json:"date"`
	Author                 string    `json:"author"`
	AuthorId               int       `json:"authorId"`
	ProfessionId           int       `json:"profession_id,omitempty"`
	Content                string    `json:"content"`
	Comments               []Comment `json:"comments"`
	Answer                 string    `json:"answer"`
	AnswerId               int       `json:"answer_id"`
	AnswerAuthor           string    `json:"answer_author"`
	AnswerAuthorId         int       `json:"answer_author_id"`
	AnswerAuthorProfession int       `json:"answer_author_profession"`
}
