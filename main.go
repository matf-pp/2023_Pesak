package main

import 
(
	"github.com/veandco/go-sdl2/sdl"
	"math/rand"
	"math"
//	"time"
	"fmt"
//	"app/mat"
)

type Materijal int

const (
	Zid Materijal = -1
	Prazno Materijal = 0
	Pesak Materijal = 1
	Voda Materijal = 2
	Metal Materijal = 3
	//itd
)

var boja = map[Materijal]uint32{
	Zid : 0xffffff,
	Prazno : 0x000000,
	Pesak : 0xffff66,
	Voda : 0x3333ff,
	Metal : 0x33334b,
//	Kamen : 0x666666,
}

const sirinaKanvasa, visinaKanvasa = 240, 144
const brojPikselaPoCestici = 4
//golang nema makroe pa koristimo globalne konst (za sada?) -s

var trenutniMat Materijal = Pesak
var velicinaKursora int32 = 0
var maxKursor int32 = 8

var keystates = sdl.GetKeyboardState()

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

	var matrix[sirinaKanvasa][visinaKanvasa]Materijal
	slajs := matrixToSlice(matrix)

	//TODO izdvojiti ove dve ruzne petlje van mejna? `void zidaj(slajs);` npr -s
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

func clampCoords(x int32, y int32) (int32, int32) {
	//osigurava da tacka ne izleti iz ekrana sto se desava u raznim slucajevima -s
	// radi isto kao što piše iznad, samo u jednoj liniji /limun
	return int32(math.Min(math.Max(float64(x), 0), sirinaKanvasa - 1)), int32(math.Min(math.Max(float64(y), 0), visinaKanvasa - 1))
}

func brush(matrix [][]Materijal, x int32, y int32, state uint32) {
	// zamenio redosled if-a i for-a /limun
	// dole se povećavala veličina kursora za 2, a ovde delila sa 2, pa sam sklonio /2 i +2 na +1 /limun
	for i := -velicinaKursora; i <= velicinaKursora; i++ {
		for j := -velicinaKursora; j <= velicinaKursora; j++{
			tx, ty := clampCoords(x/brojPikselaPoCestici+i, y/brojPikselaPoCestici+j)
			if state == 1 && matrix[tx][ty] == Prazno {
				matrix[tx][ty] = trenutniMat
			}
			if state == 4 && matrix[tx][ty] != Zid {
				matrix[tx][ty] = Prazno
			}
		}
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
				// kako si to zamisljao lepse? -s
				if keystates[sdl.SCANCODE_ESCAPE] != 0 {
					running = false
				}
				if keystates[sdl.SCANCODE_0] != 0 {
					trenutniMat = Prazno
				}
				if keystates[sdl.SCANCODE_1] != 0 {
					trenutniMat = Pesak
				}
				if keystates[sdl.SCANCODE_2] != 0 {
					trenutniMat = Voda
				}
				if keystates[sdl.SCANCODE_3] != 0 {
					trenutniMat = Metal
				}
				if keystates[sdl.SCANCODE_DOWN] != 0 {
					// povećavalo se za 2, a onda delilo sa 2 u petlji gore, pa sam smanjio na -1 i skinuo /2 /limun
					if velicinaKursora > 0 {
						velicinaKursora = velicinaKursora - 1
					}
				}
				if keystates[sdl.SCANCODE_UP] != 0 {
					if velicinaKursora < maxKursor {
						velicinaKursora = velicinaKursora + 1
					}
				}
			default:
				//null
		}
	}

	var x, y int32
	var state uint32
	x, y, state = sdl.GetMouseState()
	fmt.Printf("x: %d\t", x)
	fmt.Printf("y: %d\t", y)
	fmt.Printf("xpx: %d\t", x/brojPikselaPoCestici)
	fmt.Printf("ypx: %d\t", y/brojPikselaPoCestici)
	fmt.Printf("mb: %d\t", state)
	fmt.Printf("materijal: %d\t", trenutniMat)
	fmt.Printf("velicina: %d\n", velicinaKursora)

	brush(matrix, x, y, state)	

	return running

}

func update(matrix [][]Materijal) {
	for j := visinaKanvasa-3; j > 0; j-- {
		for i := 1; i < sirinaKanvasa-1; i++ {

			// pesak
			if matrix[i][j] == Pesak {
				//flutter := rand.Intn(2)
				//predlazem da lelujanje ostavimo za malo kasnije da pojednostavimo ovo dok ne sredimo lepo -s
				if matrix[i][j+1] == Prazno || matrix[i][j+1] == Voda{
//					if flutter == 1 {
						matrix[i][j], matrix[i][j+1] = matrix[i][j+1], matrix[i][j]
//					}//	ovde prvi put a nadalje puno puta postaje ocigledna potreba za nekom swap
					// funkcijom ali ionako ovome slede radikalne promene, ne bih trosio vreme
					// sada na to... -s 
				} else {
					var sgn int
					if rand.Intn(2) == 1 {
						sgn = 1
					} else {
						sgn = -1
					}//nasumice biramo na koju stranu prvo ide da izbegnemo pristrasno padanje

					if (matrix[i+sgn][j+1] == Prazno || matrix[i+sgn][j+1] == Voda) && (i+sgn > 0) && (i+sgn < sirinaKanvasa) {
						matrix[i][j], matrix[i+sgn][j+1] = matrix[i+sgn][j+1], matrix[i][j]
					} else if (matrix[i-sgn][j+1] == Prazno || matrix[i-sgn][j+1] == Voda) && (i+sgn > 0) && (i+sgn < sirinaKanvasa) {
						matrix[i][j], matrix[i-sgn][j+1] = matrix[i-sgn][j+1], matrix[i][j]
					} 
				}
			}
			// voda
			// ovo je teže nego što smo mislili
			//mozda teze nego sto si /ti/ mislio -s
			if matrix[i][j] == Voda {
//				flutter := rand.Intn(2)
				if matrix[i][j+1] == Prazno {
//					if flutter == 1 {
						matrix[i][j], matrix[i][j+1] = Prazno, Voda
//					}
				} else {
					var sgn int = rand.Intn(3)
					sgn = sgn - 1
					// -1, 0, 1

					if (matrix[i+sgn][j+1] == Prazno) && (i+sgn > 0) && (i+sgn < sirinaKanvasa) {
						matrix[i][j], matrix[i+sgn][j+1] = Prazno, Voda
					} else if (matrix[i-sgn][j+1] == Prazno) && (i-sgn > 0) && (i-sgn < sirinaKanvasa) {
						matrix[i][j], matrix[i-sgn][j+1] = Prazno, Voda
					} else if (matrix[i+sgn][j] == Prazno) && (i+sgn > 0) && (i+sgn < sirinaKanvasa) {
						matrix[i][j], matrix[i+sgn][j] = Prazno, Voda
					} else if (matrix[i-sgn][j] == Prazno) && (i-sgn > 0) && (i-sgn < sirinaKanvasa) {
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
			surface.FillRect(&rect, boja[matrix[i][j]])
		}
	}
}

func matrixToSlice(matrix [sirinaKanvasa][visinaKanvasa]Materijal) [][]Materijal {

	slajs := make([][]Materijal, len(matrix))

	for i := 0; i < len(matrix); i++ {
		kolona := make([]Materijal, len(matrix[i]))
		for j := 0; j < len(matrix[i]); j++ {
			kolona[j] = matrix[i][j]
		}
		slajs[i] = kolona
	}

	return slajs
	
}
