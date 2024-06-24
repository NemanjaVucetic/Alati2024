package service

import (
	"alati/model"
	"alati/repo"
	"log"
)

type ConfigGroupService struct {
	repo   model.ConfigGroupRepository
	logger *log.Logger
}

func NewConfigGroupService(repo *repo.ConfigGroupRepo, logger *log.Logger) ConfigGroupService {
	return ConfigGroupService{
		repo:   repo,
		logger: logger,
	}
}

func (s ConfigGroupService) GetAll() ([]model.ConfigGroup, error) {
	return s.repo.GetAll()
}

func (s ConfigGroupService) Get(id string) (*model.ConfigGroup, error) {
	return s.repo.Get(id)
}

func (s ConfigGroupService) Add(c *model.ConfigGroup, id string) (*model.ConfigGroup, error) {
	group, err := s.repo.Put(c, id)
	if err != nil {
		return nil, err
	}
	return group, nil
}

func (s ConfigGroupService) Delete(id string) error {
	err := s.repo.Delete(id)
	if err != nil {
		return err
	}
	return nil
}

func (s ConfigGroupService) AddConfigToGroup(group model.ConfigGroup, config model.Config, id string) error {
	err := s.repo.AddConfigToGroup(group, config, id)
	if err != nil {
		return err
	}
	return nil
}

func (s ConfigGroupService) RemoveConfigFromGroup(group model.ConfigGroup, config model.Config, id string) error {
	err := s.repo.RemoveConfigFromGroup(group, config, id)
	if err != nil {
		return err
	}
	return nil
}

func (s ConfigGroupService) GetConfigsByLabels(prefixGroup string, prefixConf string) ([]model.Config, error) {
	configs, err := s.repo.GetConfigsByLabels(prefixGroup, prefixConf)
	if err != nil {
		return nil, err
	}
	return configs, nil
}

func (s ConfigGroupService) DeleteConfigsByLabels(prefixGroup string, prefixConf string) error {
	err := s.repo.DeleteConfigsByLabels(prefixGroup, prefixConf)
	if err != nil {
		return err
	}
	return nil
}
