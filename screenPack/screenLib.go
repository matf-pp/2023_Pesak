package screenPack

import (
	"main/mat"
	"main/matrixPack"

	"github.com/fstanis/screenresolution"
	"github.com/veandco/go-sdl2/sdl"
)

var AutoFitScreen = true

const SirinaUIMargine = 10
const VisinaUIMargine = 10
const SirinaDugmeta = 40
const VisinaDugmeta = 20

// njanja: ovo ćemo da menjamo ako treba
var MarginaZaGumbad int32 = 2*SirinaUIMargine + SirinaDugmeta
var SirinaProzora = matrixPack.SirinaKan * matrixPack.BrPiksPoCestici + MarginaZaGumbad
var VisinaProzora = matrixPack.VisinaKan * matrixPack.BrPiksPoCestici

var KursorPoslednjiX = int32(0)
var KursorPoslednjiY = int32(0)

// var velicinaKursora int32 = 4
var VelicinaKursora int32 = 8
var MaxKursor int32 = 32

var TrenutniMat mat.Materijal = mat.Pesak


// takozvano dinamičko skaliranje ekrana ili nešto ne znam lupio sam
// ako ovo ikada u praksi izbaci nešto što ne staje u ekran javite mi da ga sredim ali mislim da je to besmislen posao
func FitToScreen(screenPercentage int) (int32, int32, int32) {
	resolution := screenresolution.GetPrimary()
	// adjustedScale := int32((float64(screenPercentage) / float64(100)) * float64(resolution.Height) / float64(matrixPack.VisinaKan))
	adjustedScale := int32((float64(screenPercentage) / float64(200)) * float64(resolution.Height) / float64(matrixPack.VisinaKan))

	return adjustedScale, matrixPack.SirinaKan * adjustedScale, matrixPack.VisinaKan * adjustedScale
}

func ProveriPritisakNaGumb(matrix, bafer [][]mat.Cestica, x, y int32) {
	//njanja: ovo je detekcija klika na gumb
	if x > SirinaProzora-MarginaZaGumbad+SirinaUIMargine && x < SirinaProzora-SirinaUIMargine {
		// njanja: TODO namestiti da se ređaju u više kolona ako baš mora //mora -s
		// materijali
		if y < (VisinaUIMargine+VisinaDugmeta)*int32(len(mat.Boja)-1) && y%(VisinaUIMargine+VisinaDugmeta) > VisinaUIMargine {
			TrenutniMat = mat.Materijal(y / (VisinaUIMargine + VisinaDugmeta))
		}

		// njanja: hardkodovan broj specijalnih dugmića
		// PAUZA
		if y > VisinaProzora-3*(VisinaDugmeta+VisinaUIMargine) && y < VisinaProzora-3*(VisinaDugmeta+VisinaUIMargine)+VisinaDugmeta {
			matrixPack.Pause = !matrixPack.Pause
		}
		// SEJV
		if y > VisinaProzora-2*(VisinaDugmeta+VisinaUIMargine) && y < VisinaProzora-2*(VisinaDugmeta+VisinaUIMargine)+VisinaDugmeta {
			SaveImage(matrix, int(matrixPack.BrPiksPoCestici))
			sdl.ShowSimpleMessageBox(sdl.MESSAGEBOX_INFORMATION, "pesak", "sačuvan B)", nil)
		}
		// njanja: nz je l ovo najpametniji način ali radi
		// RESET
		if y > VisinaProzora-1*(VisinaDugmeta+VisinaUIMargine) && y < VisinaProzora-1*(VisinaDugmeta+VisinaUIMargine)+VisinaDugmeta {
			for j := 0; j < matrixPack.VisinaKan; j++ {
				for i := 0; i < matrixPack.SirinaKan; i++ {
					matrix[i][j] = mat.NewCestica(mat.Prazno)
					bafer[i][j] = matrix[i][j]
				}
			}
			matrixPack.ZazidajMatricu(matrix)
		}
	}
}