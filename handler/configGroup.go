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
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

type ConfigGroupHandler struct {
	service       service.ConfigGroupService
	serviceConfig service.ConfigService
	logger        *log.Logger
}

func NewConfigGroupHandler(service service.ConfigGroupService, serviceConfig service.ConfigService, logger *log.Logger) ConfigGroupHandler {
	return ConfigGroupHandler{
		service:       service,
		serviceConfig: serviceConfig,
		logger:        logger,
	}
}

func decodeBodyCG(r io.Reader) (*model.ConfigGroup, error) {
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()

	var rt model.ConfigGroup
	if err := dec.Decode(&rt); err != nil {
		return nil, err
	}
	return &rt, nil
}

func decodeBodyLabels(r io.Reader) (*map[string]string, error) {
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()

	var rt map[string]string
	if err := dec.Decode(&rt); err != nil {
		return nil, err
	}
	return &rt, nil
}

// Get retrieves a configuration group by name and version
// swagger:route GET /config-groups/{name}/{version} ConfigGroup getConfigGroup
//
// Retrieves a configuration group by name and version.
//
// Responses:
//
//	200: ConfigGroup
//	400: BadRequest
//	404: NotFound
func (c ConfigGroupHandler) Get(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]
	version := mux.Vars(r)["version"]

	i := "configGroups/%s/%s"
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

// Add creates a new configuration group
// swagger:route POST /config-groups ConfigGroup addConfigGroup
//
// Creates a new configuration group.
//
// Consumes:
// - application/json
//
// Responses:
//
//	201: ConfigGroup
//	400: BadRequest
func (c ConfigGroupHandler) Add(w http.ResponseWriter, req *http.Request) {
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

	rt, err := decodeBodyCG(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	c.service.Add(rt)

	renderJSON(w, rt)
}

// Delete removes a configuration group by name and version
// swagger:route DELETE /config-groups/{name}/{version} ConfigGroup deleteConfigGroup
//
// Removes a configuration group by name and version.
//
// Responses:
//
//	200: string
//	400: BadRequest
//	500: InternalServerError
func (c ConfigGroupHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	version := vars["version"]

	i := "configGroups/%s/%s"
	id := fmt.Sprintf(i, name, version)

	err := c.service.Delete(id)
	if err != nil {
		http.Error(w, "Failed to delete config group", http.StatusInternalServerError)
		return
	}

	renderJSON(w, "Deleted")
}

// AddConfToGroup adds a configuration to a configuration group
// swagger:route POST /config-groups/{nameG}/{versionG}/configs/{nameC}/{versionC} ConfigGroup addConfigToGroup
//
// Adds a configuration to a configuration group.
//
// Responses:
//
//	200: string
//	400: BadRequest
func (c ConfigGroupHandler) AddConfToGroup(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	nameG := vars["nameG"]
	versionGStr := vars["versionG"]
	nameC := vars["nameC"]
	versionCStr := vars["versionC"]

	groupString := "configGroups/%s/%s"
	confString := "config/%s/%s"

	gStr := fmt.Sprintf(groupString, nameG, versionGStr)
	cStr := fmt.Sprintf(confString, nameC, versionCStr)

	group, _ := c.service.Get(gStr)
	conf, _ := c.serviceConfig.Get(cStr)

	err := c.service.AddConfigToGroup(*group, *conf)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	renderJSON(w, "success Put")
}

// RemoveConfFromGroup removes a configuration from a configuration group by key
// swagger:route DELETE /config-groups/{nameG}/{versionG}/configs/{nameC}/{versionC} ConfigGroup removeConfigFromGroup
//
// Removes a configuration from a configuration group by key.
//
// Responses:
//
//	200: string
//	400: BadRequest

func (c ConfigGroupHandler) RemoveConfFromGroup(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	nameG := vars["nameG"]
	versionGStr := vars["versionG"]
	nameC := vars["nameC"]
	versionCStr := vars["versionC"]

	i := "configGroups/%s/%s"
	id := fmt.Sprintf(i, nameG, versionGStr)

	t := "config/%s/%s"
	idc := fmt.Sprintf(t, nameC, versionCStr)

	config, _ := c.serviceConfig.Get(idc)
	group, _ := c.service.Get(id)

	err := c.service.RemoveConfigFromGroup(*group, *config)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	renderJSON(w, "success Put")
}

// GetConfigsByLabels retrieves configurations from a configuration group by labels
// swagger:route POST /config-groups/{name}/{version}/labels ConfigGroup getConfigsByLabels
//
// Retrieves configurations from a configuration group by labels.
//
// Responses:
//
//	200: []Config
//	400: BadRequest
func (c ConfigGroupHandler) GetConfigsByLabels(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	nameG := vars["nameG"]
	versionG := vars["versionG"]
	original := vars["labels"]
	nameC := vars["nameC"]
	versionC := vars["versionC"]
	labels := strings.ReplaceAll(original, ";", "/")

	prefixGroup := fmt.Sprintf("configGroups/%s/%s", nameG, versionG)
	prefixConf := fmt.Sprintf("configGroups/%s/%s/config/%s", nameG, versionG, labels)

	if labels == "" {
		prefixConf = prefixGroup
	}
	if nameC != "" {
		prefixConf = prefixConf + "/" + nameC
	}
	if versionC != "" {
		prefixConf = prefixConf + "/" + versionC
	}

	conf, err := c.service.GetConfigsByLabels(prefixGroup, prefixConf)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	renderJSON(w, conf)
}

// DeleteConfigsByLabels deletes configurations from a configuration group by labels
// swagger:route DELETE /config-groups/{name}/{version}/labels ConfigGroup deleteConfigsByLabels
//
// Deletes configurations from a configuration group by labels.
//
// Responses:
//
//	200: string
//	400: BadRequest
func (c ConfigGroupHandler) DeleteConfigsByLabels(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]
	version := mux.Vars(r)["version"]

	versionInt, err := strconv.Atoi(version)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	group, err := c.service.Get(name, versionInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	labels, err := decodeBodyLabels(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := c.service.DeleteConfigsByLabels(prefixGroup, prefixConf)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	renderJSON(w, "deleted")
}
