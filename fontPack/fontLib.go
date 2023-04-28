package fontPack

import (
	"main/mat"
	"main/screenPack"
	"main/matrixPack"

	"fmt"
	"strconv"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

const FontPath = "./assets/Minecraft.ttf"
const FontSize = 40

func SetFont() *ttf.Font {
	font, err := ttf.OpenFont(FontPath, int(screenPack.VisinaProzora)/FontSize)
	if err != nil {
		panic(err)
	}

	return font
}

func TextMaker(font *ttf.Font, surface *sdl.Surface, matrica [][]mat.Cestica) (*sdl.Surface) {
	// ovo je za tekst
	var infoText = ""
	// PESAK
	if screenPack.KursorPoslednjiX < matrixPack.SirinaKan*matrixPack.BrPiksPoCestici {
		var poslednjiPiksel = matrica[screenPack.KursorPoslednjiX/matrixPack.BrPiksPoCestici][screenPack.KursorPoslednjiY/matrixPack.BrPiksPoCestici]
		infoText = mat.Ime[poslednjiPiksel.Materijal] + " @ " + fmt.Sprintf("%.2f", float32((-27315 + int32(poslednjiPiksel.Temperatura))/100)) + "C, SekMat: " + mat.Ime[poslednjiPiksel.SekMat] + ", Ticker: " + strconv.Itoa(int(poslednjiPiksel.Ticker))

		// UI
		// njanja: ovo se sigurno i ovde i u pritisku na gumb može mnogo lepše rešiti nekim funkcionalnim pristupom :DERP: longterm ali todo, vrv kad sređujem dva reda gumbeta
	} else {
		if screenPack.KursorPoslednjiY < (screenPack.VisinaUIMargine+screenPack.VisinaDugmeta)*int32(len(mat.Boja)-1) && screenPack.KursorPoslednjiY%(screenPack.VisinaUIMargine+screenPack.VisinaDugmeta) > screenPack.VisinaUIMargine {
			infoText = mat.Ime[mat.Materijal(screenPack.KursorPoslednjiY/(screenPack.VisinaUIMargine+screenPack.VisinaDugmeta))]
		}

		// PAUZA
		if screenPack.KursorPoslednjiY > screenPack.VisinaProzora-3*(screenPack.VisinaDugmeta+screenPack.VisinaUIMargine) && screenPack.KursorPoslednjiY < screenPack.VisinaProzora-3*(screenPack.VisinaDugmeta+screenPack.VisinaUIMargine)+screenPack.VisinaDugmeta {
			infoText = "Pause"
		}
		// SEJV
		if screenPack.KursorPoslednjiY > screenPack.VisinaProzora-2*(screenPack.VisinaDugmeta+screenPack.VisinaUIMargine) && screenPack.KursorPoslednjiY < screenPack.VisinaProzora-2*(screenPack.VisinaDugmeta+screenPack.VisinaUIMargine)+screenPack.VisinaDugmeta {
			infoText = "Save"
		}
		// RESET
		if screenPack.KursorPoslednjiY > screenPack.VisinaProzora-1*(screenPack.VisinaDugmeta+screenPack.VisinaUIMargine) && screenPack.KursorPoslednjiY < screenPack.VisinaProzora-1*(screenPack.VisinaDugmeta+screenPack.VisinaUIMargine)+screenPack.VisinaDugmeta {
			infoText = "Clear"
		}

	}

	// njanja todo napraviti fju koja ludi hex int za boju pretvara u rgba vrednost
	text, err := font.RenderUTF8Blended(infoText, sdl.Color{R: 255, G: 0, B: 0, A: 255})
	if err == nil {
		err = text.Blit(nil, surface, &sdl.Rect{X: 10, Y: 10, W: 0, H: 0})
		if err != nil {
			panic(err)
		}
	}

	return text
}