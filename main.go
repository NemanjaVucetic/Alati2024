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

	//test {
	repoC := repo.NewConfigInMemRepository()
	serviceC := service.NewConfigService(repoC)
	repoG := repo.NewConfigGroupInMemRepository()
	serviceG := service.NewConfigGroupService(repoG)

	params := make(map[string]string)
	params["username"] = "pera"
	params["port"] = "5432"
	config := model.Config{
		Name:    "db_config",
		Version: 2,
		Params:  params,
	}
	serviceC.Add(config)
	h := handler.NewConfigHandler(serviceC)

	hG := handler.NewConfigGruopHandler(serviceG)
	// }

	router := mux.NewRouter()

	router.HandleFunc("/configs/{name}/{version}", h.Get).Methods("GET")
	router.HandleFunc("/configGroups/{name}/{version}", hG.Get).Methods("GET")
	router.HandleFunc("/configs/", h.Add).Methods("POST")
	router.HandleFunc("/configGroups/", hG.Add).Methods("POST")
	router.HandleFunc("/configs/{name}/{version}", h.Delete).Methods("DELETE")
	router.HandleFunc("/configGroups/{name}/{version}", hG.Delete).Methods("DELETE")
	router.HandleFunc("/configGroups/{nameG}/{versionG}/config/{nameC}/{versionC}", hG.AddConfToGroup).Methods("PUT")
	router.HandleFunc("/configGroups/{nameG}/{versionG}/{nameC}/{versionC}", hG.RemoveConfFromGroup).Methods("PUT")

	http.ListenAndServe("0.0.0.0:8000", router)
}
