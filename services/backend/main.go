package main

import (
	"net/http"
	"os"
	"strings"

	"github.com/indiependente/barcode/pkg/logging"
	"github.com/indiependente/barcode/services/backend/barcodegen/barcode128"
	"github.com/indiependente/barcode/services/backend/handlers"
	"github.com/indiependente/barcode/services/backend/store/redis"
	"github.com/julienschmidt/httprouter"
)

const serviceName = "barcode_backend"

func main() {
	var logger *logging.Logger
	debug := strings.EqualFold(os.Getenv("LOG_LEVEL"), "DEBUG")

	if debug {
		logger = logging.GetLogger(serviceName, logging.DEBUG)
	} else {
		logger = logging.GetLogger(serviceName, logging.INFO)
	}

	redisAddr := os.Getenv("REDIS_ADDR")
	redisPort := os.Getenv("REDIS_PORT")
	store := redis.InitStore(redisAddr + ":" + redisPort)
	bcgen := &barcode128.Code128Barcoder{}
	srv := &handlers.BarcodeServer{
		Bcg:    bcgen,
		Logger: logger,
		Store:  store,
	}

	router := httprouter.New()
	router.GET("/", srv.Index)
	router.GET("/barcode128/:code", srv.GetBarcode)

	logger.Info("Starting server on port 8080...")
	logger.Fatal("server stopped unexpectedly", http.ListenAndServe(":8080", router))
}
