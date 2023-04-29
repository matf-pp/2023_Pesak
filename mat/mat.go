package mat

import (
	//"fmt"
	//"math"
	"math/rand"
)

type Materijal int

const (
	Prazno    Materijal = 0
	Metal     Materijal = 1
	Led       Materijal = 2
	Kamen     Materijal = 3
	Drvo      Materijal = 4
	Sljunak   Materijal = 5
	Pesak     Materijal = 6
	So        Materijal = 7
	Rdja      Materijal = 254
	Lava      Materijal = 8
	Voda      Materijal = 9
	SlanaVoda Materijal = 255
	Para      Materijal = 10
	Vatra     Materijal = 11
	TecniAzot Materijal = 12
	Plazma    Materijal = 13
	Toplo     Materijal = 14
	Hladno    Materijal = 15
	Zid       Materijal = 256
)

var Ime = map[Materijal]string{
	Prazno:    "Prazno",
	Metal:     "Metal",
	Led:       "Led",
	Kamen:     "Kamen",
	Drvo:      "Drvo",
	Sljunak:   "Sljunak",
	Pesak:     "Pesak",
	So:        "So",
	Rdja:      "Rdja",
	Lava:      "Lava",
	Voda:      "Voda",
	SlanaVoda: "SlanaVoda",
	Para:      "Para",
	Vatra:     "Vatra",
	TecniAzot: "Tecni Azot",
	Plazma:    "Plazma",
	Toplo:     "Toplo",
	Hladno:    "Hladno",
	Zid:       "Zid",
}

var Boja = map[Materijal]uint32{
	Prazno:    0x000000,
	Metal:     0x33334b,
	Led:       0xaaaaff,
	Kamen:     0x999988,
	Drvo:      0x994400,
	Sljunak:   0x888877,
	Pesak:     0xffff66,
	So:        0xeeeeee,
	Rdja:      0x6f0f2b,
	Lava:      0xff6600,
	Voda:      0x3333ff,
	SlanaVoda: 0x4444ff,
	Para:      0x6666ff,
	Vatra:     0xd73502,
	TecniAzot: 0x99ff99,
	Plazma:    0xff99ff,
	Toplo:     0xff0000,
	Hladno:    0x00ffff,
	Zid:       0xffffff,
}

// pls samo ne dirajte ovo a ako zelite nesto drugo nazvati gustina promenite naziv ovoga i sve njegove pojave u kodu, hvala -s
var Gustina = map[Materijal]int32{
	Prazno:    0,
	Metal:     0,
	Led:       0,
	Kamen:     0,
	Drvo:      0,
	Sljunak:   5,
	Pesak:     5,
	So:        5,
	Rdja:      5,
	Lava:      4,
	Voda:      2,
	SlanaVoda: 3,
	Para:      -5,
	Vatra:     -5,
	TecniAzot: 3,
	Plazma:    0,
	Zid:       0,
}
var GustinaBoja = map[Materijal]uint32 {
	Prazno: 	0xc8c8c8,
	Metal:  	0x00ff00,
	Led:    	0x004600,
	Kamen:  	0x00b400,
	Pesak:  	0x007800,
	So:        	0x00a000,
	Rdja:      	0x00ff00,
	Lava:   	0x00c800,
	Voda:   	0x005000,
	SlanaVoda: 	0x005a00,
	Para:   	0xc800c8,
	TecniAzot: 	0x006400,
	Plazma:    	0xff00ff,
	Zid:		0,
}

// ToplotnaProvodljivost
var Lambda = map[Materijal]uint64{
	Prazno: 	26,      // 0.026
	Metal:  	50200,   // 50.2
	Led:    	1600,    // 1.6
	Kamen:  	288800,  // 288.8
	Pesak:  	2050,    // 2.05
	So:        	6000,	 // 6
	Rdja:      	50200,   // 50.2
	Lava:   	1300000, // 1300
	Voda:   	600,     // 0.6
	SlanaVoda:  600,	 // 0.6
	Para:   	16,      // 0.016
	TecniAzot:  25,		 // 0.025
	Plazma:    	1500000, // lupio sam broj zato što https://adsabs.harvard.edu/full/1962SvA.....5..495I /limun
	Zid:       	0,
}

