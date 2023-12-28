package kits

import (
	"fmt"
	"reflect"
	"strings"
)

const (
	DEFAULT = 0
	SET     = 1
	CONVERT = 2
)

type node struct {
	src, dst reflect.Value
	bind     uint8
	child    map[string]*node
}

func bind2struct[T any](src map[string]any, v *T) error {
	reg_table := make(map[string]*node)
	fmt.Println(reflect.ValueOf(src["name"]).Type().String())
	err := _prepare_bind(&reg_table, src, reflect.ValueOf(v).Elem())
	if err == nil {
		_bind_value(reg_table)
	}
	return err
}

// src should is map[string]any,dst should be struct
func _prepare_bind(reg_table *map[string]*node, src map[string]any, dst reflect.Value) error {
	filedcount := dst.NumField()
	if filedcount > 0 {
		var (
			i       = 0
			df      reflect.StructField
			dfvl    reflect.Value
			tagname string
			mapval  reflect.Value
			err     error
		)
		for i < filedcount {
			df = dst.Type().Field(i)
			if df.Anonymous {
				continue
			}
			for k, v := range src {
				tagname = compare_name(k, df.Name, df.Tag.Get("flag"))
				if len(tagname) == 0 {
					continue
				}
				if _, ok := (*reg_table)[tagname]; ok {
					return str_error("redefined tagname " + tagname)
				}
				mapval = reflect.ValueOf(v)
				dfvl = dst.Field(i)
				if !dfvl.CanSet() {
					return str_error("tag " + tagname + " can't set")
				}
				(*reg_table)[tagname] = &node{child: make(map[string]*node)}
				if mapval.Type().Kind() == reflect.Map && df.Type.Kind() == reflect.Struct {
					err = _prepare_bind(&(*reg_table)[tagname].child, v.(map[string]any), dfvl)
					if err != nil {
						return err
					}
				} else {
					if mapval.Type().AssignableTo(dfvl.Type()) {
						(*reg_table)[tagname].src = mapval
						(*reg_table)[tagname].dst = dfvl
						(*reg_table)[tagname].bind = SET
					} else if mapval.CanConvert(dfvl.Type()) {
						(*reg_table)[tagname].src = mapval
						(*reg_table)[tagname].dst = dfvl
						(*reg_table)[tagname].bind = CONVERT
					} else {
						return str_error("src can't assignable to dst;tagname " + tagname + " " + mapval.Type().String() + " " + dfvl.Type().String())
					}
				}
			}

			i++
		}
	}
	return nil
}
func _bind_value(reg_table map[string]*node) {
	for _, e := range reg_table {
		if e.bind == SET {
			e.dst.Set(e.src)
		} else if e.bind == DEFAULT {
			_bind_value(e.child)
		} else if e.bind == CONVERT {
			e.dst.Set(e.src.Convert(e.dst.Type()))
		}
	}
}
func compare_name(mapname, name, tagname string) string {
	if tagname == "!" {
		return ""
	}
	index := strings.IndexByte(tagname, ';')
	if index >= 0 {
		tagname = tagname[:index]
	}
	if len(tagname) == 0 {
		tagname = strings.ToLower(name)
	}
	if tagname == mapname {
		return tagname
	}
	return ""
}

// unmarsha but not cover to the value already existed
func UnmarshaAgain[T any](fun func([]byte, any) error, v []byte, dst *T) error {
	ans := make(map[string]any)
	err := fun(v, &ans)
	if err == nil {
		err = bind2struct(ans, dst)
	}
	return err
}
