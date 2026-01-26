package main

import (
	"bytes"
	"extractor/licenser"
	"fmt"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"time"

	"github.com/gen2brain/go-fitz"
	"github.com/mechiko/utility"
)

type SubImager interface {
	SubImage(r image.Rectangle) image.Image
}

func main() {

	var files []string

	lic, err := licenser.New(licenser.MAC, "")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v", lic)

	start := time.Now()
	root := "."
	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			utility.MessageBox("ошибка", fmt.Sprintf("%v", err))
			return err
		}
		if filepath.Ext(path) == ".pdf" {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		utility.MessageBox("ошибка", fmt.Sprintf("%v", err))
		panic(err)
	}
	mm := make(map[string][]string)
	count := 0
	for _, file := range files {
		doc, err := fitz.New(file)
		if err != nil {
			utility.MessageBox("ошибка лицензии", fmt.Sprintf("%v\n%s", err, "перезапустите программу"))
			panic(err)
		}
		// Extract pages as images
		fmt.Println(file)
		mm[file] = make([]string, 0)
		count += doc.NumPage()
		for n := 0; n < doc.NumPage(); n++ {
			img, err := doc.Image(n)
			if err != nil {
				utility.MessageBox("ошибка", fmt.Sprintf("%v", err))
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
					utility.MessageBox("ошибка", fmt.Sprintf("%v", err))
					panic(err)
				}
			} else {
				err = png.Encode(&b, img)
				if err != nil {
					utility.MessageBox("ошибка", fmt.Sprintf("%v", err))
					panic(err)
				}
			}
			s, err := decode(b.Bytes())
			if err != nil {
				mm[file] = append(mm[file], fmt.Sprintf("%d - %v", n+1, err))
				fmt.Printf("error %d - %v\n", n+1, err)
				fn := fmt.Sprintf("%d_%s.png", n+1, filepath.Base(file))
				errFile := os.WriteFile(fn, b.Bytes(), 0644)
				if errFile != nil {
					utility.MessageBox("Ошибка записи файла:", fmt.Sprintf("%v", err))
				}
				continue
			}
			if len(s) > 1 {
				mm[file] = append(mm[file], s[1:])
			} else {
				mm[file] = append(mm[file], s)
			}
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
				f.Close()
				utility.MessageBox("Ошибка записи файла:", fmt.Sprintf("%v", err))
				panic(err)
			}
		}
		f.Close()
	}
	fmt.Printf("затрачено времени %s на %d марок", time.Since(start), count)
}
