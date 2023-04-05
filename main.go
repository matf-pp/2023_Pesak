package main

import 
(
	"github.com/veandco/go-sdl2/sdl"
	"math/rand"
//	"math"
//	"time"
	"fmt"
)

//ZA BOGA MILOGA molim vas dvojicu kada stignete i kada odlucite da nam vise nisu potrebni,
//obrisete sve zakomentarisane delove koda jer ovo sada izgleda uzasno
//a mislim da ima potencijala da bude relativno uredno i sazeto, makar mejn

const sirinaKanvasa, visinaKanvasa = 240, 144
const brojPikselaPoCestici = 10
	//golang nema makroe i ne moze praviti nizove/matrice dinamicke duzine
	//ali jedna fantasticna stvar je ta sto konstante rade poso makroa:
	//	niz moze primiti konstantnu promenjivu (ironican termin) za dimenziju
	//	dovoljno su zilave da ne moramo vise ni raditi kastovanje tipova, sirina prolazi za int, uint, int32, float, itd itd
	// B)
	//e sad sta mislimo o globalnim varijablama je druga prica...

func main() {

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

	var matrix[sirinaKanvasa][visinaKanvasa]int
	slajs := matrixToSlice(matrix)

//	var event sdl.Event
	running := true
	for running {
		// r := rand.Intn(800-3)
		// for i := 0; i < 3; i++ {
		// 	for j := 0; j < 3; j++ {
		// 		matrix[r+i][0+j] = 1
		// 	}
		// }

		running = pollEvents(slajs)
		update(slajs)
		render(slajs, surface)

		window.UpdateSurface()
		
		//time.Sleep(1 * time.Millisecond)
	}

}

func pollEvents(matrix [][]int) bool{
	/*var x int32
	var y int32
	var state uint32
	for event = sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch event.(type) {
		case *sdl.QuitEvent:
			running = false
		case *sdl.MouseMotionEvent:
			x, y, state = sdl.GetMouseState()
			fmt.Printf("%d %d %d\n", x, y, state)
		}
	}
		if(state == 1) {
		for i := x; i < x + 3 && i >= 0 && i < 798; i++ {
			for j := y; j < y + 3 && j >= 0 && j < 598; j++ {
				matrix[i][j] = 1
			}
		}
	}*///pollEvents
	running := true

	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent(){
		switch event.(type){
			case *sdl.QuitEvent:
				running = false
				break
//			case *sdl.MouseButtonEvent, *sdl.MouseMotionEvent:
//				var x, y int32
//				var state uint32
//				x, y, state = sdl.GetMouseState()
//				fmt.Printf("%d %d %d\n", x, y, state)
//				if state == 1{
//					if matrix[x/brojPikselaPoCestici][y/brojPikselaPoCestici] == 0 {
//						matrix[x/brojPikselaPoCestici][y/brojPikselaPoCestici] = 1
//					}
//				}
			default:
				/* code */
		}
	}

	var x, y int32
	var state uint32
	x, y, state = sdl.GetMouseState()
	fmt.Printf("%d %d %d\n", x, y, state)
	if state == 1 {
		if matrix[x/brojPikselaPoCestici][y/brojPikselaPoCestici] == 0 {
			matrix[x/brojPikselaPoCestici][y/brojPikselaPoCestici] = 1
		}
	}

	return running

}

func update(matrix [][]int){
	/*for i := 1; i < sirinaKanvasa-3; i++ {
		for j := 600-2; j >= 0; j-- {
			r := rand.Float64()
			sign := -1
			if math.Mod(r, 2) == 0 {
				sign = 1
			}
			if matrix[i][j] == 1 && matrix[i][j+1] == 0 {
				matrix[i][j] = 0
				matrix[i][j+1] = 1
			} else if matrix[i][j] == 1 && matrix[i+sign][j+1] == 0 {
				matrix[i][j] = 0
				matrix[i+sign][j+1] = 1
			} else if matrix[i][j] == 1 && matrix[i-sign][j+1] == 0 {
				matrix[i][j] = 0
				matrix[i-sign][j+1] = 1
			}
		}
	}*///update

	for j := visinaKanvasa-3; j > 0; j-- {
		for i := 1; i < sirinaKanvasa-1; i++ {
			if matrix[i][j] == 1 {
				if matrix[i][j+1] == 0 {
					matrix[i][j], matrix[i][j+1] = 0, 1
					//ako moze pasti direkt neka padne
				} else {
					//u suprotnom gleda moze li dijagonalu
					var sgn int
					if rand.Intn(2) == 1 {
						sgn = 1
					} else {
						sgn = -1
					}//nasumice biramo na koju stranu prvo ide da izbegnemo pristrasno padanje

					if (matrix[i+sgn][j+1] == 0) && (i+sgn > 0) && (i+sgn < sirinaKanvasa) {
						matrix[i][j], matrix[i+sgn][j+1] = 0, 1
					} else if (matrix[i-sgn][j+1] == 0) && (i+sgn > 0) && (i+sgn < sirinaKanvasa) {
						matrix[i][j], matrix[i-sgn][j+1] = 0, 1
					} 
				}
			}
		}
	}

}

func render(matrix [][]int, surface *sdl.Surface){

	/*for i := 0; i < 800-3+1; i++ {
		for j := 600-2; j >= 0; j-- {
			rect := sdl.Rect{int32(i), int32(j), 1, 1}
			if(matrix[i][j] == 0) {
				surface.FillRect(&rect, 0)
			}
			if(matrix[i][j] == 1) {
				surface.FillRect(&rect, 0xffff66)
			}
		}
	}*///render

	for i := 0; i < sirinaKanvasa; i++ {
		for j := 0; j < visinaKanvasa; j++ {
			rect := sdl.Rect{int32(i*brojPikselaPoCestici), int32(j*brojPikselaPoCestici), brojPikselaPoCestici,brojPikselaPoCestici}
			switch matrix[i][j]{
				case 1:
					surface.FillRect(&rect, 0xffff66)
				default:
					surface.FillRect(&rect, 0)
			}
		}
	}
	for i := 0; i < sirinaKanvasa; i++ {
		rect := sdl.Rect{int32(i*brojPikselaPoCestici), 0, brojPikselaPoCestici, brojPikselaPoCestici}
		surface.FillRect(&rect, 0x663366)
		rect = sdl.Rect{int32(i*brojPikselaPoCestici), (visinaKanvasa-1)*brojPikselaPoCestici, brojPikselaPoCestici, brojPikselaPoCestici}
		surface.FillRect(&rect, 0x663366)
	}
	for j := 0; j < visinaKanvasa; j++ {
		rect := sdl.Rect{0, int32(j*brojPikselaPoCestici), brojPikselaPoCestici, brojPikselaPoCestici}
		surface.FillRect(&rect, 0x663366)
		rect = sdl.Rect{(sirinaKanvasa-1)*brojPikselaPoCestici, int32(j*brojPikselaPoCestici), brojPikselaPoCestici, brojPikselaPoCestici}
		surface.FillRect(&rect, 0x663366)
	}

}

func matrixToSlice(matrix [sirinaKanvasa][visinaKanvasa]int) [][]int {
	//pakao brate moj u hristu valjda je ovo najbolji nacin?
	//ako nije promenicemo, nije kao da struktura celog projekta zavisi od ovoga (:

	slajs := make([][]int, len(matrix))

	for i := 0; i < len(matrix); i++ {
		kolona := make([]int, len(matrix[i]))
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