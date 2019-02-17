package barcodegen

import "image"

type Barcoder interface {
	Barcode(data []byte) (image.Image, error)
}
