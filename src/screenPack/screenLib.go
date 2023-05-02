package screenPack

import (
	"main/src/mat"
	"main/src/matrixPack"

	"math"

	"github.com/fstanis/screenresolution"
	"github.com/veandco/go-sdl2/sdl"
)

var AutoFitScreen = true

//brate ono dimenzije margine i dugmadi pored platna
var SirinaUIMargine int32 = 10
var VisinaUIMargine int32 = 10
var SirinaDugmeta int32 = 40
var VisinaDugmeta int32 = 20

// njanja: gumb magija ne radi kad nije u mejnu stignite ako hoćete
// ja sada: https://cdn.discordapp.com/emojis/1068966756556738590.webp

//BrojMaterijala brate sve ovo je ja mislim dovoljno recito nazvano i ja nemam vise ni snage ni volje da
//igram po melodiji ovog ukletog ocenjivaca nastavicu neki drugi dan ili nek neko preuzme -s
var BrojMaterijala = len(mat.Boja) + 2
var BrojSpecijalnihGumbadi int32 = 3
var BrojGumbadiPoKoloni int32 = VisinaProzora/(VisinaDugmeta+VisinaUIMargine) - (BrojSpecijalnihGumbadi)
var BrojKolona int32 = int32(math.Ceil(float64(BrojMaterijala) / float64(BrojGumbadiPoKoloni)))
var MarginaZaGumbad = 2*SirinaUIMargine + SirinaDugmeta
var SirinaProzora = matrixPack.SirinaKan*matrixPack.BrPiksPoCestici + MarginaZaGumbad
var VisinaProzora = matrixPack.VisinaKan * matrixPack.BrPiksPoCestici

var KursorPoslednjiX = int32(0)
var KursorPoslednjiY = int32(0)

// var velicinaKursora int32 = 4
var VelicinaKursora int32 = 8
var MaxKursor int32 = 64
var GUIBoja uint32 = 0x111122

var TrenutniMat mat.Materijal = mat.Pesak

// takozvano dinamičko skaliranje ekrana ili nešto ne znam lupio sam
// ako ovo ikada u praksi izbaci nešto što ne staje u ekran javite mi da ga sredim ali mislim da je to besmislen posao
func FitToScreen(screenPercentage int) (int32, int32, int32) {
	resolution := screenresolution.GetPrimary()
	adjustedScale := int32((float64(screenPercentage) / float64(100)) * float64(resolution.Height) / float64(matrixPack.VisinaKan))

	VisinaUIMargine = int32(float64(VisinaUIMargine*adjustedScale*matrixPack.VisinaKan) / 720)
	SirinaUIMargine = int32(float64(SirinaUIMargine*adjustedScale*matrixPack.VisinaKan) / 720)
	VisinaDugmeta = int32(float64(VisinaDugmeta*adjustedScale*matrixPack.VisinaKan) / 720)
	SirinaDugmeta = int32(float64(SirinaDugmeta*adjustedScale*matrixPack.VisinaKan) / 720)

	return adjustedScale, matrixPack.SirinaKan * adjustedScale, matrixPack.VisinaKan * adjustedScale
}

func ProveriPritisakNaGumb(matrix [][]mat.Cestica, x, y int32) {
	if x > SirinaProzora-MarginaZaGumbad+SirinaUIMargine && x < SirinaProzora-SirinaUIMargine {
		// materijali
		if y < (VisinaUIMargine+VisinaDugmeta)*int32(len(mat.Boja)-1) && y%(VisinaUIMargine+VisinaDugmeta) > VisinaUIMargine {
			TrenutniMat = mat.Materijal(y / (VisinaUIMargine + VisinaDugmeta))
		}

		// njanja: hardkodovan broj specijalnih dugmića hvala bogu
		// BRUSH SHAPE
		if y > VisinaProzora-4*(VisinaDugmeta+VisinaUIMargine) && y < VisinaProzora-4*(VisinaDugmeta+VisinaUIMargine)+VisinaDugmeta {
			matrixPack.KruzniBrush = !matrixPack.KruzniBrush
		}

		// PAUZA
		if y > VisinaProzora-3*(VisinaDugmeta+VisinaUIMargine) && y < VisinaProzora-3*(VisinaDugmeta+VisinaUIMargine)+VisinaDugmeta {
			matrixPack.Pause = !matrixPack.Pause
		}

		// SEJV
		if y > VisinaProzora-2*(VisinaDugmeta+VisinaUIMargine) && y < VisinaProzora-2*(VisinaDugmeta+VisinaUIMargine)+VisinaDugmeta {
			sdl.ShowSimpleMessageBox(sdl.MESSAGEBOX_INFORMATION, "pesak", "čuvamo sliku.. molim vas ne gasite struju", nil)
			go SaveImage(matrix, int(matrixPack.BrPiksPoCestici))
			//sdl.ShowSimpleMessageBox(sdl.MESSAGEBOX_INFORMATION, "pesak", "sačuvan B)", nil)
		}

		// RESET
		if y > VisinaProzora-1*(VisinaDugmeta+VisinaUIMargine) && y < VisinaProzora-1*(VisinaDugmeta+VisinaUIMargine)+VisinaDugmeta {
			for j := 0; j < matrixPack.VisinaKan; j++ {
				for i := 0; i < matrixPack.SirinaKan; i++ {
					matrix[i][j] = mat.NewCestica(mat.Prazno)
				}
			}
			matrixPack.ZazidajMatricu(matrix)
		}

	}
}

