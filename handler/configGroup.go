package handler

import (
	"alati/model"
	"alati/service"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"mime"
	"net/http"
	"strconv"
)

type ConfigGroupHandler struct {
	service service.ConfigGroupService
}

func NewConfigGruopHandler(service service.ConfigGroupService) ConfigGroupHandler {
	return ConfigGroupHandler{
		service: service,
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

func (c ConfigGroupHandler) Get(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]
	version := mux.Vars(r)["version"]

	versionInt, err := strconv.Atoi(version)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	config, err := c.service.Get(name, versionInt)

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

func (c ConfigGroupHandler) Add(w http.ResponseWriter, req *http.Request) {
	contentType := req.Header.Get("Content-Type")
	mediatype, _, err := mime.ParseMediaType(contentType)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if mediatype != " application/json" {
		err := errors.New("expect application/json Content-Type")
		http.Error(w, err.Error(), http.StatusUnsupportedMediaType)
		return
	}

	rt, err := decodeBodyCG(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	con := model.ConfigGroup{
		Name:    rt.Name,
		Version: rt.Version,
		Configs: rt.Configs,
	}
	c.service.Add(con)

	renderJSON(w, rt)
}

func (c ConfigGroupHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	versionStr := vars["version"]

	version, err := strconv.Atoi(versionStr)
	if err != nil {
		http.Error(w, "Invalid version", http.StatusBadRequest)
		return
	}

	err = c.service.Delete(name, version)
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

	versionG, err := strconv.Atoi(versionGStr)
	if err != nil {
		http.Error(w, "Invalid version", http.StatusBadRequest)
		return
	}

	versionC, err := strconv.Atoi(versionCStr)
	if err != nil {
		http.Error(w, "Invalid version", http.StatusBadRequest)
		return
	}

	key := fmt.Sprintf("%s/%d", nameC, versionC)

	group, _ := c.service.Get(nameG, versionG)
	conf, _ := service.ConfigService{}.Get(nameC, versionC)
	group.Configs[key] = &conf

	renderJSON(w, "success Put")
}

func (c ConfigGroupHandler) RemoveConfFromGroup(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	nameG := vars["nameG"]
	versionGStr := vars["versionG"]
	nameC := vars["nameC"]
	versionCStr := vars["versionC"]

	versionG, err := strconv.Atoi(versionGStr)
	if err != nil {
		http.Error(w, "Invalid version", http.StatusBadRequest)
		return
	}

	versionC, err := strconv.Atoi(versionCStr)
	if err != nil {
		http.Error(w, "Invalid version", http.StatusBadRequest)
		return
	}

	key := fmt.Sprintf("%s/%d", nameC, versionC)

	group, _ := c.service.Get(nameG, versionG)
	delete(group.Configs, key)

	renderJSON(w, "success Put")
}
