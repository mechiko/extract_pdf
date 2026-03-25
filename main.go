package main

import (
	"errors"
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

const outDir = ".out"

var pdfPath = ".data"

const debug = false

func main() {
	var files []string
	start := time.Now()
	if !exists(pdfPath) {
		pdfPath = "."
	}
	// Create path if it doesn't exist (0755 is standard permissions)
	err := os.MkdirAll(outDir, 0755)
	if err != nil {
		utility.MessageBox("ошибка", fmt.Sprintf("%v", err))
		return
	}
	err = filepath.Walk(pdfPath, func(path string, info os.FileInfo, err error) error {
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
			utility.MessageBox("ошибка", fmt.Sprintf("%v", err))
			panic(err)
		}
		file = filepath.Base(file)
		mm[file] = make([]string, 0)
		count += doc.NumPage()
		for n := 0; n < doc.NumPage(); n++ {
			img, err := doc.Image(n)
			if err != nil {
				utility.MessageBox("ошибка", fmt.Sprintf("%v", err))
				panic(err)
			}
			arr, err := parsePage(img, n)
			if err != nil {
				utility.MessageBox(fmt.Sprintf("ошибка разбора страницы %d", n), fmt.Sprintf("%v", err))
				panic(err)
			}
			fmt.Printf("файл %s страница %d найдено марок %d\n", file, n, len(arr))
			mm[file] = append(mm[file], arr...)
		}
	}
	countKM := 0
	for k, v := range mm {
		countKM = countKM + len(v)
		k = filepath.Join(outDir, k)
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
	fmt.Printf("затрачено времени %s на %d марок\n", time.Since(start), countKM)
}

func htmlWrite(doc *fitz.Document, page int) {
	html, err := doc.HTML(page, true)
	if err != nil {
		panic(err)
	}
	f, err := os.Create(filepath.Join(outDir, fmt.Sprintf("test%03d.html", page)))
	if err != nil {
		panic(err)
	}
	_, err = f.WriteString(html)
	if err != nil {
		panic(err)
	}
	f.Close()
}

func writeImagePng(img image.Image, file string) error {
	fileName := filepath.Join(outDir, fmt.Sprintf("%s.png", file))
	f, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	defer f.Close()
	if err := png.Encode(f, img); err != nil {
		return fmt.Errorf("%w", err)
	}
	return nil
}

func exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true // Path exists
	}
	if errors.Is(err, os.ErrNotExist) {
		return false // Path does not exist
	}
	// Path might exist, but there was another error (e.g., permission denied)
	return false
}
