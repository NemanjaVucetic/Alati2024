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

func (s ConfigGroupService) Add(c model.ConfigGroup) error {
	err := s.repo.Add(c)
	if err != nil {
		return err
	}
	return nil
}

func (s ConfigGroupService) Delete(name string, version int) error {
	err := s.repo.Delete(name, version)
	if err != nil {
		return err
	}
	return nil
}

func (s ConfigGroupService) AddConfigToGroup(group model.ConfigGroup, config model.Config) error {
	err := s.repo.AddConfigToGroup(group, config)
	if err != nil {
		return err
	}
	return nil
}

func (s ConfigGroupService) RemoveConfigFromGroup(group model.ConfigGroup, key string) error {
	err := s.repo.RemoveConfigFromGroup(group, key)
	if err != nil {
		return err
	}
	return nil
}
