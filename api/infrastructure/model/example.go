package model

type Example struct {
	Id          string `db:"id" json:"id,omitempty"`
	Title       string `db:"title" json:"title,omitempty"`
	Description string `db:"description" json:"description,omitempty"`
}

type CreateExampleDto struct {
	Title       string `db:"title" json:"title,omitempty"`
	Description string `db:"description" json:"description,omitempty"`
}

type ManyExampleResponseDto struct {
	Id    string `db:"id" json:"id"`
	Title string `db:"title" json:"title"`
}
