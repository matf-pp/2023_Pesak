//gravityPack je paket koji pomaže u radu sa gravitacijom i smerom gravitacije
package gravityPack

import (

)

//Obrnuto odredjuje smer gravitacije
var Obrnuto = 1
//GRuka određuje da li je mis centar gravitacije
var GRuka bool = false
//GTacka određuje da li je fiksna tačka centar gravitacije
var GTacka bool = false
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
	dy := (y-yMis)
	dx := (x-xMis)
	koeff := 2//2.4142135623730950488016
	if dy >= koeff * (dx) && koeff * (dy) < (dx) {
		return 8
	}
	if koeff * (dy) >= (dx) && koeff * (dy) < -(dx) {
		return 7
	}
	if koeff * (dy) >= -(dx) && dy < -koeff * (dx) {
		return 6
	}
	if dy >= -koeff * (dx) && dy >= koeff * (dx) {
		return 5
	} 
	if dy < koeff * (dx) && koeff * (dy) >= (dx) {
		return 4
	}
	if koeff * (dy) < (dx) && koeff * (dy) >= -(dx) {
		return 3
	}
	if koeff * (dy) < -(dx) && (dy) >= -koeff * (dx) {
		return 2
	}
	return 1
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