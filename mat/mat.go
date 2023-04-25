package mat

import (
//	"fmt"
	"math"
	"math/rand"
)

type Materijal int

const tempRadius = 3
const usporenje = 1000
// const tezinaTempCestice = 1000
// const tezinaTempOkoline = 1
// const delilacTezina = tezinaTempCestice + tezinaTempOkoline
// treba mi bolje rešenje, npr da paket za pravljenje kanvasa sadrži širinu i visinu i f-je za kanvas /limun
const sirinaKanvasa, visinaKanvasa = 240, 144

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
	Kamen:  0x999977,
	Lava:   0xff6600,
	Led:    0xccccff,
	Para:   0x9999cc,
}

var Gustina = map[Materijal]int32 {
	Zid:    0,
	Prazno: 1220, 	// 0.01225 						0
	Pesak:  163100, // 1.631 						5
	Voda:   100000, // 1 							3
	Metal:  786000, // 7.860 						(čelik) 0
	Kamen:  260000, // 2.600 						5
	Lava:   310000, // 3.100 						4
	Led:    91700, 	// 0.917 						0
	Para:   598, 	// 0 (5.98 × 10^(–4) g cm^(–3)) -5
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

type FaznaPromena struct {
	Nize Materijal
	Vise Materijal
	TackaTopljenja float64
	TackaKljucanja float64
}

var MapaFaza = map[Materijal]FaznaPromena{
	// sto vise razmisljam o tome to bih vise insistirao da temp bude u neoznacenim intidzerima tj kelvinima, makar jedan kelvin bio uint32(10) ili uint32(100) radi preciznosti -s
	Zid:    {Zid, Zid, -math.MaxFloat64, math.MaxFloat64},
	Prazno: {Prazno, Prazno, -math.MaxFloat64, math.MaxFloat64},
	Pesak:  {Pesak, Lava, -math.MaxFloat64, 1700},
	Voda:   {Led, Para, 0, 100},
	Metal:  {Metal, Lava, -math.MaxFloat64, 1500},
	Kamen:  {Kamen, Lava, -math.MaxFloat64, 1300},
	Lava:   {Lava, Lava, math.MaxFloat64, math.MaxFloat64},
	Led:    {Led, Voda, -math.MaxFloat64, 0},
	Para:   {Voda, Para, 100, math.MaxFloat64},
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


type Cestica struct {
	Materijal   Materijal
	Temperatura float64
	SekMat      Materijal
	Ticker		int8
}

func NewCestica(materijal Materijal) Cestica {
	zrno := Cestica{
		Materijal:   materijal,
		Temperatura: 40,
		SekMat:      Prazno,
		Ticker:		 8,	//za rdju gorivo itd, opada po principu nuklearnog raspada (svaki frejm ima x% sanse da ga dekrementira, na 0 prelazi u drugo stanje)
	}
	if materijal == Voda {
		zrno.Temperatura = 10
	}
	if materijal == Led {
		zrno.SekMat = Voda
		zrno.Temperatura = -30
	}
	if materijal == Para {
		zrno.SekMat = Voda
		zrno.Temperatura = 500
	}
	if materijal == Lava {
		zrno.SekMat = Kamen
		zrno.Temperatura = 5000
	}
	return zrno
}

// Update(matrix, bafer, i, j, matrix[i][j].Materijal)
// nisam u stanju da procenim trenutno treba li ovo prebaciti u treci fajl ili Boga pitaj gde -s
func Update(matrix [][]Cestica, bafer [][]Cestica, i int, j int) {
	trenutna := matrix[i][j]

	// temperatura//da ovo tvoje je kul ali prebrzo provodi da bih testirao kako se ponasa, nezameri molim te -s
	temperatura := trenutna.Temperatura
	for k := int(math.Max(float64(i-tempRadius), 0)); k < int(math.Min(float64(i+tempRadius), sirinaKanvasa)); k++ {
		for l := int(math.Max(float64(j-tempRadius), 0)); l < int(math.Min(float64(j+tempRadius), visinaKanvasa)); l++ {
			if math.Abs(float64(k-i)) + math.Abs(float64(l-j)) >= tempRadius || temperatura == matrix[k][l].Temperatura {
				continue
			}
			gusTrenutna := float64(Gustina[trenutna.Materijal])
			gusSused := float64(Gustina[matrix[k][l].Materijal])
			ostatak := (gusTrenutna * temperatura + gusSused * matrix[k][l].Temperatura)/(gusTrenutna + gusSused) - temperatura
			temperatura += ostatak/usporenje
			matrix[k][l].Temperatura -= ostatak/usporenje
			bafer[k][l].Temperatura -= ostatak/usporenje
		}
	}
	// novaTemperatura := float64(0)
	// brojac := 0
	// for k := -1; k < 2; k++ {
	// 	for l := -1; l < 2; l++ {
	// 		if matrix[i+k][j+l].Materijal != Prazno && matrix[i+k][j+l].Materijal != Zid{
	// 			novaTemperatura += matrix[i+k][j+l].Temperatura
	// 			brojac++
	// 		}
	// 	}
	// }
	// novaTemperatura = novaTemperatura / float64(brojac)
	matrix[i][j].Temperatura = temperatura
	bafer[i][j].Temperatura = temperatura
	//interesantan fenomen - voda i para bi trebalo da se mimoidju zbog razlike u gustini, medjutim:
	//hladna voda koja pada se greje u kontaktu sa parom te isparava, menjajuci pravac kretanja ka gore
	//istovremeno topla para koja se dize u kontaktu sa hladnijom vodom se kondenzuje i pocinje da pada
	//ovime se postize taj efekat da se voda i para naizgled ne zamene,vec sudare na frontu i zaglave
	//medjutim one,naprotiv, zamene i poziciju i temperaturu. Verujem da ce se efekat izgubiti kada implementiram dijagonalni pad.
	//kada se prethodne dve linije zakomentarisu, voda i para, kao i svaki padajuci element sa parom,
	//ponasaju se normalno
	//-s

	if matrix[i][j].Materijal == Prazno {
		return
	}


	//faza
	sekmat := trenutna.SekMat
	materijal := trenutna.Materijal
	temperaturaZaFaze := trenutna.Temperatura
	if materijal == Lava {
		if temperaturaZaFaze < MapaFaza[sekmat].TackaKljucanja{
			matrix[i][j].Materijal = sekmat
			bafer[i][j].Materijal = sekmat
		}
	} else {
		if temperaturaZaFaze < MapaFaza[materijal].TackaTopljenja {
			matrix[i][j].Materijal = MapaFaza[materijal].Nize
			bafer[i][j].Materijal = MapaFaza[materijal].Nize
		} else if temperaturaZaFaze > MapaFaza[materijal].TackaKljucanja {
			matrix[i][j].Materijal = MapaFaza[materijal].Vise
			matrix[i][j].SekMat = materijal
			bafer[i][j].Materijal = MapaFaza[materijal].Vise
			bafer[i][j].SekMat = materijal
		}		
	}

	//gorenje
	if Zapaljiv[materijal] {
		//TODO
	}

	//pomeranje ce morati biti redna petlja a ne paralelna sa ostalim efektima (zameni se pesak s naftom a tek onda nafta eksplodira i sta onda) -s
	//decko ti bulaznis -s
	//padanje
	pomeren := false
	astanje := AStanje[materijal]
	if (astanje & 0b0001) != 0 {
		gornji := matrix[i][j-1]
		donji := matrix[i][j+1]
		if (AStanje[gornji.Materijal] & 0b0001 != 0) && Gustina[gornji.Materijal] > Gustina[materijal] {
			if bafer[i][j-1] == matrix[i][j-1]{
				bafer[i][j], bafer[i][j-1] = matrix[i][j-1], matrix[i][j]
				bafer[i][j].Temperatura, bafer[i][j-1].Temperatura = matrix[i][j-1].Temperatura, matrix[i][j].Temperatura
				pomeren = true
			}
		} else if (AStanje[donji.Materijal] != 0b0000) && Gustina[donji.Materijal] < Gustina[materijal] {
			if bafer[i][j+1] == matrix[i][j+1]{	
				bafer[i][j], bafer[i][j+1] = matrix[i][j+1], matrix[i][j]
				bafer[i][j].Temperatura, bafer[i][j+1].Temperatura = matrix[i][j+1].Temperatura, matrix[i][j].Temperatura
				pomeren = true
			}
		}
	}
	if pomeren {
		return
	}/**/
	// napiši šta ovo radi /limun
	if (astanje & 0b0010) != 0 {
		// ?? /limun
		sgn := rand.Intn(2)-1
//		fmt.Printf("%d ", sgn)				//!!!!!!???????
		gd := matrix[i+sgn][j-1]
		gl := matrix[i-sgn][j-1]
		dd := matrix[i+sgn][j+1]
		dl := matrix[i-sgn][j+1]
		if AStanje[gd.Materijal]&0b0010 != 0 && Gustina[gd.Materijal] > Gustina[materijal] && bafer[i+sgn][j-1] == matrix[i+sgn][j-1] { //pakao kolika linija braco moja u Hristu
			bafer[i][j], bafer[i+sgn][j-1] = matrix[i+sgn][j-1], matrix[i][j]
			bafer[i][j].Temperatura, bafer[i+sgn][j-1].Temperatura = matrix[i+sgn][j-1].Temperatura, matrix[i][j].Temperatura
			pomeren = true
		} else if AStanje[gl.Materijal]&0b0010 != 0 && Gustina[gl.Materijal] > Gustina[materijal] && bafer[i-sgn][j-1] == matrix[i-sgn][j-1] {
			bafer[i][j], bafer[i-sgn][j-1] = matrix[i-sgn][j-1], matrix[i][j]
			bafer[i][j].Temperatura, bafer[i-sgn][j-1].Temperatura = matrix[i-sgn][j-1].Temperatura, matrix[i][j].Temperatura
			pomeren = true
		} else if AStanje[dd.Materijal]&0b0010 != 0 && Gustina[dd.Materijal] < Gustina[materijal] && bafer[i+sgn][j+1] == matrix[i+sgn][j+1] { //pakao 2 braco moja mucena
			bafer[i][j], bafer[i+sgn][j+1] = matrix[i+sgn][j+1], matrix[i][j]
			bafer[i][j].Temperatura, bafer[i+sgn][j+1].Temperatura = matrix[i+sgn][j+1].Temperatura, matrix[i][j].Temperatura
			pomeren = true
		} else if AStanje[dl.Materijal]&0b0010 != 0 && Gustina[dl.Materijal] < Gustina[materijal] && bafer[i-sgn][j+1] == matrix[i-sgn][j+1] {
			bafer[i][j], bafer[i-sgn][j+1] = matrix[i-sgn][j+1], matrix[i][j]
			bafer[i][j].Temperatura, bafer[i-sgn][j+1].Temperatura = matrix[i-sgn][j+1].Temperatura, matrix[i][j].Temperatura
			pomeren = true
		}

	}/**///ovaj deo ima neki bajas na levo tj desno nmg to veceras

	if pomeren {
		return
	}

/**
	if (astanje & 0b0100) != 0 {
		sgnn := rand.Intn(2)*2-1
		desni := matrix[i+sgnn][j]
		levi := matrix[i-sgnn][j]
		if AStanje[desni.Materijal]&0b0100 != 0 && bafer[i+sgnn][j] == matrix[i+sgnn][j] {
			bafer[i][j], bafer[i+sgnn][j] = matrix[i+sgnn][j], matrix[i][j]
		} else if AStanje[levi.Materijal]&0b0100 != 0 && bafer[i-sgnn][j] == matrix[i-sgnn][j] {
			bafer[i][j], bafer[i-sgnn][j] = matrix[i-sgnn][j], matrix[i][j]			
		}
	}/**/
	//ceo ovaj deo ima neki prokleti ne znam ni ja ta ako vas ne mrzi bacite pogled ako vas mrzi nemojte, imate druge stvari koje mozete raditi -s


}
