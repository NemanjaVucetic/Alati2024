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
	ctx2, span := s.Tracer.Start(ctx, "s.GetConfig")
	defer span.End()
	return s.repo.Get(id, ctx2)
}

func (s ConfigService) GetAll(ctx context.Context) ([]model.Config, error) {
	ctx2, span := s.Tracer.Start(ctx, "s.GetAllConfigs")
	defer span.End()
	return s.repo.GetAll(ctx2)
}

func (s ConfigService) Add(c *model.Config, id string, ctx context.Context) (*model.Config, error) {
	ctx2, span := s.Tracer.Start(ctx, "s.AddConfig")
	defer span.End()
	put, err := s.repo.Put(c, id, ctx2)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	return put, nil
}

func (s ConfigService) Delete(id string, ctx context.Context) error {
	ctx2, span := s.Tracer.Start(ctx, "s.DeleteConfig")
	defer span.End()
	err := s.repo.Delete(id, ctx2)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	return nil
}

func (s ConfigService) DeleteAll(ctx context.Context) error {
	ctx2, span := s.Tracer.Start(ctx, "s.DeleteAllConfigs")
	defer span.End()
	err := s.repo.DeleteAll(ctx2)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	return nil
}
