package main

import (
	"main/brushPack"
	"main/fontPack"
	"main/mat"
	"main/matrixPack"
	"main/screenPack"

	"fmt"
	"math/rand"

	"github.com/veandco/go-sdl2/mix"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

// Njanjavi uradi ovo pls /limun
// njanja: todo dodati keyed fields u rectovima da nas ne smaraju žute linije (ili bar mene jer koliko sam shvatio vi ovo bukvalno pišete u noutpedu)
// njanja: todo izbaciti sve ove gluposti koje redovno menjamo u konfig fajl i staviti da bude gitignorovan

// njanja: pazite ovo
const korisnikNijeNanja = false
const korisnikJeLimun = false

// njanja: ovo je loša praksa majmuni
// e a reci je l si provalio bukvalno je kao `using` u cpp -s
var boja = mat.Boja
var gus = mat.Lambda

// FPS cap, kontam da je zgodno za testiranje staviti neki nizak, 0 = unlimited
var fpsCap = 0

var pozadinaGuia uint32 = 0x111122

var keystates = sdl.GetKeyboardState()

var mutirana = false
var zvuk = 100

func main() {
	// koji procenat ekrana želimo da nam igrica zauzme (probajte da ukucate 0 ili -50 ili tako nešto wild) (spojler: radiće)
	if screenPack.AutoFitScreen {
		matrixPack.BrPiksPoCestici, screenPack.SirinaProzora, screenPack.VisinaProzora = screenPack.FitToScreen(80)
	}

	screenPack.MarginaZaGumbad = screenPack.BrojKolona*(screenPack.SirinaDugmeta+screenPack.SirinaUIMargine) + screenPack.SirinaUIMargine
	screenPack.SirinaProzora += screenPack.MarginaZaGumbad

	// njanja: da vidimo hoće li ovo raditi lepo
	var font *ttf.Font
	//var text *sdl.Surface
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

	err = mix.Init(mix.INIT_MP3)
	if err != nil {
		panic(err)
	}
	defer mix.Quit()

	err = mix.OpenAudio(44100, mix.INIT_MP3, 2, 1024)
	if err != nil {
		panic(err)
	}

	var mus *mix.Music
	mus, err = mix.LoadMUS("audio/bitstorm.mp3")
	err = mus.Play(-1)
	if err != nil {
		panic(err)
	}

	window := screenPack.CreateWindow()
	defer window.Destroy()

	// njanja: ovo vratiti u neki screenpack nešto ako baš treba mada ne kontam poentu funkcija koje se izvršavaju bukvalno jednom i koje služe samo da kod ne bi bio u mejnu
	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}
	defer renderer.Destroy()

	texture, err := renderer.CreateTexture(uint32(sdl.PIXELFORMAT_RGB24), sdl.TEXTUREACCESS_STATIC, matrixPack.SirinaKan, matrixPack.VisinaKan)
	if err != nil {
		panic(err)
	}
	defer texture.Destroy()

	pixels := make([]byte, matrixPack.SirinaKan*matrixPack.VisinaKan*3)

	font = fontPack.SetFont()
	defer font.Close()

	// njanja: diskord integracija oterana u poseban fajl jer je mrvicu čonki
	// ume da pravi probleme i male blek skrinove na početku dok se ne konektuje ili ne tajmautuje pa ga pozivam kao korutinu
	go connectToDiscord()

	var matrica [][]mat.Cestica = matrixPack.NapraviSlajs()

	matrica = matrixPack.ZazidajMatricu(matrica)

	// enejblujemo dropove, stavite ovo gde hoćete
	sdl.EventState(sdl.DROPFILE, sdl.ENABLE)

	running := true
	for running {
		// fps counter
		var startTime = sdl.GetTicks64()

		running = pollEvents(matrica)
		if !matrixPack.Pause {
			update(matrica)
		}
		matrixPack.Render(matrica, renderer, texture, pixels, screenPack.SirinaProzora-screenPack.MarginaZaGumbad, screenPack.VisinaProzora)

		if matrixPack.ResetSound {
			err = mus.Play(-1)
			matrixPack.ResetSound = false
		}

		if mutirana {
			_ = mix.VolumeMusic(0)
		} else {
			_ = mix.VolumeMusic(zvuk)
		}

		// njanja: ovo je vraćalo surface? nz je l to potrebno ???
		screenPack.RenderujGumbZaSveMaterijale(renderer)

		plejGumb := screenPack.CreatePlayGumb()
		// todo: izbeći ove konverzije i samo koristiti rgba direktno
		var hexColor uint32
		if matrixPack.Pause {
			hexColor = 0x00ff00
		} else {
			hexColor = 0xffa500
		}
		renderer.SetDrawColor(uint8(hexColor>>16), uint8(hexColor>>8), uint8(hexColor), 255)
		renderer.FillRect(&plejGumb)

		sejvGumb := screenPack.CreateSaveGumb()
		hexColor = 0x0000ff
		renderer.SetDrawColor(uint8(hexColor>>16), uint8(hexColor>>8), uint8(hexColor), 255)
		renderer.FillRect(&sejvGumb)

		resetGumb := screenPack.CreateResetGumb()
		hexColor = 0xff0000
		renderer.SetDrawColor(uint8(hexColor>>16), uint8(hexColor>>8), uint8(hexColor), 255)
		renderer.FillRect(&resetGumb)

		if matrixPack.TxtMode {
			fontPack.TextMaker(font, renderer, matrica)
		}

		window.UpdateSurface()

		// njanja: ovo renderuje krug oko četkice, ne može u brush fju jer se ona ne zove u svakom frejmu
		// takođe kursor flikeruje jer bude na kratko izbrisan kad se pozove updatesurface par linija iznad pa dok ne bude ponovo nacrtan
		// ako mene pitate mislim da daje odličan retro look
		cetkica := sdl.Rect{X: screenPack.KursorPoslednjiX - screenPack.VelicinaKursora*matrixPack.BrPiksPoCestici, Y: screenPack.KursorPoslednjiY - screenPack.VelicinaKursora*matrixPack.BrPiksPoCestici, W: int32(2 * screenPack.VelicinaKursora * matrixPack.BrPiksPoCestici), H: int32(2 * screenPack.VelicinaKursora * matrixPack.BrPiksPoCestici)}
		renderer.SetDrawColor(255, 255, 255, 255)
		renderer.DrawRect(&cetkica)
		renderer.SetDrawColor(uint8(pozadinaGuia>>16), uint8(pozadinaGuia>>8), uint8(pozadinaGuia), 255)
		renderer.Present()

		if fpsCap > 0 {
			expectedFrameTime := uint64(1000 / fpsCap)
			realFrameTime := sdl.GetTicks64() - startTime
			if expectedFrameTime > realFrameTime {
				// o moj bože molim vas jedan jedini int ko je mislio da je ovo dobra ideja
				sdl.Delay(uint32(expectedFrameTime - realFrameTime))
			}
		}
		fmt.Printf("FPS: %d\n", int(1000.0/float64(sdl.GetTicks64()-startTime)))
	}

}

