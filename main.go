package main

import (
	"fmt"
	"main/mat"
	"math"
	"strconv"

	"github.com/veandco/go-sdl2/sdl"

	"time"

	"github.com/hugolgst/rich-go/client"
)

// njanja: ovo je loša praksa majmuni
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

var velicinaKursora int32 = 0

// var velicinaKursora int32 = 4
var maxKursor int32 = 32
var pause bool = false

var tempMode bool = false
var tempColorMultiplier float64 = 3

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

	// njanja: diskord integracija smislićemo šta ćemo s njom
	err = client.Login("1100118057147437207")
	if err != nil {
		panic(err)
	}

	now := time.Now()
	client.SetActivity(client.Activity{
		State:      "bleja",
		Details:    "kontemplira pesak",
		LargeImage: "bleja",
		LargeText:  "je l se učitalo ovo",
		Timestamps: &client.Timestamps{
			Start: &now,
		},
		Buttons: []*client.Button{
			&client.Button{
				Label: "priključi se",
				Url:   "https://github.com/matf-pp/2023_Pesak",
			},
		},
	})

	// zašto bafer? /limun
	var matrix [sirinaKanvasa][visinaKanvasa]mat.Cestica
	slajs := matrixToSlice(matrix)
	var bafer [sirinaKanvasa][visinaKanvasa]mat.Cestica
	bajs := matrixToSlice(bafer)

	//TODO izdvojiti ove dve ruzne petlje van mejna? `void mat.Zidaj(slajs);` npr -s
	// zašto bafer nema zidove? /limun
	// njanja: obrišite komentar ako ste izdvojili petlje slatkiši
	slajs = zazidajMatricu(slajs)

	// enejblujemo dropove, stavite ovo gde hoćete
	sdl.EventState(sdl.DROPFILE, sdl.ENABLE)

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
		// stavi ti u f-ju pošto razumeš o čemu se radi /limun
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

// dodaje zidove oko matrice /limun
func zazidajMatricu(matrix [][]mat.Cestica) [][]mat.Cestica {
	for i := 0; i < sirinaKanvasa; i++ {
		matrix[i][0] = mat.NewCestica(mat.Zid)
		matrix[i][visinaKanvasa-1] = mat.NewCestica(mat.Zid)
	}
	for j := 0; j < visinaKanvasa; j++ {
		matrix[0][j] = mat.NewCestica(mat.Zid)
		matrix[sirinaKanvasa-1][j] = mat.NewCestica(mat.Zid)
	}

	return matrix
}

// menja i vraća (x, y) koordinate tako da se nalaze na ekranu /limun
func clampCoords(x int32, y int32) (int32, int32) {
	return int32(math.Min(math.Max(float64(x), 0), sirinaKanvasa-1)), int32(math.Min(math.Max(float64(y), 0), visinaKanvasa-1))
}

