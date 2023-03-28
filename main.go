package main

import 
(
	"github.com/veandco/go-sdl2/sdl"
	"math/rand"
	"math"
	//"time"
	"fmt"
)

func main() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow("test", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		800, 600, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	surface, err := window.GetSurface()
	if err != nil {
		panic(err)
	}
	surface.FillRect(nil, 0)

	var matrix[800][600]int

	var event sdl.Event
	running := true
	for running {
		// r := rand.Intn(800-3)
		// for i := 0; i < 3; i++ {
		// 	for j := 0; j < 3; j++ {
		// 		matrix[r+i][0+j] = 1
		// 	}
		// }

		for i := 1; i < 800-3; i++ {
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
		}
		for i := 0; i < 800-3+1; i++ {
			for j := 600-2; j >= 0; j-- {
				rect := sdl.Rect{int32(i), int32(j), 1, 1}
				if(matrix[i][j] == 0) {
					surface.FillRect(&rect, 0)
				}
				if(matrix[i][j] == 1) {
					surface.FillRect(&rect, 0xffff66)
				}
			}
		}

		var x int32
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
		}

		window.UpdateSurface()
		
		//time.Sleep(1 * time.Millisecond)
	}
}