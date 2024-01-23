package kits

import (
	"fmt"
	"reflect"
	"strings"
)

const (
	SET = iota + 1
	CONVERT
	PASS
)

type str_err string

func (s str_err) Error() string {
	return string(s)
}

type action struct {
	src, dst reflect.Value
	_type    uint8
}
type Binder struct {
	flagname string
}

func NewBinder(name string) *Binder {
	return &Binder{flagname: name}
}
func (s *Binder) match(src map[string]any, dst reflect.Value) (acts []action, err error) {
	tp := dst.Type()
	if tp.Kind() != reflect.Struct {
		err = str_err("dst is must be struct,not " + tp.Kind().String())
		return
	}
	count := dst.NumField()
	if len(src) > 0 && count > 0 {
		var (
			tagname string
			// ok            bool
			vtp, vkeytp reflect.Type
			vvl         reflect.Value
			field       reflect.StructField
			_acts       []action
			vmap        map[string]any
		)
		acts = make([]action, 0, count)
		for i := 0; i < count; i++ {
			field = tp.Field(i)
			tagname = field.Tag.Get(s.flagname)
			if len(tagname) == 0 {
				tagname = strings.ToLower(field.Name)
			}
			fmt.Println("[debug] find " + tagname)
			if ve, ok := src[tagname]; ok {
				vvl = reflect.ValueOf(ve)
				vtp = vvl.Type()
				if vtp.AssignableTo(field.Type) {
					acts = append(acts, action{src: vvl, dst: dst.Field(i), _type: SET})
				} else if vtp.ConvertibleTo(field.Type) {
					acts = append(acts, action{src: vvl, dst: dst.Field(i), _type: CONVERT})
				} else {
					//if is map and struct
					if vtp.Kind() == reflect.Map && field.Type.Kind() == reflect.Struct {
						vkeytp = vtp.Key()
						if vkeytp.Kind() == reflect.Interface {
							vmap = transfermap(ve.(map[any]any))
						} else if vkeytp.Kind() == reflect.String {
							vmap = ve.(map[string]any)
						} else {
							err = str_err("key must be string not " + vkeytp.Kind().String())
							return
						}
						_acts, err = s.match(vmap, dst.Field(i))
						if err != nil {
							return
						}
						if len(_acts) > 0 {
							acts = append(acts, _acts...)
						}
					}
				}
			}
		}

	} else {
		fmt.Println("[debug]", len(src), count)
	}
	return
}
func (s *Binder) Bind(src map[string]any, dst any) (err error) {
	vl := reflect.ValueOf(dst)
	vtp := vl.Type()
	if vtp.Kind() != reflect.Pointer {
		err = str_err("dst muste be pointer")
		return
	}
	vtp = vtp.Elem()
	if vtp.Kind() == reflect.Pointer {
		err = str_err("dst mutst be map or struct dst")
		return
	}
	var acts []action
	acts, err = s.match(src, vl.Elem())
	if err == nil {
		bindaction(acts)
	}
	return
}
func bindaction(target []action) {
	for _, e := range target {
		if e._type == CONVERT {
			e.src = e.src.Convert(e.dst.Type())
		} else if e._type == PASS {
			continue
		}
		e.dst.Set(e.src)
	}
}
func transfermap(src map[any]any) map[string]any {
	var ans = make(map[string]any)
	for k, v := range src {
		ans[k.(string)] = v
	}
	return ans
}
