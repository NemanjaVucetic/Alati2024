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

// Get retrieves a configuration by name and version
// swagger:route GET /configs/{name}/{version} Config getConfig
//
// Retrieves a configuration by name and version.
//
// Responses:
//
//		200: Config
//		400: BadRequest
//		404: NotFound
//	 500: InternalServerError
func (c ConfigHandler) Get(w http.ResponseWriter, r *http.Request) {
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

	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}

// Add creates a new configuration
// swagger:route POST /configs Config addConfig
//
// Creates a new configuration.
//
// Consumes:
// - application/json
//
// Responses:
//
//	201: Config
//	400: BadRequest
func (c ConfigHandler) Add(w http.ResponseWriter, req *http.Request) {
	contentType := req.Header.Get("Content-Type")
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

	config, err := c.service.Add(rt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	renderJSON(w, config)
}

// Delete removes a configuration by name and version
// swagger:route DELETE /configs/{name}/{version} Config deleteConfig
//
// Removes a configuration by name and version.
//
// Responses:
//
//	200: string
//	400: BadRequest
//	500: InternalServerError
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
