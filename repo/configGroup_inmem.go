package repo

import (
	"alati/model"
	"errors"
	"fmt"
)

type ConfigGroupInMemRepository struct {
	configs map[string]model.ConfigGroup
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

func NewConfigGroupInMemRepository() model.ConfigGroupRepository {
	return ConfigGroupInMemRepository{
		configs: make(map[string]model.ConfigGroup),
	}
}
