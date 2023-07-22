package common

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

var (
	SYS_STDOUT         = os.Stdout
	Default_Err_Prefix = Prefix{Prefix: "error", Time: true}
	Time_Format        = time.Layout
)

type Prefix struct {
	Prefix string   `json:"prefix"`
	Time   bool     `json:"time"`
	Tags   []string `json:"tags"`
}

func make_prefix(data *Prefix, format string) string {
	ans := ""
	if len(data.Prefix) > 0 {
		ans = "[" + data.Prefix + "]"
	}
	if data.Time && len(data.Tags) > 0 {
		ans += " " + time.Now().Format(format) + " " + strings.Join(data.Tags, " ")
	} else if data.Time || (data.Tags == nil || len(data.Tags) == 0) {
		if data.Time {
			ans += " " + time.Now().Format(format)
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
func OutputWithPrefix(pfx *Prefix, format string, args ...any) {
	format = make_prefix(pfx, Time_Format) + format
	if format[len(format)-1] != '\n' {
		format += "\n"
	}
	fmt.Fprintf(SYS_STDOUT, format, args...)
}

// read config
type Cnf interface {
	Marshal() ([]byte, error)
	Unmarshal([]byte) error
}

func ReadConfig(filepath string, v Cnf) error {
	content, err := ioutil.ReadFile(filepath)
	if err == nil {
		err = v.Unmarshal(content)
	}
	return err
}

// the independent outputer 22 july 2023
type Outputer struct {
	prefix Prefix
	finfo  *os.File
	Format string //format time string
}

func NewOutputer(out string, pre Prefix) *Outputer {
	var finfo_ *os.File
	var err error
	if strings.ToLower(out) == "std" {
		finfo_ = os.Stdout
	} else {
		finfo_, err = os.OpenFile(out, os.O_WRONLY|os.O_CREATE, 0640)
		if err != nil {
			OutputWithPrefix(&Default_Err_Prefix, "set outputer failed,%s", err.Error())
		}
	}
	return &Outputer{prefix: pre, finfo: finfo_}
}

func (s *Outputer) Output(format string, args ...any) {
	if format[len(format)-1] != '\n' {
		format += "\n"
	}
	fmt.Fprintf(s.finfo, format, args...)
}
func (s *Outputer) OutputWithPrefix(format string, args ...any) {
	format = make_prefix(&s.prefix, s.Format) + format
	if format[len(format)-1] != '\n' {
		format += "\n"
	}
	fmt.Fprintf(s.finfo, format, args...)
}
