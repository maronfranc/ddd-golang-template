package dto

type CreateExampleDto struct {
	Title       string `json:"title,omitempty" bson:"title,omitempty"`
	Description string `json:"description,omitempty" bson:"description,omitempty"`
}

type CreateExampleResponseDto struct {
	Id          string `json:"id,omitempty" bson:"id,omitempty"`
	Title       string `json:"title,omitempty" bson:"title,omitempty"`
	Description string `json:"description,omitempty" bson:"description,omitempty"`
}