// 0000 nece on nidje
// ---1 pada direkt
// --1- pada dijagonalno
// -1-- curi horizontalno
// 1--- pomera se nasumicno svuda
var AStanje = map[Materijal]int{
	Prazno:    0b1111,
	Metal:     0b0000,
	Led:       0b0000,
	Kamen:     0b0000,
	Drvo:      0b0000,
	Sljunak:   0b0001,
	Pesak:     0b0011,
	So:        0b0011,
	Rdja:      0b0011,
	Lava:      0b0111,
	Voda:      0b0111,
	SlanaVoda: 0b0111,
	Para:      0b0111,
	Vatra:     0b0111,
	TecniAzot: 0b0111,
	Plazma:    0b1111,
	Zid:       0b0000,
}

type FaznaPromena struct {
	Nize           Materijal
	Vise           Materijal
	TackaTopljenja uint64
	TackaKljucanja uint64
}

var MapaFaza = map[Materijal]FaznaPromena{

	//k(c) = c+273.15
	//c(k) = k–273.15

	//100 int32		=	1.00k
	//130000 int32	=	1300.00k

	//MinTemp = 0.00k = -273.15c = int32(-27315)
	//maxtemp = 8000.00c = int32(800000)

	//	materijali	{Nize,	Vise,	TackaT,		TackaK}
	Prazno:    {TecniAzot, Plazma, 5315, 627315},
	Metal:     {Metal, Lava, MinTemp, 177315}, //1500.00c
	Led:       {Led, Voda, MinTemp, 27315},    //0.00c
	Kamen:     {Kamen, Lava, MinTemp, 157315}, //1300.00c
	Drvo:      {Drvo, Vatra, MinTemp, 87315}, //600.00c spontano zapaljenje
	Sljunak:   {Sljunak, Lava, MinTemp, 157312}, //kamen
	Pesak:     {Pesak, Lava, MinTemp, 197315}, //1700.00c
	So:        {So, Lava, MinTemp, 107315},    //800.00c
	Rdja:      {Rdja, Lava, MinTemp, 177315},  // metal
	Lava:      {Lava, Lava, MinTemp, MaxTemp},
	Voda:      {Led, Para, 27315, 37315},          //0.00c, 100.00c
	SlanaVoda: {Led, Para, 25315, 37315},          //-20.00c, 100c
	Para:      {Voda, Para, 37315, MaxTemp},       //100.00c
	Vatra:     {Prazno, Plazma, 57315, 527315},    //300.00c, 5000.00c
	TecniAzot: {TecniAzot, Prazno, MinTemp, 7315}, //-200.00c
	Plazma:    {Vatra, Plazma, 527315, MaxTemp},  //5000.00c
	Zid:       {Zid, Zid, MinTemp, MaxTemp},
}

const MinTemp uint64 = 0      // 0.00k
const MaxTemp uint64 = 827315 //8000.00c

var Zapaljiv = map[Materijal]bool{
	Prazno:    false,
	Metal:     false,
	Led:       false,
	Kamen:     false,
	Drvo:      true, ///omggggg sooo truueeee bestieeeeee slayyyy queeen
	Sljunak:   false,
	Pesak:     false,
	So:        false,
	Rdja:      false,
	Lava:      false,
	Voda:      false,
	SlanaVoda: false,
	Para:      false,
	Vatra:     false, ///da li je voda mokra xDDD
	TecniAzot: false,
	Plazma:    false,
	Zid:       false,
}

type Cestica struct {
	Materijal   Materijal
	Temperatura uint64
	BaferTemp   uint64
	SekMat      Materijal
	Ticker      int32
}

