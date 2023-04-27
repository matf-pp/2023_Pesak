package mat

import (
	//"main/brushPack"

	//	"fmt"
	//	"math"
	"math/rand"
)

type Materijal int

const tempRadius = 3
const usporenje = 1000

// const tezinaTempCestice = 1000
// const tezinaTempOkoline = 1
// const delilacTezina = tezinaTempCestice + tezinaTempOkoline

const (
	Prazno    Materijal = 0
	Metal     Materijal = 1
	Led       Materijal = 2
	Kamen     Materijal = 3
	Pesak     Materijal = 4
	Lava      Materijal = 5
	Voda      Materijal = 6
	Para      Materijal = 7
	TecniAzot Materijal = 8
	Plazma    Materijal = 9
	Zid       Materijal = 256
)

var Ime = map[Materijal]string{
	Prazno:    "Prazno",
	Metal:     "Metal",
	Led:       "Led",
	Kamen:     "Kamen",
	Pesak:     "Pesak",
	Lava:      "Lava",
	Voda:      "Voda",
	Para:      "Para",
	TecniAzot: "Tecni Azot",
	Plazma:    "Plazma",
	Zid:       "Zid",

}

var Boja = map[Materijal]uint32{
	Prazno:    0x000000,
	Metal:     0x33334b,
	Led:       0xaaaaff,
	Kamen:     0x999977,
	Pesak:     0xffff66,
	Lava:      0xff6600,
	Voda:      0x3333ff,
	Para:      0x6666ff,
	TecniAzot: 0x99ff99,
	Plazma:    0xff99ff,
	Zid:       0xffffff,
}

// mapa gustina neće raditi kako treba,
// ako nije preveliki problem hteo bih
// da bude tačna gustina * 10^5 /limun
var Gustina = map[Materijal]int32 {
	Prazno:		0,
	Metal:		0,
	Led:		0,
	Kamen:		6,
	Pesak:		5,
	Lava:		4,
	Voda:		3,
	Para:		-5,
	TecniAzot:	3,
	Plazma:		0,
	Zid:		0,
}

// ToplotnaProvodljivost
//moze li se ovo preimenovati u ToplodnaProvodljivost? da bude citkiji kod?
// "ToplotnaProvodljivost" zauzima mnogo mesta /limun
var Lambda = map[Materijal]int32 {
	Prazno: 26,			// 0,026
	Pesak:  2050,		// 2.05
	Voda:   600,		// 0,6
	Metal:  50200,		// 50.2
	Kamen:  288800,		// 288.8
	Lava:   1300000,	// 1300
	Led:    1600,		// 1,6
	Para:   16,			// 0.016
//
//	Prazno:
//	Metal:
//	Led:
//	Kamen:
//	Pesak:
//	Lava:
//	Voda:
//	Para:
//	TecniAzot:
//	Plazma:
//	Zid:
}

// 0000 nece on nidje
// ---1 pada direkt
// --1- pada dijagonalno
// -1-- curi horizontalno
// 1--- pomera se nasumicno svuda
var AStanje = map[Materijal]int{
	Prazno:		0b1111,
	Metal:		0b0000,
	Led:		0b0000,
	Kamen:		0b0001,
	Pesak:		0b0011,
	Lava:		0b0111,
	Voda:		0b0111,
	Para:		0b0111,
	TecniAzot:	0b0111,
	Plazma:		0b1111,
	Zid:		0b0000,
}

type FaznaPromena struct {
	Nize           Materijal
	Vise           Materijal
	TackaTopljenja uint32
	TackaKljucanja uint32
}

var MapaFaza = map[Materijal]FaznaPromena{

	//k(c) = c+273.15
	//c(k) = k–273.15

	//100 int32		=	1.00k
	//130000 int32	=	1300.00k

	//MinTemp = 0.00k = -273.15c = int32(-27315)
	//maxtemp = 8000.00c = int32(800000)

//	materijali	{nize,	Vise,	TackaT,		TackaK}
	Prazno:		{Prazno, Prazno, MinTemp, MaxTemp},
	Metal:		{Metal, Lava, MinTemp, 177315}, //1500.00c
	Led:		{Led, Voda, MinTemp, 27315}, //0.00c
	Kamen:		{Kamen, Lava, MinTemp, 157315}, //1300.00c
	Pesak:		{Pesak, Lava, MinTemp, 197315}, //1700.00c
	Lava:		{Lava, Lava, MinTemp, MaxTemp},
	Voda:		{Led, Para, 27315, 37315}, //0.00c, 100.00c
	Para:		{Voda, Para, 37315, MaxTemp}, //100.00c
	TecniAzot:	{TecniAzot, Prazno, MinTemp, 7315}, //-200.00c
	Plazma:		{Prazno, Plazma, 650000, MaxTemp}, //6773.15c
	Zid:		{Zid, Zid, MinTemp, MaxTemp},
}

const MinTemp uint32 = 0 // 0.00k
const MaxTemp uint32 = 827315 //8000.00c

