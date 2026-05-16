package example

import (
	"github.com/maronfranc/poc-golang-ddd/infrastructure/database"
	"github.com/maronfranc/poc-golang-ddd/infrastructure/model"
)

type Repository struct{}

func NewRepository() *Repository {
	return &Repository{}
}

const TABLE_NAME = "examples"

func (er *Repository) GetMany(page, limit int) (*[]model.ManyExampleResponseDto, int) {
	r, total, _ := database.SelectManyAndCount[model.ManyExampleResponseDto](
		TABLE_NAME, page, limit)
	return &r, *total
}

func (er *Repository) GetById(id string) (*model.Example, error) {
	return database.SelectById[model.Example](TABLE_NAME, id)
}

func (er *Repository) Create(e *model.CreateExampleDto) (*model.Example, error) {
	id, err := database.InsertReturningId(TABLE_NAME, e)
	if err != nil {
		return nil, err
	}

	r := &model.Example{
		Id:          id,
		Title:       e.Title,
		Description: e.Description,
	}
	return r, nil
}

func (er *Repository) UpdateById(id string, e *model.CreateExampleDto) error {
	err := database.UpdateById(TABLE_NAME, id, e, []string{})
	return err
}

func (er *Repository) DeleteById(id string) error {
	err := database.DeleteById(TABLE_NAME, id)
	return err
}
