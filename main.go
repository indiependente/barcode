package main

import (
	"net/http"

	"github.com/indiependente/barcode/pkg/barcodegen/barcode128"
	"github.com/indiependente/barcode/pkg/handlers"
	"github.com/indiependente/barcode/pkg/logging"
)

const serviceName = "barcode"

func main() {

	logger := logging.GetLogger(serviceName, logging.DEBUG)
	bcgen := &barcode128.Code128Barcoder{}
	srv := &handlers.BarcodeServer{
		Bcg:    bcgen,
		Logger: logger,
	}

	http.HandleFunc("/barcode128", srv.GetBarcode)
	logger.Info("Starting server on port 8080...")
	logger.Fatal("server stopped unexpectedly", http.ListenAndServe(":8080", nil))
}
