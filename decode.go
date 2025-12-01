package main

import (
	"bytes"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"

	"github.com/mechiko/dmxing"
	"github.com/mechiko/dmxing/datamatrix"
)

func decode(bb []byte) (string, error) {
	img, _, err := image.Decode(bytes.NewReader(bb))
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}
	// prepare BinaryBitmap
	bmp, err := dmxing.NewBinaryBitmapFromImage(img)
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}

	// decode image
	datamatrixReader := datamatrix.NewDataMatrixReader()
	result, err := datamatrixReader.Decode(bmp, nil)
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}
	return result.GetText(), nil
}
