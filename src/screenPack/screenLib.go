//Package screenPack sadrzi razne f-je vezane za prikaz slike
package screenPack

import (
	"main/src/mat"
	"main/src/matrixPack"

	"math"

	"github.com/fstanis/screenresolution"
	"github.com/veandco/go-sdl2/sdl"
)

//AutoFitScreen odredjuje da li se ekran transformise automatski
var AutoFitScreen = true

//SirinaUIMargine je konstanta sirine ui margine
var SirinaUIMargine int32 = 10
//VisinaUIMargine je konstanta visine ui margine
var VisinaUIMargine int32 = 10
//SirinaDugmeta je konstanta sirine dugmeta
var SirinaDugmeta int32 = 40
//VisinaDugmeta je konstanta visine dugmeta
var VisinaDugmeta int32 = 20

//BrojMaterijala je broj svih postojecih materijala + promena temperature na toplije i na hladnije
var BrojMaterijala = len(mat.Boja) + 2
//BrojSpecijalnihGumbadi je broj opcija koje ne pripadaju materijalima niti promeni temperature: Pause, Save i Clear
var BrojSpecijalnihGumbadi int32 = 3
//BrojGumbadiPoKoloni je broj materijala i promena temperatura po kolonama
var BrojGumbadiPoKoloni int32 = VisinaProzora/(VisinaDugmeta+VisinaUIMargine) - (BrojSpecijalnihGumbadi)
//BrojKolona odredjuje broj kolona koje ce Gumbadi sadrzati
var BrojKolona int32 = int32(math.Ceil(float64(BrojMaterijala) / float64(BrojGumbadiPoKoloni)))
//MarginaZaGumbad je margina koja ostavlja mesta izmedju platna i dugmica
var MarginaZaGumbad = 2*SirinaUIMargine + SirinaDugmeta
//SirinaProzora je sirina celog prozora
var SirinaProzora = matrixPack.SirinaKan*matrixPack.BrPiksPoCestici + MarginaZaGumbad
//VisinaProzora je visina celog prozora
var VisinaProzora = matrixPack.VisinaKan * matrixPack.BrPiksPoCestici

//KursorPoslednjiX je posledja x koordinata misa
var KursorPoslednjiX = int32(0)
//KursorPoslednjiY je posledja y koordinata misa
var KursorPoslednjiY = int32(0)

//VelicinaKursora je velicina cetke
var VelicinaKursora int32 = 8
//MaxKursor je maksimalna velicina cetke
var MaxKursor int32 = 64
//GUIBoja je pozadinska boja navigacije
var GUIBoja uint32 = 0x111122

//TrenutniMat je poslednji materijal koji je korisnik izabran
var TrenutniMat mat.Materijal = mat.Pesak

//FitToScreen je f-ja koja vraca promenjenu skalu, promenjenu visinu i sirinu kanvasa
func FitToScreen(screenPercentage int) (int32, int32, int32) {
	resolution := screenresolution.GetPrimary()
	adjustedScale := int32((float64(screenPercentage) / float64(100)) * float64(resolution.Height) / float64(matrixPack.VisinaKan))

	VisinaUIMargine = int32(float64(VisinaUIMargine*adjustedScale*matrixPack.VisinaKan) / 720)
	SirinaUIMargine = int32(float64(SirinaUIMargine*adjustedScale*matrixPack.VisinaKan) / 720)
	VisinaDugmeta = int32(float64(VisinaDugmeta*adjustedScale*matrixPack.VisinaKan) / 720)
	SirinaDugmeta = int32(float64(SirinaDugmeta*adjustedScale*matrixPack.VisinaKan) / 720)

	return adjustedScale, matrixPack.SirinaKan * adjustedScale, matrixPack.VisinaKan * adjustedScale
}

