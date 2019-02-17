package main

import (
	"log"
	"net/http"

	"github.com/indiependente/barcode/barcodegen/barcode128"
	"github.com/indiependente/barcode/handlers"
)

func main() {
	bcgen := &barcode128.Code128Barcoder{}
	srv := &handlers.BarcodeServer{
		Bcg: bcgen,
	}

	http.HandleFunc("/barcode128", srv.GetBarcode)
	log.Printf("Starting server on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
