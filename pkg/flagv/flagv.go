package flagv

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type flagnode struct {
	ops         []Option
	default_val any
	intptr      *int
	floatptr    *float64
	stringptr   *string
	boolptr     *bool
	iscall      bool //iscall this flag
	_type       uint
	usage       string
}

const (
	INT    uint = uint(reflect.Int)
	FLOAT  uint = uint(reflect.Float64)
	BOOL   uint = uint(reflect.Bool)
	STRING uint = uint(reflect.String)
)

type FlagSet struct {
	kv      map[string]string //store node,after parse and bind it
	bind_kv map[string]*flagnode
	ops     []string //store kv with options
}

type Option func(string, bool, *FlagSet) //key,iscall,root

func (s *FlagSet) Parse(args []string) (err error) {
	if len(args) > 0 {
		var (
			argname           string
			callarg, callhelp bool
		)
		for i := 0; i < len(args); i++ {
			if args[i] == "--help" {
				callhelp = true
				continue
			}
			if strings.HasPrefix(args[i], "--") {
				argname = args[i][2:]
				callarg = true
			} else if strings.HasPrefix(args[i], "-") {
				argname = args[i][1:]
				callarg = true
			} else if callarg {
				//set argval
				s.kv[argname] = args[i]
				argname = ""
				callarg = false
			}
			if callarg {
				if len(argname) > 0 {
					s.kv[argname] = ""
					// fmt.Println("set args " + argname)
				} else {
					err = errors.New("empty argname")
					return
				}
			}
		}
		// var tp reflect.Type
		//must wait first bind end and then start do options
		for key, val := range s.bind_kv {
			if len(val.ops) == 0 {
				if v, ok := s.kv[key]; ok {
					// fmt.Println("set value "+key, v)
					val.setvraw(v)
					val.iscall = true
				} else {
					//use default value
					fmt.Println("set default " + key)
					val.setv(val.default_val)
				}
			}
		}
		var _node *flagnode
		//do option
		for _, key := range s.ops {
			_node = s.bind_kv[key]
			for _, op := range _node.ops {
				op(key, _node.iscall, s)
			}
		}
		if callhelp {
			Usage(s)
			os.Exit(0)
		}
	}
	return
}
func (s *FlagSet) StringVar(p *string, name string, value string, usage string, ops ...Option) {
	if _, ok := s.bind_kv[name]; ok {
		panic("redefined " + name)
	}
	s.bind_kv[name] = &flagnode{ops: ops, default_val: value, usage: usage, _type: STRING, stringptr: p}
	s.ops = append(s.ops, name)
}
func (s *FlagSet) IntVar(p *int, name string, value int, usage string, ops ...Option) {
	if _, ok := s.bind_kv[name]; ok {
		panic("redefined " + name)
	}
	s.bind_kv[name] = &flagnode{ops: ops, default_val: value, usage: usage, _type: INT, intptr: p}
	s.ops = append(s.ops, name)
}
func (s *FlagSet) BoolVar(p *bool, name string, value bool, usage string, ops ...Option) {
	if _, ok := s.bind_kv[name]; ok {
		panic("redefined " + name)
	}
	s.bind_kv[name] = &flagnode{ops: ops, default_val: value, usage: usage, _type: BOOL, boolptr: p}
	s.ops = append(s.ops, name)
}
func (s *FlagSet) Float64Var(p *float64, name string, value float64, usage string, ops ...Option) {
	if _, ok := s.bind_kv[name]; ok {
		panic("redefined " + name)
	}
	s.bind_kv[name] = &flagnode{ops: ops, default_val: value, usage: usage, _type: FLOAT, floatptr: p}
	s.ops = append(s.ops, name)
}
func NewFlagV() *FlagSet {
	return &FlagSet{kv: make(map[string]string), bind_kv: make(map[string]*flagnode)}
}

// options

