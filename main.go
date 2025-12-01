package main

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"os"
	"path/filepath"

	"github.com/gen2brain/go-fitz"
)

type SubImager interface {
	SubImage(r image.Rectangle) image.Image
}

func main() {

	var files []string

	root := "."
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) == ".pdf" {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	mm := make(map[string][]string)
	for _, file := range files {
		doc, err := fitz.New(file)
		if err != nil {
			panic(err)
		}
		// Extract pages as images
		fmt.Println(file)
		mm[file] = make([]string, 0)
		for n := 0; n < doc.NumPage(); n++ {
			img, err := doc.Image(n)
			if err != nil {
				panic(err)
			}
			bounds := img.Bounds()
			if n == 0 {
				fmt.Printf("размеры этикетки %dx%d\n", bounds.Dx(), bounds.Dy())
			}
			var b bytes.Buffer
			if bounds.Dx() > 400 && bounds.Dy() > 700 {
				cropSize := image.Rect(0, 0, 300, 300)
				cropSize = cropSize.Add(image.Point{60, 150})
				croppedImage := img.SubImage(cropSize)
				err = png.Encode(&b, croppedImage)
				if err != nil {
					panic(err)
				}
			} else {
				err = png.Encode(&b, img)
				if err != nil {
					panic(err)
				}
			}
			s, err := decode(b.Bytes())
			if err != nil {
				mm[file] = append(mm[file], fmt.Sprintf("%d - %v", n+1, err))
				fmt.Printf("error %d - %v\n", n+1, err)
				fn := fmt.Sprintf("%d_%s.png", n+1, file)
				errFile := os.WriteFile(fn, b.Bytes(), 0644)
				if errFile != nil {
					fmt.Println("Ошибка записи файла:", errFile)
				}
				continue
			}
			mm[file] = append(mm[file], s[1:])
		}
	}
	for k, v := range mm {
		f, err := os.Create(k + ".csv")
		if err != nil {
			panic(err)
		}
		for _, str := range v {
			_, err = f.Write([]byte(str + "\n"))
			if err != nil {
				panic(err)
			}
		}
		defer f.Close()
	}
}
