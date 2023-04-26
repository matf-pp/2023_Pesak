package mat

import (
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
// treba mi bolje rešenje, npr da paket za pravljenje kanvasa sadrži širinu i visinu i f-je za kanvas /limun
// mogu li se ovi kom sada obrisati ako su suvisni? malo je sizofreno ne pratim vise sta je sta
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

var Ime = map[Materijal]string{
	Zid:    "Zid",
	Prazno: "Prazno",
	Pesak:  "Pesak",
	Voda:   "Voda",
	Metal:  "Metal",
	Kamen:  "Kamen",
	Lava:   "Lava",
	Led:    "Led",
	Para:   "Para",
}

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
	Zid:	0,
	Prazno: 1225, 	// 0.01225 			0
	Pesak:  163100, // 1.631 			5
	Voda:   100000, // 1 				3
	Metal:  786000, // 7.860 	(čelik) 0
	Kamen:  260000, // 2.600 			5
	Lava:   310000, // 3.100 			4
	Led:    91700, 	// 0.917 			0
	Para:   598, 	// 0.005.98	   	   -5
}

// ToplotnaProvodljivost
var Lambda = map[Materijal]int32 {
	Prazno: 26,			// 0,026
	Pesak:  2050,		// 2.05
	Voda:   600,		// 0,6
	Metal:  50200,		// 50.2
	Kamen:  288800,		// 288.8
	Lava:   1300000,	// 1300
	Led:    1600,		// 1,6
	Para:   16,			// 0.016
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
	Nize           Materijal
	Vise           Materijal
	TackaTopljenja int32
	TackaKljucanja int32
}

var MapaFaza = map[Materijal]FaznaPromena{

	//k(c) = c+273.15
	//c(k) = k–273.15

	//100 int32		=	1.00c
	//130000 int32	=	1300.00c

	//MinTemp = 0.00k = -273.15c = int32(-27315)
	//maxtemp = 8000.00c = int32(800000)

	//	materijali	{nize,	Vise,	TackaT,		TackaK}
	Zid:    {Zid, Zid, MinTemp, MaxTemp},
	Prazno: {Prazno, Prazno, MinTemp, MaxTemp},
	Pesak:  {Pesak, Lava, MinTemp, 170000},
	Voda:   {Led, Para, 0, 10000},
	Metal:  {Metal, Lava, MinTemp, 150000},
	Kamen:  {Kamen, Lava, MinTemp, 130000},
	Lava:   {Lava, Lava, MinTemp, MaxTemp},
	Led:    {Led, Voda, MinTemp, 0},
	Para:   {Voda, Para, 10000, MaxTemp},
}

const MinTemp int32 = -27315
const MaxTemp int32 = 800000

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
	Temperatura int32
	SekMat      Materijal
	Ticker      int8
}

func NewCestica(materijal Materijal) Cestica {
	zrno := Cestica{
		Materijal:   materijal,
		Temperatura: 2000,
		SekMat:      Prazno,
		Ticker:      8, //za rdju gorivo itd, opada po principu nuklearnog raspada (svaki frejm ima x% sanse da ga dekrementira, na 0 prelazi u drugo stanje)
	}
	if materijal == Voda {
		zrno.Temperatura = 2000
	}
	if materijal == Led {
		zrno.SekMat = Voda
		zrno.Temperatura = -3000
	}
	if materijal == Para {
		zrno.SekMat = Voda
		zrno.Temperatura = 15000
	}
	if materijal == Lava {
		zrno.SekMat = Kamen
		zrno.Temperatura = 200000
	}
	return zrno
}

// Update(matrix, bafer, i, j, matrix[i][j].Materijal)
// nisam u stanju da procenim trenutno treba li ovo prebaciti u treci fajl ili Boga pitaj gde -s
func UpdateTemp(matrix [][]Cestica, bafer [][]Cestica, i int, j int) {

	if matrix[i][j].Materijal == Prazno || matrix[i][j].Materijal == Zid {
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

	var brojacKomsija int32 = 0
	for k := -1; k < 2; k++ {
		for l := -1; l < 2; l++ {
			if matrix[i+k][j+l].Materijal != Prazno && matrix[i+k][j+l].Materijal != Zid {
				brojacKomsija++
			}
		}
	}
	deliTemp := trenutna.Temperatura / brojacKomsija
	for k := -1; k < 2; k++ {
		for l := -1; l < 2; l++ {
			if matrix[i+k][j+l].Materijal != Prazno && matrix[i+k][j+l].Materijal != Zid {
				bafer[i+k][j+l].Temperatura += deliTemp
				if bafer[i+k][j+l].Temperatura < MinTemp {
					//					err := fmt.Sprintf("Temperatura cestice na poziciji [%d][%d] je van minimalne granice: %d \< %d\n", i+k, j+l, bafer[i+k][j+l].Temperatura < MinTemp)
					//					panic(err)
				}
				if bafer[i+k][j+l].Temperatura > MaxTemp {
					//					err := fmt.Sprintf("Temperatura cestice na poziciji [%d][%d] je van maksimalne granice: %d \> %d\n", i+k, j+l, bafer[i+k][j+l].Temperatura < MaxTemp)
					//					panic(err)
					//ako vas je dibagovanje dovelo ovde, moguce je da samo treba povecati MaxTemp ali razmislite o implikacijama, mozda ne valja racunanje temperature i negde krsimo termodinamiku
				}
			}
		}
	}

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

	if matrix[i][j].Materijal == Prazno {
		return
	}
	if matrix[i][j].Materijal == Zid {
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
	//pomeranje ce morati biti redna petlja a ne paralelna sa ostalim efektima (zameni se pesak s naftom a tek onda nafta eksplodira i sta onda) -s
	//decko ti bulaznis -s
	//padanje

	if matrix[i][j].Materijal == Prazno || matrix[i][j].Materijal == Zid {
		return
	}

	trenutna := matrix[i][j]
	pomeren := false
	astanje := AStanje[trenutna.Materijal]
	smer := 0
	if Gustina[trenutna.Materijal] > 1225 {
		smer = 1
	} else {
		smer = -1
	}
	//				{0, 1}		{0, 2}	{-1, 1}
	rFaktor := rand.Intn(2)*2 - 1

	if (astanje & 0b0001) != 0 {
		komsija := matrix[i][j+smer]
		//												( 1  *      G[v] = 2             <  1  *      g[ps] =  5) == True
		//                                              (-1  *      G[v] = 2             < -1  *      g[pr] = -5) == True
		if (AStanje[komsija.Materijal]&0b0001 != 0) && smer*int(Gustina[komsija.Materijal]) < smer*int(Gustina[trenutna.Materijal]) { ///ovde samo dodati || bafer[i][j+smer].Materijal == Prazno za blokovsko padanje, slicno u ostalim delovima ove f je
			if bafer[i][j+smer] == komsija {
				bafer[i][j+smer] = trenutna
				bafer[i][j] = komsija
				pomeren = true
			}
		}
	} //ovo ne radi bas uvek a nmg da provalim sto i kako? iskreno mng bi mi znacilo da nemanja uradi da haverom preko cestice vidimo njene promenjive
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

	if !pomeren {
		bafer[i][j] = trenutna
	}

}
