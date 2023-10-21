package kits

import (
	"encoding/binary"
	"fmt"
	"os"
)

type EnumMode interface {
	Int() uint64
	Bytes() []byte
}
type Mode string

func pow(src uint64, step uint8) uint64 {
	if src == 0 {
		return 0
	} else if step == 0 {
		return 1
	}
	var i uint8
	var ans uint64 = 1
	for i = 0; i < step; i++ {
		ans *= src
	}
	return ans
}
func (s *Mode) Int() uint64 {
	var (
		ans            uint64 = 0
		code, lastcode uint8
	)
	for k, ele := range *s {
		code = uint8(ele - '0')
		if code < 16 {
			if k%2 == 0 {
				lastcode += code
			} else {
				lastcode += code * 16
				ans += uint64(lastcode) * pow(256, uint8(k/2))
				lastcode = 0
			}

		} else {
			return 0
		}
	}
	if lastcode > 0 {
		ans += uint64(lastcode) * pow(256, uint8(len(*s)/2))
	}
	return ans
}
func (s *Mode) Bytes() []byte {
	var (
		lastcode, code uint8
	)
	ans := make([]byte, len(*s)/2)
	for k, ele := range *s {
		code = uint8(ele - '0')
		if code < 16 {
			if k%2 == 0 {
				lastcode += code
			} else {
				lastcode += code * 16
				ans[k/2] = lastcode
				lastcode = 0
			}

		} else {
			return nil
		}
	}
	if lastcode > 0 {
		ans[(len(*s)/2)-1] = lastcode
	}
	return ans
}

// enum value limit 0~4
type SmallMode string

func (s *SmallMode) Int() uint64 {
	var (
		ans            uint64 = 0
		code, lastcode uint8
	)
	for k, ele := range *s {
		code = uint8(ele - '0')
		if code < 4 {
			lastcode += code * uint8(pow(4, uint8(k%4)))
			if k%4 == 3 {
				ans += uint64(lastcode) * pow(256, uint8(k/4))
				lastcode = 0
			}
			// fmt.Printf("code is %d,lastcode %d\n", code, lastcode)
		} else {
			return 0
		}
	}
	if lastcode > 0 {
		ans += uint64(lastcode) * pow(256, uint8(len(*s)/4))
	}
	return ans
}

func (s *SmallMode) Bytes() []byte {
	var (
		lastcode, code uint8
		ans            []byte
	)
	if len(*s)%4 == 0 {
		ans = make([]byte, len(*s)/4)
	} else {
		ans = make([]byte, len(*s)/4+1)
	}

	for k, ele := range *s {
		code = uint8(ele - '0')
		if code < 16 {
			lastcode += code * uint8(pow(4, uint8(k%4)))
			// fmt.Printf("code is %d,lastcode %d\n", code, lastcode)
			if k%4 == 3 {
				ans[k/4] = lastcode
				lastcode = 0
			}

		} else {
			return nil
		}
	}
	if lastcode > 0 {
		ans[(len(*s) / 4)] = lastcode
	}
	return ans
}

type CustomMode struct {
	Split_Num uint8
	Val       []byte
}

func (s *CustomMode) Int() uint64 {
	return binary.LittleEndian.Uint64(s.Val)
}

func (s *CustomMode) Bytes() []byte {
	return s.Val
}
func GetCustomMode(v string, splitval uint8) EnumMode {
	var (
		lastcode, code, step uint8
		ans                  []byte = make([]byte, 8)
	)
	switch splitval {
	case 2:
		step = 8
	case 4:
		step = 4
	case 16:
		step = 2
	default:
		fmt.Fprintf(os.Stderr, "%d not vaild split number", splitval)
		return nil
	}
	for k, ele := range v {
		code = uint8(ele - '0')
		if code < splitval {
			lastcode += code * uint8(pow(uint64(splitval), uint8(k)%step))
			// fmt.Printf("code %v lastcode %v\n", code, lastcode)
			if uint8(k)%step == step-1 {
				ans[uint8(k)/step] = lastcode
				lastcode = 0
			}
		} else {
			fmt.Fprintf(os.Stderr, "code %d large than %d", code, splitval)
			return nil
		}
	}
	if lastcode > 0 {
		ans[len(v)/int(step)] = lastcode
	}
	return &CustomMode{Val: ans}
}

// recover origin enum value
func ParseMode(v uint64, split uint8) []uint8 {
	var step uint8
	switch split {
	case 2:
		step = 8
	case 4:
		step = 4
	case 16:
		step = 2
	default:
		return nil
	}
	var (
		code, lang uint8
		ansbuffer  []uint8
		startpos   uint8 = 0
	)
	if v%256 == 0 {
		ansbuffer = make([]uint8, step*uint8(v/256))
	} else {
		ansbuffer = make([]uint8, step*uint8(v/256+1))
	}
	// fmt.Printf("step %d,split %d\n", step, split)
	for v > 0 { //split to packet for parse
		code = uint8(v % 256)
		v /= 256
		lang = parsemode(code, step, split, ansbuffer[startpos:])
		// fmt.Printf("code %d,lang %d\n", code, lang)
		startpos += lang
	}
	return ansbuffer[:startpos]
}
func parsemode(v, step, split uint8, buffer []uint8) uint8 {
	var i, lang uint8
	lang = 0
	for i = 0; i < step; i++ {
		if v > 0 {
			buffer[i] = v % split
			v /= split
			lang++
		} else {
			break
		}
	}
	return lang
}
