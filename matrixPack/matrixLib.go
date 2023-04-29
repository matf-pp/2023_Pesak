package matrixPack

import (
	"main/mat"
	"main/mathPack"

	"github.com/veandco/go-sdl2/sdl"
)

// mora ovo ovde /limun
var BrPiksPoCestici int32 = 9000
const SirinaKan, VisinaKan = 360, 216

var Pause bool = false
var TMode bool = false
var NMode bool = false
var DMode bool = false
var TxtMode bool = true

func ClampCoords(x int32, y int32) (int32, int32) {
	return mathPack.MinInt32(mathPack.MaxInt32(x, 0), SirinaKan-1),
		mathPack.MinInt32(mathPack.MaxInt32(y, 0), VisinaKan-1)
}

func ZazidajMatricu(matrix [][]mat.Cestica) [][]mat.Cestica {
	for i := 0; i < SirinaKan; i++ {
		matrix[i][0], matrix[i][VisinaKan-1] = mat.NewCestica(mat.Zid), mat.NewCestica(mat.Zid)
		matrix[i][1], matrix[i][VisinaKan-2] = mat.NewCestica(mat.Zid), mat.NewCestica(mat.Zid)
	}
	for j := 0; j < VisinaKan; j++ {
		matrix[0][j], matrix[SirinaKan-1][j] = mat.NewCestica(mat.Zid), mat.NewCestica(mat.Zid)
		matrix[1][j], matrix[SirinaKan-2][j] = mat.NewCestica(mat.Zid), mat.NewCestica(mat.Zid)
	}
	return matrix
}

func NapraviSlajs() [][]mat.Cestica {
	slajs := make([][]mat.Cestica, SirinaKan)
	for i := 0; i < SirinaKan; i++ {
		kolona := make([]mat.Cestica, VisinaKan)
		for j := 0; j < VisinaKan; j++ {
			kolona[j] = mat.NewCestica(mat.Prazno)
		}
		slajs[i] = kolona
	}
	return slajs
}

func Render(matrix [][]mat.Cestica, surface *sdl.Surface) {
	for i := 0; i < SirinaKan; i++ {
		for j := 0; j < VisinaKan; j++ {
			// njanja: braaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaate je l stvarno toliko teško kompajleru da implicitno konvertuje int u int32
			rect := sdl.Rect{int32(i) * BrPiksPoCestici, int32(j) * BrPiksPoCestici, BrPiksPoCestici, BrPiksPoCestici}
			if TMode {
				bojaTemp := IzracunajTempBoju(matrix[i][j])
				surface.FillRect(&rect, bojaTemp)
			} else if DMode {
				gustTemp := mat.GustinaBoja[matrix[i][j].Materijal]
				surface.FillRect(&rect, gustTemp)
			} else {
				surface.FillRect(&rect, mat.Boja[matrix[i][j].Materijal])
			}
		}
	}
}

var MinTempRendered uint64 = 29315
var MaxTempRendered uint64 = 29316

func IzracunajTempBoju(zrno mat.Cestica) uint32 {
	//	minTemp := mat.MinTemp
	//	maxTemp := mat.MaxTemp

	temperatura := zrno.Temperatura

	// tMin         temp                  tMax
	// 0            xx                    255

	//	(temp - tMin) / (tMax - tMin) = xx / 255
	// xx = 255(temp-tMin)/(tMax-tMin)

	var crvenaKomponenta uint32 = uint32(255 * (temperatura - MinTempRendered) / (MaxTempRendered - MinTempRendered))
	var plavaKomponenta uint32 = uint32(255 - crvenaKomponenta)
	var zelenaKomponenta uint32 = uint32(63)

	if zrno.Materijal == mat.Prazno {
        crvenaKomponenta, plavaKomponenta, zelenaKomponenta = crvenaKomponenta/2, plavaKomponenta/2, zelenaKomponenta/2
    }

	var boja uint32 = (crvenaKomponenta*256+zelenaKomponenta)*256 + plavaKomponenta
	return boja
	/**/
}

func IzbrojiCesticeKamenLavu(matrix [][]mat.Cestica) (int, int, int) {
	x, y, z := 0, 0, 0
	for j := 1; j < VisinaKan-1; j++ {
		for i := 1; i < SirinaKan-1; i++ {
			if matrix[i][j].Materijal != mat.Prazno {
				x++
			}
			if matrix[i][j].Materijal == mat.Lava {
				y++
			}
			if matrix[i][j].Materijal == mat.Kamen {
				z++
			}
		}
	}
	return x, y, z
}