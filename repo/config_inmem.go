package repo

import (
	"alati/model"
)

type ConfigInMemRepository struct {
	config map[string]model.Config
}

func (c2 ConfigInMemRepository) Get(key string) model.Config {
	c := c2.config[key]
	return c
}

func (c2 ConfigInMemRepository) Add(c model.Config) {
	c2.config[c.GenerateKey()] = c
}

func (c2 ConfigInMemRepository) Delete(key string) {
	delete(c2.config, key)
}

func NewConfigInMemRepository() model.ConfigRepository {
	return ConfigInMemRepository{
		config: make(map[string]model.Config),
	}
}
