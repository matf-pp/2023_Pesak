package main

import (
	"github.com/veandco/go-sdl2/sdl"
	"math/rand"
	"math"
	//	"time"
	"fmt"

	"main/mat"
)

var boja = mat.Boja

const sirinaKanvasa, visinaKanvasa = 240, 144
const brojPikselaPoCestici = 6

// const brojPikselaPoCestici = 8
// TODO automatska detekcija rezolucije ekrana
// njanja: može ali na početku maina i ne ovde tkd vidi šta ćeš sa njom npr možemo da pomerimo sve constove u mejn nz je l bi se vama to svidelo :creepysmirk
// njanja: tld vidi šestu liniju u mejnu
var sirinaEkrana = 0
var visinaEkrana = 0

// njanja: mislim da je ovo korisno bar meni jeste
const sirinaUIMargine = 10
const visinaUIMargine = 20
const sirinaDugmeta = 60
const visinaDugmeta = 30
const marginaZaGumbad = 2*sirinaUIMargine + sirinaDugmeta
const sirinaProzora = sirinaKanvasa*brojPikselaPoCestici + marginaZaGumbad
const visinaProzora = visinaKanvasa * brojPikselaPoCestici

var trenutniMat mat.Materijal = mat.Pesak

// var velicinaKursora int32 = 0
var velicinaKursora int32 = 4
var maxKursor int32 = 32
var pause bool = false

// njanja: ovo mi treba da bih znao gde da displejujem konture kursora dok se ne pomera
var kursorPoslednjiX = int32(sirinaEkrana / 2)
var kursorPoslednjiY = int32(sirinaEkrana / 2)

var keystates = sdl.GetKeyboardState()

func main() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	// njanja: automatska detekcija ekrana
	var desktop sdl.DisplayMode
	desktop, err := sdl.GetDesktopDisplayMode(0)
	sirinaEkrana = int(desktop.W)
	visinaEkrana = int(desktop.H)

	// njanja: dodao marginu za gumbad
	window, err := sdl.CreateWindow("test", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		sirinaKanvasa*brojPikselaPoCestici+marginaZaGumbad, visinaKanvasa*brojPikselaPoCestici, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	surface, err := window.GetSurface()
	if err != nil {
		panic(err)
	}
	surface.FillRect(nil, 0)

	// njanja: proverite grešku ako hoćete štreberi
	renderer, _ := window.GetRenderer()

	// zašto bafer? /limun
	var matrix [sirinaKanvasa][visinaKanvasa]mat.Cestica
	slajs := matrixToSlice(matrix)
	var bafer [sirinaKanvasa][visinaKanvasa]mat.Cestica
	bajs := matrixToSlice(bafer)

	//TODO izdvojiti ove dve ruzne petlje van mejna? `void mat.Zidaj(slajs);` npr -s
	// zašto bafer nema zidove? /limun
	slajs = zazidajMatricu(slajs)
	bajs = zazidajMatricu(bajs)

	running := true
	for running {
		running = pollEvents(slajs, bajs)
		if !pause {
			updateCanvas(slajs, bajs)
		}
		render(slajs, surface)

		// njanja: ovo renderuje gumbad za sve materijale
		for i, _ := range boja {
			gumb := sdl.Rect{int32(sirinaProzora - sirinaUIMargine - sirinaDugmeta), int32(visinaUIMargine + i*(visinaDugmeta+visinaUIMargine)), sirinaDugmeta, visinaDugmeta}
			surface.FillRect(&gumb, boja[i])
		}

		// njanja: ovo renderuje dodatnu gumbad
		// možda može u fju ne znam
		plejGumb := sdl.Rect{int32(sirinaProzora - sirinaUIMargine - sirinaDugmeta), int32(visinaProzora - 3*visinaUIMargine - 3*visinaDugmeta), sirinaDugmeta, visinaDugmeta}
		if pause {
			surface.FillRect(&plejGumb, 0x00ff00)
		} else {
			surface.FillRect(&plejGumb, 0xffa500)
		}

		sejvGumb := sdl.Rect{int32(sirinaProzora - sirinaUIMargine - sirinaDugmeta), int32(visinaProzora - 2*visinaUIMargine - 2*visinaDugmeta), sirinaDugmeta, visinaDugmeta}
		surface.FillRect(&sejvGumb, 0x0000ff)

		resetGumb := sdl.Rect{int32(sirinaProzora - sirinaUIMargine - sirinaDugmeta), int32(visinaProzora - visinaUIMargine - visinaDugmeta), sirinaDugmeta, visinaDugmeta}
		surface.FillRect(&resetGumb, 0xff0000)

		window.UpdateSurface()

		// njanja: ovo renderuje krug oko četkice, ne može u brush fju jer se ona ne zove u svakom frejmu
		// takođe kursor flikeruje jer bude na kratko izbrisan kad se pozove updatesurface par linija iznad pa dok ne bude ponovo nacrtan
		// ako mene pitate mislim da daje odličan retro look
		// pored toga što treperi takođe čini da uvidim problem da se displejuje preko UIja a slobodno probajte da selektujete materijal baš baš velikom četkicom
		cetkica := sdl.Rect{kursorPoslednjiX - velicinaKursora*brojPikselaPoCestici, kursorPoslednjiY - velicinaKursora*brojPikselaPoCestici, int32(2 * velicinaKursora * brojPikselaPoCestici), int32(2 * velicinaKursora * brojPikselaPoCestici)}
		renderer.SetDrawColor(255, 255, 255, 255)
		renderer.DrawRect(&cetkica)
		renderer.Present()
	}

}

