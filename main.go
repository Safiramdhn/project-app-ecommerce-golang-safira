package main

import (
	"net/http"

	"github.com/Safiramdhn/project-app-ecommerce-golang-safira/router"
	"go.uber.org/zap"
)

func main() {
	r, log, PORT, err := router.InitRouter()
	if err != nil {
		log.Fatal("init router failed", zap.Error(err))
		return
	}
	defer log.Sync()

	log.Info("Server started on port", zap.String("port", PORT))
	if err := http.ListenAndServe(":"+PORT, r); err != nil {
		log.Fatal("listen and serve failed", zap.Error(err))
	}
}
