package common

import (
	"fmt"
	"os"
	"plugin"
	"reflect"

	"github.com/oswaldoooo/cmicro/kits"
)

// name is function's name which in your plugin go filem, target is tell plugin loador your function how to made
// here is an example
// var testone func(string)string
// var test_one PluginOption={"TestOne",testone}
type PluginOption struct {
	Name   string
	Target any //store function
}

// load plugin,return the real load count
func LoadPlugin(src *plugin.Plugin, opts ...*PluginOption) (count int) { //return successful find number
	if src == nil || len(opts) == 0 {
		return
	}
	var (
		ottp reflect.Type
		sy   plugin.Symbol
		err  error
	)
	for _, opt := range opts {
		ottp = reflect.TypeOf(opt.Target)
		if len(opt.Name) > 0 && ottp.Kind() == reflect.Func {
			//filter options,name must be set,and target must be func
			sy, err = src.Lookup(opt.Name)
			if err == nil {
				if kits.Tofunc(sy, &opt.Target) {
					count++
				} else {
					fmt.Fprintln(os.Stderr, "[tofunc failed]")
				}
			} else {
				fmt.Fprintln(os.Stderr, "[look up error]", err.Error())
			}
		} else {
			fmt.Fprintln(os.Stderr, "option is invaild")
		}
	}
	return
}
