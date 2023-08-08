package characters

type Character struct {
	WorldId        string                 `json:"worldId" dynamodbav:"worldId"`
	Id             string                 `json:"id" dynamodbav:"id"`
	Name           string                 `json:"name" dynamodbav:"name"`
	Intro          string                 `json:"intro" dynamodbav:"intro"`
	Attributes     map[string]interface{} `json:"attributes"`
	AttributesJSON string                 `json:"-" dynamodbav:"AttributesJSON"`
	CoverImage     string                 `json:"coverImage" dynamodbav:"coverImage"`
	OwnerId        string                 `json:"ownerId" dynamodbav:"ownerId"`
	CreatedAt      string                 `json:"createdAt" dynamodbav:"createdAt"`
	UpdatedAt      string                 `json:"updatedAt" dynamodbav:"updatedAt"`
}
