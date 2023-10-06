package infrastructure

import (
	"github.com/maronfranc/poc-golang-ddd/application/dto"
)

type ExampleRepository struct{}

const TABLE_NAME = "examples"

func (er *ExampleRepository) GetMany(page, limit int) (*[]dto.ManyExampleResponseDto, int) {
	r, total := SelectPagination[dto.ManyExampleResponseDto](
		TABLE_NAME, []string{"id", "title"}, page, limit)
	return &r, total
}
func (er *ExampleRepository) GetById(id string) (*dto.CreateExampleResponseDto, error) {
	return SelectById[dto.CreateExampleResponseDto](
		TABLE_NAME, id, []string{"id", "title", "description"})
}
func (er *ExampleRepository) Create(e *dto.CreateExampleDto) (*dto.CreateExampleResponseDto, error) {
	id, err := InsertReturningId(TABLE_NAME, e)
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
func (er *ExampleRepository) UpdateById(id string, e *dto.CreateExampleDto) error {
	err := UpdateById(TABLE_NAME, id, e)
	return err
}
func (er *ExampleRepository) DeleteById(id string) error {
	err := DeleteById(TABLE_NAME, id)
	return err
}
