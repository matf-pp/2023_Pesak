//Package musicPack je paket za pustanje muzike
package musicPack

import (
	"github.com/veandco/go-sdl2/mix"
)

//Mutirana odredjuje da li je muzika utisana na 0 ili nije
var Mutirana = false
//Zvuk odredjuje jacinu muzike
var Zvuk = 33

//MusicInit inicira muziku
func MusicInit() {
	err := mix.Init(mix.INIT_MP3)
	if err != nil {
		panic(err)
	}
}
//OpenAudio otvara audio na koji se stavljaju pesme
func OpenAudio() {
	err := mix.OpenAudio(44100, mix.INIT_MP3, 2, 1024)
	if err != nil {
		panic(err)
	}
}
//LoadMusic uzima putanju do mp3 fajla i postavlja sadrzaj fajla u mus promenljivu
func LoadMusic(file string) *mix.Music {
	var mus *mix.Music
	mus, _ = mix.LoadMUS(file)
	err := mus.Play(-1)
	if err != nil {
		panic(err)
	}
	return mus
}