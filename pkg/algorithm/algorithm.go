package algorithm

func Binary_Search[T any](src []T, borad int, target T, cmp func(int, T) int) (int, bool) {
	if borad == 0 {
		return 0, true
	}
	var left, right = 0, borad
	var mid int
	// _, f, line, _ := runtime.Caller(1)
	// fmt.Printf("origin left %d right %d call from %s:%d\n", left, right, f, line)
	for left < right-1 {
		mid = (left + right) / 2
		// fmt.Printf("compare target %d %v %v\n", mid, src[mid], target)
		switch cmp(mid, target) {
		case 0:
			return mid, false
		case 1: //target <src[mid]
			right = mid
		case -1: //target > src[mid]
			left = mid
		}
	}
	if left == right {
		return left, true
	}
	switch cmp(left, target) {
	case 0:
		return left, false
	case 1:
		return left, true
	case -1:
		return left + 1, true
	default:
		panic("unknown status")
	}

}

var NextCap func(old int) int = default_extend

// 手动扩张，自定义扩容算法修改NextCap
func Append[T any](src []T) []T {
	old := len(src)
	newcap := NextCap(old)
	if newcap <= old {
		panic("append failed new cap small than old cap")
	}
	newcp := make([]T, newcap)
	copy(newcp[:old], src)
	return newcp
}
func default_extend(old int) int {
	newcap := 2 * old
	if old > 512 {
		newcap = old + old/2
		if old < 2048 {
			old += old * (2048 - old) * 3 / 4
		}
	}
	return newcap
}
