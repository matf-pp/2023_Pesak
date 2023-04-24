package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"main/mat"
	"os"
	"path/filepath"
)

// TODO SAVING prompt može da piše do not turn off the console ili tako nešto smešno ali nek piše
// TODO nek neko obeleži sve komentare ovde kao moje - njanja (umem da koristim find and replace ali sam malo nervozan sada već)

// TODO skejl apovati rezoluciju slike tako da mečuje rezoluciju ekrana
func save_image(matrix [][]mat.Cestica, width, height int) {
	// pravimo praznu sliku pa je popunjavamo
	newImg := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			var hexColor = mat.Boja[matrix[x][y].Materijal]
			var RGBColor color.RGBA
			RGBColor.R = uint8((hexColor >> 16) & 0xFF)
			RGBColor.G = uint8((hexColor >> 8) & 0xFF)
			RGBColor.B = uint8(hexColor & 0xFF)
			RGBColor.A = uint8(255)
			newImg.SetRGBA(x, y, RGBColor)
		}
	}

	// sejv deo
	// ako nemamo dir
	imgDir := "images"
	if err := os.MkdirAll(imgDir, os.ModePerm); err != nil {
		fmt.Println("Failed to create image directory:", err)
		return
	}

	// sklepamo najmanji indeks nmg kasnim na večeru
	fileName := ""
	index := 1
	for {
		fileName = fmt.Sprintf("pesak_%d.png", index)
		filePath := filepath.Join(imgDir, fileName)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			break
		}
		index++
	}

	// jeeeej
	file, err := os.Create(filepath.Join(imgDir, fileName))
	if err != nil {
		log.Panic("Failed to create output file")
	}
	defer file.Close()

	err = png.Encode(file, newImg)
	if err != nil {
		log.Panic("Failed to encode PNG")
	}

}
