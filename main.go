package main

import 
(
	"github.com/veandco/go-sdl2/sdl"
//	"math/rand"
	"math"
//	"time"
	"fmt"
	"main/mat"
)

var boja = mat.Boja

const sirinaKanvasa, visinaKanvasa = 240, 144
//const brojPikselaPoCestici = 4
const brojPikselaPoCestici = 8
//TODO automatska detekcija rezolucije ekrana

var trenutniMat mat.Materijal = mat.Pesak
//var velicinaKursora int32 = 0
var velicinaKursora int32 = 4
var maxKursor int32 = 32
var pause bool = true

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

	var matrix[sirinaKanvasa][visinaKanvasa] mat.Cestica
	slajs := matrixToSlice(matrix)
	var bafer[sirinaKanvasa][visinaKanvasa] mat.Cestica
	bajs := matrixToSlice(bafer)

	//TODO izdvojiti ove dve ruzne petlje van mejna? `void mat.Zidaj(slajs);` npr -s
	for i := 0; i < sirinaKanvasa; i++ {
		slajs[i][0] = mat.NewCestica(mat.Zid)
		slajs[i][visinaKanvasa-1] = mat.NewCestica(mat.Zid)
	}
	for j := 0; j < visinaKanvasa; j++ {
		slajs[0][j] = mat.NewCestica(mat.Zid)
		slajs[sirinaKanvasa-1][j] = mat.NewCestica(mat.Zid)
	}

	running := true
	for running {
		running = pollEvents(slajs, bajs)
		if !pause {
			updateCanvas(slajs, bajs)
		}
		render(slajs, surface)

		window.UpdateSurface()
	}

}

func clampCoords(x int32, y int32) (int32, int32) {
	//osigurava da tacka ne izleti iz ekrana sto se desava u raznim slucajevima -s
	// radi isto kao što piše iznad, samo u jednoj liniji /limun
	// lep funkcionalni pristup, kolega -s
	return int32(math.Min(math.Max(float64(x), 0), sirinaKanvasa - 1)), int32(math.Min(math.Max(float64(y), 0), visinaKanvasa - 1))
}

func brush(matrix [][]mat.Cestica, bafer [][]mat.Cestica, x int32, y int32, state uint32) {
	//TODO za srednji klik da uzme materijal na koj mis trenutno pokazuje i postavi ga kao trenutni
	//ukoliko nije u pitanju Zid ili Prazno. Nije mi pri ruci mis, mrzi me da trazim koj je to stejt -s
	// zamenio redosled if-a i for-a /limun
	// dole se povećavala veličina kursora za 2, a ovde delila sa 2, pa sam sklonio /2 i +2 na +1 /limun 
	for i := -velicinaKursora; i <= velicinaKursora; i++ {
		for j := -velicinaKursora; j <= velicinaKursora; j++ {
			tx, ty := clampCoords(x/brojPikselaPoCestici+i, y/brojPikselaPoCestici+j)
			if state == 1 && matrix[tx][ty].Materijal == mat.Prazno {
				matrix[tx][ty] = mat.NewCestica(trenutniMat)
				bafer[tx][ty] = matrix[tx][ty]

			}
			if (state == 4 || (state == 1 && trenutniMat == mat.Prazno)) && matrix[tx][ty].Materijal != mat.Zid && matrix[tx][ty].Materijal!= mat.Prazno{
				matrix[tx][ty] = mat.NewCestica(mat.Prazno)
				bafer[tx][ty] = matrix[tx][ty]
			}
		}
	}
}

