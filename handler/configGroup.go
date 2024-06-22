package handler

import (
	"alati/model"
	"alati/service"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"log"
	"mime"
	"net/http"
	"strings"
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

// @Summary Get a configuration group
// @Description Retrieves a configuration group by name and version
// @Tags configGroups
// @Produce json
// @Param name path string true "Name of the configuration group"
// @Param version path int true "Version of the configuration group"
// @Success 200 {object} model.ConfigGroup
// @Failure 400 {string} string "Invalid input"
// @Failure 404 {string} string "Configuration group not found"
// @Failure 500 {string} string "Internal server error"
// @Router /configGroups/{name}/{version} [get]
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

	w.Header().Set("Content−Type", "application/json")
	w.Write(resp)
}

// @Summary Get all configuration groups
// @Description Retrieves all configuration groups
// @Tags configGroups
// @Produce json
// @Success 200 {array} model.ConfigGroup
// @Failure 500 {string} string "Internal server error"
// @Router /configGroups/ [get]
func (c *ConfigGroupHandler) GetAll(rw http.ResponseWriter, h *http.Request) {
	allProducts, err := c.service.GetAll()

	if err != nil {
		http.Error(rw, "Database exception", http.StatusInternalServerError)
		c.logger.Fatal("Database exception: ", err)
	}

	renderJSON(rw, allProducts)

}

// @Summary Add a new configuration group
// @Description Adds a new configuration group
// @Tags configGroups
// @Accept json
// @Produce json
// @Param configGroup body model.ConfigGroup true "Configuration group to add"
// @Success 201 {object} model.ConfigGroup
// @Failure 400 {string} string "Invalid input"
// @Failure 415 {string} string "Unsupported media type"
// @Failure 500 {string} string "Internal server error"
// @Router /configGroups/ [post]
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

// @Summary Delete a configuration group
// @Description Deletes a configuration group by name and version
// @Tags configGroups
// @Produce json
// @Param name path string true "Name of the configuration group"
// @Param version path int true "Version of the configuration group"
// @Success 200 {string} string "Deleted"
// @Failure 500 {string} string "Internal server error"
// @Router /configGroups/{name}/{version} [delete]
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

// @Summary Add a configuration to a group
// @Description Adds a configuration to a group by their names and versions
// @Tags configGroups
// @Produce json
// @Param nameG path string true "Name of the configuration group"
// @Param versionG path int true "Version of the configuration group"
// @Param nameC path string true "Name of the configuration"
// @Param versionC path int true "Version of the configuration"
// @Success 200 {string} string "success Put"
// @Failure 500 {string} string "Internal server error"
// @Router /configGroups/{nameG}/{versionG}/configs/{nameC}/{versionC} [put]
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
		return
	}

	renderJSON(w, "success Put")
}

// @Summary Remove a configuration from a group
// @Description Removes a configuration from a group by their names and versions
// @Tags configGroups
// @Produce json
// @Param nameG path string true "Name of the configuration group"
// @Param versionG path int true "Version of the configuration group"
// @Param nameC path string true "Name of the configuration"
// @Param versionC path int true "Version of the configuration"
// @Success 200 {string} string "success Put"
// @Failure 500 {string} string "Internal server error"
// @Router /configGroups/{nameG}/{versionG}/configs/{nameC}/{versionC} [put]
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
		return
	}

	renderJSON(w, "success Put")
}

// @Summary Get configurations by labels
// @Description Retrieves configurations by labels within a group
// @Tags configGroups
// @Produce json
// @Param nameG path string true "Name of the configuration group"
// @Param versionG path int true "Version of the configuration group"
// @Param labels path string true "Labels of the configuration"
// @Param nameC path string false "Name of the configuration"
// @Param versionC path int false "Version of the configuration"
// @Success 200 {object} []model.Config
// @Failure 400 {string} string "Invalid input"
// @Failure 500 {string} string "Internal server error"
// @Router /configGroups/{nameG}/{versionG}/configs/{labels}/{nameC}/{versionC} [get]
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

// @Summary Delete configurations by labels
// @Description Deletes configurations by labels within a group
// @Tags configGroups
// @Produce json
// @Param nameG path string true "Name of the configuration group"
// @Param versionG path int true "Version of the configuration group"
// @Param labels path string true "Labels of the configuration"
// @Param nameC path string false "Name of the configuration"
// @Param versionC path int false "Version of the configuration"
// @Success 200 {string} string "deleted"
// @Failure 400 {string} string "Invalid input"
// @Failure 500 {string} string "Internal server error"
// @Router /configGroups/{nameG}/{versionG}/configs/{labels}/{nameC}/{versionC} [patch]
func (c ConfigGroupHandler) DeleteConfigsByLabels(w http.ResponseWriter, r *http.Request) {
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

	err := c.service.DeleteConfigsByLabels(prefixGroup, prefixConf)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	renderJSON(w, "deleted")
}
