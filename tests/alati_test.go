package tests

import (
	"alati/handler"
	"alati/repo"
	"alati/service"
	"bytes"
	"go.opentelemetry.io/otel"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux" // Zamislite da koristimo gorilla/mux za rute
)

func newTestHandler() handler.ConfigHandler {
	logger := log.New(io.Discard, "", log.LstdFlags)
	tracer := otel.Tracer("test")
	repos, _ := repo.NewConfigRepo(logger, tracer)
	services := service.NewConfigService(repos, logger, tracer)
	return handler.NewConfigHandler(services, logger, tracer)
}

func TestGetConfig(t *testing.T) {
	configHandler := newTestHandler() // Zamijeni sa stvarnim handlerom

	// Simuliramo HTTP zahtev
	req, err := http.NewRequest("GET", "/configs/db_config/2", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Kreiramo testni ResponseWriter da uhvatimo odgovor
	rr := httptest.NewRecorder()

	// Kreiramo ruter koji će obrađivati naš zahtev (ovde pretpostavljamo da imamo router kreiran u ConfigHandler)
	router := mux.NewRouter()
	router.HandleFunc("/configs/{name}/{version}", configHandler.Get).Methods("GET")

	// Serviramo zahtev na naš ruter
	router.ServeHTTP(rr, req)

	// Proveravamo status kod
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status %v, but got %v", http.StatusOK, status)
	}
}

func TestDeleteConfig(t *testing.T) {

	configHandler := newTestHandler()
	// Simulacija zahteva za brisanje konfiguracije
	req, err := http.NewRequest("DELETE", "/configs/config_name/config_version", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Testni ResponseWriter za hvatanje odgovora
	rr := httptest.NewRecorder()

	// Kreiramo ruter sa odgovarajućim hendlerom
	router := mux.NewRouter()
	router.HandleFunc("/configs/{name}/{version}", configHandler.Delete).Methods("DELETE")

	// Serviramo zahtev na ruter
	router.ServeHTTP(rr, req)

	// Proveravamo status kod
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status %v, but got %v", http.StatusOK, status)
	}
}

func TestAddConfig(t *testing.T) {

	configHandler := newTestHandler()
	// Simulacija zahteva za dodavanje nove konfiguracije
	jsonBody := []byte(`{"name": "config_name", "version": "config_version", "data": {...}}`)
	req, err := http.NewRequest("POST", "/configs/", bytes.NewBuffer(jsonBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Testni ResponseWriter za hvatanje odgovora
	rr := httptest.NewRecorder()

	// Kreiramo ruter sa odgovarajućim hendlerom
	router := mux.NewRouter()
	router.HandleFunc("/configs/", configHandler.Add).Methods("POST")

	// Serviramo zahtev na ruter
	router.ServeHTTP(rr, req)

	// Proveravamo status kod
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status %v, but got %v", http.StatusOK, status)
	}
}

func TestDeleteAllConfigs(t *testing.T) {

	configHandler := newTestHandler()
	// Simulacija zahteva za brisanje svih konfiguracija
	req, err := http.NewRequest("DELETE", "/configs/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Testni ResponseWriter za hvatanje odgovora
	rr := httptest.NewRecorder()

	// Kreiramo ruter sa odgovarajućim hendlerom
	router := mux.NewRouter()
	router.HandleFunc("/configs/", configHandler.DeleteAll).Methods("DELETE")

	// Serviramo zahtev na ruter
	router.ServeHTTP(rr, req)

	// Proveravamo status kod
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status %v, but got %v", http.StatusOK, status)
	}
}