// (left==right)?iftrue:ifalse
// delive: if set delive if value is false not do any thing
func Equal[T any](left string, right any, iftrue, ifalse T, delive bool) Option {
	rtp := reflect.TypeOf(iftrue)
	if rtp.Kind() == reflect.Pointer || rtp.Kind() == reflect.Array || rtp.Kind() == reflect.Slice || rtp.Kind() == reflect.Map || rtp.Kind() == reflect.Struct || rtp.Kind() == reflect.Chan {
		panic("right value is pointer,not support " + rtp.Kind().String() + " be rightvalue")
	}
	tp := reflect.TypeOf(right)
	if tp.Kind() == reflect.Pointer || tp.Kind() == reflect.Array || tp.Kind() == reflect.Slice || tp.Kind() == reflect.Map || tp.Kind() == reflect.Struct || tp.Kind() == reflect.Chan {
		panic("right value is pointer,not support " + rtp.Kind().String() + " be rightvalue")
	}

	return func(s1 string, iscall bool, fs *FlagSet) {
		if !iscall {
			//set default
			nownode := fs.bind_kv[s1]
			var rightval any
			node := fs.bind_kv[left]
			var uselink bool = false
			if node._type != uint(tp.Kind()) {
				if tp.Kind() == reflect.String {
					uselink = true
				} else {
					panic(fmt.Sprintf("value type is not same left %d right %s", node._type, tp.Kind().String()))
				}
			}

			if uselink {
				if _, ok := fs.bind_kv[right.(string)]; !ok {
					panic("right value link is not be registered")
				}
				rv := fs.bind_kv[right.(string)]
				switch rv._type {
				case STRING:
					rightval = *rv.stringptr
				case INT:
					rightval = *rv.intptr
				case BOOL:
					rightval = *rv.boolptr
				case FLOAT:
					rightval = *rv.floatptr
				}
			} else {
				rightval = right
			}
			fmt.Println("[debug]", left, rightval, node.getv())
			if rightval == node.getv() {
				nownode.default_val = iftrue
			} else if !delive {
				nownode.default_val = ifalse
			}
		}
	}
}
func BiggerThan[T int | int32 | int64 | float32 | float64 | string](left string, right any, iftrue, ifalse T, delive bool) Option {
	rtp := reflect.TypeOf(iftrue)
	if rtp.Kind() == reflect.Pointer || rtp.Kind() == reflect.Array || rtp.Kind() == reflect.Slice || rtp.Kind() == reflect.Map || rtp.Kind() == reflect.Struct || rtp.Kind() == reflect.Chan {
		panic("right value is pointer,not support " + rtp.Kind().String() + " be rightvalue")
	}
	tp := reflect.TypeOf(right)
	if tp.Kind() != reflect.Int && tp.Kind() != reflect.Float64 && tp.Kind() != reflect.String {
		panic("right value is pointer,not support " + rtp.Kind().String() + " be rightvalue")
	}
	return func(s string, b bool, fs *FlagSet) {
		if !b {
			node := fs.bind_kv[left]
			var rnode, nownode *flagnode
			nownode = fs.bind_kv[s]
			if nownode._type != uint(rtp.Kind()) {
				panic("type not same")
			}
			if tp.Kind() == reflect.String {
				rnode = fs.bind_kv[right.(string)]
				switch rnode._type {
				case INT, FLOAT:
				default:
					panic("right value type is not number")
				}
			} else {
				// val := right
				rnode = &flagnode{_type: uint(tp.Kind())}
				switch rnode._type {
				case INT:
					var val int = right.(int)
					rnode.intptr = &val
				case FLOAT:
					var val float64 = right.(float64)
					rnode.floatptr = &val
				}
			}
			var ans bool

			switch node._type {
			case INT:
				// fmt.Println(node._type, *node.intptr, rnode._type)
				if rnode._type != INT {
					ans = *node.intptr > int(*rnode.floatptr)
				} else {
					ans = *node.intptr > *rnode.intptr
				}
			case FLOAT:
				// fmt.Println(node._type, *node.floatptr, rnode._type)
				if rnode._type != INT {
					ans = *node.floatptr > float64(*rnode.intptr)
				} else {
					ans = *node.floatptr > *rnode.floatptr
				}
			}
			if ans {
				// reflect.ValueOf(nownode).Elem().Set(reflect.ValueOf(iftrue))
				nownode.default_val = iftrue
			} else if !delive {
				// reflect.ValueOf(nownode).Elem().Set(reflect.ValueOf(ifalse))
				nownode.default_val = ifalse
			}
		}
	}
}
func SmallThan[T int | int32 | int64 | float32 | float64 | string](left string, right any, iftrue, ifalse T, delive bool) Option {
	rtp := reflect.TypeOf(iftrue)
	if rtp.Kind() == reflect.Pointer || rtp.Kind() == reflect.Array || rtp.Kind() == reflect.Slice || rtp.Kind() == reflect.Map || rtp.Kind() == reflect.Struct || rtp.Kind() == reflect.Chan {
		panic("right value is pointer,not support " + rtp.Kind().String() + " be rightvalue")
	}
	tp := reflect.TypeOf(right)
	if tp.Kind() != reflect.Int && tp.Kind() != reflect.Float64 && tp.Kind() != reflect.String {
		panic("right value is pointer,not support " + rtp.Kind().String() + " be rightvalue")
	}
	return func(s string, b bool, fs *FlagSet) {
		if !b {
			node := fs.bind_kv[left]
			var rnode, nownode *flagnode
			nownode = fs.bind_kv[s]
			if nownode._type != uint(rtp.Kind()) {
				panic("type not same")
			}
			if tp.Kind() == reflect.String {
				rnode = fs.bind_kv[right.(string)]
				switch rnode._type {
				case INT, FLOAT:
				default:
					panic("right value type is not number")
				}
			} else {
				// val := right
				rnode = &flagnode{_type: uint(tp.Kind())}
				switch rnode._type {
				case INT:
					var val int = right.(int)
					rnode.intptr = &val
				case FLOAT:
					var val float64 = right.(float64)
					rnode.floatptr = &val
				}
			}
			var ans bool

			switch node._type {
			case INT:
				// fmt.Println(node._type, *node.intptr, rnode._type)
				if rnode._type != INT {
					ans = *node.intptr < int(*rnode.floatptr)
				} else {
					ans = *node.intptr < *rnode.intptr
				}
			case FLOAT:
				// fmt.Println(node._type, *node.floatptr, rnode._type)
				if rnode._type != INT {
					ans = *node.floatptr < float64(*rnode.intptr)
				} else {
					ans = *node.floatptr < *rnode.floatptr
				}
			}
			if ans {
				// reflect.ValueOf(nownode).Elem().Set(reflect.ValueOf(iftrue))
				nownode.default_val = iftrue
			} else if !delive {
				// reflect.ValueOf(nownode).Elem().Set(reflect.ValueOf(ifalse))
				nownode.default_val = ifalse
			}
		}
	}
}
func (s *flagnode) getv() any {
	switch s._type {
	case STRING:
		return *s.stringptr
	case INT:
		return *s.intptr
	case FLOAT:
		return *s.floatptr
	case BOOL:
		return *s.boolptr
	}
	return nil
}
func (s *flagnode) setv(v any) {
	switch s._type {
	case STRING:
		*s.stringptr = v.(string)
	case INT:
		*s.intptr = v.(int)
	case FLOAT:
		*s.floatptr = v.(float64)
	case BOOL:
		*s.boolptr = v.(bool)
	}
}
func (s *flagnode) setvraw(v string) {
	var err error
	switch s._type {
	case STRING:
		*s.stringptr = v
	case INT:
		*s.intptr, err = strconv.Atoi(v)
	case FLOAT:
		*s.floatptr, err = strconv.ParseFloat(v, 64)
	case BOOL:
		b := s.default_val.(bool)
		if b {
			*s.boolptr = false
		} else {
			*s.boolptr = true
		}
	}
	if err != nil {
		panic("set value failed " + err.Error())
	}
}

func Usage(fs *FlagSet) {
	usage := "commandline help:\n"
	if len(fs.bind_kv) > 0 {
		for k, v := range fs.bind_kv {
			usage += "-" + k + "\n  " + v.usage + " (default " + v.defautlval_str() + ")" + "\n"
		}
	}
	fmt.Print(usage)
}
func (s *flagnode) defautlval_str() string {
	switch s._type {
	case INT:
		return strconv.Itoa(s.default_val.(int))
	case FLOAT:
		return strconv.FormatFloat(s.default_val.(float64), 'f', -1, 64)
	case BOOL:
		return strconv.FormatBool(s.default_val.(bool))
	case STRING:
		return s.default_val.(string)
	}
	return ""
}
