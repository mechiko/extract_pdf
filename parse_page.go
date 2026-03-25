package main

import (
	"fmt"
	"image"
)

func parsePage(img *image.RGBA, numPage int) (out []string, err error) {
	out = make([]string, 0)
	bounds := img.Bounds()
	if numPage == 0 {
		fmt.Printf("размеры этикетки %dx%d\n", bounds.Dx(), bounds.Dy())
	}
	if bounds.Dx() > 2000 && bounds.Dy() > 3000 {
		arr, err := parsePageBig(img, numPage)
		if err != nil {
			return out, fmt.Errorf("%w", err)
		}
		out = append(out, arr...)
	} else {
		arr, err := parsePageBig(img, numPage)
		if err != nil {
			return out, fmt.Errorf("%w", err)
		}
		out = append(out, arr...)
	}
	return out, nil
}

func parsePageBig(img *image.RGBA, numPage int) (out []string, err error) {
	out = make([]string, 0)
	prefix := fmt.Sprintf("page_%d_1", numPage)
	if s, err := decodeRect(img, 70, 170, prefix); err != nil {
		return out, fmt.Errorf("%w", err)
	} else {
		out = append(out, s)
	}
	prefix = fmt.Sprintf("page_%d_2", numPage)
	if s, err := decodeRect(img, 70, 885, prefix); err != nil {
		return out, fmt.Errorf("%w", err)
	} else {
		out = append(out, s)
	}
	prefix = fmt.Sprintf("page_%d_3", numPage)
	if s, err := decodeRect(img, 70, 1595, prefix); err != nil {
		return out, fmt.Errorf("%w", err)
	} else {
		out = append(out, s)
	}
	prefix = fmt.Sprintf("page_%d_4", numPage)
	if s, err := decodeRect(img, 70, 2315, prefix); err != nil {
		return out, fmt.Errorf("%w", err)
	} else {
		out = append(out, s)
	}

	return out, nil
}

func parsePageSingle(img *image.RGBA, _ int) (out []string, err error) {
	out = make([]string, 0)
	s, err := decode(img)
	if err != nil {
		return out, fmt.Errorf("%w", err)
	}
	out = append(out, s)
	return out, nil
}

const shiftRect = 280

func decodeRect(img *image.RGBA, x, y int, prefix string) (out string, err error) {
	cropSize := image.Rect(x, y, x+shiftRect, y+shiftRect)
	croppedImage := img.SubImage(cropSize)
	if debug {
		if err := writeImagePng(croppedImage, prefix); err != nil {
			return out, fmt.Errorf("%w", err)
		}
	}
	s, err := decode(croppedImage)
	if err != nil {
		return out, fmt.Errorf("%w", err)
	}
	return s, nil
}
