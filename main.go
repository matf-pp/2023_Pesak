package main

import (
	"fmt"
	"main/mat"
	"math"
	"math/rand"
	"strconv"

	"github.com/fstanis/screenresolution"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

// njanja: todo dodati keyed fields u rectovima da nas ne smaraju žute linije (ili bar mene jer koliko sam shvatio vi ovo bukvalno pišete u noutpedu)
// njanja: todo izbaciti sve ove gluposti koje redovno menjamo u konfig fajl i staviti da bude gitignorovan

// njanja: pazite ovo
const korisnikJeLimun = false

// njanja: ovo je loša praksa majmuni
// e a reci je l si provalio bukvalno je kao `using` u cpp -s
var boja = mat.Boja
var gus = mat.Lambda

//inace mislim da jeovo znak da ove dve mape treba da budu u ovom fajlu
//jer se tamo ionako nekoriste? ako se ne varam -s
// njanja: neeeeee ovaj fajl je već dovoljno veliki

const sirinaKanvasa, visinaKanvasa = 240, 144

// FPS cap, kontam da je zgodno za testiranje staviti neki nizak, 0 = unlimited
var fpsCap = 60

// njanja: ako hoćete da eksperimentišete samo stavite ovo na false ali mislim da nema razloga samo promenite veličinu kanvasa
var autoFitScreen = true

// njanja: ovo sada menjamo tako da ekran prekrije određen procenat korisnikovog ekrana (vidi početak mejna)
var brojPikselaPoCestici int32 = 9000

const sirinaUIMargine = 10
const visinaUIMargine = 10
const sirinaDugmeta = 40
const visinaDugmeta = 20

// njanja: ovo ćemo da menjamo ako treba
var marginaZaGumbad int32 = 2*sirinaUIMargine + sirinaDugmeta

var sirinaProzora = sirinaKanvasa*brojPikselaPoCestici + marginaZaGumbad
var visinaProzora = visinaKanvasa * brojPikselaPoCestici

// njanja: prebačeno na nula jer svakako bude pregaženo u prvom frejmu pa da ne zbunjuje
// a za ovaj int32 dabogda vam nešto nešto
var kursorPoslednjiX = int32(0)
var kursorPoslednjiY = int32(0)

var keystates = sdl.GetKeyboardState()

var trenutniMat mat.Materijal = mat.Pesak

// var velicinaKursora int32 = 4
var velicinaKursora int32 = 8
var maxKursor int32 = 32

var pause bool = false
var tempMode bool = false
var normalMode bool = false
var densityMode bool = false
var textMode bool = true

var tempColorMultiplier int32 = 3

const fontPath = "./assets/Minecraft.ttf"

// njanja: ide unazad verujte mi na reč
const fontSize = 40

func main() {
	// koji procenat ekrana želimo da nam igrica zauzme (probajte da ukucate 0 ili -50 ili tako nešto wild) (spojler: radiće)
	if autoFitScreen {
		brojPikselaPoCestici, sirinaProzora, visinaProzora = fitToScreen(50)
	}

	// njanja: gumb magija ne radi kad nije u mejnu stignite ako hoćete
	// ja sada: https://cdn.discordapp.com/emojis/1068966756556738590.webp
	var brojMaterijala = len(mat.Boja) + 2
	var brojSpecijalnihGumbadi int32 = 3
	var brojGumbadiPoKoloni int32 = visinaProzora/(visinaDugmeta+visinaUIMargine) - (brojSpecijalnihGumbadi)
	var brojKolona int32 = int32(math.Ceil(float64(brojMaterijala) / float64(brojGumbadiPoKoloni)))
	marginaZaGumbad = brojKolona*(sirinaDugmeta+sirinaUIMargine) + sirinaUIMargine
	sirinaProzora += marginaZaGumbad

	// njanja: da vidimo hoće li ovo raditi lepo
	var font *ttf.Font
	var text *sdl.Surface
	err := ttf.Init()
	if err != nil {
		panic(err)
	}
	defer ttf.Quit()

	err = sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		panic(err)
	}
	defer sdl.Quit()

	// njanja: dodao marginu za gumbad
	window, err := sdl.CreateWindow("pesak", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		sirinaProzora, visinaProzora, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	surface, err := window.GetSurface()
	if err != nil {
		panic(err)
	}
	surface.FillRect(nil, 0)

	renderer, err := window.GetRenderer()
	if err != nil {
		panic(err)
	}

	font, err = ttf.OpenFont(fontPath, int(visinaProzora)/fontSize)
	if err != nil {
		panic(err)
	}
	defer font.Close()

	// njanja: diskord integracija oterana u poseban fajl jer je mrvicu čonki
	// ume da pravi probleme i male blek skrinove na početku dok se ne konektuje ili ne tajmautuje pa ga pozivam kao korutinu
	go connectToDiscord()

	var matrica [][]mat.Cestica = napraviSlajs()
	var bafer [][]mat.Cestica = napraviSlajs()

	matrica = zazidajMatricu(matrica)

	// enejblujemo dropove, stavite ovo gde hoćete
	sdl.EventState(sdl.DROPFILE, sdl.ENABLE)

	running := true
	for running {
		// fps counter
		var startTime = sdl.GetTicks64()

		running = pollEvents(matrica, bafer)
		if !pause {
			update(matrica, bafer)
		}
		render(matrica, surface)

		// njanja: ovo renderuje gumbad za sve materijale
		var counter int32 = 1
		for i, _ := range boja {
			gumb := sdl.Rect{int32(sirinaProzora - marginaZaGumbad + ((int32(i)%brojKolona)*(sirinaDugmeta+sirinaUIMargine) + sirinaUIMargine)), int32(visinaUIMargine + int32(i)/brojKolona*(visinaDugmeta+visinaUIMargine)), sirinaDugmeta, visinaDugmeta}
			surface.FillRect(&gumb, boja[i])
			counter++
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

		// njanja: i ovo i četkicu moram nekako da izbacim iz mejna
		// ovo je za tekst
		if textMode {
			var infoText = ""
			// PESAK
			if kursorPoslednjiX < sirinaKanvasa*brojPikselaPoCestici {
				var poslednjiPiksel = matrica[kursorPoslednjiX/brojPikselaPoCestici][kursorPoslednjiY/brojPikselaPoCestici]
				infoText = mat.Ime[poslednjiPiksel.Materijal] + " @ " + fmt.Sprintf("%.2f", float32(poslednjiPiksel.Temperatura/100)) + "C, SekMat: " + mat.Ime[poslednjiPiksel.SekMat] + ", Ticker: " + strconv.Itoa(int(poslednjiPiksel.Ticker))

				// UI
				// njanja: ovo se sigurno i ovde i u pritisku na gumb može mnogo lepše rešiti nekim funkcionalnim pristupom :DERP: longterm ali todo, vrv kad sređujem dva reda gumbeta
			} else {
				if kursorPoslednjiY < (visinaUIMargine+visinaDugmeta)*int32(len(boja)-1) && kursorPoslednjiY%(visinaUIMargine+visinaDugmeta) > visinaUIMargine {
					infoText = mat.Ime[mat.Materijal(kursorPoslednjiY/(visinaUIMargine+visinaDugmeta))]
				}

				// PAUZA
				if kursorPoslednjiY > visinaProzora-3*(visinaDugmeta+visinaUIMargine) && kursorPoslednjiY < visinaProzora-3*(visinaDugmeta+visinaUIMargine)+visinaDugmeta {
					infoText = "Pause"
				}
				// SEJV
				if kursorPoslednjiY > visinaProzora-2*(visinaDugmeta+visinaUIMargine) && kursorPoslednjiY < visinaProzora-2*(visinaDugmeta+visinaUIMargine)+visinaDugmeta {
					infoText = "Save"
				}
				// RESET
				if kursorPoslednjiY > visinaProzora-1*(visinaDugmeta+visinaUIMargine) && kursorPoslednjiY < visinaProzora-1*(visinaDugmeta+visinaUIMargine)+visinaDugmeta {
					infoText = "Clear"
				}

			}

			// njanja todo napraviti fju koja ludi hex int za boju pretvara u rgba vrednost
			text, err = font.RenderUTF8Blended(infoText, sdl.Color{R: 255, G: 0, B: 0, A: 255})
			if err == nil {
				err = text.Blit(nil, surface, &sdl.Rect{X: 10, Y: 10, W: 0, H: 0})
				if err != nil {
					panic(err)
				}
			}
			defer text.Free()

		}

		window.UpdateSurface()

		// njanja: ovo renderuje krug oko četkice, ne može u brush fju jer se ona ne zove u svakom frejmu
		// takođe kursor flikeruje jer bude na kratko izbrisan kad se pozove updatesurface par linija iznad pa dok ne bude ponovo nacrtan
		// ako mene pitate mislim da daje odličan retro look
		cetkica := sdl.Rect{kursorPoslednjiX - velicinaKursora*brojPikselaPoCestici, kursorPoslednjiY - velicinaKursora*brojPikselaPoCestici, int32(2 * velicinaKursora * brojPikselaPoCestici), int32(2 * velicinaKursora * brojPikselaPoCestici)}
		renderer.SetDrawColor(255, 255, 255, 255)
		renderer.DrawRect(&cetkica)
		renderer.Present()

		if fpsCap > 0 {
			expectedFrameTime := uint64(1000 / fpsCap)
			realFrameTime := sdl.GetTicks64() - startTime
			if expectedFrameTime > realFrameTime {
				// o moj bože molim vas jedan jedini int ko je mislio da je ovo dobra ideja
				sdl.Delay(uint32(expectedFrameTime - realFrameTime))
			}
		}
		//fmt.Printf("FPS: %d\n", int(1000.0/float64(sdl.GetTicks64()-startTime)))
	}

}

// takozvano dinamičko skaliranje ekrana ili nešto ne znam lupio sam
// ako ovo ikada u praksi izbaci nešto što ne staje u ekran javite mi da ga sredim ali mislim da je to besmislen posao
func fitToScreen(screenPercentage int) (int32, int32, int32) {
	resolution := screenresolution.GetPrimary()
	adjustedScale := int32((float64(screenPercentage) / float64(100)) * float64(resolution.Height) / float64(visinaKanvasa))
	return adjustedScale, sirinaKanvasa * adjustedScale, visinaKanvasa * adjustedScale
}

// dodaje zidove oko matrice /limun
func zazidajMatricu(matrix [][]mat.Cestica) [][]mat.Cestica {
	for i := 0; i < sirinaKanvasa; i++ {
		matrix[i][0], matrix[i][visinaKanvasa-1] = mat.NewCestica(mat.Zid), mat.NewCestica(mat.Zid)
		matrix[i][1], matrix[i][visinaKanvasa-2] = mat.NewCestica(mat.Zid), mat.NewCestica(mat.Zid)
	}
	for j := 0; j < visinaKanvasa; j++ {
		matrix[0][j], matrix[sirinaKanvasa-1][j] = mat.NewCestica(mat.Zid), mat.NewCestica(mat.Zid)
		matrix[1][j], matrix[sirinaKanvasa-2][j] = mat.NewCestica(mat.Zid), mat.NewCestica(mat.Zid)
	}
	return matrix
}

// menja i vraća (x, y) koordinate tako da se nalaze na ekranu /limun
func clampCoords(x int32, y int32) (int32, int32) {
	return int32(math.Min(math.Max(float64(x), 0), sirinaKanvasa-1)),
		int32(math.Min(math.Max(float64(y), 0), visinaKanvasa-1))
}

func brush(matrix [][]mat.Cestica, bafer [][]mat.Cestica, x int32, y int32, state uint32) {
	//TODO za srednji klik da uzme materijal na koj mis trenutno pokazuje i postavi ga kao trenutni
	//ukoliko nije u pitanju Zid ili Prazno. Nije mi pri ruci mis, mrzi me da trazim koj je to stejt -s
	//a jeste sabani mogli ste ovo trideset puta uraditi danas -s
	if x > sirinaKanvasa*brojPikselaPoCestici {
		return
	}
	if state != 1 && state != 4 {
		return
	}

	if state == 1 {
		for i := -velicinaKursora; i <= velicinaKursora; i++ {
			for j := -velicinaKursora; j <= velicinaKursora; j++ {
				tx, ty := clampCoords(x/brojPikselaPoCestici+i, y/brojPikselaPoCestici+j)
				if matrix[tx][ty].Materijal == mat.Prazno || (trenutniMat == mat.Prazno && matrix[tx][ty].Materijal != mat.Zid) {
					matrix[tx][ty] = mat.NewCestica(trenutniMat)
					bafer[tx][ty] = mat.NewCestica(trenutniMat)
				}
			}
		}
	}

	if state == 4 {
		for i := -velicinaKursora; i <= velicinaKursora; i++ {
			for j := -velicinaKursora; j <= velicinaKursora; j++ {
				tx, ty := clampCoords(x/brojPikselaPoCestici+i, y/brojPikselaPoCestici+j)
				if matrix[tx][ty].Materijal != mat.Zid { //napomenuo bih da prazne cestice ovde brisemo i pravimo opet da bismo resetovali temp, inace bi bilo efikasnije samo postaviti im Materijal na Prazno, NAGADJAM
					matrix[tx][ty] = mat.NewCestica(mat.Prazno)
					bafer[tx][ty] = mat.NewCestica(mat.Prazno)
				}
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

		case *sdl.KeyboardEvent:
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
				normalMode = false
				tempMode = true
				densityMode = false
			}
			if keystates[sdl.SCANCODE_D] != 0 {
				normalMode = false
				tempMode = false
				densityMode = true
			}
			if keystates[sdl.SCANCODE_N] != 0 {
				normalMode = true
				tempMode = false
				densityMode = false
			}
			if keystates[sdl.SCANCODE_V] != 0 {
				textMode = !textMode
			}

		// njanja: za ovo mi je potreban diskretan klik a ne frejm sa dugmetom dole
		// p.s. hoćemo da ostavimo komentare ristoviću da ih vidi
		//boze pomogi -s
		// paa, barem imajte naznaku za svaku liniju čiji je čiji /limun
		case *sdl.MouseButtonEvent:
			if t.State == sdl.PRESSED {
				proveriPritisakNaGumb(matrix, bafer, t.X, t.Y)
			}

		// drag and drop slike je odmah učitava
		case *sdl.DropEvent:
			dropEvent := event.(*sdl.DropEvent)
			if dropEvent.Type == sdl.DROPFILE {
				filePath := string(dropEvent.File)
				// njanja: todo dodati support za bmp i webp

				err := ucitajSliku(filePath, matrix, bafer)
				if err != nil {
					sdl.ShowSimpleMessageBox(sdl.MESSAGEBOX_INFORMATION, "pesak", "rade samo png jpg bmp webp slike", nil)
				}

			}
		default:
			//null
		}
	}

	var x, y int32
	var state uint32
	x, y, state = sdl.GetMouseState()
	if x > 0 && x < sirinaProzora {
		kursorPoslednjiX = x
	}
	if y > 0 && y < visinaProzora {
		kursorPoslednjiY = y
	}

	if korisnikJeLimun {
		fmt.Printf("x: %d ", x)
		fmt.Printf("y: %d\t", y)
		fmt.Printf("xpx: %d ", x/brojPikselaPoCestici)
		fmt.Printf("ypx: %d\t", y/brojPikselaPoCestici)
		fmt.Printf("mb: %d\t", state)
		fmt.Printf("mat.Materijal: %d\t", trenutniMat)
		fmt.Printf("velicina: %d\t", velicinaKursora)
		fmt.Printf("pauza: %t\n", pause)
	}

	brush(matrix, bafer, x, y, state)

	return running

}

func proveriPritisakNaGumb(matrix, bafer [][]mat.Cestica, x, y int32) {
	//njanja: ovo je detekcija klika na gumb
	if x > sirinaProzora-marginaZaGumbad+sirinaUIMargine && x < sirinaProzora-sirinaUIMargine {
		// njanja: TODO namestiti da se ređaju u više kolona ako baš mora //mora
		// materijali
		if y < (visinaUIMargine+visinaDugmeta)*int32(len(boja)-1) && y%(visinaUIMargine+visinaDugmeta) > visinaUIMargine {
			trenutniMat = mat.Materijal(y / (visinaUIMargine + visinaDugmeta))
		}

		// njanja: hardkodovan broj specijalnih dugmića
		// PAUZA
		if y > visinaProzora-3*(visinaDugmeta+visinaUIMargine) && y < visinaProzora-3*(visinaDugmeta+visinaUIMargine)+visinaDugmeta {
			pause = !pause
		}
		// SEJV
		if y > visinaProzora-2*(visinaDugmeta+visinaUIMargine) && y < visinaProzora-2*(visinaDugmeta+visinaUIMargine)+visinaDugmeta {
			saveImage(matrix, int(brojPikselaPoCestici))
			sdl.ShowSimpleMessageBox(sdl.MESSAGEBOX_INFORMATION, "pesak", "sačuvan B)", nil)
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
			zazidajMatricu(bafer) //ova linija suvisna, takodje razmisli da koristis f ju napraviSlajs (na dnu fajla), takodje mozda da zazidavanje zovemo u njoj? pa izbacimo iz mejna i odavde, takodje ovo je naaaaaaajduzi komentar poput one Ivonine naaaaajduze lasice na celom svetu samo cekam da mu eksplodiraju oci -s
		}
	}
}
func update(matrix [][]mat.Cestica, bafer [][]mat.Cestica) {

	for j := 1; j < visinaKanvasa-1; j++ {
		for i := 1; i < sirinaKanvasa-1; i++ {
			bafer[i][j].Temperatura = 0
		}
	} // da, mora redno. izdvojte u fje ako vam se ne svidja kod -s
	for j := 1; j < visinaKanvasa-1; j++ {
		for i := 1; i < sirinaKanvasa-1; i++ {
			mat.UpdateTemp(matrix, bafer, i, j)
		}
	}
	minTempRendered = mat.MaxTemp
	maxTempRendered = mat.MinTemp
	for j := 1; j < visinaKanvasa-1; j++ {
		for i := 1; i < sirinaKanvasa-1; i++ {
			matrix[i][j].Temperatura = bafer[i][j].Temperatura
			temperatura := matrix[i][j].Temperatura
			if temperatura+1 > maxTempRendered {
				maxTempRendered = temperatura + 1
			}
			if temperatura < minTempRendered {
				minTempRendered = temperatura
			}
			//todo smisli sta sa tikerima
		}
	}

	for j := 1; j < visinaKanvasa-1; j++ {
		for i := 1; i < sirinaKanvasa-1; i++ {
			mat.UpdatePhaseOfMatter(matrix, bafer, i, j)
		}
	}

	ja := make([]int, visinaKanvasa)
	for j := range ja {
		ja[j] = j
	}
	ia := make([]int, sirinaKanvasa)
	for i := range ia {
		ia[i] = i
	}
	rand.Shuffle(len(ja), func(i, j int) { ja[i], ja[j] = ja[j], ja[i] })
	rand.Shuffle(len(ia), func(i, j int) { ia[i], ia[j] = ia[j], ia[i] })

	for j := 1; j < visinaKanvasa-1; j++ {
		for i := 1; i < sirinaKanvasa-1; i++ {
			mat.UpdatePosition(matrix, bafer, ia[i], ja[j])
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
			// njanja: braaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaate je l stvarno toliko teško kompajleru da implicitno konvertuje int u int32
			rect := sdl.Rect{int32(i) * brojPikselaPoCestici, int32(j) * brojPikselaPoCestici, brojPikselaPoCestici, brojPikselaPoCestici}
			if tempMode {
				bojaTemp := izracunajTempBoju(matrix[i][j])
				surface.FillRect(&rect, bojaTemp)
			} else if densityMode {
				gustTemp := izracunajGustBoju(float64(gus[matrix[i][j].Materijal]))
				surface.FillRect(&rect, gustTemp)
			} else {
				surface.FillRect(&rect, boja[matrix[i][j].Materijal])
			}
		}
	}
}

// todo probao bih alternativu da napravim -s
// onda stavi pravi #TODO, kolega /limun
// xDDD
/*
func izracunajTempBoju(temp int32) uint32 {
	temp *= tempColorMultiplier
	temp /= 100
	if temp > 0 {
		temp = int32Min(temp, 255)
		temp = (255-temp) << 8 + (255 << 16)
	} else if temp < 0 {
		temp *= -1
		temp = int32Min(temp, 255)
		temp = (255-temp) << 8  + 255
	} else {
		temp = 230
		temp += (230 << 8) + (230 << 16)
	}

	hexadeca := strconv.FormatUint(uint64(temp), 16)
	tempBoja, err := strconv.ParseUint(hexadeca, 16, 32)
	if err != nil {
		panic(err)
	}

	return uint32(tempBoja)
}
*/
var minTempRendered int32 = 20
var maxTempRendered int32 = 21

func izracunajTempBoju(zrno mat.Cestica) uint32 {

	//	minTemp := mat.MinTemp
	//	maxTemp := mat.MaxTemp

	temperatura := zrno.Temperatura

	// tMin         temp                  tMax
	// 0            xx                    255

	//	(temp - tMin) / (tMax - tMin) = xx / 255
	// xx = 255(temp-tMin)/(tMax-tMin)

	var crvenaKomponenta uint32 = uint32(255 * (temperatura - minTempRendered) / (maxTempRendered - minTempRendered))
	var plavaKomponenta uint32 = uint32(255 - crvenaKomponenta)
	var zelenaKomponenta uint32 = 0

	var boja uint32 = (crvenaKomponenta*256+zelenaKomponenta)*256 + plavaKomponenta
	return boja
	/**/
}
func izracunajGustBoju(gust float64) uint32 {
	if gust > 1220 {
		gustInt := int32(math.Max(math.Min(gust*255/400000, 255), 0))
		gustInt = (gustInt << 8)
		gust = float64(gustInt)
	} else if gust < 1220 {
		gust = math.Min(gust*255/1000, 255)
		gust += float64(int32(gust) << 16)
	} else {
		gust = (200 << 16) + (200 << 8) + 200
	}

	hexadeca := strconv.FormatUint(uint64(gust), 16)
	gustBoja, err := strconv.ParseUint(hexadeca, 16, 32)
	if err != nil {
		panic(err)
	}

	return uint32(gustBoja)
}

func napraviSlajs() [][]mat.Cestica {
	slajs := make([][]mat.Cestica, sirinaKanvasa)
	for i := 0; i < sirinaKanvasa; i++ {
		kolona := make([]mat.Cestica, visinaKanvasa)
		for j := 0; j < visinaKanvasa; j++ {
			kolona[j] = mat.NewCestica(mat.Prazno)
		}
		slajs[i] = kolona
	}
	return slajs
}
