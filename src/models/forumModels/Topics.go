package forumModels

type Topics struct {
	TopicId       int    `json:"topic_id"`
	Topic         string `json:"topic"`
	SubTopicId    int    `json:"sub_topic_id`
	SubTopic      string `json:"sub_topic"`
	ParentTopicId int    `json:"parent_topic_id"`
}
