package main

import 
(
	"github.com/veandco/go-sdl2/sdl"
	"math/rand"
//	"math"
//	"time"
	"fmt"
)

//(lenguedž kolega) molim vas dvojicu kada stignete i kada odlucite da nam vise nisu potrebni,
//obrisete sve zakomentarisane delove koda jer ovo sada izgleda uzasno
//a mislim da ima potencijala da bude relativno uredno i sazeto, makar mejn

type Materijal int

const (
    Zid Materijal = -1
    Prazno Materijal = 0
    Pesak Materijal = 1
    Voda Materijal = 2
    Kamen Materijal = 3
    Metal Materijal = 4
)

var boja = map[Materijal]uint32{
    Zid : 0xffffff,
    Prazno : 0x000000,
    Pesak : 0xffff66,
    Voda : 0x3333ff,
    Kamen : 0x666666,
    Metal : 0x33334b,
}

type Cestica struct{

    materijal Materijal

}

const sirinaKanvasa, visinaKanvasa = 240, 144
const brojPikselaPoCestici = 4
	//golang nema makroe i ne moze praviti nizove/matrice dinamicke duzine
	//ali jedna fantasticna stvar je ta sto konstante rade poso makroa:
	//	niz moze primiti konstantnu promenjivu (ironican termin) za dimenziju
	//	dovoljno su zilave da ne moramo vise ni raditi kastovanje tipova, sirina prolazi za int, uint, int32, float, itd itd
	// B)
	//e sad sta mislimo o globalnim varijablama je druga prica...