func brush(matrix [][]mat.Cestica, bafer [][]mat.Cestica, x int32, y int32, state uint32) {
	//TODO za srednji klik da uzme materijal na koj mis trenutno pokazuje i postavi ga kao trenutni
	//ukoliko nije u pitanju Zid ili Prazno. Nije mi pri ruci mis, mrzi me da trazim koj je to stejt -s
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
			if keystates[sdl.SCANCODE_UP] != 0 {
				if velicinaKursora < maxKursor {
					velicinaKursora = velicinaKursora + 1
				}
			}
			if keystates[sdl.SCANCODE_P] != 0 {
				pause = !pause
			}
			if keystates[sdl.SCANCODE_T] != 0 {
				tempMode = !tempMode
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
				proveriPritisakNaGumb(matrix, bafer, t.X, t.Y)
			}

		// drag and drop slike je odmah učitava
		case *sdl.DropEvent:
			dropEvent := event.(*sdl.DropEvent)
			if dropEvent.Type == sdl.DROPFILE {
				filePath := string(dropEvent.File)
				obradiSliku(filePath, sirinaKanvasa, visinaKanvasa, matrix, bafer)
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

	fmt.Printf("x: %d\t", x)
	fmt.Printf("y: %d\t", y)
	fmt.Printf("xpx: %d\t", x/brojPikselaPoCestici)
	fmt.Printf("ypx: %d\t", y/brojPikselaPoCestici)
	fmt.Printf("mb: %d\t", state)
	fmt.Printf("mat.Materijal: %d\t", trenutniMat)
	fmt.Printf("velicina: %d\t", velicinaKursora)
	fmt.Printf("pauza: %t\n", pause)
	fmt.Printf("tempMode: %t\n", tempMode)

	brush(matrix, bafer, x, y, state)

	return running

}

func proveriPritisakNaGumb(matrix, bafer [][]mat.Cestica, x, y int32) {
	//njanja: ovo je detekcija klika na gumb
	if x > sirinaProzora-marginaZaGumbad+sirinaUIMargine && x < sirinaProzora-sirinaUIMargine {
		// njanja: TODO namestiti da se ređaju u više kolona ako baš mora
		// materijali
		if y < (visinaUIMargine+visinaDugmeta)*int32(len(boja)-1) && y%(visinaUIMargine+visinaDugmeta) > visinaUIMargine {
			trenutniMat = mat.Materijal(y / (visinaUIMargine + visinaDugmeta))
		}

		// njanja: hardkodovan broj specijalnih dugmića
		// PAUZA
		if y > visinaProzora-3*(visinaDugmeta+visinaUIMargine) && y < visinaProzora-3*(visinaDugmeta+visinaUIMargine)+visinaDugmeta {
			pause = !pause
		}
		// njanja: neki medvedić dobrog srca nek rinejmuje ovu funkciju da prati konvenciju kamilju ili ću ja sutra dobro ajde TODO
		// SEJV
		if y > visinaProzora-2*(visinaDugmeta+visinaUIMargine) && y < visinaProzora-2*(visinaDugmeta+visinaUIMargine)+visinaDugmeta {
			save_image(matrix, sirinaKanvasa, visinaKanvasa)
		}
		// njanja: nz je l ovo najpametniji način ali radi
		// RESET
		if y > visinaProzora-1*(visinaDugmeta+visinaUIMargine) && y < visinaProzora-1*(visinaDugmeta+visinaUIMargine)+visinaDugmeta {
			for j := 0; j < visinaKanvasa; j++ {
				for i := 0; i < sirinaKanvasa; i++ {
					matrix[i][j] = mat.NewCestica(mat.Prazno)
					bafer[i][j] = matrix[i][j]
				}
			}
			zazidajMatricu(matrix)
			zazidajMatricu(bafer)
		}
	}
}
func updateCanvas(matrix [][]mat.Cestica, bafer [][]mat.Cestica) {
	/* Problem rešen! Imali smo bag gde se čestice ne iscrtavaju kako treba skroz
	   dole i skroz desno; pisalo je visinaKanvasa-2 a ne visinaKanvasa-1
	   /limun
	*/
	for j := 1; j < visinaKanvasa-1; j++ {
		for i := 1; i < sirinaKanvasa-1; i++ {
			mat.Update(matrix, bafer, i, j)
		}
	}
	for j := 1; j < visinaKanvasa-1; j++ {
		for i := 1; i < sirinaKanvasa-1; i++ {
			matrix[i][j] = bafer[i][j]
		}
	}
}

func render(matrix [][]mat.Cestica, surface *sdl.Surface) {
	for i := 0; i < sirinaKanvasa; i++ {
		for j := 0; j < visinaKanvasa; j++ {
			rect := sdl.Rect{int32(i * brojPikselaPoCestici), int32(j * brojPikselaPoCestici), brojPikselaPoCestici, brojPikselaPoCestici}
			if !tempMode {
				surface.FillRect(&rect, boja[matrix[i][j].Materijal])
			} else {
				bojaTemp := izracunajTempBoju(matrix[i][j].Temperatura)
				surface.FillRect(&rect, bojaTemp)
			}
		}
	}
}

func izracunajTempBoju(temp float64) uint32 {
	temp *= tempColorMultiplier
	if temp > 0 {
		temp = math.Min(float64(temp), 255)
		temp = float64(int32(256-temp)<<8) + (255 << 16)
	} else if temp < 0 {
		temp *= -3
		temp = math.Min(float64(temp), 255)
		temp = float64(int32(256-temp)<<8) + 255
	}

	hexadeca := strconv.FormatUint(uint64(temp), 16)
	tempBoja, err := strconv.ParseUint(hexadeca, 16, 32)
	if err != nil {
		panic(err)
	}

	return uint32(tempBoja)
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
