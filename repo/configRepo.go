package repo

import (
	"alati/model"
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/consul/api"
	"log"
	"os"
	"strconv"


	"github.com/hashicorp/consul/api"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

)

type ConfigRepo struct {
	cli    *api.Client
	logger *log.Logger
	Tracer trace.Tracer
}

func NewConfigRepo(logger *log.Logger, tracer trace.Tracer) (*ConfigRepo, error) {
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
		Tracer: tracer,
	}, nil
}

func (conf *ConfigRepo) Get(id string, ctx context.Context) (*model.Config, error) {
	_, span := conf.Tracer.Start(ctx, "r.GetConfig")
	defer span.End()
	kv := conf.cli.KV()

	pair, _, err := kv.Get(id, nil)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	if pair == nil {
		span.SetStatus(codes.Error, err.Error())
		return nil, nil
	}

	c := &model.Config{}
	err = json.Unmarshal(pair.Value, c)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return c, nil
}

func (conf *ConfigRepo) GetAll(ctx context.Context) ([]model.Config, error) {
	_, span := conf.Tracer.Start(ctx, "r.GetAllConfig")
	defer span.End()
	kv := conf.cli.KV()
	data, _, err := kv.List(all, nil)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	var configs []model.Config
	for _, pair := range data {
		var product model.Config
		err = json.Unmarshal(pair.Value, &product)
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
			return nil, err
		}
		configs = append(configs, product)
	}

	return configs, nil
}


func (conf *ConfigRepo) Put(c *model.Config, id string, ctx context.Context) (*model.Config, error) {
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

	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}


	data, err := json.Marshal(c)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	confKeyValue := &api.KVPair{Key: constructKey(c.Name, strconv.Itoa(c.Version)), Value: data}
	_, err = kv.Put(confKeyValue, nil)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return c, nil
}

func (conf *ConfigRepo) Delete(id string, ctx context.Context) error {
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

func (conf *ConfigRepo) DeleteAll(ctx context.Context) error {
	_, span := conf.Tracer.Start(ctx, "r.DeleteAllConfig")
	defer span.End()
	kv := conf.cli.KV()

	_, err := kv.DeleteTree(all, nil)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	return nil
}
