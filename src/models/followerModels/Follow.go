package followerModels

type Follow struct {
	RelationshipId int `json:"relationship_id"`
	ParentId int `json:"parent_id"`
	ParentUser string `json:"parent_user"`
	SubId int `json:"sub_id"`
	SubUser string `json:"sub_user"`
	Following bool `json:"following"`
	Followed bool `json:"followed"`
}