func pollEvents(matrix [][]mat.Cestica) bool {
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
				screenPack.TrenutniMat = mat.Materijal(10)
			}
			if keystates[sdl.SCANCODE_1] != 0 {
				screenPack.TrenutniMat = mat.Materijal(1)
			}
			if keystates[sdl.SCANCODE_2] != 0 {
				screenPack.TrenutniMat = mat.Materijal(2)
			}
			if keystates[sdl.SCANCODE_3] != 0 {
				screenPack.TrenutniMat = mat.Materijal(3)
			}
			if keystates[sdl.SCANCODE_4] != 0 {
				screenPack.TrenutniMat = mat.Materijal(4)
			}
			if keystates[sdl.SCANCODE_5] != 0 {
				screenPack.TrenutniMat = mat.Materijal(5)
			}
			if keystates[sdl.SCANCODE_6] != 0 {
				screenPack.TrenutniMat = mat.Materijal(6)
			}
			if keystates[sdl.SCANCODE_7] != 0 {
				screenPack.TrenutniMat = mat.Materijal(7)
			}
			if keystates[sdl.SCANCODE_8] != 0 {
				screenPack.TrenutniMat = mat.Materijal(8)
			}
			if keystates[sdl.SCANCODE_9] != 0 {
				screenPack.TrenutniMat = mat.Materijal(9)
			}
			if keystates[sdl.SCANCODE_LEFTBRACKET] != 0 {
				screenPack.TrenutniMat = mat.Toplo
			}
			if keystates[sdl.SCANCODE_RIGHTBRACKET] != 0 {
				screenPack.TrenutniMat = mat.Hladno
			}
			if keystates[sdl.SCANCODE_DOWN] != 0 {
				if screenPack.VelicinaKursora > 1 {
					screenPack.VelicinaKursora = screenPack.VelicinaKursora - 1
				}
			}
			if keystates[sdl.SCANCODE_UP] != 0 {
				if screenPack.VelicinaKursora < screenPack.MaxKursor {
					screenPack.VelicinaKursora = screenPack.VelicinaKursora + 1
				}
			}
			if keystates[sdl.SCANCODE_P] != 0 {
				matrixPack.Pause = !matrixPack.Pause
			}
			if keystates[sdl.SCANCODE_T] != 0 {
				matrixPack.NMode = false
				matrixPack.TMode = true
				matrixPack.DMode = false
			}
			if keystates[sdl.SCANCODE_D] != 0 {
				matrixPack.NMode = false
				matrixPack.TMode = false
				matrixPack.DMode = true
			}
			if keystates[sdl.SCANCODE_N] != 0 {
				matrixPack.NMode = true
				matrixPack.TMode = false
				matrixPack.DMode = false
			}
			if keystates[sdl.SCANCODE_V] != 0 {
				matrixPack.TxtMode = !matrixPack.TxtMode
			}
			if keystates[sdl.SCANCODE_R] != 0 {
				matrixPack.ResetSound = !matrixPack.ResetSound
			}
			if keystates[sdl.SCANCODE_M] != 0 {
				mutirana = !mutirana
			}
			if keystates[sdl.SCANCODE_Z] != 0 {
				if zvuk > 0 {
					zvuk -= 5
				}
			}
			if keystates[sdl.SCANCODE_X] != 0 {
				// do daske! /limun
				zvuk += 5
			}

		// njanja: za ovo mi je potreban diskretan klik a ne frejm sa dugmetom dole
		// p.s. hoćemo da ostavimo komentare ristoviću da ih vidi
		//boze pomogi -s
		// paa, barem imajte naznaku za svaku liniju čiji je čiji /limun
		case *sdl.MouseButtonEvent:
			if t.State == sdl.PRESSED {
				screenPack.ProveriPritisakNaGumb(matrix, t.X, t.Y)
			}

		// drag and drop slike je odmah učitava
		case *sdl.DropEvent:
			dropEvent := event.(*sdl.DropEvent)
			if dropEvent.Type == sdl.DROPFILE {
				filePath := string(dropEvent.File)
				// njanja: todo dodati support za bmp i webp

				err := screenPack.UcitajSliku(filePath, matrix)
				if err != nil {
					sdl.ShowSimpleMessageBox(sdl.MESSAGEBOX_INFORMATION, "pesak", "rade samo png jpg bmp jbg webp itd slike", nil)
				}

			}
		default:
			//null
		}
	}

	var x, y int32
	var state uint32
	x, y, state = sdl.GetMouseState()
	if x > 0 && x < screenPack.SirinaProzora {
		screenPack.KursorPoslednjiX = x
	}
	if y > 0 && y < screenPack.VisinaProzora {
		screenPack.KursorPoslednjiY = y
	}

	if korisnikNijeNanja {
		//		fmt.Printf("x: %d ", x)
		//		fmt.Printf("y: %d\t", y)
		//		fmt.Printf("xpx: %d ", x/matrixPack.BrPiksPoCestici)
		//		fmt.Printf("ypx: %d\t", y/matrixPack.BrPiksPoCestici)
		fmt.Printf("mb: %d\t", state)
		fmt.Printf("mat.Materijal: %d\t", screenPack.TrenutniMat)
		fmt.Printf("velicina: %d\t", screenPack.VelicinaKursora)
		fmt.Printf("pauza: %t\t", matrixPack.Pause)
		if !korisnikJeLimun {
			fmt.Printf("brCestica: %d\n", brCestica)
			//			fmt.Printf("brLave: %d\n", brLave)
			//			fmt.Printf("brKamena: %d\n", brKamena)
		} else {
			fmt.Printf("\n")
		}
	}

	brushPack.Brush(matrix, x, y, state)

	return running

}

