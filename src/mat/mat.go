// Package mat sadzi osobine svih materijala (u globalnim mapama) i njihove interakcije (u funkcijama Update_) nalaze se ovde
package mat

import (
	"main/src/gravityPack"

	"math/rand"

	"github.com/theodesp/unionfind"
)

// KursorPoslednjiX je posledja x koordinata misa
var KursorPoslednjiX = int32(0)

// KursorPoslednjiY je posledja y koordinata misa
var KursorPoslednjiY = int32(0)

// IzabraniJezik: српски, srpski, engleski, poljski, tagalog, ...
var IzabraniJezik = 1

// BrJezika je broj jezika
var BrJezika = 5

// PoslMat je poslednji materijal koji se prikazuje
var PoslMat = 18

//velRupe je velicina crne rupe
var VelRupe = 8

// Materijal je gradivna jedinica celog projekta
type Materijal int

// materijali sa kojima radimo
const (
	Prazno    	Materijal = 0
	Metal     	Materijal = 1
	Led       	Materijal = 2
	SlaniLed 	Materijal = 251
	Kamen     	Materijal = 3
	Drvo      	Materijal = 4
	Biljka    	Materijal = 5
	Sljunak   	Materijal = 6
	Pesak     	Materijal = 7
	So        	Materijal = 8
	Rdja      	Materijal = 252
	Lava      	Materijal = 9
	Voda      	Materijal = 10
	Zejtin    	Materijal = 11
	Kiselina  	Materijal = 12
	SlanaVoda 	Materijal = 253
	Para      	Materijal = 13
	Vatra     	Materijal = 14
	Dim       	Materijal = 254
	Radijacija 	Materijal = 255
	TecniAzot 	Materijal = 15
	Plazma    	Materijal = 16
	Toplo     	Materijal = 17
	Hladno    	Materijal = 18
	Zid       	Materijal = 256
)

// Ime materijala koje se ispisuje pri haverovanju misem preko cestice ili dugmeta
var Ime = map[Materijal][]string{
	// српски, srpski, engleski, poljski, tagalog
	Prazno:    	{"Празно", "Prazno", "Nothing", "Nic", "Wala", "لا شيئ"},
	Metal:     	{"Метал", "Metal", "Metal", "Metal", "Metal", "معدن"},
	Led:       	{"Лед", "Led", "Ice", "Lód", "Yelo", "ملح"},
	SlaniLed:  	{"Слани лед", "Slani led", "Salty ice", "Słony lód", "Maalat na yelo", "جليد مالح"},
	Kamen:     	{"Камен", "Kamen", "Rock", "Skała", "Bato", "صخر"},
	Drvo:      	{"Дрво", "Drvo", "Wood", "Drewno", "Kahoy", "خشب"},
	Biljka:    	{"Биљка", "Biljka", "Plant", "Roślina", "Halaman", "نبات"},
	Sljunak:   	{"Шљунак", "Šljunak", "Gravel", "Żwir", "Graba", "حصى"},
	Pesak:     	{"Песак", "Pesak", "Sand", "Piasek", "Buhangin", "رمل"},
	So:        	{"Со", "So", "Salt", "Sól", "Asin", "ملح"},
	Rdja:      	{"Рђа", "Rđa", "Rust", "Rdza", "Kalawang", "الصدأ"},
	Lava:      	{"Лава", "Lava", "Lava", "Lawa", "Lava", "حمم بركانية"},
	Voda:      	{"Вода", "Voda", "Water", "Woda", "Tubig", "ماء"},
	Zejtin:    	{"Зејтин", "Zejtin", "Oil", "Olej", "Krudo", "زيت"},
	Kiselina:  	{"Киселина", "Kiselina", "Acid", "Kwas", "Asido", "حمض"},
	SlanaVoda: 	{"Слана вода", "Slana voda", "Saltwater", "Słona woda", "Tubig alat", "ماء مالح"},
	Para:      	{"Пара", "Para", "Steam", "Para", "Singaw", "بخار"},
	Vatra:     	{"Ватра", "Vatra", "Fire", "Ogień", "Apoy", "نار"},
	Dim:       	{"Дим", "Dim", "Smoke", "Dym", "Usok", "دخان"},
	Radijacija: {"Радијација", "Radijacija", "Radiation", "Radiacja", "?", "?"},
	TecniAzot: 	{"Течни азот", "Tečni Azot", "Liquid nitrogen", "Ciekły azot", "Liquid nitrogen", "نيتروجين سائل"},
	Plazma:    	{"Плазма", "Plazma", "Plasma", "Plazma", "Plasma", "بلازما"},
	Toplo:     	{"Топло", "Toplo", "Warm", "Ciepłe", "Painitin", "دافئ"},
	Hladno:    	{"Хладно", "Hladno", "Cold", "Zimne", "Palamigin", "بارد"},
	Zid:       	{"Зид", "Zid", "Wall", "Ściana", "Pader", "حائط"},
}

