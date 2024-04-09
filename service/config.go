package service

import (
	"alati/model"
)

type ConfigService struct {
	repo model.ConfigRepository
}

func NewConfigService(repo model.ConfigRepository) ConfigService {
	return ConfigService{
		repo: repo,
	}
}

func (s ConfigService) Get(name string, version int) (model.Config, error) {
	return s.repo.Get(name, version)
}

func (s ConfigService) Add(c model.Config) {
	s.repo.Add(c)
}

func (s ConfigService) Delete(name string, version int) error {
	err := s.repo.Delete(name, version)
	if err != nil {
		return err
	}
	return nil
}
