package main

import (
	"alati/handler"
	"alati/repo"
	"alati/service"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.opentelemetry.io/otel"
)

func TestGetConfig(t *testing.T) {

	logger := log.New(io.Discard, "", log.LstdFlags)
	tracer := otel.Tracer("test")
	repo, _ := repo.NewConfigRepo(logger, tracer)
	service := service.NewConfigService(repo, logger, tracer)
	h := handler.NewConfigHandler(service, logger, tracer)

	server := httptest.NewServer(http.HandlerFunc(h.GetAll))
	resp, err := http.Get(server.URL)
	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200 but got %d", resp.StatusCode)
	}

}
