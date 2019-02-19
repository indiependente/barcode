package handlers

import (
	"fmt"
	"image"
	"image/png"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/indiependente/barcode/pkg/barcodegen"
	"github.com/indiependente/barcode/pkg/logging"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
)

type BarcodeServer struct {
	Bcg    barcodegen.Barcoder
	Logger *logging.Logger
}

func (srv *BarcodeServer) Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
}

func (srv *BarcodeServer) GetBarcode(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	srv.Logger.Debug("GetBarcode request received")
	start := time.Now()

	code := params.ByName("code")
	l := srv.Logger.CodeID(code)

	// sanity checks
	if len(code) < 1 {
		err := errors.Wrap(ErrMissingCode, "Could not read code parameter in URL")
		l.Error("processing failed", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if !isValidData(code) {
		err := errors.Wrap(barcodegen.ErrInvalidData, "Could not generate barcode")
		l.Error("processing failed", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// generate image
	img, err := srv.Bcg.Barcode([]byte(code))
	if err != nil {
		err := errors.Wrap(err, "Could not convert to barcode")
		l.Error("processing failed", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = writePNG(w, http.StatusOK, img)
	if err != nil {
		err := errors.Wrap(err, "Could not write response")
		l.Error("processing failed", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	l.ElapsedTime(time.Since(start)).Debug("GetBarcode request processed")
}

func writePNG(w io.Writer, status int, img image.Image) error {
	err := png.Encode(w, img)
	if err != nil {
		return errors.Wrap(err, "Could not encode to png")
	}
	return nil
}

func isValidData(data string) bool {
	if len(data) != 20 {
		return false
	}
	for _, c := range data {
		_, err := strconv.ParseInt(string(c), 10, 32)
		if err != nil {
			return false
		}
	}
	return true
}
