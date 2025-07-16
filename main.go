package main

import (
	"log"

	router "github.com/bncunha/erp-api/src/api"
	controller "github.com/bncunha/erp-api/src/api/controllers"
	"github.com/bncunha/erp-api/src/application/service"
	"github.com/bncunha/erp-api/src/infrastructure/persistence"
	"github.com/bncunha/erp-api/src/infrastructure/repository"
	config "github.com/bncunha/erp-api/src/main"
)

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Erro buscar variaveis de ambiente", err)
	}


	persistence := persistence.NewPersistence(config)
	db, err := persistence.ConnectDb()
	if err != nil {
		log.Fatal("Erro ao conectar no banco de dados", err)
	}
	defer persistence.CloseConnection(db)

	repository := repository.NewRepository(db)
	repository.SetupRepositories()

	service := service.NewApplicationService(repository)
	service.SetupServices()
	
	controller := controller.NewController(service)
	controller.SetupControllers()

	r := router.NewRouter(controller)
	r.SetupRoutes()
	r.Start()
}