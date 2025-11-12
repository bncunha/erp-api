package main

import (
	router "github.com/bncunha/erp-api/src/api"
	controller "github.com/bncunha/erp-api/src/api/controllers"
	"github.com/bncunha/erp-api/src/application/ports"
	"github.com/bncunha/erp-api/src/application/service"
	"github.com/bncunha/erp-api/src/application/usecase"
	"github.com/bncunha/erp-api/src/infrastructure/bcrypt"
	"github.com/bncunha/erp-api/src/infrastructure/logs"
	"github.com/bncunha/erp-api/src/infrastructure/observability"
	"github.com/bncunha/erp-api/src/infrastructure/persistence"
	"github.com/bncunha/erp-api/src/infrastructure/repository"
	config "github.com/bncunha/erp-api/src/main"
	"github.com/newrelic/go-agent/v3/newrelic"
)

func main() {
	bcrypt := bcrypt.NewBcrypt()
	logs.NewLogs()
	logs.Logger.Infof("Iniciando aplicação")

	config, err := config.LoadConfig()
	if err != nil {
		logs.Logger.Fatalf("Erro buscar variaveis de ambiente", err)
	}

	observability := observability.NewObservability(observability.NewNewRelicObservability())
	nrl := observability.GetApp().(*newrelic.Application)
	logs.Logger.AddHook(logs.NewRelicHook(nrl))
	err = observability.SetupObservability(config)
	if err != nil {
		logs.Logger.Fatalf("Erro ao configurar observabilidade", err)
	}

	persistence := persistence.NewPersistence(config)
	db, err := persistence.ConnectDb()
	if err != nil {
		logs.Logger.Fatalf("Erro ao conectar no banco de dados", err)
	}
	defer persistence.CloseConnection(db)

	repository := repository.NewRepository(db)
	repository.SetupRepositories()

	useCase := usecase.NewApplicationUseCase(repository)
	useCase.SetupUseCases()

	ports := ports.NewPorts(bcrypt)

	service := service.NewApplicationService(repository, useCase, ports)
	service.SetupServices()

	controller := controller.NewController(service)
	controller.SetupControllers()

	r := router.NewRouter(controller)
	r.SetupCors(config.APP_ENV)
	r.SetupRoutes()
	r.Start()
}
