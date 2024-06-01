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

func (c *ConfigGroupHandler) GetAll(rw http.ResponseWriter, h *http.Request) {
	allProducts, err := c.service.GetAll()

	if err != nil {
		http.Error(rw, "Database exception", http.StatusInternalServerError)
		c.logger.Fatal("Database exception: ", err)
	}

	renderJSON(rw, allProducts)

}

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
