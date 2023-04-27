package screenPack

import (
	"main/mat"
	"main/matrixPack"

	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"math"
	"os"

	"github.com/chai2010/webp"
	"github.com/sergeymakinen/go-bmp"
)

// euklidska razdaljina između boja
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

// otvara sliku risajzuje je i pretvori je u matricu pescanih boja velicine kanvasa
func UcitajSliku(filePath string, matrix, bafer [][]mat.Cestica) error {
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

	// tražimo najbolji materijal i postavljamo materijale u matriks
	for y := 0; y < newHeight; y++ {
		for x := 0; x < newWidth; x++ {
			color := resized.At(x, y).(color.RGBA)
			var min_dist float64 = Distance(color, uint32(mat.Prazno))
			var best_mat mat.Materijal = mat.Prazno
			for materijal, matBoja := range mat.Boja {
				if materijal != mat.Zid {
					var dist = Distance(color, matBoja)
					if dist < min_dist {
						min_dist = dist
						best_mat = materijal
					}
				}
			}

			// njanja: spera kaže da ne moram da apdejtujem oba ali ja mu ne verujem
			matrix[x][y] = mat.NewCestica(best_mat)
			bafer[x][y] = matrix[x][y]
		}
	}

	matrixPack.ZazidajMatricu(matrix)
	matrixPack.ZazidajMatricu(bafer)

	return nil
}
