package worlds

type World struct {
	UserId     string   `json:"userId" dynamodbav:"userId"`
	Id         string   `json:"id" dynamodbav:"id"`
	Name       string   `json:"name" dynamodbav:"name"`
	Intro      string   `json:"intro" dynamodbav:"intro"`
	Genres     []string `json:"genres" dynamodbav:"genres,stringset,omitempty"`
	CoverImage string   `json:"coverImage" dynamodbav:"coverImage"`
	CreatedAt  string   `json:"createdAt" dynamodbav:"createdAt"`
	UpdatedAt  string   `json:"updatedAt" dynamodbav:"updatedAt"`
}
