package main

import (
	"fmt"
	"image"
	_ "image/png"
	"strings"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/datamatrix"
)

func decode(img image.Image) (string, error) {
	bmp, err := gozxing.NewBinaryBitmapFromImage(img)
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}
	qrReader := datamatrix.NewDataMatrixReader()
	result, err := qrReader.Decode(bmp, nil)
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}
	return trimPrefixGS(result.GetText()), nil
}

func trimPrefixGS(code string) string {
	return strings.TrimPrefix(code, "\x1d")
}