// Boja cestice (za neke materijale se zove funkcija koja u obzir uzima druge osobine)
var Boja = map[Materijal]uint32{
	Prazno:    	0x000000,
	Metal:     	0x33334b,
	Led:       	0xaaaaff,
	SlaniLed:  	0xababff,
	Kamen:     	0x999988,
	Drvo:      	0x994400,
	Biljka:    	0x00ff00,
	Sljunak:   	0x888877,
	Pesak:     	0xffff66,
	So:        	0xeeeeee,
	Rdja:      	0x6f0f2b,
	Lava:      	0xff6600,
	Voda:      	0x3333ff,
	Zejtin:    	0x3b3131,
	Kiselina:  	0xb0bf1a,
	SlanaVoda: 	0x4444ff,
	Para:      	0x6666ff,
	Vatra:     	0xd73502,
	Dim:       	0x222222,
	Radijacija:	0xffb84d,
	TecniAzot: 	0x99ff99,
	Plazma:    	0x5f007f,
	Toplo:     	0xff0000,
	Hladno:    	0x00ffff,
	Zid:       	0xffffff,
}

// Gustina je broj koji odredjuje prioritet plutanja
var Gustina = map[Materijal]int32{
	Prazno:    	0,
	Metal:     	0,
	Led:       	0,
	SlaniLed:  	0,
	Kamen:     	0,
	Drvo:      	0,
	Biljka:    	0,
	Sljunak:   	5,
	Pesak:     	5,
	So:        	5,
	Rdja:      	5,
	Lava:      	4,
	Voda:      	2,
	Zejtin:    	1,
	Kiselina:  	2,
	SlanaVoda: 	3,
	Para:       -3,
	Vatra:     	-4,
	Dim:       	-5,
	Radijacija: -5,
	TecniAzot: 	3,
	Plazma:    	0,
	Zid:       	0,
}

// GustinaBoja je boja cestica u takozvanom Gustina modu
var GustinaBoja = map[Materijal]uint32{
	Prazno:   	0xc8c8c8,
	Metal:     	0x00ff00,
	Led:       	0x004600,
	SlaniLed:  	0x004000,
	Kamen:     	0x00b400,
	Drvo:      	0x00b400,
	Biljka:    	0x00b400,
	Sljunak:   	0x00a000,
	Pesak:     	0x007800,
	So:        	0x00a000,
	Rdja:      	0x00ff00,
	Lava:      	0x00c800,
	Voda:      	0x005000,
	Zejtin:    	0x00aa00,
	Kiselina:  	0x005800,
	SlanaVoda: 	0x005a00,
	Para:      	0xc800c8,
	Vatra:     	0xc800c8,
	Dim:       	0xfa00fa,
	Radijacija:	0xfa00fa,
	TecniAzot: 	0x006400,
	Plazma:    	0xff00ff,
	Zid:       	0,
}

