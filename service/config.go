package service

import (
	"alati/model"
	"fmt"
)

type ConfigService struct {
	repo model.ConfigRepository
}

func NewConfigService(repo model.ConfigRepository) ConfigService {
	return ConfigService{
		repo: repo,
	}
}

func (s ConfigService) Get(key string) model.Config {
	return s.repo.Get(key)
}

func (s ConfigService) Add(c model.Config) {
	s.repo.Add(c)
}
func (s ConfigService) Delete(key string) {
	s.repo.Delete(key)
}

func (s ConfigService) Hello() {
	fmt.Println("poz")
}
