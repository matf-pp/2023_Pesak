package mat

import (
	"math"
)

type Materijal int

const (
	Zid    Materijal = 256
	Prazno Materijal = 0
	Pesak  Materijal = 1
	Voda   Materijal = 2
	Metal  Materijal = 3
	Kamen  Materijal = 4
	Lava   Materijal = 5
	Led    Materijal = 6
	Para   Materijal = 7
)

var Boja = map[Materijal]uint32{
	Zid:    0xffffff,
	Prazno: 0x000000,
	Pesak:  0xffff66,
	Voda:   0x3333ff,
	Metal:  0x33334b,
	Kamen:  0x666666,
	Lava:   0xff6600,
	Led:    0xccccff,
	Para:   0x9999cc,
}

var Gustina = map[Materijal]int32{
	Zid:    0,
	Prazno: 0,
	Pesak:  5,
	Voda:   3,
	Metal:  0,
	Kamen:  5,
	Lava:   4,
	Led:    0,
	Para:   -5,
}

// 0000 nece on nidje
// 0001 pada direkt
// 0010 pada dijagonalno
// 0100 curi horizontalno
var AStanje = map[Materijal]int{
	Zid:    0b0000,
	Prazno: 0b1111,
	Pesak:  0b0011,
	Voda:   0b0111,
	Metal:  0b0000,
	Kamen:  0b0001,
	Lava:   0b0111,
	Led:    0b0000,
	Para:   0b0111,
}

var FaznaPromena = map[Materijal][4]uint32{
	//TODO pretvoriti [4]int u struct{Materijal, int, int,Materijal}
	Zid:    {uint32(Zid), 0, math.MaxUint32, uint32(Zid)},
	Prazno: {uint32(Prazno), 0, math.MaxUint32, uint32(Zid)},
	Pesak:  {uint32(Pesak), 0, 1986, uint32(Lava)}, //1986
	Voda:   {uint32(Led), 273, 373, uint32(Para)},
	Metal:  {uint32(Metal), 0, 1811, uint32(Lava)}, //1811
	Kamen:  {uint32(Kamen), 0, 1473, uint32(Lava)}, //1473
	Lava:   {uint32(Lava), 1300, math.MaxUint32, uint32(Lava)},
	Led:    {uint32(Led), 0, 273, uint32(Voda)},
	Para:   {uint32(Voda), 373, math.MaxUint32, uint32(Para)},
}

var Zapaljiv = map[Materijal]bool{
	Zid:    false,
	Prazno: false,
	Pesak:  false,
	Voda:   false,
	Metal:  false,
	Kamen:  false,
	Lava:   false,
	Led:    false,
	Para:   false,
}

/*	ako pozelimo da ga bas objektno krkamo to bi izgledalo ovako nekako
type Materijali struct {
	boja uint32
	padajuc bool
	tecan bool
	tezina int
	zapaljiv bool
}

type Pesakk struct {
	Materijali
	temperatura uint32
}

func noviPesakk() Pesakk {
	zrno := Pesakk{
		Materijali : Materijali{
			boja : 0xffff66,
			padajuc : true,
			tecan : false,
			tezina : 5,
			zapaljiv : false,
			},
			temperatura : 293, //20 celzijusa
	}
	return zrno
}
*/

type Cestica struct {
	Materijal   Materijal
	Temperatura uint32
	SekMat      Materijal
}

func NewCestica(materijal Materijal) Cestica {
	zrno := Cestica{
		Materijal:   materijal,
		Temperatura: 293, //20 celzijusa
		SekMat:      Prazno,
	} //ovaj deo if materijal je uzasno ruzan -s
	if materijal == Led {
		zrno.SekMat = Voda
		zrno.Temperatura = 253 //-20
	}
	if materijal == Para {
		zrno.SekMat = Voda
		zrno.Temperatura = 383 //110
	}
	if materijal == Lava {
		zrno.SekMat = Kamen
		zrno.Temperatura = 1400
	}
	return zrno
}

// Update(matrix, bafer, i, j, matrix[i][j].Materijal)
// nisam u stanju da procenim trenutno treba li ovo prebaciti u treci fajl ili Boga pitaj gde -s
func Update(matrix [][]Cestica, bafer [][]Cestica, i int, j int) {

	if matrix[i][j].Materijal == Prazno {
		return
	}

	materijal := matrix[i][j].Materijal
	temperatura := matrix[i][j].Temperatura
	//	sekmat := matrix[i][j].SekMat

	//faza
	if materijal != Lava {
		if temperatura < FaznaPromena[materijal][1] {
			matrix[i][j].Materijal = Materijal(FaznaPromena[materijal][1])
			matrix[i][j].Temperatura = matrix[i][j].Temperatura
			matrix[i][j].SekMat = materijal
		} else if temperatura > FaznaPromena[materijal][2] {
			matrix[i][j].Materijal = Materijal(FaznaPromena[materijal][3])
			matrix[i][j].Temperatura = matrix[i][j].Temperatura
			matrix[i][j].SekMat = materijal
		}
	} else {
		if temperatura < FaznaPromena[materijal][1] {
			matrix[i][j].Materijal = matrix[i][j].SekMat
		}
	}

	//gorenje
	if Zapaljiv[materijal] {
		//TODO
	}

	//pomeranje ce morati biti redna petlja a ne paralelna sa ostalim efektima (zameni se pesak s naftom a tek onda nafta eksplodira i sta onda) -s
	//padanje
	pomeren := false
	padanje := AStanje[materijal]
	if (padanje & 0b0001) != 0 {
		gornji := matrix[i][j-1]
		gornjiMat := gornji.Materijal
		donji := matrix[i][j+1]
		donjiMat := donji.Materijal
		if AStanje[gornjiMat]&0b0001 != 0 && Gustina[gornjiMat] > Gustina[materijal] {
			// Molim? Kako ovo radi? /limun
			bafer[i][j].Materijal, bafer[i][j-1].Materijal = bafer[i][j].Materijal, bafer[i][j-1].Materijal
			bafer[i][j].Temperatura, bafer[i][j-1].Temperatura = bafer[i][j].Temperatura, bafer[i][j-1].Temperatura
			bafer[i][j].SekMat, bafer[i][j-1].SekMat = bafer[i][j].SekMat, bafer[i][j-1].SekMat
			// ovo ništa ne radi, bafer[i][j] = bafer[i][j-1] = bafer[i][j]... /limun
			//			bafer[i][j-1].Materijal = bafer[i][j].Materijal
			//			bafer[i][j-1].Temperatura = bafer[i][j].Temperatura
			//			bafer[i][j-1].SekMat = bafer[i][j].SekMat
			pomeren = true
		} else if AStanje[donjiMat] != 0b0000 && Gustina[donjiMat] < Gustina[materijal] {
			// isti problem ovde /limun
			bafer[i][j].Materijal, bafer[i][j+1].Materijal = bafer[i][j+1].Materijal, bafer[i][j].Materijal
			bafer[i][j].Temperatura, bafer[i][j+1].Temperatura = bafer[i][j+1].Temperatura, bafer[i][j].Temperatura
			bafer[i][j].SekMat, bafer[i][j+1].SekMat = bafer[i][j+1].SekMat, bafer[i][j].SekMat
			// isto, isto, isto /limun
			//			bafer[i][j+1].Materijal = bafer[i][j].Materijal
			//			bafer[i][j+1].Temperatura = bafer[i][j].Temperatura
			//			bafer[i][j+1].SekMat = bafer[i][j].SekMat
			pomeren = true
		}
	}

	if pomeren {

	}

}