var Zapaljiv = map[Materijal]bool{
	Prazno:		false,
	Metal:		false,
	Led:		false,
	Kamen:		false,
	Pesak:		false,
	Lava:		false,
	Voda:		false,
	Para:		false,
	TecniAzot:	false,
	Plazma:		false,
	Zid:		false,
}

type Cestica struct {
	Materijal   Materijal
	Temperatura uint32
	SekMat      Materijal
	Ticker      int8
}

func NewCestica(materijal Materijal) Cestica {
	zrno := Cestica{
		Materijal:   materijal,
		Temperatura: 29315, //20.00c
		SekMat:      Prazno,
		Ticker:      8, //za rdju gorivo itd, opada po principu nuklearnog raspada (svaki frejm ima x% sanse da ga dekrementira, na 0 prelazi u drugo stanje)
	}
	if materijal == Led {
		zrno.SekMat = Voda
		zrno.Temperatura = 24315 //-30.00c
	}
	if materijal == Para {
		zrno.SekMat = Voda
		zrno.Temperatura = 42315 //150.00c
	}
	if materijal == Lava {
		zrno.SekMat = Kamen
		zrno.Temperatura = 227315 //2000.00c
	}
	if materijal == TecniAzot {
		zrno.Temperatura = 2315 //-250.00c
	}
	if materijal == Plazma {
		zrno.Temperatura = 727315 //7000.00c
	}
	return zrno
}

// Update(matrix, bafer, i, j, matrix[i][j].Materijal)
// nisam u stanju da procenim trenutno treba li ovo prebaciti u treci fajl ili Boga pitaj gde -s
func UpdateTemp(matrix [][]Cestica, bafer [][]Cestica, i int, j int) {

	if matrix[i][j].Materijal == Prazno || matrix[i][j].Materijal == Zid {
		bafer[i][j].Temperatura = 29315
		return
	}
	trenutna := matrix[i][j]

	// temperatura
	//da ovo tvoje je kul ali prebrzo provodi da bih testirao kako se ponasa, ne zameri molim te -s
	/**
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
	matrix[i][j].Temperatura = temperatura
	bafer[i][j].Temperatura = temperatura

	if bafer[i][j].Temperatura < MinTemp {
		err := fmt.Sprintf("Temperatura cestice na poziciji [%d][%d] je van minimalne granice: %d < %d\n", i+k, j+l, bafer[i+k][j+l].Temperatura < MinTemp)
		panic(err)
	}
	if bafer[i][j].Temperatura > MaxTemp {
		err := fmt.Sprintf("Temperatura cestice na poziciji [%d][%d] je van maksimalne granice: %d > %d\n", i+k, j+l, bafer[i+k][j+l].Temperatura < MaxTemp)
		panic(err)
		//ako vas je dibagovanje dovelo ovde, moguce je da samo treba povecati MaxTemp ali razmislite o implikacijama, mozda ne valja racunanje temperature i negde krsimo termodinamiku
	}

	/**/
	/**/



	temperatura := float64(trenutna.Temperatura)
	parcePice := temperatura/9
	for k := -1; k < 2; k++ {
		for l := -1; l < 2; l++ {
			if matrix[i+k][j+l].Materijal != Prazno && matrix[i+k][j+l].Materijal != Zid {
				bafer[i+k][j+l].Temperatura += uint32(parcePice)
				temperatura = temperatura - parcePice
			}
		}
	}
	bafer[i][j].Temperatura += uint32(temperatura)

	/**/
	//interesantan fenomen - voda i para bi trebalo da se mimoidju zbog razlike u gustini, medjutim:
	//hladna voda koja pada se greje u kontaktu sa parom te isparava, menjajuci pravac kretanja ka gore
	//istovremeno topla para koja se dize u kontaktu sa hladnijom vodom se kondenzuje i pocinje da pada
	//ovime se postize taj efekat da se voda i para naizgled ne zamene,vec sudare na frontu i zaglave
	//medjutim one,naprotiv, zamene i poziciju i temperaturu. Verujem da ce se efekat izgubiti kada implementiram dijagonalni pad.
	//kada se prethodne dve linije zakomentarisu, voda i para, kao i svaki padajuci element sa parom,
	//ponasaju se normalno
	//-s
	// luka molim te obrisi ovaj blok komentara kad ga procitas, poenta je bila samo da ne mislis da je nesto izbagovano...
	// !!!!!!!!!!!!!!!!! ^^^^^
	// !!!!!!!!!!!!!!!!! |||||

}

