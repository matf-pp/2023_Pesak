package brushPack

import (
	"main/mat"
	"main/matrixPack"
	"main/screenPack"
)

func Brush(matrix [][]mat.Cestica, bafer [][]mat.Cestica, x int32, y int32, state uint32) {
	//TODO za srednji klik da uzme materijal na koj mis trenutno pokazuje i postavi ga kao trenutni
	if x > matrixPack.SirinaKan * matrixPack.BrPiksPoCestici {
		return
	}
	if state != 1 && state != 2 && state != 4 {
		return
	}

	if state == 1 {
		for i := -screenPack.VelicinaKursora; i <= screenPack.VelicinaKursora; i++ {
			for j := -screenPack.VelicinaKursora; j <= screenPack.VelicinaKursora; j++ {
				tx, ty := matrixPack.ClampCoords(x/matrixPack.BrPiksPoCestici+i, y/matrixPack.BrPiksPoCestici+j)
				if matrix[tx][ty].Materijal == mat.Prazno || (screenPack.TrenutniMat == mat.Prazno && matrix[tx][ty].Materijal != mat.Zid) {
					matrix[tx][ty] = mat.NewCestica(screenPack.TrenutniMat)
					bafer[tx][ty] = mat.NewCestica(screenPack.TrenutniMat)
				}
			}
		}
	}

	if state == 2 {
		tx, ty := matrixPack.ClampCoords(x/matrixPack.BrPiksPoCestici, y/matrixPack.BrPiksPoCestici)
		screenPack.TrenutniMat = matrix[tx][ty].Materijal
	}

	if state == 4 {
		for i := -screenPack.VelicinaKursora; i <= screenPack.VelicinaKursora; i++ {
			for j := -screenPack.VelicinaKursora; j <= screenPack.VelicinaKursora; j++ {
				tx, ty := matrixPack.ClampCoords(x/matrixPack.BrPiksPoCestici+i, y/matrixPack.BrPiksPoCestici+j)
				if matrix[tx][ty].Materijal != mat.Zid {
					matrix[tx][ty] = mat.NewCestica(mat.Prazno)
					bafer[tx][ty] = mat.NewCestica(mat.Prazno)
				}
			}
		}
	}
}