// AStanje je odredjeno pravilima:
// 0000 ne pomera se
// ---1 pada direktno
// --1- pada dijagonalno
// -1-- curi horizontalno
// 1--- pomera se nasumicno svuda
// u zavisnosti od bitova materijal se ponasa drugacije u funkciji UpdatePosition
var AStanje = map[Materijal]int{
	Prazno:    	0b1111,
	Metal:     	0b0000,
	Led:       	0b0000,
	SlaniLed:  	0b0000,
	Kamen:     	0b0000,
	Drvo:      	0b0000,
	Biljka:    	0b0000,
	Sljunak:   	0b0001,
	Pesak:     	0b0011,
	So:        	0b0011,
	Rdja:      	0b0011,
	Lava:      	0b0111,
	Voda:      	0b0111,
	Zejtin:    	0b0111,
	SlanaVoda: 	0b0111,
	Kiselina:  	0b0111,
	Para:      	0b0111,
	Vatra:     	0b0111,
	Dim:       	0b0111,
	Radijacija: 0b0111,
	TecniAzot: 	0b0111,
	Plazma:    	0b1111,
	Zid:       	0b0000,
}

// FaznaPromena odredjuje pri kojim temperaturama koj materijal prelazi u koj drugi
type FaznaPromena struct {
	Nize           Materijal
	Vise           Materijal
	TackaTopljenja uint64
	TackaKljucanja uint64
}

// MapaFaza zavisi od FaznaPromena strukture
var MapaFaza = map[Materijal]FaznaPromena{

	//k(c) = c+273.15
	//c(k) = k–273.15

	//100 int32		=	1.00k
	//130000 int32	=	1300.00k

	//MinTemp = 0.00k = -273.15c = int32(-27315)
	//maxtemp = 8000.00c = int32(800000)

	//	materijali	{Nize,	Vise,	TackaT,		TackaK}
	Prazno:    	{TecniAzot, Plazma, 5315, 627315},
	Metal:     	{Metal, Lava, MinTemp, 177315},        //1500.00c
	Led:       	{Led, Voda, MinTemp, 27315},           //0.00c
	SlaniLed:  	{SlaniLed, SlanaVoda, MinTemp, 25315}, //-20.00c
	Kamen:     	{Kamen, Lava, MinTemp, 157315},        //1300.00c
	Drvo:      	{Drvo, Vatra, MinTemp, 77315},         //500.00c spontano zapaljenje
	Biljka:    	{Biljka, Vatra, MinTemp, 87312},       //600.00c -||-
	Sljunak:   	{Sljunak, Lava, MinTemp, 157312},      //kamen
	Pesak:     	{Pesak, Lava, MinTemp, 197315},        //1700.00c
	So:        	{So, Lava, MinTemp, 107315},           //800.00c
	Rdja:      	{Rdja, Lava, MinTemp, 177315},         // metal
	Lava:      	{Lava, Lava, MinTemp, MaxTemp},
	Voda:      	{Led, Para, 27315, 37315},       //0.00c, 100.00c
	Zejtin:    	{Zejtin, Vatra, MinTemp, 67315}, //TODO: mast? 400.00c
	SlanaVoda: 	{SlaniLed, Para, 25315, 37315},  //-20.00c, 100c
	Kiselina:  	{Kiselina, Kiselina, MinTemp, MaxTemp},
	Para:      	{Voda, Para, 37315, MaxTemp},       //100.00c
	Vatra:     	{Dim, Plazma, 57315, 527315},       //300.00c, 5000.00c
	Dim:       	{Prazno, Vatra, 32315, MaxTemp},    //50.00c, 600.00c
	Radijacija:	{Prazno, Prazno, 32315, MaxTemp},
	TecniAzot: 	{TecniAzot, Prazno, MinTemp, 7315}, //-200.00c
	Plazma:    	{Vatra, Plazma, 527315, MaxTemp},   //5000.00c
	Zid:       	{Zid, Zid, MinTemp, MaxTemp},
}

// MinTemp koju dozvoljavamo
const MinTemp uint64 = 0 // 0.00k
// MaxTemp koju dozvoljavamo
const MaxTemp uint64 = 827315 //8000.00c

