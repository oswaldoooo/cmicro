package kits

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"errors"
	"flag"
	"strings"
)

type str_error string

func (s str_error) Error() string {
	return string(s)
}

// reflect do function
func ReflectInfo(v, dst any, args ...any) (err error) {
	vl := reflect.ValueOf(v)
	if reflect.TypeOf(v).Kind() != reflect.Func {
		err = str_error("input v is not function")
		return
	}
	if vl.Type().NumIn() > 0 && len(args) == vl.Type().NumIn() {
		tp := vl.Type()
		var (
			tpt     reflect.Type
			argsarr = make([]reflect.Value, 0, tp.NumIn())
		)
		for i := 0; i < len(args); i++ {
			tpt = tp.In(i)
			if tpt.Kind() != reflect.TypeOf(args[i]).Kind() {
				err = str_error("args not correct")
				return
			} else if tpt.Kind() == reflect.Pointer && tpt.Elem().Kind() != reflect.TypeOf(args[i]).Elem().Kind() {
				err = str_error("args not correct")
				return
			}
			argsarr = append(argsarr, reflect.ValueOf(args[i]))
		}
		vlarr := vl.Call(argsarr)
		if dst != nil && reflect.TypeOf(dst).Kind() == reflect.Pointer {
			var min_len int = len(vlarr)
			anstp := reflect.ValueOf(dst)
			if min_len == 1 && vlarr[0].Type().Kind() != reflect.Interface && vlarr[0].Kind() == anstp.Elem().Type().Kind() { //single value set
				anstp.Elem().Set(vlarr[0])
			} else if min_len > 0 && anstp.Elem().Type().Kind() == reflect.Struct { //many value set to struct
				anstpt := anstp.Elem().Type()
				var count = anstpt.NumField()
				var i, j int = 0, 0 //i dst,j vlarr
				var (
					anstype, vltype reflect.Type
					ok              bool
				)
				for j < min_len && i < count {
					anstype = anstpt.Field(i).Type
					vltype = vlarr[j].Type()
					if anstype.Kind() == vltype.Kind() && vltype.Kind() != reflect.Pointer && vltype.Kind() != reflect.Interface {
						if vltype.Kind() == reflect.Struct {

							// fmt.Println(i, j, reflect.Indirect(anstp.Elem().Field(i).Addr()).CanSet())
							ok, err = struct_set(anstp.Elem().Field(i), vlarr[j])
							if err != nil {
								return
							} else if !ok {
								err = str_error("struct don't compare")
								return
							}
						} else {
							// fmt.Println(anstp.Elem().Field(i).CanSet())
							anstp.Elem().Field(i).Set(vlarr[j])
						}

						j++
					}
					i++
				}
			} else {
				// fmt.Println("[debug]", min_len, anstp.Type().Kind(), reflect.TypeOf(dst).Elem().Kind())
				err = str_error("ans don't compare")
			}
			// fmt.Println(vlarr[0].Type().Implements(anstp.Type().Elem().Field(0).Type), anstp.Kind().String())
		} else if dst == nil && len(vlarr) == 0 {
			return
		} else {
			err = str_error("dst is null")
		}
		// fmt.Println(vlarr)
	} else {
		err = str_error("args not correct")
	}
	return
}

// compare and set struct
func struct_set(dst, src reflect.Value) (ok bool, err error) {
	dn, sn := dst.NumField(), src.NumField()
	ok = false               //issame ?
	if dn != sn || dn == 0 { //not same struct
		return
	}
	var (
		vp reflect.Value
	)
	for i := 0; i < dn; i++ {
		vp = dst.Field(i)
		if dst.Field(i).Type().Kind() != src.Field(i).Type().Kind() {
			return
		} else if dst.Field(i).Kind() == reflect.Pointer { //don't support pointer transport
			return
		} else if dst.Field(i).Kind() == reflect.Struct {
			ok, err = struct_set(dst.Field(i), src.Field(i))
			if !ok {
				return
			}
		} else if vp.CanSet() {
			vp.Set(src.Field(i))
		} else {
			err = str_error("dst " + strconv.FormatInt(int64(i), 10) + " field can't be set,elem ")
			return
		}
	}
	ok = true
	for i := 0; i < dn; i++ {
		dst.Field(i).Set(src.Field(i))
	}
	return
}

// function type set
func Tofunc[T any](src any, dst *T) bool {
	dtp := reflect.TypeOf(dst).Elem()
	// dvl := reflect.ValueOf(dst).Elem()
	stp := reflect.TypeOf(src)
	svl := reflect.ValueOf(src)
	// fmt.Println(dtp.Kind().String(), dvl.IsNil(), reflect.TypeOf(src).Kind())
	if svl.IsNil() {
		return false
	} else if stp.Kind() != reflect.Func || dtp.Kind() != reflect.Func {
		if svl.CanConvert(dtp) {
			var ok bool
			*dst, ok = src.(T)
			if ok {
				return ok
			}
		}
		panic("args must be func type")
	} else if (stp.NumIn() != dtp.NumIn()) || (stp.NumOut() != dtp.NumOut()) { //compare in and out args number
		return false
	}
	if stp.NumIn() > 0 { //compare more details
		for i := 0; i < stp.NumIn(); i++ {
			if stp.In(i) != dtp.In(i) {
				return false
			}
		}
	}
	if stp.NumOut() > 0 {
		for i := 0; i < stp.NumOut(); i++ {
			if stp.Out(i) != dtp.Out(i) {
				return false
			}
		}
	}
	var ok bool
	*dst, ok = src.(T)
	return ok
}

