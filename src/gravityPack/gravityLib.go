//gravityPack je paket koji pomaže u radu sa gravitacijom i smerom gravitacije
package gravityPack

import (
	"math"
)

//Obrnuto odredjuje smer gravitacije
var Obrnuto = 1
//GRuka određuje da li je mis centar gravitacije
var GRuka bool = false
//GTacka određuje da li je fiksna tačka centar gravitacije
var GTacka bool = false
//CrnaRupa određuje da li je crna rupa uključena
var CrnaRupa bool = false
//(CentarGravitacijeX, CentarGravitacijeY) je tačka u kojoj je miš bio kada je pritisnut Q
var CentarGravitacijeX int
var CentarGravitacijeY int

/*
	4  5  6
	3  -  7
	2  1  8
*/
//ProveriOktant vraća oktant kome pripada Cestica u odnosu na položaj miša
func ProveriOktant(x int, y int, xMis int, yMis int) int {
	dy := float64(y-yMis)
	dx := float64(x-xMis)
	
	alfa := math.Atan2(dy, dx)
	alfaDeg := alfa * 180 / math.Pi

	oktant := int(math.Ceil(alfaDeg/45)) + 3 % 8
	if CrnaRupa {
		oktant = (oktant + 1)
	}
	if oktant == 0 {
		oktant = 8
	}

	return oktant
}
//GdePada vraća smerI i smerJ prema kome Cestica pada
func GdePadaDole(oktant int) (int, int) {
	switch oktant {
		case 1: return 0, 1
		case 2: return -1, 1
		case 3: return -1, 0
		case 4: return -1, -1
		case 5: return 0, -1
		case 6: return 1, -1
		case 7: return 1, 0
		case 8:	return 1, 1
		default: return 0, 1
	}
}
//GdePadaUkosoDole vraća smerILevo, smerJLevo, smerIDesno, smerJDesno
func GdePadaUkosoDole(oktant int) (int, int, int, int) {
	switch oktant {
		case 1: return -1, 1,		1, 1
		case 2: return -1, 0,		0, 1
		case 3: return -1, -1,		-1, 1
		case 4: return 0, -1,		-1, 0
		case 5: return 1, -1,		-1, -1
		case 6: return 1, 0,		0, -1
		case 7: return 1, 1,		1, -1
		case 8:	return 0, 1,		1, 0
		default: return -1, 1,		1, 1
	}
}
//GdeIdeLevoDesno vraća smerILevo2, smerJLevo2, smerIDesno2, smerJDesno2
func GdeIdeLevoDesno(oktant int) (int, int, int, int) {
	switch oktant {
		case 1: return -1, 0,		1, 0
		case 2: return -1, -1,		1, 1
		case 3: return 0, -1,		0, 1
		case 4: return 1, -1,		-1, 1
		case 5: return 1, 0,		-1, 0
		case 6: return 1, 1,		-1, -1
		case 7: return 0, 1,		0, -1
		case 8:	return -1, 1,		1, -1
		default: return -1, 0,		1, 0
	}
}
//UpadaUCrnuRupu vraća da li koordinata upada u crnu rupu veličine četkice
func UpadaUCrnuRupu1(i int, j int, velKur int) bool {
	x := i - CentarGravitacijeX
	y := j - CentarGravitacijeY
	if x*x + y*y <= velKur*velKur {
		return true
	}

	return false
}

func UpadaUCrnuRupu2(i int, j int, velKur int, a int, b int) bool {
	x := i - a
	y := j - b
	if x*x + y*y <= velKur*velKur {
		return true
	}

	return false
}