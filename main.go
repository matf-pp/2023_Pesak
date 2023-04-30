package main

import (
	"main/brushPack"
	"main/fontPack"
	"main/mat"
	"main/matrixPack"
	"main/musicPack"
	"main/screenPack"

	"fmt"
	"math/rand"

	// ako mix zabaguje, $export CGO_CFLAGS=-I/usr/include/SDL2 	/limun
	"github.com/veandco/go-sdl2/mix"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

// njanja: pazite ovo
const korisnikNijeNanja = false
const korisnikJeLimun = false

// njanja: ovo je loša praksa majmuni
// e a reci je l si provalio bukvalno je kao `using` u cpp -s
var boja = mat.Boja

// FPS cap, kontam da je zgodno za testiranje staviti neki nizak, 0 = unlimited
// 250 FPS-a ?!?!?!?! /limun
var fpsCap = 120

var pozadinaGuia uint32 = 0x111122

var keystates = sdl.GetKeyboardState()

func main() {
	// koji procenat ekrana želimo da nam igrica zauzme (probajte da ukucate 0 ili -50 ili tako nešto wild) (spojler: radiće)
	if screenPack.AutoFitScreen {
		matrixPack.BrPiksPoCestici, screenPack.SirinaProzora, screenPack.VisinaProzora = screenPack.FitToScreen(40)
	}

	screenPack.UpdateRazmere()

	var font = fontPack.FontInit()
	defer ttf.Quit()

	screenPack.InitEverything()
	defer sdl.Quit()

	musicPack.MusicInit()
	defer mix.Quit()
	musicPack.OpenAudio()
	mus := musicPack.LoadMusic("audio/bitstorm.mp3")

	window := screenPack.CreateWindow()
	defer window.Destroy()

	renderer := screenPack.CreateRenderer(window)
	defer renderer.Destroy()

	texture := screenPack.CreateTexture(renderer)
	defer texture.Destroy()

	pixels := make([]byte, matrixPack.SirinaKan*matrixPack.VisinaKan*3)

	font = fontPack.SetFont()
	defer font.Close()

	go connectToDiscord()

	var matrica [][]mat.Cestica = matrixPack.NapraviSlajs()

	matrica = matrixPack.ZazidajMatricu(matrica)

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
			_ = mus.Play(-1)
			matrixPack.ResetSound = false
		}

		if musicPack.Mutirana {
			_ = mix.VolumeMusic(0)
		} else {
			_ = mix.VolumeMusic(musicPack.Zvuk)
		}

		// njanja: ovo je vraćalo surface? nz je l to potrebno ???
		// ne znam ni ja /limun
		screenPack.RenderujGumbZaSveMaterijale(renderer)
		var hexColor uint32

		brushGumb := screenPack.CreateSpecialGumb(4)
		hexColor = 0x00ffff
		renderer.SetDrawColor(uint8(hexColor>>16), uint8(hexColor>>8), uint8(hexColor), 255)
		renderer.FillRect(&brushGumb)

		plejGumb := screenPack.CreateSpecialGumb(3)
		// todo: izbeći ove konverzije i samo koristiti rgba direktno
		if matrixPack.Pause {
			hexColor = 0x00ff00
		} else {
			hexColor = 0xffa500
		}
		renderer.SetDrawColor(uint8(hexColor>>16), uint8(hexColor>>8), uint8(hexColor), 255)
		renderer.FillRect(&plejGumb)

		sejvGumb := screenPack.CreateSpecialGumb(2)
		hexColor = 0x0000ff
		renderer.SetDrawColor(uint8(hexColor>>16), uint8(hexColor>>8), uint8(hexColor), 255)
		renderer.FillRect(&sejvGumb)

		resetGumb := screenPack.CreateSpecialGumb(1)
		hexColor = 0xff0000
		renderer.SetDrawColor(uint8(hexColor>>16), uint8(hexColor>>8), uint8(hexColor), 255)
		renderer.FillRect(&resetGumb)

		if matrixPack.TxtMode {
			fontPack.TextMaker(font, renderer, matrica)
		}

		//window.UpdateSurface()

		brushPack.OblikCetkice(matrixPack.KruzniBrush, renderer)

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
		//fmt.Printf("FPS: %d\n", int(1000.0/float64(sdl.GetTicks64()-startTime)))
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

		// ne volim ovo /limun
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
			if keystates[sdl.SCANCODE_C] != 0 {
				for j := 0; j < matrixPack.VisinaKan; j++ {
					for i := 0; i < matrixPack.SirinaKan; i++ {
						matrix[i][j] = mat.NewCestica(mat.Prazno)
					}
				}
				matrixPack.ZazidajMatricu(matrix)

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
				musicPack.Mutirana = !musicPack.Mutirana
			}
			if keystates[sdl.SCANCODE_Z] != 0 {
				if musicPack.Zvuk > 0 {
					musicPack.Zvuk -= 5
				}
			}
			if keystates[sdl.SCANCODE_X] != 0 {
				// do daske! /limun
				// njanja: e isk msm da je difolt preglasan
				// ja držim nisko zvuk pa je možda do toga /limun
				musicPack.Zvuk += 5
			}
			if keystates[sdl.SCANCODE_B] != 0 {
				matrixPack.KruzniBrush = !matrixPack.KruzniBrush
			}
			if keystates[sdl.SCANCODE_LSHIFT] != 0 {
				brushPack.ShiftOn = true
			} else {
				brushPack.ShiftOn = false
			}
			if keystates[sdl.SCANCODE_G] != 0 {
				mat.Obrnuto *= -1
			}

		// njanja: za ovo mi je potreban diskretan klik a ne frejm sa dugmetom dole
		// p.s. hoćemo da ostavimo komentare ristoviću da ih vidi
		// boze pomogi -s
		// paa, barem imajte naznaku za svaku liniju čiji je čiji /limun
		case *sdl.MouseButtonEvent:
			if t.State == sdl.PRESSED {
				screenPack.ProveriPritisakNaGumb(matrix, t.X, t.Y)
			}

		case *sdl.MouseWheelEvent:
			if brushPack.ShiftOn {
				if t.Y < 0 {
					screenPack.SledeciMaterijal()
				} else {
					screenPack.PrethodniMaterijal()
				}
			} else {
				if t.Y > 0 {
					if screenPack.VelicinaKursora < screenPack.MaxKursor {
						screenPack.VelicinaKursora = screenPack.VelicinaKursora + 1
					}
				} else {
					if screenPack.VelicinaKursora > 1 {
						screenPack.VelicinaKursora = screenPack.VelicinaKursora - 1
					}
				}
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