// if struct copy,dst must be large side,src must be small side
func memcopy[T, B any](dst *T, src *B) (err error) {
	dtp, stp := reflect.TypeOf(dst), reflect.TypeOf(src)
	dvl, svl := reflect.ValueOf(dst), reflect.ValueOf(src)
	if dvl != svl {
		if dtp == stp {
			if dvl.Elem().CanSet() {
				dvl.Elem().Set(svl.Elem())
			}
		} else if dtp.Elem().Kind() == stp.Elem().Kind() && dtp.Elem().Kind() == reflect.Struct {
			//todo: struct set value is not complete yet
			dtp, stp = dtp.Elem(), stp.Elem()
			dvl, svl = dvl.Elem(), svl.Elem()
			var (
				i, j                 = 0, 0
				dcount, scount       = dtp.NumField(), stp.NumField()
				parse_pos      []int = make([]int, scount)
				dftp, sftp     reflect.Type
			)
			for i < dcount && j < scount {
				dftp, sftp = dtp.Field(i).Type, stp.Field(j).Type
			compare:
				if dftp == sftp {
					if dftp.Kind() == reflect.Pointer {
						dftp, sftp = dftp.Elem(), sftp.Elem()
						goto compare
					} else if dftp.Kind() == reflect.Struct {

					}
					if dvl.Field(i).CanSet() {
						parse_pos[j] = i
						j++
					}
				}
				i++
			}
			if j != scount {
				err = str_error("parse failed")
			} else {
				for i := 0; i < scount; i++ {
					dvl.Field(parse_pos[i]).Set(svl.Field(i))
				}
			}
		} else {
			err = str_error("type not copiable")
		}
	}
	return
}

// problem: copy need sure whether support, if don't support, can't copy any val to dst's field
const ( //method
	SETVAL         = 0x01
	SETVAL_        = 0x02 //it's pointer need elem() and set
	SETVAL_STRUCT  = 0x03 //it's struct instance
	SETVAL_STRUCT_ = 0x04 //it's struct pointer,need elem
)

type copy_option struct {
	pos                int
	method             uint8
	child_copy_options []copy_option
}

func StructCopy[T, B any](dst *T, src *B) error {
	struct_copy(reflect.TypeOf(dst).Elem(), reflect.TypeOf(src).Elem(), reflect.ValueOf(dst).Elem(), reflect.ValueOf(src).Elem())
	return nil
}

// copy data with compatible
func struct_copy(dtp, stp reflect.Type, dvl, svl reflect.Value) (err error) {
	if opts := struct_copy_prepare(dtp, stp, dvl); opts != nil {
		if !struct_copy_stmt(dvl, svl, opts) {
			err = str_error("parse failed")
		}
	} else {
		err = str_error("can't copy src to dst")
	}
	return
}

func struct_copy_prepare(dtp, stp reflect.Type, dvl reflect.Value) []copy_option {
	var (
		ans            []copy_option
		dcount, scount int
	)
	if dtp.Kind() != reflect.Struct || stp.Kind() != reflect.Struct {
		fmt.Fprintln(os.Stderr, "src or dst is not struct")
		return nil
	}
	dcount, scount = dtp.NumField(), stp.NumField()
	if scount > dcount {
		fmt.Fprintln(os.Stderr, "src field number can't over dst field number")
		return nil
	}
	ans = make([]copy_option, 0, scount)
	var (
		i, j             = 0, 0
		umask      uint8 = 0
		dftp, sftp reflect.Type
		dfvl       reflect.Value
		child_co   []copy_option = nil
	)
	for i < dcount && j < scount {
		dftp, sftp = dtp.Field(i).Type, stp.Field(j).Type
		dfvl = dvl.Field(i)
		umask = 0
		if dftp == sftp { //type same
			if dfvl.CanSet() {
				ans = append(ans, copy_option{i, SETVAL, nil})
				j++
			} else if dftp.Kind() == reflect.Pointer && dfvl.Elem().CanSet() {
				ans = append(ans, copy_option{i, SETVAL_, nil})
				j++
			}
		} else if dftp.Kind() == sftp.Kind() { //type don't same
			if dftp.Kind() == reflect.Pointer && dftp.Elem().Kind() == reflect.Struct && sftp.Elem().Kind() == reflect.Struct {
				// this is *struct and *struct
				umask = 1
				dftp, sftp = dftp.Elem(), sftp.Elem()
				dfvl = dfvl.Elem()
			}
			child_co = struct_copy_prepare(dftp, sftp, dfvl)
			if child_co != nil {
				ans = append(ans, copy_option{i, SETVAL_STRUCT + umask, child_co})
				j++
			}
		}
		i++
	}
	if len(ans) < scount { //src can't copy complete
		fmt.Fprintln(os.Stderr, "[iherit error]failed iherit")
		return nil
	}
	return ans
}

