package barcode128

import (
	"image"
	"strconv"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/code128"
	"github.com/indiependente/barcode/barcodegen"
	"github.com/pkg/errors"
)

type Code128Barcoder struct{}

func (cb *Code128Barcoder) Barcode(data []byte) (image.Image, error) {
	strdata := string(data)
	if !isValidData(strdata) {
		return nil, errors.Wrapf(barcodegen.ErrInvalidData, "Could not generate barcode: %s", data)
	}
	// Create the barcode
	code, err := code128.Encode(strdata)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not encode data: %s", data)
	}

	// Scale the barcode to 200x200 pixels
	scaled, err := barcode.Scale(code, 200, 200)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not scale barcode: %s", data)
	}
	return scaled, nil
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
