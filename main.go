package main

import (
	"alati/handler"
	"alati/model"
	"alati/repo"
	"alati/service"
	"context"
	"github.com/gorilla/mux"
	"golang.org/x/time/rate"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	//test {
	repoC := repo.NewConfigInMemRepository()
	serviceC := service.NewConfigService(repoC)
	repoG := repo.NewConfigGroupInMemRepository()
	serviceG := service.NewConfigGroupService(repoG)
	h := handler.NewConfigHandler(serviceC)
	hG := handler.NewConfigGroupHandler(serviceG, serviceC)

	params := make(map[string]string)
	params["username"] = "pera"
	params["port"] = "5432"

	config := model.Config{
		Name:    "db_config",
		Version: 2,
		Params:  params,
	}
	config2 := model.Config{
		Name:    "db_config2",
		Version: 3,
		Params:  params,
	}

	configMap := make(map[string]model.Config)
	configMap["conf1"] = config
	configMap["conf2"] = config2

	group := model.ConfigGroup{
		Name:    "db_cg",
		Version: 2,
		Configs: configMap,
	}

	serviceC.Add(config2)
	serviceC.Add(config)
	serviceG.Add(group)

	// }

	router := mux.NewRouter()
	limiter := rate.NewLimiter(0.167, 10)

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	router.Handle("/configs/{name}/{version}", handler.RateLimit(limiter, h.Get)).Methods(http.MethodGet)
	router.Handle("/configGroups/{name}/{version}", handler.RateLimit(limiter, hG.Get)).Methods(http.MethodGet)
	router.Handle("/configs/", handler.RateLimit(limiter, h.Add)).Methods(http.MethodPost)
	router.Handle("/configGroups/", handler.RateLimit(limiter, hG.Add)).Methods(http.MethodPost)
	router.Handle("/configs/{name}/{version}", handler.RateLimit(limiter, h.Delete)).Methods(http.MethodDelete)
	router.Handle("/configGroups/{name}/{version}", handler.RateLimit(limiter, hG.Delete)).Methods(http.MethodDelete)
	router.Handle("/configGroups/{nameG}/{versionG}/config/{nameC}/{versionC}", handler.RateLimit(limiter, hG.AddConfToGroup)).Methods(http.MethodPut)
	router.Handle("/configGroups/{nameG}/{versionG}/{nameC}/{versionC}", handler.RateLimit(limiter, hG.RemoveConfFromGroup)).Methods(http.MethodPut)

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
