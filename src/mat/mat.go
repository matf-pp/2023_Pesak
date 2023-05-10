//Package mat sadzi osobine svih materijala (u globalnim mapama) i njihove interakcije (u funkcijama Update_) nalaze se ovde
package mat

import (
	"math/rand"
)

//Obrnuto odredjuje smer gravitacije
var Obrnuto = 1
//GravityRukavica određuje da li je rukavica za gravitaciju uključena ili isključena
var GravityRukavica bool = false

//KursorPoslednjiX je posledja x koordinata misa
var KursorPoslednjiX = int32(0)
//KursorPoslednjiY je posledja y koordinata misa
var KursorPoslednjiY = int32(0)

//IzabraniJezik: srpski, engleski, ...
var IzabraniJezik = 0

//Materijal je gradivna jedinica celog projekta
type Materijal int

//materijali sa kojima radimo
const (
	Prazno    Materijal = 0
	Metal     Materijal = 1
	Led       Materijal = 2
	Kamen     Materijal = 3
	Drvo      Materijal = 4
	Sljunak   Materijal = 5
	Pesak     Materijal = 6
	So        Materijal = 7
	Rdja      Materijal = 253
	Lava      Materijal = 8
	Voda      Materijal = 9
	Zejtin    Materijal = 10
	Kiselina  Materijal = 11
	SlanaVoda Materijal = 254
	Para      Materijal = 12
	Vatra     Materijal = 13
	Dim       Materijal = 255
	TecniAzot Materijal = 14
	Plazma    Materijal = 15
	Toplo     Materijal = 16
	Hladno    Materijal = 17
	Zid       Materijal = 256
)

//Ime materijala koje se ispisuje pri haverovanju misem preko cestice ili dugmeta
var Ime = map[Materijal][]string {
			  // srpski   engleski
	Prazno:    {"Prazno", "Nothing"},
	Metal:     {"Metal", "Metal"},
	Led:       {"Led", "Ice"},
	Kamen:     {"Kamen", "Rock"},
	Drvo:      {"Drvo", "Wood"},
	Sljunak:   {"Sljunak", "Gravel"},
	Pesak:     {"Pesak", "Sand"},
	So:        {"So", "Salt"},
	Rdja:      {"Rdja", "Rust"},
	Lava:      {"Lava", "Lava"},
	Voda:      {"Voda", "Water"},
	Zejtin:    {"Zejtin", "Olive oil"},
	Kiselina:  {"Kiselina", "Acid"},
	SlanaVoda: {"Slana voda", "Saltwater"},
	Para:      {"Para", "Steam"},
	Vatra:     {"Vatra", "Fire"},
	Dim:       {"Dim", "Smoke"},
	TecniAzot: {"Tecni Azot", "Liquid nitrogen"},
	Plazma:    {"Plazma", "Plazma"},
	Toplo:     {"Toplo", "Warm"},
	Hladno:    {"Hladno", "Cold"},
	Zid:       {"Zid", "Wall"},
}

//Boja cestice (za neke materijale se zove funkcija koja u obzir uzima druge osobine)
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
	Zejtin:    0x3b3131,
	Kiselina:  0xb0bf1a,
	SlanaVoda: 0x4444ff,
	Para:      0x6666ff,
	Vatra:     0xd73502,
	Dim:       0x222222,
	TecniAzot: 0x99ff99,
	Plazma:    0x5f007f,
	Toplo:     0xff0000,
	Hladno:    0x00ffff,
	Zid:       0xffffff,
}

//Gustina je broj koji odredjuje prioritet plutanja
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
	Zejtin:    1,
	Kiselina:  2,
	SlanaVoda: 3,
	Para:      -3,
	Vatra:     -4,
	Dim:       -5,
	TecniAzot: 3,
	Plazma:    0,
	Zid:       0,
}

