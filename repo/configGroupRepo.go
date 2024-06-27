package repo

import (
	"alati/model"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/hashicorp/consul/api"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type ConfigGroupRepo struct {
	cli    *api.Client
	logger *log.Logger
	Tracer trace.Tracer
}

func NewConfigGroupRepo(logger *log.Logger, tracer trace.Tracer) (*ConfigGroupRepo, error) {
	db := os.Getenv("DB")
	dbport := os.Getenv("DBPORT")

	config := api.DefaultConfig()
	config.Address = fmt.Sprintf("%s:%s", db, dbport)
	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}

	return &ConfigGroupRepo{
		cli:    client,
		logger: logger,
		Tracer: tracer,
	}, nil
}

func (conf *ConfigGroupRepo) Get(id string, ctx context.Context) (*model.ConfigGroup, error) {
	_, span := conf.Tracer.Start(ctx, "r.GetConfig")
	defer span.End()
	kv := conf.cli.KV()

	pair, _, err := kv.Get(id, nil)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	if pair == nil {
		return nil, nil
	}

	c := &model.ConfigGroup{}
	err = json.Unmarshal(pair.Value, c)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return c, nil
}

func (conf *ConfigGroupRepo) GetAll(ctx context.Context) ([]model.ConfigGroup, error) {
	_, span := conf.Tracer.Start(ctx, "r.GetAllConfig")
	defer span.End()
	kv := conf.cli.KV()
	data, _, err := kv.List(allGroups, nil)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	var configs []model.ConfigGroup
	for _, pair := range data {
		var co model.ConfigGroup
		err = json.Unmarshal(pair.Value, &co)
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
			return nil, err
		}
		configs = append(configs, co)
	}

	return configs, nil
}

func (conf *ConfigGroupRepo) Put(c *model.ConfigGroup, id string, ctx context.Context) (*model.ConfigGroup, error) {
	_, span := conf.Tracer.Start(ctx, "r.AddConfig")
	defer span.End()
	kv := conf.cli.KV()
	value, _, err := kv.Get(id, nil)
	if value == nil {
		idReal, _ := json.Marshal(id)
		confKeyValue := &api.KVPair{Key: id, Value: idReal}
		kv.Put(confKeyValue, nil)
	} else {
		return nil, nil
	}

	data, err := json.Marshal(c)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	confKeyValue := &api.KVPair{Key: constructKeyGroup(c.Name, strconv.Itoa(c.Version)), Value: data}
	_, err = kv.Put(confKeyValue, nil)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return c, nil
}

func (conf *ConfigGroupRepo) Delete(id string, ctx context.Context) error {
	_, span := conf.Tracer.Start(ctx, "r.DeleteConfig")
	defer span.End()
	kv := conf.cli.KV()

	_, err := kv.Delete(id, nil)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	return nil
}

func (conf *ConfigGroupRepo) AddConfigToGroup(group model.ConfigGroup, config model.Config, id string, ctx context.Context) error {
	_, span := conf.Tracer.Start(ctx, "r.AddConfigToGroup")
	defer span.End()
	key := constructKeyInGroup(group, config)
	group.Configs[key] = config

	_, err := conf.Put(&group, id, ctx)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	return nil
}
func (conf *ConfigGroupRepo) RemoveConfigFromGroup(group model.ConfigGroup, config model.Config, id string, ctx context.Context) error {
	_, span := conf.Tracer.Start(ctx, "r.RemoveConfigFromGroup")
	defer span.End()
	key := constructKeyInGroup(group, config)
	delete(group.Configs, key)

	_, err := conf.Put(&group, id, ctx)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	return nil
}

func (conf *ConfigGroupRepo) GetConfigsByLabels(prefixGroup string, prefixConf string, ctx context.Context) ([]model.Config, error) {
	_, span := conf.Tracer.Start(ctx, "r.GetCOnfigsByLabels")
	defer span.End()
	kv := conf.cli.KV()

	data, _, err := kv.List(prefixGroup, nil)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	var allConfigs []model.Config
	for _, pair := range data {
		var configGroup model.ConfigGroup
		err = json.Unmarshal(pair.Value, &configGroup)
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
			return nil, err
		}

		for key, config := range configGroup.Configs {
			if strings.HasPrefix(key, prefixConf) {
				allConfigs = append(allConfigs, config)
			}
		}
	}

	return allConfigs, nil
}

func (conf *ConfigGroupRepo) DeleteConfigsByLabels(prefixGroup string, prefixConf string, ctx context.Context) error {
	_, span := conf.Tracer.Start(ctx, "r.DeleteCOnfigsByLabels")
	defer span.End()
	kv := conf.cli.KV()

	data, _, err := kv.List(prefixGroup, nil)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	var configGroup model.ConfigGroup
	for _, pair := range data {
		err = json.Unmarshal(pair.Value, &configGroup)
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
			return err
		}

		for key, _ := range configGroup.Configs {
			if strings.HasPrefix(key, prefixConf) {
				delete(configGroup.Configs, key)
			}
		}
	}
	group, err := json.Marshal(configGroup)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	confKeyValue := &api.KVPair{Key: constructKeyGroup(configGroup.Name, strconv.Itoa(configGroup.Version)), Value: group}
	_, err = kv.Put(confKeyValue, nil)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	return nil
}
