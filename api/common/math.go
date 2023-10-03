package common

func Pow[T int | int16 | int32 | uint8 | uint16 | uint32](src T, step uint8) T {
	if src == 0 {
		return 0
	}
	if step == 0 {
		return 1
	}
	var ans T = 1
	var i uint8
	for i = 0; i < step; i++ {
		ans *= src
	}
	return ans
}
