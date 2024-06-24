package handler

import (
	"alati/model"
	"alati/service"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"

	"github.com/gorilla/mux"
)

type ConfigHandler struct {
	service service.ConfigService
	logger  *log.Logger
}

func NewConfigHandler(service service.ConfigService, logger *log.Logger) ConfigHandler {
	return ConfigHandler{
		service: service,
		logger:  logger,
	}
}

func decodeBody(r io.Reader) (*model.Config, error) {
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()

	var rt model.Config
	if err := dec.Decode(&rt); err != nil {
		return nil, err
	}
	return &rt, nil
}

func renderJSON(w http.ResponseWriter, v interface{}) {
	js, err := json.Marshal(v)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

// @Summary Get a configuration
// @Description Retrieves a configuration by name and version
// @Tags configs
// @Produce json
// @Param name path string true "Name of the configuration"
// @Param version path int true "Version of the configuration"
// @Success 200 {object} model.Config
// @Failure 400 {string} string "Invalid input"
// @Failure 404 {string} string "Configuration not found"
// @Failure 500 {string} string "Internal server error"
// @Router /configs/{name}/{version} [get]
func (c ConfigHandler) Get(w http.ResponseWriter, r *http.Request) {
	//time.Sleep(9 * time.Second)
	name := mux.Vars(r)["name"]
	version := mux.Vars(r)["version"]

	i := "config/%s/%s"
	id := fmt.Sprintf(i, name, version)

	config, err := c.service.Get(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	resp, err := json.Marshal(config)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content−Type", "application/json")
	w.Write(resp)
}

// @Summary Get all configurations
// @Description Retrieves all configurations
// @Tags configs
// @Produce json
// @Success 200 {array} model.Config
// @Failure 500 {string} string "Internal server error"
// @Router /configs/ [get]
func (c *ConfigHandler) GetAll(rw http.ResponseWriter, h *http.Request) {
	allProducts, err := c.service.GetAll()

	if err != nil {
		http.Error(rw, "Database exception", http.StatusInternalServerError)
		c.logger.Fatal("Database exception: ", err)
	}

	renderJSON(rw, allProducts)

}

// @Summary Add a new configuration
// @Description Adds a new configuration
// @Tags configs
// @Accept json
// @Produce json
// @Param config body model.Config true "Configuration to add"
// @Success 201 {object} model.Config
// @Failure 400 {string} string "Invalid input"
// @Failure 403 {string} string "Idempotentcy protection"
// @Failure 415 {string} string "Unsupported media type"
// @Failure 500 {string} string "Internal server error"
// @Router /configs/ [post]
func (c ConfigHandler) Add(w http.ResponseWriter, req *http.Request) {
	contentType := req.Header.Get("Content-Type")
	idempotency_key := req.Header.Get("idempotency-key")
	mediatype, _, err := mime.ParseMediaType(contentType)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if mediatype != "application/json" {
		err := errors.New("expect application/json Content-Type")
		http.Error(w, err.Error(), http.StatusUnsupportedMediaType)
		return
	}

	rt, err := decodeBody(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	config, err := c.service.Add(rt, idempotency_key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if config == nil && err == nil {
		http.Error(w, "Idempotency protection", http.StatusForbidden)
		return
	}

	renderJSON(w, config)
}

// @Summary Delete a configuration
// @Description Deletes a configuration by name and version
// @Tags configs
// @Produce json
// @Param name path string true "Name of the configuration"
// @Param version path int true "Version of the configuration"
// @Success 200 {string} string "Deleted"
// @Failure 500 {string} string "Internal server error"
// @Router /configs/{name}/{version} [delete]
func (c ConfigHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	version := vars["version"]

	i := "config/%s/%s"
	id := fmt.Sprintf(i, name, version)

	err := c.service.Delete(id)
	if err != nil {
		http.Error(w, "Failed to delete config", http.StatusInternalServerError)
		return
	}

	renderJSON(w, "Deleted")
}

func (c ConfigHandler) DeleteAll(rw http.ResponseWriter, h *http.Request) {

	err := c.service.DeleteAll()
	if err != nil {
		http.Error(rw, "Database exception", http.StatusInternalServerError)
		c.logger.Fatal("Database exception:", err)
		return
	}

	rw.WriteHeader(http.StatusOK)
}
