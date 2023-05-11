//gravityPack je paket koji pomaže u radu sa gravitacijom i smerom gravitacije
package gravityPack

import (

)

//Obrnuto odredjuje smer gravitacije
var Obrnuto = 1
//GRuka određuje da li je mis centar gravitacije
var GRuka bool = false

/*
	4  5  6
	3  -  7
	2  1  8
*/
//ProveriOktant vraća oktant kome pripada Cestica u odnosu na položaj miša
func ProveriOktant(x int, y int, xMis int, yMis int) int {
	if (y - yMis) >= 3 * (x - xMis) && 3 * (y - yMis) < (x - xMis) {
		return 8
	}
	if 3 * (y - yMis) >= (x - xMis) && 3 * (y - yMis) < -(x - xMis) {
		return 7
	}
	if 3 * (y - yMis) >= -(x - xMis) && y - yMis < -3 * (x - xMis) {
		return 6
	}
	if y - yMis >= -3 * (x - xMis) && y - yMis >= 3 * (x - xMis) {
		return 5
	} 
	if y - yMis < 3 * (x - xMis) && 3 * (y - yMis) >= (x - xMis) {
		return 4
	}
	if 3 * (y - yMis) < (x - xMis) && 3 * (y - yMis) >= -(x - xMis) {
		return 3
	}
	if 3 * (y - yMis) < -(x - xMis) && (y - yMis) >= -3 * (x - xMis) {
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