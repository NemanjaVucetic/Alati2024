package repo

import (
	"alati/model"
	"errors"
	"fmt"
)

type ConfigGroupInMemRepository struct {
	configs map[string]model.ConfigGroup
	//configRepo model.ConfigRepository
}

func (c ConfigGroupInMemRepository) Add(config model.ConfigGroup) error {
	key := fmt.Sprintf("%s/%d", config.Name, config.Version)
	c.configs[key] = config
	return nil
}

func (c ConfigGroupInMemRepository) Get(name string, version int) (model.ConfigGroup, error) {
	key := fmt.Sprintf("%s/%d", name, version)
	config, ok := c.configs[key]
	if !ok {
		return model.ConfigGroup{}, errors.New("config not found")
	}
	return config, nil
}

func (c ConfigGroupInMemRepository) Delete(name string, version int) error {
	key := fmt.Sprintf("%s/%d", name, version)

	if _, ok := c.configs[key]; !ok {
		return errors.New("config not found")
	}

	delete(c.configs, key)

	return nil
}

func (c ConfigGroupInMemRepository) AddConfigToGroup(group model.ConfigGroup, config model.Config) error {
	key := fmt.Sprintf("%s/%d", config.Name, config.Version)
	group.Configs[key] = config
	return nil
}

func (c ConfigGroupInMemRepository) RemoveConfigFromGroup(group model.ConfigGroup, key string) error {
	fmt.Println(key)
	delete(group.Configs, key)
	return nil
}

func (c ConfigGroupInMemRepository) GetConfigsByLabels(group model.ConfigGroup, labels *map[string]string) ([]model.Config, error) {
	var filteredConfigs []model.Config

	for _, conf := range group.Configs {
		for keyC, valueC := range conf.Labels {
			for keyL, valueL := range *labels {
				if keyC == keyL && valueC == valueL {
					filteredConfigs = append(filteredConfigs, conf)
				}
			}
		}
	}

	return filteredConfigs, nil
}

func (c ConfigGroupInMemRepository) DeleteConfigsByLabels(group model.ConfigGroup, labels *map[string]string) error {

	for _, conf := range group.Configs {
		for keyC, valueC := range conf.Labels {
			for keyL, valueL := range *labels {
				if keyC == keyL && valueC == valueL {
					key := fmt.Sprintf("%s/%d", conf.Name, conf.Version)
					//c.configRepo.Delete(conf.Name, conf.Version)
					c.RemoveConfigFromGroup(group, key)
				}
			}
		}
	}

	return nil
}

func NewConfigGroupInMemRepository() model.ConfigGroupRepository {
	return ConfigGroupInMemRepository{
		configs: make(map[string]model.ConfigGroup),
	}
}
