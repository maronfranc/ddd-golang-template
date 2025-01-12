package example

import (
	"github.com/maronfranc/poc-golang-ddd/domain/dto"
	"github.com/maronfranc/poc-golang-ddd/infrastructure"
)

type Repository struct{}

const TABLE_NAME = "examples"

func (er *Repository) GetMany(page, limit int) (*[]dto.ManyExampleResponseDto, int) {
	r, total := infrastructure.SelectPagination[dto.ManyExampleResponseDto](
		TABLE_NAME, []string{"id", "title"}, page, limit)
	return &r, total
}
func (er *Repository) GetById(id string) (*dto.CreateExampleResponseDto, error) {
	return infrastructure.SelectById[dto.CreateExampleResponseDto](
		TABLE_NAME, id, []string{"id", "title", "description"})
}
func (er *Repository) Create(e *dto.CreateExampleDto) (*dto.CreateExampleResponseDto, error) {
	id, err := infrastructure.InsertReturningId(TABLE_NAME, e)
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
	err := infrastructure.UpdateById(TABLE_NAME, id, e)
	return err
}
func (er *Repository) DeleteById(id string) error {
	err := infrastructure.DeleteById(TABLE_NAME, id)
	return err
}
