package router

import (
	"github.com/Safiramdhn/project-app-ecommerce-golang-safira/database"
	"github.com/Safiramdhn/project-app-ecommerce-golang-safira/handlers"
	"github.com/Safiramdhn/project-app-ecommerce-golang-safira/repository"
	"github.com/Safiramdhn/project-app-ecommerce-golang-safira/service"
	"github.com/Safiramdhn/project-app-ecommerce-golang-safira/util"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func InitRouter() (*chi.Mux, *zap.Logger, string, error) {
	r := chi.NewRouter()
	config, err := util.InitConfig()
	if err != nil {
		return nil, nil, "", err
	}
	logger := util.InitLog(config)

	logger.Info("Starting database connection")
	db, err := database.InitDatabase(config)
	if err != nil {
		return nil, nil, "", err
	}

	repositories := repository.NewMainRepository(db, logger)
	services := service.NewMainService(repositories, logger)
	handlers := handlers.NewMainHandler(services, logger, config)

	r.Route("/api", func(r chi.Router) {
		r.Post("/register", handlers.UserHandler.RegisterHanlder)
		r.Get("/login", handlers.UserHandler.LoginHandler)
	})

	return r, logger, config.Port, nil
}
