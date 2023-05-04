//f-je za pravljenje platna, racunanje boje cestica, crtanje platna
//varijiable koje odredjuju sta ce se crtati na ekranu
package matrixPack

import (
	"main/src/mat"
	"main/src/mathPack"

	"unsafe"
	"strconv"

	"github.com/veandco/go-sdl2/sdl"
)

var BrPiksPoCestici int32 = 9000
const SirinaKan, VisinaKan = 240, 135

//Pause odredjuje da li je igra pauzirana
var Pause = false
//TMode je temperaturni mod
var TMode = false
//NMode je normalan mod
var NMode = false
//DMode je gustina(density) mod
var DMode = false
//TxtMode je tekst mod(nezavisan od TMode, NMode i DMode)
var TxtMode = true
//ResetSound pusta pesmu od pocetka
var ResetSound = false
//KruzniBrush odredjuje da li je Brush krug ili kvadrat
var KruzniBrush = true

//ClampCoords prima koordinate i ukoliko su one van kanvasa vraca njihove projekcije na ivice
func ClampCoords(x int32, y int32) (int32, int32) {
	return mathPack.MinInt32(mathPack.MaxInt32(x, 0), SirinaKan-1),
		mathPack.MinInt32(mathPack.MaxInt32(y, 0), VisinaKan-1)
}

//ZazidajMatricu dodaje dva sloja piksela Zida oko celog platna
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

//NapraviSlajs je konstruktor matrice Cestica
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

//Render slika matricu na ekran
func Render(matrix [][]mat.Cestica, renderer *sdl.Renderer, texture *sdl.Texture, pixels []byte, simWidth, simHeight int32) {

	var count int
	var hexColor uint32
	for i := 0; i < VisinaKan; i++ {
		for j := 0; j < SirinaKan; j++ {
			if TMode {
				bojaTemp := IzracunajTempBoju(matrix[j][i])
				hexColor = bojaTemp
			} else if DMode {
				gustTemp := mat.GustinaBoja[matrix[j][i].Materijal]
				hexColor = gustTemp
			} else {
				if matrix[j][i].Materijal == mat.Vatra || matrix[j][i].Materijal == mat.Dim {
					hexColor = IzracunajBoju(matrix[j][i])
				} else if matrix[j][i].Materijal == mat.Drvo && matrix[j][i].Temperatura > 47315 {
					hexColor = IzracunajBoju(matrix[j][i])
				} else {
					hexColor = mat.Boja[matrix[j][i].Materijal]
				}
			}

			pixels[count] = byte((hexColor >> 16) & 0xFF)
			pixels[count+1] = byte((hexColor >> 8) & 0xFF)
			pixels[count+2] = byte(hexColor & 0xFF)
			count = count + 3
		}
	}

	var pitch int32 = SirinaKan * 3
	texture.Update(nil, unsafe.Pointer(&pixels[0]), int(pitch))
	renderer.Clear()
	renderer.Copy(texture, nil, &sdl.Rect{X: 0, Y: 0, W: simWidth, H: simHeight})
}

//MinTempRendered je minimalna promenljiva granica za temeraturu koju renderujemo u TMode
var MinTempRendered uint64 = 29315
//MaxTempRendered je maksimalna promenljiva granica za temeraturu koju renderujemo u TMode
var MaxTempRendered uint64 = 29316

//IzracunajBoju racuna boju cestica za one koje se specilajno racunaju
func IzracunajBoju(zrno mat.Cestica) uint32 {

	var boja uint32

	if zrno.Materijal == mat.Dim {
		if zrno.Ticker > 40 {
			return 0x323232
		} else if zrno.Ticker < 0 {
			return 0xaca696
		}
		boje := [6]uint32{0x484343, 0x534c4c, 0x5c4343, 0x63635d, 0x6d6d5c}
		boja = boje[zrno.Ticker/10]
	}
	if zrno.Materijal == mat.Vatra {
		if zrno.Ticker > 11 {
			return 0xfac000
		} else if zrno.Ticker < 0 {
			return 0x400500
		}
		boje := [12]uint32{0x801100, 0xb62203, 0xd73502, 0xfc6400, 0xfc6400, 0xff7500, 0xff7500, 0xff7500, 0xfac000, 0xfac000, 0xfac000, 0xfac000}
		boja = boje[zrno.Ticker]
	} else if zrno.Materijal == mat.Drvo && zrno.Temperatura > 47315 { //200.00c
		temperatura := zrno.Temperatura
		var crvena = uint32(0x99 * (87315 - temperatura + 47315) / 87315)
		var plava = uint32(0 * (87315 - temperatura + 47315) / 87315)
		var zelena = uint32(0x44 * (87315 - temperatura + 47315) / 87315)
		boja = (crvena*256+zelena)*256 + plava
	}

	return boja

}

//IzracunajTempBoju racuna boju cestica u zavisnosti od njihove temperature
func IzracunajTempBoju(zrno mat.Cestica) uint32 {
	temperatura := zrno.Temperatura

	// tMin         temp                  tMax
	// 0            xx                    255

	//	(temp - tMin) / (tMax - tMin) = xx / 255
	// xx = 255(temp-tMin)/(tMax-tMin)

	var crvena uint32 = uint32(255 * (temperatura - MinTempRendered) / (MaxTempRendered - MinTempRendered))
	var plava uint32 = uint32(255 - crvena)
	var zelena uint32 = uint32(63)

	if zrno.Materijal == mat.Prazno {
		crvena, plava, zelena = crvena/2, plava/2, zelena/2
	}
	var boja = (crvena*256+zelena)*256 + plava

	return boja
}

//IzracunajGustBoju racuna boju cestica u zavisnosti od njihove "gustine"
func IzracunajGustBoju(gust int32) uint32 {
	if gust > 0 {
		gust *= 255 / 10
		gust = mathPack.MaxInt32(mathPack.MinInt32(int32(gust), 255), 0)
		gust <<= 8
	} else if gust < 0 {
		gust *= -255 / 10
		gust = mathPack.MinInt32(int32(gust), 255)
		gust += gust << 16
	} else {
		gust = (200 << 16) + (200 << 8) + 200
	}

	hexadeca := strconv.FormatInt(int64(gust), 16)
	gustBoja, err := strconv.ParseInt(hexadeca, 16, 32)
	if err != nil {
		panic(err)
	}

	return uint32(gustBoja)
}
