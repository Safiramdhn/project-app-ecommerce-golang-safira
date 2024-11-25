package router

import (
	"github.com/Safiramdhn/project-app-ecommerce-golang-safira/database"
	"github.com/Safiramdhn/project-app-ecommerce-golang-safira/handlers"
	"github.com/Safiramdhn/project-app-ecommerce-golang-safira/middleware"
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
	middleware := middleware.NewMiddleware(logger, config)

	r.Route("/api", func(r chi.Router) {
		r.Post("/register", handlers.UserHandler.RegisterHanlder)
		r.Get("/login", handlers.UserHandler.LoginHandler)

		r.Get("/categories", handlers.CategoryHandler.GetAllCategoryHandler)

		r.Route("/products", func(r chi.Router) {
			r.Get("/", handlers.ProductHandler.GetAllProductHandler)
			r.Get("/{id}", handlers.ProductHandler.GetProductByIdHandler)
			r.Get("/recommendation", handlers.RecommendationHandler.GetRecommendationsHandler)
			r.Get("/banner", handlers.RecommendationHandler.GetBannerProduct)
			r.Get("/weekly-promo", handlers.ProductHandler.GetWeeklyPromotionsHandler)
		})

		r.With(middleware.AuthMiddleware).Route("/wishlist", func(r chi.Router) {
			r.Post("/add", handlers.WishlistHandler.AddWishlistHandler)
			r.Get("/", handlers.WishlistHandler.GetWishlistHandler)
			r.Delete("/remove/{id}", handlers.WishlistHandler.RemoveProductFromWishlistHandler)
		})

		r.With(middleware.AuthMiddleware).Route("/user", func(r chi.Router) {
			r.Route("/address", func(r chi.Router) {
				r.Post("/", handlers.AddressHandler.AddAddressHandler)
				r.Put("/{id}", handlers.AddressHandler.UpdateAddressHandler)
				r.Patch("/{id}", handlers.AddressHandler.SetDefaultAddressHandler)
				r.Delete("/{id}", handlers.AddressHandler.DeleteAddressHandler)
				r.Get("/", handlers.AddressHandler.GetAllAddressesHandler)
				r.Get("/{id}", handlers.AddressHandler.GetAddressByIdHandler)
			})
		})

		r.With(middleware.AuthMiddleware).Route("/cart", func(r chi.Router) {
			r.Post("/add-item", handlers.CartHandler.AddToCartHandler)
			r.Get("/", handlers.CartHandler.GetUserCart)
			r.Delete("/remove-item/{id}", handlers.CartHandler.DeleteItemHandler)
			r.Put("/update-item/{id}", handlers.CartHandler.UpdateCartItemHandler)
		})

		r.With(middleware.AuthMiddleware).Route("/orders", func(r chi.Router) {
			r.Post("/", handlers.OrderHandler.CreateOrderHanlder)
			r.Get("/", handlers.OrderHandler.GetOrderHistoryHandler)
			r.Get("/{id}", handlers.OrderHandler.GetOrderDetailsHandler)
		})
	})

	return r, logger, config.Port, nil
}
