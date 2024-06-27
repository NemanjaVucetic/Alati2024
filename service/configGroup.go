package service

import (
	"alati/model"
	"alati/repo"
	"context"
	"log"

	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type ConfigGroupService struct {
	repo   model.ConfigGroupRepository
	logger *log.Logger
	Tracer trace.Tracer
}

func NewConfigGroupService(repo *repo.ConfigGroupRepo, logger *log.Logger, tracer trace.Tracer) ConfigGroupService {
	return ConfigGroupService{
		repo:   repo,
		logger: logger,
		Tracer: tracer,
	}
}

func (s ConfigGroupService) GetAll(ctx context.Context) ([]model.ConfigGroup, error) {
	ctx2, span := s.Tracer.Start(ctx, "s.GetAllConfig")
	defer span.End()
	return s.repo.GetAll(ctx2)
}

func (s ConfigGroupService) Get(id string, ctx context.Context) (*model.ConfigGroup, error) {
	ctx2, span := s.Tracer.Start(ctx, "s.GetConfig")
	defer span.End()
	return s.repo.Get(id, ctx2)
}

func (s ConfigGroupService) Add(c *model.ConfigGroup, id string, ctx context.Context) (*model.ConfigGroup, error) {
	ctx2, span := s.Tracer.Start(ctx, "s.AddConfig")
	defer span.End()
	group, err := s.repo.Put(c, id, ctx2)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	return group, nil
}

func (s ConfigGroupService) Delete(id string, ctx context.Context) error {
	ctx2, span := s.Tracer.Start(ctx, "s.DeleteConfig")
	defer span.End()
	err := s.repo.Delete(id, ctx2)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	return nil
}

func (s ConfigGroupService) AddConfigToGroup(group model.ConfigGroup, config model.Config, id string, ctx context.Context) error {
	ctx2, span := s.Tracer.Start(ctx, "s.AddConfigToGroup")
	defer span.End()
	err := s.repo.AddConfigToGroup(group, config, id, ctx2)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	return nil
}

func (s ConfigGroupService) RemoveConfigFromGroup(group model.ConfigGroup, config model.Config, id string, ctx context.Context) error {
	ctx2, span := s.Tracer.Start(ctx, "s.RemoveConfigFromGroup")
	defer span.End()
	err := s.repo.RemoveConfigFromGroup(group, config, id, ctx2)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	return nil
}

func (s ConfigGroupService) GetConfigsByLabels(prefixGroup string, prefixConf string, ctx context.Context) ([]model.Config, error) {
	ctx2, span := s.Tracer.Start(ctx, "s.GetConfigByLabels")
	defer span.End()
	configs, err := s.repo.GetConfigsByLabels(prefixGroup, prefixConf, ctx2)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	return configs, nil
}

func (s ConfigGroupService) DeleteConfigsByLabels(prefixGroup string, prefixConf string, ctx context.Context) error {
	ctx2, span := s.Tracer.Start(ctx, "s.DeleteConfigByLabels")
	defer span.End()
	err := s.repo.DeleteConfigsByLabels(prefixGroup, prefixConf, ctx2)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	return nil
}
