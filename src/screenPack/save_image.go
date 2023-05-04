package screenPack

import (
	"main/src/mat"
	"main/src/matrixPack"

	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
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
	// ^ potencijalno nam treba veća slika jer enkriptujemo podatke

	newImg := image.NewRGBA(image.Rect(0, 0, targetWidth, targetHeight))
	var hexColor uint32
	var RGBColor color.RGBA
	for y := 0; y < targetHeight; y++ {
		for x := 0; x < targetWidth; x++ {

			if matrix[x/scaleFactor][y/scaleFactor].Materijal == mat.Vatra || matrix[x/scaleFactor][y/scaleFactor].Materijal == mat.Drvo || matrix[x/scaleFactor][y/scaleFactor].Materijal == mat.Dim {

				hexColor = matrixPack.IzracunajBoju(matrix[x/scaleFactor][y/scaleFactor])
			} else {
				hexColor = mat.Boja[matrix[x/scaleFactor][y/scaleFactor].Materijal]
			}

			RGBColor.R = uint8((hexColor >> 16) & 0xFF)
			RGBColor.G = uint8((hexColor >> 8) & 0xFF)
			RGBColor.B = uint8(hexColor & 0xFF)
			RGBColor.A = uint8(255)
			newImg.SetRGBA(x, y, RGBColor)
		}
	}

	// pravimo dir za sejv ako ga nemamo
	imgDir := "res/images"
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

	// njanja: ovo je jezivo sporo i moraćemo da ga kompresujemo plus ubrzamo samo nz kako
	metadata := ""
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			metadata += fmt.Sprintf("%d:%d:%d;", matrix[j][i].Temperatura, matrix[j][i].SekMat, matrix[j][i].Ticker)
		}
	}

	encdec(true, metadata, filepath.Join(imgDir, fileName))

}
