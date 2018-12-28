package main

import (
	"fmt"
	"github.com/JaredReisinger/drone-plugin-helper/env"
	// "os"
	// "reflect"
	// "strings"
)

const (
	envPrefix string = "PLUGIN_"
)

func main() {
	// fmt.Println("this is a test")
	// env := getEnviron()
	vars := env.Extract(envPrefix)
	fmt.Println(vars)

	cfg := &config{}
	err := env.Parse(vars, cfg)
	if err != nil {
		fmt.Printf("error: %+v\n", err)
		return
	}
	fmt.Printf("parsed: %+v\n", cfg)
	// inspect(reflect.ValueOf(cfg), 0)
}

// func getEnviron() map[string]string {
// 	env := make(map[string]string)
// 	for _, e := range os.Environ() {
// 		// fmt.Println(e)
// 		if strings.HasPrefix(e, envPrefix) {
// 			// fmt.Println(e)
// 			// split the key/value, and save the key _without_ the PLUGIN_ prefix.  (normalize to lower case?)
// 			keyValue := strings.SplitN(e, "=", 2)
// 			key := strings.ToLower(strings.TrimPrefix(keyValue[0], envPrefix))
// 			// fmt.Printf("%q: %q\n", key, keyValue[1])
// 			env[key] = keyValue[1]
// 		}
// 	}
// 	return env
// }
//
// func inspect(val reflect.Value, depth int) {
// 	if depth > 10 {
// 		fmt.Println("TOO DEEP!")
// 		return
// 	}
// 	pad := strings.Repeat("  ", depth)
// 	// val := reflect.ValueOf(i)
// 	kind := val.Kind()
// 	fmt.Printf("%s%s (settable: %v)\n", pad, kind, val.CanSet())
//
// 	switch kind {
//
// 	case reflect.Ptr:
// 		inspect(val.Elem(), depth+1)
//
// 	case reflect.Struct:
// 		structT := val.Type()
// 		fields := val.NumField()
// 		// fmt.Printf("%s-fields: %d\n", pad, fields)
// 		for fi := 0; fi < fields; fi++ {
// 			sf := structT.Field(fi)
// 			flag, ok := sf.Tag.Lookup("flag")
// 			// fmt.Printf("%s-%s (flag: %q %v, pkg: %q, anon: %v, tag: %+v)\n", pad, sf.Name, flag, ok, sf.PkgPath, sf.Anonymous, sf.Tag)
// 			fmt.Printf("%s-%s %q %v\n", pad, sf.Name, flag, ok)
// 			inspect(val.Field(fi), depth+1)
// 		}
//
// 	case reflect.Slice:
// 		sliceT := val.Type()
// 		fmt.Printf("%s %s\n", pad, sliceT.Elem().Kind())
// 		// inspect(sliceT.Elem(), depth+1)
// 	}
// }

type config struct {
	Bool   bool   `flag:"--bool"`
	Int    int    `flag:"--int"`
	Int8   int8   `flag:"--int8"`
	Int16  int16  `flag:"--int16"`
	Int32  int32  `flag:"--int32"`
	Int64  int64  `flag:"--int64"`
	Uint   uint   `flag:"--uint"`
	Uint8  uint8  `flag:"--uint8"`
	Uint16 uint16 `flag:"--uint16"`
	Uint32 uint32 `flag:"--uint32"`
	Uint64 uint64 `flag:"--uint64"`

	String string `flag:"--string"`

	Map map[string]string
}

func (c *config) foo() {

}

func (c *config) Bar() {

}