var brCestica int = 0
var brLave int = 0
var brKamena int = 0

func update(matrix [][]mat.Cestica) {
	brCestica, brKamena, brLave = matrixPack.IzbrojiCesticeKamenLavu(matrix)

	for j := 1; j < matrixPack.VisinaKan-1; j++ {
		for i := 1; i < matrixPack.SirinaKan-1; i++ {
			mat.UpdateTemp(matrix, i, j)
		}
	}

	matrixPack.MinTempRendered = mat.MaxTemp
	matrixPack.MaxTempRendered = mat.MinTemp
	for j := 1; j < matrixPack.VisinaKan-1; j++ {
		for i := 1; i < matrixPack.SirinaKan-1; i++ {
			matrix[i][j].Temperatura = matrix[i][j].BaferTemp
			matrix[i][j].BaferTemp = 0
			temperatura := matrix[i][j].Temperatura
			if temperatura+1 > matrixPack.MaxTempRendered {
				matrixPack.MaxTempRendered = temperatura + 1
			}
			if temperatura < matrixPack.MinTempRendered {
				matrixPack.MinTempRendered = temperatura
			}
			//todo smisli sta sa tikerima
		}
	}
	for j := 1; j < matrixPack.VisinaKan-1; j++ {
		for i := 1; i < matrixPack.SirinaKan-1; i++ {
			mat.UpdatePhaseOfMatter(matrix, i, j)
		}
	}

	ja := make([]int, matrixPack.VisinaKan)
	for j := range ja {
		ja[j] = j
	}
	ia := make([]int, matrixPack.SirinaKan)
	for i := range ia {
		ia[i] = i
	}
	rand.Shuffle(len(ja), func(i, j int) { ja[i], ja[j] = ja[j], ja[i] })
	rand.Shuffle(len(ia), func(i, j int) { ia[i], ia[j] = ia[j], ia[i] })

	for j := 1; j < matrixPack.VisinaKan-1; j++ {
		for i := 1; i < matrixPack.SirinaKan-1; i++ {
			mat.UpdatePosition(matrix, ia[i], ja[j])
		}
	}
}
