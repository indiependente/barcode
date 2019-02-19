package handlers

import (
	"image"
	"image/png"
	"io"
	"net/http"
	"time"

	"github.com/indiependente/barcode/pkg/barcodegen"
	"github.com/indiependente/barcode/pkg/logging"
	"github.com/pkg/errors"
)

type BarcodeServer struct {
	Bcg    barcodegen.Barcoder
	Logger *logging.Logger
}

func (srv *BarcodeServer) GetBarcode(w http.ResponseWriter, r *http.Request) {
	srv.Logger.Debug("request received")
	start := time.Now()
	err := srv.sendBarcode(srv.Bcg, w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	srv.Logger.ElapsedTime(time.Since(start)).Debug("request processed")
}

func (srv *BarcodeServer) sendBarcode(bcgen barcodegen.Barcoder, w io.Writer, r *http.Request) error {
	code, ok := r.URL.Query()["code"]
	if !ok || len(code) < 1 {
		return errors.Wrap(ErrMissingCode, "Could not read code parameter in URL")
	}
	srv.Logger.Code(code[0]).Debug("code received")
	data := []byte(code[0])
	img, err := bcgen.Barcode(data)
	if err != nil {
		return errors.Wrap(err, "Could not convert to barcode")
	}
	if err := writePNG(w, http.StatusOK, img); err != nil {
		return errors.Wrapf(err, "Could not send barcode")
	}
	return nil
}

func writePNG(w io.Writer, status int, img image.Image) error {
	err := png.Encode(w, img)
	if err != nil {
		return errors.Wrap(err, "Could not encode to png")
	}
	return nil
}
