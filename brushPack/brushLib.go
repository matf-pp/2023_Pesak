package brushPack

import (
	"main/mat"
	"main/matrixPack"
	"main/screenPack"

	"math"

	"github.com/veandco/go-sdl2/sdl"
)

const KruzniBrush = false

func Brush(matrix [][]mat.Cestica, x int32, y int32, state uint32) {
	//TODO za srednji klik da uzme materijal na koj mis trenutno pokazuje i postavi ga kao trenutni
	//ukoliko nije u pitanju Zid ili Prazno. Nije mi pri ruci mis, mrzi me da trazim koj je to stejt -s
	//a jeste sabani mogli ste ovo trideset puta uraditi danas -s
	if x > matrixPack.SirinaKan*matrixPack.BrPiksPoCestici {
		return
	}
	if state != 1 && state != 2 && state != 4 {
		return
	}

	if state == 1 {
		for i := -screenPack.VelicinaKursora; i < screenPack.VelicinaKursora; i++ {
			for j := -screenPack.VelicinaKursora; j < screenPack.VelicinaKursora; j++ {
				tx, ty := matrixPack.ClampCoords(x/matrixPack.BrPiksPoCestici+i, y/matrixPack.BrPiksPoCestici+j)
				if screenPack.TrenutniMat != mat.Toplo && screenPack.TrenutniMat != mat.Hladno {
					if matrix[tx][ty].Materijal == mat.Prazno || (screenPack.TrenutniMat == mat.Prazno && matrix[tx][ty].Materijal != mat.Zid) {
						matrix[tx][ty] = mat.NewCestica(screenPack.TrenutniMat)
					}
				} else {
					if screenPack.TrenutniMat == mat.Toplo && matrix[tx][ty].Materijal != mat.Zid {
						if matrix[tx][ty].Temperatura+1000 > mat.MaxTemp {
							matrix[tx][ty].Temperatura = mat.MaxTemp
						} else {
							matrix[tx][ty].Temperatura += 1000
						}
					} else if screenPack.TrenutniMat == mat.Hladno && matrix[tx][ty].Materijal != mat.Zid {
						if matrix[tx][ty].Temperatura-1000 > mat.MaxTemp {
							matrix[tx][ty].Temperatura = mat.MinTemp
						} else {
							matrix[tx][ty].Temperatura -= 1000
						}
					}
				}
			}
		}
	}

	if state == 2 {
		tx, ty := matrixPack.ClampCoords(x/matrixPack.BrPiksPoCestici, y/matrixPack.BrPiksPoCestici)
		screenPack.TrenutniMat = matrix[tx][ty].Materijal
	}

	if state == 4 {
		for i := -screenPack.VelicinaKursora; i < screenPack.VelicinaKursora; i++ {
			for j := -screenPack.VelicinaKursora; j < screenPack.VelicinaKursora; j++ {
				tx, ty := matrixPack.ClampCoords(x/matrixPack.BrPiksPoCestici+i, y/matrixPack.BrPiksPoCestici+j)
				if matrix[tx][ty].Materijal != mat.Zid {
					//napomenuo bih da prazne cestice ovde brisemo i pravimo opet da bismo resetovali temp
					//inace bi bilo efikasnije samo postaviti im Materijal na Prazno, NAGADJAM
					//takodje mozda je brze izmeniti polja cestice nego praviti novu, ne znam, ostavio bih to bencmarkingu
					matrix[tx][ty] = mat.NewCestica(mat.Prazno)
				}
			}
		}
	}
}

func OblikCetkice(KruzniBrush bool, renderer *sdl.Renderer) {
	// njanja: koliko odvratne zagrade
	if KruzniBrush {
		// krug
		renderer.SetDrawColor(255, 255, 255, 255)
		radius := int(screenPack.VelicinaKursora * matrixPack.BrPiksPoCestici)
		numSegments := int(math.Ceil(float64(radius) / 2.0))
		for i := 0; i < numSegments; i++ {
			angle1 := float64(i) / float64(numSegments) * math.Pi * 2.0
			angle2 := float64(i+1) / float64(numSegments) * math.Pi * 2.0
			x1 := float64(screenPack.KursorPoslednjiX) + float64(radius)*math.Cos(angle1)
			y1 := float64(screenPack.KursorPoslednjiY) + float64(radius)*math.Sin(angle1)
			x2 := float64(screenPack.KursorPoslednjiX) + float64(radius)*math.Cos(angle2)
			y2 := float64(screenPack.KursorPoslednjiY) + float64(radius)*math.Sin(angle2)
			renderer.DrawLine(int32(x1), int32(y1), int32(x2), int32(y2))
		}
	} else {
		// kvadrat
		renderer.SetDrawColor(255, 255, 255, 255)
		cetkica := sdl.Rect{X: screenPack.KursorPoslednjiX - screenPack.VelicinaKursora*matrixPack.BrPiksPoCestici, Y: screenPack.KursorPoslednjiY - screenPack.VelicinaKursora*matrixPack.BrPiksPoCestici, W: int32(2 * screenPack.VelicinaKursora * matrixPack.BrPiksPoCestici), H: int32(2 * screenPack.VelicinaKursora * matrixPack.BrPiksPoCestici)}
		renderer.DrawRect(&cetkica)
	}
}