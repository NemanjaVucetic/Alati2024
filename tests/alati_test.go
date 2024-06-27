package tests

import (
	"alati/handler"
	"alati/model"
	"alati/repo"
	"alati/service"
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel"
)

func newTestHandler() handler.ConfigHandler {
	logger := log.New(io.Discard, "", log.LstdFlags)
	tracer := otel.Tracer("test")
	repos, _ := repo.NewConfigRepo(logger, tracer)
	services := service.NewConfigService(repos, logger, tracer)
	return handler.NewConfigHandler(services, logger, tracer)
}

func TestGet(t *testing.T) {
	handler := newTestHandler()

	req := httptest.NewRequest(http.MethodGet, "/configs/db_config/2", nil)
	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/configs/{name}/{version}", handler.Get).Methods("GET")
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `{"Name":"db_config","Version":2}`
	if rr.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestGetAll(t *testing.T) {
	handler := newTestHandler()

	req := httptest.NewRequest(http.MethodGet, "/configs/", nil)
	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/configs/", handler.GetAll).Methods("GET")
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `[{"Name":"test","Version":1}]`
	if rr.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestAdd(t *testing.T) {
	handler := newTestHandler()

	config := model.Config{Name: "new_config", Version: 1}
	body, _ := json.Marshal(config)
	req := httptest.NewRequest(http.MethodPost, "/configs/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("idempotency_key", "test_key")
	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/configs/", handler.Add).Methods("POST")
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	expected := `{"Name":"new_config","Version":1}`
	if rr.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestDelete(t *testing.T) {
	handler := newTestHandler()

	req := httptest.NewRequest(http.MethodDelete, "/configs/test/1", nil)
	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/configs/{name}/{version}", handler.Delete).Methods("DELETE")
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `"Deleted"`
	if rr.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}

}

func TestDeleteAll(t *testing.T) {
	handler := newTestHandler()

	req := httptest.NewRequest(http.MethodDelete, "/configs/", nil)
	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/configs/", handler.DeleteAll).Methods("DELETE")
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}