// Zapaljiv daje true ako je zapaljiv
var Zapaljiv = map[Materijal]bool{
	Prazno:    	false,
	Metal:     	false,
	Led:       	false,
	SlaniLed:  	false,
	Kamen:     	false,
	Drvo:      	true,
	Biljka:    	true,
	Sljunak:   	false,
	Pesak:     	false,
	So:        	false,
	Rdja:      	false,
	Lava:      	false,
	Voda:      	false,
	Zejtin:    	true,
	SlanaVoda: 	false,
	Para:      	false,
	Vatra:     	false,
	Dim:       	false,
	Radijacija: false,
	TecniAzot: 	false,
	Kiselina:  	false,
	Plazma:    	false,
	Zid:       	false,
}

// Cestica je struktura u kojoj cuvamo sve potrebne informacije o cestici
type Cestica struct {
	Materijal   Materijal
	Temperatura uint64
	BaferTemp   uint64
	SekMat      Materijal
	Ticker      int32
}

// NewCestica prima Materijal, konstruise novu Cesticu i vraca je
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
	if materijal == SlaniLed {
		zrno.SekMat = SlanaVoda
		zrno.Temperatura = 24315 //-30.00c
	}
	if materijal == Drvo {
		zrno.Ticker = 64
	}
	if materijal == Biljka {
		zrno.Ticker = 16
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
	if materijal == Radijacija {
		zrno.Temperatura = 47315 //200.00c
		zrno.Ticker = radijacijaTiker
	}
	if materijal == Lava {
		zrno.SekMat = Sljunak
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

// UpdateTemp prima matricu Cestica i koordinate jedne, na osnovu trenutnih temperatura sebe i suseda racuna narednu temperaturu, smesta je u BaferTemp (unutar mejna se BaferTemp primenjuje)
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
				if matrix[i+k][j+l].Materijal == Prazno && matrix[i][j].Materijal == Prazno {
					matrix[i+k][j+l].BaferTemp += uint64(parcePice / 10)
					temperatura = temperatura - uint64(parcePice / 10)
				}
				if matrix[i+k][j+l].Materijal == Prazno || matrix[i][j].Materijal == Prazno {
					matrix[i+k][j+l].BaferTemp += uint64(parcePice / 100)
					temperatura = temperatura - uint64(parcePice / 100)
				} else if matrix[i][j].Materijal == Vatra && matrix[i+k][j+l].Materijal == Voda {
					matrix[i+k][j+l].BaferTemp += uint64(parcePice / 100)
					temperatura = temperatura - uint64(parcePice / 100)
				} else {
					matrix[i+k][j+l].BaferTemp += uint64(parcePice / 2)
					temperatura = temperatura - uint64(parcePice / 2)
				}
			}
		}
	}
	matrix[i][j].BaferTemp += temperatura
}

const vatraTiker = 16
const dimTiker = 64
const radijacijaTiker = 64