// jesi li na ovo mislio? /limun
// "slice" da se ne bi mešao sa postojećim "slajs" /limun
func zazidajMatricu(slice [][]mat.Cestica) [][]mat.Cestica {
	for i := 0; i < sirinaKanvasa; i++ {
		slice[i][0] = mat.NewCestica(mat.Zid)
		slice[i][visinaKanvasa-1] = mat.NewCestica(mat.Zid)
	}
	for j := 0; j < visinaKanvasa; j++ {
		slice[0][j] = mat.NewCestica(mat.Zid)
		slice[sirinaKanvasa-1][j] = mat.NewCestica(mat.Zid)
	}

	return slice
}

func clampCoords(x int32, y int32) (int32, int32) {
	//osigurava da tacka ne izleti iz ekrana sto se desava u raznim slucajevima -s
	// radi isto kao što piše iznad, samo u jednoj liniji /limun
	// lep funkcionalni pristup, kolega -s
	return int32(math.Min(math.Max(float64(x), 0), sirinaKanvasa-1)), int32(math.Min(math.Max(float64(y), 0), visinaKanvasa-1))
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
			if (state == 4 || (state == 1 && trenutniMat == mat.Prazno)) && matrix[tx][ty].Materijal != mat.Zid && matrix[tx][ty].Materijal != mat.Prazno {
				matrix[tx][ty] = mat.NewCestica(mat.Prazno)
				bafer[tx][ty] = matrix[tx][ty]
			}
		}
	}
}

func pollEvents(matrix [][]mat.Cestica, bafer [][]mat.Cestica) bool {
	running := true
	keystates = sdl.GetKeyboardState()

	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		// njanja: switch event.(type) -> switch t := event.(type)
		switch t := event.(type) {
		case *sdl.QuitEvent:
			running = false
			break

		case *sdl.KeyboardEvent:
			// može ovo lepše #TODO
			// kako si to zamisljao lepse? -s
			// ignorisite ovo, mislim da pricam sam sa sobom (nz ciji je prvi kom)
			// prvi kom je moj, tako da ne moraš da ga ignorišeš -s
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

		// njanja: za ovo mi je potreban diskretan klik a ne frejm sa dugmetom dole
		// p.s. hoćemo da ostavimo komentare ristoviću da ih vidi
		// paa, barem imajte naznaku za svaku liniju čiji je čiji /limun
		/*
			ili prosto napravi višelinijski komentar pa na početku/kraju stavi ime
			/limun
		*/
		case *sdl.MouseButtonEvent:
			if t.State == sdl.PRESSED {
				proveriPritisakNaGumb(t.X, t.Y)
			}

		default:
			//null
		}
	}

	var x, y int32
	var state uint32
	x, y, state = sdl.GetMouseState()
	kursorPoslednjiX = x
	kursorPoslednjiY = y

	// njanja: ovo sam zakomentarisao da me ne smara
	// ovo sam odkomentarisao da me smara /limun
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

