package main

import (
	"bytes"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"

	"github.com/mechiko/dmxing"
	"github.com/mechiko/dmxing/datamatrix"
	"github.com/mechiko/dmxing/qrcode/decoder"
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
	hints := map[dmxing.DecodeHintType]interface{}{dmxing.DecodeHintType_PURE_BARCODE: decoder.ErrorCorrectionLevel_L}
	datamatrixReader := datamatrix.NewDataMatrixReader()
	result, err := datamatrixReader.Decode(bmp, hints)
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}
	return result.GetText(), nil
}
