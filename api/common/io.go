package common

import (
	"fmt"
	"os"
	"strings"
	"time"
)

var (
	SYS_STDOUT = os.Stdout
)

type Prefix struct {
	Prefix string   `json:"prefix"`
	Time   bool     `json:"time"`
	Tags   []string `json:"tags"`
}

func make_prefix(data *Prefix) string {
	ans := ""
	if len(data.Prefix) > 0 {
		ans = "[" + data.Prefix + "]"
	}
	if data.Time && len(data.Tags) > 0 {
		ans += " " + time.Now().Format(time.Kitchen) + " " + strings.Join(data.Tags, " ")
	} else if data.Time || (data.Tags == nil || len(data.Tags) == 0) {
		if data.Time {
			ans += " " + time.Now().Format(time.Kitchen)
		} else {
			ans += " " + strings.Join(data.Tags, " ")
		}
	}
	ans += " >> "
	return ans
}

func SetRelease(filepath string, isflush bool, mode os.FileMode) error {
	if isflush { //flush the target file
		_, err := os.Stat(filepath)
		if err == nil {
			err = os.Remove(filepath)
			if err != nil {
				return err
			}
		}
	}
	fid, err := os.OpenFile(filepath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, mode)
	if err == nil {
		SYS_STDOUT = fid
	}
	return err
}
func SetDebug() {
	SYS_STDOUT = os.Stdout
}
func Output(format string, args ...any) {
	if format[len(format)-1] != '\n' {
		format += "\n"
	}
	fmt.Fprintf(SYS_STDOUT, format, args...)
}
func OuptutWithPrefix(pfx *Prefix, format string, args ...any) {
	format = make_prefix(pfx) + format
	if format[len(format)-1] != '\n' {
		format += "\n"
	}
	fmt.Fprintf(SYS_STDOUT, format, args...)
}
