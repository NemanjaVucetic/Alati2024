package main

import (
	"alati/handler"
	"alati/model"
	"alati/repo"
	"alati/service"
	"github.com/gorilla/mux"
	"net/http"
)

func main() {
	repo := repo.NewConfigInMemRepository()
	service := service.NewConfigService(repo)

	params := make(map[string]string)
	params["username"] = "pera"
	params["port"] = "5432"
	config := model.Config{
		Name:    "db_config",
		Version: 2,
		Params:  params,
	}
	service.Add(config)
	h := handler.NewConfigHandler(service)

	router := mux.NewRouter()

	router.HandleFunc("/configs/{name}/{version}", h.Get).Methods("GET")

	http.ListenAndServe("0.0.0.0:8000", router)
}
