package service

import (
	"alati/model"
	"context"
	"log"

	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type ConfigService struct {
	repo   model.ConfigRepository
	logger *log.Logger
	Tracer trace.Tracer
}

func NewConfigService(repo model.ConfigRepository, logger *log.Logger, tracer trace.Tracer) ConfigService {
	return ConfigService{
		repo:   repo,
		logger: logger,
		Tracer: tracer,
	}
}

func (s ConfigService) Get(id string, ctx context.Context) (*model.Config, error) {
	_, span := s.Tracer.Start(ctx, "GetConfig")
	defer span.End()
	return s.repo.Get(id, ctx)
}

func (s ConfigService) GetAll(ctx context.Context) ([]model.Config, error) {
	_, span := s.Tracer.Start(ctx, "GetAllConfigs")
	defer span.End()
	return s.repo.GetAll(ctx)
}

func (s ConfigService) Add(c *model.Config, id string, ctx context.Context) (*model.Config, error) {
	_, span := s.Tracer.Start(ctx, "AddConfig")
	defer span.End()
	put, err := s.repo.Put(c, id, ctx)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	return put, nil
}

func (s ConfigService) Delete(id string, ctx context.Context) error {
	_, span := s.Tracer.Start(ctx, "DeleteConfig")
	defer span.End()
	err := s.repo.Delete(id, ctx)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	return nil
}

func (s ConfigService) DeleteAll(ctx context.Context) error {
	_, span := s.Tracer.Start(ctx, "DeleteAllConfigs")
	defer span.End()
	err := s.repo.DeleteAll(ctx)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	return nil
}
