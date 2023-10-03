package common

import (
	"strconv"
	"strings"
)

func Inet_Aton(src string) uint32 {
	var (
		ans      uint32 = 0
		tinumber uint64
		err      error
	)
	if len(src) < 7 {
		return 0
	}
	ptrarr := strings.Split(src, ".")
	if len(ptrarr) != 4 {
		return 0
	}
	for k, ele := range ptrarr {
		if len(ele) < 1 {
			return 0
		}
		tinumber, err = strconv.ParseUint(ele, 10, 8)
		if err != nil {
			return 0
		}
		ans += uint32(tinumber) * Pow[uint32](256, uint8(k))
	}
	return ans
}