//GustinaBoja je boja cestica u takozvanom Gustina modu
var GustinaBoja = map[Materijal]uint32 {
	Prazno: 	0xc8c8c8,
	Metal:  	0x00ff00,
	Led:    	0x004600,
	Kamen:  	0x00b400,
	Drvo:      	0x00b400,
	Sljunak:   	0x00a000,
	Pesak:  	0x007800,
	So:        	0x00a000,
	Rdja:      	0x00ff00,
	Lava:   	0x00c800,
	Voda:   	0x005000,
	Zejtin:     0x00aa00,
	Kiselina:   0x005800,
	SlanaVoda: 	0x005a00,
	Para:   	0xc800c8,
	Vatra:     	0xc800c8,
	Dim:        0xfa00fa,
	TecniAzot: 	0x006400,
	Plazma:    	0xff00ff,
	Zid:		0,
}

//AStanje je odredjeno pravilima:
//0000 ne pomera se
//---1 pada direktno
//--1- pada dijagonalno
//-1-- curi horizontalno
//1--- pomera se nasumicno svuda
//u zavisnosti od bitova materijal se ponasa drugacije u funkciji UpdatePosition
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
	Zejtin:    0b0111,
	SlanaVoda: 0b0111,
	Kiselina:  0b0111,
	Para:      0b0111,
	Vatra:     0b0111,
	Dim:       0b0111,
	TecniAzot: 0b0111,
	Plazma:    0b1111,
	Zid:       0b0000,
}

//FaznaPromena odredjuje pri kojim temperaturama koj materijal prelazi u koj drugi
type FaznaPromena struct {
	Nize           Materijal
	Vise           Materijal
	TackaTopljenja uint64
	TackaKljucanja uint64
}


//MapaFaza zavisi od FaznaPromena strukture
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
	Zejtin:    {Zejtin, Vatra, MinTemp, 67315}, //TODO: mast? 400.00c
	SlanaVoda: {Led, Para, 25315, 37315},          //-20.00c, 100c
	Kiselina:  {Kiselina, Kiselina, MinTemp, MaxTemp},
	Para:      {Voda, Para, 37315, MaxTemp},       //100.00c
	Vatra:     {Dim, Plazma, 57315, 527315},    //300.00c, 5000.00c
	Dim:       {Prazno, Vatra, 32315, MaxTemp}, //50.00c, 600.00c
	TecniAzot: {TecniAzot, Prazno, MinTemp, 7315}, //-200.00c
	Plazma:    {Vatra, Plazma, 527315, MaxTemp},  //5000.00c
	Zid:       {Zid, Zid, MinTemp, MaxTemp},
}

//MinTemp koju dozvoljavamo
const MinTemp uint64 = 0      // 0.00k
//MaxTemp koju dozvoljavamo
const MaxTemp uint64 = 827315 //8000.00c

//Zapaljiv daje true ako je zapaljiv
var Zapaljiv = map[Materijal]bool{
	Prazno:    false,
	Metal:     false,
	Led:       false,
	Kamen:     false,
	Drvo:      true,
	Sljunak:   false,
	Pesak:     false,
	So:        false,
	Rdja:      false,
	Lava:      false,
	Voda:      false,
	Zejtin:    true,
	SlanaVoda: false,
	Para:      false,
	Vatra:     false,
	Dim:       false,
	TecniAzot: false,
	Kiselina:  false,
	Plazma:    false,
	Zid:       false,
}

//Cestica je struktura u kojoj cuvamo sve potrebne informacije o cestici
type Cestica struct {
	Materijal   Materijal
	Temperatura uint64
	BaferTemp   uint64
	SekMat      Materijal
	Ticker      int32
}