// UpdatePhaseOfMatter vrsi promenu cestica iz jednog u drugi materijal, ukoliko je to potrebno
// agregatno stanje u odnosu na temperaturu
// nagrizanje kiseline
// gorenje zapaljivih materijala
// itd
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

	if materijal == Radijacija {
		if trenutna.Ticker < 0 {
			matrix[i][j].Materijal = Prazno
		} else if trenutna.Ticker > radijacijaTiker {
			matrix[i][j].Ticker = radijacijaTiker
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
				if matrix[i+k][j+l].Materijal == Voda || matrix[i+k][j+l].Materijal == Para {
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
			} else if sekmat == Kamen {
				matrix[i][j].Materijal = Sljunak
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
			if matrix[i][j].Materijal == Vatra && matrix[i][j].Ticker > vatraTiker {
				matrix[i][j].Ticker = vatraTiker
			}
		} else if temperatura > MapaFaza[materijal].TackaKljucanja {
			matrix[i][j].Materijal = MapaFaza[materijal].Vise
			//			if matrix[i][j].SekMat == SlanaVoda {
			//				matrix[i][j].Materijal = SlanaVoda
			//			}
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

		if sekmat != Vatra {
			for k := -1; k < 2; k++ {
				for l := -1; l < 2; l++ {
					if matrix[i+k][j+l].Materijal == Vatra {
						matrix[i][j].SekMat = Vatra
					}
				}
			}
		}

		if sekmat == Vatra {
			for k := -1; k < 2; k++ {
				random := rand.Intn(10)
				if random > 6 {
					matrix[i][j].Ticker = matrix[i][j].Ticker - 1
					if matrix[i+k][j-1*gravityPack.Obrnuto].Materijal == Prazno {
						matrix[i+k][j-1*gravityPack.Obrnuto].Materijal = Vatra
						matrix[i+k][j-1*gravityPack.Obrnuto].Temperatura = 87315 //600.00c
						matrix[i+k][j-1*gravityPack.Obrnuto].Ticker = vatraTiker
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
			if !imaVazduha {
				matrix[i][j].SekMat = Prazno
			}
		}

		if matrix[i][j].Ticker < 1 {
			matrix[i][j].Materijal = Vatra
			matrix[i][j].SekMat = Prazno
			matrix[i][j].Ticker = vatraTiker
		}

	}

	if materijal == Biljka {
		brojKomsVoda := 1
		for k := -1; k < 2; k++ {
			for l := -1; l < 2; l++ {
				if matrix[i+k][j+l].Materijal == Voda {
					brojKomsVoda *= 10
				}
			}
		}
		for k := -1; k < 2; k++ {
			for l := -1; l < 2; l++ {
				raste := rand.Intn(100000000/brojKomsVoda)
				if raste == 0 {
					if matrix[i+k][j+l].Materijal == Voda {
						matrix[i+k][j+l].Materijal = Biljka
						matrix[i+k][j+l].Ticker = 16
					}
				}
			}
		}
	}

	if matrix[i][j].Ticker < 0 {
		matrix[i][j].Materijal = Prazno
	}

}

var komponentePovezanostiVode = unionfind.New(10000)

var Nedostizna = Cestica{
	Materijal:   Voda,
	Temperatura: 29315,
	BaferTemp:   0,
	SekMat:      Prazno,
	Ticker:      1023,
}

// UpdatePosition radi promenu pozicije cestice ukoliko je moguce i potrebno
func UpdatePosition(matrix [][]Cestica, i int, j int) {
	//padanje

	if matrix[i][j].Materijal == Zid {
		return
	}

	// check za crnu rupu
	if gravityPack.CrnaRupa {
		if matrix[i][j].Materijal != Radijacija {
			if gravityPack.GTacka {
				if gravityPack.UpadaUCrnuRupu1(6*i, 6*j, 6*VelRupe) {
					random := rand.Intn(36*VelRupe)
					if random != 0 {
						matrix[i][j] = NewCestica(Prazno)
					} else {
						matrix[i][j] = NewCestica(Radijacija)
					}
				}
			} else if gravityPack.GRuka {
				if gravityPack.UpadaUCrnuRupu2(6*i, 6*j, 6*VelRupe, int(KursorPoslednjiX), int(KursorPoslednjiY)) {
					random := rand.Intn(36*VelRupe)
					if random != 0 {
						matrix[i][j] = NewCestica(Prazno)
					} else {
						matrix[i][j] = NewCestica(Radijacija)
					}
				}
			}
		}
	}

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

	crnaRupaMn := 1
	if gravityPack.CrnaRupa && (gravityPack.GTacka || gravityPack.GRuka) {
		if smer == -1 && matrix[i][j].Materijal != Radijacija {
			smer = 1
			crnaRupaMn = -1
		}
	} else if matrix[i][j].Materijal == Radijacija {
		astanje = 0b1000
	}

	oktant := 1
	smerI, smerJ := 0, 1
	smerILevo, smerJLevo, smerIDesno, smerJDesno := -1, 1, 1, 1
	smerILevo2, smerJLevo2, smerIDesno2, smerJDesno2 := -1, 0, 1, 0
	if gravityPack.GRuka || gravityPack.GTacka {
		if gravityPack.GRuka {
			oktant = gravityPack.ProveriOktant(i*6, j*6, int(KursorPoslednjiX), int(KursorPoslednjiY))
		} else {
			oktant = gravityPack.ProveriOktant(i*6 , j*6, int(gravityPack.CentarGravitacijeX), int(gravityPack.CentarGravitacijeY))
		}
		smerI, smerJ = gravityPack.GdePadaDole(oktant)
		smerILevo, smerJLevo, smerIDesno, smerJDesno = gravityPack.GdePadaUkosoDole(oktant)
		smerILevo2, smerJLevo2, smerIDesno2, smerJDesno2 = gravityPack.GdeIdeLevoDesno(oktant)
	}
	smerI *= gravityPack.Obrnuto * smer
	smerJ *= gravityPack.Obrnuto * smer
	smerILevo *= gravityPack.Obrnuto * smer
	smerJLevo *= gravityPack.Obrnuto * smer
	smerIDesno *= gravityPack.Obrnuto * smer
	smerJDesno *= gravityPack.Obrnuto * smer

	if (astanje & 0b1000) != 0 {
		lRand := rand.Intn(3) - 1
		rRand := rand.Intn(3) - 1
		komsija := matrix[i+lRand][j+rRand]
		if komsija.Materijal == Prazno {
			matrix[i][j] = komsija
			matrix[i+lRand][j+rRand] = trenutna
			pomeren = true
		}
		return
	}

	if pomeren {
		return
	}

	// pomeranje čvrstih ka crnoj rupi
	if gravityPack.CrnaRupa && (gravityPack.GTacka || gravityPack.GRuka) {
		komsija := matrix[i+smerI][j+smerJ]

		if komsija.Materijal != Zid {
			matrix[i+smerI][j+smerJ] = trenutna
			matrix[i][j] = komsija
		}
	}
	if pomeren {
		return
	}

	if (astanje & 0b0001) != 0 {
		komsija := matrix[i+smerI][j+smerJ]
		//												( 1  *      G[v] = 2             <  1  *      g[ps] =  5) == True
		//                                              (-1  *      G[v] = 2             < -1  *      g[pr] = -5) == True
		if (AStanje[komsija.Materijal]&0b0001 != 0) && smer*crnaRupaMn*int(Gustina[komsija.Materijal]) < smer*crnaRupaMn*int(Gustina[trenutna.Materijal]) {
			matrix[i+smerI][j+smerJ] = trenutna
			matrix[i][j] = komsija
			pomeren = true
		}
	}
	if pomeren {
		return
	}

	if gravityPack.CrnaRupa && (gravityPack.GTacka || gravityPack.GRuka) {
		trenutna = matrix[i][j]
		rFaktor := rand.Intn(2)*2 - 1 //{-1, 1}
		if rFaktor == -1 {
			rFaktor = smerILevo
			smerJ = smerJLevo
			komsija := matrix[i+rFaktor][j+smerJ]
			if komsija.Materijal != Zid {
				matrix[i+rFaktor][j+smerJ] = trenutna
				matrix[i][j] = komsija
				pomeren = true
				return
			}
		} else {
			rFaktor = smerIDesno
			smerJ = smerJDesno
			komsija := matrix[i+rFaktor][j+smerJ]
			if komsija.Materijal != Zid {
				matrix[i+rFaktor][j+smerJ] = trenutna
				matrix[i][j] = komsija
				pomeren = true
				return
			}
		}
	}
	if pomeren {
		return
	}

	/**/
	if (astanje & 0b0010) != 0 {
		trenutna = matrix[i][j]
		rFaktor := rand.Intn(2)*2 - 1 //{-1, 1}
		// koristi rFaktor da izabere smerILevo ili smerIDesno
		// ali ujedno i postavlja rFaktor na izabrani
		// smerJ samo prati rFaktor /limun
		if rFaktor == -1 {
			rFaktor = smerILevo
			smerJ = smerJLevo
			komsija := matrix[i+rFaktor][j+smerJ]
			if (AStanje[komsija.Materijal]&0b0010 != 0) && smer*crnaRupaMn*int(Gustina[komsija.Materijal]) < smer*crnaRupaMn*int(Gustina[trenutna.Materijal]) {
				matrix[i+rFaktor][j+smerJ] = trenutna
				matrix[i][j] = komsija
				pomeren = true
				return
			}
		} else {
			rFaktor = smerIDesno
			smerJ = smerJDesno
			komsija := matrix[i+rFaktor][j+smerJ]
			if (AStanje[komsija.Materijal]&0b0010 != 0) && smer*crnaRupaMn*int(Gustina[komsija.Materijal]) < smer*crnaRupaMn*int(Gustina[trenutna.Materijal]) {
				matrix[i+rFaktor][j+smerJ] = trenutna
				matrix[i][j] = komsija
				pomeren = true
				return
			}
		}
	}

	if gravityPack.CrnaRupa && (gravityPack.GTacka || gravityPack.GRuka) {
		trenutna = matrix[i][j]
		rFaktor := rand.Intn(2)*2 - 1 //{-1, 1}
		lFaktor := 0
		if rFaktor == -1 {
			rFaktor = smerILevo2
			lFaktor = smerJLevo2
		} else {
			rFaktor = smerIDesno2
			lFaktor = smerJDesno2
		}
		komsija1 := matrix[i+rFaktor][j+lFaktor]
		komsijaDalji1 := matrix[i+rFaktor+rFaktor][j+lFaktor+lFaktor]
		komsija2 := matrix[i-rFaktor][j-lFaktor]
		komsijaDalji2 := matrix[i-rFaktor-rFaktor][j-lFaktor-lFaktor]
		if komsija1.Materijal != Zid {
			if komsijaDalji1.Materijal != Zid {
				matrix[i+rFaktor+rFaktor][j+lFaktor+lFaktor], matrix[i][j] = trenutna, matrix[i+rFaktor+rFaktor][j+lFaktor+lFaktor]
			} else {
				matrix[i+rFaktor][j+lFaktor], matrix[i][j] = trenutna, matrix[i+rFaktor][j+lFaktor]
			}
		} else if komsija2.Materijal != Zid {
			if komsijaDalji2.Materijal != Zid {
				matrix[i-rFaktor-rFaktor][j-lFaktor-lFaktor], matrix[i][j] = trenutna, matrix[i-rFaktor-rFaktor][j-lFaktor-lFaktor]
			} else {
				matrix[i-rFaktor][j-lFaktor], matrix[i][j] = trenutna, matrix[i-rFaktor][j-lFaktor]
			}
		}
		pomeren = true
	}
	if pomeren {
		return
	}

	// proverava moze li se zameniti ne samo sa horizontalnim susedom vec i sa dva odjednom, da bi se brze iznivelisala
	if (astanje & 0b0100) != 0 {
		trenutna = matrix[i][j]
		rFaktor := rand.Intn(2)*2 - 1 //{-1, 1}
		lFaktor := 0
		if rFaktor == -1 {
			rFaktor = smerILevo2
			lFaktor = smerJLevo2
		} else {
			rFaktor = smerIDesno2
			lFaktor = smerJDesno2
		}
		komsija1 := matrix[i+rFaktor][j+lFaktor]
		komsijaDalji1 := matrix[i+rFaktor+rFaktor][j+lFaktor+lFaktor]
		komsija2 := matrix[i-rFaktor][j-lFaktor]
		komsijaDalji2 := matrix[i-rFaktor-rFaktor][j-lFaktor-lFaktor]
		if (AStanje[komsija1.Materijal] & 0b0100) != 0 {
			if (AStanje[komsijaDalji1.Materijal] & 0b0100) != 0 {
				matrix[i+rFaktor+rFaktor][j+lFaktor+lFaktor], matrix[i][j] = trenutna, matrix[i+rFaktor+rFaktor][j+lFaktor+lFaktor]
			} else {
				matrix[i+rFaktor][j+lFaktor], matrix[i][j] = trenutna, matrix[i+rFaktor][j+lFaktor]
			}
		} else if (AStanje[komsija2.Materijal] & 0b0100) != 0 {
			if (AStanje[komsijaDalji2.Materijal] & 0b0100) != 0 {
				matrix[i-rFaktor-rFaktor][j-lFaktor-lFaktor], matrix[i][j] = trenutna, matrix[i-rFaktor-rFaktor][j-lFaktor-lFaktor]
			} else {
				matrix[i-rFaktor][j-lFaktor], matrix[i][j] = trenutna, matrix[i-rFaktor][j-lFaktor]
			}
		}
		pomeren = true

	}

	if pomeren {
		return
	}
}
