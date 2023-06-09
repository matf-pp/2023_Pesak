// Package fontPack sluzi za podesavanje fonta i ispisivanje teksta na ekranu
package fontPack

import (
	"main/src/mat"
	"main/src/matrixPack"
	"main/src/screenPack"
	"main/src/languagePack"

	"fmt"
	"strconv"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

var fontPath = "./res/fonts/PixeloidSans-Bold.ttf"
var FontSize = 48
const outlineSize = 2

// SetFont ne prima nista; vraca font
func SetFont() *ttf.Font {
	font, err := ttf.OpenFont(fontPath, int(screenPack.VisinaProzora)/FontSize)
	if err != nil {
		panic(err)
	}
	return font
}

// FontInit nagadjam da inicijalizuje font
func FontInit() *ttf.Font {
	var font *ttf.Font
	err := ttf.Init()
	if err != nil {
		panic(err)
	}
	return font
}

// TextMaker prima font, renderer i matricu Cestica; ne vraca nista; ispisuje odgovarajuci tekst na ekranu
//ko je ovaj komentar napisao debilan je sta znaci "odgovarajuci tekst"
//f ja ispisuje info o cestici na kojoj je kursor, ukoliko je dibag modupaljen.
func TextMaker(font *ttf.Font, renderer *sdl.Renderer, matrica [][]mat.Cestica) {
	var infoText = ""
	// PESAK
	if mat.KursorPoslednjiX < matrixPack.SirinaKan*matrixPack.BrPiksPoCestici {
		var poslednjiPiksel = matrica[mat.KursorPoslednjiX/matrixPack.BrPiksPoCestici][mat.KursorPoslednjiY/matrixPack.BrPiksPoCestici]
		infoText = "M: " + mat.Ime[poslednjiPiksel.Materijal][mat.IzabraniJezik] + " @ " + fmt.Sprintf("%.2f", float32((-27315+int32(poslednjiPiksel.Temperatura))/100)) + "C, SM: " + mat.Ime[poslednjiPiksel.SekMat][mat.IzabraniJezik] + ", T: " + strconv.Itoa(int(poslednjiPiksel.Ticker))// + ", Speed: " + fmt.Sprintf("%.1f", float64(matrixPack.FpsCap)/60.0) + "x"

		// UI
	} else {
		if mat.KursorPoslednjiY < (screenPack.VisinaUIMargine+screenPack.VisinaDugmeta)*int32(len(mat.Boja)-1) && mat.KursorPoslednjiY%(screenPack.VisinaUIMargine+screenPack.VisinaDugmeta) > screenPack.VisinaUIMargine {
			visinaY := mat.KursorPoslednjiY / (screenPack.VisinaUIMargine + screenPack.VisinaDugmeta)
			someMat := mat.Materijal(visinaY)
			if len(mat.Ime[someMat]) > 0 {
				infoText = mat.Ime[someMat][mat.IzabraniJezik]
			}
		}
		// CHANGE BRUSH
		if mat.KursorPoslednjiY > screenPack.VisinaProzora-5*(screenPack.VisinaDugmeta+screenPack.VisinaUIMargine) && mat.KursorPoslednjiY < screenPack.VisinaProzora-5*(screenPack.VisinaDugmeta+screenPack.VisinaUIMargine)+screenPack.VisinaDugmeta {
			infoText = languagePack.BrzinaTekst[mat.IzabraniJezik] + fmt.Sprintf("%.1f", float64((matrixPack.FpsCap*2)%210)/60.0) + "x"
		}
		// CHANGE BRUSH
		if mat.KursorPoslednjiY > screenPack.VisinaProzora-4*(screenPack.VisinaDugmeta+screenPack.VisinaUIMargine) && mat.KursorPoslednjiY < screenPack.VisinaProzora-4*(screenPack.VisinaDugmeta+screenPack.VisinaUIMargine)+screenPack.VisinaDugmeta {
			infoText = languagePack.PromeniOblikCetkeTekst[mat.IzabraniJezik]
		}
		// PAUZA
		if mat.KursorPoslednjiY > screenPack.VisinaProzora-3*(screenPack.VisinaDugmeta+screenPack.VisinaUIMargine) && mat.KursorPoslednjiY < screenPack.VisinaProzora-3*(screenPack.VisinaDugmeta+screenPack.VisinaUIMargine)+screenPack.VisinaDugmeta {
			infoText = languagePack.ZaustaviTekst[mat.IzabraniJezik]
		}
		// SEJV
		if mat.KursorPoslednjiY > screenPack.VisinaProzora-2*(screenPack.VisinaDugmeta+screenPack.VisinaUIMargine) && mat.KursorPoslednjiY < screenPack.VisinaProzora-2*(screenPack.VisinaDugmeta+screenPack.VisinaUIMargine)+screenPack.VisinaDugmeta {
			infoText = languagePack.SacuvajTekst[mat.IzabraniJezik]
		}
		// RESET
		if mat.KursorPoslednjiY > screenPack.VisinaProzora-1*(screenPack.VisinaDugmeta+screenPack.VisinaUIMargine) && mat.KursorPoslednjiY < screenPack.VisinaProzora-1*(screenPack.VisinaDugmeta+screenPack.VisinaUIMargine)+screenPack.VisinaDugmeta {
			infoText = languagePack.OcistiTekst[mat.IzabraniJezik]
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
			renderer.Copy(texture, nil, &sdl.Rect{X: 3*matrixPack.BrPiksPoCestici - outlineSize/2, Y: 3*matrixPack.BrPiksPoCestici - outlineSize/2, W: width, H: height})
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
