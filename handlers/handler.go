package handlers

import (
	"image"
	"image/png"
	"net/http"

	"github.com/indiependente/barcode/barcodegen"
	"github.com/pkg/errors"
)

type BarcodeServer struct {
	Bcg barcodegen.Barcoder
}

func (srv *BarcodeServer) GetBarcode(w http.ResponseWriter, r *http.Request) {
	err := sendBarcode(srv.Bcg, w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func sendBarcode(bcgen barcodegen.Barcoder, w http.ResponseWriter, r *http.Request) error {
	code, ok := r.URL.Query()["code"]
	if !ok || len(code) < 1 {
		return errors.Wrap(ErrMissingCode, "Could not read code parameter in URL")
	}
	data := []byte(code[0])
	img, err := bcgen.Barcode(data)
	if err != nil {
		return errors.Wrap(err, "Could not convert to barcode")
	}
	if err := writePNG(w, 200, img); err != nil {
		return errors.Wrapf(err, "Could not send barcode")
	}
	return nil
}

func writePNG(w http.ResponseWriter, status int, img image.Image) error {
	err := png.Encode(w, img)
	if err != nil {
		return errors.Wrap(err, "Could not encode to png")
	}
	return nil
}
