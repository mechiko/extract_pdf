package main

import (
	"bytes"
	"fmt"
	"image"
	_ "image/png"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/datamatrix"
)

func decode1(bb []byte) (string, error) {
	// open and decode image file

	img, _, err := image.Decode(bytes.NewReader(bb))
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}

	// prepare BinaryBitmap
	bmp, err := gozxing.NewBinaryBitmapFromImage(img)
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}

	// decode image
	qrReader := datamatrix.NewDataMatrixReader()
	result, err := qrReader.Decode(bmp, nil)
	return result.GetText(), err
}
