package example

import (
	"github.com/maronfranc/poc-golang-ddd/domain/dto"
)

var exampleRepository = Repository{}

type Service struct{}

func (es *Service) GetMany(page, limit int) (*[]dto.ManyExampleResponseDto, int) {
	return exampleRepository.GetMany(page, limit)
}

func (es *Service) GetById(id string) (*dto.CreateExampleResponseDto, error) {
	return exampleRepository.GetById(id)
}

func (es *Service) Create(e *dto.CreateExampleDto) (*dto.CreateExampleResponseDto, error) {
	return exampleRepository.Create(e)
}

func (es *Service) UpdateById(id string, e *dto.CreateExampleDto) error {
	return exampleRepository.UpdateById(id, e)
}

func (es *Service) DeleteById(id string) error {
	return exampleRepository.DeleteById(id)
}
