package mathPack

func AbsInt(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func MinInt(x int, y int) int {
	if x < y {
		return x
	}
	return y
}

func MaxInt(x int, y int) int {
	if x > y {
		return x
	}
	return y
}

func AbsInt32(x int32) int32 {
	if x < 0 {
		return -x
	}
	return x
}

func MinInt32(x int32, y int32) int32 {
	if x < y {
		return x
	}
	return y
}

func MaxInt32(x int32, y int32) int32 {
	if x > y {
		return x
	}
	return y
}