//NewCestica prima Materijal, konstruise novu Cesticu i vraca je
func NewCestica(materijal Materijal) Cestica {
	zrno := Cestica{
		Materijal:   materijal,
		Temperatura: 29315, //20.00c
		BaferTemp:   0,
		SekMat:      Prazno,
		Ticker:      1023, //za rdju, gorivo i druge, opada po principu nuklearnog raspada
						   //(svaki frejm ima x% sanse da ga dekrementira, na 0 prelazi u drugo stanje)
	}
	if materijal == Led {
		zrno.SekMat = Voda
		zrno.Temperatura = 24315 //-30.00c
	}
	if materijal == Drvo {
		zrno.Ticker = 64
	}
	if materijal == Zejtin {
		zrno.Ticker = 4
	}
	if materijal == Para {
		zrno.SekMat = Voda
		zrno.Temperatura = 42315 //150.00c
	}
	if materijal == Vatra {
		zrno.Temperatura = 77315 //500.00c
		zrno.Ticker = vatraTiker
	}
	if materijal == Dim {
		zrno.Temperatura = 47315 //200.00c
		zrno.Ticker = dimTiker
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

//UpdateTemp prima matricu Cestica i koordinate jedne, na osnovu trenutnih temperatura sebe i suseda racuna narednu temperaturu, smesta je u BaferTemp (unutar mejna se BaferTemp primenjuje)
func UpdateTemp(matrix [][]Cestica, i int, j int) {
	if matrix[i][j].Materijal == Zid {
		matrix[i][j].BaferTemp = 29315
		return
	}
	trenutna := matrix[i][j]

	temperatura := trenutna.Temperatura
	parcePice := float32(temperatura) / 9
	for k := -1; k < 2; k++ {
		for l := -1; l < 2; l++ {
			if matrix[i+k][j+l].Materijal != Zid {
				if matrix[i+k][j+l].Materijal == Prazno && matrix[i][j].Materijal == Prazno{
					matrix[i+k][j+l].BaferTemp += uint64(parcePice/10)
					temperatura = temperatura - uint64(parcePice/10)
				}
				if matrix[i+k][j+l].Materijal == Prazno || matrix[i][j].Materijal == Prazno {
					matrix[i+k][j+l].BaferTemp += uint64(parcePice/100)
					temperatura = temperatura - uint64(parcePice/100)
				} else if matrix[i][j].Materijal == Vatra && matrix[i+k][j+l].Materijal == Voda {
					matrix[i+k][j+l].BaferTemp += uint64(parcePice/100)
					temperatura = temperatura - uint64(parcePice/100)
				} else {
					matrix[i+k][j+l].BaferTemp += uint64(parcePice/2)
					temperatura = temperatura - uint64(parcePice/2)
				}
			}
		}
	}
	matrix[i][j].BaferTemp += temperatura
}

const vatraTiker = 16
const dimTiker = 64

//UpdatePhaseOfMatter vrsi promenu cestica iz jednog u drugi materijal, ukoliko je to potrebno
//agregatno stanje u odnosu na temperaturu
//nagrizanje kiseline
//gorenje zapaljivih materijala
//itd
func UpdatePhaseOfMatter(matrix [][]Cestica, i int, j int) {
// ono sto je bitno je da radi (:
	if matrix[i][j].Materijal == Zid {
		return
	}

	trenutna := matrix[i][j]
	sekmat := trenutna.SekMat
	materijal := trenutna.Materijal
	temperatura := trenutna.Temperatura

	if materijal == Dim {
		if trenutna.Ticker < 0 {
			matrix[i][j].Materijal = Prazno
		} else if trenutna.Ticker > dimTiker {
			matrix[i][j].Ticker = dimTiker
		} else {
			matrix[i][j].Ticker--
		}
	}

	if materijal == Vatra {
		if trenutna.Ticker == 0 {
			matrix[i][j].Materijal = Dim
			matrix[i][j].Ticker = dimTiker
		} else if trenutna.Ticker > vatraTiker {
			matrix[i][j].Ticker = vatraTiker
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
			if matrix[i][j].Materijal == Vatra && matrix[i][j].Ticker > vatraTiker{
				matrix[i][j].Ticker = vatraTiker
			}
		} else if temperatura > MapaFaza[materijal].TackaKljucanja {
			matrix[i][j].Materijal = MapaFaza[materijal].Vise
			if matrix[i][j].SekMat == SlanaVoda {
				matrix[i][j].Materijal = SlanaVoda
			}
			matrix[i][j].SekMat = materijal
		}
	}

	if materijal == So {
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
						matrix[i][j].Ticker--
					}
				} else if matrix[i+k][j+l].Materijal == Voda {
					randBr := rand.Intn(7)
					if randBr > 5 {
						matrix[i][j].Ticker--
					}
				}
			}
		}
		if matrix[i][j].Ticker < 0 {
			matrix[i][j].Materijal = Rdja
			matrix[i][j].Ticker = 511
		}
	}

	if materijal == Kiselina {
		for k := -1; k < 2; k++ {
			for l := -1; l < 2; l++ {
				komsa := matrix[i+k][j+l]
				if komsa.Materijal != Kiselina && komsa.Materijal != Zid && komsa.Materijal != Prazno {
					matrix[i+k][j+l].Ticker -= 10
					matrix[i+k][j+l].Temperatura += 100
					matrix[i][j].Ticker -= 2
				}
			}
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
					if matrix[i+k][j-1*Obrnuto].Materijal == Prazno {
						matrix[i+k][j-1*Obrnuto].Materijal = Vatra
						matrix[i+k][j-1*Obrnuto].Temperatura = 87315 //600.00c
						matrix[i+k][j-1*Obrnuto].Ticker = vatraTiker
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
			matrix[i][j].Ticker = vatraTiker
		}

	}

	if matrix[i][j].Ticker < 0 {
		matrix[i][j].Materijal = Prazno
	}

}