func pollEvents(matrix [][]mat.Cestica, bafer [][]mat.Cestica) bool {
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
				// ignorisite ovo, mislim da pricam sam sa sobom (nz ciji je prvi kom)
				if keystates[sdl.SCANCODE_ESCAPE] != 0 {
					running = false
				}
				if keystates[sdl.SCANCODE_0] != 0 {
					trenutniMat = mat.Prazno
				}
				if keystates[sdl.SCANCODE_1] != 0 {
					trenutniMat = mat.Pesak
				}
				if keystates[sdl.SCANCODE_2] != 0 {
					trenutniMat = mat.Voda
				}
				if keystates[sdl.SCANCODE_3] != 0 {
					trenutniMat = mat.Metal
				}
				if keystates[sdl.SCANCODE_4] != 0 {
					trenutniMat = mat.Kamen
				}
				if keystates[sdl.SCANCODE_5] != 0 {
					trenutniMat = mat.Lava
				}
				if keystates[sdl.SCANCODE_6] != 0 {
					trenutniMat = mat.Led
				}
				if keystates[sdl.SCANCODE_7] != 0 {
					trenutniMat = mat.Para
				}
				if keystates[sdl.SCANCODE_DOWN] != 0 {
					if velicinaKursora > 0 {
						velicinaKursora = velicinaKursora - 1
					}
				}
				if keystates[sdl.SCANCODE_P] != 0 {
					pause = !pause
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
	fmt.Printf("mat.Materijal: %d\t", trenutniMat)
	fmt.Printf("velicina: %d\t", velicinaKursora)
	fmt.Printf("pauza: %t\n", pause)

	brush(matrix, bafer, x, y, state)	

	return running

}

func updateCanvas(matrix [][]mat.Cestica, bafer [][]mat.Cestica) {
	for j := 1; j < visinaKanvasa-2; j++ {
		for i := 1; i < sirinaKanvasa-2; i++ {
			mat.Update(matrix, bafer, i, j)


/**			// mat.Pesak
			if matrix[i][j] == mat.Pesak {
				//flutter := rand.Intn(2)
				//predlazem da lelujanje ostavimo za malo kasnije da pojednostavimo ovo dok ne sredimo lepo -s
				if matrix[i][j+1] == mat.Prazno || matrix[i][j+1] == mat.Voda{
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

					if (matrix[i+sgn][j+1] == mat.Prazno || matrix[i+sgn][j+1] == mat.Voda) && (i+sgn > 0) && (i+sgn < sirinaKanvasa) {
						matrix[i][j], matrix[i+sgn][j+1] = matrix[i+sgn][j+1], matrix[i][j]
					} else if (matrix[i-sgn][j+1] == mat.Prazno || matrix[i-sgn][j+1] == mat.Voda) && (i+sgn > 0) && (i+sgn < sirinaKanvasa) {
						matrix[i][j], matrix[i-sgn][j+1] = matrix[i-sgn][j+1], matrix[i][j]
					} 
				}
			}
			// mat.Voda
			// ovo je teže nego što smo mislili
			//mozda teze nego sto si /ti/ mislio -s
			if matrix[i][j] == mat.Voda {
//				flutter := rand.Intn(2)
				if matrix[i][j+1] == mat.Prazno {
//					if flutter == 1 {
						matrix[i][j], matrix[i][j+1] = mat.Prazno, mat.Voda
//					}
				} else {
					var sgn int = rand.Intn(3)
					sgn = sgn - 1
					// -1, 0, 1

					if (matrix[i+sgn][j+1] == mat.Prazno) && (i+sgn > 0) && (i+sgn < sirinaKanvasa) {
						matrix[i][j], matrix[i+sgn][j+1] = mat.Prazno, mat.Voda
					} else if (matrix[i-sgn][j+1] == mat.Prazno) && (i-sgn > 0) && (i-sgn < sirinaKanvasa) {
						matrix[i][j], matrix[i-sgn][j+1] = mat.Prazno, mat.Voda
					} else if (matrix[i+sgn][j] == mat.Prazno) && (i+sgn > 0) && (i+sgn < sirinaKanvasa) {
						matrix[i][j], matrix[i+sgn][j] = mat.Prazno, mat.Voda
					} else if (matrix[i-sgn][j] == mat.Prazno) && (i-sgn > 0) && (i-sgn < sirinaKanvasa) {
						matrix[i][j], matrix[i-sgn][j] = mat.Prazno, mat.Voda
					}
				}
			}/**/
		}
	}

	for j := 1; j < visinaKanvasa-2; j++ {
		for i := 1; i < sirinaKanvasa-2; i++ {
			matrix[i][j] = bafer[i][j]
		}
	}
}

func render(matrix [][]mat.Cestica, surface *sdl.Surface) {
	for i := 0; i < sirinaKanvasa; i++ {
		for j := 0; j < visinaKanvasa; j++ {
			rect := sdl.Rect{int32(i*brojPikselaPoCestici), int32(j*brojPikselaPoCestici), brojPikselaPoCestici,brojPikselaPoCestici}
			surface.FillRect(&rect, boja[matrix[i][j].Materijal])
		}
	}
}

func matrixToSlice(matrix [sirinaKanvasa][visinaKanvasa]mat.Cestica) [][]mat.Cestica {

	slajs := make([][]mat.Cestica, len(matrix))

	for i := 0; i < len(matrix); i++ {
		kolona := make([]mat.Cestica, len(matrix[i]))
		for j := 0; j < len(matrix[i]); j++ {
			kolona[j] = matrix[i][j]
		}
		slajs[i] = kolona
	}

	return slajs
	
}