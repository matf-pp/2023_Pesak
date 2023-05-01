package fontPack

import (
	"main/src/mat"
	"main/src/matrixPack"
	"main/src/screenPack"

	"fmt"
	"strconv"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

const FontPath = "./res/fonts/Minecraft.ttf"
const FontSize = 40
const OutlineSize = 2

func SetFont() *ttf.Font {
	font, err := ttf.OpenFont(FontPath, int(screenPack.VisinaProzora)/FontSize)
	if err != nil {
		panic(err)
	}
	return font
}

func FontInit() *ttf.Font {
	var font *ttf.Font
	err := ttf.Init()
	if err != nil {
		panic(err)
	}
	return font
}

func TextMaker(font *ttf.Font, renderer *sdl.Renderer, matrica [][]mat.Cestica) {
	var infoText = ""
	// PESAK
	if screenPack.KursorPoslednjiX < matrixPack.SirinaKan*matrixPack.BrPiksPoCestici {
		var poslednjiPiksel = matrica[screenPack.KursorPoslednjiX/matrixPack.BrPiksPoCestici][screenPack.KursorPoslednjiY/matrixPack.BrPiksPoCestici]
		infoText = mat.Ime[poslednjiPiksel.Materijal] + " @ " + fmt.Sprintf("%.2f", float32((-27315+int32(poslednjiPiksel.Temperatura))/100)) + "C, SekMat: " + mat.Ime[poslednjiPiksel.SekMat] + ", Ticker: " + strconv.Itoa(int(poslednjiPiksel.Ticker))

		// UI
	} else {
		if screenPack.KursorPoslednjiY < (screenPack.VisinaUIMargine+screenPack.VisinaDugmeta)*int32(len(mat.Boja)-1) && screenPack.KursorPoslednjiY%(screenPack.VisinaUIMargine+screenPack.VisinaDugmeta) > screenPack.VisinaUIMargine {
			infoText = mat.Ime[mat.Materijal(screenPack.KursorPoslednjiY/(screenPack.VisinaUIMargine+screenPack.VisinaDugmeta))]
		}

		// PAUZA
		if screenPack.KursorPoslednjiY > screenPack.VisinaProzora-4*(screenPack.VisinaDugmeta+screenPack.VisinaUIMargine) && screenPack.KursorPoslednjiY < screenPack.VisinaProzora-4*(screenPack.VisinaDugmeta+screenPack.VisinaUIMargine)+screenPack.VisinaDugmeta {
			infoText = "Change Brush"
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

	font.SetOutline(2)
	text, err := font.RenderUTF8Blended(infoText, sdl.Color{R: 0, G: 0, B: 0, A: 255})
	if err == nil {
		texture, err := renderer.CreateTextureFromSurface(text)
		if err == nil {
			_, _, width, height, err := texture.Query()
			if err != nil {
				panic(err)
			}
			renderer.Copy(texture, nil, &sdl.Rect{X: 3*matrixPack.BrPiksPoCestici - OutlineSize/2, Y: 3*matrixPack.BrPiksPoCestici - OutlineSize/2, W: width, H: height})
		}
		defer texture.Destroy()
	}
	defer text.Free()

	font.SetOutline(0)
	text, err = font.RenderUTF8Blended(infoText, sdl.Color{R: 255, G: 255, B: 255, A: 255})
	if err == nil {
		texture, err := renderer.CreateTextureFromSurface(text)
		if err == nil {
			_, _, width, height, err := texture.Query()
			if err != nil {
				panic(err)
			}
			renderer.Copy(texture, nil, &sdl.Rect{X: 3 * matrixPack.BrPiksPoCestici, Y: 3 * matrixPack.BrPiksPoCestici, W: width, H: height})
		}
		defer texture.Destroy()
	}
	defer text.Free()
}