func NewCestica(materijal Materijal) Cestica {
	zrno := Cestica{
		Materijal:   materijal,
		Temperatura: 29315, //20.00c
		BaferTemp:   0,
		SekMat:      Prazno,
		Ticker:      1023, //za rdju gorivo itd, opada po principu nuklearnog raspada (svaki frejm ima x% sanse da ga dekrementira, na 0 prelazi u drugo stanje)
	}
	if materijal == Led {
		zrno.SekMat = Voda
		zrno.Temperatura = 24315 //-30.00c
	}
	if materijal == Drvo {
		zrno.Ticker = 64
	}
	if materijal == Para {
		zrno.SekMat = Voda
		zrno.Temperatura = 42315 //150.00c
	}
	if materijal == Vatra {
		zrno.Temperatura = 77315 //500.00c
		zrno.Ticker = 8
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

func UpdateTemp(matrix [][]Cestica, i int, j int) {
	if matrix[i][j].Materijal == Zid {
		matrix[i][j].BaferTemp = 29315
		return
	}
	trenutna := matrix[i][j]

	/**/
	temperatura := trenutna.Temperatura
	parcePice := float32(temperatura) / 9
	for k := -1; k < 2; k++ {
		for l := -1; l < 2; l++ {
			if matrix[i+k][j+l].Materijal != Zid {
				if matrix[i+k][j+l].Materijal == Prazno || matrix[i][j].Materijal == Prazno {
					matrix[i+k][j+l].BaferTemp += uint64(parcePice/100)
					temperatura = temperatura - uint64(parcePice/100)
				} else if matrix[i][j].Materijal == Voda && matrix[i+k][j+l].Materijal == Vatra {
					matrix[i+k][j+l].BaferTemp += uint64(parcePice)
					temperatura = temperatura - uint64(parcePice)
				} else if matrix[i][j].Materijal == Vatra && matrix[i+k][j+l].Materijal == Voda {
					matrix[i+k][j+l].BaferTemp += uint64(parcePice/100)
					temperatura = temperatura - uint64(parcePice/100)
				} else {
					matrix[i+k][j+l].BaferTemp += uint64(parcePice)
					temperatura = temperatura - uint64(parcePice)
				}
			}
		}
	}
	matrix[i][j].BaferTemp += temperatura
	/**/
}

func UpdatePhaseOfMatter(matrix [][]Cestica, i int, j int) {

	if matrix[i][j].Materijal == Zid {
		return
	}

	trenutna := matrix[i][j]
	sekmat := trenutna.SekMat
	materijal := trenutna.Materijal
	temperatura := trenutna.Temperatura

	if materijal == Vatra {
		if trenutna.Ticker == 0 {
			matrix[i][j].Materijal = Prazno
		} else if trenutna.Ticker > 10 {
			matrix[i][j].Ticker = 10
		} else {
			matrix[i][j].Ticker--
		}
		for k := -1; k < 2; k++ {
				for l := -1; l < 2; l++ {
					if matrix[i+k][j+l].Materijal == Voda || matrix[i+k][j+l].Materijal == Para{
						matrix[i][j].Materijal = Prazno
					}
				}
			}
	}

	if materijal == Lava {
		if temperatura < MapaFaza[sekmat].TackaKljucanja {
			if sekmat == Rdja {
				matrix[i][j].Materijal = Metal
				matrix[i][j].Ticker = 127
			} else {
				matrix[i][j].Materijal = sekmat
			}
		}
	} else if materijal == SlanaVoda {
		if temperatura < MapaFaza[materijal].TackaTopljenja {
			matrix[i][j].Materijal = MapaFaza[materijal].Nize
			matrix[i][j].SekMat = SlanaVoda
		} else if temperatura > MapaFaza[materijal].TackaKljucanja {
			rFaktor := rand.Intn(2)*2 - 1

			if matrix[i][j-1].Materijal == Prazno {
				matrix[i][j].Materijal = So
				matrix[i][j-1].Materijal = Para
				matrix[i][j-1].Temperatura = matrix[i][j].Temperatura
			} else if matrix[i+rFaktor][j-1].Materijal == Prazno {
				matrix[i][j].Materijal = So
				matrix[i+rFaktor][j-1].Materijal = Para
				matrix[i+rFaktor][j-1].Temperatura = matrix[i][j].Temperatura
			} else if matrix[i-rFaktor][j-1].Materijal == Prazno {
				matrix[i][j].Materijal = So
				matrix[i-rFaktor][j-1].Materijal = Para
				matrix[i-rFaktor][j-1].Temperatura = matrix[i][j].Temperatura
			} else {
				return
			}

		}
	} else {
		if temperatura < MapaFaza[materijal].TackaTopljenja {
			matrix[i][j].Materijal = MapaFaza[materijal].Nize
		} else if temperatura > MapaFaza[materijal].TackaKljucanja {
			matrix[i][j].Materijal = MapaFaza[materijal].Vise
			if matrix[i][j].SekMat == SlanaVoda {
				matrix[i][j].Materijal = SlanaVoda
			}
			matrix[i][j].SekMat = materijal
		}
	}

	if materijal == So { //Gospode oprosti mi za ovaj blok koda bio sam mlad i naivan nisam znao za bolje -s
		rFaktor := rand.Intn(2)*2 - 1

		if matrix[i][j+1].Materijal == Voda {
			matrix[i][j].Materijal = Prazno
			matrix[i][j+1].Materijal = SlanaVoda
		} else if matrix[i+rFaktor][j+1].Materijal == Voda {
			matrix[i][j].Materijal = Prazno
			matrix[i+rFaktor][j+1].Materijal = SlanaVoda
		} else if matrix[i-rFaktor][j+1].Materijal == Voda {
			matrix[i][j].Materijal = Prazno
			matrix[i-rFaktor][j+1].Materijal = SlanaVoda
		} else if matrix[i+rFaktor][j].Materijal == Voda {
			matrix[i][j].Materijal = Prazno
			matrix[i+rFaktor][j].Materijal = SlanaVoda
		} else if matrix[i-rFaktor][j].Materijal == Voda {
			matrix[i][j].Materijal = Prazno
			matrix[i-rFaktor][j].Materijal = SlanaVoda
		} else if matrix[i+rFaktor][j-1].Materijal == Voda {
			matrix[i][j].Materijal = Prazno
			matrix[i+rFaktor][j-1].Materijal = SlanaVoda
		} else if matrix[i-rFaktor][j-1].Materijal == Voda {
			matrix[i][j].Materijal = Prazno
			matrix[i-rFaktor][j-1].Materijal = SlanaVoda
		} else if matrix[i][j+1].Materijal == Voda {
			matrix[i][j].Materijal = Prazno
			matrix[i][j+1].Materijal = SlanaVoda
		}
	}

	if materijal == Metal {
		for k := -1; k < 2; k++ {
			for l := -1; l < 2; l++ {
				if matrix[i+k][j+l].Materijal == SlanaVoda {
					randBr := rand.Intn(7)
					if randBr > 3 {
						matrix[i][j].Ticker -= 1
					}
				} else if matrix[i+k][j+l].Materijal == Voda {
					randBr := rand.Intn(7)
					if randBr > 5 {
						matrix[i][j].Ticker -= 1
					}
				}
			}
		}
		if matrix[i][j].Ticker < 0 {
			matrix[i][j].Materijal = Rdja
		}
	}

	//gorenje
	if Zapaljiv[materijal] {
		
		if sekmat != Vatra{
			for k := -1; k < 2; k++ {
				for l := -1; l < 2; l++ {
					if matrix[i+k][j+l].Materijal == Vatra {
						matrix[i][j].SekMat = Vatra
					}
				}
			}
		}

		if sekmat == Vatra {
			for k:=-1; k < 2; k++ {
				random := rand.Intn(10)
				if random > 6 {
				matrix[i][j].Ticker = matrix[i][j].Ticker - 1
					if matrix[i+k][j-1].Materijal == Prazno {
						matrix[i+k][j-1].Materijal = Vatra
						matrix[i+k][j-1].Temperatura = 87315 //600.00c
						matrix[i+k][j-1].Ticker = 8
					}
				}
			}
			imaVazduha := false
			for k := -1; k < 2; k++ {
				for l := -1; l < 2; l++ {
					if matrix[i+k][j+l].Materijal == Voda {
						matrix[i][j].SekMat = Prazno
					}
					if matrix[i+k][j+l].Materijal == Prazno {
						imaVazduha = true
					}
				}
			}
			if !imaVazduha{
				matrix[i][j].SekMat = Prazno
			}
		}
		if matrix[i][j].Ticker < 1 {
			matrix[i][j].Materijal = Vatra
			matrix[i][j].SekMat = Prazno
			matrix[i][j].Ticker = 8
		}
	}

}

func UpdatePosition(matrix [][]Cestica, i int, j int) {
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

	if (astanje & 0b1000) != 0 {
		lRand := rand.Intn(3) - 1
		rRand := rand.Intn(3) - 1
		komsija := matrix[i+lRand][j+rRand]
		if komsija.Materijal == Prazno {
			matrix[i][j] = komsija
			matrix[i+lRand][j+rRand] = trenutna
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
		if (AStanje[komsija.Materijal]&0b0001 != 0) && smer*int(Gustina[komsija.Materijal]) < smer*int(Gustina[trenutna.Materijal]) { ///ovde samo dodati || bafer[i][j+smer].Materijal == Prazno za blokovsko padanje, slicno u ostalim delovima ove f je //ovaj komentar je zastareo i odnosi se na neku davno zaboravljenu arhitekturu projekta zakopanu tu negde izmedju Atlantide i Drazinog groba
			matrix[i][j+smer] = trenutna
			matrix[i][j] = komsija
			pomeren = true
		}
	}
	if pomeren {
		return
	}

	/**/
	if (astanje & 0b0010) != 0 {
		rFaktor := rand.Intn(2)*2 - 1 //{-1, 1}
		komsija1 := matrix[i+rFaktor][j+smer]
		if (AStanje[komsija1.Materijal]&0b0010 != 0) && smer*int(Gustina[komsija1.Materijal]) < smer*int(Gustina[trenutna.Materijal]) {
			matrix[i+rFaktor][j+smer] = trenutna
			matrix[i][j] = komsija1
			pomeren = true
			return
		}
		komsija2 := matrix[i-rFaktor][j+smer]
		if (AStanje[komsija2.Materijal]&0b0010 != 0) && smer*int(Gustina[komsija2.Materijal]) < smer*int(Gustina[trenutna.Materijal]) {
			matrix[i-rFaktor][j+smer] = trenutna
			matrix[i][j] = komsija2
			pomeren = true
			return
		}
	}
	/**/
	if (astanje & 0b0100) != 0 {
		rFaktor := rand.Intn(2)*2 - 1 //{-1, 1}
		if matrix[i+rFaktor][j].Materijal == Prazno {
			if matrix[i+rFaktor+rFaktor][j].Materijal == Prazno {
				matrix[i+rFaktor+rFaktor][j], matrix[i][j] = trenutna, matrix[i+rFaktor+rFaktor][j]
			} else {
				matrix[i+rFaktor][j], matrix[i][j] = trenutna, matrix[i+rFaktor][j]
			}
		} else if matrix[i-rFaktor][j].Materijal == Prazno {
			if matrix[i-rFaktor-rFaktor][j].Materijal == Prazno {
				matrix[i-rFaktor-rFaktor][j], matrix[i][j] = trenutna, matrix[i-rFaktor-rFaktor][j]
			} else {
				matrix[i-rFaktor][j], matrix[i][j] = trenutna, matrix[i-rFaktor][j]
			}
		}
		pomeren = true

	}
	/**/

	if pomeren {
		return
	}

}
