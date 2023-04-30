package screenPack

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

// čuva kanvas kao sliku apskejlovanu na veličinu ekrana
func SaveImage(matrix [][]mat.Cestica, scaleFactor int) {
	// pravimo praznu sliku pa je popunjavamo
	width := len(matrix)
	height := len(matrix[0])

	targetWidth := width * scaleFactor
	targetHeight := height * scaleFactor

	newImg := image.NewRGBA(image.Rect(0, 0, targetWidth, targetHeight))
	for y := 0; y < targetHeight; y++ {
		for x := 0; x < targetWidth; x++ {
			var hexColor = mat.Boja[matrix[x/scaleFactor][y/scaleFactor].Materijal]
			var RGBColor color.RGBA
			RGBColor.R = uint8((hexColor >> 16) & 0xFF)
			RGBColor.G = uint8((hexColor >> 8) & 0xFF)
			RGBColor.B = uint8(hexColor & 0xFF)
			RGBColor.A = uint8(255)
			newImg.SetRGBA(x, y, RGBColor)
		}
	}

	// pravimo dir za sejv ako ga nemamo
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
	// njanja: radiće uskoro
	/*
		metadata := ""
		for i := 0; i < width; i++ {
			for j := 0; j < height; j++ {
				metadata += fmt.Sprintf("%d:%d:%d;", matrix[i][j].Temperatura, matrix[i][j].SekMat, matrix[i][j].Ticker)
			}
		}


		encdec(true, metadata, filepath.Join(imgDir, fileName))
	*/
}
