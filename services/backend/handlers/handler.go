package handlers

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/indiependente/barcode/pkg/logging"
	"github.com/indiependente/barcode/pkg/store"
	"github.com/indiependente/barcode/services/backend/barcodegen"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
)

type BarcodeServer struct {
	Bcg    barcodegen.Barcoder
	Store  store.Storer
	Logger *logging.Logger
}

func (srv *BarcodeServer) Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n") // nolint: errcheck
}

func (srv *BarcodeServer) GetBarcode(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	srv.Logger.Info("GetBarcode request received")
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

	// retrieve from store
	imgByte, err := srv.Store.Get(code)
	if err != nil {
		l.Error("get from store failed", err)
	}

	var img image.Image

	if imgByte != nil && err == nil { // use the cached version
		img, _, err = image.Decode(bytes.NewReader(imgByte))
		if err != nil {
			l.Error("failed image decoding", err)
		}
		l.Debug("image retrieved from cache")
	} else {
		// generate image
		img, err = srv.Bcg.Barcode([]byte(code))
		if err != nil {
			err := errors.Wrap(err, "Could not convert to barcode")
			l.Error("processing failed", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		l.Debug("barcode generated")
		go srv.storeImg(code, img, l)
	}

	// send image
	err = writePNG(w, http.StatusOK, img)
	if err != nil {
		err := errors.Wrap(err, "Could not write response")
		l.Error("processing failed", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	l.Debug("image sent")
	l.ElapsedTime(time.Since(start)).Info("GetBarcode request processed")
}

func writePNG(w io.Writer, status int, img image.Image) error {
	err := png.Encode(w, img)
	if err != nil {
		return errors.Wrap(err, "Could not encode to png")
	}
	return nil
}

func (srv *BarcodeServer) storeImg(code string, img image.Image, l logging.LogChainer) {
	var buf bytes.Buffer
	err := png.Encode(&buf, img)
	if err != nil {
		l.Error("image encoding to png failed", err)
	}
	l.Debug("encoded to PNG")
	// save image
	err = srv.Store.Put(code, buf.Bytes())
	if err != nil {
		l.Error("put to store failed", err)
	} else {
		l.Info("image stored")
	}
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