func CreateWindow() *sdl.Window {
	window, err := sdl.CreateWindow("pesak", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		SirinaProzora, VisinaProzora, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	return window
}

func CreateSurface(window *sdl.Window) *sdl.Surface {
	surface, err := window.GetSurface()
	if err != nil {
		panic(err)
	}
	surface.FillRect(nil, 0x111122)
	return surface
}

func CreateRenderer(window *sdl.Window) *sdl.Renderer {
	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}

	return renderer
}

// njanja: mislim da se više ne ređaju u više kolona otkad je struktura promenjena TODO
func RenderujGumbZaSveMaterijale(renderer *sdl.Renderer) {
	for i, _ := range mat.Boja {
		if i == TrenutniMat {
			gumb := sdl.Rect{X: int32(SirinaProzora - MarginaZaGumbad + ((int32(i)%BrojKolona)*(SirinaDugmeta+SirinaUIMargine) + SirinaUIMargine) - SirinaUIMargine/3),
				Y: int32(VisinaUIMargine+int32(i)/BrojKolona*(VisinaDugmeta+VisinaUIMargine)) - VisinaUIMargine/3, W: SirinaDugmeta + 2*SirinaUIMargine/3, H: VisinaDugmeta + 2*VisinaUIMargine/3}
			renderer.SetDrawColor(uint8(mat.Boja[i]>>16), uint8(mat.Boja[i]>>8), uint8(mat.Boja[i]), 255)
			renderer.FillRect(&gumb)

			gumb = sdl.Rect{X: int32(SirinaProzora - MarginaZaGumbad + ((int32(i)%BrojKolona)*(SirinaDugmeta+SirinaUIMargine) + SirinaUIMargine)),
				Y: int32(VisinaUIMargine + int32(i)/BrojKolona*(VisinaDugmeta+VisinaUIMargine)), W: SirinaDugmeta, H: VisinaDugmeta}
			renderer.SetDrawColor(uint8(mat.Boja[i]>>16)/3*2, uint8(mat.Boja[i]>>8)/3*2, uint8(mat.Boja[i])/3*2, 255)
			renderer.FillRect(&gumb)

			continue
		}
		gumb := sdl.Rect{X: int32(SirinaProzora - MarginaZaGumbad + ((int32(i)%BrojKolona)*(SirinaDugmeta+SirinaUIMargine) + SirinaUIMargine)),
			Y: int32(VisinaUIMargine + int32(i)/BrojKolona*(VisinaDugmeta+VisinaUIMargine)), W: SirinaDugmeta, H: VisinaDugmeta}
		renderer.SetDrawColor(uint8(mat.Boja[i]>>16), uint8(mat.Boja[i]>>8), uint8(mat.Boja[i]), 255)
		renderer.FillRect(&gumb)
	}
	// zašto je ovo vraćalo surface
	// ne sećam se, možda je bio neki check /limun
}

func CreateSpecialGumb(index int32) sdl.Rect {
	gumb := sdl.Rect{X: int32(SirinaProzora - SirinaUIMargine - SirinaDugmeta),
		Y: int32(VisinaProzora - index*VisinaUIMargine - index*VisinaDugmeta), W: SirinaDugmeta, H: VisinaDugmeta}
	return gumb
}

func CreateBrushGumb() sdl.Rect {
	brushGumb := sdl.Rect{X: int32(SirinaProzora - SirinaUIMargine - SirinaDugmeta),
		Y: int32(VisinaProzora - 4*VisinaUIMargine - 4*VisinaDugmeta), W: SirinaDugmeta, H: VisinaDugmeta}
	return brushGumb
}

func CreatePlayGumb() sdl.Rect {
	plejGumb := sdl.Rect{X: int32(SirinaProzora - SirinaUIMargine - SirinaDugmeta),
		Y: int32(VisinaProzora - 3*VisinaUIMargine - 3*VisinaDugmeta), W: SirinaDugmeta, H: VisinaDugmeta}
	return plejGumb
}

func CreateSaveGumb() sdl.Rect {
	sejvGumb := sdl.Rect{X: int32(SirinaProzora - SirinaUIMargine - SirinaDugmeta),
		Y: int32(VisinaProzora - 2*VisinaUIMargine - 2*VisinaDugmeta), W: SirinaDugmeta, H: VisinaDugmeta}
	return sejvGumb
}

func CreateResetGumb() sdl.Rect {
	resetGumb := sdl.Rect{X: int32(SirinaProzora - SirinaUIMargine - SirinaDugmeta),
		Y: int32(VisinaProzora - VisinaUIMargine - VisinaDugmeta), W: SirinaDugmeta, H: VisinaDugmeta}
	return resetGumb
}

func CreateTexture(renderer *sdl.Renderer) *sdl.Texture {
	texture, err := renderer.CreateTexture(uint32(sdl.PIXELFORMAT_RGB24), sdl.TEXTUREACCESS_STATIC, matrixPack.SirinaKan, matrixPack.VisinaKan)
	if err != nil {
		panic(err)
	}
	return texture
}

func InitEverything() {
	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		panic(err)
	}
}

func UpdateRazmere() {
	MarginaZaGumbad = BrojKolona*(SirinaDugmeta+SirinaUIMargine) + SirinaUIMargine
	SirinaProzora += MarginaZaGumbad
}

func SledeciMaterijal() {
	if TrenutniMat < 16 {
		TrenutniMat++
		return
	}
	TrenutniMat = 0
	return
}

func PrethodniMaterijal() {
	if TrenutniMat > 0 {
		TrenutniMat--
		return
	}
	TrenutniMat = 16
	return
}
