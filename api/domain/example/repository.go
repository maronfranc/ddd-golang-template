package example

import (
	"github.com/maronfranc/poc-golang-ddd/domain/dto"
	"github.com/maronfranc/poc-golang-ddd/infrastructure/database"
)

type Repository struct{}

const TABLE_NAME = "examples"

func (er *Repository) GetMany(page, limit int) (*[]dto.ManyExampleResponseDto, int) {
	r, total, _ := database.SelectManyAndCount[dto.ManyExampleResponseDto](
		TABLE_NAME, page, limit)
	return &r, *total
}
func (er *Repository) GetById(id string) (*dto.CreateExampleResponseDto, error) {
	return database.SelectById[dto.CreateExampleResponseDto](TABLE_NAME, id)
}
func (er *Repository) Create(e *dto.CreateExampleDto) (*dto.CreateExampleResponseDto, error) {
	id, err := database.InsertReturningId(TABLE_NAME, e)
	if err != nil {
		return nil, err
	}
	r := &dto.CreateExampleResponseDto{
		Id:          id,
		Title:       e.Title,
		Description: e.Description,
	}
	return r, nil
}

func (er *Repository) UpdateById(id string, e *dto.CreateExampleDto) error {
	err := database.UpdateById(TABLE_NAME, id, e, []string{})
	return err
}

func (er *Repository) DeleteById(id string) error {
	err := database.DeleteById(TABLE_NAME, id)
	return err
}