func UpdatePhaseOfMatter(matrix [][]Cestica, bafer [][]Cestica, i int, j int) {

	if matrix[i][j].Materijal == Prazno || matrix[i][j].Materijal == Zid{
		return
	}

	trenutna := matrix[i][j]
	sekmat := trenutna.SekMat
	materijal := trenutna.Materijal
	temperatura := trenutna.Temperatura

	if materijal == Lava {
		if temperatura < MapaFaza[sekmat].TackaKljucanja {
			matrix[i][j].Materijal = sekmat
			bafer[i][j].Materijal = sekmat
		}
	} else {
		if temperatura < MapaFaza[materijal].TackaTopljenja {
			matrix[i][j].Materijal = MapaFaza[materijal].Nize
			bafer[i][j].Materijal = MapaFaza[materijal].Nize
		} else if temperatura > MapaFaza[materijal].TackaKljucanja {
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

}

func UpdatePosition(matrix [][]Cestica, bafer [][]Cestica, i int, j int) {
	//padanje

	if matrix[i][j].Materijal == Prazno || matrix[i][j].Materijal == Zid {
		return
	}

	trenutna := matrix[i][j]
	pomeren := false
	astanje := AStanje[trenutna.Materijal]
	smer := 0
	if Gustina[trenutna.Materijal] > 0 {
		smer = 1
	} else {
		smer = -1
	}
	//				{0, 1}		{0, 2}	{-1, 1}
	rFaktor := rand.Intn(2)*2 - 1

	if (astanje & 0b1000) != 0 {
		lRand := rand.Intn(3)-1
		rRand := rand.Intn(3)-1
		komsija := matrix[i+lRand][j+rRand]
		komsijaBuff := bafer[i+lRand][j+rRand]
		if komsija.Materijal == Prazno && komsijaBuff.Materijal == Prazno {
			bafer[i][j] = komsija
			bafer[i+lRand][j+rRand] = trenutna
			pomeren = true
		}
	}

	if pomeren {
		return
	}

	// dangerzone: start /limun
	if (astanje & 0b0001) != 0 {
		komsija := matrix[i][j+smer]
		//												( 1  *      G[v] = 2             <  1  *      g[ps] =  5) == True
		//                                              (-1  *      G[v] = 2             < -1  *      g[pr] = -5) == True
		if (AStanje[komsija.Materijal]&0b0001 != 0) && smer*int(Gustina[komsija.Materijal]) < smer*int(Gustina[trenutna.Materijal]) { ///ovde samo dodati || bafer[i][j+smer].Materijal == Prazno za blokovsko padanje, slicno u ostalim delovima ove f je
			if bafer[i][j+smer] == komsija {
				matrix[i][j+smer] = trenutna
				bafer[i][j+smer] = trenutna
				matrix[i][j] = komsija
				bafer[i][j] = komsija
				pomeren = true
			}
		}
	}
	// dangerzone: end /limun
	//ovo ne radi bas uvek a nmg da provalim sto i kako? iskreno mng bi mi znacilo da nemanja uradi da haverom preko cestice vidimo njene promenjive
	if pomeren {
		return
	}

	/**/
	if (astanje & 0b0010) != 0 {
		komsija1 := matrix[i+rFaktor][j+smer]
		if (AStanje[komsija1.Materijal]&0b0010 != 0) && smer*int(Gustina[komsija1.Materijal]) < smer*int(Gustina[trenutna.Materijal]) {
			if bafer[i+rFaktor][j+smer] == komsija1 {
				bafer[i+rFaktor][j+smer] = trenutna
				bafer[i][j] = komsija1
				pomeren = true
				return
			}
		}
		komsija2 := matrix[i-rFaktor][j+smer]
		if (AStanje[komsija2.Materijal]&0b0010 != 0) && smer*int(Gustina[komsija2.Materijal]) < smer*int(Gustina[trenutna.Materijal]) {
			if bafer[i-rFaktor][j+smer] == komsija2 {
				bafer[i-rFaktor][j+smer] = trenutna
				bafer[i][j] = komsija2
				pomeren = true
				return
			}
		}
	}
	/**/
	if (astanje & 0b0100) != 0 {

		if matrix[i+rFaktor][j].Materijal == Prazno && bafer[i+rFaktor][j].Materijal == Prazno {
			if matrix[i+rFaktor+rFaktor][j].Materijal == Prazno && bafer[i+rFaktor+rFaktor][j].Materijal == Prazno {
				bafer[i+rFaktor+rFaktor][j] = trenutna
				bafer[i][j] = matrix[i+rFaktor+rFaktor][j]
			} else {
				bafer[i+rFaktor][j] = trenutna
				bafer[i][j] = matrix[i+rFaktor][j]
			}
		} else if matrix[i-rFaktor][j].Materijal == Prazno && bafer[i-rFaktor][j].Materijal == Prazno {
			if matrix[i-rFaktor-rFaktor][j].Materijal == Prazno && bafer[i-rFaktor-rFaktor][j].Materijal == Prazno {
				bafer[i-rFaktor-rFaktor][j] = trenutna
				bafer[i][j] = matrix[i-rFaktor-rFaktor][j]
			} else {
				bafer[i-rFaktor][j] = trenutna
				bafer[i][j] = matrix[i-rFaktor][j]
			}
		}
		pomeren = true

	}

	if pomeren {
		return
	}

	if !pomeren {
		bafer[i][j] = trenutna
	}

}
