package service

import (
	"alati/model"
	"log"
)

type ConfigService struct {
	repo   model.ConfigRepository
	logger *log.Logger
}

func NewConfigService(repo model.ConfigRepository, logger *log.Logger) ConfigService {
	return ConfigService{
		repo:   repo,
		logger: logger,
	}
}

func (s ConfigService) Get(id string) (*model.Config, error) {
	return s.repo.Get(id)
}

func (s ConfigService) GetAll() ([]model.Config, error) {
	return s.repo.GetAll()
}

func (s ConfigService) Add(c *model.Config, id string) (*model.Config, error) {
	put, err := s.repo.Put(c, id)
	if err != nil {
		return nil, err
	}
	return put, nil
}

func (s ConfigService) Delete(id string) error {
	err := s.repo.Delete(id)
	if err != nil {
		return err
	}
	return nil
}

func (s ConfigService) DeleteAll() error {
	err := s.repo.DeleteAll()
	if err != nil {
		return err
	}
	return nil
}
