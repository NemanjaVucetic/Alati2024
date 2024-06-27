package handler

import (
	"alati/model"
	"alati/service"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"log"
	"mime"
	"net/http"
	"strings"


	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

)

type ConfigGroupHandler struct {
	service       service.ConfigGroupService
	serviceConfig service.ConfigService
	logger        *log.Logger
	Tracer        trace.Tracer
}

func NewConfigGroupHandler(service service.ConfigGroupService, serviceConfig service.ConfigService,
	logger *log.Logger, Tracer trace.Tracer) ConfigGroupHandler {

	return ConfigGroupHandler{
		service:       service,
		serviceConfig: serviceConfig,
		logger:        logger,
		Tracer:        Tracer,
	}
}

func (c ConfigGroupHandler) decodeBodyCG(r io.Reader, ctx context.Context) (*model.ConfigGroup, context.Context, error) {
	cont, span := c.Tracer.Start(ctx, "decodeBody")
	defer span.End()
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()

	var rt model.ConfigGroup
	if err := dec.Decode(&rt); err != nil {
		return nil, cont, err
	}
	return &rt, cont, nil
}

func (c *ConfigGroupHandler) renderJSON(w http.ResponseWriter, v interface{}, ctx context.Context) {
	_, span := c.Tracer.Start(ctx, "renderJSON")
	defer span.End()
	js, err := json.Marshal(v)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
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
	ctx, span := c.Tracer.Start(r.Context(), "h.GetGroup")
	defer span.End()
	name := mux.Vars(r)["name"]
	version := mux.Vars(r)["version"]

	i := "configGroups/%s/%s"
	id := fmt.Sprintf(i, name, version)

	config, err := c.service.Get(id, ctx)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	resp, err := json.Marshal(config)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
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
func (c *ConfigGroupHandler) GetAll(rw http.ResponseWriter, r *http.Request) {
	ctx, span := c.Tracer.Start(r.Context(), "h.GetAllGroups")
	defer span.End()
	allProducts, err := c.service.GetAll(ctx)

	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		http.Error(rw, "Database exception", http.StatusInternalServerError)
		c.logger.Fatal("Database exception: ", err)
	}

	c.renderJSON(rw, allProducts, ctx)

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
	ctx, span := c.Tracer.Start(req.Context(), "h.AddGroup")
	defer span.End()
	contentType := req.Header.Get("Content-Type")
	mediatype, _, err := mime.ParseMediaType(contentType)

	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if mediatype != "application/json" {
		span.SetStatus(codes.Error, err.Error())
		err := errors.New("expect application/json Content-Type")
		http.Error(w, err.Error(), http.StatusUnsupportedMediaType)
		return
	}

	rt, cont, err := c.decodeBodyCG(req.Body, ctx)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}


	group, err := c.service.Add(rt, idempotency_key, cont)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if group == nil && err == nil {
		span.SetStatus(codes.Error, err.Error())
		http.Error(w, "Idempotency protection", http.StatusForbidden)
		return
	}

	c.renderJSON(w, group, cont)
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
	ctx, span := c.Tracer.Start(r.Context(), "h.DeleteGroup")
	defer span.End()
	vars := mux.Vars(r)
	name := vars["name"]
	version := vars["version"]

	i := "configGroups/%s/%s"
	id := fmt.Sprintf(i, name, version)

	err := c.service.Delete(id, ctx)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		http.Error(w, "Failed to delete config group", http.StatusInternalServerError)
		return
	}

	c.renderJSON(w, "Deleted", ctx)
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
	ctx, span := c.Tracer.Start(r.Context(), "h.AddConfToGroup")
	defer span.End()
	vars := mux.Vars(r)
	nameG := vars["nameG"]
	versionGStr := vars["versionG"]
	nameC := vars["nameC"]
	versionCStr := vars["versionC"]

	groupString := "configGroups/%s/%s"
	confString := "config/%s/%s"

	gStr := fmt.Sprintf(groupString, nameG, versionGStr)
	cStr := fmt.Sprintf(confString, nameC, versionCStr)

	group, _ := c.service.Get(gStr, ctx)
	conf, _ := c.serviceConfig.Get(cStr, ctx)


	err := c.service.AddConfigToGroup(*group, *conf, idempotency_key, ctx)

	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return
	}

	c.renderJSON(w, "success Put", ctx)
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
	ctx, span := c.Tracer.Start(r.Context(), "h.RemoveConfFromGroup")
	defer span.End()
	vars := mux.Vars(r)
	nameG := vars["nameG"]
	versionGStr := vars["versionG"]
	nameC := vars["nameC"]
	versionCStr := vars["versionC"]

	i := "configGroups/%s/%s"
	id := fmt.Sprintf(i, nameG, versionGStr)

	t := "config/%s/%s"
	idc := fmt.Sprintf(t, nameC, versionCStr)

	config, _ := c.serviceConfig.Get(idc, ctx)
	group, _ := c.service.Get(id, ctx)


	err := c.service.RemoveConfigFromGroup(*group, *config, idempotency_key, ctx)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return
	}

	c.renderJSON(w, "success Put", ctx)
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
	ctx, span := c.Tracer.Start(r.Context(), "h.GetConfigsByLabels")
	defer span.End()
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

	conf, err := c.service.GetConfigsByLabels(prefixGroup, prefixConf, ctx)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	c.renderJSON(w, conf, ctx)
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
	ctx, span := c.Tracer.Start(r.Context(), "h.DeleteConfigsByLabels")
	defer span.End()
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

	err := c.service.DeleteConfigsByLabels(prefixGroup, prefixConf, ctx)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	c.renderJSON(w, "deleted", ctx)
}
