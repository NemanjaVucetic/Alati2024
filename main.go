// Alati2024
//
//	Title: Alati2024
//
//	Schemes: http
//	Version: 0.0.1
//	BasePath: /
//
//	Produces:
//	  - application/json
//
// swagger:meta
package main

import (
	"alati/handler"
	"alati/model"
	"alati/repo"
	"alati/service"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gorilla/mux"
	"golang.org/x/time/rate"
)

func main() {

	//test {

	logger := log.New(os.Stdout, "[config-api] ", log.LstdFlags)

	repoC, err := repo.NewConfigRepo(logger)
	if err != nil {
		logger.Fatal(err)
	}

	serviceC := service.NewConfigService(repoC, logger)

	repoG, err := repo.NewConfigGroupRepo(logger)
	if err != nil {
		logger.Fatal(err)
	}

	serviceG := service.NewConfigGroupService(repoG, logger)

	h := handler.NewConfigHandler(serviceC, logger)
	hG := handler.NewConfigGroupHandler(serviceG, serviceC, logger)

	params := make(map[string]string)
	params["param1"] = "param1"
	params["param2"] = "param2"

	labels := make(map[string]string)
	labels["l1"] = "v1"
	labels["l2"] = "v2"

	labels2 := make(map[string]string)
	labels2["l1"] = "v1"

	config := model.Config{
		Name:    "db_config",
		Version: 2,
		Params:  params,
		Labels:  labels,
	}
	config2 := model.Config{
		Name:    "db_config2",
		Version: 3,
		Params:  params,
		Labels:  labels2,
	}

	configMap := make(map[string]model.Config)
	configMap["configGroups/db_cg/2/config/l1:v1/l2:v2/db_config/2"] = config
	configMap["configGroups/db_cg/2/config/l1:v1/db_config2/3"] = config2

	group := model.ConfigGroup{
		Name:    "db_cg",
		Version: 2,
		Configs: configMap,
	}

	serviceC.Add(&config2)
	serviceC.Add(&config)
	serviceG.Add(&group)

	// }

	router := mux.NewRouter()
	limiter := rate.NewLimiter(0.167, 10)

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	router.Handle("/configs/{name}/{version}", handler.RateLimit(limiter, h.Get)).Methods(http.MethodGet)
	router.Handle("/configs/", handler.RateLimit(limiter, h.Add)).Methods(http.MethodPost)
	router.Handle("/configs/{name}/{version}", handler.RateLimit(limiter, h.Delete)).Methods(http.MethodDelete)
	router.Handle("/configs/", handler.RateLimit(limiter, h.DeleteAll)).Methods(http.MethodDelete)

	router.Handle("/configGroups/{name}/{version}", handler.RateLimit(limiter, hG.Get)).Methods(http.MethodGet)
	router.Handle("/configGroups/", handler.RateLimit(limiter, hG.Add)).Methods(http.MethodPost)
	router.Handle("/configGroups/{name}/{version}", handler.RateLimit(limiter, hG.Delete)).Methods(http.MethodDelete)
	router.Handle("/configGroups/{nameG}/{versionG}/configs/{nameC}/{versionC}", handler.RateLimit(limiter, hG.AddConfToGroup)).Methods(http.MethodPut)
	router.Handle("/configGroups/{nameG}/{versionG}/{nameC}/{versionC}", handler.RateLimit(limiter, hG.RemoveConfFromGroup)).Methods(http.MethodPut)
	router.Handle("/configGroups/{name}/{version}", handler.RateLimit(limiter, hG.GetConfigsByLabels)).Methods(http.MethodPost)
	router.Handle("/configGroups/{name}/{version}", handler.RateLimit(limiter, hG.DeleteConfigsByLabels)).Methods(http.MethodPut)

	// SwaggerUI
	optionsDevelopers := middleware.SwaggerUIOpts{SpecURL: "swagger.yaml"}
	developerDocumentationHandler := middleware.SwaggerUI(optionsDevelopers, nil)
	router.Handle("/docs", developerDocumentationHandler)

	srv := &http.Server{Addr: "0.0.0.0:8000", Handler: router}

	go func() {
		log.Println("server_starting")
		if err := srv.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				log.Fatal(err)
			}
		}
	}()

	<-quit

	log.Println("service_shutting_down")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
	log.Println(" server stopped")
}