// materijal koji nastaje levim klikom
// 0-vazduh 1-pesak 2-kamen 3-voda
var mat Materijal
var keystates = sdl.GetKeyboardState()
func main() {
	// pesak default
	mat = Pesak
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow("test", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		sirinaKanvasa*brojPikselaPoCestici, visinaKanvasa*brojPikselaPoCestici, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	surface, err := window.GetSurface()
	if err != nil {
		panic(err)
	}
	surface.FillRect(nil, 0)

	var matrix[sirinaKanvasa][visinaKanvasa]Materijal
	slajs := matrixToSlice(matrix)

	for i := 0; i < sirinaKanvasa; i++ {
        slajs[i][0] = Zid
        slajs[i][visinaKanvasa-1] = Zid
    }
    for j := 0; j < visinaKanvasa; j++ {
        slajs[0][j] = Zid
        slajs[sirinaKanvasa-1][j] = Zid
    }

	running := true
	for running {
		running = pollEvents(slajs)
		update(slajs)
		render(slajs, surface)

		window.UpdateSurface()
	}

}

func pollEvents(matrix [][]Materijal) bool {
	running := true
	keystates = sdl.GetKeyboardState()

	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent(){
		switch event.(type){
			case *sdl.QuitEvent:
				running = false
				break
			case *sdl.KeyboardEvent:
				// može ovo lepše #TODO
				if keystates[sdl.SCANCODE_ESCAPE] != 0 {
					running = false
				}
				if keystates[sdl.SCANCODE_0] != 0 {
					mat = 0
				}
				if keystates[sdl.SCANCODE_1] != 0 {
					mat = 1
				}
				if keystates[sdl.SCANCODE_2] != 0 {
					mat = 2
				}
				if keystates[sdl.SCANCODE_3] != 0 {
					mat = 3
				}
			default:
				/* code */
		}
	}

	var x, y int32
	var state uint32
	x, y, state = sdl.GetMouseState()
	fmt.Printf("%d %d %d\n", x, y, state)
	if x < 0 {
		x = 0
	} else if x > sirinaKanvasa * brojPikselaPoCestici-1 {
		x = sirinaKanvasa * brojPikselaPoCestici-1
	}
	if y < 0 {
		y = 0
	} else if y > visinaKanvasa * brojPikselaPoCestici-1 {
		y = visinaKanvasa * brojPikselaPoCestici-1
	}
	if state == 1 {
		if matrix[x/brojPikselaPoCestici][y/brojPikselaPoCestici] == 0 {
			matrix[x/brojPikselaPoCestici][y/brojPikselaPoCestici] = mat
		}
	}

	return running

}

func update(matrix [][]Materijal) {
	for j := visinaKanvasa-3; j > 0; j-- {
		for i := 1; i < sirinaKanvasa-1; i++ {

			// pesak
			if matrix[i][j] == Pesak {
				flutter := rand.Intn(2)
				if matrix[i][j+1] == Prazno {
					if flutter == 1 {
						matrix[i][j], matrix[i][j+1] = Prazno, Pesak
					}
					//ako moze pasti direkt neka padne
				} else {
					//u suprotnom gleda moze li dijagonalu
					var sgn int
					if rand.Intn(2) == 1 {
						sgn = 1
					} else {
						sgn = -1
					}//nasumice biramo na koju stranu prvo ide da izbegnemo pristrasno padanje

					if (matrix[i+sgn][j+1] == Prazno) && (i+sgn > 0) && (i+sgn < sirinaKanvasa) {
						matrix[i][j], matrix[i+sgn][j+1] = Prazno, Pesak
					} else if (matrix[i-sgn][j+1] == Prazno) && (i+sgn > 0) && (i+sgn < sirinaKanvasa) {
						matrix[i][j], matrix[i-sgn][j+1] = Prazno, Pesak
					} 
				}
			}
			// voda
			// ovo je teže nego što smo mislili
			if matrix[i][j] == Voda {
				flutter := rand.Intn(2)
				if matrix[i][j+1] == Prazno {
					if flutter == 1 {
						matrix[i][j], matrix[i][j+1] = Prazno, Voda
					}
					//ako moze pasti direkt neka padne
				} else {
					var sgn int
					if rand.Intn(2) == 1 {
						sgn = 1
					} else {
						sgn = -1
					}

					if (matrix[i+sgn][j+1] == Prazno) && (i+sgn > 0) && (i+sgn < sirinaKanvasa) {
						matrix[i][j], matrix[i+sgn][j+1] = Prazno, Voda
					} else if (matrix[i-sgn][j+1] == Prazno) && (i+sgn > 0) && (i+sgn < sirinaKanvasa) {
						matrix[i][j], matrix[i-sgn][j+1] = Prazno, Voda
					} else if (matrix[i+sgn][j] == Prazno) && (i+sgn > 0) && (i+sgn < sirinaKanvasa) {
						matrix[i][j], matrix[i+sgn][j] = Prazno, Voda
					} else if (matrix[i-sgn][j] == Prazno) && (i+sgn > 0) && (i+sgn < sirinaKanvasa) {
						matrix[i][j], matrix[i-sgn][j] = Prazno, Voda
					}
				}
			}
		}
	}
}

func render(matrix [][]Materijal, surface *sdl.Surface) {
	for i := 0; i < sirinaKanvasa; i++ {
		for j := 0; j < visinaKanvasa; j++ {
			rect := sdl.Rect{int32(i*brojPikselaPoCestici), int32(j*brojPikselaPoCestici), brojPikselaPoCestici,brojPikselaPoCestici}
			// ovo je lepo
			surface.FillRect(&rect, boja[matrix[i][j]])
		}
	}
}

func matrixToSlice(matrix [sirinaKanvasa][visinaKanvasa]Materijal) [][]Materijal {
	// ma biće dobro
	slajs := make([][]Materijal, len(matrix))

	for i := 0; i < len(matrix); i++ {
		kolona := make([]Materijal, len(matrix[i]))
		for j := 0; j < len(matrix[i]); j++ {
			kolona[j] = matrix[i][j]
		}
		slajs[i] = kolona
	}

	return slajs
	//dakle poenta je pravimo slajs cele matrice na pocetku mejna
	//i dalje sve sto radimo radimo sa slajsom
	//jer ne znam ni ja valja onda radi poreferenci a ne vrednosti
	//pa mozemo imati izdvojene funkcije majku mu njegovu
	//nakon kratkog razmatranja kapiram da bi redosled ova dva komentara trebalo biti zamenjen, #TODO

}