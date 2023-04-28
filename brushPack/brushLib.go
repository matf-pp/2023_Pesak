package brushPack

import (
	"main/mat"
	"main/matrixPack"
	"main/screenPack"
)

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
		for i := -screenPack.VelicinaKursora; i <= screenPack.VelicinaKursora; i++ {
			for j := -screenPack.VelicinaKursora; j <= screenPack.VelicinaKursora; j++ {
				tx, ty := matrixPack.ClampCoords(x/matrixPack.BrPiksPoCestici+i, y/matrixPack.BrPiksPoCestici+j)
				if screenPack.TrenutniMat != mat.Toplo && screenPack.TrenutniMat != mat.Hladno{
					if matrix[tx][ty].Materijal == mat.Prazno || (screenPack.TrenutniMat == mat.Prazno && matrix[tx][ty].Materijal != mat.Zid) {
						matrix[tx][ty] = mat.NewCestica(screenPack.TrenutniMat)
					}
				} else {
					if screenPack.TrenutniMat == mat.Toplo && matrix[tx][ty].Materijal != mat.Prazno && matrix[tx][ty].Materijal != mat.Zid {
						matrix[tx][ty].Temperatura += 100
					} else if screenPack.TrenutniMat == mat.Hladno && matrix[tx][ty].Materijal != mat.Prazno && matrix[tx][ty].Materijal != mat.Zid {
						matrix[tx][ty].Temperatura -= 100
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
		for i := -screenPack.VelicinaKursora; i <= screenPack.VelicinaKursora; i++ {
			for j := -screenPack.VelicinaKursora; j <= screenPack.VelicinaKursora; j++ {
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