// totest:
func struct_copy_stmt(dvl, svl reflect.Value, opts []copy_option) bool {
	if opts == nil || len(opts) == 0 || len(opts) != svl.NumField() {
		return false
	}
	for key, opt := range opts {
		switch opt.method {
		case SETVAL:
			dvl.Field(opt.pos).Set(svl.Field(key))
		case SETVAL_:
			dvl.Field(opt.pos).Elem().Set(svl.Field(key).Elem())
		case SETVAL_STRUCT:
			if !struct_copy_stmt(dvl.Field(opt.pos), svl.Field(key), opt.child_copy_options) {
				return false
			}
		case SETVAL_STRUCT_:
			if !struct_copy_stmt(dvl.Field(opt.pos).Elem(), svl.Field(key).Elem(), opt.child_copy_options) {
				return false
			}
		default:
			return false
		}
	}
	return true
}

//bind to flag
var (
	Stderr = os.Stderr
	Stdout = os.Stdout
)

func FlagParse[T any](fs *flag.FlagSet, v *T) error {
	if reflect.TypeOf(v).Elem().Kind() != reflect.Struct {
		return errors.New("input args is not struct")
	}
	vl := reflect.ValueOf(v).Elem()
	if !vl.CanSet() {
		return errors.New("v is not can set")
	}
	return flagparse(reflect.ValueOf(fs), vl)
}
func flagparse(fs reflect.Value, vl reflect.Value) error {
	tp := vl.Type()
	fn := tp.NumField()
	if fn == 0 {
		return nil
	}
	var (
		ff          reflect.StructField
		fvl, ffvl   reflect.Value
		ftp         reflect.Type
		bindptr     reflect.Value
		flagname    string
		default_val reflect.Value
		argslist    []reflect.Value = make([]reflect.Value, 4)
	)
	// argslist[0] = vl
	for i := 0; i < fn; i++ {
		fvl = vl.Field(i)
		ff = tp.Field(i)
		ftp = ff.Type
		switch ftp.Kind() {
		case reflect.Int:
			ffvl = fs.MethodByName("IntVar")
			default_val = reflect.ValueOf(0)
		case reflect.Float64:
			ffvl = fs.MethodByName("Float64Var")
			default_val = reflect.ValueOf(0.00)
		case reflect.Bool:
			ffvl = fs.MethodByName("BoolVar")
			default_val = reflect.ValueOf(false)
		case reflect.String:
			ffvl = fs.MethodByName("StringVar")
			default_val = reflect.ValueOf("")
		case reflect.Struct:
			flagparse(fs, fvl)
			continue
		default:
			fmt.Fprintln(Stderr, "don't support", ftp.Kind())
			continue
		}

		flagname = ff.Tag.Get("flag")
		if len(flagname) == 0 {
			flagname = strings.ToLower(ff.Name)
		} else {
			if newdefault, ok := parseflagtag(&flagname, ftp.Kind()); ok {
				default_val = newdefault
			}
		}
		//debug
		fmt.Fprintln(Stdout, "use flag", flagname)
		//
		bindptr = fvl.Addr()
		argslist[0] = bindptr
		argslist[1] = reflect.ValueOf(flagname)
		argslist[2] = default_val
		argslist[3] = reflect.ValueOf(strings.ReplaceAll(flagname, "_", " "))
		ffvl.Call(argslist)
	}
	return nil
}
func parseflagtag(tagstr *string, tp reflect.Kind) (ans reflect.Value, ok bool) {
	if strings.Count(*tagstr, ";") == 1 {
		next := strings.IndexByte(*tagstr, ';')
		val := (*tagstr)[next+1:]
		if len(val) > 0 {
			fmt.Fprintln(Stdout, "find value", val)
			switch tp {
			case reflect.Int:
				newvl, err := strconv.Atoi(val)
				if err == nil {
					ok = true
					ans = reflect.ValueOf(newvl)
				}
			case reflect.Bool:
				newvl, err := strconv.ParseBool(val)
				if err == nil {
					ok = true
					ans = reflect.ValueOf(newvl)
				}
			case reflect.Float64:
				newvl, err := strconv.ParseFloat(val, 10)
				if err == nil {
					ok = true
					ans = reflect.ValueOf(newvl)
				}
			case reflect.String:
				ans = reflect.ValueOf(val)
				ok = true
			}
		}
		*tagstr = (*tagstr)[:next]
	}
	return
}
