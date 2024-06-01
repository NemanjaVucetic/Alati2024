package repo

import (
	"alati/model"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/consul/api"
	"log"
	"os"
	"strconv"
)

type ConfigRepo struct {
	cli    *api.Client
	logger *log.Logger
}

func NewConfigRepo(logger *log.Logger) (*ConfigRepo, error) {
	db := os.Getenv("DB")
	dbport := os.Getenv("DBPORT")

	config := api.DefaultConfig()
	config.Address = fmt.Sprintf("%s:%s", db, dbport)
	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}

	return &ConfigRepo{
		cli:    client,
		logger: logger,
	}, nil
}

func (conf *ConfigRepo) Get(id string) (*model.Config, error) {
	kv := conf.cli.KV()

	pair, _, err := kv.Get(id, nil)
	if err != nil {
		return nil, err
	}

	if pair == nil {
		return nil, nil
	}

	c := &model.Config{}
	err = json.Unmarshal(pair.Value, c)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (conf *ConfigRepo) GetAll() ([]model.Config, error) {
	kv := conf.cli.KV()
	data, _, err := kv.List(all, nil)
	if err != nil {
		return nil, err
	}

	var configs []model.Config
	for _, pair := range data {
		var product model.Config
		err = json.Unmarshal(pair.Value, &product)
		if err != nil {
			return nil, err
		}
		configs = append(configs, product)
	}

	return configs, nil
}

func (conf *ConfigRepo) Put(c *model.Config) (*model.Config, error) {
	kv := conf.cli.KV()

	data, err := json.Marshal(c)
	if err != nil {
		return nil, err
	}

	confKeyValue := &api.KVPair{Key: constructKey(c.Name, strconv.Itoa(c.Version)), Value: data}
	_, err = kv.Put(confKeyValue, nil)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (conf *ConfigRepo) Delete(id string) error {
	kv := conf.cli.KV()

	_, err := kv.Delete(id, nil)
	if err != nil {
		return err
	}

	return nil
}

func (conf *ConfigRepo) DeleteAll() error {
	kv := conf.cli.KV()

	_, err := kv.DeleteTree(all, nil)
	if err != nil {
		return err
	}

	return nil
}
