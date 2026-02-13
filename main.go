package main

import (
	router "github.com/bncunha/erp-api/src/api"
	controller "github.com/bncunha/erp-api/src/api/controllers"
	"github.com/bncunha/erp-api/src/application/ports"
	"github.com/bncunha/erp-api/src/application/service"
	"github.com/bncunha/erp-api/src/application/usecase"
	"github.com/bncunha/erp-api/src/infrastructure/bcrypt"
	email_brevo "github.com/bncunha/erp-api/src/infrastructure/email/brevo"
	"github.com/bncunha/erp-api/src/infrastructure/logs"
	"github.com/bncunha/erp-api/src/infrastructure/observability"
	"github.com/bncunha/erp-api/src/infrastructure/persistence"
	"github.com/bncunha/erp-api/src/infrastructure/repository"
	config "github.com/bncunha/erp-api/src/main"
)

func main() {
	bcrypt := bcrypt.NewBcrypt()
	logs.NewLogs()
	logs.Logger.Infof("Iniciando aplicação")

	config, err := config.LoadConfig()
	if err != nil {
		logs.Logger.Fatalf("Erro buscar variaveis de ambiente", err)
	}

	obs := observability.NewObservability(observability.NewNewRelicObservability())
	err = obs.SetupObservability(config)
	if err != nil {
		logs.Logger.Fatalf("Erro ao configurar observabilidade", err)
	}

	persistence := persistence.NewPersistence(config)
	db, err := persistence.ConnectDb()
	if err != nil {
		logs.Logger.Fatalf("Erro ao conectar no banco de dados", err)
	}
	defer persistence.CloseConnection(db)

	emailBrevo := email_brevo.NewEmailBrevo(email_brevo.EmailBrevoConfig{ApiKey: config.BREVO_API_KEY})

	repository := repository.NewRepository(db)
	repository.SetupRepositories()

	ports := ports.NewPorts(bcrypt, emailBrevo)

	useCase := usecase.NewApplicationUseCase(repository, config, ports)
	useCase.SetupUseCases()

	service := service.NewApplicationService(repository, useCase, ports)
	service.SetupServices()

	controller := controller.NewController(service)
	controller.SetupControllers()

	r := router.NewRouter(controller, obs)
	r.SetupCors(config.APP_ENV)
	r.SetupRoutes()
	r.Start()
}
