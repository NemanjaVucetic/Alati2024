package main

import (
	"alati/model"
	repo "alati/repo"
	"fmt"
)
import service "alati/service"

func main() {
	repo := repo.NewConfigInMemRepository()
	service := service.NewConfigService(repo)
	service.Hello()

	mapa := make(map[string]string)

	mapa["kljuc1"] = "cao"
	mapa["kljuc2"] = "cao"
	mapa["kljuc3"] = "cao"

	c := model.NewConfig("cao", 2.12, mapa)
	conf := *c
	//fmt.Println(c.GenerateKey())

	repo.Add(conf)
	repo.Delete(conf.GenerateKey())
	fmt.Println(repo.Get("cao2.12"))

}
