//Package mathPack sadrzi Abs, Min i Maks funkcije za tipove int i int32
package mathPack

//AbsInt racuna apsolutnu vrednost promenjive tipa int
func AbsInt(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

//MinInt vraca manji od dva broja tipa int
func MinInt(x int, y int) int {
	if x < y {
		return x
	}
	return y
}

//MaxInt vraca veci od dva broja tipa int
func MaxInt(x int, y int) int {
	if x > y {
		return x
	}
	return y
}

//AbsInt32 racuna apsolutnu vrednost promenjive tipa int32
func AbsInt32(x int32) int32 {
	if x < 0 {
		return -x
	}
	return x
}

//MinInt32 vraca manji od dva broja tipa int32
func MinInt32(x int32, y int32) int32 {
	if x < y {
		return x
	}
	return y
}

//MaxInt32 vraca veci od dva broja tipa int32
func MaxInt32(x int32, y int32) int32 {
	if x > y {
		return x
	}
	return y
}