package service

import "alati/model"

type ConfigGroupService struct {
	repo model.ConfigGroupRepository
}

func NewConfigGroupService(repo model.ConfigGroupRepository) ConfigGroupService {
	return ConfigGroupService{
		repo: repo,
	}
}

func (s ConfigGroupService) Get(name string, version int) (model.ConfigGroup, error) {
	return s.repo.Get(name, version)
}

func (s ConfigGroupService) Add(c model.ConfigGroup) {
	s.repo.Add(c)
}

func (s ConfigGroupService) Delete(name string, version int) {
	err := s.repo.Delete(name, version)
	if err != nil {
		return
	}
}