//ProveriPritisakNaGumb proverava da li se desio pritisak vezan za neki deo Gumba
func ProveriPritisakNaGumb(matrix [][]mat.Cestica, x, y int32) {
	if x > SirinaProzora-MarginaZaGumbad+SirinaUIMargine && x < SirinaProzora-SirinaUIMargine {
		if y < (VisinaUIMargine+VisinaDugmeta)*int32(len(mat.Boja)-1) && y%(VisinaUIMargine+VisinaDugmeta) > VisinaUIMargine {
			TrenutniMat = mat.Materijal(y / (VisinaUIMargine + VisinaDugmeta))
		}

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

//CreateWindow pravi novi prozor i onda ga vraca
func CreateWindow() *sdl.Window {
	window, err := sdl.CreateWindow("pesak", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		SirinaProzora, VisinaProzora, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	return window
}

//CreateSurface pravi novu povrsinu i onda je vraca
func CreateSurface(window *sdl.Window) *sdl.Surface {
	surface, err := window.GetSurface()
	if err != nil {
		panic(err)
	}
	surface.FillRect(nil, 0x111122)
	return surface
}

//CreateRenderer pravi novi rednerer i onda ga vraca
func CreateRenderer(window *sdl.Window) *sdl.Renderer {
	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}

	return renderer
}

//RenderujGumbZaSveMaterijale renderuje Gumb za svaki od materijala prisutnih u navigaciji
func RenderujGumbZaSveMaterijale(renderer *sdl.Renderer) {
	for i := range mat.Boja {
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
}

//CreateSpecialGumb pravi poseban Gumb i vraca ga
func CreateSpecialGumb(index int32) sdl.Rect {
	gumb := sdl.Rect{X: int32(SirinaProzora - SirinaUIMargine - SirinaDugmeta),
		Y: int32(VisinaProzora - index*VisinaUIMargine - index*VisinaDugmeta), W: SirinaDugmeta, H: VisinaDugmeta}
	return gumb
}

//CreateBrushGumb pravi Gumb za Brush i vraca ga
func CreateBrushGumb() sdl.Rect {
	brushGumb := sdl.Rect{X: int32(SirinaProzora - SirinaUIMargine - SirinaDugmeta),
		Y: int32(VisinaProzora - 4*VisinaUIMargine - 4*VisinaDugmeta), W: SirinaDugmeta, H: VisinaDugmeta}
	return brushGumb
}

//CreatePlayGumb pravi Gumb za dugme koje pauzira i pokrece igru, i vraca ga
func CreatePlayGumb() sdl.Rect {
	plejGumb := sdl.Rect{X: int32(SirinaProzora - SirinaUIMargine - SirinaDugmeta),
		Y: int32(VisinaProzora - 3*VisinaUIMargine - 3*VisinaDugmeta), W: SirinaDugmeta, H: VisinaDugmeta}
	return plejGumb
}

//CreateSaveGumb pravi Gumb za dugme koje cuva kanvas kao sliku png formata i vraca ga
func CreateSaveGumb() sdl.Rect {
	sejvGumb := sdl.Rect{X: int32(SirinaProzora - SirinaUIMargine - SirinaDugmeta),
		Y: int32(VisinaProzora - 2*VisinaUIMargine - 2*VisinaDugmeta), W: SirinaDugmeta, H: VisinaDugmeta}
	return sejvGumb
}

//CreateResetGumb pravi Gumb za dugme koje resetuje celu matricu na materijal Prazno
func CreateResetGumb() sdl.Rect {
	resetGumb := sdl.Rect{X: int32(SirinaProzora - SirinaUIMargine - SirinaDugmeta),
		Y: int32(VisinaProzora - VisinaUIMargine - VisinaDugmeta), W: SirinaDugmeta, H: VisinaDugmeta}
	return resetGumb
}

//CreateTexture pravi teksturu i vraca je
func CreateTexture(renderer *sdl.Renderer) *sdl.Texture {
	texture, err := renderer.CreateTexture(uint32(sdl.PIXELFORMAT_RGB24), sdl.TEXTUREACCESS_STATIC, matrixPack.SirinaKan, matrixPack.VisinaKan)
	if err != nil {
		panic(err)
	}
	return texture
}

//InitEverything inicira ceo sdl
func InitEverything() {
	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		panic(err)
	}
}

//UpdateRazmere menja razmere tako da kanvas ne prelazi u navigiju
func UpdateRazmere() {
	MarginaZaGumbad = BrojKolona*(SirinaDugmeta+SirinaUIMargine) + SirinaUIMargine
	SirinaProzora += MarginaZaGumbad
}

//SledeciMaterijal se poziva prilikom skrola na dole
func SledeciMaterijal() {
	if TrenutniMat < 16 {
		TrenutniMat++
		return
	}
	TrenutniMat = 0
	return
}

//PrethodniMaterijal se poziva prilikom skrola na gore
func PrethodniMaterijal() {
	if TrenutniMat > 0 {
		TrenutniMat--
		return
	}
	TrenutniMat = 16
	return
}
