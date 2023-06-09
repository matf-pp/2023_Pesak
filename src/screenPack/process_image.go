//Package screenPack sadrzi razne f-je vezane za prikaz slike
package screenPack

import (
	"main/src/mat"
	"main/src/matrixPack"

	"strconv"
	"strings"
	"image"
	"image/color"
	//blank import
	_ "image/jpeg"
	//blank import
	_ "image/png"
	"math"
	"os"

	"github.com/chai2010/webp"
	"github.com/sergeymakinen/go-bmp"
)

//Distance prima boju u RGBA formatu i boju u uint32 formatu i vraca euklidsku razdaljinu izmenju njih
func Distance(c1 color.RGBA, hexColor uint32) float64 {
	var c2 color.RGBA
	c2.R = uint8((hexColor >> 16) & 0xFF)
	c2.G = uint8((hexColor >> 8) & 0xFF)
	c2.B = uint8(hexColor & 0xFF)

	rDelta := int(c1.R) - int(c2.R)
	gDelta := int(c1.G) - int(c2.G)
	bDelta := int(c1.B) - int(c2.B)
	return math.Sqrt(float64(rDelta*rDelta + gDelta*gDelta + bDelta*bDelta))
}

//UcitajSliku otvara sliku risajzuje je i pretvori je u matricu pescanih boja velicine kanvasa
func UcitajSliku(filePath string, matrix [][]mat.Cestica) error {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	img, format, err := image.Decode(file)
	if err != nil {
		// njanja: ovo je MOŽDA bloat ali se baš smorim kad ne mogu da učitam bmp/webp
		switch format {
		case "bmp":
			img, err = bmp.Decode(file)
		case "webp":
			img, err = webp.Decode(file)
		}
		if err != nil {
			return err
		}
	}

	// sakupljamo trenutne dimenzije pa izracunamo nove
	bounds := img.Bounds()
	origWidth := bounds.Max.X
	origHeight := bounds.Max.Y
	zeljenaSirina := len(matrix)
	zeljenaVisina := len(matrix[0])

	aspectRatio := float64(origWidth) / float64(origHeight)
	newWidth := zeljenaSirina
	newHeight := int(float64(zeljenaSirina) / aspectRatio)
	if newHeight > zeljenaVisina {
		newHeight = zeljenaVisina
		newWidth = int(float64(zeljenaVisina) * aspectRatio)
	}

	// upscale (nearest neighbour)
	scaleX := float64(bounds.Max.X) / float64(newWidth)
	scaleY := float64(bounds.Max.Y) / float64(newHeight)
	resized := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))

	for y := 0; y < newHeight; y++ {
		for x := 0; x < newWidth; x++ {
			srcX := int(float64(x)*scaleX + 0.5)
			srcY := int(float64(y)*scaleY + 0.5)
			srcColor := img.At(srcX, srcY)
			resized.Set(x, y, srcColor)
		}
	}

	// njanja: ovo je ekstremno kul napisaću možda nekad uljudne komentare
	// ugl koristimo steganografiju da bismo na slici čuvali temperaturu i to
	saveData := ""
	if format == "png" {
		saveData = encdec(false, "", filePath)
	}
	hasMetadata := false
	var pixelList []string
	if saveData != "" {
		hasMetadata = true
		pixelList = strings.Split(saveData, ";")
	}

	i := 0
	// tražimo najbolji materijal i postavljamo materijale u matriks
	for y := 0; y < newHeight; y++ {
		for x := 0; x < newWidth; x++ {
			color := resized.At(x, y).(color.RGBA)
			var minDist = Distance(color, uint32(mat.Prazno))
			var bestMat mat.Materijal = mat.Prazno
			for materijal, matBoja := range mat.Boja {
				if materijal != mat.Zid {
					var dist = Distance(color, matBoja)
					if dist < minDist {
						minDist = dist
						bestMat = materijal
					}
				}
			}
			matrix[x][y] = mat.NewCestica(bestMat)

			if hasMetadata {
				pixelData := strings.Split(pixelList[i], ":")
				temp, _ := strconv.Atoi(pixelData[0])
				sekmat, _ := strconv.Atoi(pixelData[1])
				ticker, _ := strconv.Atoi(pixelData[2])
				matrix[x][y].Temperatura = uint64(temp)
				matrix[x][y].SekMat = mat.Materijal(sekmat)
				matrix[x][y].Ticker = int32(ticker)
				i++
			}
		}
	}

	matrixPack.ZazidajMatricu(matrix)

	return nil
}