func proveriPritisakNaGumb(x int32, y int32) {
	//njanja: ovo je detekcija klika na gumb
	if x > sirinaProzora-marginaZaGumbad+sirinaUIMargine && x < sirinaProzora-sirinaUIMargine {
		// materijali
		if y < (visinaUIMargine+visinaDugmeta)*int32(len(boja)-1) && y%(visinaUIMargine+visinaDugmeta) > visinaUIMargine {
			trenutniMat = mat.Materijal(y / (visinaUIMargine + visinaDugmeta))
		}

		// njanja: hardkodovan broj specijalnih dugmića
		if y > visinaProzora-3*(visinaDugmeta+visinaUIMargine) && y < visinaProzora-3*(visinaDugmeta+visinaUIMargine)+visinaDugmeta {
			pause = !pause
		}
		if y > visinaProzora-2*(visinaDugmeta+visinaUIMargine) && y < visinaProzora-2*(visinaDugmeta+visinaUIMargine)+visinaDugmeta {
			// sejv
		}
		if y > visinaProzora-1*(visinaDugmeta+visinaUIMargine) && y < visinaProzora-1*(visinaDugmeta+visinaUIMargine)+visinaDugmeta {
			// reset
		}
	}
}

func updateCanvas(matrix [][]mat.Cestica, bafer [][]mat.Cestica) {
	for j := 1; j < visinaKanvasa-2; j++ {
		for i := 1; i < sirinaKanvasa-2; i++ {
			mat.Update(matrix, bafer, i, j)

						// mat.Pesak
			  			if matrix[i][j].Materijal == mat.Pesak {
			  				//flutter := rand.Intn(2)
			  				//predlazem da lelujanje ostavimo za malo kasnije da pojednostavimo ovo dok ne sredimo lepo -s
			  				if matrix[i][j+1].Materijal == mat.Prazno || matrix[i][j+1].Materijal == mat.Voda{
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

			  					if (matrix[i+sgn][j+1].Materijal == mat.Prazno || matrix[i+sgn][j+1].Materijal == mat.Voda) && (i+sgn > 0) && (i+sgn < sirinaKanvasa) {
			  						matrix[i][j], matrix[i+sgn][j+1] = matrix[i+sgn][j+1], matrix[i][j]
			  					} else if (matrix[i-sgn][j+1].Materijal == mat.Prazno || matrix[i-sgn][j+1].Materijal == mat.Voda) && (i+sgn > 0) && (i+sgn < sirinaKanvasa) {
			  						matrix[i][j], matrix[i-sgn][j+1] = matrix[i-sgn][j+1], matrix[i][j]
			  					}
			  				}
			  			}
			  			// mat.Voda
			  			// ovo je teže nego što smo mislili
			  			//mozda teze nego sto si /ti/ mislio -s
			  			if matrix[i][j].Materijal == mat.Voda {
			  //				flutter := rand.Intn(2)
			  				if matrix[i][j+1].Materijal == mat.Prazno {
			  //					if flutter == 1 {
			  						matrix[i][j].Materijal, matrix[i][j+1].Materijal = mat.Prazno, mat.Voda
			  //					}
			  				} else {
			  					var sgn int = rand.Intn(3)
			  					sgn = sgn - 1
			  					// -1, 0, 1

			  					if (matrix[i+sgn][j+1].Materijal == mat.Prazno) && (i+sgn > 0) && (i+sgn < sirinaKanvasa) {
			  						matrix[i][j].Materijal, matrix[i+sgn][j+1].Materijal = mat.Prazno, mat.Voda
			  					} else if (matrix[i-sgn][j+1].Materijal == mat.Prazno) && (i-sgn > 0) && (i-sgn < sirinaKanvasa) {
			  						matrix[i][j].Materijal, matrix[i-sgn][j+1].Materijal = mat.Prazno, mat.Voda
			  					} else if (matrix[i+sgn][j].Materijal == mat.Prazno) && (i+sgn > 0) && (i+sgn < sirinaKanvasa) {
			  						matrix[i][j].Materijal, matrix[i+sgn][j].Materijal = mat.Prazno, mat.Voda
			  					} else if (matrix[i-sgn][j].Materijal == mat.Prazno) && (i-sgn > 0) && (i-sgn < sirinaKanvasa) {
			  						matrix[i][j].Materijal, matrix[i-sgn][j].Materijal = mat.Prazno, mat.Voda
			  					}
			  				}
			  			}
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
			rect := sdl.Rect{int32(i * brojPikselaPoCestici), int32(j * brojPikselaPoCestici), brojPikselaPoCestici, brojPikselaPoCestici}
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