//UpdatePosition radi promenu pozicije cestice ukoliko je moguce i potrebno
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
	
	if GravityRukavica {
		if j*6 < int(KursorPoslednjiY) {
			Obrnuto = 1
		} else {
			Obrnuto = -1
		}
	}
	smer *= Obrnuto

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

	if (astanje & 0b0001) != 0 {
		komsija := matrix[i][j+smer]
		//												( 1  *      G[v] = 2             <  1  *      g[ps] =  5) == True
		//                                              (-1  *      G[v] = 2             < -1  *      g[pr] = -5) == True
		if (AStanje[komsija.Materijal]&0b0001 != 0) && Obrnuto*smer*int(Gustina[komsija.Materijal]) < Obrnuto*smer*int(Gustina[trenutna.Materijal]) { ///ovde samo dodati || bafer[i][j+smer].Materijal == Prazno za blokovsko padanje, slicno u ostalim delovima ove f je //ovaj komentar je zastareo i odnosi se na neku davno zaboravljenu arhitekturu projekta zakopanu tu negde izmedju Atlantide i Drazinog groba
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
		if (AStanje[komsija1.Materijal]&0b0010 != 0) && Obrnuto*smer*int(Gustina[komsija1.Materijal]) < Obrnuto*smer*int(Gustina[trenutna.Materijal]) {
			matrix[i+rFaktor][j+smer] = trenutna
			matrix[i][j] = komsija1
			pomeren = true
			return
		}
		komsija2 := matrix[i-rFaktor][j+smer]
		if (AStanje[komsija2.Materijal]&0b0010 != 0) && Obrnuto*smer*int(Gustina[komsija2.Materijal]) < Obrnuto*smer*int(Gustina[trenutna.Materijal]) {
			matrix[i-rFaktor][j+smer] = trenutna
			matrix[i][j] = komsija2
			pomeren = true
			return
		}
	}
	// proverava moze li se zameniti ne samo sa horizontalnim susedom vec i sa dva odjednom, da bi se brze iznivelisala
	if (astanje & 0b0100) != 0 {
		rFaktor := rand.Intn(2)*2 - 1 //{-1, 1}
		komsa1 := matrix[i+rFaktor][j]
		komsa2 := matrix[i-rFaktor][j]
		komsa11 := matrix[i+rFaktor+rFaktor][j]
		komsa22 :=matrix[i-rFaktor-rFaktor][j]
		if (AStanje[komsa1.Materijal] & 0b0100) != 0 {
			if (AStanje[komsa11.Materijal] & 0b0100) != 0 {
				matrix[i+rFaktor+rFaktor][j], matrix[i][j] = trenutna, matrix[i+rFaktor+rFaktor][j]
			} else {
				matrix[i+rFaktor][j], matrix[i][j] = trenutna, matrix[i+rFaktor][j]
			}
		} else if (AStanje[komsa2.Materijal] & 0b0100) != 0 {
			if (AStanje[komsa22.Materijal] & 0b0100) != 0 {
				matrix[i-rFaktor-rFaktor][j], matrix[i][j] = trenutna, matrix[i-rFaktor-rFaktor][j]
			} else {
				matrix[i-rFaktor][j], matrix[i][j] = trenutna, matrix[i-rFaktor][j]
			}
		}
		pomeren = true

	}

	if pomeren {
		return
	}

}